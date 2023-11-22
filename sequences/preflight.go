package sequences

import (
	"mosaicmfg.com/ps-postprocess/gcode"
	"strings"
)

type lookaheadType int

const (
	lookaheadStart lookaheadType = iota
	lookaheadLayerChange
	lookaheadMaterialChange
)

type lookahead struct {
	lookaheadType lookaheadType
	nextX         float64
	nextY         float64
	nextZ         float64
}

type sequencesPreflight struct {
	preheat               PreheatHints
	totalLayers           int
	totalTime             int
	startSequenceNextPos  lookahead
	layerChangeNextPos    []lookahead
	materialChangeNextPos []lookahead
}

func preflight(inpath string) (sequencesPreflight, error) {
	results := sequencesPreflight{
		layerChangeNextPos:    make([]lookahead, 0),
		materialChangeNextPos: make([]lookahead, 0),
	}
	position := gcode.PositionTracker{}

	currentLookaheads := make([]lookahead, 0)

	commitCurrentLookaheads := func() {
		if len(currentLookaheads) == 0 {
			return
		}
		for _, currentLookahead := range currentLookaheads {
			switch currentLookahead.lookaheadType {
			case lookaheadStart:
				results.startSequenceNextPos = currentLookahead
			case lookaheadLayerChange:
				results.layerChangeNextPos = append(results.layerChangeNextPos, currentLookahead)
			case lookaheadMaterialChange:
				results.materialChangeNextPos = append(results.materialChangeNextPos, currentLookahead)
			}
		}
		currentLookaheads = make([]lookahead, 0)
	}

	addLookahead := func(lookaheadType lookaheadType) {
		currentLookaheads = append(currentLookaheads, lookahead{
			lookaheadType: lookaheadType,
			nextX:         float64(position.CurrentX),
			nextY:         float64(position.CurrentY),
			nextZ:         float64(position.CurrentZ),
		})
	}

	err := gcode.ReadByLine(inpath, func(line gcode.Command, lineNum int) error {
		position.TrackInstruction(line)

		if line.Command == "M104" {
			// - first extruder temperature in the print
			// - max extruder temperature in the print
			if temp, ok := line.Params["s"]; ok {
				if results.preheat.Extruder == 0 {
					results.preheat.Extruder = temp
				}
				if temp > results.preheat.ExtruderMax {
					results.preheat.ExtruderMax = temp
				}
			}
		} else if line.Command == "M109" {
			// - first extruder temperature in the print
			// - max extruder temperature in the print
			if temp, ok := line.Params["s"]; ok {
				if results.preheat.Extruder == 0 {
					results.preheat.Extruder = temp
				}
				if temp > results.preheat.ExtruderMax {
					results.preheat.ExtruderMax = temp
				}
			} else if temp, ok = line.Params["r"]; ok {
				if results.preheat.Extruder == 0 {
					results.preheat.Extruder = temp
				}
				if temp > results.preheat.ExtruderMax {
					results.preheat.ExtruderMax = temp
				}
			}
			return nil
		} else if results.preheat.Bed == 0 && line.Command == "M140" {
			// - first print bed temperature in the print
			if temp, ok := line.Params["s"]; ok {
				results.preheat.Bed = temp
			}
			return nil
		} else if results.preheat.Bed == 0 && line.Command == "M190" {
			// - first print bed temperature in the print
			if temp, ok := line.Params["s"]; ok {
				results.preheat.Bed = temp
			} else if temp, ok = line.Params["r"]; ok {
				results.preheat.Bed = temp
			}
			return nil
		} else if results.preheat.Chamber == 0 && line.Command == "M141" {
			// - first chamber temperature in the print
			if temp, ok := line.Params["s"]; ok {
				results.preheat.Chamber = temp
			}
			return nil
		} else if results.preheat.Chamber == 0 && line.Command == "M191" {
			// - first chamber temperature in the print
			if temp, ok := line.Params["s"]; ok {
				results.preheat.Chamber = temp
			} else if temp, ok = line.Params["r"]; ok {
				results.preheat.Chamber = temp
			}
			return nil
		}

		if len(currentLookaheads) > 0 {
			// logic: keep applying Z changes, and commit when we see X and/or Y change
			if line.IsLinearMove() {
				needsCommit := false
				if z, ok := line.Params["z"]; ok {
					for i := 0; i < len(currentLookaheads); i++ {
						currentLookaheads[i].nextZ = float64(z)
					}
				}
				if x, ok := line.Params["x"]; ok {
					for i := 0; i < len(currentLookaheads); i++ {
						currentLookaheads[i].nextX = float64(x)
					}
					needsCommit = true
				}
				if y, ok := line.Params["y"]; ok {
					for i := 0; i < len(currentLookaheads); i++ {
						currentLookaheads[i].nextY = float64(y)
					}
					needsCommit = true
				}
				if needsCommit {
					commitCurrentLookaheads()
				}
				return nil
			}
		}

		if line.Command != "" {
			return nil
		}

		if line.Comment == "LAYER_CHANGE" {
			results.totalLayers++
		} else if strings.HasPrefix(line.Comment, "estimated printing time (normal mode) = ") {
			timeEstimate, err := gcode.ParseTimeString(line.Comment)
			if err != nil {
				return err
			}
			results.totalTime = int(timeEstimate)
		} else if line.Raw == endOfStartPlaceholder {
			addLookahead(lookaheadStart)
		} else if strings.HasPrefix(line.Raw, layerChangePrefix) {
			addLookahead(lookaheadLayerChange)
		} else if strings.HasPrefix(line.Raw, materialChangePrefix) {
			addLookahead(lookaheadMaterialChange)
		}
		return nil
	})
	commitCurrentLookaheads()
	return results, err
}
