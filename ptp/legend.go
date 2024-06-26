package ptp

import (
	"encoding/json"
	"fmt"
	"math"
	"sort"
)

const (
	maxDecimalsFeedrate    = 1
	maxDecimalsFanSpeed    = 0
	maxDecimalsTemperature = 1
	maxDecimalsLayerHeight = 4
)

type bufferData struct {
	Offset uint32 `json:"offset"`
	Size   uint32 `json:"size"`
}

type legendHeader struct {
	Version          int        `json:"version"`
	Position         bufferData `json:"position"`
	Normal           bufferData `json:"normal"`
	Index            bufferData `json:"index"`
	ExtrusionWidth   bufferData `json:"extrusionWidth"`
	LayerHeight      bufferData `json:"layerHeight"`
	IsTravel         bufferData `json:"isTravel"`
	ToolColor        bufferData `json:"toolColor"`
	PathTypeColor    bufferData `json:"pathTypeColor"`
	FeedrateColor    bufferData `json:"feedrateColor"`
	FanSpeedColor    bufferData `json:"fanSpeedColor"`
	TemperatureColor bufferData `json:"temperatureColor"`
	LayerHeightColor bufferData `json:"layerHeightColor"`
}

func (w *Writer) getLegendHeader() legendHeader {
	header := legendHeader{
		Version:          int(w.version),
		Position:         bufferData{Offset: 0, Size: w.bufferSizes["position"]},
		Normal:           bufferData{Offset: 0, Size: w.bufferSizes["normal"]},
		Index:            bufferData{Offset: 0, Size: w.bufferSizes["index"]},
		ExtrusionWidth:   bufferData{Offset: 0, Size: w.bufferSizes["extrusionWidth"]},
		LayerHeight:      bufferData{Offset: 0, Size: w.bufferSizes["layerHeight"]},
		IsTravel:         bufferData{Offset: 0, Size: w.bufferSizes["isTravel"]},
		ToolColor:        bufferData{Offset: 0, Size: w.bufferSizes["toolColor"]},
		PathTypeColor:    bufferData{Offset: 0, Size: w.bufferSizes["pathTypeColor"]},
		FeedrateColor:    bufferData{Offset: 0, Size: w.bufferSizes["feedrateColor"]},
		FanSpeedColor:    bufferData{Offset: 0, Size: w.bufferSizes["fanSpeedColor"]},
		TemperatureColor: bufferData{Offset: 0, Size: w.bufferSizes["temperatureColor"]},
		LayerHeightColor: bufferData{Offset: 0, Size: w.bufferSizes["layerHeightColor"]},
	}
	offset := headerSize
	header.Position.Offset = offset
	offset += w.bufferSizes["position"]
	header.Normal.Offset = offset
	offset += w.bufferSizes["normal"]
	header.Index.Offset = offset
	offset += w.bufferSizes["index"]
	header.ExtrusionWidth.Offset = offset
	offset += w.bufferSizes["extrusionWidth"]
	header.LayerHeight.Offset = offset
	offset += w.bufferSizes["layerHeight"]
	header.IsTravel.Offset = offset
	return header
}

type legendColors struct {
	MinFeedrateColor    [3]float32 `json:"minFeedrateColor"`
	MaxFeedrateColor    [3]float32 `json:"maxFeedrateColor"`
	MinFanSpeedColor    [3]float32 `json:"minFanSpeedColor"`
	MaxFanSpeedColor    [3]float32 `json:"maxFanSpeedColor"`
	MinTemperatureColor [3]float32 `json:"minTemperatureColor"`
	MaxTemperatureColor [3]float32 `json:"maxTemperatureColor"`
	MinLayerHeightColor [3]float32 `json:"minLayerHeightColor"`
	MaxLayerHeightColor [3]float32 `json:"maxLayerHeightColor"`
}

func getLegendColors() legendColors {
	return legendColors{
		MinFeedrateColor:    feedrateColorMin,
		MaxFeedrateColor:    feedrateColorMax,
		MinFanSpeedColor:    fanColorMin,
		MaxFanSpeedColor:    fanColorMax,
		MinTemperatureColor: temperatureColorMin,
		MaxTemperatureColor: temperatureColorMax,
		MinLayerHeightColor: layerHeightColorMin,
		MaxLayerHeightColor: layerHeightColorMax,
	}
}

type legendEntry struct {
	Label string
	Color string
}

