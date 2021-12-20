package msf

import (
    "fmt"
    "mosaicmfg.com/ps-postprocess/gcode"
    "strconv"
    "strings"
)

// round layer thicknesses and top Zs to this many decimal places
const maxZPrecision = 5

type Transition struct {
    Layer int
    From int
    To int
    TotalExtrusion float32 // total non-tower extrusion at start of transition
    TransitionLength float32 // actual transition length as specified by user
    PurgeLength float32 // amount of filament to extrude
    UsableInfill float32 // subtract this amount from the splice length
}

func (t Transition) String() string {
    return fmt.Sprintf(
        "Layer = %d, From = %d, To = %d, TotalExtrusion = %f, TransitionLength = %f, PurgeLength = %f, UsableInfill = %f",
        t.Layer,
        t.From,
        t.To,
        t.TotalExtrusion,
        t.TransitionLength,
        t.PurgeLength,
        t.UsableInfill,
    )
}

type msfPreflight struct {
    // always used
    drivesUsed []bool
    pingStarts []float32
    // used for print summary
    printSummaryStart int // line number before which to output our own print summary
    boundingBox gcode.BoundingBox
    // used for PrusaSlicer-generated towers
    towerBoundingBox gcode.BoundingBox

    // used for postprocess-generated towers
    layerTopZs []float32 // printing height of each layer (i.e. the Z value of the top of these paths)
    layerThicknesses []float32 // thickness of each layer in mm (i.e. layerTopZs[n] - layerTopZs[n-1])
    layerObjectStarts []int // number of "printing object" comments per layer
    layerObjectEnds []int // number of "stop printing object" comments per layer
    transitionsByLayer map[int][]Transition // array of Transition per layer
    transitions []Transition // same data as transitionsByLayer but flattened into 1D

    // used for side transition custom scripts
    transitionNextPositions []sideTransitionLookahead
    timeEstimate float32 // seconds
    totalLayers int
}

func (mp *msfPreflight) totalDrivesUsed() int {
    total := 0
    for _, used := range mp.drivesUsed {
        if used {
            total++
        }
    }
    return total
}

type sideTransitionLookahead struct {
    X float32
    Y float32
    Z float32
    Moved bool
}

