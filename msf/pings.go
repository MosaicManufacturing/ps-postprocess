package msf

import (
    "../gcode"
    "fmt"
    "math"
)

func getPingRetract(palette *Palette) (bool, string) {
    if palette.PingRetractDistance == 0 {
        return false, ""
    }
    return true, fmt.Sprintf("G1 E%.5f F%.1f", -palette.PingRetractDistance, palette.PingRetractFeedrate)
}

func getPingRestart(palette *Palette) (bool, string) {
    if palette.PingRestartDistance == 0 {
        return false, ""
    }
    return true, fmt.Sprintf("G1 E%.5f F%.1f", palette.PingRestartDistance, palette.PingRestartFeedrate)
}

func getDwellPause(durationMS int) string {
    sequence := ""
    for durationMS > 0 {
        if durationMS > 4000 {
            sequence += "G4 P4000" + EOL
            sequence += "G1" + EOL
            durationMS -= 4000
        } else {
            sequence += fmt.Sprintf("G4 P%d%s", durationMS, EOL)
            durationMS = 0
        }
    }
    return sequence
}

func getJogPause(durationMS int, direction gcode.Direction, state *State) string {
    const totalJogs = 5
    const feedrate = 10
    durationMM := float32(durationMS) / 1000

    x1 := state.XYZF.CurrentX
    y1 := state.XYZF.CurrentY
    x2 := x1
    y2 := y1
    switch direction {
    case gcode.North:
        y2 += durationMM
    case gcode.South:
        y2 -= durationMM
    case gcode.West:
        x2 -= durationMM
    case gcode.East:
        x2 += durationMM
    }

    sequence := ""
    for i := 0; i < totalJogs; i++ {
        sequence += fmt.Sprintf("G1 X%.3f Y%.3f F%d%s", x2, y2, feedrate, EOL)
        sequence += fmt.Sprintf("G1 X%.3f Y%.3f F%d%s", x1, y1, feedrate, EOL)
    }

    return sequence
}

func getTowerJogPauseDirection(state *State) gcode.Direction {
    // PrusaSlicer tower extrusions always run east-west, so either jog eastward or westward,
    // depending on which edge we're currently closer to
    if math.Abs(float64(state.XYZF.CurrentX - state.TowerBoundingBox.min[0])) >
        math.Abs(float64(state.TowerBoundingBox.max[0] - state.XYZF.CurrentX)) {
        // closer to west edge -- jog westward
        return gcode.West
    } else {
        // closer to east edge -- jog eastward
        return gcode.East
    }
}

func getTowerJogPause(durationMS int, state *State) string {
    direction := getTowerJogPauseDirection(state)
    return getJogPause(durationMS, direction, state)
}

func getTowerPause(durationMS int, state *State) string {
    sequence := ""
    if useRetract, retract := getPingRetract(state.Palette); useRetract {
        sequence += retract + EOL
    }
    currentF := state.XYZF.CurrentFeedrate
    currentX := state.XYZF.CurrentX
    currentY := state.XYZF.CurrentY
    pauseX := currentX
    pauseY := currentY
    if state.Palette.PingOffTowerDistance > 0 {
        // move off the tower before pausing
        direction := getTowerJogPauseDirection(state)
        if direction == gcode.West {
            pauseX -= state.Palette.PingOffTowerDistance
        } else {
            pauseX += state.Palette.PingOffTowerDistance
        }
        move := gcode.Command{
            Raw:     fmt.Sprintf("G1 X%.3f Y%.3f F%.1f", pauseX, pauseY, currentF),
            Command: "G1",
            Params:  map[string]float32{
                "x": pauseX,
                "y": pauseY,
                "f": currentF,
            },
            Flags: map[string]bool{},
        }
        state.XYZF.TrackInstruction(move)
        sequence += move.Raw + EOL
    }
    if state.Palette.JogPauses {
        sequence += getTowerJogPause(durationMS, state)
    } else {
        sequence += getDwellPause(durationMS)
    }
    if state.Palette.PingOffTowerDistance > 0 {
        // move back onto the tower after pausing
        move := gcode.Command{
            Raw:     fmt.Sprintf("G1 X%.3f Y%.3f F%.1f", currentX, currentY, currentF),
            Command: "G1",
            Params:  map[string]float32{
                "x": currentX,
                "y": currentY,
                "f": currentF,
            },
            Flags: map[string]bool{},
        }
        state.XYZF.TrackInstruction(move)
        sequence += move.Raw + EOL
    }
    if useRestart, restart := getPingRestart(state.Palette); useRestart {
        sequence += restart + EOL
    }
    return sequence
}