func (l *legendEntry) MarshalJSON() ([]byte, error) {
	arr := []interface{}{l.Label, l.Color}
	return json.Marshal(arr)
}

type ptpLegend struct {
	Header                  legendHeader  `json:"header"`                  // header data (version, buffer offsets and sizes)
	Colors                  legendColors  `json:"colors"`                  // max/min colors for interpolated coloring
	Tool                    []legendEntry `json:"tool"`                    // legend of tools seen
	PathType                []legendEntry `json:"pathType"`                // legend of path types seen
	Feedrate                []legendEntry `json:"feedrate"`                // legend of feedrates -- needs gradation
	FanSpeed                []legendEntry `json:"fanSpeed"`                // legend of fan speeds -- possible gradation
	Temperature             []legendEntry `json:"temperature"`             // legend of temperatures -- needs gradation
	LayerHeight             []legendEntry `json:"layerHeight"`             // legend of layer heights -- needs gradation
	ZValues                 []float32     `json:"zValues"`                 // Z values for UI sliders
	LayerStartIndices       []uint32      `json:"layerStartIndices"`       // index values for rendering layer ranges
	LayerStartTravelIndices []uint32      `json:"layerStartTravelIndices"` // index values for rendering layer ranges
	HasPings                bool          `json:"hasPings"`                // for UI to show the relevant option
}

func removeDuplicateLegendEntries(legend []legendEntry) []legendEntry {
	labelsSeen := make(map[string]bool)
	uniqueLegend := make([]legendEntry, 0)
	for _, entry := range legend {
		if _, seen := labelsSeen[entry.Label]; !seen {
			labelsSeen[entry.Label] = true
			uniqueLegend = append(uniqueLegend, entry)
		}
	}
	return uniqueLegend
}

func (w *Writer) getToolLegend() []legendEntry {
	toolsSeen := setToSlice(w.state.toolsSeen, sort.Ints)
	legend := make([]legendEntry, 0, len(w.state.toolsSeen))
	for _, tool := range toolsSeen {
		legend = append(legend, legendEntry{
			Label: fmt.Sprintf("Tool %d", tool),
			Color: floatsToHex(w.toolColors[tool][0], w.toolColors[tool][1], w.toolColors[tool][2]),
		})
	}
	return legend
}

func (w *Writer) getPathTypeLegend() []legendEntry {
	legend := make([]legendEntry, 0)
	for i := PathType(0); i < pathTypeCount; i++ {
		if _, ok := w.state.pathTypesSeen[i]; ok {
			name := pathTypeNames[i]
			if i == PathTypeBrim {
				if w.brimIsSkirt {
					name = "Skirt"
				} else {
					name = "Brim"
				}
			}
			legend = append(legend, legendEntry{
				Label: name,
				Color: pathTypeColorStrings[i],
			})
		}
	}
	return legend
}

func (w *Writer) getFeedrateLegend() []legendEntry {
	feedratesSeen := setToSlice(w.state.feedratesSeen, sortFloat32Slice)
	legend := make([]legendEntry, 0, len(feedratesSeen))
	if len(feedratesSeen) <= 6 {
		for _, feedrate := range feedratesSeen {
			t := (feedrate - w.minFeedrate) / (w.maxFeedrate - w.minFeedrate)
			r := lerp(feedrateColorMin[0], feedrateColorMax[0], t)
			g := lerp(feedrateColorMin[1], feedrateColorMax[1], t)
			b := lerp(feedrateColorMin[2], feedrateColorMax[2], t)
			legend = append(legend, legendEntry{
				Label: fmt.Sprintf("%s mm/min", prepareFloatForJSON(feedrate, maxDecimalsFeedrate)),
				Color: floatsToHex(r, g, b),
			})
		}
	} else {
		step := float32(math.Round(float64(w.maxFeedrate-w.minFeedrate) / 6))
		for i := 0; i < 6; i++ {
			feedrate := (float32(i) * step) + w.minFeedrate
			t := float32(i) / 5
			r := lerp(feedrateColorMin[0], feedrateColorMax[0], t)
			g := lerp(feedrateColorMin[1], feedrateColorMax[1], t)
			b := lerp(feedrateColorMin[2], feedrateColorMax[2], t)
			legend = append(legend, legendEntry{
				Label: fmt.Sprintf("%s mm/min", prepareFloatForJSON(feedrate, maxDecimalsFeedrate)),
				Color: floatsToHex(r, g, b),
			})
		}
		legend = append(legend, legendEntry{
			Label: fmt.Sprintf("%s mm/min", prepareFloatForJSON(w.maxFeedrate, maxDecimalsFeedrate)),
			Color: floatsToHex(feedrateColorMax[0], feedrateColorMax[1], feedrateColorMax[2]),
		})
	}
	// de-duplicate legend entries with labels that are identical after rounding
	return removeDuplicateLegendEntries(legend)
}

