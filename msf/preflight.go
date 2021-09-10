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

    // initialize state
    state := NewState(palette)
    // account for a firmware purge (not part of G-code) once
    state.E.TotalExtrusion += palette.FirmwarePurge

    pingExtrusionMM := palette.GetPingExtrusion()

    err := gcode.ReadByLine(inpath, func(line gcode.Command) error {
        state.E.TrackInstruction(line)
        state.XYZF.TrackInstruction(line)
        if line.IsLinearMove() {
            if x, ok := line.Params["x"]; ok {
                if x < results.boundingBox.min[0] { results.boundingBox.min[0] = x }
                if x > results.boundingBox.max[0] { results.boundingBox.max[0] = x }
            }
            if y, ok := line.Params["y"]; ok {
                if y < results.boundingBox.min[1] { results.boundingBox.min[1] = y }
                if y > results.boundingBox.max[1] { results.boundingBox.max[1] = y }
            }
            if z, ok := line.Params["z"]; ok {
                if z < results.boundingBox.min[2] { results.boundingBox.min[2] = z }
                if z > results.boundingBox.max[2] { results.boundingBox.max[2] = z }
            }
            if state.OnWipeTower {
                // check for ping actions
                if state.CurrentlyPinging {
                    // currentlyPinging == true implies accessory mode
                    if state.E.TotalExtrusion >= state.CurrentPingStart + pingExtrusionMM {
                        // commit to the accessory ping sequence
                        results.pingStarts = append(results.pingStarts, state.CurrentPingStart)
                        state.LastPingStart = state.CurrentPingStart
                        state.CurrentlyPinging = false
                    }
                } else if state.E.TotalExtrusion >= state.LastPingStart + PingMinSpacing {
                    // attempt to start a ping sequence
                    //  - connected pings: guaranteed to finish
                    //  - accessory pings: may be "cancelled" if near the end of the transition
                    if palette.ConnectedMode {
                        results.pingStarts = append(results.pingStarts, state.E.TotalExtrusion)
                        state.LastPingStart = state.E.TotalExtrusion
                    } else {
                        state.CurrentPingStart = state.E.TotalExtrusion
                        state.CurrentlyPinging = true
                    }
                }
            }
        } else if len(line.Command) > 1 && line.Command[0] == 'T' {
            tool, err := strconv.ParseInt(line.Command[1:], 10, 32)
            if err != nil {
                return err
            }
            if state.FirstToolChange {
                state.FirstToolChange = false
            } else {
                state.CurrentTool = int(tool)
                state.CurrentlyTransitioning = true
                results.drivesUsed[state.CurrentTool] = true
                results.transitionStarts = append(results.transitionStarts, state.E.TotalExtrusion)
            }
        } else if palette.TransitionMethod == TransitionTower &&
            strings.HasPrefix(line.Comment, "TYPE:") {
            startingWipeTower := line.Comment == "TYPE:Wipe tower"
            if !state.OnWipeTower && startingWipeTower {
                // start of the actual transition being printed
                // todo: any logic needed?
            } else if state.OnWipeTower && !startingWipeTower {
                // end of the actual transition being printed
                if state.CurrentlyTransitioning {
                    state.CurrentlyTransitioning = false
                    // if we're in the middle of an accessory ping, cancel it -- too late to finish
                    state.CurrentlyPinging = false
                }
            }
            state.OnWipeTower = startingWipeTower
        }

        return nil
    })
    if results.boundingBox.min[2] > 0 {
        results.boundingBox.min[2] = 0
    }
    return results, err
}