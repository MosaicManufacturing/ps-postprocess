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

type AnnotatedCommand struct {
    gcode gcode.Command
    extrusion float32
}

type Tower struct {
    // constant or pre-calculated
    Palette *Palette
    BoundingBox gcode.BoundingBox
    Layers []TowerLayer
    BrimCount int
    BrimExtrusion float32

    // for use during output
    CurrentLayerPaths []AnnotatedCommand // feedrates, raw strings, real E values not included yet
    CurrentLayerIndex int // total transitions on this layer
    CurrentLayerTransitionIndex int // current transition on this layer
    CurrentLayerCommandIndex int // index into CurrentLayerPaths
    CurrentLayerExtrusion float32 // sum of extrusions in CurrentLayerPaths
}

func GenerateTower(palette *Palette, preflight *msfPreflight) (Tower, bool) {
    totalLayers := preflight.totalLayers + 1
    tower := Tower{
        Palette:     palette,
        BoundingBox: gcode.NewBoundingBox(),
    }

    minDensity := float64(palette.TowerMinDensity) / 100
    minFirstLayerDensity := float64(palette.TowerMinFirstLayerDensity) / 100
    maxDensity := float64(palette.TowerMaxDensity) / 100
    extrusionWidth := palette.TowerExtrusionWidth
    extrusionMultiplier := palette.TowerExtrusionMultiplier / 100

    // tower layers must have at least this much extrusion to be able to fit pings!
    minLayerExtrusion := palette.GetPingExtrusion() / float32(maxDensity) // mm
    minLayerVolume := float64(filamentLengthToVolume(minLayerExtrusion)) // mm3

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
    footprintArea := float64(0)
    for layer, transitions := range preflight.transitionsByLayer {
        layerPurgeLength := float32(0) // mm
        for _, transition := range transitions {
            layerPurgeLength += transition.PurgeLength / extrusionMultiplier
        }
        layerPurgeVolume := float64(filamentLengthToVolume(layerPurgeLength)) // mm3
        // adjust for max density
        layerPurgeVolume /= maxDensity
        // raise the volume slightly to account for errors in total toolpath extrusion
        layerPurgeVolume *= 1.05
        // ensure the layer has room for at least one ping
        layerPurgeVolume = math.Max(layerPurgeVolume, minLayerVolume)

        layerFootprintArea := float32(layerPurgeVolume) / layerThicknesses[layer]
        layerFootprintAreas[layer] = layerFootprintArea
        footprintArea = math.Max(footprintArea, float64(layerFootprintArea))
    }
    if footprintArea == 0 {
        // no dense layers == no Palette processing
        return tower, false
    }

    // 7. finalize the tower dimensions
    //    - try and maintain the current aspect ratio

    towerWidth := float64(palette.TowerSize[0])
    towerHeight := float64(palette.TowerSize[1])
    squareLength := math.Sqrt(footprintArea)
    aspectRatio := 1 / math.SqrtPhi // default to the golden ratio
    if towerWidth > 0 && towerHeight > 0 {
        // prefer to use the provided aspect ratio, but not necessarily size
        aspectRatio = math.Sqrt(towerWidth / towerHeight)
    }
    towerWidth = squareLength / aspectRatio
    towerHeight = squareLength * aspectRatio
    towerHalfHeight := float32(towerWidth) / 2
    towerHalfWidth := float32(towerHeight) / 2

    // 8. determine the density of each layer
    //    - minimum and maximum tower density, minimum first layer density
    //    - ratio of required footprint area for this layer to overall footprint of the tower
    layerDensities := make([]float32, totalLayers)
    for layer := 0; layer < totalLayers; layer++ {
        layerFootprintArea := layerFootprintAreas[layer]
        var density float64
        if layerFootprintArea > 0 {
            // dense layer
            density = float64(layerFootprintArea) / footprintArea
        } else {
            // sparse layer -- ensure enough density to fit a ping
            layerThickness := layerThicknesses[layer]
            fullLayerVolume := footprintArea * float64(layerThickness)
            density = minLayerVolume / fullLayerVolume
            // if perimeters will be added, account for it
            if density <= TowerPerimeterThreshold {
                perimeterLength := float32((4 * towerWidth) + (4 * towerHeight)) - (8 * extrusionWidth)
                perimeterExtrusion := getExtrusionLength(extrusionWidth, layerThickness, perimeterLength) * extrusionMultiplier
                perimeterVolume := float64(filamentLengthToVolume(perimeterExtrusion))
                doubleEW := 2 * float64(extrusionWidth)
                infillFootprint := footprintArea - (doubleEW * towerWidth) - (doubleEW * (towerHeight - doubleEW))
                fullInfillVolume := infillFootprint * float64(layerThickness)
                requiredInfillVolume := (minLayerVolume - perimeterVolume) * 1.1
                density = requiredInfillVolume / fullInfillVolume
            }
        }
        // adjust for user-defined density limits
        if layer == 0 {
            density = math.Max(density, minFirstLayerDensity)
        } else {
            density = math.Max(density, minDensity)
        }
        density = math.Min(density, maxDensity)
        layerDensities[layer] = float32(density)
    }

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

    // 10. determine number of first-layer brims needed
    if palette.RaftLayers == 0 {
        firstTransitionTotalE := preflight.transitions[0].TotalExtrusion
        firstTransitionTotalE += minLayerExtrusion * float32(preflight.transitions[0].Layer)
        firstTransitionTotalE += preflight.transitions[0].TransitionLength * (palette.TransitionTarget / 100)
        firstLayerThickness := tower.Layers[0].Thickness
        minFirstSpliceLength := palette.GetFirstSpliceMinLength()
        perimeterLength := (towerHalfWidth * 4) + (towerHalfHeight * 4) + (palette.TowerExtrusionWidth * 8)
        for firstTransitionTotalE < minFirstSpliceLength {
            tower.BrimCount++
            brimExtrusion := getExtrusionLength(extrusionWidth, firstLayerThickness, perimeterLength) * extrusionMultiplier
            firstTransitionTotalE += brimExtrusion
            tower.BrimExtrusion += brimExtrusion
        }
        if tower.BrimCount < palette.TowerMinBrims {
            tower.BrimCount = palette.TowerMinBrims
        }
    }

    return tower, true
}

