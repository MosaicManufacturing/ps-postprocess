package msf

import (
    "../gcode"
	"fmt"
    "math"
)

func getSideTransitionStartPosition(state *State) (x, y float32) {
    if state.Palette.SideTransitionJog {
        x = state.XYZF.CurrentX
        y = state.XYZF.CurrentY
        switch state.Palette.SideTransitionEdge {
        case gcode.North:
            y = state.Palette.PrintBedMaxY + state.Palette.SideTransitionEdgeOffset
        case gcode.South:
            y = state.Palette.PrintBedMinY - state.Palette.SideTransitionEdgeOffset
        case gcode.West:
            x = state.Palette.PrintBedMinX - state.Palette.SideTransitionEdgeOffset
        case gcode.East:
            x = state.Palette.PrintBedMaxX + state.Palette.SideTransitionEdgeOffset
        }
    } else {
        // side transition in place
        x = state.Palette.SideTransitionX
        y = state.Palette.SideTransitionY
    }
    return
}

func moveToSideTransition(transitionLength float32, state *State) (string, error) {
    startX, startY := getSideTransitionStartPosition(state)

    if state.Palette.PreSideTransitionScript != nil {
        // user script instead of built-in logic
        transitionIdx := len(state.MSF.SpliceList) - 1
        locals := state.Locals.Prepare(state.CurrentTool, map[string]float64{
            "layer": float64(state.CurrentLayer),
            "currentPrintTemperature": float64(state.Temperature.Extruder),
            "currentBedTemperature": float64(state.Temperature.Bed),
            "currentX": float64(state.XYZF.CurrentX),
            "currentY": float64(state.XYZF.CurrentY),
            "currentZ": float64(state.XYZF.CurrentZ),
            "nextX": float64(state.TransitionNextPositions[transitionIdx][0]),
            "nextY": float64(state.TransitionNextPositions[transitionIdx][1]),
            "nextZ": float64(state.TransitionNextPositions[transitionIdx][2]),
            "transitionLength": float64(transitionLength),
        })
        return evaluateScript(state.Palette.PreSideTransitionScript, locals, state)
    }

    sequence := "; move to side transition" + EOL
    travel := gcode.Command{
        Raw:     fmt.Sprintf("G1 X%.3f Y%.3f F%.1f", startX, startY, state.Palette.TravelSpeedXY),
        Command: "G1",
        Params:  map[string]float32{
            "x": startX,
            "y": startY,
            "f": state.Palette.TravelSpeedXY,
        },
        Flags: map[string]bool{},
    }
    state.TimeEstimate += estimateMoveTime(state.XYZF.CurrentX, state.XYZF.CurrentY, startX, startY, state.Palette.TravelSpeedXY)
    state.XYZF.TrackInstruction(travel)
    sequence += travel.Raw + EOL
    return sequence, nil
}

func leaveSideTransition(transitionLength float32, state *State) (string, error) {
    if state.Palette.PostSideTransitionScript != nil {
        // user script instead of built-in logic
        transitionIdx := len(state.MSF.SpliceList) - 1
        locals := state.Locals.Prepare(state.CurrentTool, map[string]float64{
            "layer": float64(state.CurrentLayer),
            "currentPrintTemperature": float64(state.Temperature.Extruder),
            "currentBedTemperature": float64(state.Temperature.Bed),
            "currentX": float64(state.XYZF.CurrentX),
            "currentY": float64(state.XYZF.CurrentY),
            "currentZ": float64(state.XYZF.CurrentZ),
            "nextX": float64(state.TransitionNextPositions[transitionIdx][0]),
            "nextY": float64(state.TransitionNextPositions[transitionIdx][1]),
            "nextZ": float64(state.TransitionNextPositions[transitionIdx][2]),
            "transitionLength": float64(transitionLength),
        })
        return evaluateScript(state.Palette.PostSideTransitionScript, locals, state)
    }

    return "; leave side transition" + EOL, nil
}

func checkSideTransitionPings(state *State) (bool, string, float32) {
    if state.E.TotalExtrusion < state.LastPingStart + PingMinSpacing {
        // not time for a ping yet
        return false, "", 0
    }

    if state.Palette.ConnectedMode {
        // connected pings can happen anywhere during the transition,
        // even at the very end
        state.MSF.AddPing(state.E.TotalExtrusion)
        state.LastPingStart = state.E.TotalExtrusion
        sequence := fmt.Sprintf("; Ping %d%s", len(state.MSF.PingList), EOL)
        sequence += "G4 P0" + EOL
        sequence += state.MSF.GetConnectedPingLine()
        return true, sequence, 0
    }

    var sequence string
    var extrusion float32
    if state.Palette.SideTransitionJog {
        sequence, extrusion = doSideTransitionOnEdgeAccessoryPing(state)
    } else {
        sequence, extrusion = doSideTransitionInPlaceAccessoryPing(state)
    }
    return true, sequence, extrusion
}

