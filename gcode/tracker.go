package gcode

type Tracker interface {
    TrackInstruction(instruction Command)
}
