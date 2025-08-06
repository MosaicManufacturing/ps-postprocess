package gcode

type ExtrusionTracker struct {
	RelativeExtrusion      bool    // true == relative, false == absolute
	TotalExtrusion         float64 // total filament consumption -- never decreases
	CurrentExtrusionValue  float64 // current position of the E axis
	PreviousExtrusionValue float64 // last position of the E axis
	LastExtrudeWasRetract  bool    // true == last E movement was negative, false == positive E
	LastRetractDistance    float64 // most recent E axis value of negative E
	CurrentRetraction      float64 // current total retraction (negative, need this much positive E to be primed)
	LastCommandWasG92      bool    // true if the last E modification was a manual position being set
}

func (et *ExtrusionTracker) TrackInstruction(instruction Command) {
	if len(instruction.Command) == 0 {
		return
	}
	if instruction.IsLinearOrArcMove() {
		if eValue, ok := instruction.Params["e"]; ok {
			eValue64 := float64(eValue)
			et.PreviousExtrusionValue = et.CurrentExtrusionValue
			et.CurrentExtrusionValue = eValue64
			if et.RelativeExtrusion {
				// relative extrusion
				et.TotalExtrusion += eValue64
				if eValue < 0 {
					// retraction
					et.LastExtrudeWasRetract = true
					et.LastRetractDistance = eValue64
					et.CurrentRetraction += eValue64
				} else if eValue > 0 {
					et.LastExtrudeWasRetract = false
					if et.CurrentRetraction+eValue64 >= 0 {
						// normal extrusion
						et.CurrentRetraction = 0
					} else {
						// restart
						et.CurrentRetraction += eValue64
					}
				}
			} else {
				// absolute extrusion
				et.TotalExtrusion += eValue64 - et.PreviousExtrusionValue
				if et.CurrentExtrusionValue < et.PreviousExtrusionValue {
					// retraction
					et.LastExtrudeWasRetract = true
					et.LastRetractDistance = eValue64 - et.PreviousExtrusionValue
					et.CurrentRetraction += et.LastRetractDistance
				} else if et.CurrentExtrusionValue > et.PreviousExtrusionValue {
					et.LastExtrudeWasRetract = false
					if et.CurrentRetraction+(eValue64-et.PreviousExtrusionValue) >= 0 {
						// normal extrusion
						et.CurrentRetraction = 0
					} else {
						// restart
						et.CurrentRetraction += eValue64 - et.PreviousExtrusionValue
					}
				}
			}
		}
	} else if setExtrusionMode, relative := instruction.IsSetExtrusionMode(); setExtrusionMode {
		et.RelativeExtrusion = relative
	} else if instruction.IsSetPosition() {
		hasParamsOrFlags := len(instruction.Params) > 0 || len(instruction.Flags) > 0
		if hasParamsOrFlags {
			if eValue, ok := instruction.Params["e"]; ok {
				et.LastCommandWasG92 = true
				et.CurrentExtrusionValue = float64(eValue)
			} else if aValue, ok := instruction.Params["a"]; ok {
				et.LastCommandWasG92 = true
				et.CurrentExtrusionValue = float64(aValue)
			} else if bValue, ok := instruction.Params["b"]; ok {
				et.LastCommandWasG92 = true
				et.CurrentExtrusionValue = float64(bValue)
			}
		} else {
			et.LastCommandWasG92 = true
			et.CurrentExtrusionValue = 0
		}
	}
}
