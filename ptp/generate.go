package ptp

import (
	"errors"
	"log"
	"regexp"
	"strconv"
	"strings"

	"mosaicmfg.com/ps-postprocess/gcode"
)

type generatorState struct {
	// used for all generation
	currentTool   int
	currentLayerZ float32
	currentE      float32
	relativeE     bool

	// only used for transition tower gradients
	extrusionSoFar   float32 // cumulative over the transition
	transitioning    bool    // when true, values below are enabled
	lastTool         int     // one behind currentTool
	purgeLength      float32 // constant for entire transition
	transitionLength float32 // constant for entire transition
	offset           float32 // constant for entire transition
	target           float32 // constant for entire transition

}

func getStartingGeneratorState() generatorState {
	return generatorState{
		currentTool:      0,
		currentLayerZ:    0,
		currentE:         0,
		relativeE:        false,
		transitioning:    false,
		extrusionSoFar:   0,
		lastTool:         0,
		purgeLength:      0,
		transitionLength: 0,
		offset:           0,
		target:           0,
	}
}

var rePtpTowerComment = regexp.MustCompile("\\(purge=(.*),transition=(.*),offset=(.*),target=(.*)\\)")

func parsePtpTowerComment(comment string) (error, float32, float32, float32, float32) {
	matches := rePtpTowerComment.FindStringSubmatch(comment)
	if len(matches) < 5 {
		return errors.New("failed to parse PTP comment"), 0, 0, 0, 0
	}
	purgeLength := float32(0)
	transitionLength := float32(0)
	offset := float32(0)
	target := float32(0)
	if asFloat, err := strconv.ParseFloat(matches[1], 32); err == nil {
		purgeLength = float32(asFloat)
	} else {
		return err, 0, 0, 0, 0
	}
	if asFloat, err := strconv.ParseFloat(matches[2], 32); err == nil {
		transitionLength = float32(asFloat)
	} else {
		return err, 0, 0, 0, 0
	}
	if asFloat, err := strconv.ParseFloat(matches[3], 32); err == nil {
		offset = float32(asFloat)
	} else {
		return err, 0, 0, 0, 0
	}
	if asFloat, err := strconv.ParseFloat(matches[4], 32); err == nil {
		target = float32(asFloat)
	} else {
		return err, 0, 0, 0, 0
	}
	return nil, purgeLength, transitionLength, offset, target
}

func (s *generatorState) startDenseTowerSegment(purgeLength, transitionLength, offset, target float32) {
	s.transitioning = true
	s.extrusionSoFar = 0
	s.purgeLength = purgeLength
	s.transitionLength = transitionLength
	s.offset = offset
	s.target = target / 100
}

func interpolateTowerColor(linearT, target float32) float32 {
	minCutoff := target - 0.1
	maxCutoff := target + 0.35
	if linearT <= minCutoff {
		return 0.0
	}
	if linearT >= maxCutoff {
		return 1.0
	}
	return (linearT - minCutoff) * (1 / (maxCutoff - minCutoff))
}

// must be called after updating extrusionSoFar
func (s *generatorState) getT() float32 {
	if !s.transitioning {
		return 0
	}
	return interpolateTowerColor((s.extrusionSoFar-s.offset)/s.purgeLength, s.target)
}

func parseArgvFloat32(arg string) (float32, error) {
	if val, err := strconv.ParseFloat(arg, 32); err != nil {
		return 0, err
	} else {
		return float32(val), nil
	}
}

