package ptp

import (
	"errors"
	"log"
	"mosaicmfg.com/ps-postprocess/gcode"
	"regexp"
	"strconv"
	"strings"
)

type generatorState struct {
	// used for all generation
	currentTool int
	currentE    float32
	relativeE   bool

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

func parsePtpTowerComment(comment string) (error, float32, float32, float32, float32) {
	re := regexp.MustCompile("\\(purge=(.*),transition=(.*),offset=(.*),target=(.*)\\)")
	matches := re.FindStringSubmatch(comment)
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
	return interpolateTowerColor(s.extrusionSoFar/s.purgeLength, s.target)
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
						writer.AddRestart()
					}
				} else if eDecreased {
					// don't add retract point until wipe sequence, if any, is complete
					if !writer.state.inWipe {
						// add retract point regardless of there being X/Y/Z movement as well
						writer.AddRetract()
					}
				}
				state.currentE = e
			}
			if f, ok := line.Params["f"]; ok {
				writer.SetFeedrate(f)
			}
			if isVisibleMove {
				if isPrintMove {
					if state.transitioning {
						t := state.getT()
						writer.AddXYZTransitionLineTo(x, y, z, state.lastTool, t)
					} else {
						writer.AddXYZPrintLineTo(x, y, z)
					}
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
				writer.SetTool(state.currentTool)
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
