package ptp

import (
	"math"
	"strconv"
	"strings"

	"mosaicmfg.com/ps-postprocess/gcode"
)

type ptpPreflight struct {
	minFeedrate    float32
	maxFeedrate    float32
	minTemperature float32
	maxTemperature float32
	minLayerHeight float32
	maxLayerHeight float32
	minFanSpeed    float32
	maxFanSpeed    float32
}

func toolpathPreflight(inpath string) (ptpPreflight, error) {
	minFeedrate := float32(math.Inf(1))
	maxFeedrate := float32(math.Inf(-1))
	minTemperature := float32(math.Inf(1))
	maxTemperature := float32(math.Inf(-1))
	minLayerHeight := float32(math.Inf(1))
	maxLayerHeight := float32(math.Inf(-1))
	minFanSpeed := float32(math.Inf(1))
	maxFanSpeed := float32(math.Inf(-1))

	currentFeedrate := float32(0)
	var et gcode.ExtrusionTracker

	err := gcode.ReadByLine(inpath, func(line gcode.Command, _ int) error {
		et.TrackInstruction(line)
		if line.IsLinearOrArcMove() {
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
				if hasMovement && !et.LastExtrudeWasRetract {
					if currentFeedrate < minFeedrate {
						minFeedrate = currentFeedrate
					}
					if currentFeedrate > maxFeedrate {
						maxFeedrate = currentFeedrate
					}
				}
			}
		} else if line.Command == "M104" {
			// temperatures
			if temp, ok := line.Params["s"]; ok && temp > 0 {
				if temp < minTemperature {
					minTemperature = temp
				}
				if temp > maxTemperature {
					maxTemperature = temp
				}
			}
		} else if line.Command == "M109" {
			// temperatures
			if temp, ok := line.Params["s"]; ok && temp > 0 {
				if temp < minTemperature {
					minTemperature = temp
				}
				if temp > maxTemperature {
					maxTemperature = temp
				}
			} else if temp, ok = line.Params["r"]; ok && temp > 0 {
				if temp < minTemperature {
					minTemperature = temp
				}
				if temp > maxTemperature {
					maxTemperature = temp
				}
			}
		} else if line.Command == "M106" {
			// fan speed
			// ignore P10, which is specifically assigned to the cooling module
			if fanIndex, ok := line.Params["p"]; !ok || fanIndex != 10 {
				if pwm, ok := line.Params["s"]; ok {
					if pwm < minFanSpeed {
						minFanSpeed = pwm
					}
					if pwm > maxFanSpeed {
						maxFanSpeed = pwm
					}
				}
			}
		} else if line.Command == "M107" {
			// fan speed
			// ignore P10, which is specifically assigned to the cooling module
			if fanIndex, ok := line.Params["p"]; !ok || fanIndex != 10 {
				minFanSpeed = 0
				if maxFanSpeed < 0 {
					maxFanSpeed = 0
				}
			}
		} else if line.Comment != "" && strings.HasPrefix(line.Comment, "HEIGHT:") {
			// layer heights
			height, err := strconv.ParseFloat(line.Comment[7:], 64)
			if err != nil {
				return err
			}
			height32 := roundZ(float32(height))
			if height32 < minLayerHeight {
				minLayerHeight = height32
			}
			if height32 > maxLayerHeight {
				maxLayerHeight = height32
			}
		}
		return nil
	})
	if err != nil {
		return ptpPreflight{}, err
	}
	results := ptpPreflight{
		minFeedrate:    minFeedrate,
		maxFeedrate:    maxFeedrate,
		minTemperature: minTemperature,
		maxTemperature: maxTemperature,
		minLayerHeight: minLayerHeight,
		maxLayerHeight: maxLayerHeight,
		minFanSpeed:    minFanSpeed,
		maxFanSpeed:    maxFanSpeed,
	}
	return results, err
}
