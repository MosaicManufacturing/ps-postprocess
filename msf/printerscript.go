package msf

import (
    "mosaicmfg.com/ps-postprocess/gcode"
	"mosaicmfg.com/ps-postprocess/printerscript"
)

func getTimeEstimate(command gcode.Command, state *State) float32 {
    if command.IsLinearMove() {
        feedrate := state.XYZF.CurrentFeedrate
        if f, ok := command.Params["f"]; ok {
            feedrate = f
        }
        nextX := state.XYZF.CurrentX
        if x, ok := command.Params["x"]; ok {
            nextX = x
        }
        nextY := state.XYZF.CurrentY
        if y, ok := command.Params["y"]; ok {
            nextY = y
        }
        return estimateMoveTime(state.XYZF.CurrentX, state.XYZF.CurrentY, nextX, nextY, feedrate)
    }
    if command.Command == "G4" {
        if ms, ok := command.Params["p"]; ok {
            // e.g. G4 P5000
            return ms / 1000
        }
        if s, ok := command.Params["s"]; ok {
            // e.g. G4 S5
            return s
        }
        return 0
    }
    return 0
}

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
        // update time estimate first
        timeEstimate := getTimeEstimate(command, state)
        if timeEstimate > 0 {
            state.TimeEstimate += timeEstimate
        }
        // update other state trackers
        state.E.TrackInstruction(command)
        state.XYZF.TrackInstruction(command)
        state.Temperature.TrackInstruction(command)
    }
    return sequence, nil
}