func (w *Writer) getFanSpeedLegend() []legendEntry {
	fanSpeedsSeen := setToSlice(w.state.fanSpeedsSeen, sort.Ints)
	legend := make([]legendEntry, 0, len(fanSpeedsSeen))
	if len(fanSpeedsSeen) == 1 && fanSpeedsSeen[0] == 0 {
		legend = append(legend, legendEntry{
			Label: "Off",
			Color: floatsToHex(fanColorMin[0], fanColorMin[1], fanColorMin[2]),
		})
	} else if len(fanSpeedsSeen) == 1 && fanSpeedsSeen[0] == 255 {
		legend = append(legend, legendEntry{
			Label: "On",
			Color: floatsToHex(fanColorMax[0], fanColorMax[1], fanColorMax[2]),
		})
	} else if len(fanSpeedsSeen) == 2 &&
		((fanSpeedsSeen[0] == 0 && fanSpeedsSeen[1] == 255) ||
			(fanSpeedsSeen[0] == 255 && fanSpeedsSeen[1] == 0)) {
		legend = append(legend, legendEntry{
			Label: "Off",
			Color: floatsToHex(fanColorMin[0], fanColorMin[1], fanColorMin[2]),
		}, legendEntry{
			Label: "On",
			Color: floatsToHex(fanColorMax[0], fanColorMax[1], fanColorMax[2]),
		})
	} else if len(fanSpeedsSeen) <= 6 {
		for _, pwmValue := range fanSpeedsSeen {
			t := float32(pwmValue) / 255
			percent := float32(math.Max(0, math.Min(100, math.Round(float64(pwmValue)*100)/255)))
			r := lerp(fanColorMin[0], fanColorMax[0], t)
			g := lerp(fanColorMin[1], fanColorMax[1], t)
			b := lerp(fanColorMin[2], fanColorMax[2], t)
			legend = append(legend, legendEntry{
				Label: fmt.Sprintf("%s%%", prepareFloatForJSON(percent, maxDecimalsFanSpeed)),
				Color: floatsToHex(r, g, b),
			})
		}
	} else {
		step := float32(math.Round(255 / 6))
		for i := 0; i < 6; i++ {
			pwmValue := float32(i) * step
			t := float32(i) / 5
			percent := float32(math.Round(float64(pwmValue)*100*10/255) / 10)
			r := lerp(fanColorMin[0], fanColorMax[0], t)
			g := lerp(fanColorMin[1], fanColorMax[1], t)
			b := lerp(fanColorMin[2], fanColorMax[2], t)
			legend = append(legend, legendEntry{
				Label: fmt.Sprintf("%s%%", prepareFloatForJSON(percent, maxDecimalsFanSpeed)),
				Color: floatsToHex(r, g, b),
			})
		}
		legend = append(legend, legendEntry{
			Label: "100%",
			Color: floatsToHex(fanColorMax[0], fanColorMax[1], fanColorMax[2]),
		})
	}
	// de-duplicate legend entries with labels that are identical after rounding
	return removeDuplicateLegendEntries(legend)
}

