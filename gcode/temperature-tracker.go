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
	if instruction.Command == "M104" {
		if temp, ok := instruction.Params["s"]; ok {
			tt.Extruder = temp
		}
	} else if instruction.Command == "M109" {
		if temp, ok := instruction.Params["s"]; ok {
			tt.Extruder = temp
		} else if temp, ok = instruction.Params["r"]; ok {
			tt.Extruder = temp
		}
	} else if instruction.Command == "M140" {
		if temp, ok := instruction.Params["s"]; ok {
			tt.Bed = temp
		}
	} else if instruction.Command == "M190" {
		if temp, ok := instruction.Params["s"]; ok {
			tt.Bed = temp
		} else if temp, ok = instruction.Params["r"]; ok {
			tt.Bed = temp
		}
	} else if instruction.Command == "M141" {
		if temp, ok := instruction.Params["s"]; ok {
			tt.Chamber = temp
		}
	} else if instruction.Command == "M191" {
		if temp, ok := instruction.Params["s"]; ok {
			tt.Chamber = temp
		} else if temp, ok = instruction.Params["r"]; ok {
			tt.Chamber = temp
		}
	}
}