func getSideTransitionInPlaceJogPauseDirection(state *State) gcode.Direction {
    bedMiddleX := state.Palette.PrintBedMaxX - state.Palette.PrintBedMinX
    bedMiddleY := state.Palette.PrintBedMaxY - state.Palette.PrintBedMinY
    x := state.XYZF.CurrentX
    y := state.XYZF.CurrentY

    if x < state.Palette.PrintBedMinX {
        // off west edge
        if y <= state.Palette.PrintBedMinY || y >= state.Palette.PrintBedMaxY {
            // too close to north/south limits to go north or south
            return gcode.East

        }
        // go in the direction in which we're further away from the edge
        if y <= bedMiddleY {
            return gcode.North
        }
        return gcode.South
    }
    if x >= state.Palette.PrintBedMaxX {
        // off east edge
        if y <= state.Palette.PrintBedMinY || y >= state.Palette.PrintBedMaxY {
            // too close to north/south limits to go north or south
            return gcode.West
        }
        // go in the direction in which we're further away from the edge
        if y <= bedMiddleY {
            return gcode.North
        }
        return gcode.South
    }
    if y <= state.Palette.PrintBedMinY || y >= state.Palette.PrintBedMaxY {
        // off north or south edge
        // go in the direction in which we're further away from the edge
        if x <= bedMiddleX {
            return gcode.East
        }
        return gcode.West
    }
    // user specified a purge location over the print bed
    // (unlikely, but possible)
    if x < bedMiddleX {
        return gcode.West
    }
    if x > bedMiddleX {
        return gcode.East
    }
    if y >= bedMiddleY {
        return gcode.North
    }
    return gcode.South
}

func getSideTransitionInPlaceJogPause(durationMS int, state *State) string {
    direction := getSideTransitionInPlaceJogPauseDirection(state)
    return getJogPause(durationMS, direction, state)
}

func doSideTransitionInPlaceAccessoryPing(state *State) (string, float32) {
    // first pause
    sequence := fmt.Sprintf("; Ping %d pause 1", len(state.MSF.PingList) + 1)
    if state.Palette.JogPauses {
        sequence += getSideTransitionInPlaceJogPause(Ping1PauseLength, state)
    } else {
        sequence += getDwellPause(Ping1PauseLength)
    }

    // extrusion between pauses
    pingStartExtrusion := state.E.TotalExtrusion
    purgeLength := state.Palette.GetPingExtrusion()
    eValue := purgeLength
    if !state.E.RelativeExtrusion {
        eValue += state.E.CurrentExtrusionValue
    }
    purge := gcode.Command{
        Raw:     fmt.Sprintf("G1 E%.5f F%.1f", eValue, state.Palette.SideTransitionPurgeSpeed),
        Command: "G1",
        Params:  map[string]float32{
            "e": eValue,
            "f": state.Palette.SideTransitionPurgeSpeed,
        },
        Flags: map[string]bool{},
    }
    state.E.TrackInstruction(purge)
    sequence += purge.Raw + EOL

    // second pause
    sequence += fmt.Sprintf("; Ping %d pause 2", len(state.MSF.PingList) + 1)
    if state.Palette.JogPauses {
        sequence += getSideTransitionInPlaceJogPause(Ping2PauseLength, state)
    } else {
        sequence += getDwellPause(Ping2PauseLength)
    }
    state.MSF.AddPingWithExtrusion(pingStartExtrusion, purgeLength)

    return sequence, purgeLength
}

