package sequences

import (
	"mosaicmfg.com/ps-postprocess/gcode"
	"strings"
)

func filterToolchangeCommands(evaluatedScript string) string {
	lines := gcode.ParseLines(evaluatedScript)
	filteredLines := make([]string, 0)

	for _, line := range lines {
		if isToolChange, _ := line.IsToolChange(); !isToolChange {
			filteredLines = append(filteredLines, line.Raw)
		}
	}

	return strings.Join(filteredLines, "\n")
}