func sideTransitionInPlace(transitionLength float32, state *State) (string, error) {
    feedrate := state.Palette.SideTransitionPurgeSpeed * 60
    transitionSoFar := float32(0)
    sequence, err := moveToSideTransition(transitionLength, state)
    if err != nil {
        return sequence, err
    }

    for transitionSoFar < transitionLength {
        if doPing, pingSequence, pingExtrusion := checkSideTransitionPings(state); doPing {
            transitionSoFar += pingExtrusion
            sequence += pingSequence
        }
        nextPurgeExtrusion := float32(math.Min(10, float64(transitionLength - transitionSoFar)))
        nextE := nextPurgeExtrusion
        if !state.E.RelativeExtrusion {
            nextE = state.E.CurrentExtrusionValue + nextPurgeExtrusion
        }
        purge := gcode.Command{
            Raw:     fmt.Sprintf("G1 E%.5f F%.1f", nextE, feedrate),
            Command: "G1",
            Params:  map[string]float32{
                "e": nextE,
                "f": feedrate,
            },
            Flags: map[string]bool{},
        }
        state.TimeEstimate += estimatePurgeTime(nextPurgeExtrusion, feedrate)
        state.E.TrackInstruction(purge)
        sequence += purge.Raw + EOL
        transitionSoFar += nextPurgeExtrusion
    }

    if doPing, pingSequence, pingExtrusion := checkSideTransitionPings(state); doPing {
        transitionSoFar += pingExtrusion
        sequence += pingSequence
    }

    leave, err := leaveSideTransition(transitionLength, state)
    if err != nil {
        return sequence, err
    }
    sequence += leave
    return sequence, nil
}

func sideTransitionOnEdge(transitionLength float32, state *State) (string, error) {
    eFeedrate := state.Palette.SideTransitionPurgeSpeed * 60
    xyFeedrate := state.Palette.SideTransitionMoveSpeed * 60
    transitionSoFar := float32(0)

    // determine next purge direction
    var nextPurgeDirection gcode.Direction
    if state.Palette.SideTransitionEdge == gcode.North || state.Palette.SideTransitionEdge == gcode.South {
        if state.Palette.PrintBedMaxX - state.XYZF.CurrentX >= state.XYZF.CurrentX - state.Palette.PrintBedMinX {
            nextPurgeDirection = gcode.East
        } else {
            nextPurgeDirection = gcode.West
        }
    } else {
        if state.Palette.PrintBedMaxY - state.XYZF.CurrentY >= state.XYZF.CurrentY - state.Palette.PrintBedMinY {
            nextPurgeDirection = gcode.North
        } else {
            nextPurgeDirection = gcode.South
        }
    }
    nextX := state.XYZF.CurrentX
    nextY := state.XYZF.CurrentY
    switch state.Palette.SideTransitionEdge {
    case gcode.North:
        nextY = state.Palette.PrintBedMaxY + state.Palette.SideTransitionEdgeOffset
    case gcode.South:
        nextY = state.Palette.PrintBedMinY - state.Palette.SideTransitionEdgeOffset
    case gcode.West:
        nextX = state.Palette.PrintBedMinX - state.Palette.SideTransitionEdgeOffset
    case gcode.East:
        nextX = state.Palette.PrintBedMaxX + state.Palette.SideTransitionEdgeOffset
    }

    // move to starting position
    sequence, err := moveToSideTransition(transitionLength, state)
    if err != nil {
        return sequence, err
    }

    dimensionOfInterest := state.Palette.PrintBedMaxX - state.Palette.PrintBedMinX
    if state.Palette.SideTransitionEdge == gcode.West || state.Palette.SideTransitionEdge == gcode.East {
        dimensionOfInterest = state.Palette.PrintBedMaxY - state.Palette.PrintBedMinY
    }
    edgeClearance := float32(15)
    if dimensionOfInterest < 50 {
        edgeClearance = 0
    } else if dimensionOfInterest < 80 {
        edgeClearance = 10
    }

    // purge until transition length is achieved
    for transitionSoFar < transitionLength {
        if doPing, pingSequence, pingExtrusion := checkSideTransitionPings(state); doPing {
            transitionSoFar += pingExtrusion
            sequence += pingSequence
        }
        switch nextPurgeDirection {
        case gcode.North:
            nextY = state.Palette.PrintBedMaxY - edgeClearance
        case gcode.South:
            nextY = state.Palette.PrintBedMinY + edgeClearance
        case gcode.West:
            nextX = state.Palette.PrintBedMinX + edgeClearance
        case gcode.East:
            nextX = state.Palette.PrintBedMaxX - edgeClearance
        }
        nextLineLength := getLineLength(state.XYZF.CurrentX, state.XYZF.CurrentY, nextX, nextY)
        nextPurgeExtrusion := nextLineLength * (eFeedrate / xyFeedrate)
        if transitionSoFar + nextPurgeExtrusion > transitionLength {
            t := (transitionLength - transitionSoFar) / nextPurgeExtrusion
            nextPurgeExtrusion = lerp(0, nextPurgeExtrusion, t)
            if nextPurgeDirection == gcode.North || nextPurgeDirection == gcode.South {
                nextY = lerp(state.XYZF.CurrentY, nextY, t)
            } else {
                nextX = lerp(state.XYZF.CurrentX, nextX, t)
            }
        }
        nextE := nextPurgeExtrusion
        if !state.E.RelativeExtrusion {
            nextE = state.E.CurrentExtrusionValue + nextPurgeExtrusion
        }
        purge := gcode.Command{
            Raw:     fmt.Sprintf("G1 X%.3f Y%.3f E%.5f F%.1f", nextX, nextY, nextE, xyFeedrate),
            Command: "G1",
            Params:  map[string]float32{
                "x": nextX,
                "y": nextY,
                "e": nextE,
                "f": xyFeedrate,
            },
            Flags: map[string]bool{},
        }
        state.TimeEstimate += estimateMoveTime(state.XYZF.CurrentX, state.XYZF.CurrentY, nextX, nextY, xyFeedrate)
        state.E.TrackInstruction(purge)
        state.XYZF.TrackInstruction(purge)
        sequence += purge.Raw + EOL
        transitionSoFar += nextPurgeExtrusion
        switch nextPurgeDirection {
        case gcode.North:
            nextPurgeDirection = gcode.South
        case gcode.South:
            nextPurgeDirection = gcode.North
        case gcode.West:
            nextPurgeDirection = gcode.East
        case gcode.East:
            nextPurgeDirection = gcode.West
        }
    }

    if doPing, pingSequence, pingExtrusion := checkSideTransitionPings(state); doPing {
        transitionSoFar += pingExtrusion
        sequence += pingSequence
    }

    leave, err := leaveSideTransition(transitionLength, state)
    if err != nil {
        return sequence, err
    }
    sequence += leave
    return sequence, nil
}