func getSideTransitionOnEdgeJogPauseDirection(state *State) gcode.Direction {
    if state.Palette.SideTransitionEdge == gcode.North ||
        state.Palette.SideTransitionEdge == gcode.South {
        if state.XYZF.CurrentX - state.Palette.PrintBedMinX >
            state.Palette.PrintBedMaxX - state.XYZF.CurrentX {
            return gcode.West
        }
        return gcode.East
    }
    if state.XYZF.CurrentY - state.Palette.PrintBedMinY >
        state.Palette.PrintBedMaxY - state.XYZF.CurrentY {
        return gcode.South
    }
    return gcode.North
}

func getSideTransitionOnEdgeJogPause(durationMS int, direction gcode.Direction, state *State) string {
    return getJogPause(durationMS, direction, state)
}

func doSideTransitionOnEdgeAccessoryPing(state *State) (string, float32) {
    jogPauseDirection := getSideTransitionOnEdgeJogPauseDirection(state)

    // first pause
    sequence := fmt.Sprintf("; Ping %d pause 1", len(state.MSF.PingList) + 1)
    if state.Palette.JogPauses {
        sequence += getSideTransitionOnEdgeJogPause(Ping1PauseLength, jogPauseDirection, state)
    } else {
        sequence += getDwellPause(Ping1PauseLength)
    }

    // extrusion between pauses
    pingStartExtrusion := state.E.TotalExtrusion
    nextX := state.XYZF.CurrentX
    nextY := state.XYZF.CurrentY
    switch jogPauseDirection {
    case gcode.North:
        if state.Palette.PrintBedMaxY - state.XYZF.CurrentY < 20 {
            nextY -= 5
        } else {
            nextY += 5
        }
    case gcode.South:
        if state.XYZF.CurrentY - state.Palette.PrintBedMinY < 20 {
            nextY += 5
        } else {
            nextY -= 5
        }
    case gcode.West:
        if state.XYZF.CurrentX - state.Palette.PrintBedMinX < 20 {
            nextX += 5
        } else {
            nextX -= 5
        }
    case gcode.East:
        if state.Palette.PrintBedMaxX - state.XYZF.CurrentX < 20 {
            nextX -= 5
        } else {
            nextX += 5
        }
    }

    const totalJogs = 5
    const feedrate = 600
    purgeLength := state.Palette.GetPingExtrusion()
    purgePerJog := purgeLength / (totalJogs * 2)

    for i := 0; i < totalJogs; i++ {
        eValue := purgePerJog
        if !state.E.RelativeExtrusion {
            eValue += state.E.CurrentExtrusionValue
        }
        // jog out
        purge := gcode.Command{
            Raw:     fmt.Sprintf("G1 X%.3f Y%.3f E%.5f F%d", nextX, nextY, eValue, feedrate),
            Command: "G1",
            Params:  map[string]float32{
                "x": nextX,
                "y": nextY,
                "e": eValue,
                "f": feedrate,
            },
            Flags: map[string]bool{},
        }
        state.E.TrackInstruction(purge)
        sequence += purge.Raw + EOL

        // jog back in
        if !state.E.RelativeExtrusion {
            eValue += purgePerJog
        }
        purge = gcode.Command{
            Raw:     fmt.Sprintf("G1 X%.3f Y%.3f E%.5f F%d", state.XYZF.CurrentX, state.XYZF.CurrentY, eValue, feedrate),
            Command: "G1",
            Params:  map[string]float32{
                "x": state.XYZF.CurrentX,
                "y": state.XYZF.CurrentY,
                "e": eValue,
                "f": feedrate,
            },
            Flags: map[string]bool{},
        }
        state.E.TrackInstruction(purge)
        sequence += purge.Raw + EOL
    }

    // second pause
    sequence += fmt.Sprintf("; Ping %d pause 2", len(state.MSF.PingList) + 1)
    if state.Palette.JogPauses {
        sequence += getSideTransitionOnEdgeJogPause(Ping2PauseLength, jogPauseDirection, state)
    } else {
        sequence += getDwellPause(Ping2PauseLength)
    }
    state.MSF.AddPingWithExtrusion(pingStartExtrusion, purgeLength)

    return sequence, purgeLength
}
