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
    "strconv"
    "strings"
)

// explaining outpath and msfpath:
// - P1:            outpath == *.msf.gcode,  msfpath == *.msf
// - P2 accessory:  outpath == *.maf.gcode,  msfpath == *.maf
// - P2 connected:  outpath == *.mcf.gcode,  [no msfpath]
// - P3 accessory:  outpath == *.gcode,      msfpath == *.json
// - P3 connected:  outpath == *.gcode,      msfpath == *.json

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

    pingExtrusionMM := palette.GetPingExtrusion()

    if len(preflight.pingStarts) > 0 {
        state.NextPingStart = preflight.pingStarts[0]
    } else {
        state.NextPingStart = posInf
    }

    didFinalSplice := false

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
                    if state.E.TotalExtrusion >= state.LastPingStart + pingExtrusionMM {
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
                        pauseSequence := getTowerPause(Ping1PauseLength, &state)
                        if err := writeLines(writer, pauseSequence); err != nil {
                            return err
                        }
                        state.LastPingStart = state.E.TotalExtrusion
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
                    if err := writeLine(writer, fmt.Sprintf("T%d ; change extruder", palette.PrintExtruder)); err != nil {
                        return err
                    }
                } else {
                    comment := fmt.Sprintf("; Transitioning to T%d from T%d", tool, state.CurrentTool)
                    if err := writeLine(writer, comment); err != nil {
                        return err
                    }
                    currentTransitionLength := palette.TransitionLengths[tool][state.CurrentTool]
                    spliceOffset := currentTransitionLength * (palette.TransitionTarget / 100)
                    extra := float32(0) // todo: account for extending transition length to maintain piece lengths
                    spliceLength := state.E.TotalExtrusion + spliceOffset + extra
                    if err := msfOut.AddSplice(state.CurrentTool, spliceLength); err != nil {
                        return err
                    }
                    state.CurrentTool = int(tool)
                    state.CurrentlyTransitioning = true
                    if palette.TransitionMethod == SideTransitions {
                        transition, err := sideTransition(currentTransitionLength, &state)
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
        } else if line.Raw == ";START_OF_PRINT" {
            state.PastStartSequence = true
        } else if line.Raw == ";LAYER_CHANGE" {
            state.CurrentLayer++
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

    // MSF 3 preheat hints
    if palette.Type == TypeP3 {
        preheatHintsPath := inpath + ".preheat"
        firstTool := msfOut.SpliceList[0].Drive
        hints := PreheatHints{
            Extruder: palette.FirstLayerTemperatures[firstTool],
            Bed:      palette.FirstLayerBedTemperatures[firstTool],
            Chamber:  0,
        }
        if err := hints.Save(preheatHintsPath); err != nil {
            return err
        }
    }

    return nil
}

func ConvertForPalette(argv []string) {
    argc := len(argv)

    if argc < 6 {
        log.Fatalln("expected 6 command-line arguments")
    }
    inpath := argv[0] // unmodified G-code file
    outpath := argv[1] // modified G-code file
    msfpath := argv[2] // supplementary MSF file, if applicable
    palettepath := argv[3] // serialized Palette data
    localsPath := argv[4] // JSON-stringified locals
    perExtruderLocalsPath := argv[5] // JSON-stringified locals

    palette, err := LoadFromFile(palettepath)
    if err != nil {
        log.Fatalln(err)
    }

    locals := sequences.NewLocals()
    if err := locals.LoadGlobal(localsPath); err != nil {
        log.Fatalln(err)
    }
    if err := locals.LoadPerExtruder(perExtruderLocalsPath); err != nil {
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
    err = paletteOutput(inpath, outpath, msfpath, &palette, &preflightResults, locals)
    if err != nil {
        log.Fatalln(err)
    }
}
