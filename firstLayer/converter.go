package firstLayer

import (
	"bufio"
	"log"
	"os"
	"strings"

	"mosaicmfg.com/ps-postprocess/gcode"
)

func UseFirstLayerSettings(argv []string) error {
	argc := len(argv)

	if argc < 3 {
		log.Fatalln("expected 3 command-line arguments")
	}

	const EOL = "\r\n"

	inPath := argv[0]                      // unmodified G-code file
	outPath := argv[1]                     // modified G-code file
	firstLayerStyleSettingsPath := argv[2] // style settings that are effected by the first tool

	// load style settings effected by first tool from JSON
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
		return err
	}

	// computer the style settings values to be used in first layer
	usedFirstLayerValues := FirstLayer{
		BedTemperature: 0,
		ZOffset:        0,
	}
	for key := range toolUsedInFirstLayer {
		if usedFirstLayerValues.BedTemperature < firstLayerStyleSettings.BedTemperature[key] {
			usedFirstLayerValues.BedTemperature = firstLayerStyleSettings.BedTemperature[key]
		}
		if usedFirstLayerValues.ZOffset < firstLayerStyleSettings.ZOffsetPerExt[key] {
			usedFirstLayerValues.ZOffset = firstLayerStyleSettings.ZOffsetPerExt[key]
		}
	}

	// create out file
	outfile, createErr := os.Create(outPath)
	if createErr != nil {
		return createErr
	}
	defer outfile.Close()

	writer := bufio.NewWriter(outfile)
	defer writer.Flush()

	output := ""
	err = gcode.ReadByLine(inPath, func(line gcode.Command, linenNum int) error {
		if line.Command == "M140" {
			// bed temp
			if _, ok := line.Params["s"]; ok {
				line.Params["s"] = usedFirstLayerValues.BedTemperature
				line.Raw = ""
				output += line.String() + EOL
				return nil
			}
		} else if value, ok := line.Params["z"]; ok {
			// z-offset
			line.Params["z"] = value + usedFirstLayerValues.ZOffset
			line.Raw = ""
			output += line.String() + EOL
			return nil
		}
		output += line.Raw + EOL
		return nil
	})

	if err != nil {
		return err
	}

	// write g-code to outPath file
	if _, err := writer.WriteString(output + EOL); err != nil {
		return err
	}

	usedFirstLayerValues.Save(outPath + ".firstLayerResults")
	return nil
}
