package msf

import (
    "fmt"
    "math"
    "mosaicmfg.com/ps-postprocess/gcode"
)

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
    if math.Abs(float64(state.XYZF.CurrentX - state.TowerBoundingBox.Min[0])) >
        math.Abs(float64(state.TowerBoundingBox.Max[0] - state.XYZF.CurrentX)) {
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
    if retractDistance := state.Palette.RetractDistance[state.CurrentTool]; retractDistance != 0 {
        retractFeedrate := state.Palette.RetractFeedrate[state.CurrentTool]
        sequence += getRetract(state, retractDistance, retractFeedrate)
    } else if state.Palette.UseFirmwareRetraction {
        sequence += getFirmwareRetract()
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
        sequence += getXYTravel(state, pauseX, pauseY, currentF, "")
    }
    if state.Palette.JogPauses {
        sequence += getTowerJogPause(durationMS, state)
    } else {
        sequence += getDwellPause(durationMS)
    }
    state.TimeEstimate += float32(durationMS / 1000)
    if state.Palette.PingOffTowerDistance > 0 {
        // move back onto the tower after pausing
        sequence += getXYTravel(state, currentX, currentY, currentF, "")
    }
    if restartDistance := state.Palette.RestartDistance[state.CurrentTool]; restartDistance != 0 {
        restartFeedrate := state.Palette.RestartFeedrate[state.CurrentTool]
        sequence += getRestart(state, restartDistance, restartFeedrate)
    } else if state.Palette.UseFirmwareRetraction {
        sequence += getFirmwareRestart()
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
    state.TimeEstimate += float32(Ping1PauseLength / 1000)

    // extrusion between pauses
    pingStartExtrusion := state.E.TotalExtrusion
    purgeLength := state.PingExtrusion
    sequence += getPurge(state, purgeLength, state.Palette.SideTransitionPurgeSpeed * 60)

    // second pause
    sequence += fmt.Sprintf("; Ping %d pause 2", len(state.MSF.PingList) + 1)
    if state.Palette.JogPauses {
        sequence += getSideTransitionInPlaceJogPause(Ping2PauseLength, state)
    } else {
        sequence += getDwellPause(Ping2PauseLength)
    }
    state.TimeEstimate += float32(Ping2PauseLength / 1000)
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
    state.TimeEstimate += float32(Ping1PauseLength / 1000)

    // extrusion between pauses
    pingStartExtrusion := state.E.TotalExtrusion
    currentX := state.XYZF.CurrentX
    currentY := state.XYZF.CurrentY
    nextX := currentX
    nextY := currentY
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
    purgeLength := state.PingExtrusion
    purgePerJog := purgeLength / (totalJogs * 2)

    for i := 0; i < totalJogs; i++ {
        // jog out
        sequence += getXYExtrusion(state, nextX, nextY, purgePerJog, feedrate)
        // jog back in
        sequence += getXYExtrusion(state, currentX, currentY, purgePerJog, feedrate)
    }

    // second pause
    sequence += fmt.Sprintf("; Ping %d pause 2", len(state.MSF.PingList) + 1)
    if state.Palette.JogPauses {
        sequence += getSideTransitionOnEdgeJogPause(Ping2PauseLength, jogPauseDirection, state)
    } else {
        sequence += getDwellPause(Ping2PauseLength)
    }
    state.TimeEstimate += float32(Ping2PauseLength / 1000)
    state.MSF.AddPingWithExtrusion(pingStartExtrusion, purgeLength)

    return sequence, purgeLength
}
