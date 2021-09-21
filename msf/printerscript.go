package msf

import (
    "../gcode"
	"../printerscript"
)

func evaluateScript(script printerscript.Tree, locals map[string]float64, state *State) (string, error) {
    opts := printerscript.InterpreterOptions{
        EOL:             EOL,
        TrailingNewline: true,
        Locals:          locals,
    }
    result, err := printerscript.EvaluateTree(script, opts)
    if err != nil {
        return "", err
    }
    sequence := result.Output
    commands := gcode.ParseLines(sequence)
    for _, command := range commands {
        state.E.TrackInstruction(command)
        state.XYZF.TrackInstruction(command)
        state.Temperature.TrackInstruction(command)
       // TODO: add to time estimate
    }
    return sequence, nil
}