func (t *Tower) layerNeedsPerimeters(layer int, density float32) bool {
    if layer < t.Palette.RaftLayers {
        // raft layers never get perimeters
        return false
    }
    if layer == t.Palette.RaftLayers &&
        (t.Palette.TowerFirstLayerPerimeters || t.BrimCount > 0) {
        // first non-raft layer -- force a perimeter if desired by user or if using brims
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

    t.CurrentLayerPaths = make([]AnnotatedCommand, 0)
    t.CurrentLayerExtrusion = 0
    currentXMin := xMin
    currentXMax := xMax
    currentYMin := yMin
    currentYMax := yMax

    layerThickness := t.Layers[layer].Thickness
    density := t.Layers[layer].Density
    extrusionWidth := t.Palette.TowerExtrusionWidth
    extrusionMultiplier :=  t.Palette.TowerExtrusionMultiplier / 100
    if layer < t.Palette.RaftLayers {
        extrusionWidth = t.Palette.RaftExtrusionWidth
    }
    addPerimeters := t.layerNeedsPerimeters(layer, density)

    // create perimeters

    if addPerimeters {
        perimeterCount := TowerPerimeterCount
        if layer == 0 && t.BrimCount > 0 {
            perimeterCount += t.BrimCount
            inflation := extrusionWidth * float32(t.BrimCount)
            currentXMin -= inflation
            currentYMin -= inflation
            currentXMax += inflation
            currentYMax += inflation
        }
        for i := 0; i < perimeterCount; i++ {
            // travel to southeast corner
            nextX, nextY := currentXMax, currentYMin
            travel := AnnotatedCommand{
                gcode: gcode.Command{
                    Command: "G1",
                    Params:  map[string]float32{
                        "x": nextX,
                        "y": nextY,
                    },
                },
            }
            // for next extrusion line's length
            fromX, fromY := nextX, nextY

            // extrude to northeast corner
            nextX, nextY = currentXMax, currentYMax
            lineLength := getLineLength(fromX, fromY, nextX, nextY) // mm
            deltaE := getExtrusionLength(extrusionWidth, layerThickness, lineLength) * extrusionMultiplier
            t.CurrentLayerExtrusion += deltaE
            extrudeNorth := AnnotatedCommand{
                gcode: gcode.Command{
                    Command: "G1",
                    Params:  map[string]float32{
                        "x": nextX,
                        "y": nextY,
                    },
                },
                extrusion: deltaE,
            }
            fromX, fromY = nextX, nextY

            // extrude to northwest corner
            nextX, nextY = currentXMin, currentYMax
            lineLength = getLineLength(fromX, fromY, nextX, nextY) // mm
            deltaE = getExtrusionLength(extrusionWidth, layerThickness, lineLength) * extrusionMultiplier
            t.CurrentLayerExtrusion += deltaE
            extrudeWest := AnnotatedCommand{
                gcode: gcode.Command{
                    Command: "G1",
                    Params:  map[string]float32{
                        "x": nextX,
                        "y": nextY,
                    },
                },
                extrusion: deltaE,
            }
            fromX, fromY = nextX, nextY

            // extrude to southwest corner
            nextX, nextY = currentXMin, currentYMin
            lineLength = getLineLength(fromX, fromY, nextX, nextY) // mm
            deltaE = getExtrusionLength(extrusionWidth, layerThickness, lineLength) * extrusionMultiplier
            t.CurrentLayerExtrusion += deltaE
            extrudeSouth := AnnotatedCommand{
                gcode: gcode.Command{
                    Command: "G1",
                    Params:  map[string]float32{
                        "x": nextX,
                        "y": nextY,
                    },
                },
                extrusion: deltaE,
            }
            fromX, fromY = nextX, nextY

            // extrude to southeast corner
            nextX, nextY = currentXMax, currentYMin
            lineLength = getLineLength(fromX, fromY, nextX, nextY) // mm
            deltaE = getExtrusionLength(extrusionWidth, layerThickness, lineLength) * extrusionMultiplier
            t.CurrentLayerExtrusion += deltaE
            extrudeEast := AnnotatedCommand{
                gcode: gcode.Command{
                    Command: "G1",
                    Params:  map[string]float32{
                        "x": nextX,
                        "y": nextY,
                    },
                },
                extrusion: deltaE,
            }
            fromX, fromY = nextX, nextY

            t.CurrentLayerPaths = append(t.CurrentLayerPaths,
                travel,
                extrudeNorth,
                extrudeWest,
                extrudeSouth,
                extrudeEast,
            )
            // step inward by 1 extrusion width
            currentXMin += extrusionWidth
            currentYMin += extrusionWidth
            currentXMax -= extrusionWidth
            currentYMax -= extrusionWidth
        }
        // step outward slightly to produce an infill-perimeter overlap
        overlap := extrusionWidth * (t.Palette.InfillPerimeterOverlap / 100)
        currentXMin -= overlap
        currentYMin -= overlap
        currentXMax += overlap
        currentYMax += overlap
    }

    // create infill

    stride := (1 / t.Layers[layer].Density) * extrusionWidth
    if layer < t.Palette.RaftLayers {
        stride = t.Palette.RaftStride
    }
    axisAlignedStride := float32(math.Sqrt(float64(stride * stride * 2)))

    firstLine := true
    var lastX, lastY float32 // populated once firstLine == false
    needsMoreLines := true
    xBoundReached := false // once reached, step southwest vertex north instead of west
    yBoundReached := false // once reached, step northeast vertex west instead of north
    printSouthwest := true // direction of the next extrusion line ("back" or "forth")
    reverseLayer := layer % 2 == 1 // mirror every other layer's X coordinates

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

        if reverseLayer {
            // mirror layer's X values
            x1 = currentXMax + currentXMin - x1
            x2 = currentXMax + currentXMin - x2
        }

        // travel to (x1, y1)
        travel := AnnotatedCommand{
            gcode: gcode.Command{
                Command: "G1",
                Params:  map[string]float32{
                    "x": x1,
                    "y": y1,
                },
            },
        }
        if !firstLine && layer < t.Palette.RaftLayers {
            // raft layers should have continuous extrusion after the initial travel
            lineLength := getLineLength(lastX, lastY, x1, y1) // mm
            travel.extrusion = getExtrusionLength(extrusionWidth, layerThickness, lineLength) * extrusionMultiplier
            t.CurrentLayerExtrusion += travel.extrusion
        }
        firstLine = false
        lastX, lastY = x1, y1

        // extrude to (x2, y2)
        extrude := AnnotatedCommand{
            gcode: gcode.Command{
                Command: "G1",
                Params:  map[string]float32{
                    "x": x2,
                    "y": y2,
                },
            },
        }
        lineLength := getLineLength(lastX, lastY, x2, y2) // mm
        extrude.extrusion = getExtrusionLength(extrusionWidth, layerThickness, lineLength) * extrusionMultiplier
        t.CurrentLayerExtrusion += extrude.extrusion
        lastX, lastY = x2, y2

        t.CurrentLayerPaths = append(t.CurrentLayerPaths, travel, extrude)

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
    annotatedTravel := t.CurrentLayerPaths[t.CurrentLayerCommandIndex]
    travel := annotatedTravel.gcode
    if annotatedTravel.extrusion > 0 {
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
    if state.CurrentlyPinging {
        sequence += t.checkTowerPingEnd(state, true)
    }
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
    annotatedCommand := t.CurrentLayerPaths[t.CurrentLayerCommandIndex]
    command := annotatedCommand.gcode
    commandExtrusion := annotatedCommand.extrusion

    // when printing a segment, all commands use the print feedrate
    // so as not to alternate feedrates constantly
    command.Params["f"] = printFeedrate
    if commandExtrusion > 0 {
        // extrusion command
        if state.E.RelativeExtrusion {
            command.Params["e"] = commandExtrusion
        } else {
            command.Params["e"] = state.E.CurrentExtrusionValue + commandExtrusion
        }
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

func (t *Tower) isAccessoryPingStartConditionMet(state *State, segmentExtrusionSoFar, totalSegmentExtrusion float32) bool {
    // have we extruded enough since the last ping to warrant starting the next one?
    if state.E.TotalExtrusion < state.NextPingStart {
        return false
    }

    // if we start a ping now, can we fit the required extrusion in the remainder of the segment?
    if segmentExtrusionSoFar + (state.PingExtrusion * 0.9) >= totalSegmentExtrusion {
        return false
    }

    // do we want to avoid starting a ping right at the beginning of the segment?
    return (totalSegmentExtrusion < (state.PingExtrusion * 2.2)) || // small (sparse) tower segments can start a ping right away
        (segmentExtrusionSoFar >= state.PingExtrusion) // larger (dense) tower segments must wait ~20 mm to start a ping
}

func (t *Tower) checkTowerPingStart(state *State, segmentExtrusionSoFar, totalSegmentExtrusion float32) string {
    totalExtrusion := state.E.TotalExtrusion
    sequence := ""
    if state.Palette.ConnectedMode {
        if totalExtrusion >= state.NextPingStart {
            // connected pings
            state.MSF.AddPing(totalExtrusion)
            state.NextPingStart = totalExtrusion + PingMinSpacing
            sequence += fmt.Sprintf("; Ping %d%s", len(state.MSF.PingList) + 1, EOL)
            state.MSF.AddPing(totalExtrusion)
            sequence += "G4 P0" + EOL
            sequence += state.MSF.GetConnectedPingLine()
        }
    } else {
        if t.isAccessoryPingStartConditionMet(state, segmentExtrusionSoFar, totalSegmentExtrusion) {
            // start the accessory ping sequence
            sequence += fmt.Sprintf("; Ping %d pause 1%s", len(state.MSF.PingList) + 1, EOL)
            sequence += getTowerPause(Ping1PauseLength, state)
            state.CurrentPingStart = totalExtrusion
            state.NextPingStart = totalExtrusion + PingMinSpacing
            state.CurrentlyPinging = true
        }
    }
    return sequence
}

func (t *Tower) checkTowerPingEnd(state *State, force bool) string {
    totalExtrusion := state.E.TotalExtrusion
    nextPingEnd := state.CurrentPingStart + state.PingExtrusion
    finish := force || totalExtrusion >= nextPingEnd
    if !finish {
        // tower segment is not finished, and we haven't extruded PingExtrusion yet,
        // but we may be better off finishing the ping anyway
        if t.CurrentLayerCommandIndex + 1 < len(t.CurrentLayerPaths) {
            nextPathExtrusion := t.CurrentLayerPaths[t.CurrentLayerCommandIndex + 1].extrusion
            if math.Abs(float64(nextPingEnd + 0.5 - totalExtrusion)) <
                math.Abs(float64(totalExtrusion + nextPathExtrusion - 0.5)) {
                // the next path would put us further from PingExtrusion (in absolute value)
                // than we currently are -- finish the ping now to increase chance of detection
                finish = true
            }
        }
    }
    sequence := ""
    if finish {
        // finish the accessory ping sequence
        sequence += fmt.Sprintf("; Ping %d pause 2%s", len(state.MSF.PingList) + 1, EOL)
        sequence += getTowerPause(Ping2PauseLength, state)
        state.MSF.AddPingWithExtrusion(state.CurrentPingStart, totalExtrusion - state.CurrentPingStart)
        state.LastPingStart = state.CurrentPingStart
        state.NextPingStart = state.CurrentPingStart + PingMinSpacing
        state.CurrentlyPinging = false
    }
    return sequence
}

func (t *Tower) checkTowerPingActions(state *State, segmentExtrusionSoFar, totalSegmentExtrusion float32) string {
    if state.CurrentlyPinging {
        return t.checkTowerPingEnd(state, false)
    } else {
        return t.checkTowerPingStart(state, segmentExtrusionSoFar, totalSegmentExtrusion)
    }
}

func (t *Tower) getNextDenseSegmentPaths(state *State) string {
    transitionInfo := t.GetCurrentTransitionInfo()
    requiredPurge := transitionInfo.PurgeLength
    if t.CurrentLayerIndex == 0 && t.CurrentLayerTransitionIndex == 0 {
        requiredPurge += t.BrimExtrusion
    }
    // if this layer is denser than expected, distribute the extra extrusion
    // equally between transitions on this layer
    if t.CurrentLayerIsDense() {
        transitionCount := 0
        totalRequiredPurge := float32(0)
        for _, transition := range t.Layers[t.CurrentLayerIndex].Transitions {
            transitionCount++
            totalRequiredPurge += transition.PurgeLength
        }
        if t.CurrentLayerExtrusion > totalRequiredPurge {
            extra := t.CurrentLayerExtrusion - totalRequiredPurge
            requiredPurge += extra / float32(transitionCount)
        }
    }
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
    for (totalPurge < requiredPurge || t.CurrentLayerTransitionIndex == thisLayerTransitions - 1) &&
      t.CurrentLayerCommandIndex < len(t.CurrentLayerPaths) {
        sequence += t.checkTowerPingActions(state, totalPurge, requiredPurge)
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
    totalPurge := float32(0)
    for t.CurrentLayerCommandIndex < len(t.CurrentLayerPaths) {
        sequence += t.checkTowerPingActions(state, totalPurge, t.CurrentLayerExtrusion)
        commandString, commandExtrusion := t.getNextPath(state, printFeedrate)
        sequence += commandString
        totalPurge += commandExtrusion
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
