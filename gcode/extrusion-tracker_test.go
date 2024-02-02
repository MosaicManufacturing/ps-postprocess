package gcode

import "testing"

func parseLines(lines []string) []Command {
	commands := make([]Command, 0, len(lines))
	for _, line := range lines {
		command := ParseLine(line)
		commands = append(commands, command)
	}
	return commands
}

func trackCommands(lines []string) ExtrusionTracker {
	var et ExtrusionTracker
	commands := parseLines(lines)
	for _, command := range commands {
		et.TrackInstruction(command)
	}
	return et
}

func expectTotalExtrusion(t *testing.T, et *ExtrusionTracker, totalExtrusion float32) {
	if et.TotalExtrusion != totalExtrusion {
		t.Logf("expected TotalExtrusion = %f, got %f", totalExtrusion, et.TotalExtrusion)
		t.Fail()
	}
}

func Test_AbsoluteENoRetract(t *testing.T) {
	et := trackCommands([]string{
		"M82",
		"G1 E1",
		"G1 E2",
		"G1 E3",
	})
	expectTotalExtrusion(t, &et, 3)
}

func Test_RelativeENoRetract(t *testing.T) {
	et := trackCommands([]string{
		"M83",
		"G1 E1",
		"G1 E2",
		"G1 E3",
	})
	expectTotalExtrusion(t, &et, 6)
}

func Test_AbsoluteEWithRetract(t *testing.T) {
	// retract length == restart length
	et := trackCommands([]string{
		"M82",
		"G1 E1",
		"G1 E2",
		"G1 E3",
		"G1 E2",
		"G1 E3",
		"G1 E4",
	})
	expectTotalExtrusion(t, &et, 4)

	// retract length > restart length
	et = trackCommands([]string{
		"M82",
		"G1 E1",
		"G1 E2",
		"G1 E3",
		"G1 E2",
		"G1 E2.5",
		"G1 E3",
		"G1 E4",
	})
	expectTotalExtrusion(t, &et, 4)

	// retract length < restart length
	et = trackCommands([]string{
		"M82",
		"G1 E1",
		"G1 E2",
		"G1 E3",
		"G1 E2",
		"G1 E3.5",
		"G1 E4",
	})
	expectTotalExtrusion(t, &et, 4)
}

func Test_RelativeEWithRetract(t *testing.T) {
	// retract length == restart length
	et := trackCommands([]string{
		"M83",
		"G1 E1",
		"G1 E2",
		"G1 E3",
		"G1 E-1",
		"G1 E1",
		"G1 E4",
	})
	expectTotalExtrusion(t, &et, 10)

	// retract length > restart length
	et = trackCommands([]string{
		"M83",
		"G1 E1",
		"G1 E2",
		"G1 E3",
		"G1 E-1",
		"G1 E0.5",
		"G1 E1",
		"G1 E2",
	})
	expectTotalExtrusion(t, &et, 8.5)

	// retract length < restart length
	et = trackCommands([]string{
		"M83",
		"G1 E1",
		"G1 E2",
		"G1 E3",
		"G1 E-1",
		"G1 E1.5",
		"G1 E2",
	})
	expectTotalExtrusion(t, &et, 8.5)
}

func Test_AbsoluteEWithG92(t *testing.T) {
	// E == 0
	et := trackCommands([]string{
		"M82",
		"G1 E1",
		"G1 E2",
		"G92 E0",
		"G1 E3",
	})
	expectTotalExtrusion(t, &et, 5)

	// E != 0
	et = trackCommands([]string{
		"M82",
		"G1 E1",
		"G1 E2",
		"G92 E1",
		"G1 E3",
	})
	expectTotalExtrusion(t, &et, 4)
}

func Test_RelativeEWithG92(t *testing.T) {
	// E == 0
	et := trackCommands([]string{
		"M83",
		"G1 E1",
		"G1 E2",
		"G92 E0",
		"G1 E3",
	})
	expectTotalExtrusion(t, &et, 6)

	// E != 0
	et = trackCommands([]string{
		"M83",
		"G1 E1",
		"G1 E2",
		"G92 E1",
		"G1 E3",
	})
	expectTotalExtrusion(t, &et, 6)
}

func Test_ModeChanges(t *testing.T) {
	et := trackCommands([]string{
		"M82",    // 0
		"G1 E1",  // 1
		"G1 E2",  // 2
		"G1 E1",  // 2
		"G92 E1", // 2 (with 1 mm retraction)
		"G1 E3",  // 3 (with 0 retraction)
		"M83",    // 3
		"G1 E1",  // 4
		"G1 E1",  // 5
		"G1 E2",  // 7
		"G1 E-1", // 7
		"G1 E1",  // 7
		"G92 E0", // 7
		"G1 E1",  // 8
	})
	expectTotalExtrusion(t, &et, 8)
}
