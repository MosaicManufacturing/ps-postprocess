package msf

// TransitionLookahead contains lookahead data for the XYZ coordinates of the start
// point of the first print line after a transition. A print line is defined by
// a linear move command with an X and/or Y parameter and positive extrusion.
// MovedXY and MovedZ are not affected by the print line that ends the lookahead.
type TransitionLookahead struct {
	X       float32 // X position of the next print line after the transition
	Y       float32 // Y position of the next print line after the transition
	Z       float32 // Z position of the next print line after the transition
	MovedXY bool    // true iff X or Y movement was seen during lookahead process
	MovedZ  bool    // true iff Z movement was seen during lookahead process
}

// ScriptLookahead contains XYZ coordinates representing the nextX/nextY/nextZ
// PrinterScript locals for "after side transition" scripts. TODO: use me
type ScriptLookahead struct {
	X            float32
	Y            float32
	Z            float32
	InitializedX bool
	InitializedY bool
	InitializedZ bool
}
