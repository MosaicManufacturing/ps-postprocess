package msf

import (
    "../gcode"
    "errors"
    "fmt"
    "math"
)

type TowerLayer struct {
    TopZ float32 // Z value of the top of this layer
    Thickness float32 // height of extruded paths
    Density float32 // 0..1
    Transitions []Transition
}

func (l TowerLayer) String() string {
    return fmt.Sprintf("TopZ = %.2f mm, Thickness = %.2f mm, Density = %.1f%%, Transitions = %d", l.TopZ, l.Thickness, l.Density * 100, len(l.Transitions))
}

type Tower struct {
    // constant or pre-calculated
    Palette *Palette
    BoundingBox gcode.BoundingBox
    Layers []TowerLayer

    // for use during output
    CurrentLayerPaths []gcode.Command // feedrates, raw strings, real E values not included yet
    CurrentLayerIndex int // total transitions on this layer
    CurrentLayerTransitionIndex int // current transition on this layer
    CurrentLayerCommandIndex int // index into CurrentLayerPaths
}

func GenerateTower(palette *Palette, preflight *msfPreflight) (Tower, bool) {
    totalLayers := preflight.totalLayers + 1
    tower := Tower{
        Palette:                       palette,
        BoundingBox:                   gcode.NewBoundingBox(),
        Layers:                        nil,
        CurrentLayerPaths:             nil,
        CurrentLayerIndex:             0,
        CurrentLayerTransitionIndex:   0,
        CurrentLayerCommandIndex:      0,
    }

    minDensity := float64(palette.TowerMinDensity) / 100
    minFirstLayerDensity := float64(palette.TowerMinFirstLayerDensity) / 100
    maxDensity := float64(palette.TowerMaxDensity) / 100
    extrusionMultiplier := palette.TowerExtrusionMultiplier / 100

    // tower must have at least this many mm3 of extrusion to be able to fit pings!
    minLayerVolume := float64(filamentLengthToVolume(palette.GetPingExtrusion())) / maxDensity

    // 1. determine the number of transitions required on each layer

    layerTransitionCounts := make([]int, totalLayers)
    for layer, transitions := range preflight.transitionsByLayer {
        layerTransitionCounts[layer] = len(transitions)
    }

    // 2. discard sparse layers
    //    - start from the top and work downwards
    //    - discard sparse layers until a layer with more than zero transitions is reached

    for i := totalLayers - 1; i >= 0; i-- {
        if layerTransitionCounts[i] == 0 {
            totalLayers--
        } else {
            break
        }
    }
    if totalLayers == 0 {
        // no dense layers == no Palette processing
        return tower, false
    }
    layerTransitionCounts = layerTransitionCounts[:totalLayers]

    // 3. determine the thickness and top Z of each layer
    //    - also store the physical Z height of each layer, for G-code output

    layerThicknesses := make([]float32, totalLayers)
    layerTopZs := make([]float32, totalLayers)
    lastTopZ := float32(0)
    for i := 0; i < totalLayers; i++ {
        topZ := preflight.layerTopZs[i]
        thickness := preflight.layerThicknesses[i]
        if thickness > topZ - lastTopZ {
            // ensure layer is no thicker than this Z minus the previous Z
            thickness = topZ - lastTopZ
        }
        layerTopZs[i] = topZ
        layerThicknesses[i] = thickness
        lastTopZ = topZ
    }

    // 4. determine the volume of filament required on each layer
    //    - number of transitions on the layer
    //    - thickness of the layer
    //    - transition purge lengths

    // 5. determine the 2D area (footprint) required by each layer
    //    - layer volume and layer thickness
    //    - minimum and maximum tower density, minimum first layer density
    //    - extrusion multiplier/flow rate, extrusion width

    // 6. determine the overall footprint size of the tower
    //    - calculated as the greatest 2D footprint of all the layers

    layerFootprintAreas := make([]float32, totalLayers)
    minFootprintArea := float64(0)
    for layer, transitions := range preflight.transitionsByLayer {
        layerPurgeLength := float32(0) // mm
        for _, transition := range transitions {
            layerPurgeLength += transition.PurgeLength / extrusionMultiplier
        }
        layerPurgeVolume := filamentLengthToVolume(layerPurgeLength) // mm3
        // adjust for max density
        layerPurgeVolume /= float32(maxDensity)
        // raise the volume slightly to account for errors in total toolpath extrusion
        layerPurgeVolume *= 1.05
        // ensure the layer has room for at least one ping
        layerPurgeVolume = float32(math.Max(float64(layerPurgeVolume), minLayerVolume))

        layerFootprintArea := layerPurgeVolume / layerThicknesses[layer]
        layerFootprintAreas[layer] = layerFootprintArea
        if layerPurgeVolume > 0 {
            minFootprintArea = math.Max(minFootprintArea, float64(layerFootprintArea))
        }
    }
    if minFootprintArea == 0 {
        // no dense layers == no Palette processing
        return tower, false
    }

    // 7. determine the density of each layer
    //    - minimum and maximum tower density, minimum first layer density
    //    - ratio of required footprint area for this layer to overall footprint of the tower
    layerDensities := make([]float32, totalLayers)
    for layer := 0; layer < totalLayers; layer++ {
        footprintArea := layerFootprintAreas[layer]
        density := float64(footprintArea) / minFootprintArea
        if layer == 0 {
            density = math.Max(density, minFirstLayerDensity)
        } else {
            density = math.Max(density, minDensity)
        }
        density = math.Min(density, maxDensity)
        layerDensities[layer] = float32(density)
    }

    // 8. finalize the tower dimensions
    //    - try and maintain the current aspect ratio

    towerWidth := float64(palette.TowerSize[0])
    towerHeight := float64(palette.TowerSize[1])
    squareLength := math.Sqrt(minFootprintArea)
    aspectRatio := 1 / math.SqrtPhi // default to the golden ratio
    if towerWidth > 0 && towerHeight > 0 {
        // prefer to use the provided aspect ratio, but not necessarily size
        aspectRatio = math.Sqrt(towerWidth / towerHeight)
    }
    towerHalfHeight := float32(squareLength / aspectRatio) / 2
    towerHalfWidth := float32(squareLength * aspectRatio) / 2

    // 9. store everything relevant

    tower.Layers = make([]TowerLayer, totalLayers)
    for layer := 0; layer < totalLayers; layer++ {
        tower.Layers[layer] = TowerLayer{
            TopZ:      layerTopZs[layer],
            Thickness: layerThicknesses[layer],
            Density:   layerDensities[layer],
            Transitions: preflight.transitionsByLayer[layer],
        }
    }

    tower.BoundingBox.Min[0] = palette.TowerPosition[0] - towerHalfWidth
    tower.BoundingBox.Max[0] = palette.TowerPosition[0] + towerHalfWidth
    tower.BoundingBox.Min[1] = palette.TowerPosition[1] - towerHalfHeight
    tower.BoundingBox.Max[1] = palette.TowerPosition[1] + towerHalfHeight
    tower.BoundingBox.Min[2] = layerTopZs[0] - layerThicknesses[0]
    tower.BoundingBox.Max[2] = layerTopZs[len(layerTopZs)-1]

    return tower, true
}

