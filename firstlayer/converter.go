package firstlayer

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"strings"

	"mosaicmfg.com/ps-postprocess/gcode"
)

func DetermineToolsUsedInTheFirstLayer(inPath string, firstLayerStyleSettingsPath string) (error, FirstLayer) {
	// load style settings affected by first tool from JSON
	firstLayerStyleSettings, err := LoadFirstLayerStylesFromFile(firstLayerStyleSettingsPath)
	if err != nil {
		log.Fatalln(err)
	}

	// all the tools used in first layer
	toolUsedInFirstLayer := make(map[int]bool)
	const layerChangeComment = ";LAYER_CHANGE"
	// layerChangeComment appears once before initial layer change
	layer := -1
	err = gcode.ReadByLine(inPath, func(line gcode.Command, _ int) error {
		if layer > 1 {
			return gcode.ErrEarlyExit
		} else if isToolChange, tool := line.IsToolChange(); isToolChange {
			toolUsedInFirstLayer[tool] = true
		} else if strings.HasPrefix(line.Raw, layerChangeComment) {
			layer += 1
		}
		return nil
	})

	if err != nil {
		return err, FirstLayer{}
	}

	// compute the style settings values to be used in first layer
	var negInf = float32(math.Inf(-1))
	usedFirstLayerValues := FirstLayer{
		BedTemperature: 0,
		ZOffset:        negInf,
	}
	for key := range toolUsedInFirstLayer {
		if usedFirstLayerValues.BedTemperature < firstLayerStyleSettings.BedTemperature[key] {
			usedFirstLayerValues.BedTemperature = firstLayerStyleSettings.BedTemperature[key]
		}
		if usedFirstLayerValues.ZOffset < firstLayerStyleSettings.ZOffsetPerExt[key] {
			usedFirstLayerValues.ZOffset = firstLayerStyleSettings.ZOffsetPerExt[key]
		}
	}

	// Check if the final values are valid
	if usedFirstLayerValues.ZOffset == negInf {
		return fmt.Errorf("invalid ZOffset: %v", usedFirstLayerValues.ZOffset), FirstLayer{}
	}

	return nil, usedFirstLayerValues

}

func UseFirstLayerSettings(argv []string) error {
	argc := len(argv)

	if argc < 3 {
		log.Fatalln("expected 3 command-line arguments")
	}

	const EOL = "\r\n"

	inPath := argv[0]                      // unmodified G-code file
	outPath := argv[1]                     // modified G-code file
	firstLayerStyleSettingsPath := argv[2] // style settings that are affected by the first tool

	// determine the tools used in the first layer
	err, usedFirstLayerValues := DetermineToolsUsedInTheFirstLayer(inPath, firstLayerStyleSettingsPath)
	if err != nil {
		return err
	}

	// create out file
	outfile, createErr := os.Create(outPath)
	if createErr != nil {
		return createErr
	}
	writer := bufio.NewWriter(outfile)
	writeGCodeError := gcode.ReadByLine(inPath, func(line gcode.Command, linenNum int) error {
		if line.Command == "M140" {
			if _, ok := line.Params["s"]; ok {
				line.Params["s"] = usedFirstLayerValues.BedTemperature
				// set line.Raw to an empty string so that later line.String()
				// will generate the correct line rather then using line.raw
				line.Raw = ""
			}
		} else if line.Command == "M190" {
			if _, ok := line.Params["s"]; ok {
				line.Params["s"] = usedFirstLayerValues.BedTemperature
				line.Raw = ""
			} else if _, ok = line.Params["r"]; ok {
				line.Params["r"] = usedFirstLayerValues.BedTemperature
				line.Raw = ""
			}
		} else if value, ok := line.Params["z"]; ok {
			// Check if the command is one of the specified commands
			switch line.Command {
			case "G0", "G1", "G2", "G3", "G92":
				// z-offset
				line.Params["z"] = value + usedFirstLayerValues.ZOffset
				line.Raw = ""
			}
		}
		// write g-code to outPath file
		if _, err := writer.WriteString(line.String() + EOL); err != nil {
			return err
		}
		return nil
	})

	if writeGCodeError != nil {
		return err
	}
	if err := writer.Flush(); err != nil {
		return err
	}
	if err := outfile.Close(); err != nil {
		return err
	}

	usedFirstLayerValues.Save(outPath + ".firstLayerResults")
	return nil
}
