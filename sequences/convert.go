package sequences

import (
    "../gcode"
    "../printerscript"
    "bufio"
    "log"
    "os"
    "strconv"
    "strings"
)

const EOL = "\r\n"

func convert(inpath, outpath string, scripts ParsedScripts, locals Locals) error {
    outfile, createErr := os.Create(outpath)
    if createErr != nil {
        return createErr
    }
    writer := bufio.NewWriter(outfile)

    // run through the file once for summary information
    preflightResults, err := preflight(inpath)
    if err != nil {
        return err
    }
    locals.Global["totalLayers"] = float64(preflightResults.totalLayers)
    locals.Global["totalTime"] = float64(preflightResults.totalTime)

    // keep track of current state
    positionTracker := gcode.PositionTracker{}
    temperatureTracker := gcode.TemperatureTracker{}
    currentTool := 0
    currentLayer := float64(0)
    nextLayerChangeIdx := 0
    nextMaterialChangeIdx := 0

    // todo: any way to cheaply calculate timeElapsed?

    err = gcode.ReadByLine(inpath, func(line gcode.Command, _ int) error {
        // todo: hackily replace "Sliced by PrusaSlicer" with "Sliced by Canvas" now

        // update current position and/or temperature
        positionTracker.TrackInstruction(line)
        temperatureTracker.TrackInstruction(line)
        if len(line.Command) > 1 && line.Command[0] == 'T' {
            tool, err := strconv.ParseInt(line.Command[1:], 10, 32)
            if err != nil {
                return err
            }
            currentTool = int(tool)
        }

        output := line.Raw
        if line.Raw == startPlaceholder {
            opts := printerscript.InterpreterOptions{
                EOL:             EOL,
                TrailingNewline: false,
                Locals:          locals.Prepare(currentTool, map[string]float64{
                    "layer": 0,
                    "nextX": preflightResults.startSequenceNextPos.nextX,
                    "nextY": preflightResults.startSequenceNextPos.nextY,
                    "nextZ": preflightResults.startSequenceNextPos.nextZ,
                    "currentPrintTemperature": 0,
                    "currentBedTemperature": float64(temperatureTracker.Bed),
                }),
            }
            result, err := printerscript.EvaluateTree(scripts.Start, opts)
            if err != nil {
                return err
            }
            output = result.Output
        } else if line.Raw == endPlaceholder {
            opts := printerscript.InterpreterOptions{
                EOL:             EOL,
                TrailingNewline: false,
                Locals:          locals.Prepare(currentTool, map[string]float64{
                    "layer": float64(preflightResults.totalLayers),
                    "currentPrintTemperature": float64(temperatureTracker.Extruder),
                    "currentBedTemperature": float64(temperatureTracker.Bed),
                }),
            }
            result, err := printerscript.EvaluateTree(scripts.End, opts)
            if err != nil {
                return err
            }
            output = result.Output
        } else if strings.HasPrefix(line.Raw, layerChangePrefix) {
            layer, layerZ, err := parseLayerChangePlaceholder(line.Raw)
            if err != nil {
                return err
            }
            currentLayer = float64(layer)
            opts := printerscript.InterpreterOptions{
                EOL:             EOL,
                TrailingNewline: false,
                Locals:          locals.Prepare(currentTool, map[string]float64{
                    "layer": currentLayer,
                    "currentX": float64(positionTracker.CurrentX),
                    "currentY": float64(positionTracker.CurrentY),
                    "currentZ": float64(positionTracker.CurrentZ),
                    "nextX": preflightResults.layerChangeNextPos[nextLayerChangeIdx].nextX,
                    "nextY": preflightResults.layerChangeNextPos[nextLayerChangeIdx].nextY,
                    "nextZ": layerZ,
                    "currentPrintTemperature": float64(temperatureTracker.Extruder),
                    "currentBedTemperature": float64(temperatureTracker.Bed),
                }),
            }
            result, err := printerscript.EvaluateTree(scripts.LayerChange, opts)
            if err != nil {
                return err
            }
            output = result.Output
            nextLayerChangeIdx++
        } else if strings.HasPrefix(line.Raw, materialChangePrefix) {
            toTool, err := parseMaterialChangePlaceholder(line.Raw)
            if err != nil {
                return err
            }
            opts := printerscript.InterpreterOptions{
                EOL:             EOL,
                TrailingNewline: false,
                Locals:          locals.Prepare(currentTool, map[string]float64{
                    "layer": currentLayer,
                    "currentX": float64(positionTracker.CurrentX),
                    "currentY": float64(positionTracker.CurrentY),
                    "currentZ": float64(positionTracker.CurrentZ),
                    "nextX": preflightResults.layerChangeNextPos[nextMaterialChangeIdx].nextX,
                    "nextY": preflightResults.layerChangeNextPos[nextMaterialChangeIdx].nextY,
                    "nextZ": preflightResults.layerChangeNextPos[nextMaterialChangeIdx].nextZ,
                    "currentPrintTemperature": float64(temperatureTracker.Extruder),
                    "currentBedTemperature": float64(temperatureTracker.Bed),
                }),
            }
            result, err := printerscript.EvaluateTree(scripts.MaterialChange[toTool], opts)
            if err != nil {
                return err
            }
            output = result.Output
            nextMaterialChangeIdx++
        }
        if _, err := writer.WriteString(output + EOL); err != nil {
            return err
        }
        return nil
    })
    if err != nil {
        return err
    }

    if err := writer.Flush(); err != nil {
        return err
    }
    if err := outfile.Close(); err != nil {
        return err
    }
    return nil
}

func ConvertSequences(argv []string) {
    argc := len(argv)

    if argc < 5 {
        log.Fatalln("expected 5 command-line arguments")
    }
    inPath := argv[0] // unmodified G-code file
    outPath := argv[1] // modified G-code file
    scriptsPath := argv[2] // JSON-stringified scripts to swap in
    localsPath := argv[3] // JSON-stringified locals
    perExtruderLocalsPath := argv[4] // JSON-stringified locals

    scripts, err := LoadScripts(scriptsPath)
    if err != nil {
        log.Fatalln(err)
    }

    // lex and parse scripts just once now, and re-use the parse trees when evaluating
    parsedScripts, err := scripts.Parse()
    if err != nil {
        log.Fatalln(err)
    }

    // load locals that are available in all scripts
    locals := NewLocals()
    if err := locals.LoadGlobal(localsPath); err != nil {
        log.Fatalln(err)
    }
    if err := locals.LoadPerExtruder(perExtruderLocalsPath); err != nil {
        log.Fatalln(err)
    }

    err = convert(inPath, outPath, parsedScripts, locals)
    if err != nil {
        log.Fatalln(err)
    }
}