func (t *Tower) layerNeedsPerimeters(layer int, density float32) bool {
    if layer < t.Palette.RaftLayers {
        // raft layers never get perimeters
        return false
    }
    if layer == t.Palette.RaftLayers && t.Palette.TowerFirstLayerPerimeters {
        // first layer -- force a perimeter if desired by user
        return true
    }
    return density <= TowerPerimeterThreshold
}

func (t *Tower) rasterizeLayer(layer int) {
    xMin := t.BoundingBox.Min[0]
    xMax := t.BoundingBox.Max[0]
    yMin := t.BoundingBox.Min[1]
    yMax := t.BoundingBox.Max[1]
    if layer < t.Palette.RaftLayers {
        inflation := t.Palette.RaftInflation
        if layer == 0 {
            inflation *= 2
        }
        xMin -= inflation
        yMin -= inflation
        xMax += inflation
        yMax += inflation
    }

    t.CurrentLayerPaths = make([]gcode.Command, 0)
    currentXMin := xMin
    currentXMax := xMax
    currentYMin := yMin
    currentYMax := yMax

    density := t.Layers[layer].Density
    extrusionWidth := t.Palette.TowerExtrusionWidth
    if layer < t.Palette.RaftLayers {
        extrusionWidth = t.Palette.RaftExtrusionWidth
    }
    addPerimeters := t.layerNeedsPerimeters(layer, density)

    // create perimeters

    if addPerimeters {
        for i := 0; i < TowerPerimeterCount; i++ {
            t.CurrentLayerPaths = append(t.CurrentLayerPaths,
                // travel to southeast corner
                gcode.Command{
                    Command: "G1",
                    Params:  map[string]float32{
                        "x": currentXMax,
                        "y": currentYMin,
                        "f": 0,
                    },
                },
                // extrude to northeast corner
                gcode.Command{
                    Command: "G1",
                    Params:  map[string]float32{
                        "x": currentXMax,
                        "y": currentYMax,
                        "e": 0,
                        "f": 0,
                    },
                },
                // extrude to northwest corner
                gcode.Command{
                    Command: "G1",
                    Params:  map[string]float32{
                        "x": currentXMin,
                        "y": currentYMax,
                        "e": 0,
                        "f": 0,
                    },
                },
                // extrude to southwest corner
                gcode.Command{
                    Command: "G1",
                    Params:  map[string]float32{
                        "x": currentXMin,
                        "y": currentYMin,
                        "e": 0,
                        "f": 0,
                    },
                },
                // extrude to southeast corner
                gcode.Command{
                    Command: "G1",
                    Params:  map[string]float32{
                        "x": currentXMax,
                        "y": currentYMin,
                        "e": 0,
                        "f": 0,
                    },
                },
            )
            currentXMin += extrusionWidth
            currentYMin += extrusionWidth
            currentXMax -= extrusionWidth
            currentYMax -= extrusionWidth
        }
        overlap := extrusionWidth * (t.Palette.InfillPerimeterOverlap / 100)
        currentXMin -= overlap
        currentYMin -= overlap
        currentXMax += overlap
        currentYMax += overlap
    }

    // create infill

    reverse := layer % 2 == 1
    stride := (1 / t.Layers[layer].Density) * extrusionWidth
    if layer < t.Palette.RaftLayers {
        stride = t.Palette.RaftStride
    }
    axisAlignedStride := float32(math.Sqrt(float64(stride * stride * 2)))

    firstLine := true
    needsMoreLines := true
    xBoundReached := false
    yBoundReached := false
    // assume infill is drawn southwest/northeast
    // (X coordinates will be mirrored if reverse == true)
    printSouthwest := true

    neX := currentXMax
    neY := currentYMin
    swX := neX
    swY := neY

    for needsMoreLines {
        if xBoundReached {
            swY += axisAlignedStride // move the southwest corner north
        } else {
            swX -= axisAlignedStride // move the southwest corner west
        }

        if yBoundReached {
            neX -= axisAlignedStride // move the northeast corner west
        } else {
            neY += axisAlignedStride // move the northeast corner north
        }

        // check if the southwest corner was reached
        if !xBoundReached && swX < currentXMin {
            swY = currentYMin + (currentXMin - swX)
            swX = currentXMin
            xBoundReached = true
        }

        // check if the northeast corner was reached
        if !yBoundReached && neY > currentYMax {
            neX = currentXMax - (neY - currentYMax)
            neY = currentYMax
            yBoundReached = true
        }

        x1 := neX; y1 := neY
        x2 := swX; y2 := swY
        if !printSouthwest {
            // reverse direction of path
            x2, y2, x1, y1 = x1, y1, x2, y2
        }

        if reverse {
            // mirror layer's X values
            x1 = currentXMax + currentXMin - x1
            x2 = currentXMax + currentXMin - x2
        }

        // travel to (x1, y1)
        travel := gcode.Command{
            Command: "G1",
            Params:  map[string]float32{
                "x": x1,
                "y": y1,
                "f": 0,
            },
        }
        // extrude to (x2, y2)
        extrude := gcode.Command{
            Command: "G1",
            Params:  map[string]float32{
                "x": x2,
                "y": y2,
                "e": 0,
                "f": 0,
            },
        }

        if !firstLine && layer < t.Palette.RaftLayers {
            // raft layers should have continuous extrusion
            // (but the first command should still be travel)
            travel.Params["e"] = 0
        }

        t.CurrentLayerPaths = append(t.CurrentLayerPaths, travel, extrude)
        firstLine = false

        if neX - currentXMin < axisAlignedStride && currentYMax - swY < axisAlignedStride {
            // layer has been fully rasterized
            needsMoreLines = false
        } else {
            // reverse the direction of the next line
            printSouthwest = !printSouthwest
        }
    }

}

