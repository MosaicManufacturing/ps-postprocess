package ptp

import (
	"log"
	"mosaicmfg.com/ps-postprocess/gcode"
	"strconv"
	"strings"
)

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
			isPrintMove := false   // specifically print line
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
					// don't add retract point until wipe sequence, if any, is complete
					if !writer.state.inWipe {
						// add retract point regardless of there being X/Y/Z movement as well
						writer.AddRetract()
					}
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
			if line.Comment == "WIPE_START" {
				writer.state.inWipe = true
			} else if line.Comment == "WIPE_END" {
				// retract points were not added during the wipe sequence
				if writer.state.inWipe {
					// add retract point regardless of there being X/Y/Z movement as well
					writer.AddRetract()
				}
				writer.state.inWipe = false
			} else if strings.HasPrefix(line.Comment, "TYPE:") {
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