func (w *Writer) getTemperatureLegend() []legendEntry {
	temperaturesSeen := setToSlice(w.state.temperaturesSeen, sortFloat32Slice)
	legend := make([]legendEntry, 0, len(temperaturesSeen))
	if len(temperaturesSeen) <= 6 {
		for _, temperature := range temperaturesSeen {
			var r, g, b float32
			if w.maxTemperature == w.minTemperature {
				r = temperatureColorMax[0]
				g = temperatureColorMax[1]
				b = temperatureColorMax[2]
			} else {
				t := (temperature - w.minTemperature) / (w.maxTemperature - w.minTemperature)
				r = lerp(temperatureColorMin[0], temperatureColorMax[0], t)
				g = lerp(temperatureColorMin[1], temperatureColorMax[1], t)
				b = lerp(temperatureColorMin[2], temperatureColorMax[2], t)
			}
			legend = append(legend, legendEntry{
				Label: fmt.Sprintf("%s °C", prepareFloatForJSON(temperature, maxDecimalsTemperature)),
				Color: floatsToHex(r, g, b),
			})
		}
	} else {
		step := float32(math.Round(float64(w.maxTemperature-w.minTemperature) / 6))
		for i := 0; i < 6; i++ {
			temperature := (float32(i) * step) + w.minTemperature
			t := float32(i) / 5
			r := lerp(temperatureColorMin[0], temperatureColorMax[0], t)
			g := lerp(temperatureColorMin[1], temperatureColorMax[1], t)
			b := lerp(temperatureColorMin[2], temperatureColorMax[2], t)
			legend = append(legend, legendEntry{
				Label: fmt.Sprintf("%s °C", prepareFloatForJSON(temperature, maxDecimalsTemperature)),
				Color: floatsToHex(r, g, b),
			})
		}
		legend = append(legend, legendEntry{
			Label: fmt.Sprintf("%s °C", prepareFloatForJSON(w.maxTemperature, maxDecimalsTemperature)),
			Color: floatsToHex(temperatureColorMax[0], temperatureColorMax[1], temperatureColorMax[2]),
		})
	}
	// de-duplicate legend entries with labels that are identical after rounding
	return removeDuplicateLegendEntries(legend)
}

func (w *Writer) getLayerHeightLegend() []legendEntry {
	layerHeightsSeen := setToSlice(w.state.layerHeightsSeen, sortFloat32Slice)
	legend := make([]legendEntry, 0, len(layerHeightsSeen))
	if len(layerHeightsSeen) == 1 {
		legend = []legendEntry{
			{
				Label: fmt.Sprintf("%s mm", prepareFloatForJSON(layerHeightsSeen[0], maxDecimalsLayerHeight)),
				Color: floatsToHex(layerHeightColorMax[0], layerHeightColorMax[1], layerHeightColorMax[2]),
			},
		}
	} else if len(layerHeightsSeen) <= 6 {
		for _, layerHeight := range layerHeightsSeen {
			t := (layerHeight - w.minLayerHeight) / (w.maxLayerHeight - w.minLayerHeight)
			r := lerp(layerHeightColorMin[0], layerHeightColorMax[0], t)
			g := lerp(layerHeightColorMin[1], layerHeightColorMax[1], t)
			b := lerp(layerHeightColorMin[2], layerHeightColorMax[2], t)
			legend = append(legend, legendEntry{
				Label: fmt.Sprintf("%s mm", prepareFloatForJSON(layerHeight, maxDecimalsLayerHeight)),
				Color: floatsToHex(r, g, b),
			})
		}
	} else {
		step := float32(math.Round(float64(w.maxLayerHeight-w.minLayerHeight)*1000/6) / 1000)
		for i := 0; i < 6; i++ {
			layerHeight := (float32(i) * step) + w.minLayerHeight
			t := float32(i) / 5
			r := lerp(layerHeightColorMin[0], layerHeightColorMax[0], t)
			g := lerp(layerHeightColorMin[1], layerHeightColorMax[1], t)
			b := lerp(layerHeightColorMin[2], layerHeightColorMax[2], t)
			legend = append(legend, legendEntry{
				Label: fmt.Sprintf("%s mm", prepareFloatForJSON(layerHeight, maxDecimalsLayerHeight)),
				Color: floatsToHex(r, g, b),
			})
		}
		legend = append(legend, legendEntry{
			Label: fmt.Sprintf("%s mm", prepareFloatForJSON(w.maxLayerHeight, maxDecimalsLayerHeight)),
			Color: floatsToHex(layerHeightColorMax[0], layerHeightColorMax[1], layerHeightColorMax[2]),
		})
	}
	// de-duplicate legend entries with labels that are identical after rounding
	return removeDuplicateLegendEntries(legend)
}

func (w *Writer) getLegend() ([]byte, error) {
	legend := ptpLegend{
		Header:            w.getLegendHeader(),
		Colors:            getLegendColors(),
		Tool:              w.getToolLegend(),
		PathType:          w.getPathTypeLegend(),
		Feedrate:          w.getFeedrateLegend(),
		FanSpeed:          w.getFanSpeedLegend(),
		Temperature:       w.getTemperatureLegend(),
		LayerHeight:       w.getLayerHeightLegend(),
		ZValues:           w.state.layerHeights,
		LayerStartIndices: w.state.layerStartIndices,
		HasPings:          w.bufferSizes["pingPosition"] > 0,
	}
	return json.Marshal(legend)
}