func _preflight(readerFn func(callback gcode.LineCallback) error, palette *Palette) (msfPreflight, error) {
    results := msfPreflight{
        drivesUsed:        make([]bool, palette.GetInputCount()),
        pingStarts:        make([]float32, 0),
        boundingBox:       gcode.NewBoundingBox(),
        towerBoundingBox:  gcode.NewBoundingBox(),
        printSummaryStart: -1,
        totalLayers:       -1,
        transitionsByLayer: make(map[int][]Transition),
        transitions:        make([]Transition, 0),
    }

    // initialize state
    state := NewState(palette)
    // account for a firmware purge (not part of G-code) once
    state.E.TotalExtrusion += palette.FirmwarePurge
    // prepare to collect lookahead positions
    transitionNextPosition := sideTransitionLookahead{}

    transitionCount := 0
    lastTransitionLayer := 0
    lastTransitionSpliceLength := float32(0)

    // calculate available infill per transition
    currentInfillStartE := float32(-1) // < 0 indicates not to use this value

    err := readerFn(func(line gcode.Command, lineNumber int) error {
        state.E.TrackInstruction(line)
        state.XYZF.TrackInstruction(line)
        if line.IsLinearMove() {
            if x, ok := line.Params["x"]; ok {
                results.boundingBox.ExpandX(x)
            }
            if y, ok := line.Params["y"]; ok {
                results.boundingBox.ExpandY(y)
            }
            if z, ok := line.Params["z"]; ok {
                results.boundingBox.ExpandZ(z)
            }
            if palette.TransitionMethod == SideTransitions && state.CurrentlyTransitioning {
                continueLookahead := true
                if transitionNextPosition.Moved {
                    // had X/Y (and maybe Z) movement, and now we're extruding
                    //  - commit the most recent XYZ values, but ignore the ones in this command
                    if _, ok := line.Params["e"]; ok {
                        continueLookahead = false
                        results.transitionNextPositions = append(results.transitionNextPositions, transitionNextPosition)
                        transitionNextPosition = sideTransitionLookahead{}
                        state.CurrentlyTransitioning = false
                    }
                }
                if continueLookahead {
                    if x, ok := line.Params["x"]; ok {
                        transitionNextPosition.X = x
                        transitionNextPosition.Moved = true
                    }
                    if y, ok := line.Params["y"]; ok {
                        transitionNextPosition.Y = y
                        transitionNextPosition.Moved = true
                    }
                    if z, ok := line.Params["z"]; ok {
                        transitionNextPosition.Z = z
                    }
                }
            } else if state.OnWipeTower {
                if _, ok := line.Params["e"]; ok {
                    // extrusion on wipe tower -- update bounding box
                    if state.E.CurrentRetraction == 0 {
                        if x, ok := line.Params["x"]; ok {
                            results.towerBoundingBox.ExpandX(x)
                        }
                        if y, ok := line.Params["y"]; ok {
                            results.towerBoundingBox.ExpandY(y)
                        }
                    }
                }
                // check for ping actions
                if state.CurrentlyPinging {
                    // currentlyPinging == true implies accessory mode
                    if state.E.TotalExtrusion >= state.CurrentPingStart + state.PingExtrusion {
                        // commit to the accessory ping sequence
                        results.pingStarts = append(results.pingStarts, state.CurrentPingStart)
                        state.LastPingStart = state.CurrentPingStart
                        state.CurrentlyPinging = false
                    }
                } else if state.E.TotalExtrusion >= state.LastPingStart + PingMinSpacing {
                    // attempt to start a ping sequence
                    //  - connected pings: guaranteed to finish
                    //  - accessory pings: may be "cancelled" if near the end of the transition
                    if palette.ConnectedMode {
                        results.pingStarts = append(results.pingStarts, state.E.TotalExtrusion)
                        state.LastPingStart = state.E.TotalExtrusion
                    } else {
                        state.CurrentPingStart = state.E.TotalExtrusion
                        state.CurrentlyPinging = true
                    }
                }
            }
        } else if isToolChange, tool := line.IsToolChange(); isToolChange {
            if state.PastStartSequence {
                if state.FirstToolChange {
                    state.FirstToolChange = false
                    state.CurrentTool = tool
                    results.drivesUsed[state.CurrentTool] = true
                } else {
                    transitionLength := palette.GetTransitionLength(tool, state.CurrentTool)
                    spliceOffset := transitionLength * (palette.TransitionTarget / 100)
                    purgeLength := transitionLength
                    spliceLength := state.E.TotalExtrusion + spliceOffset
                    // start by subtracting usable infill from splice and purge length
                    usableInfill := float32(0)
                    if currentInfillStartE >= 0 && palette.InfillTransitioning {
                        usableInfill = state.E.TotalExtrusion - currentInfillStartE
                        if usableInfill < 0 {
                            usableInfill = 0
                        }
                        purgeLength -= usableInfill
                        spliceLength -= usableInfill
                    }
                    // safety check to ensure minimum piece lengths
                    deltaE := spliceLength - lastTransitionSpliceLength
                    // try to account for any sparse layers that will be added between the last
                    // dense layer and this one (note: sparse layer extrusion may be more than this)
                    if palette.TransitionMethod == CustomTower && results.totalLayers > lastTransitionLayer + 1 {
                        sparseLayers := results.totalLayers - (lastTransitionLayer + 1)
                        sparseLayerExtrusionEstimate := state.PingExtrusion * float32(sparseLayers)
                        deltaE += sparseLayerExtrusionEstimate
                    }
                    minSpliceLength := MinSpliceLength
                    if transitionCount == 0 {
                        minSpliceLength = palette.GetFirstSpliceMinLength()
                    }
                    if deltaE < minSpliceLength {
                        extra := minSpliceLength - deltaE
                        purgeLength += extra
                        spliceLength += extra
                        if palette.InfillTransitioning {
                            usableInfill -= extra
                            if usableInfill < 0 {
                                purgeLength += usableInfill
                                spliceLength += usableInfill
                                usableInfill = 0
                            }
                        }
                    }
                    tInfo := Transition{
                        Layer:            results.totalLayers,
                        From:             state.CurrentTool,
                        To:               tool,
                        TotalExtrusion:   state.E.TotalExtrusion,
                        TransitionLength: transitionLength,
                        PurgeLength:      purgeLength,
                        UsableInfill:     usableInfill,
                    }
                    results.transitions = append(results.transitions, tInfo)
                    if _, ok := results.transitionsByLayer[results.totalLayers]; ok {
                        results.transitionsByLayer[results.totalLayers] = append(results.transitionsByLayer[results.totalLayers], tInfo)
                    } else {
                        results.transitionsByLayer[results.totalLayers] = []Transition{tInfo}
                    }
                    transitionCount++
                    // we haven't actually inserted the purge paths yet, so state.E.TotalExtrusion is
                    // missing purgeLength mm -- account for this by subtracting from last splice length
                    lastTransitionSpliceLength = spliceLength - purgeLength
                    lastTransitionLayer = results.totalLayers
                    state.CurrentTool = tool
                    if palette.TransitionMethod != CustomTower {
                        state.CurrentlyTransitioning = true
                    }
                    results.drivesUsed[state.CurrentTool] = true
                }
            }
        } else if line.Raw == ";START_OF_PRINT" {
            state.PastStartSequence = true
        } else if line.Raw == ";LAYER_CHANGE" {
            results.totalLayers++
            results.layerTopZs = append(results.layerTopZs, 0)
            results.layerThicknesses = append(results.layerThicknesses, 0)
            results.layerObjectStarts = append(results.layerObjectStarts, 0)
            results.layerObjectEnds = append(results.layerObjectEnds, 0)
        } else if palette.TransitionMethod == CustomTower &&
            strings.HasPrefix(line.Raw, ";Z:") {
            if topZ, err := strconv.ParseFloat(line.Raw[3:], 64); err == nil {
                results.layerTopZs[results.totalLayers] = roundTo(float32(topZ), maxZPrecision)
            }
        } else if palette.TransitionMethod == CustomTower &&
            strings.HasPrefix(line.Raw, ";HEIGHT:") {
            if thickness, err := strconv.ParseFloat(line.Raw[8:], 64); err == nil {
                thickness32 := roundTo(float32(thickness), maxZPrecision)
                if thickness32 > results.layerThicknesses[results.totalLayers] {
                    results.layerThicknesses[results.totalLayers] = thickness32
                }
            }
        } else if (palette.TransitionMethod == TransitionTower || palette.InfillTransitioning) &&
            strings.HasPrefix(line.Comment, "TYPE:") {
            if line.Comment == "TYPE:Internal infill" {
                // changed to infill -- initialize accumulated value
                currentInfillStartE = state.E.TotalExtrusion
            } else {
                // changed to non-infill -- reset accumulated value
                currentInfillStartE = -1
            }
            startingWipeTower := line.Comment == "TYPE:Wipe tower"
            if !state.OnWipeTower && startingWipeTower {
                // start of the actual transition being printed
            } else if state.OnWipeTower && !startingWipeTower {
                // end of the actual transition being printed
                if state.CurrentlyTransitioning {
                    state.CurrentlyTransitioning = false
                    // if we're in the middle of an accessory ping, cancel it -- too late to finish
                    state.CurrentlyPinging = false
                }
            }
            state.OnWipeTower = startingWipeTower
        } else if results.timeEstimate == 0 &&
            strings.HasPrefix(line.Comment, "estimated printing time (normal mode) = ") {
            timeEstimate, err := gcode.ParseTimeString(line.Comment)
            if err != nil {
                return err
            }
            results.timeEstimate = timeEstimate
            results.printSummaryStart = lineNumber + 2
        } else if strings.HasPrefix(line.Comment, "stop printing object ") {
            results.layerObjectEnds[results.totalLayers]++
        } else if strings.HasPrefix(line.Comment, "printing object ") {
            results.layerObjectStarts[results.totalLayers]++
        }

        return nil
    })
    if err != nil {
        return results, err
    }
    results.totalLayers++ // switch from 0-indexing to a true count

    // invariant assertions
    if palette.TransitionMethod == CustomTower {
        if layerThicknesses := len(results.layerThicknesses); layerThicknesses != results.totalLayers {
            return results, fmt.Errorf("invariant violation: expected %d layerThicknesses, got %d", results.totalLayers, layerThicknesses)
        }
        if layerTopZs := len(results.layerTopZs); layerTopZs != results.totalLayers {
            return results, fmt.Errorf("invariant violation: expected %d layerTopZs, got %d", results.totalLayers, layerTopZs)
        }
        for i := 0; i < results.totalLayers; i++ {
            if results.layerThicknesses[i] == 0 {
                return results, fmt.Errorf("invariant violation: zero thickness at layer %d", i)
            }
            if results.layerTopZs[i] == 0 {
                return results, fmt.Errorf("invariant violation: zero height at layer %d", i)
            }
            if results.layerObjectStarts[i] == 0 {
                return results, fmt.Errorf("invariant violation: zero layer object starts at layer %d", i)
            }
            if results.layerObjectEnds[i] == 0 {
                return results, fmt.Errorf("invariant violation: zero layer object ends at layer %d", i)
            }
            if results.layerObjectStarts[i] != results.layerObjectEnds[i] {
                return results, fmt.Errorf("invariant violation: layer object count mismatch at layer %d", i)
            }
        }
    }
    if palette.ZOffset != 0 {
        for i := range results.layerTopZs {
            results.layerTopZs[i] += palette.ZOffset
        }
    }

    if palette.TransitionMethod == SideTransitions && state.CurrentlyTransitioning {
        results.transitionNextPositions = append(results.transitionNextPositions, transitionNextPosition)
    }
    if results.boundingBox.Min[2] > 0 {
        results.boundingBox.Min[2] = 0
    }
    return results, nil
}

func preflight(inpath string, palette *Palette) (msfPreflight, error) {
    readerFn := func(callback gcode.LineCallback) error {
        return gcode.ReadByLine(inpath, callback)
    }
    return _preflight(readerFn, palette)
}
