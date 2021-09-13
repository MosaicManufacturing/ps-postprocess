package msf

import (
    "../gcode"
    "bufio"
    "errors"
    "fmt"
    "io/ioutil"
    "log"
    "os"
    "strconv"
    "strings"
)

// explaining outpath and msfpath:
// - P1:            outpath == *.msf.gcode,  msfpath == *.msf
// - P2 accessory:  outpath == *.maf.gcode,  msfpath == *.maf
// - P2 connected:  outpath == *.mcf.gcode,  [no msfpath]
// - P3 accessory:  outpath == *.gcode,      msfpath == *.json
// - P3 connected:  outpath == *.gcode,      msfpath == *.json

func paletteOutput(inpath, outpath, msfpath string, palette *Palette, preflight *msfPreflight) (err error) {
    // todo: P2 will need a temp file so the final MSF can be prepended
    outfile, createErr := os.Create(outpath)
    if createErr != nil {
        err = createErr
        return
    }
    defer func() {
        if closeErr := outfile.Close(); closeErr != nil {
            err = closeErr
        }
    }()
    writer := bufio.NewWriter(outfile)
    defer func() {
        if flushErr := writer.Flush(); flushErr != nil {
            err = flushErr
        }
    }()
    msfOut := NewMSF(palette)

    // initialize state
    state := NewState(palette)
    state.MSF = &msfOut
    state.TowerBoundingBox = preflight.towerBoundingBox
    // account for a firmware purge (not part of G-code) once
    state.E.TotalExtrusion += palette.FirmwarePurge

    pingExtrusionMM := palette.GetPingExtrusion()

    if len(preflight.pingStarts) > 0 {
        state.NextPingStart = preflight.pingStarts[0]
    } else {
        state.NextPingStart = posInf
    }

    err = gcode.ReadByLine(inpath, func(line gcode.Command) error {
        state.E.TrackInstruction(line)
        state.XYZF.TrackInstruction(line)
        if line.IsLinearMove() {
            if err := writeLine(writer, line.Raw); err != nil {
                return err
            }
            if state.OnWipeTower {
                // check for ping actions
                if state.CurrentlyPinging {
                    // currentlyPinging == true implies accessory mode
                    if state.E.TotalExtrusion >= state.LastPingStart + pingExtrusionMM {
                        // finish the accessory ping sequence
                        comment := fmt.Sprintf("; Ping %d pause 2", len(msfOut.PingList) + 1)
                        if err := writeLine(writer, comment); err != nil {
                            return err
                        }
                        if useRetract, retract := getPingRetract(palette); useRetract {
                            if err := writeLine(writer, retract); err != nil {
                                return err
                            }
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
                        if useRestart, restart := getPingRestart(palette); useRestart {
                            if err := writeLine(writer, restart); err != nil {
                                return err
                            }
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
                        if err := writeLine(writer, "G4 P0"); err != nil {
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
                        if useRetract, retract := getPingRetract(palette); useRetract {
                            if err := writeLine(writer, retract); err != nil {
                                return err
                            }
                        }
                        pauseSequence := getTowerPause(Ping1PauseLength, &state)
                        if err := writeLines(writer, pauseSequence); err != nil {
                            return err
                        }
                        state.LastPingStart = state.E.TotalExtrusion
                        if useRestart, restart := getPingRestart(palette); useRestart {
                            if err := writeLine(writer, restart); err != nil {
                                return err
                            }
                        }
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
                state.FirstToolChange = false
                if err := writeLine(writer, fmt.Sprintf("T%d ; change extruder", palette.PrintExtruder)); err != nil {
                    return err
                }
            } else {
                currentTransitionLength := palette.TransitionLengths[tool][state.CurrentTool]
                spliceOffset := currentTransitionLength * (palette.TransitionTarget / 100)
                if err := msfOut.AddSplice(state.CurrentTool, state.E.TotalExtrusion + spliceOffset); err != nil {
                    return err
                }
                state.CurrentTool = int(tool)
                state.CurrentlyTransitioning = true
                if palette.TransitionMethod == SideTransitions {
                    transition := sideTransition(currentTransitionLength, &state)
                    if err := writeLines(writer, transition); err != nil {
                        return err
                    }
                    state.CurrentlyTransitioning = false
                }
            }
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
    err = msfOut.AddLastSplice(state.CurrentTool, state.E.TotalExtrusion)
    if err != nil {
        return err
    }
    msfStr, err := msfOut.CreateMSF()
    if err != nil {
        return err
    }
    return ioutil.WriteFile(msfpath, []byte(msfStr), 0644)
}

func ConvertForPalette(argv []string) {
    argc := len(argv)

    if argc < 4 {
        log.Fatalln("expected 4 command-line arguments")
    }
    inpath := argv[0] // unmodified G-code file
    outpath := argv[1] // modified G-code file
    msfpath := argv[2] // supplementary MSF file, if applicable
    palettepath := argv[3] // serialized Palette data

    palette, err := LoadFromFile(palettepath)
    if err != nil {
        fmt.Print(err)
        log.Fatalln(err)
    }

    // preflight: run through the G-code once to determine all necessary
    // information for performing modifications

    // - drives used
    // - splice lengths -- check early if any splices will be too short
    // - number of pings
    // - bounding box
    preflightResults, err := preflight(inpath, &palette)
    if err != nil {
        log.Fatalln(err)
    }

    // output: run through the G-code once and apply modifications
    // using information determined in preflight

    // - start of print O commands
    // - add initial toolchange to Palette extruder
    // - remove toolchange commands
    // - accessory pings (two pauses with precise-ish amount of E between them)
    // - connected pings
    // - print summary in footer
    err = paletteOutput(inpath, outpath, msfpath, &palette, &preflightResults)
    if err != nil {
        log.Fatalln(err)
    }
}
