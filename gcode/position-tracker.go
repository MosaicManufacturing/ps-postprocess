package gcode

type PositionTracker struct {
	CurrentX        float32
	CurrentY        float32
	CurrentZ        float32
	CurrentFeedrate float32
}

func (pt *PositionTracker) TrackInstruction(instruction Command) {
	if len(instruction.Command) == 0 {
		return
	}
	if instruction.IsLinearMove() || instruction.IsArcMove() {
		if x, ok := instruction.Params["x"]; ok {
			pt.CurrentX = x
		}
		if y, ok := instruction.Params["y"]; ok {
			pt.CurrentY = y
		}
		if z, ok := instruction.Params["z"]; ok {
			pt.CurrentZ = z
		}
		if f, ok := instruction.Params["f"]; ok {
			pt.CurrentFeedrate = f
		}
	} else if instruction.IsHome() {
		if instruction.Flags["x"] || instruction.Flags["y"] || instruction.Flags["z"] {
			// flags present == only home some axes
			if instruction.Flags["x"] {
				pt.CurrentX = 0
			}
			if instruction.Flags["y"] {
				pt.CurrentY = 0
			}
			if instruction.Flags["z"] {
				pt.CurrentZ = 0
			}
		} else {
			// no flags == home all axes
			pt.CurrentX = 0
			pt.CurrentY = 0
			pt.CurrentZ = 0
		}
	}
}
