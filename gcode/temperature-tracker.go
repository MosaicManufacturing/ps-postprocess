package gcode

type TemperatureTracker struct {
	Extruder float32
	Bed      float32
	Chamber  float32
}

func (tt *TemperatureTracker) TrackInstruction(instruction Command) {
	if len(instruction.Command) == 0 {
		return
	}
	if instruction.Command == "M104" || instruction.Command == "M109" {
		if s, ok := instruction.Params["s"]; ok {
			tt.Extruder = s
		}
	} else if instruction.Command == "M140" || instruction.Command == "M190" {
		if s, ok := instruction.Params["s"]; ok {
			tt.Bed = s
		}
	} else if instruction.Command == "M141" || instruction.Command == "M191" {
		if s, ok := instruction.Params["s"]; ok {
			tt.Chamber = s
		}
	}
}
