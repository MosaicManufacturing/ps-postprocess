package msf

import (
    "../gcode"
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
    towerBoundingBox bbox
}

func emptyBBox() bbox {
    return bbox{
        min: [3]float32{posInf, posInf, posInf},
        max: [3]float32{negInf, negInf, negInf},
    }
}

func (b *bbox) expandX(x float32) {
    if x < b.min[0] { b.min[0] = x }
    if x > b.max[0] { b.max[0] = x }
}

func (b *bbox) expandY(y float32) {
    if y < b.min[1] { b.min[1] = y }
    if y > b.max[1] { b.max[1] = y }
}

func (b *bbox) expandZ(z float32) {
    if z < b.min[2] { b.min[2] = z }
    if z > b.max[2] { b.max[2] = z }
}

func preflight(inpath string, palette *Palette) (msfPreflight, error) {
    results := msfPreflight{
        drivesUsed:       make([]bool, palette.GetInputCount()),
        transitionStarts: make([]float32, 0),
        pingStarts:       make([]float32, 0),
        boundingBox:      emptyBBox(),
        towerBoundingBox: emptyBBox(),
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
                results.boundingBox.expandX(x)
            }
            if y, ok := line.Params["y"]; ok {
                results.boundingBox.expandY(y)
            }
            if z, ok := line.Params["z"]; ok {
                results.boundingBox.expandZ(z)
            }
            if state.OnWipeTower {
                if _, ok := line.Params["e"]; ok {
                    // extrusion on wipe tower -- update bounding box
                    if state.E.CurrentRetraction == 0 {
                        if x, ok := line.Params["x"]; ok {
                            results.towerBoundingBox.expandX(x)
                        }
                        if y, ok := line.Params["y"]; ok {
                            results.towerBoundingBox.expandY(y)
                        }
                    }
                }
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