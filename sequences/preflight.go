package sequences

import (
    "../gcode"
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
    nextX float64
    nextY float64
    nextZ float64
}

type sequencesPreflight struct {
    totalLayers int
    totalTime int
    startSequenceNextPos lookahead
    layerChangeNextPos []lookahead
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
                results.layerChangeNextPos = append(results.materialChangeNextPos, currentLookahead)
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
        } else if line.Raw == startPlaceholder {
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