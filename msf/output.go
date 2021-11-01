package msf

import (
    "../gcode"
    "../sequences"
    "bufio"
    "errors"
    "fmt"
    "io/ioutil"
    "log"
    "os"
    "strings"
)

func paletteOutput(inpath, outpath, msfpath string, palette *Palette, preflight *msfPreflight, locals sequences.Locals) error {
    outfile, createErr := os.Create(outpath)
    if createErr != nil {
        return createErr
    }
    writer := bufio.NewWriter(outfile)
    msfOut := NewMSF(palette)

    // initialize state
    state := NewState(palette)
    state.MSF = &msfOut
    state.TowerBoundingBox = preflight.towerBoundingBox
    for _, position := range preflight.transitionNextPositions {
        state.TransitionNextPositions = append(state.TransitionNextPositions, [3]float32{
            position.X, position.Y, position.Z,
        })
    }
    locals.Global["totalTime"] = float64(preflight.timeEstimate)
    locals.Global["totalLayers"] = float64(preflight.totalLayers)

    state.Locals = locals
    // account for a firmware purge (not part of G-code) once
    state.E.TotalExtrusion += palette.FirmwarePurge
    state.TimeEstimate = preflight.timeEstimate

    if len(preflight.pingStarts) > 0 {
        state.NextPingStart = preflight.pingStarts[0]
    } else {
        state.NextPingStart = PingMinSpacing
    }

    if palette.TransitionMethod == CustomTower {
        tower, needsTower := GenerateTower(palette, preflight)
        if !needsTower {
            log.Fatalln("should not have generated a tower!")
        }
        state.Tower = &tower
    }

    didFinalSplice := false // used to prevent calling msfOut.AddLastSplice multiple times
    needsSparseLayers := false // used at end-of-layer to postpone sparse layer insertion slightly

    err := gcode.ReadByLine(inpath, func(line gcode.Command, lineNumber int) error {
        if lineNumber == preflight.printSummaryStart {
            if err := msfOut.AddLastSplice(state.CurrentTool, state.E.TotalExtrusion); err != nil {
                return err
            }
            didFinalSplice = true // make sure not to do this again at EOF
            // insert our (more accurate) print summary
            summary := getPrintSummary(&msfOut, state.TimeEstimate)
            if err := writeLines(writer, summary); err != nil {
                return err
            }
        }
        state.E.TrackInstruction(line)
        state.XYZF.TrackInstruction(line)
        state.Temperature.TrackInstruction(line)
        if line.IsLinearMove() {
            if err := writeLine(writer, line.Raw); err != nil {
                return err
            }
            if state.OnWipeTower {
                // check for ping actions
                if state.CurrentlyPinging {
                    // currentlyPinging == true implies accessory mode
                    if state.E.TotalExtrusion >= state.LastPingStart + state.PingExtrusion {
                        // finish the accessory ping sequence
                        comment := fmt.Sprintf("; Ping %d pause 2", len(msfOut.PingList) + 1)
                        if err := writeLine(writer, comment); err != nil {
                            return err
                        }
                        pauseSequence := getTowerPause(Ping2PauseLength, &state)
                        if err := writeLines(writer, pauseSequence); err != nil {
                            return err
                        }
                        actualPingExtrusion := state.E.TotalExtrusion - state.LastPingStart
                        msfOut.AddPingWithExtrusion(state.LastPingStart, actualPingExtrusion)
                        if len(msfOut.PingList) < len(preflight.pingStarts) {
                            state.NextPingStart = preflight.pingStarts[len(msfOut.PingList)]
                        } else {
                            state.NextPingStart = posInf
                        }
                        state.CurrentlyPinging = false
                    }
                } else if state.E.TotalExtrusion >= state.NextPingStart {
                    // attempt to start a ping sequence
                    //  - connected pings: guaranteed to finish
                    //  - accessory pings: may be "cancelled" if near the end of the transition
                    if palette.ConnectedMode {
                        comment := fmt.Sprintf("; Ping %d", len(msfOut.PingList) + 1)
                        if err := writeLine(writer, comment); err != nil {
                            return err
                        }
                        msfOut.AddPing(state.E.TotalExtrusion)
                        if err := writeLine(writer, palette.ClearBufferCommand); err != nil {
                            return err
                        }
                        pingLine := msfOut.GetConnectedPingLine()
                        if err := writeLines(writer, pingLine); err != nil {
                            return err
                        }
                        if len(msfOut.PingList) < len(preflight.pingStarts) {
                            state.NextPingStart = preflight.pingStarts[len(msfOut.PingList)]
                        } else {
                            state.NextPingStart = posInf
                        }
                        state.LastPingStart = state.E.TotalExtrusion
                    } else {
                        // start the accessory ping sequence
                        comment := fmt.Sprintf("; Ping %d pause 1", len(msfOut.PingList) + 1)
                        if err := writeLine(writer, comment); err != nil {
                            return err
                        }
                        pauseSequence := getTowerPause(Ping1PauseLength, &state)
                        if err := writeLines(writer, pauseSequence); err != nil {
                            return err
                        }
                        state.LastPingStart = state.E.TotalExtrusion
                        state.CurrentlyPinging = true
                    }
                }
            }
        } else if needsSparseLayers && len(line.Command) == 0 && strings.HasPrefix(line.Comment, "printing object") {
            if err := writeLine(writer, "; Sparse tower layer"); err != nil {
                return err
            }
            layerPaths, err := state.Tower.GetNextSegment(&state, false)
            if err != nil {
                return err
            }
            if err := writeLines(writer, layerPaths); err != nil {
                return err
            }
            // should we double up and print the next sparse layer now too?
            if !state.Tower.CurrentLayerIsDense() {
                if err := writeLine(writer, "; Sparse tower layer"); err != nil {
                    return err
                }
                // need to manually move up to next layer
                topZ := state.Tower.CurrentLayerTopZ()
                zLift := getZTravel(&state, topZ, "")
                if err := writeLines(writer, zLift); err != nil {
                    return err
                }
                layerPaths, err = state.Tower.GetNextSegment(&state, false)
                if err != nil {
                    return err
                }
                if err := writeLines(writer, layerPaths); err != nil {
                    return err
                }
            }
            needsSparseLayers = false
        } else if isToolChange, tool := line.IsToolChange(); isToolChange {
            if state.PastStartSequence {
                if state.FirstToolChange {
                    state.FirstToolChange = false
                    if err := writeLine(writer, fmt.Sprintf("T%d ; change extruder", palette.PrintExtruder)); err != nil {
                        return err
                    }
                    state.CurrentTool = tool
                    if err := writeLine(writer, fmt.Sprintf("; Printing with input %d", state.CurrentTool)); err != nil {
                        return err
                    }
                } else {
                    comment := fmt.Sprintf("; Printing with input %d", tool)
                    if err := writeLine(writer, comment); err != nil {
                        return err
                    }
                    if palette.TransitionMethod == CustomTower {
                        if err := writeLine(writer, "; Dense tower segment"); err != nil {
                            return err
                        }
                        currentTransition := state.Tower.GetCurrentTransitionInfo()
                        spliceOffset := currentTransition.TransitionLength * (palette.TransitionTarget / 100)
                        spliceLength := state.E.TotalExtrusion + spliceOffset - currentTransition.UsableInfill
                        if len(msfOut.SpliceList) == 0 {
                            spliceLength += state.Tower.BrimExtrusion
                        }
                        if err := msfOut.AddSplice(state.CurrentTool, spliceLength); err != nil {
                            return err
                        }
                        state.CurrentTool = tool
                        state.CurrentlyTransitioning = true
                        transition, err := state.Tower.GetNextSegment(&state, true)
                        if err != nil {
                            return err
                        }
                        if err := writeLines(writer, transition); err != nil {
                            return err
                        }
                        state.CurrentlyTransitioning = false
                    } else {
                        currentTransition := preflight.transitions[len(msfOut.SpliceList)]
                        currentPurgeLength := currentTransition.PurgeLength
                        spliceOffset := currentTransition.TransitionLength * (palette.TransitionTarget / 100)
                        spliceLength := state.E.TotalExtrusion + spliceOffset - currentTransition.UsableInfill
                        if palette.TransitionMethod == SideTransitions {
                            extra := msfOut.GetRequiredExtraSpliceLength(spliceLength)
                            if extra > 0 {
                                currentPurgeLength += extra
                                spliceLength += extra
                            }
                        }
                        if err := msfOut.AddSplice(state.CurrentTool, spliceLength); err != nil {
                            return err
                        }
                        state.CurrentTool = tool
                        state.CurrentlyTransitioning = true
                        if palette.TransitionMethod == SideTransitions {
                            transition, err := sideTransition(currentPurgeLength, &state)
                            if err != nil {
                                return err
                            }
                            if err := writeLines(writer, transition); err != nil {
                                return err
                            }
                            state.CurrentlyTransitioning = false
                        }
                    }
                }
            }
        } else if line.Raw == ";START_OF_PRINT" {
            state.PastStartSequence = true
            return writeLine(writer, line.Raw)
        } else if line.Raw == ";LAYER_CHANGE" {
            state.CurrentLayer++
            if palette.TransitionMethod == CustomTower && !state.Tower.IsComplete() {
                // check for sparse layer insertion
                if state.Tower.NeedsSparseLayers(state.CurrentLayer) {
                    needsSparseLayers = true
                }
            }
            return writeLine(writer, line.Raw)
        } else if palette.TransitionMethod == TransitionTower &&
            strings.HasPrefix(line.Comment, "TYPE:") {
            if err := writeLine(writer, line.Raw); err != nil {
                return err
            }
            startingWipeTower := line.Comment == "TYPE:Wipe tower"
            if !state.OnWipeTower && startingWipeTower {
                // start of the actual transition being printed
            } else if state.OnWipeTower && !startingWipeTower {
                // end of the actual transition being printed
                if state.CurrentlyPinging {
                    return errors.New("incomplete ping occurred")
                }
            }
            state.OnWipeTower = startingWipeTower
        } else {
            return writeLine(writer, line.Raw)
        }
        return nil
    })
    if err != nil {
        return err
    }
    if !didFinalSplice {
        if err := msfOut.AddLastSplice(state.CurrentTool, state.E.TotalExtrusion); err != nil {
            return err
        }
        didFinalSplice = true
    }
    if palette.Type == TypeP2 && palette.ConnectedMode {
        // .mcf.gcode -- append footer
        if err := writeLines(writer, msfOut.GetMSF2Footer()); err != nil {
            return err
        }
    }
    // finalize outfile now
    if err := writer.Flush(); err != nil {
        return err
    }
    if err := outfile.Close(); err != nil {
        return err
    }
    if palette.Type == TypeP2 && palette.ConnectedMode {
        // .mcf.gcode -- prepend header instead of writing to separate file
        header := msfOut.GetMSF2Header()
        if err := prependFile(outpath, header); err != nil {
            return err
        }
    } else {
        msfStr, err := msfOut.CreateMSF()
        if err != nil {
            return err
        }
        if err := ioutil.WriteFile(msfpath, []byte(msfStr), 0644); err != nil {
            return err
        }
    }

    return nil
}
