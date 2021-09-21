package msf

import (
    "../gcode"
    "strconv"
    "strings"
)

type msfPreflight struct {
    drivesUsed []bool
    transitionStarts []float32
    pingStarts []float32
    boundingBox gcode.BoundingBox
    towerBoundingBox gcode.BoundingBox
    timeEstimate float32 // seconds
    printSummaryStart int // line number before which to output print summary
}

func preflight(inpath string, palette *Palette) (msfPreflight, error) {
    results := msfPreflight{
        drivesUsed:       make([]bool, palette.GetInputCount()),
        transitionStarts: make([]float32, 0),
        pingStarts:       make([]float32, 0),
        boundingBox:      gcode.NewBoundingBox(),
        towerBoundingBox: gcode.NewBoundingBox(),
        printSummaryStart: -1,
    }

    // initialize state
    state := NewState(palette)
    // account for a firmware purge (not part of G-code) once
    state.E.TotalExtrusion += palette.FirmwarePurge

    pingExtrusionMM := palette.GetPingExtrusion()

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
            if state.OnWipeTower {
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
            if state.FirstToolChange {
                // todo: only set state.FirstToolChange = false if we're past the start sequence
                //  (i.e. account for start sequences containing toolchanges)
                state.FirstToolChange = false
            } else {
                state.CurrentTool = int(tool)
                state.CurrentlyTransitioning = true
                results.drivesUsed[state.CurrentTool] = true
                results.transitionStarts = append(results.transitionStarts, state.E.TotalExtrusion)
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
    if results.boundingBox.Min[2] > 0 {
        results.boundingBox.Min[2] = 0
    }
    return results, err
}