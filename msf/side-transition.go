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

func moveToSideTransition(state *State) string {
    startX, startY := getSideTransitionStartPosition(state)
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
    return sequence
}

func leaveSideTransition() string {
    return "; leave side transition" + EOL
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

func sideTransitionInPlace(transitionLength float32, state *State) string {
    feedrate := state.Palette.SideTransitionPurgeSpeed
    transitionSoFar := float32(0)
    sequence := moveToSideTransition(state)

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

    sequence += leaveSideTransition()
    return sequence
}

func sideTransitionOnEdge(transitionLength float32, state *State) string {
    eFeedrate := state.Palette.SideTransitionPurgeSpeed
    xyFeedrate := state.Palette.SideTransitionFeedrate
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
    sequence := moveToSideTransition(state)

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

    sequence += leaveSideTransition()
    return sequence
}

func sideTransition(transitionLength float32, state *State) string {
    sequence := ";TYPE:Side transition" + EOL
    if state.Palette.SideTransitionJog {
        sequence += sideTransitionOnEdge(transitionLength, state)
    } else {
        sequence += sideTransitionInPlace(transitionLength, state)
    }
    return sequence
}
