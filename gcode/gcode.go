package gcode

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
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
	return gcc.Command == "G0" || gcc.Command == "G1"
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

func ParseLine(raw string) Command {
	params := Params{}
	flags := Flags{}
	if len(strings.Trim(raw, " ")) == 0 {
		return NewCommand(raw, "", "", params, flags)
	}

	// collect any comments
	commentSplit := strings.Split(raw, ";")
	line := strings.Trim(commentSplit[0], " ")
	comment := strings.Trim(strings.Join(commentSplit[1:], ";"), " ")

	// split the line at spaces
	argv := strings.Fields(line)
	argc := len(argv)

	// get the command
	command := ""
	if len(argv) > 0 {
		command = strings.ToUpper(argv[0])
	}

	// collect arguments
	for i := 1; i < argc; i++ {
		key := strings.ToLower(argv[i][0:1])
		value := argv[i][1:]
		if len(value) == 0 {
			flags[key] = true
		} else {
			floatValue, err := strconv.ParseFloat(value, 32)
			if err == nil {
				params[key] = float32(floatValue)
			}
		}
	}

	return NewCommand(raw, command, comment, params, flags)
}

func (gcc Command) ToString() string {
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

type LineCallback func(Command)

func ReadByLine (path string, callback LineCallback) (err error) {
	infile, openErr := os.Open(path)
	if openErr != nil {
		err = openErr
		return
	}
	defer func() {
		if closeErr := infile.Close(); closeErr != nil {
			err = closeErr
		}
	}()
	reader := bufio.NewReader(infile)
	for {
		line, isPrefix, readErr := reader.ReadLine()
		if readErr == io.EOF {
			break
		}
		if readErr != nil {
			err = readErr
			return
		}
		if isPrefix {
			var fragment []byte
			for isPrefix {
				fragment, isPrefix, readErr = reader.ReadLine()
				if readErr == io.EOF {
					break
				}
				if readErr != nil {
					err = readErr
					return
				}
				line = append(line, fragment...)
			}
		}
		gcode := ParseLine(string(line))
		callback(gcode)
	}
	return nil
}

