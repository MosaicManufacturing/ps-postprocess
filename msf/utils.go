package msf

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"mosaicmfg.com/ps-postprocess/gcode"
	"os"
	"path"
	"strings"
)

func roundTo(value float32, maxDecimalPlaces int) float32 {
	multiplier := math.Pow(10, float64(maxDecimalPlaces))
	rounded := math.Round(float64(value)*multiplier) / multiplier
	return float32(rounded)
}

func intToHexString(value uint, minHexDigits int) string {
	return fmt.Sprintf("%0*x", minHexDigits, value)
}

func int16ToHexString(value int16) string {
	return intToHexString(uint(uint16(value)), 4)
}

func floatToHexString(value float32) string {
	bits := math.Float32bits(value)
	buf := make([]byte, 4)
	binary.BigEndian.PutUint32(buf, bits)
	asUint := uint(binary.BigEndian.Uint32(buf))
	return intToHexString(asUint, 8)
}

func replaceSpaces(input string) string {
	return strings.ReplaceAll(input, " ", "_")
}

func truncate(input string, length int) string {
	if len(input) <= length {
		return input
	}
	return input[:length]
}

func msfVersionToO21(major, minor uint) string {
	versionNumber := (major * 10) + minor
	return fmt.Sprintf("O21 D%s%s", intToHexString(versionNumber, 4), EOL)
}

func writeLine(writer *bufio.Writer, line string) error {
	_, err := writer.WriteString(line + EOL)
	return err
}

func writeLines(writer *bufio.Writer, lines string) error {
	_, err := writer.WriteString(lines)
	return err
}

func getLineLength(x1, y1, x2, y2 float32) float32 {
	dx := float64(x2 - x1)
	dy := float64(y2 - y1)
	return float32(math.Sqrt(dx*dx + dy*dy))
}

func estimateMoveTime(x1, y1, x2, y2, feedrate float32) float32 {
	lineLength := getLineLength(x1, y1, x2, y2)
	mmPerS := feedrate / 60
	return lineLength / mmPerS
}

func estimateZMoveTime(z1, z2, feedrate float32) float32 {
	lineLength := float32(math.Abs(float64(z2 - z1)))
	mmPerS := feedrate / 60
	return lineLength / mmPerS
}

func estimatePurgeTime(eDelta, feedrate float32) float32 {
	mmPerS := feedrate / 60
	return float32(math.Abs(float64(eDelta))) / mmPerS
}

func lerp(minVal, maxVal, t float32) float32 {
	boundedT := float32(math.Max(0, math.Min(1, float64(t))))
	return ((1 - boundedT) * minVal) + (t * maxVal)
}

const filamentRadius = 1.75 / 2
const filamentPiRSquared = math.Pi * filamentRadius * filamentRadius

func filamentLengthToVolume(length float32) float32 {
	// V = Pi * (r^2) * h
	return filamentPiRSquared * length
}

func filamentVolumeToLength(volume float32) float32 {
	// h = V / (Pi * (r^2))
	return volume / filamentPiRSquared
}

func getExtrusionVolume(extrusionWidth, layerHeight, length float32) float32 {
	// https://manual.slic3r.org/advanced/flow-math -- "Extruding on top of a surface"
	area := (extrusionWidth-layerHeight)*layerHeight + math.Pi*float32(math.Pow(float64(layerHeight)/2, 2))
	return area * length
}

func getExtrusionLength(extrusionWidth, layerHeight, length float32) float32 {
	volume := getExtrusionVolume(extrusionWidth, layerHeight, length)
	return filamentVolumeToLength(volume)
}

// getTowerExtrusionCorrectionFactor returns the difference, as a ratio, between the
// naive estimate of a toolpath's cross-sectional area that treats it as a rectangle
// (width * height) and the more accurate, smaller value that treats it as a rectangle
// with semicircular end caps (the formula below).
func getTowerExtrusionCorrectionFactor(extrusionWidth, layerHeight float32) float64 {
	var x, y float64
	if extrusionWidth >= layerHeight {
		x, y = float64(extrusionWidth), float64(layerHeight)
	} else {
		x, y = float64(layerHeight), float64(extrusionWidth)
	}
	naiveArea := x * y
	informedArea := y*(x-y) + y*y*math.Pi/4
	return naiveArea / informedArea
}

func getPrintSummary(msf *MSF, timeEstimate float32) string {
	totalFilament := msf.GetTotalFilamentLength()
	filamentByDrive := msf.GetFilamentLengthsByDrive()

	summary := "; According to Chroma:" + EOL

	// total filament length
	summary += fmt.Sprintf("; filament total [mm] = %.5f%s", totalFilament, EOL)

	// filament lengths by drive
	for drive, length := range filamentByDrive {
		if length > 0 {
			summary += fmt.Sprintf(";    T%d filament = %.5f%s", drive+1, length, EOL)
		}
	}

	// time estimate
	summary += fmt.Sprintf("; estimated printing time = %s%s", gcode.GetTimeString(timeEstimate), EOL)
	summary += EOL

	return summary
}

func prependFile(filepath, content string) error {
	// create a temporary file
	tempfile, err := ioutil.TempFile(path.Dir(filepath), "")
	if err != nil {
		return err
	}

	// write prepended content first
	if _, err := tempfile.WriteString(content); err != nil {
		return err
	}

	// now append content of original file
	reader, err := os.Open(filepath)
	if err != nil {
		return err
	}
	if _, err := io.Copy(tempfile, reader); err != nil {
		return err
	}
	if err := reader.Close(); err != nil {
		return err
	}

	// finalize and close temporary file
	if err := tempfile.Sync(); err != nil {
		return err
	}
	if err := tempfile.Close(); err != nil {
		return err
	}

	// overwrite original file with temporary one
	return os.Rename(tempfile.Name(), filepath)
}
