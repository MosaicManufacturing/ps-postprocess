package gcode

import (
    "strconv"
    "strings"
)

func normalizeNewlines(input string) string {
    // replace CR LF \r\n (Windows) with LF \n (Unix)
    input = strings.ReplaceAll(input, "\r\n", "\n")
    // replace CF \r (Mac Classic) with LF \n (Unix)
    return strings.ReplaceAll(input, "\r", "\n")
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

func ParseLines(raw string) []Command {
    raw = normalizeNewlines(raw)
    lines := strings.Split(raw, "\n")
    commands := make([]Command, 0, len(lines))
    for _, line := range lines {
        commands = append(commands, ParseLine(line))
    }
    return commands
}
