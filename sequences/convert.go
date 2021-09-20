package sequences

import (
    "../gcode"
    "../printerscript"
    "bufio"
    "fmt"
    "log"
    "os"
    "strings"
)

const EOL = "\r\n"

func convert(inpath, outpath string, scripts ParsedScripts, locals map[string]float64) error {
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
    locals["totalLayers"] = float64(preflightResults.totalLayers)
    locals["totalTime"] = float64(preflightResults.totalTime)

    // keep track of current state
    positionTracker := gcode.PositionTracker{}
    currentLayer := float64(0)
    currentPrintTemperature := float64(0)
    currentBedTemperature := float64(0)

    // todo: any way to cheaply calculate timeElapsed?

    err = gcode.ReadByLine(inpath, func(line gcode.Command, _ int) error {
        // update current position and/or temperature
        positionTracker.TrackInstruction(line)
        if line.Command == "M104" || line.Command == "M109" {
            if s, ok := line.Params["s"]; ok {
                currentPrintTemperature = float64(s)
            }
        } else if line.Command == "M140" || line.Command == "M190" {
            if s, ok := line.Params["s"]; ok {
                currentBedTemperature = float64(s)
            }
        }

        output := line.Raw
        if line.Raw == startPlaceholder {
            fmt.Println("start sequence found")
            opts := printerscript.InterpreterOptions{
                EOL:             EOL,
                TrailingNewline: false,
                Locals:          MergeLocals(locals, map[string]float64{
                    "layer": 0,
                    // todo: nextX/Y/Z
                    "currentPrintTemperature": 0,
                    "currentBedTemperature": currentBedTemperature,
                }),
            }
            result, err := printerscript.EvaluateTree(scripts.Start, opts)
            if err != nil {
                return err
            }
            output = result.Output
        } else if line.Raw == endPlaceholder {
            fmt.Println("end sequence found")
            opts := printerscript.InterpreterOptions{
                EOL:             EOL,
                TrailingNewline: false,
                Locals:          MergeLocals(locals, map[string]float64{
                    "layer": float64(preflightResults.totalLayers),
                    // todo: nextX/Y/Z
                    "currentPrintTemperature": currentPrintTemperature,
                    "currentBedTemperature": currentBedTemperature,
                }),
            }
            result, err := printerscript.EvaluateTree(scripts.End, opts)
            if err != nil {
                return err
            }
            output = result.Output
        } else if strings.HasPrefix(line.Raw, layerChangePrefix) {
            fmt.Println("layer change sequence found")
            layer, layerZ, err := parseLayerChangePlaceholder(line.Raw)
            fmt.Printf("layer = %d, layerZ = %f \n", layer, layerZ)
            if err != nil {
                return err
            }
            currentLayer = float64(layer)
            opts := printerscript.InterpreterOptions{
                EOL:             EOL,
                TrailingNewline: false,
                Locals:          MergeLocals(locals, map[string]float64{
                    "layer": currentLayer,
                    "currentX": float64(positionTracker.CurrentX),
                    "currentY": float64(positionTracker.CurrentY),
                    "currentZ": float64(positionTracker.CurrentZ),
                    // todo: nextX/Y
                    "nextZ": layerZ,
                    "currentPrintTemperature": currentPrintTemperature,
                    "currentBedTemperature": currentBedTemperature,
                }),
            }
            result, err := printerscript.EvaluateTree(scripts.LayerChange, opts)
            if err != nil {
                return err
            }
            output = result.Output
        } else if strings.HasPrefix(line.Raw, materialChangePrefix) {
            fmt.Println("material change sequence found")
            toTool, err := parseMaterialChangePlaceholder(line.Raw)
            if err != nil {
                return err
            }
            opts := printerscript.InterpreterOptions{
                EOL:             EOL,
                TrailingNewline: false,
                Locals:          MergeLocals(locals, map[string]float64{
                    "layer": currentLayer,
                    "currentX": float64(positionTracker.CurrentX),
                    "currentY": float64(positionTracker.CurrentY),
                    "currentZ": float64(positionTracker.CurrentZ),
                    // todo: nextX/Y/Z
                    "currentPrintTemperature": currentPrintTemperature,
                    "currentBedTemperature": currentBedTemperature,
                }),
            }
            result, err := printerscript.EvaluateTree(scripts.MaterialChange[toTool], opts)
            if err != nil {
                return err
            }
            output = result.Output
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

    if argc < 4 {
        log.Fatalln("expected 4 command-line arguments")
    }
    inpath := argv[0] // unmodified G-code file
    outpath := argv[1] // modified G-code file
    scriptspath := argv[2] // JSON-stringified scripts to swap in
    localspath := argv[3] // JSON-stringified locals

    scripts, err := LoadScripts(scriptspath)
    if err != nil {
        log.Fatalln(err)
    }

    // lex and parse scripts just once now, and re-use the parse trees when evaluating
    parsedScripts, err := scripts.Parse()
    if err != nil {
        log.Fatalln(err)
    }

    // load locals that are available in all scripts
    locals, err := LoadLocals(localspath)
    if err != nil {
        log.Fatalln(err)
    }

    err = convert(inpath, outpath, parsedScripts, locals)
    if err != nil {
        log.Fatalln(err)
    }
}
