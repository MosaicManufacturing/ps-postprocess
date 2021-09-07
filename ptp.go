package main

import (
    "./gcode"
    "./ptp"
    "errors"
    "log"
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

func convertPathType(hint string) ptp.PathType {
    switch hint {
    case "Perimeter":
        return ptp.PathTypeInnerPerimeter
    case "External perimeter":
        fallthrough
    case "Overhang perimeter":
        return ptp.PathTypeOuterPerimeter
    case "Internal infill":
        return ptp.PathTypeInfill
    case "Solid infill":
        fallthrough
    case "Top solid infill":
        fallthrough
    case "Ironing":
        return ptp.PathTypeSolidLayer
    case "Bridge infill":
        return ptp.PathTypeBridge
    case "Gap fill":
        return ptp.PathTypeGapFill
    case "Skirt":
        fallthrough
    case "Skirt/Brim":
        return ptp.PathTypeBrim
    case "Support material":
        return ptp.PathTypeSupport
    case "Support material interface":
        return ptp.PathTypeSupportInterface
    case "Wipe tower":
        return ptp.PathTypeTransition
    case "Custom":
        return ptp.PathTypeStartSequence
    }
    return ptp.PathTypeUnknown
}

func generateToolpath(argv []string) {
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
    writer := ptp.NewWriter(outpath, brimIsSkirt, toolColors)
    // TODO: set these bounds from actual information
    writer.SetFeedrateBounds(0, 9000)
    writer.SetTemperatureBounds(0, 300)
    writer.SetLayerHeightBounds(0.2, 0.3)
    if err := writer.Initialize(); err != nil {
        log.Fatalln(err)
    }

    currentE := float32(0)
    relativeE := false
    err = gcode.ReadByLine(inpath, func(line gcode.Command) {
        isSetExtrusionMode, isRelativeE := line.IsSetExtrusionMode()
        if isSetExtrusionMode {
            relativeE = isRelativeE
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
                log.Fatalln(err)
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
                    log.Fatalln(err)
                }
                writer.SetExtrusionWidth(float32(width))
            } else if strings.HasPrefix(line.Comment, "HEIGHT:") {
                // layer height hints
                height, err := strconv.ParseFloat(line.Comment[7:], 32)
                if err != nil {
                    log.Fatalln(err)
                }
                writer.SetLayerHeight(float32(height))
            }
        }
    })
    if err != nil {
        log.Fatalln(err)
    }
    if err := writer.Finalize(); err != nil {
        log.Fatalln(err)
    }
}