func sideTransitionCustom(transitionLength float32, state *State) (string, error) {
    sequence, err := moveToSideTransition(transitionLength, state)
    if err != nil {
        return sequence, err
    }
    transitionSoFar := float32(0)

    if doPing, pingSequence, pingExtrusion := checkSideTransitionPings(state); doPing {
        transitionSoFar += pingExtrusion
        sequence += pingSequence
    }

    locals := state.Locals.Prepare(state.CurrentTool, map[string]float64{
        "layer": float64(state.CurrentLayer),
        "currentPrintTemperature": float64(state.Temperature.Extruder),
        "currentBedTemperature": float64(state.Temperature.Bed),
        "currentX": float64(state.XYZF.CurrentX),
        "currentY": float64(state.XYZF.CurrentY),
        "currentZ": float64(state.XYZF.CurrentZ),
        "transitionLength": float64(transitionLength),
    })
    result, err := evaluateScript(state.Palette.SideTransitionScript, locals, state)
    if err != nil {
        return sequence, err
    }
    sequence += result

    transitionSoFar = transitionLength
    if doPing, pingSequence, pingExtrusion := checkSideTransitionPings(state); doPing {
        transitionSoFar += pingExtrusion
        sequence += pingSequence
    }

    leave, err := leaveSideTransition(transitionLength, state)
    if err != nil {
        return sequence, err
    }
    sequence += leave
    return sequence, nil
}

func sideTransition(transitionLength float32, state *State) (string, error) {
    sequence := ";TYPE:Side transition" + EOL
    if state.Palette.SideTransitionScript != nil {
        result, err := sideTransitionCustom(transitionLength, state)
        if err != nil {
            return sequence, err
        }
        sequence += result
    } else if state.Palette.SideTransitionJog {
        result, err := sideTransitionOnEdge(transitionLength, state)
        if err != nil {
            return sequence, err
        }
        sequence += result
    } else {
        result, err := sideTransitionInPlace(transitionLength, state)
        if err != nil {
            return sequence, err
        }
        sequence += result
    }
    return sequence, nil
}
