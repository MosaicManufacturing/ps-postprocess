package gcode

type ExtrusionTracker struct {
    RelativeExtrusion bool // true == relative, false == absolute
    TotalExtrusion float32 // total filament consumption -- never decreases
    CurrentExtrusionValue float32 // current position of the E axis
    PreviousExtrusionValue float32 // last position of the E axis
    LastExtrudeWasRetract bool // true == last E movement was negative, false == positive E
    LastRetractDistance float32 // most recent E axis value of negative E
    CurrentRetraction float32 // current total retraction (negative, need this much positive E to be primed)
    LastCommandWasG92 bool // true if the last E modification was a manual position being set
}

func (et *ExtrusionTracker) TrackInstruction(instruction Command) {
    if len(instruction.Command) == 0 {
        return
    }
    if instruction.IsLinearMove() || instruction.IsArcMove() {
        if eValue, ok := instruction.Params["e"]; ok {
            et.PreviousExtrusionValue = et.CurrentExtrusionValue
            et.CurrentExtrusionValue = eValue
            if et.RelativeExtrusion {
                // relative extrusion
                et.TotalExtrusion += eValue
                if eValue < 0 {
                    // retraction
                    et.LastExtrudeWasRetract = true
                    et.LastRetractDistance = eValue
                    et.CurrentRetraction += eValue
                } else if eValue > 0 {
                    et.LastExtrudeWasRetract = false
                    if et.CurrentRetraction + eValue >= 0 {
                        // normal extrusion
                        et.CurrentRetraction = 0
                    } else {
                        // restart
                        et.CurrentRetraction += eValue
                    }
                }
            } else {
                // absolute extrusion
                et.TotalExtrusion += eValue - et.PreviousExtrusionValue
                if et.CurrentExtrusionValue < et.PreviousExtrusionValue {
                    // retraction
                    et.LastExtrudeWasRetract = true
                    et.LastRetractDistance = eValue - et.PreviousExtrusionValue
                    et.CurrentRetraction += et.LastRetractDistance
                } else if et.CurrentExtrusionValue > et.PreviousExtrusionValue {
                    et.LastExtrudeWasRetract = false
                    if et.CurrentRetraction + (eValue - et.PreviousExtrusionValue) >= 0 {
                        // normal extrusion
                        et.CurrentRetraction = 0
                    } else {
                        // restart
                        et.CurrentRetraction += eValue - et.PreviousExtrusionValue
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
                et.CurrentExtrusionValue = eValue
            } else if aValue, ok := instruction.Params["a"]; ok {
                et.LastCommandWasG92 = true
                et.CurrentExtrusionValue = aValue
            } else if bValue, ok := instruction.Params["b"]; ok {
                et.LastCommandWasG92 = true
                et.CurrentExtrusionValue = bValue
            }
        } else {
            et.LastCommandWasG92 = true
            et.CurrentExtrusionValue = 0
        }
    }
}