func generateToolpath(argv []string) error {
	argc := len(argv)

	if argc != 7 {
		return errors.New("expected 7 command-line arguments")
	}
	inpath := argv[0]
	outpath := argv[1]
	initialExtrusionWidth, err := parseArgvFloat32(argv[2])
	if err != nil {
		return err
	}
	initialLayerHeight, err := parseArgvFloat32(argv[3])
	zOffset, err := parseArgvFloat32(argv[4])
	if err != nil {
		return err
	}
	brimIsSkirt := argv[5] == "true"
	toolColors, err := parseToolColors(argv[6])
	if err != nil {
		return err
	}
	preflight, err := toolpathPreflight(inpath)
	if err != nil {
		return err
	}

	writer := NewWriter(outpath, initialExtrusionWidth, initialLayerHeight, zOffset, brimIsSkirt, toolColors)
	writer.SetFeedrateBounds(preflight.minFeedrate, preflight.maxFeedrate)
	writer.SetTemperatureBounds(preflight.minTemperature, preflight.maxTemperature)
	writer.SetLayerHeightBounds(preflight.minLayerHeight, preflight.maxLayerHeight)
	if err = writer.Initialize(); err != nil {
		return err
	}

	state := getStartingGeneratorState()
	err = gcode.ReadByLine(inpath, func(line gcode.Command, _ int) error {
		if setExtrusionMode, relative := line.IsSetExtrusionMode(); setExtrusionMode {
			state.relativeE = relative
			state.currentE = 0
		} else if line.IsSetPosition() {
			if e, ok := line.Params["e"]; ok {
				state.currentE = e
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
				eIncreased := e > state.currentE
				eDecreased := e < state.currentE
				if state.relativeE {
					eIncreased = e > 0
					eDecreased = e < 0
				}
				if state.transitioning {
					deltaE := e - state.currentE
					if state.relativeE {
						deltaE = e
					}
					state.extrusionSoFar += deltaE
				}
				if eIncreased {
					if isVisibleMove {
						isPrintMove = true
					} else {
						if err = writer.AddRestart(); err != nil {
							return err
						}
					}
				} else if eDecreased {
					// don't add retract point until wipe sequence, if any, is complete
					if !writer.state.inWipe {
						// add retract point regardless of there being X/Y/Z movement as well
						if err = writer.AddRetract(); err != nil {
							return err
						}
					}
				}
				state.currentE = e
			}
			if f, ok := line.Params["f"]; ok {
				if err = writer.SetFeedrate(f); err != nil {
					return err
				}
			}
			if isVisibleMove {
				if isPrintMove {
					if state.transitioning {
						t := state.getT()
						if err = writer.AddXYZTransitionLineTo(x, y, z, state.lastTool, t); err != nil {
							return err
						}
					} else {
						if err = writer.AddXYZPrintLineTo(x, y, z); err != nil {
							return err
						}
					}
				} else {
					if err = writer.AddXYZTravelTo(x, y, z); err != nil {
						return err
					}
				}
			}
		} else if line.Command == "M106" {
			// ignore P10, which is specifically assigned to the cooling module
			if fanIndex, ok := line.Params["p"]; !ok || fanIndex != 10 {
				if pwm, ok := line.Params["s"]; ok {
					if err = writer.SetFanSpeed(int(pwm)); err != nil {
						return err
					}
				}
			}
		} else if line.Command == "M107" {
			// ignore P10, which is specifically assigned to the cooling module
			if fanIndex, ok := line.Params["p"]; !ok || fanIndex != 10 {
				if err = writer.SetFanSpeed(0); err != nil {
					return err
				}
			}
		} else if line.Command == "M104" {
			if temp, ok := line.Params["s"]; ok {
				if err = writer.SetTemperature(temp); err != nil {
					return err
				}
			}
		} else if line.Command == "M109" {
			if temp, ok := line.Params["s"]; ok {
				if err = writer.SetTemperature(temp); err != nil {
					return err
				}
			} else if temp, ok = line.Params["r"]; ok {
				if err = writer.SetTemperature(temp); err != nil {
					return err
				}
			}
		} else if isToolChange, tool := line.IsToolChange(); isToolChange {
			if err = writer.SetTool(tool); err != nil {
				return err
			}
		} else if line.Command == "M135" {
			if t, ok := line.Params["t"]; ok {
				if err = writer.SetTool(int(t)); err != nil {
					return err
				}
			}
		} else if line.Command == "O31" {
			if err = writer.AddPing(); err != nil {
				return err
			}
		} else if line.Comment != "" {
			if line.Comment == "WIPE_START" {
				writer.state.inWipe = true
			} else if line.Comment == "WIPE_END" {
				// retract points were not added during the wipe sequence
				if writer.state.inWipe {
					// add retract point regardless of there being X/Y/Z movement as well
					if err = writer.AddRetract(); err != nil {
						return err
					}
				}
				writer.state.inWipe = false
			} else if strings.HasPrefix(line.Comment, "Z:") {
				z, err := strconv.ParseFloat(line.Comment[2:], 32)
				if err != nil {
					return err
				}
				state.currentLayerZ = float32(z)
				// TODO: do we need to explicitly also trigger a layer change at end sequence?
			} else if IsPathTypeComment(line) {
				// path type hints
				pathType := convertPathType(line.Comment[5:])
				if err = writer.SetPathType(pathType); err != nil {
					return err
				}
			} else if IsWidthComment(line) {
				// extrusion width hints
				width, err := strconv.ParseFloat(line.Comment[6:], 32)
				if err != nil {
					return err
				}
				if err = writer.SetExtrusionWidth(float32(width)); err != nil {
					return err
				}
			} else if IsHeightComment(line) {
				// layer height hints
				height, err := strconv.ParseFloat(line.Comment[7:], 32)
				if err != nil {
					return err
				}
				if err = writer.SetLayerHeight(roundZ(float32(height))); err != nil {
					return err
				}
			} else if strings.HasPrefix(line.Comment, "PTP_TYPE:") {
				err, purgeLength, transitionLength, offset, target := parsePtpTowerComment(line.Comment)
				if err != nil {
					return err
				}
				state.startDenseTowerSegment(purgeLength, transitionLength, offset, target)
			} else if strings.HasPrefix(line.Comment, "PTP_END") {
				state.transitioning = false
			} else if strings.HasPrefix(line.Comment, "Printing with input ") {
				tool, err := strconv.ParseInt(line.Comment[20:], 10, 32)
				if err != nil {
					return err
				}
				state.lastTool = state.currentTool
				state.currentTool = int(tool)
				if err = writer.SetTool(state.currentTool); err != nil {
					return err
				}
			} else if line.Raw == ";END OF LAYER CHANGE SEQUENCE" {
				if err = writer.LayerChange(state.currentLayerZ); err != nil {
					return err
				}
			}
		}

		// calculate bounding box
		if writer.state.currentPathType != PathTypeTravel &&
			writer.state.currentPathType != PathTypeSequence &&
			writer.state.currentPathType != PathTypeUnknown {
			x, y, z := writer.GetCurrentPosition()
			currentExtrusionRadius := (writer.state.currentExtrusionWidth) / 2
			// x: account for the extrusion radius
			writer.state.boundingBox.Min.X = MinFloat32(writer.state.boundingBox.Min.X, x-currentExtrusionRadius)
			writer.state.boundingBox.Max.X = MaxFloat32(writer.state.boundingBox.Max.X, x+currentExtrusionRadius)
			// y: account for the extrusion radius
			writer.state.boundingBox.Min.Y = MinFloat32(writer.state.boundingBox.Min.Y, y-currentExtrusionRadius)
			writer.state.boundingBox.Max.Y = MaxFloat32(writer.state.boundingBox.Max.Y, y+currentExtrusionRadius)
			// z
			// min: calculate min from the bottom of each path
			writer.state.boundingBox.Min.Z = MinFloat32(writer.state.boundingBox.Min.Z, z-writer.state.currentLayerHeight)
			writer.state.boundingBox.Max.Z = MaxFloat32(writer.state.boundingBox.Max.Z, z)
		}

		return nil
	})

	if err != nil {
		return err
	}

	// write bounding box info as a JSON to outPath.summary
	summaryPath := inpath + ".summary"
	summary := Summary{
		BoundingBox: writer.state.boundingBox,
	}

	if err = summary.Save(summaryPath); err != nil {
		return err
	}

	return writer.Finalize()
}

func GenerateToolpath(argv []string) {
	if err := generateToolpath(argv); err != nil {
		log.Fatalln(err)
	}
}
