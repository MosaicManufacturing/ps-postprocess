package msf

import (
    "../gcode"
    "math"
    "strconv"
    "strings"
)

type bbox struct {
    min [3]float32
    max [3]float32
}

type msfPreflight struct {
    drivesUsed []bool
    transitionStarts []float32
    pingStarts []float32
    boundingBox bbox
}

func preflight(inpath string, palette *Palette) (msfPreflight, error) {
    results := msfPreflight{
        drivesUsed:       make([]bool, palette.GetInputCount()),
        transitionStarts: make([]float32, 0),
        pingStarts:       make([]float32, 0),
        boundingBox:      bbox{
            min: [3]float32{float32(math.Inf(1)), float32(math.Inf(1)), float32(math.Inf(1))},
            max: [3]float32{float32(math.Inf(-1)), float32(math.Inf(-1)), float32(math.Inf(-1))},
        },
    }

    firstToolChange := true // don't treat the first T command as a toolchange
    currentTool := 0
    currentX := float32(0)
    currentY := float32(0)
    currentZ := float32(0)
    eTracker := gcode.ExtrusionTracker{}
    // account for a firmware purge (not part of G-code) once
    eTracker.TotalExtrusion += palette.FirmwarePurge

    currentlyTransitioning := false
    onWipeTower := false
    pingExtrusionMM := float32(PingExtrusion)
    if palette.Type == TypeP1 {
        pingExtrusionMM = PingExtrusionCounts / palette.GetPulsesPerMM()
    }
    lastPingStart := float32(0)
    currentlyPinging := false
    currentPingStart := float32(0)

    err := gcode.ReadByLine(inpath, func(line gcode.Command) error {
        eTracker.TrackInstruction(line)
        if line.IsLinearMove() {
            if x, ok := line.Params["x"]; ok {
                currentX = x
                if x < results.boundingBox.min[0] { results.boundingBox.min[0] = x }
                if x > results.boundingBox.max[0] { results.boundingBox.max[0] = x }
            }
            if y, ok := line.Params["y"]; ok {
                currentY = y
                if y < results.boundingBox.min[1] { results.boundingBox.min[1] = y }
                if y > results.boundingBox.max[1] { results.boundingBox.max[1] = y }
            }
            if z, ok := line.Params["z"]; ok {
                currentZ = z
                if z < results.boundingBox.min[2] { results.boundingBox.min[2] = z }
                if z > results.boundingBox.max[2] { results.boundingBox.max[2] = z }
            }
            if onWipeTower {
                // check for ping actions
                if currentlyPinging {
                    // currentlyPinging == true implies accessory mode
                    if eTracker.TotalExtrusion >= currentPingStart + pingExtrusionMM {
                        // commit to the accessory ping sequence
                        results.pingStarts = append(results.pingStarts, currentPingStart)
                        lastPingStart = currentPingStart
                        currentlyPinging = false
                    }
                } else if eTracker.TotalExtrusion >= lastPingStart + PingMinSpacing {
                    // attempt to start a ping sequence
                    //  - connected pings: guaranteed to finish
                    //  - accessory pings: may be "cancelled" if near the end of the transition
                    if palette.ConnectedMode {
                        results.pingStarts = append(results.pingStarts, eTracker.TotalExtrusion)
                        lastPingStart = eTracker.TotalExtrusion
                    } else {
                        currentPingStart = eTracker.TotalExtrusion
                        currentlyPinging = true
                    }
                }
            }
        } else if len(line.Command) > 1 && line.Command[0] == 'T' {
            tool, err := strconv.ParseInt(line.Command[1:], 10, 32)
            if err != nil {
                return err
            }
            if firstToolChange {
                firstToolChange = false
            } else {
                currentTool = int(tool)
                currentlyTransitioning = true
                results.drivesUsed[currentTool] = true
                results.transitionStarts = append(results.transitionStarts, eTracker.TotalExtrusion)
            }
        } else if strings.HasPrefix(line.Comment, "TYPE:") {
            startingWipeTower := line.Comment == "TYPE:Wipe tower"
            if !onWipeTower && startingWipeTower {
                // start of the actual transition being printed
                // todo: any logic needed?
            } else if onWipeTower && !startingWipeTower {
                // end of the actual transition being printed
                if currentlyTransitioning {
                    currentlyTransitioning = false
                    // if we're in the middle of an accessory ping, cancel it -- too late to finish
                    currentlyPinging = false
                }
            }
            onWipeTower = startingWipeTower
        }

        return nil
    })
    if results.boundingBox.min[2] > 0 {
        results.boundingBox.min[2] = 0
    }
    return results, err
}