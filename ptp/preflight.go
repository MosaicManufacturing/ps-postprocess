package ptp

import (
	"math"
	"mosaicmfg.com/ps-postprocess/gcode"
	"strconv"
	"strings"
)

type ptpPreflight struct {
	minFeedrate       float32
	maxFeedrate       float32
	minTemperature    float32
	maxTemperature    float32
	minLayerThickness float32
	maxLayerThickness float32
}

func toolpathPreflight(inpath string) (ptpPreflight, error) {
	minFeedrate := float32(math.Inf(1))
	maxFeedrate := float32(math.Inf(-1))
	minTemperature := float32(math.Inf(1))
	maxTemperature := float32(math.Inf(-1))
	minLayerThickness := float32(math.Inf(1))
	maxLayerThickness := float32(math.Inf(-1))
	currentFeedrate := float32(0)

	err := gcode.ReadByLine(inpath, func(line gcode.Command, _ int) error {
		if line.IsLinearMove() {
			// feedrates
			if f, ok := line.Params["f"]; ok {
				currentFeedrate = f
			}
			if _, ok := line.Params["e"]; ok {
				hasMovement := false
				if _, ok := line.Params["x"]; ok {
					hasMovement = true
				}
				if _, ok := line.Params["y"]; ok {
					hasMovement = true
				}
				if hasMovement {
					if currentFeedrate < minFeedrate {
						minFeedrate = currentFeedrate
					}
					if currentFeedrate > maxFeedrate {
						maxFeedrate = currentFeedrate
					}
				}
			}
		} else if line.Command == "M104" || line.Command == "M109" {
			// temperatures
			if temp, ok := line.Params["s"]; ok {
				if temp < minTemperature {
					minTemperature = temp
				}
				if temp > maxTemperature {
					maxTemperature = temp
				}
			}
		} else if line.Comment != "" && strings.HasPrefix(line.Comment, "HEIGHT:") {
			// layer heights
			height, err := strconv.ParseFloat(line.Comment[7:], 64)
			if err != nil {
				return err
			}
			height32 := roundZ(float32(height))
			if height32 < minLayerThickness {
				minLayerThickness = height32
			}
			if height32 > maxLayerThickness {
				maxLayerThickness = height32
			}
		}
		return nil
	})
	if err != nil {
		return ptpPreflight{}, err
	}
	results := ptpPreflight{
		minFeedrate:       minFeedrate,
		maxFeedrate:       maxFeedrate,
		minTemperature:    minTemperature,
		maxTemperature:    maxTemperature,
		minLayerThickness: minLayerThickness,
		maxLayerThickness: maxLayerThickness,
	}
	return results, err
}
