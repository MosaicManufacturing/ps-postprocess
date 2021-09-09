package main

import (
    "./gcode"
    "./palette"
    "bufio"
    "errors"
    "fmt"
    "io/ioutil"
    "log"
    "math"
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

const EOL = "\r\n"

type bbox struct {
    min [3]float32
    max [3]float32
}

type msfPreflight struct {
    drivesUsed []bool
    transitionStarts []float32
    transitionEnds []float32
    pingStarts []float32
    boundingBox bbox
}

func palettePreflight(inpath string, pal *palette.Palette) (msfPreflight, error) {
    results := msfPreflight{
        drivesUsed:       make([]bool, pal.GetInputCount()),
        transitionStarts: make([]float32, 0),
        transitionEnds:   make([]float32, 0),
        pingStarts:       make([]float32, 0),
        boundingBox:      bbox{
            min: [3]float32{float32(math.Inf(1)), float32(math.Inf(1)), float32(math.Inf(1))},
            max: [3]float32{float32(math.Inf(-1)), float32(math.Inf(-1)), float32(math.Inf(-1))},
        },
    }

    firstToolChange := true // don't treat the first T command as a toolchange
    currentTool := 0
    currentX := float32(0)
    currentY := float32(0)
    currentZ := float32(0)
    eTracker := gcode.ExtrusionTracker{}

    currentlyTransitioning := false
    onWipeTower := false
    pingExtrusionMM := float32(palette.PingExtrusion)
    if pal.Type == palette.TypeP1 {
        pingExtrusionMM = palette.PingExtrusionCounts / pal.GetPulsesPerMM()
    }
    lastPingStart := float32(0)
    currentlyPinging := false
    currentPingStart := float32(0)

    err := gcode.ReadByLine(inpath, func(line gcode.Command) error {
        eTracker.TrackInstruction(line)
        if line.IsLinearMove() {
            if x, ok := line.Params["x"]; ok {
                currentX = x
                if x < results.boundingBox.min[0] { results.boundingBox.min[0] = x }
                if x > results.boundingBox.max[0] { results.boundingBox.max[0] = x }
            }
            if y, ok := line.Params["y"]; ok {
                currentY = y
                if y < results.boundingBox.min[1] { results.boundingBox.min[1] = y }
                if y > results.boundingBox.max[1] { results.boundingBox.max[1] = y }
            }
            if z, ok := line.Params["z"]; ok {
                currentZ = z
                if z < results.boundingBox.min[2] { results.boundingBox.min[2] = z }
                if z > results.boundingBox.max[2] { results.boundingBox.max[2] = z }
            }
            if onWipeTower {
                // check for ping actions
                if currentlyPinging {
                    if eTracker.TotalExtrusion >= currentPingStart + pingExtrusionMM {
                        // commit to the accessory ping sequence
                        results.pingStarts = append(results.pingStarts, currentPingStart)
                        lastPingStart = currentPingStart
                        currentlyPinging = false
                    }
                } else if eTracker.TotalExtrusion >= lastPingStart + palette.PingMinSpacing {
                    // attempt to start a ping sequence
                    //  - connected pings: guaranteed to finish
                    //  - accessory pings: may be "cancelled" if near the end of the transition
                    if pal.ConnectedMode {
                        results.pingStarts = append(results.pingStarts, eTracker.TotalExtrusion)
                        lastPingStart = eTracker.TotalExtrusion
                    } else {
                        currentPingStart = eTracker.TotalExtrusion
                        currentlyPinging = true
                    }
                }
            }
        } else if len(line.Command) > 1 && line.Command[0] == 'T' {
            tool, err := strconv.ParseInt(line.Command[1:], 10, 32)
            if err != nil {
                return err
            }
            if firstToolChange {
                firstToolChange = false
            } else {
                currentTool = int(tool)
                currentlyTransitioning = true
                results.drivesUsed[currentTool] = true
                results.transitionStarts = append(results.transitionStarts, eTracker.TotalExtrusion)
            }
        } else if strings.HasPrefix(line.Comment, "TYPE:") {
            startWipeTower := line.Comment == "TYPE:Wipe tower"
            if !onWipeTower && startWipeTower {
                // start of the actual transition being printed
                // todo: any logic needed?
            } else if onWipeTower && !startWipeTower {
                // end of the actual transition being printed
                if currentlyTransitioning {
                    results.transitionEnds = append(results.transitionEnds, eTracker.TotalExtrusion)
                    currentlyTransitioning = false
                    // if we're in the middle of an accessory ping, cancel it -- too late to finish
                    currentlyPinging = false
                }
            }
            onWipeTower = startWipeTower
        }

        return nil
    })
    if len(results.transitionStarts) != len(results.transitionEnds) {
        return results, errors.New("mismatch between transition starts and ends")
    }
    if results.boundingBox.min[2] > 0 {
        results.boundingBox.min[2] = 0
    }
    return results, err
}

func getAccessoryPause(durationMS int) string {
    str := ""
    for durationMS > 0 {
        if durationMS > 4000 {
            str += "G4 P4000" + EOL
            str += "G1" + EOL
            durationMS -= 4000
        } else {
            str += fmt.Sprintf("G4 P%d%s", durationMS, EOL)
            durationMS = 0
        }
    }
    return str
}

func writeLine(writer *bufio.Writer, line string) error {
    _, err := writer.WriteString(line + EOL)
    return err
}

func writeLines(writer *bufio.Writer, lines string) error {
    _, err := writer.WriteString(lines)
    return err
}

func paletteOutput(inpath, outpath, msfpath string, pal *palette.Palette, preflight *msfPreflight) (err error) {
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
    msf := palette.NewMSF(pal)

    firstToolChange := true // don't treat the first T command as a toolchange
    currentTool := 0
    eTracker := gcode.ExtrusionTracker{}

    currentlyTransitioning := false
    onWipeTower := false
    pingExtrusionMM := float32(palette.PingExtrusion)
    if pal.Type == palette.TypeP1 {
        pingExtrusionMM = palette.PingExtrusionCounts / pal.GetPulsesPerMM()
    }
    lastPingStart := float32(0)
    nextPingStart := float32(math.Inf(1))
    if len(preflight.pingStarts) > 0 {
        nextPingStart = preflight.pingStarts[0]
    }
    currentlyPinging := false

    err = gcode.ReadByLine(inpath, func(line gcode.Command) error {
        eTracker.TrackInstruction(line)
        if line.IsLinearMove() {
            if err := writeLine(writer, line.Raw); err != nil {
                return err
            }
            if onWipeTower {
                // check for ping actions
                if currentlyPinging {
                    if eTracker.TotalExtrusion >= lastPingStart + pingExtrusionMM {
                        // finish the accessory ping sequence
                        comment := fmt.Sprintf("; Ping %d pause 2", len(msf.PingList) + 1)
                        if err := writeLine(writer, comment); err != nil {
                            return err
                        }
                        pauseSequence := getAccessoryPause(palette.Ping2PauseLength)
                        if err := writeLines(writer, pauseSequence); err != nil {
                            return err
                        }
                        actualPingExtrusion := eTracker.TotalExtrusion - lastPingStart
                        msf.AddPingWithExtrusion(lastPingStart, actualPingExtrusion)
                        if len(msf.PingList) < len(preflight.pingStarts) {
                            nextPingStart = preflight.pingStarts[len(msf.PingList)]
                        } else {
                            nextPingStart = float32(math.Inf(1))
                        }
                        currentlyPinging = false
                    }
                } else if eTracker.TotalExtrusion >= nextPingStart {
                    // attempt to start a ping sequence
                    //  - connected pings: guaranteed to finish
                    //  - accessory pings: may be "cancelled" if near the end of the transition
                    if pal.ConnectedMode {
                        comment := fmt.Sprintf("; Ping %d", len(msf.PingList) + 1)
                        if err := writeLine(writer, comment); err != nil {
                            return err
                        }
                        pingLine := msf.GetConnectedPingLine()
                        if err := writeLines(writer, pingLine); err != nil {
                            return err
                        }
                        msf.AddPing(eTracker.TotalExtrusion)
                        if len(msf.PingList) < len(preflight.pingStarts) {
                            nextPingStart = preflight.pingStarts[len(msf.PingList)]
                        } else {
                            nextPingStart = float32(math.Inf(1))
                        }
                        lastPingStart = eTracker.TotalExtrusion
                    } else {
                        // start the accessory ping sequence
                        comment := fmt.Sprintf("; Ping %d pause 1", len(msf.PingList) + 1)
                        if err := writeLine(writer, comment); err != nil {
                            return err
                        }
                        pauseSequence := getAccessoryPause(palette.Ping1PauseLength)
                        if err := writeLines(writer, pauseSequence); err != nil {
                            return err
                        }
                        lastPingStart = eTracker.TotalExtrusion
                        currentlyPinging = true
                    }
                }
            }
        } else if len(line.Command) > 1 && line.Command[0] == 'T' {
            tool, err := strconv.ParseInt(line.Command[1:], 10, 32)
            if err != nil {
                return err
            }
            if firstToolChange {
                firstToolChange = false
                if err := writeLine(writer, fmt.Sprintf("T%d ; change extruder", pal.PrintExtruder)); err != nil {
                    return err
                }
            } else {
                currentTransitionIdx := len(msf.SpliceList)
                currentTransitionLength := preflight.transitionEnds[currentTransitionIdx] - preflight.transitionStarts[currentTransitionIdx]
                spliceOffset := currentTransitionLength * (pal.TransitionTarget / 100)
                if err := msf.AddSplice(currentTool, eTracker.TotalExtrusion + spliceOffset); err != nil {
                   return err
                }
                currentTool = int(tool)
                currentlyTransitioning = true
            }
        } else if strings.HasPrefix(line.Comment, "TYPE:") {
            if err := writeLine(writer, line.Raw); err != nil {
                return err
            }
            startWipeTower := line.Comment == "TYPE:Wipe tower"
            if !onWipeTower && startWipeTower {
                // start of the actual transition being printed
                // todo: any logic needed?
            } else if onWipeTower && !startWipeTower {
                // end of the actual transition being printed
                if currentlyPinging {
                    return errors.New("incomplete ping occurred")
                }
            }
            onWipeTower = startWipeTower
        } else {
            return writeLine(writer, line.Raw)
        }
        return nil
    })
    if err != nil {
        return err
    }
    err = msf.AddLastSplice(currentTool, eTracker.TotalExtrusion)
    if err != nil {
        return err
    }
    msfStr, err := msf.CreateMSF()
    if err != nil {
        return err
    }
    return ioutil.WriteFile(msfpath, []byte(msfStr), 0644)
}

func convertForPalette(argv []string) {
    argc := len(argv)

    if argc < 4 {
        log.Fatalln("expected 4 command-line arguments")
    }
    inpath := argv[0] // unmodified G-code file
    outpath := argv[1] // modified G-code file
    msfpath := argv[2] // supplementary MSF file, if applicable
    palettepath := argv[3] // serialized palette data

    pal, err := palette.LoadFromFile(palettepath)
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
    preflight, err := palettePreflight(inpath, &pal)
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
    err = paletteOutput(inpath, outpath, msfpath, &pal, &preflight)
    if err != nil {
        log.Fatalln(err)
    }
}
