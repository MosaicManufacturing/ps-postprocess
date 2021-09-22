package msf

import (
    "../gcode"
    "math"
)

type TowerLayer struct {
    TopZ float32 // Z value of the top of this layer
    Thickness float32 // height of extruded paths
    Density float32 // 0..1
}

type Tower struct {
    // constant or pre-calculated
    Palette *Palette
    BoundingBox gcode.BoundingBox
    Layers []TowerLayer

    // for use during output
    CurrentLayerPaths []gcode.Command // feedrates, raw strings, real E values not included yet
    CurrentLayerTransitions int // total transitions on this layer
    CurrentLayerTransitionIndex int // current transition on this layer
    CurrentTransitionCommandIndex int // index into CurrentLayerPaths
}

func filamentVolume(filamentDiameter, length float32) float32 {
    // V = Pi * (r^2) * h
    radius := filamentDiameter / 2
    return math.Pi * radius * radius * length
}

func GenerateTower(palette *Palette, preflight *msfPreflight) (Tower, bool) {
    totalLayers := len(preflight.transitionsByLayer)
    tower := Tower{
        Palette:                       palette,
        BoundingBox:                   gcode.NewBoundingBox(),
        Layers:                        nil,
        CurrentLayerPaths:             nil,
        CurrentLayerTransitions:       0,
        CurrentLayerTransitionIndex:   0,
        CurrentTransitionCommandIndex: 0,
    }

    // tower must have at least this many mm3 of extrusion to be able to fit pings!
    minLayerVolume := filamentVolume(1.75, palette.GetPingExtrusion())

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
    for i := 0; i < totalLayers; i++ {
        layerThicknesses[i] = preflight.layerThicknesses[i]
        layerTopZs[i] = preflight.layerTopZs[i]
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
        layerPurgeVolume := float32(0) // mm3
        for _, transition := range transitions {
            purgeVolume := filamentVolume(1.75, transition.purgeLength)
            // adjust volume based on extrusion multiplier
            purgeVolume /= palette.TowerExtrusionMultiplier
            layerPurgeVolume += purgeVolume
        }
        // adjust for max density
        layerPurgeVolume /= palette.TowerMaxDensity
        // raise the volume slightly to account for errors in total toolpath extrusion
        layerPurgeVolume *= 1.05
        // ensure the layer has room for at least one ping
        layerPurgeVolume = float32(math.Max(float64(layerPurgeVolume), float64(minLayerVolume / palette.TowerMaxDensity)))

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
            density = math.Max(density, float64(palette.TowerMinFirstLayerDensity))
        } else {
            density = math.Max(density, float64(palette.TowerMinDensity))
        }
        density = math.Min(density, float64(palette.TowerMaxDensity))
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
    towerHalfWidth := float32(math.Round(squareLength / aspectRatio)) / 2
    towerHalfHeight := float32(math.Round(squareLength * aspectRatio)) / 2

    // 9. store everything relevant

    tower.Layers = make([]TowerLayer, totalLayers)
    for layer := 0; layer < totalLayers; layer++ {
        tower.Layers[layer] = TowerLayer{
            TopZ:      layerTopZs[layer],
            Thickness: layerThicknesses[layer],
            Density:   layerDensities[layer],
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
        xMin -= t.Palette.RaftInflation
        yMin -= t.Palette.RaftInflation
        xMax += t.Palette.RaftInflation
        yMax += t.Palette.RaftInflation
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

        t.CurrentLayerPaths = append(t.CurrentLayerPaths,
            // travel to (x1, y1)
            gcode.Command{
                Command: "G1",
                Params:  map[string]float32{
                    "x": x1,
                    "y": y1,
                    "f": 0,
                },
            },
            // extrude to (x2, y2)
            gcode.Command{
                Command: "G1",
                Params:  map[string]float32{
                    "x": x2,
                    "y": y2,
                    "e": 0,
                    "f": 0,
                },
            },
        )

        if neX - currentXMin < axisAlignedStride && currentYMax - swY < axisAlignedStride {
            // layer has been fully rasterized
            needsMoreLines = false
        } else {
            // reverse the direction of the next line
            printSouthwest = !printSouthwest
        }
    }

}
