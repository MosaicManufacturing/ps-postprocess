package msf

import "mosaicmfg.com/ps-postprocess/gcode"

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
		Params: map[string]float32{
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
		Params: map[string]float32{
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

func getFirmwareRetract() string {
	return "G10" + EOL
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
		Params: map[string]float32{
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

func getFirmwareRestart() string {
	return "G11" + EOL
}

func getZTravel(state *State, toZ float32, comment string) string {
	if state.XYZF.CurrentZ == toZ {
		return ""
	}
	zTravel := gcode.Command{
		Command: "G1",
		Comment: comment,
		Params: map[string]float32{
			"z": toZ,
			"f": state.Palette.TravelSpeedZ,
		},
	}
	state.TimeEstimate += estimateZMoveTime(state.XYZF.CurrentZ, toZ, state.Palette.TravelSpeedZ)
	state.XYZF.TrackInstruction(zTravel)
	return zTravel.String() + EOL
}

func getXYTravel(state *State, toX, toY, feedrate float32, comment string) string {
	if state.XYZF.CurrentX == toX && state.XYZF.CurrentY == toY {
		return ""
	}
	xyTravel := gcode.Command{
		Command: "G1",
		Comment: comment,
		Params: map[string]float32{
			"x": toX,
			"y": toY,
			"f": feedrate,
		},
	}
	state.TimeEstimate += estimateMoveTime(state.XYZF.CurrentX, state.XYZF.CurrentY, toX, toY, feedrate)
	state.XYZF.TrackInstruction(xyTravel)
	return xyTravel.String() + EOL
}

func getPurge(state *State, distance, feedrate float32) string {
	if distance == 0 {
		return ""
	}
	eValue := distance
	if !state.E.RelativeExtrusion {
		eValue += state.E.CurrentExtrusionValue
	}
	purge := gcode.Command{
		Command: "G1",
		Params: map[string]float32{
			"e": eValue,
			"f": feedrate,
		},
	}
	state.TimeEstimate += estimatePurgeTime(distance, feedrate)
	state.E.TrackInstruction(purge)
	return purge.String() + EOL
}

func getXYExtrusion(state *State, toX, toY, distance, feedrate float32) string {
	eValue := distance
	if !state.E.RelativeExtrusion {
		eValue += state.E.CurrentExtrusionValue
	}
	purge := gcode.Command{
		Command: "G1",
		Params: map[string]float32{
			"x": toX,
			"y": toY,
			"e": eValue,
			"f": feedrate,
		},
	}
	state.TimeEstimate += estimateMoveTime(state.XYZF.CurrentX, state.XYZF.CurrentY, toX, toY, feedrate)
	state.XYZF.TrackInstruction(purge)
	state.E.TrackInstruction(purge)
	return purge.String() + EOL
}
