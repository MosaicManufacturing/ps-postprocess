package msf

import (
    "../gcode"
    "strconv"
    "strings"
)

type Transition struct {
    Layer int
    From int
    To int
    TransitionLength float32 // actual transition length as specified by user
    PurgeLength float32 // amount of filament to extrude
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
    layerThicknesses []float32
    layerTopZs []float32
    transitionsByLayer map[int][]Transition

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

func preflight(inpath string, palette *Palette) (msfPreflight, error) {
    results := msfPreflight{
        drivesUsed:        make([]bool, palette.GetInputCount()),
        pingStarts:        make([]float32, 0),
        boundingBox:       gcode.NewBoundingBox(),
        towerBoundingBox:  gcode.NewBoundingBox(),
        printSummaryStart: -1,
        totalLayers:       -1,
        transitionsByLayer: make(map[int][]Transition),
    }

    // initialize state
    state := NewState(palette)
    // account for a firmware purge (not part of G-code) once
    state.E.TotalExtrusion += palette.FirmwarePurge
    // prepare to collect lookahead positions
    transitionNextPosition := sideTransitionLookahead{}

    pingExtrusionMM := palette.GetPingExtrusion()

    transitionCount := 0
    lastTransitionSpliceLength := float32(0)

    err := gcode.ReadByLine(inpath, func(line gcode.Command, lineNumber int) error {
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
                    // had X/Y (and maybe Z) movement and now we're extruding
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
                    if state.E.TotalExtrusion >= state.CurrentPingStart + pingExtrusionMM {
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
        } else if len(line.Command) > 1 && line.Command[0] == 'T' {
            tool, err := strconv.ParseInt(line.Command[1:], 10, 32)
            if err != nil {
                return err
            }
            if state.PastStartSequence {
                if state.FirstToolChange {
                    state.FirstToolChange = false
                } else {
                    transitionLength := palette.GetTransitionLength(int(tool), state.CurrentTool)
                    spliceOffset := transitionLength * (palette.TransitionTarget / 100)
                    purgeLength := transitionLength
                    spliceLength := state.E.TotalExtrusion + (transitionLength * spliceOffset)
                    deltaE := spliceLength - lastTransitionSpliceLength
                    minSpliceLength := MinSpliceLength
                    if transitionCount == 0 {
                        minSpliceLength = palette.GetFirstSpliceMinLength()
                    }
                    if deltaE < minSpliceLength {
                        extra := minSpliceLength - deltaE
                        purgeLength += extra
                        spliceLength += extra
                    }
                    tInfo := Transition{
                        Layer:            results.totalLayers,
                        From:             state.CurrentTool,
                        To:               int(tool),
                        TransitionLength: transitionLength,
                        PurgeLength:      purgeLength,
                    }
                    if _, ok := results.transitionsByLayer[results.totalLayers]; ok {
                        results.transitionsByLayer[results.totalLayers] = append(results.transitionsByLayer[results.totalLayers], tInfo)
                    } else {
                        results.transitionsByLayer[results.totalLayers] = []Transition{tInfo}
                    }
                    transitionCount++
                    lastTransitionSpliceLength = spliceLength - purgeLength // we haven't generated the purges yet
                    state.CurrentTool = int(tool)
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
        } else if palette.TransitionMethod == CustomTower &&
            strings.HasPrefix(line.Raw, ";Z:") {
            if topZ, err := strconv.ParseFloat(line.Raw[3:], 32); err == nil {
                results.layerTopZs = append(results.layerTopZs, float32(topZ))
            }
        } else if palette.TransitionMethod == CustomTower &&
            strings.HasPrefix(line.Raw, ";HEIGHT:") {
            if thickness, err := strconv.ParseFloat(line.Raw[8:], 32); err == nil {
                results.layerThicknesses = append(results.layerThicknesses, float32(thickness))
            }
        } else if palette.TransitionMethod == TransitionTower &&
            strings.HasPrefix(line.Comment, "TYPE:") {
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
        }

        return nil
    })
    if palette.TransitionMethod == SideTransitions && state.CurrentlyTransitioning {
        results.transitionNextPositions = append(results.transitionNextPositions, transitionNextPosition)
    }
    if results.boundingBox.Min[2] > 0 {
        results.boundingBox.Min[2] = 0
    }
    return results, err
}