func (t *Tower) IsComplete() bool {
    return t.CurrentLayerIndex >= len(t.Layers)
}

func (t *Tower) CurrentLayerTopZ() float32 {
    if t.IsComplete() {
        // no more layers -- N/A
        return -1
    }
    return t.Layers[t.CurrentLayerIndex].TopZ
}

func (t *Tower) CurrentLayerIsDense() bool {
    if t.IsComplete() {
        // no more layers -- N/A
        return false
    }
    return len(t.Layers[t.CurrentLayerIndex].Transitions) > 0
}

func (t *Tower) NeedsSparseLayers(nextLayer int) bool {
    // if true, tower has not been printed to the current layer height yet
    // -- at least one sparse layer should be inserted
    return nextLayer > t.CurrentLayerIndex
}

func (t *Tower) GetCurrentTransitionInfo() *Transition {
    if t.CurrentLayerIndex >= len(t.Layers) {
        return nil
    }
    if t.CurrentLayerTransitionIndex >= len(t.Layers[t.CurrentLayerIndex].Transitions) {
        return nil
    }
    return &t.Layers[t.CurrentLayerIndex].Transitions[t.CurrentLayerTransitionIndex]
}

func (t *Tower) moveToTower(state *State) (string, error) {
    sequence := ";TYPE:Wipe tower" + EOL
    sequence += fmt.Sprintf(";WIDTH:%s%s", gcode.FormatFloat(float64(t.Palette.TowerExtrusionWidth)), EOL)
    sequence += fmt.Sprintf(";HEIGHT:%s%s", gcode.FormatFloat(float64(t.Layers[t.CurrentLayerIndex].Thickness)), EOL)

    // next tower command should always be a travel
    travel := t.CurrentLayerPaths[t.CurrentLayerCommandIndex]
    if _, ok := travel.Params["e"]; ok {
        return "", errors.New("tower segment started with extrusion, not travel")
    }
    travel.Params["f"] = state.Palette.TravelSpeedXY
    travel.Comment = "move to tower"
    t.CurrentLayerCommandIndex++ // use up the command
    sequence += travel.String() + EOL

    state.TimeEstimate += estimateMoveTime(state.XYZF.CurrentX, state.XYZF.CurrentY, travel.Params["x"], travel.Params["y"], travel.Params["f"])
    state.XYZF.TrackInstruction(travel)

    // z-lift down if needed
    if topZ := t.CurrentLayerTopZ(); state.XYZF.CurrentZ > topZ {
        sequence += getZTravel(state, topZ, "restore layer Z")
    }

    if state.E.CurrentRetraction < 0 {
        // un-retract
        sequence += getRestart(state, state.E.CurrentRetraction, state.Palette.RestartFeedrate[state.CurrentTool])
    }
    return sequence, nil
}

