package msf

import (
    "../gcode"
	"fmt"
    "math"
)

func getSideTransitionStartPosition(currentX, currentY float32) (float32, float32) {
    // todo
    return 0, 0
}

func moveToSideTransition(transitionLength float32, palette *Palette) string {
    // todo
    return ""
}

func checkSideTransitionPings(transitionSoFar, transitionLength float32, palette *Palette) string {
    // todo
    return ""
}

func sideTransitionInPlace(transitionLength float32, palette *Palette, eTracker *gcode.ExtrusionTracker) string {
    feedrate := palette.SideTransitionPurgeSpeed
    transitionSoFar := float32(0)
    sequence := ""

    for transitionSoFar < transitionLength {
        sequence += checkSideTransitionPings(transitionSoFar, transitionLength, palette)
        nextPurgeExtrusion := float32(math.Min(10, float64(transitionLength - transitionSoFar)))
        nextE := nextPurgeExtrusion
        if !eTracker.RelativeExtrusion {
            nextE = eTracker.CurrentExtrusionValue + nextPurgeExtrusion
        }
        purge := gcode.Command{
            Raw:     fmt.Sprintf("G1 E%.5f F%.1f", nextE, feedrate),
            Command: "G1",
            Params:  map[string]float32{
                "e": nextE,
                "f": feedrate,
            },
            Flags: map[string]bool{},
        }
        eTracker.TrackInstruction(purge)
        sequence += purge.Raw + EOL
    }

    // todo
    return sequence
}

func sideTransitionOnEdge(transitionLength float32, palette *Palette, eTracker *gcode.ExtrusionTracker) string {
    // todo
    return ""
}

func sideTransition(transitionLength float32, palette *Palette, eTracker *gcode.ExtrusionTracker) string {
    // todo: move to side and do the transition purge
    //   - make sure to track all of this with eTracker, or the next splice will be very short!

    sequence := moveToSideTransition(transitionLength, palette)
    if palette.SideTransitionJog {
        sequence += sideTransitionOnEdge(transitionLength, palette, eTracker)
    } else {
        sequence += sideTransitionInPlace(transitionLength, palette, eTracker)
    }
    return sequence
}
