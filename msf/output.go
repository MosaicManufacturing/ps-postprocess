package msf

import (
	"bufio"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"mosaicmfg.com/ps-postprocess/gcode"
	"mosaicmfg.com/ps-postprocess/ptp"
	"mosaicmfg.com/ps-postprocess/sequences"
)

func _paletteOutput(
	readerFn func(callback gcode.LineCallback) error,
	writer *bufio.Writer,
	msfOut *MSF,
	palette *Palette,
	preflight *msfPreflight,
	locals sequences.Locals,
) error {
	// initialize state
	state := NewState(palette)
	state.MSF = msfOut
	state.TowerBoundingBox = preflight.towerBoundingBox
	state.TransitionNextPositions = preflight.transitionNextPositions
	locals.Global["totalTime"] = float64(preflight.timeEstimate)
	locals.Global["totalLayers"] = float64(preflight.totalLayers)

	state.Locals = locals
	// account for a firmware purge (not part of G-code) once
	state.E.TotalExtrusion += palette.FirmwarePurge
	state.TimeEstimate = preflight.timeEstimate

	if len(preflight.pingStarts) > 0 {
		state.NextPingStart = preflight.pingStarts[0]
	} else {
		state.NextPingStart = PingMinSpacing
	}

	if palette.TransitionMethod == CustomTower {
		tower, needsTower := GenerateTower(palette, preflight)
		if !needsTower {
			log.Fatalln("should not have generated a tower!")
		}
		state.Tower = &tower
	}

	didFinalSplice := false             // used to prevent calling msfOut.AddLastSplice multiple times
	upcomingSparseLayer := false        // used for special-case wipe sequence handling
	upcomingDoubledSparseLayer := false // used for special-case layer change handling

	insertNonDoubledSparseLayer := func() error {
		if err := writeLine(writer, "; Sparse tower layer"); err != nil {
			return err
		}
		retractDistance := palette.RetractDistance[state.CurrentTool]
		retractFeedrate := palette.RetractFeedrate[state.CurrentTool]
		if retractDistance != 0 {
			if retract := getRetract(&state, retractDistance, retractFeedrate); len(retract) > 0 {
				if err := writeLines(writer, retract); err != nil {
					return err
				}
			}
		} else if palette.UseFirmwareRetraction {
			retract := getFirmwareRetract()
			if err := writeLines(writer, retract); err != nil {
				return err
			}
		}
		if reset := resetEAxis(&state); len(reset) > 0 {
			if err := writeLines(writer, reset); err != nil {
				return err
			}
		}
		zLiftTarget := state.XYZF.CurrentZ + palette.ZLift[state.CurrentTool]
		if zLift := getZTravel(&state, zLiftTarget, "lift Z"); len(zLift) > 0 {
			if err := writeLines(writer, zLift); err != nil {
				return err
			}
		}
		layerPaths, err := state.Tower.GetNextSegment(&state, false)
		if err != nil {
			return err
		}
		if !state.Tower.IsComplete() && !state.Tower.CurrentLayerIsDense() {
			upcomingDoubledSparseLayer = true
		}
		return writeLines(writer, layerPaths)
	}

	restorePathType := func() error {
		if state.CurrentPathTypeLine != "" {
			if err := writeLine(writer, state.CurrentPathTypeLine); err != nil {
				return err
			}
		}
		if state.CurrentWidthLine != "" {
			if err := writeLine(writer, state.CurrentWidthLine); err != nil {
				return err
			}
		}
		return nil
	}

	err := readerFn(func(line gcode.Command, lineNumber int) error {
		if lineNumber == preflight.printSummaryStart {
			if err := msfOut.AddLastSplice(state.CurrentTool, state.E.TotalExtrusion); err != nil {
				return err
			}
			didFinalSplice = true // make sure not to do this again at EOF
			// insert our (more accurate) print summary
			summary := getPrintSummary(msfOut, state.TimeEstimate)
			if err := writeLines(writer, summary); err != nil {
				return err
			}
		}
		if upcomingDoubledSparseLayer && line.IsLinearMove() &&
			(line.Comment == "retract" || line.Comment == "lift Z") {
			// filter out these commands as we've already included them as needed
			return nil
		} else if upcomingDoubledSparseLayer && line.IsSetPosition() &&
			line.Comment == "reset extrusion distance" {
			// filter out these commands as we've already included them as needed
			return nil
		} else if line.IsLinearMove() && line.Comment == "retract" && state.E.LastExtrudeWasRetract {
			// avoid double-retraction after toolchange
			return nil
		} else if ptp.IsPathTypeComment(line) {
			state.CurrentPathTypeLine = line.Raw
		} else if ptp.IsWidthComment(line) {
			state.CurrentWidthLine = line.Raw
		} else {
			// update state
			state.E.TrackInstruction(line)
			state.XYZF.TrackInstruction(line)
			state.Temperature.TrackInstruction(line)
		}
		if state.NeedsPostTransitionZAdjust && line.IsLinearMove() {
			_, hasX := line.Params["x"]
			_, hasY := line.Params["y"]
			_, hasZ := line.Params["z"]
			eParam, hasE := line.Params["e"]
			isPrintLine := (hasX || hasY) && hasE
			isRestart := !(hasX || hasY || hasZ) && hasE && state.E.CurrentRetraction+eParam == 0
			if isPrintLine || isRestart {
				// restore pre-transition Z height immediately before doing a print line
				currentF := state.XYZF.CurrentFeedrate
				if err := writeLines(writer, getZTravel(&state, state.PostTransitionZ, "restore layer Z")); err != nil {
					return err
				}
				state.NeedsPostTransitionZAdjust = false
				state.PostTransitionZ = 0
				// restore most recent F value, as Z travel likely changed it
				// (not needed for restart command, which always includes an F value)
				if !isRestart {
					feedrateAdjustment := getFeedrateAdjust(&state, currentF)
					if len(feedrateAdjustment) > 0 {
						if err := writeLines(writer, feedrateAdjustment); err != nil {
							return err
						}
					}
				}
			}
		}
		if line.IsLinearMove() {
			if err := writeLine(writer, line.Raw); err != nil {
				return err
			}
			if state.OnWipeTower && state.Palette.SupportsPings() {
				// check for ping actions
				if state.CurrentlyPinging {
					// currentlyPinging == true implies accessory mode
					if state.E.TotalExtrusion >= state.LastPingStart+state.PingExtrusion {
						// finish the accessory ping sequence
						comment := fmt.Sprintf("; Ping %d pause 2%s", len(msfOut.PingList)+1, EOL)
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
						comment := fmt.Sprintf("; Ping %d", len(msfOut.PingList)+1)
						if err := writeLine(writer, comment); err != nil {
							return err
						}
						msfOut.AddPing(state.E.TotalExtrusion)
						if err := writeLine(writer, palette.ClearBufferCommand); err != nil {
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
						comment := fmt.Sprintf("; Ping %d pause 1%s", len(msfOut.PingList)+1, EOL)
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
		} else if upcomingDoubledSparseLayer && line.IsSetPosition() {
			return nil
		} else if isToolChange, tool := line.IsToolChange(); isToolChange {
			if state.PastStartSequence {
				if state.FirstToolChange {
					state.FirstToolChange = false
					var toolChangeLine string
					// for element replace extruder index with filament id
					if palette.Type == TypeElement {
						filamentId := palette.FilamentIds[tool] - 1
						toolChangeLine = fmt.Sprintf("; Printing with input %d", filamentId)
						fmt.Printf(" palette.FilamentIds: %d\n", palette.FilamentIds)
						fmt.Printf("tool: %d\n", tool)
						fmt.Printf("filamentId: %d\n", filamentId)
					} else {
						toolChangeLine = fmt.Sprintf("; Printing with input %d", palette.PrintExtruder)
					}
					if err := writeLine(writer, toolChangeLine); err != nil {
						return err
					}
					state.CurrentTool = tool
					if err := writeLine(writer, fmt.Sprintf("; Printing with input %d", state.CurrentTool)); err != nil {
						return err
					}
				} else if tool != state.CurrentTool {
					comment := fmt.Sprintf("; Printing with input %d", tool)
					if err := writeLine(writer, comment); err != nil {
						return err
					}
					if palette.Type == TypeElement {
						filamentId := palette.FilamentIds[tool] - 1
						// output a tool change
						toolChangeLine := fmt.Sprintf("; Printing with input %d", filamentId)
						if err := writeLine(writer, toolChangeLine); err != nil {
							return err
						}
						spliceLength := state.E.TotalExtrusion
						if err := msfOut.AddSplice(state.CurrentTool, spliceLength); err != nil {
							return err
						}
						state.CurrentTool = tool
					} else {
						if palette.TransitionMethod == CustomTower {
							if err := writeLine(writer, "; Dense tower segment"); err != nil {
								return err
							}
							currentTransition := state.Tower.GetCurrentTransitionInfo()
							spliceOffset := currentTransition.TransitionLength * (palette.TransitionTarget / 100)
							// if purge length is more than transition length, the extra purge is there
							// to ensure minimum piece lengths are maintained, so the difference between
							// the two should be included on the end of the previous tool's splice
							preTransitionAdd := currentTransition.PurgeLength - currentTransition.TransitionLength
							if preTransitionAdd < 0 {
								preTransitionAdd = 0
							}
							spliceOffset += preTransitionAdd
							spliceLength := state.E.TotalExtrusion + spliceOffset - currentTransition.UsableInfill

							ptpPurgeLength := currentTransition.PurgeLength
							ptpTransitionLength := currentTransition.TransitionLength
							ptpOffset := float32(0)

							if len(msfOut.SpliceList) == 0 {
								spliceLength += state.Tower.BrimExtrusion
								ptpOffset = state.Tower.BrimExtrusion
							}

							if err := msfOut.AddSplice(state.CurrentTool, spliceLength); err != nil {
								return err
							}
							state.CurrentTool = tool
							state.CurrentlyTransitioning = true
							ptpComment := getPtpStartComment(
								ptpPurgeLength,
								ptpTransitionLength,
								ptpOffset,
								palette.TransitionTarget,
							)
							if err := writeLines(writer, ptpComment); err != nil {
								return err
							}
							transition, err := state.Tower.GetNextSegment(&state, true)
							upcomingDoubledSparseLayer = false
							if err != nil {
								return err
							}
							if err := writeLines(writer, transition); err != nil {
								return err
							}
							state.CurrentlyTransitioning = false
							if err := writeLines(writer, getPtpEndComment()); err != nil {
								return err
							}
						} else {
							currentTransition := preflight.transitions[len(msfOut.SpliceList)]
							currentPurgeLength := currentTransition.PurgeLength
							spliceOffset := currentTransition.TransitionLength * (palette.TransitionTarget / 100)
							spliceLength := state.E.TotalExtrusion + spliceOffset - currentTransition.UsableInfill
							if palette.TransitionMethod == SideTransitions {
								extra := msfOut.GetRequiredExtraSpliceLength(spliceLength)
								if extra > 0 {
									currentPurgeLength += extra
									spliceLength += extra
								}
							}
							if err := msfOut.AddSplice(state.CurrentTool, spliceLength); err != nil {
								return err
							}
							state.CurrentTool = tool
							state.CurrentlyTransitioning = true
							if palette.TransitionMethod == SideTransitions {
								transition, err := sideTransition(currentPurgeLength, &state)
								if err != nil {
									return err
								}
								if err := writeLines(writer, transition); err != nil {
									return err
								}
								state.CurrentlyTransitioning = false
							}
						}
						if palette.TransitionMethod != TransitionTower {
							if err := restorePathType(); err != nil {
								return err
							}
						}
					}
				}
			}
		} else if line.Raw == ";START_OF_PRINT" {
			state.PastStartSequence = true
			return writeLine(writer, line.Raw)
		} else if line.Raw == ";LAYER_CHANGE" {
			state.CurrentLayer++
			state.CurrentLayerObject = 0
			state.CurrentLayerTotalObjects = preflight.layerObjectStarts[state.CurrentLayer]
			return writeLine(writer, line.Raw)
		} else if upcomingSparseLayer && line.Raw == ";WIPE_END" {
			upcomingSparseLayer = false
			// insert deferred sparse layer now
			return insertNonDoubledSparseLayer()
		} else if palette.TransitionMethod == CustomTower &&
			strings.HasPrefix(line.Comment, "stop printing object") {
			// sparse tower layer
			if err := writeLine(writer, line.Raw); err != nil {
				return err
			}
			state.CurrentLayerObject++
			if state.CurrentLayerObject == state.CurrentLayerTotalObjects {
				if !state.Tower.IsComplete() && !state.Tower.CurrentLayerIsDense() &&
					state.CurrentLayer == state.Tower.CurrentLayerIndex {
					if palette.Wipe[state.CurrentTool] {
						// need to look for ;WIPE_END
						upcomingSparseLayer = true
						return nil
					} else {
						// can start sparse layer immediately
						return insertNonDoubledSparseLayer()
					}
				}
			}
		} else if palette.TransitionMethod == CustomTower &&
			strings.HasPrefix(line.Comment, "printing object") {
			// doubled sparse tower layer
			if !state.Tower.IsComplete() &&
				upcomingDoubledSparseLayer &&
				state.CurrentLayer == state.Tower.CurrentLayerIndex &&
				!state.Tower.CurrentLayerIsDense() {
				if err := writeLine(writer, "; Doubled sparse tower layer"); err != nil {
					return err
				}
				layerPaths, err := state.Tower.GetNextSegment(&state, false)
				if err != nil {
					return err
				}
				upcomingDoubledSparseLayer = false
				if err := writeLines(writer, layerPaths); err != nil {
					return err
				}
			}
			return writeLine(writer, line.Raw)
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

	return nil
}

func paletteOutput(inpath, outpath, msfpath string, palette *Palette, preflight *msfPreflight, locals sequences.Locals) error {
	outfile, createErr := os.Create(outpath)
	if createErr != nil {
		return createErr
	}
	writer := bufio.NewWriter(outfile)
	msfOut := NewMSF(palette)

	readerFn := func(callback gcode.LineCallback) error {
		return gcode.ReadByLine(inpath, callback)
	}

	err := _paletteOutput(readerFn, writer, &msfOut, palette, preflight, locals)
	if err != nil {
		return err
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

	return nil
}
