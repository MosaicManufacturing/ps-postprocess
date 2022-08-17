package msf

import "fmt"

type Transition struct {
	Layer            int     // 0-indexed
	From             int     // tool, 0-indexed
	To               int     // tool, 0-indexed
	TotalExtrusion   float32 // total print extrusion at start of transition (excluding earlier transitions)
	TransitionLength float32 // actual transition length as specified by user
	PurgeLength      float32 // amount of filament to extrude during this transition
	UsableInfill     float32 // subtract this amount from the splice length to transition in infill
}

func (t Transition) String() string {
	return fmt.Sprintf(
		"Layer = %d, From = %d, To = %d, TotalExtrusion = %f, TransitionLength = %f, PurgeLength = %f, UsableInfill = %f",
		t.Layer,
		t.From,
		t.To,
		t.TotalExtrusion,
		t.TransitionLength,
		t.PurgeLength,
		t.UsableInfill,
	)
}
