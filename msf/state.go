package msf

import "../gcode"

type State struct {
    Palette *Palette // reference stored here to reduce arguments passed to routines
    MSF *MSF // reference stored here to reduce arguments passed to routines

    E gcode.ExtrusionTracker
    XYZF gcode.PositionTracker
    CurrentX float32
    CurrentY float32
    CurrentZ float32

    FirstToolChange bool // don't treat the first T command as a toolchange
    CurrentTool int
    CurrentlyTransitioning bool
    OnWipeTower bool

    LastPingStart float32
    CurrentlyPinging bool
    CurrentPingStart float32
    NextPingStart float32
}

func NewState(palette *Palette) State {
    return State{
        Palette: palette,
        FirstToolChange: true,
    }
}