func (t *Tower) leaveTower(state *State, retractDistance float32) string {
    sequence := ""
    if retractDistance != 0 {
        // restore any retraction from before tower was started
        sequence += getRetract(state, retractDistance, state.Palette.RetractFeedrate[state.CurrentTool])
    }
    sequence += resetEAxis(state)
    if state.Palette.ZLift[state.CurrentTool] > 0 {
        // lift z
        sequence += getZTravel(state, state.XYZF.CurrentZ + state.Palette.ZLift[state.CurrentTool], "lift Z")
    }
    return sequence
}

func (t *Tower) getNextPath(state *State, printFeedrate float32) (string, float32) {
    extrusionWidth :=  t.Palette.TowerExtrusionWidth
    layerHeight := t.Layers[t.CurrentLayerIndex].Thickness
    extrusionMultiplier :=  t.Palette.TowerExtrusionMultiplier / 100

    command := t.CurrentLayerPaths[t.CurrentLayerCommandIndex]
    commandExtrusion := float32(0)

    // when printing a segment, all commands use the print feedrate
    // so as not to alternate feedrates constantly
    command.Params["f"] = printFeedrate
    if _, ok := command.Params["e"]; ok {
        // extrusion command
        lineLength := getLineLength(state.XYZF.CurrentX, state.XYZF.CurrentY, command.Params["x"], command.Params["y"]) // mm
        deltaE := getExtrusionLength(extrusionWidth, layerHeight, lineLength) * extrusionMultiplier
        if state.E.RelativeExtrusion {
            command.Params["e"] = deltaE
        } else {
            command.Params["e"] = state.E.CurrentExtrusionValue + deltaE
        }
        commandExtrusion = deltaE
    } else {
        // travel command
    }
    currentX := state.XYZF.CurrentX
    currentY := state.XYZF.CurrentY
    currentFeedrate := state.XYZF.CurrentFeedrate

    state.TimeEstimate += estimateMoveTime(currentX, currentY, command.Params["x"], command.Params["y"], command.Params["f"])
    state.XYZF.TrackInstruction(command)
    state.E.TrackInstruction(command)

    // optimize output size by taking advantage of sticky parameters
    if currentFeedrate == command.Params["f"] {
        delete(command.Params, "f")
    }

    sequence := command.String() + EOL

    t.CurrentLayerCommandIndex++

    return sequence, commandExtrusion
}

