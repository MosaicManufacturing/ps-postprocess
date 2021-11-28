package ptp

import (
    "errors"
    "log"
    "math"
    "mosaicmfg.com/ps-postprocess/gcode"
    "strconv"
    "strings"
)

func parseToolColors(serialized string) ([][3]float32, error) {
    perToolVals := strings.Split(serialized, "|")
    toolColors := make([][3]float32, 0, len(perToolVals))
    for _, colors := range perToolVals {
        rgbParts := strings.Split(colors, ",")
        if len(rgbParts) != 3 {
            return nil, errors.New("expected 3 components for RGB value")
        }
        r, rErr := strconv.ParseFloat(rgbParts[0], 32)
        if rErr != nil {
            return nil, rErr
        }
        g, gErr := strconv.ParseFloat(rgbParts[1], 32)
        if gErr != nil {
            return nil, gErr
        }
        b, bErr := strconv.ParseFloat(rgbParts[2], 32)
        if bErr != nil {
            return nil, bErr
        }
        thisToolColors := [3]float32{float32(r), float32(g), float32(b)}
        toolColors = append(toolColors, thisToolColors)
    }
    return toolColors, nil
}

func convertPathType(hint string) PathType {
    switch hint {
    case "Perimeter":
        return PathTypeInnerPerimeter
    case "External perimeter":
        fallthrough
    case "Overhang perimeter":
        return PathTypeOuterPerimeter
    case "Internal infill":
        return PathTypeInfill
    case "Solid infill":
        fallthrough
    case "Top solid infill":
        return PathTypeSolidLayer
    case "Bridge infill":
        return PathTypeBridge
    case "Ironing":
        return PathTypeIroning
    case "Gap fill":
        return PathTypeGapFill
    case "Skirt":
        fallthrough
    case "Skirt/Brim":
        return PathTypeBrim
    case "Support material":
        return PathTypeSupport
    case "Support material interface":
        return PathTypeSupportInterface
    case "Wipe tower":
        return PathTypeTransition
    case "Side transition":
        return PathTypeTransition
    case "Custom":
        return PathTypeStartSequence
    }
    return PathTypeUnknown
}

type ptpPreflight struct {
    minFeedrate float32
    maxFeedrate float32
    minTemperature float32
    maxTemperature float32
    minLayerHeight float32
    maxLayerHeight float32
}

func toolpathPreflight(inpath string) (ptpPreflight, error) {
    minFeedrate := float32(math.Inf(1)); maxFeedrate := float32(math.Inf(-1))
    minTemperature := float32(math.Inf(1)); maxTemperature := float32(math.Inf(-1))
    minLayerHeight := float32(math.Inf(1)); maxLayerHeight := float32(math.Inf(-1))
    currentFeedrate := float32(0)

    err := gcode.ReadByLine(inpath, func(line gcode.Command, _ int) error {
        if line.IsLinearMove() {
            // feedrates
            if f, ok := line.Params["f"]; ok {
                currentFeedrate = f
            }
            if _, ok := line.Params["e"]; ok {
                hasMovement := false
                if _, ok := line.Params["x"]; ok {
                    hasMovement = true
                }
                if _, ok := line.Params["y"]; ok {
                    hasMovement = true
                }
                if hasMovement {
                    if currentFeedrate < minFeedrate {
                        minFeedrate = currentFeedrate
                    }
                    if currentFeedrate > maxFeedrate {
                        maxFeedrate = currentFeedrate
                    }
                }
            }
        } else if line.Command == "M104" || line.Command == "M109" {
            // temperatures
            if temp, ok := line.Params["s"]; ok {
                if temp < minTemperature {
                    minTemperature = temp
                }
                if temp > maxTemperature {
                    maxTemperature = temp
                }
            }
        } else if line.Comment != "" &&  strings.HasPrefix(line.Comment, "HEIGHT:") {
            // layer heights
            height, err := strconv.ParseFloat(line.Comment[7:], 64)
            if err != nil {
                return err
            }
            height32 := roundZ(float32(height))
            if height32 < minLayerHeight {
                minLayerHeight = height32
            }
            if height32 > maxLayerHeight {
                maxLayerHeight = height32
            }
        }
        return nil
    })
    if err != nil {
        return ptpPreflight{}, err
    }
    results := ptpPreflight{
        minFeedrate:    minFeedrate,
        maxFeedrate:    maxFeedrate,
        minTemperature: minTemperature,
        maxTemperature: maxTemperature,
        minLayerHeight: minLayerHeight,
        maxLayerHeight: maxLayerHeight,
    }
    return results, err
}

