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

func convert(inpath, outpath string, scripts ParsedScripts) error {
    outfile, createErr := os.Create(outpath)
    if createErr != nil {
        return createErr
    }
    writer := bufio.NewWriter(outfile)

    err := gcode.ReadByLine(inpath, func(line gcode.Command, _ int) error {
        output := line.Raw
        if line.Raw == startPlaceholder {
            fmt.Println("start sequence found")
            opts := printerscript.InterpreterOptions{
                EOL:             EOL,
                TrailingNewline: false,
                Locals:          map[string]float64{
                    // todo: populate printer/style settings (constants)
                    // todo: populate current state values (dynamic)
                    "layer": 0,
                    // todo: nextX/Y/Z
                    "currentPrintTemperature": 0,
                },
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
                Locals:          map[string]float64{
                    // todo: populate printer/style settings (constants)
                    // todo: populate current state values (dynamic)
                    // todo: layer (equal to totalLayers)
                    // todo: nextX/Y/Z
                },
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
            opts := printerscript.InterpreterOptions{
                EOL:             EOL,
                TrailingNewline: false,
                Locals:          map[string]float64{
                    // todo: populate printer/style settings (constants)
                    // todo: populate current state values (dynamic)
                    "layer": float64(layer),
                    // todo: currentX/Y/Z
                    // todo: nextX/Y
                    "nextZ": layerZ,
                },
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
                Locals:          map[string]float64{
                    // todo: populate printer/style settings (constants)
                    // todo: populate current state values (dynamic)
                    // todo: layer
                    // todo: currentX/Y/Z
                    // todo: nextX/Y/Z
                },
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
        log.Fatalln(err)
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

    if argc < 3 {
        log.Fatalln("expected 3 command-line arguments")
    }
    inpath := argv[0] // unmodified G-code file
    outpath := argv[1] // modified G-code file
    scriptspath := argv[2] // JSON-stringified scripts to swap in

    scripts, err := LoadScripts(scriptspath)
    if err != nil {
        log.Fatalln(err)
    }

    // lex and parse scripts just once now, and re-use the parse trees when evaluating
    parsedScripts, err := scripts.Parse()
    if err != nil {
        log.Fatalln(err)
    }

    err = convert(inpath, outpath, parsedScripts)
    if err != nil {
        log.Fatalln(err)
    }
}
