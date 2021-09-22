package gcode

import (
	"fmt"
	"strings"
)

type Params map[string]float32
type Flags map[string]bool

type Command struct {
	Raw string
	Command string
	Comment string
	Params Params
	Flags Flags
}

func NewCommand(raw, command, comment string, params Params, flags Flags) Command {
	return Command{
		Raw:     raw,
		Command: command,
		Comment: comment,
		Params:  params,
		Flags:   flags,
	}
}

func (gcc Command) IsLinearMove() bool {
	// slight optimization: G1 is much more common, so check for that first
	return gcc.Command == "G1" || gcc.Command == "G0"
}

func (gcc Command) IsArcMove() bool {
	return gcc.Command == "G2" || gcc.Command == "G3"
}

func (gcc Command) IsHome() bool {
	return gcc.Command == "G28"
}

func (gcc Command) IsSetExtrusionMode() (bool, bool) {
	isSetExtrusionMode := gcc.Command == "M82" || gcc.Command == "M83"
	isModeRelative := gcc.Command == "M83"
	return isSetExtrusionMode, isModeRelative
}

func (gcc Command) IsSetPosition() bool {
	return gcc.Command == "G92"
}

func (gcc Command) String() string {
	line := ""
	if gcc.Command != "" {
		line += gcc.Command
		for param, value := range gcc.Params {
			line += fmt.Sprintf(" %s%f", strings.ToUpper(param), value)
		}
		for flag := range gcc.Flags {
			line += fmt.Sprintf(" %s", strings.ToUpper(flag))
		}
	}
	if len(gcc.Comment) > 0 {
		line += fmt.Sprintf("; %s", gcc.Comment)
	}
	return line
}