func GenerateToolpath(argv []string) {
    argc := len(argv)

    if argc != 4 {
        log.Fatalln("expected 4 command-line arguments")
    }
    inpath := argv[0]
    outpath := argv[1]
    brimIsSkirt := argv[2] == "true"
    toolColors, err := parseToolColors(argv[3])
    if err != nil {
        log.Fatalln(err)
    }
    preflight, err := toolpathPreflight(inpath)
    if err != nil {
        log.Fatalln(err)
    }

    writer := NewWriter(outpath, brimIsSkirt, toolColors)
    writer.SetFeedrateBounds(preflight.minFeedrate, preflight.maxFeedrate)
    writer.SetTemperatureBounds(preflight.minTemperature, preflight.maxTemperature)
    writer.SetLayerHeightBounds(preflight.minLayerHeight, preflight.maxLayerHeight)
    if err := writer.Initialize(); err != nil {
        log.Fatalln(err)
    }

    currentE := float32(0)
    relativeE := false
    err = gcode.ReadByLine(inpath, func(line gcode.Command, _ int) error {
        if setExtrusionMode, relative := line.IsSetExtrusionMode(); setExtrusionMode {
            relativeE = relative
            currentE = 0
        } else if line.IsSetPosition() {
            if e, ok := line.Params["e"]; ok {
                currentE = e
            }
        } else if line.IsLinearMove() {
            isVisibleMove := false // either print line or travel line
            isPrintMove := false // specifically print line
            x, y, z := writer.GetCurrentPosition()
            if lineX, ok := line.Params["x"]; ok {
                x = lineX
                isVisibleMove = true
            }
            if lineY, ok := line.Params["y"]; ok {
                y = lineY
                isVisibleMove = true
            }
            if lineZ, ok := line.Params["z"]; ok {
                z = lineZ
                isVisibleMove = true
            }
            if e, ok := line.Params["e"]; ok {
                eIncreased := e > currentE
                eDecreased := e < currentE
                if relativeE {
                    eIncreased = e > 0
                    eDecreased = e < 0
                }
                if eIncreased {
                    if isVisibleMove {
                        isPrintMove = true
                    } else {
                        writer.AddRestart()
                    }
                } else if eDecreased {
                    // add retract point regardless of there being X/Y/Z movement as well
                    writer.AddRetract()
                }
                currentE = e
            }
            if f, ok := line.Params["f"]; ok {
                writer.SetFeedrate(f)
            }
            if isVisibleMove {
                if isPrintMove {
                    writer.AddXYZPrintLineTo(x, y, z)
                } else {
                    writer.AddXYZTravelTo(x, y, z)
                }
            }
        } else if line.Command == "M106" {
            if pwm, ok := line.Params["s"]; ok {
                writer.SetFanSpeed(int(pwm))
            }
        } else if line.Command == "M107" {
            writer.SetFanSpeed(0)
        } else if line.Command == "M104" || line.Command == "M109" {
            if temp, ok := line.Params["s"]; ok {
                writer.SetTemperature(temp)
            }
        } else if len(line.Command) > 1 && line.Command[0] == 'T' {
            tool, err := strconv.ParseInt(line.Command[1:], 10, 32)
            if err != nil {
                return err
            }
            writer.SetTool(int(tool))
        } else if line.Command == "M135" {
            if t, ok := line.Params["t"]; ok {
                writer.SetTool(int(t))
            }
        } else if line.Comment != "" {
            if strings.HasPrefix(line.Comment, "TYPE:") {
                // path type hints
                pathType := convertPathType(line.Comment[5:])
                writer.SetPathType(pathType)
            } else if strings.HasPrefix(line.Comment, "WIDTH:") {
                // extrusion width hints
                width, err := strconv.ParseFloat(line.Comment[6:], 32)
                if err != nil {
                    return err
                }
                writer.SetExtrusionWidth(float32(width))
            } else if strings.HasPrefix(line.Comment, "HEIGHT:") {
                // layer height hints
                height, err := strconv.ParseFloat(line.Comment[7:], 32)
                if err != nil {
                    return err
                }
                writer.SetLayerHeight(roundZ(float32(height)))
            } else if strings.HasPrefix(line.Comment, "Printing with input ") {
                tool, err := strconv.ParseInt(line.Comment[20:], 10, 32)
                if err != nil {
                    return err
                }
                writer.SetTool(int(tool))
            }
        }
        return nil
    })
    if err != nil {
        log.Fatalln(err)
    }
    if err := writer.Finalize(); err != nil {
        log.Fatalln(err)
    }
}
