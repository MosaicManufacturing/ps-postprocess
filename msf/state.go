package msf

import (
	"mosaicmfg.com/ps-postprocess/gcode"
	"mosaicmfg.com/ps-postprocess/sequences"
)

type State struct {
	Palette       *Palette // reference stored here to reduce arguments passed to routines
	MSF           *MSF     // reference stored here to reduce arguments passed to routines
	Tower         *Tower
	PingExtrusion float32 // stored to avoid re-calculating every time

	CurrentLayer int
	E            gcode.ExtrusionTracker
	XYZF         gcode.PositionTracker
	Temperature  gcode.TemperatureTracker
	TimeEstimate float32

	PastStartSequence          bool
	FirstToolChange            bool // don't treat the first T command as a toolchange
	CurrentTool                int
	CurrentlyTransitioning     bool
	NeedsPostTransitionZAdjust bool
	PostTransitionZ            float32
	OnWipeTower                bool
	TowerBoundingBox           gcode.BoundingBox

	// for maintaining print preview comment state after transitions
	CurrentPathTypeLine string
	CurrentWidthLine    string

	LastPingStart    float32
	CurrentlyPinging bool
	CurrentPingStart float32
	NextPingStart    float32

	TransitionNextPositions []SideTransitionLookahead
	Locals                  sequences.Locals // for PrinterScript side transition sequences
}

func NewState(palette *Palette) State {
	return State{
		Palette:         palette,
		FirstToolChange: true,
		CurrentLayer:    -1,
		PingExtrusion:   palette.GetPingExtrusion(),
	}
}