func (t *Tower) getNextDenseSegmentPaths(state *State) string {
    transitionInfo := t.GetCurrentTransitionInfo()
    totalPurge := float32(0)

    printFeedrate := t.Palette.TowerSpeed[transitionInfo.To] * 60
    if t.Palette.TowerSpeed[transitionInfo.From] < t.Palette.TowerSpeed[transitionInfo.To] {
        // use the slower of the two material settings for this transition
        printFeedrate = t.Palette.TowerSpeed[transitionInfo.From] * 60
    }

    sequence := ""

    // last segment of the layer: finish the layer
    // all other segments: extrude just the purge length of this segment
    thisLayerTransitions := len(t.Layers[t.CurrentLayerIndex].Transitions)
    for (totalPurge < transitionInfo.PurgeLength || t.CurrentLayerTransitionIndex == thisLayerTransitions - 1) &&
      t.CurrentLayerCommandIndex < len(t.CurrentLayerPaths) {
        commandString, commandExtrusion := t.getNextPath(state, printFeedrate)
        sequence += commandString
        totalPurge += commandExtrusion
    }

    return sequence
}

func (t *Tower) getNextSparseLayerPaths(state *State) string {
    printFeedrate := t.Palette.TowerSpeed[state.CurrentTool] * 60

    sequence := ""

    // sparse layer: do the entire layer
    for t.CurrentLayerCommandIndex < len(t.CurrentLayerPaths) {
        commandString, _ := t.getNextPath(state, printFeedrate)
        sequence += commandString
    }

    return sequence
}

func (t *Tower) GetNextSegment(state *State, expectingDense bool) (string, error) {
    if t.CurrentLayerTransitionIndex == 0 && len(t.CurrentLayerPaths) == 0 {
        t.rasterizeLayer(t.CurrentLayerIndex)
    }

    // assertions for current layer being dense/sparse
    if expectingDense {
        if !t.CurrentLayerIsDense() {
            return "", fmt.Errorf("expected dense layer but layer %d is sparse", t.CurrentLayerIndex)
        }
    } else {
        if t.CurrentLayerIsDense() {
            return "", fmt.Errorf("expected sparse layer but layer %d is dense", t.CurrentLayerIndex)
        }
    }

    currentRetraction := state.E.CurrentRetraction

    sequence, err := t.moveToTower(state)
    if err != nil {
        return "", err
    }

    // TODO: add tower brims if first segment of first layer
    //  - respect user's minimum brim count
    //  - auto-increase brim count to ensure minimum first piece length

    if expectingDense {
        sequence += t.getNextDenseSegmentPaths(state)
    } else {
        sequence += t.getNextSparseLayerPaths(state)
    }

    sequence += t.leaveTower(state, currentRetraction)

    // move to the next transition on this layer
    t.CurrentLayerTransitionIndex++
    if t.CurrentLayerTransitionIndex >= len(t.Layers[t.CurrentLayerIndex].Transitions) {
        // move to the first transition on the next layer
        t.CurrentLayerIndex++
        t.CurrentLayerTransitionIndex = 0
        t.CurrentLayerCommandIndex = 0
        t.CurrentLayerPaths = nil
    }

    return sequence, nil
}
