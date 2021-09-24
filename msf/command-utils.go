package msf

import "../gcode"

func resetEAxis(state *State) string {
    if state.E.RelativeExtrusion {
        // G92 not needed in relative mode
        return ""
    }
    if state.E.CurrentExtrusionValue == 0 {
        // G92 is redundant
        return ""
    }

    // reset extrusion distance
    reset := gcode.Command{
        Command: "G92",
        Params:  map[string]float32{
            "e": 0,
        },
        Comment: "reset extrusion distance",
    }
    state.E.TrackInstruction(reset)
    return reset.String() + EOL
}

func getRetract(state *State, distance, feedrate float32) string {
    if distance == 0 {
        return ""
    }
    if distance < 0 {
        distance = -distance
    }

    eParam := -distance
    if !state.E.RelativeExtrusion {
        eParam = state.E.CurrentExtrusionValue - distance
    }
    retract := gcode.Command{
        Command: "G1",
        Params:  map[string]float32{
            "e": eParam,
            "f": feedrate,
        },
        Comment: "retract",
    }
    state.TimeEstimate += estimatePurgeTime(distance, feedrate)
    state.XYZF.TrackInstruction(retract)
    state.E.TrackInstruction(retract)
    return retract.String() + EOL
}

func getRestart(state *State, distance, feedrate float32) string {
    if distance == 0 {
        return ""
    }
    if distance < 0 {
        distance = -distance
    }

    eParam := distance
    if !state.E.RelativeExtrusion {
        eParam = state.E.CurrentExtrusionValue + distance
    }
    restart := gcode.Command{
        Command: "G1",
        Params:  map[string]float32{
            "e": eParam,
            "f": feedrate,
        },
        Comment: "unretract",
    }
    state.TimeEstimate += estimatePurgeTime(distance, feedrate)
    state.XYZF.TrackInstruction(restart)
    state.E.TrackInstruction(restart)
    return restart.String() + EOL
}
