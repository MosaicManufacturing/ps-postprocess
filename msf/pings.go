package msf

import (
    "fmt"
    "math"
)

func getPingRetract(palette *Palette) (bool, string) {
    if palette.PingRetractDistance == 0 {
        return false, ""
    }
    return true, fmt.Sprintf("G1 E%.5f F%.1f", -palette.PingRetractDistance, palette.PingRetractFeedrate)
}

func getPingRestart(palette *Palette) (bool, string) {
    if palette.PingRestartDistance == 0 {
        return false, ""
    }
    return true, fmt.Sprintf("G1 E%.5f F%.1f", palette.PingRestartDistance, palette.PingRestartFeedrate)
}

func getDwellPause(durationMS int) string {
    sequence := ""
    for durationMS > 0 {
        if durationMS > 4000 {
            sequence += "G4 P4000" + EOL
            sequence += "G1" + EOL
            durationMS -= 4000
        } else {
            sequence += fmt.Sprintf("G4 P%d%s", durationMS, EOL)
            durationMS = 0
        }
    }
    return sequence
}

func getTowerJogPause(state *State, durationMS int) string {
    durationMM := float32(durationMS) / 1000

    // PrusaSlicer tower extrusions always run east-west, so either jog eastward or westward,
    // depending on which edge we're currently closer to

    x1 := state.XYZF.CurrentX
    y1 := state.XYZF.CurrentY
    x2 := x1
    y2 := y1
    if math.Abs(float64(x1 - state.TowerBoundingBox.min[0])) > math.Abs(float64(state.TowerBoundingBox.max[0] - x1)) {
        // closer to west edge -- jog westward
        x2 -= durationMM
    } else {
        // closer to east edge -- jog eastward
        x2 += durationMM
    }

    const totalJogs = 5
    const feedrate = 10

    sequence := ""
    for i := 0; i < totalJogs; i++ {
        sequence += fmt.Sprintf("G1 X%.3f Y%.3f F%d%s", x2, y2, feedrate, EOL)
        sequence += fmt.Sprintf("G1 X%.3f Y%.3f F%d%s", x1, y1, feedrate, EOL)
    }

    return sequence
}

func doSideTransitionInPlaceAccessoryPing(transitionSoFar, transitionLength float32, state *State) (string, float32) {
   // todo: implement
    return "", 0
}

func doSideTransitionOnEdgeAccessoryPing(transitionSoFar, transitionLength float32, state *State) (string, float32) {
   // todo: implement
    return "", 0
}
