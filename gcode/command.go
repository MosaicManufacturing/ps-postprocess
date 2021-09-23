package gcode

import (
	"fmt"
	"math"
	"sort"
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

func FormatFloat(value float64) string {
	// round to 5 decimal places first
	value = math.Round(value * 10e5) / 10e5
	// output with exactly 5 decimal places
	valStr := fmt.Sprintf("%.5f", value)
	// remove trailing zeros, and the decimal point if we reach it
	valStr = strings.TrimRight(strings.TrimRight(valStr, "0"), ".")
	// special-case for numbers that were printed as 0.00000
	if len(valStr) == 0 {
		valStr = "0"
	}
	return valStr
}

func scoreParamKey(letter string) int {
	switch letter {
	case "X":
		return 5
	case "Y":
		return 4
	case "Z":
		return 3
	case "E":
		return 2
	case "F":
		return 1
	default:
		return 0
	}
}

func (gcc Command) String() string {
	if len(gcc.Raw) > 0 {
		return gcc.Raw
	}

	line := ""
	if gcc.Command != "" {
		line += gcc.Command
		paramsAndFlags := make([]string, 0, len(gcc.Params) + len(gcc.Flags))

		for param, value := range gcc.Params {
			paramString := fmt.Sprintf("%s%s", strings.ToUpper(param), FormatFloat(float64(value)))
			paramsAndFlags = append(paramsAndFlags, paramString)
		}
		for flag := range gcc.Flags {
			flagString := strings.ToUpper(flag)
			paramsAndFlags = append(paramsAndFlags, flagString)
		}
		sort.Slice(paramsAndFlags, func(i, j int) bool {
			// sorting logic:
			// X, Y, Z, E, F, then alphabetical
			iKey := paramsAndFlags[i][0:1]; iScore := scoreParamKey(iKey)
			jKey := paramsAndFlags[j][0:1]; jScore := scoreParamKey(jKey)
			if iScore == 0 && jScore == 0 {
				// just alphabetical
				return iKey < jKey
			}
			// at least one element with a priority
			return iScore > jScore
		})
		if len(paramsAndFlags) > 0 {
			line += " " + strings.Join(paramsAndFlags, " ")
		}
	}
	if len(gcc.Comment) > 0 {
		if len(line) > 0 {
			line += " "
		}
		line += fmt.Sprintf("; %s", gcc.Comment)
	}
	return line
}
