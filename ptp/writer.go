package ptp

import (
	"bufio"
	"errors"
	"fmt"
	"io/ioutil"
	"math"
	"os"
)

type writerState struct {
	lastLineWasPrint       bool    // if true, corner triangles will be created to join the lines
	printLineBuffered      bool    // if true, a print line from prevX/Y/Z to currentX/Y/Z needs to be output
	transitionLineBuffered bool    // if true, use `bufferedFromTool` to create a gradient
	bufferedFromTool       int     // used for interpolation when outputting a buffered transition line
	bufferedT              float32 // used for interpolation when outputting a buffered transition line
	inWipe                 bool    // if true, "wipe and retract" commands are occurring
	currentX               float32
	currentY               float32
	currentZ               float32
	prevX                  float32
	prevY                  float32
	prevZ                  float32
	currentExtrusionWidth  float32
	currentLayerHeight     float32
	currentTool            int // == travelTool when outputting travel moves!
	currentPathType        PathType
	currentFeedrate        float32
	currentFanSpeed        int
	currentTemperature     float32

	// buffered data while adding travel paths (to be able to revert to the previous values after)
	travelling                   bool
	travelBufferedPathType       PathType
	travelBufferedExtrusionWidth float32
	travelBufferedLayerHeight    float32
	travelBufferedTool           int

	// track the display height of each layer (for UI sliders to use), but for each geometry
	// also track the actual start index of the data for this layer, so that we can render
	// a subset of the geometry based on index ranges rather than Z-height clipping.
	//
	// for a print with N layers, 0..N-1,
	// - layer 0 is the start sequence
	// - layer 1 is the first of the print
	// - ...
	// - layer N is the last layer of the print (+ end sequence)
	layerHeights      []float32 // [0] == 0, [1] == first layer height, [N + 1] == [N]
	layerStartIndices []uint32  // index of first vertex in layer

	// sets used to track unique values seen, for generating the legend
	toolsSeen        map[int]bool
	pathTypesSeen    map[PathType]bool
	feedratesSeen    map[float32]bool
	fanSpeedsSeen    map[int]bool
	temperaturesSeen map[float32]bool
	layerHeightsSeen map[float32]bool
}

func getStartingWriterState(initialExtrusionWidth, initialLayerHeight float32) writerState {
	return writerState{
		currentExtrusionWidth: initialExtrusionWidth,
		currentLayerHeight:    initialLayerHeight,
		layerHeights:          []float32{0}, // initial state is "in the start sequence"
		layerStartIndices:     []uint32{0},
		toolsSeen:             make(map[int]bool),
		pathTypesSeen:         make(map[PathType]bool),
		feedratesSeen:         make(map[float32]bool),
		fanSpeedsSeen:         make(map[int]bool),
		temperaturesSeen:      make(map[float32]bool),
		layerHeightsSeen:      make(map[float32]bool),
	}
}

type Writer struct {
	version     uint8
	paths       map[string]string
	files       map[string]*os.File
	writers     map[string]*bufio.Writer
	bufferSizes map[string]uint32

	// bounds for interpolated color scales
	minFeedrate    float32
	maxFeedrate    float32
	minTemperature float32
	maxTemperature float32
	minLayerHeight float32
	maxLayerHeight float32

	brimIsSkirt bool         // if true, PathTypeBrim will be referred to as Skirt
	toolColors  [][3]float32 // array of [r, g, b] floats in range 0..1
	state       writerState
}

func NewWriter(outpath string, initialExtrusionWidth, initialLayerHeight float32, brimIsSkirt bool, toolColors [][3]float32) Writer {
	return Writer{
		version: ptpVersion,
		paths: map[string]string{
			"main":             outpath,
			"legend":           fmt.Sprintf("%s.%s", outpath, "legend"),
			"normal":           fmt.Sprintf("%s.%s", outpath, "normal"),
			"index":            fmt.Sprintf("%s.%s", outpath, "index"),
			"extrusionWidth":   fmt.Sprintf("%s.%s", outpath, "extrusionWidth"),
			"layerHeight":      fmt.Sprintf("%s.%s", outpath, "layerHeight"),
			"isTravel":         fmt.Sprintf("%s.%s", outpath, "isTravel"),
			"retractPosition":  fmt.Sprintf("%s.%s", outpath, "retractPosition"),
			"indexAtRetract":   fmt.Sprintf("%s.%s", outpath, "indexAtRetract"),
			"restartPosition":  fmt.Sprintf("%s.%s", outpath, "restartPosition"),
			"indexAtRestart":   fmt.Sprintf("%s.%s", outpath, "indexAtRestart"),
			"pingPosition":     fmt.Sprintf("%s.%s", outpath, "pingPosition"),
			"indexAtPing":      fmt.Sprintf("%s.%s", outpath, "indexAtPing"),
			"toolColor":        fmt.Sprintf("%s.%s", outpath, "toolColor"),
			"pathTypeColor":    fmt.Sprintf("%s.%s", outpath, "pathTypeColor"),
			"feedrateColor":    fmt.Sprintf("%s.%s", outpath, "feedrateColor"),
			"fanSpeedColor":    fmt.Sprintf("%s.%s", outpath, "fanSpeedColor"),
			"temperatureColor": fmt.Sprintf("%s.%s", outpath, "temperatureColor"),
			"layerHeightColor": fmt.Sprintf("%s.%s", outpath, "layerHeightColor"),
		},
		files: map[string]*os.File{
			"main":             nil,
			"normal":           nil,
			"index":            nil,
			"extrusionWidth":   nil,
			"layerHeight":      nil,
			"isTravel":         nil,
			"retractPosition":  nil,
			"indexAtRetract":   nil,
			"restartPosition":  nil,
			"indexAtRestart":   nil,
			"pingPosition":     nil,
			"indexAtPing":      nil,
			"toolColor":        nil,
			"pathTypeColor":    nil,
			"feedrateColor":    nil,
			"fanSpeedColor":    nil,
			"temperatureColor": nil,
			"layerHeightColor": nil,
		},
		writers: map[string]*bufio.Writer{
			"main":             nil,
			"normal":           nil,
			"index":            nil,
			"extrusionWidth":   nil,
			"layerHeight":      nil,
			"isTravel":         nil,
			"retractPosition":  nil,
			"indexAtRetract":   nil,
			"restartPosition":  nil,
			"indexAtRestart":   nil,
			"pingPosition":     nil,
			"indexAtPing":      nil,
			"toolColor":        nil,
			"pathTypeColor":    nil,
			"feedrateColor":    nil,
			"fanSpeedColor":    nil,
			"temperatureColor": nil,
			"layerHeightColor": nil,
		},
		bufferSizes: map[string]uint32{
			"position":         0,
			"normal":           0,
			"index":            0,
			"extrusionWidth":   0,
			"layerHeight":      0,
			"isTravel":         0,
			"retractPosition":  0,
			"indexAtRetract":   0,
			"restartPosition":  0,
			"indexAtRestart":   0,
			"pingPosition":     0,
			"indexAtPing":      0,
			"toolColor":        0,
			"pathTypeColor":    0,
			"feedrateColor":    0,
			"fanSpeedColor":    0,
			"temperatureColor": 0,
			"layerHeightColor": 0,
		},
		minFeedrate:    0,
		maxFeedrate:    0,
		minTemperature: 0,
		maxTemperature: 0,
		minLayerHeight: 0,
		maxLayerHeight: 0,
		brimIsSkirt:    brimIsSkirt,
		toolColors:     toolColors,
		state:          getStartingWriterState(initialExtrusionWidth, initialLayerHeight),
	}
}

func (w *Writer) SetFeedrateBounds(min, max float32) {
	w.minFeedrate = min
	w.maxFeedrate = max
}

func (w *Writer) SetTemperatureBounds(min, max float32) {
	w.minTemperature = min
	w.maxTemperature = max
}

func (w *Writer) SetLayerHeightBounds(min, max float32) {
	w.minLayerHeight = min
	w.maxLayerHeight = max
}

func (w *Writer) Initialize() error {
	if w.maxFeedrate < w.minFeedrate || w.minFeedrate < 0 || w.maxFeedrate <= 0 {
		return errors.New("invalid feedrate bounds for creating legend")
	}
	if w.maxTemperature < w.minTemperature || w.minTemperature < 0 || w.maxTemperature <= 0 {
		return errors.New("invalid temperature bounds for creating legend")
	}
	if w.maxLayerHeight < w.minLayerHeight || w.minLayerHeight < 0 || w.maxLayerHeight <= 0 {
		return errors.New("invalid layer height bounds for creating legend")
	}

	filenamesToOpen := []string{
		"main",
		"normal",
		"index",
		"extrusionWidth",
		"layerHeight",
		"isTravel",
		"retractPosition",
		"indexAtRetract",
		"restartPosition",
		"indexAtRestart",
		"pingPosition",
		"indexAtPing",
		"toolColor",
		"pathTypeColor",
		"feedrateColor",
		"fanSpeedColor",
		"temperatureColor",
		"layerHeightColor",
	}
	for _, filename := range filenamesToOpen {
		if err := openForWrite(w, filename); err != nil {
			return err
		}
	}
	return w.writeHeader()
}

func (w *Writer) writeHeader() error {
	buf := make([]byte, headerSize)
	buf[0] = w.version // only first byte of header is used
	_, err := w.writers["main"].Write(buf)
	return err
}

func (w *Writer) flushLineBuffers() error {
	if w.state.printLineBuffered {
		if err := w.outputPrintLine(); err != nil {
			return err
		}
		w.state.printLineBuffered = false
	}
	return nil
}

func (w *Writer) Finalize() error {
	// flush any remaining buffers
	if err := w.flushLineBuffers(); err != nil {
		return err
	}

	w.updateLayerStartIndices()

	// close the temp files
	filenamesToClose := []string{
		"normal",
		"index",
		"extrusionWidth",
		"layerHeight",
		"isTravel",
		"retractPosition",
		"indexAtRetract",
		"restartPosition",
		"indexAtRestart",
		"pingPosition",
		"indexAtPing",
		"toolColor",
		"pathTypeColor",
		"feedrateColor",
		"fanSpeedColor",
		"temperatureColor",
		"layerHeightColor",
	}
	for _, filename := range filenamesToClose {
		if err := flushAndClose(w, filename); err != nil {
			return err
		}
	}

	// concatenate the files
	filenamesToConcatenate := []string{
		"normal",
		"index",
		"extrusionWidth",
		"layerHeight",
		"isTravel",
	}
	for _, filename := range filenamesToConcatenate {
		if err := concatOntoWriter(w, "main", filename); err != nil {
			return err
		}
	}

	// write legend and commit main file
	if err := w.saveLegend(); err != nil {
		return err
	}
	return flushAndClose(w, "main")
}

func (w *Writer) saveLegend() error {
	legend, err := w.getLegend()
	if err != nil {
		return err
	}
	return ioutil.WriteFile(w.paths["legend"], legend, 0644)
}

func (w *Writer) writePosition(x, y, z float32) error {
	if err := writeFloat32LE(w.writers["main"], x); err != nil {
		return err
	}
	if err := writeFloat32LE(w.writers["main"], y); err != nil {
		return err
	}
	if err := writeFloat32LE(w.writers["main"], z); err != nil {
		return err
	}
	w.bufferSizes["position"] += floatBytes * 3
	return nil
}

func (w *Writer) writeNormal(x, y, z float32) error {
	if err := writeFloat32LE(w.writers["normal"], x); err != nil {
		return err
	}
	if err := writeFloat32LE(w.writers["normal"], y); err != nil {
		return err
	}
	if err := writeFloat32LE(w.writers["normal"], z); err != nil {
		return err
	}
	w.bufferSizes["normal"] += floatBytes * 3
	return nil
}

func (w *Writer) writeIndex(idx uint32) error {
	if err := writeUint32LE(w.writers["index"], idx); err != nil {
		return err
	}
	w.bufferSizes["index"] += uint32Bytes
	return nil
}

func (w *Writer) writeExtrusionWidth(width float32) error {
	if err := writeFloat32LE(w.writers["extrusionWidth"], width); err != nil {
		return err
	}
	w.bufferSizes["extrusionWidth"] += floatBytes
	return nil
}

func (w *Writer) writeLayerHeight(height float32) error {
	if err := writeFloat32LE(w.writers["layerHeight"], height); err != nil {
		return err
	}
	w.bufferSizes["layerHeight"] += floatBytes
	return nil
}

func (w *Writer) writeIsTravel(isTravel bool) error {
	val := uint8(0)
	if isTravel {
		val = 1
	}
	if err := writeUint8(w.writers["isTravel"], val); err != nil {
		return err
	}
	w.bufferSizes["isTravel"] += uint8Bytes
	return nil
}

func (w *Writer) writeRetractPosition(x, y, z float32) error {
	if err := writeFloat32LE(w.writers["retractPosition"], x); err != nil {
		return err
	}
	if err := writeFloat32LE(w.writers["retractPosition"], y); err != nil {
		return err
	}
	if err := writeFloat32LE(w.writers["retractPosition"], z); err != nil {
		return err
	}
	w.bufferSizes["retractPosition"] += floatBytes * 3
	return nil
}

func (w *Writer) writeIndexAtRetract(idx uint32) error {
	if err := writeFloat32LE(w.writers["indexAtRetract"], float32(idx)); err != nil {
		return err
	}
	w.bufferSizes["indexAtRetract"] += floatBytes
	return nil
}

func (w *Writer) writeRestartPosition(x, y, z float32) error {
	if err := writeFloat32LE(w.writers["restartPosition"], x); err != nil {
		return err
	}
	if err := writeFloat32LE(w.writers["restartPosition"], y); err != nil {
		return err
	}
	if err := writeFloat32LE(w.writers["restartPosition"], z); err != nil {
		return err
	}
	w.bufferSizes["restartPosition"] += floatBytes * 3
	return nil
}

func (w *Writer) writeIndexAtRestart(idx uint32) error {
	if err := writeFloat32LE(w.writers["indexAtRestart"], float32(idx)); err != nil {
		return err
	}
	w.bufferSizes["indexAtRestart"] += floatBytes
	return nil
}

func (w *Writer) writePingPosition(x, y, z float32) error {
	if err := writeFloat32LE(w.writers["pingPosition"], x); err != nil {
		return err
	}
	if err := writeFloat32LE(w.writers["pingPosition"], y); err != nil {
		return err
	}
	if err := writeFloat32LE(w.writers["pingPosition"], z); err != nil {
		return err
	}
	w.bufferSizes["pingPosition"] += floatBytes * 3
	return nil
}

func (w *Writer) writeIndexAtPing(idx uint32) error {
	if err := writeFloat32LE(w.writers["indexAtPing"], float32(idx)); err != nil {
		return err
	}
	w.bufferSizes["indexAtPing"] += floatBytes
	return nil
}

func (w *Writer) writeToolColor(toTool, fromTool int, t float32) error {
	var r, g, b float32
	if toTool < 0 {
		r = travelColor[0]
		g = travelColor[1]
		b = travelColor[2]
	} else if t >= 1 {
		r = w.toolColors[toTool][0]
		g = w.toolColors[toTool][1]
		b = w.toolColors[toTool][2]
	} else if t <= 0 {
		r = w.toolColors[fromTool][0]
		g = w.toolColors[fromTool][1]
		b = w.toolColors[fromTool][2]
	} else {
		r = lerp(w.toolColors[fromTool][0], w.toolColors[toTool][0], t)
		g = lerp(w.toolColors[fromTool][1], w.toolColors[toTool][1], t)
		b = lerp(w.toolColors[fromTool][2], w.toolColors[toTool][2], t)
	}
	if err := writeFloat32LE(w.writers["toolColor"], r); err != nil {
		return err
	}
	if err := writeFloat32LE(w.writers["toolColor"], g); err != nil {
		return err
	}
	if err := writeFloat32LE(w.writers["toolColor"], b); err != nil {
		return err
	}
	w.bufferSizes["toolColor"] += floatBytes * 3
	return nil
}

func (w *Writer) writePathTypeColor(pathType PathType) error {
	if err := writeFloat32LE(w.writers["pathTypeColor"], pathTypeColors[pathType][0]); err != nil {
		return err
	}
	if err := writeFloat32LE(w.writers["pathTypeColor"], pathTypeColors[pathType][1]); err != nil {
		return err
	}
	if err := writeFloat32LE(w.writers["pathTypeColor"], pathTypeColors[pathType][2]); err != nil {
		return err
	}
	w.bufferSizes["pathTypeColor"] += floatBytes * 3
	return nil
}

func (w *Writer) writeFeedrateColor(feedrate float32) error {
	t := float32(1)
	if w.maxFeedrate > w.minFeedrate {
		t = (feedrate - w.minFeedrate) / (w.maxFeedrate - w.minFeedrate)
	}
	if err := writeFloat32LE(w.writers["feedrateColor"], t); err != nil {
		return err
	}
	w.bufferSizes["feedrateColor"] += floatBytes
	return nil
}

func (w *Writer) writeFanSpeedColor(pwmValue int) error {
	t := float32(pwmValue) / 255
	if err := writeFloat32LE(w.writers["fanSpeedColor"], t); err != nil {
		return err
	}
	w.bufferSizes["fanSpeedColor"] += floatBytes
	return nil
}

func (w *Writer) writeTemperatureColor(temperature float32) error {
	t := float32(1)
	if w.maxTemperature > w.minTemperature {
		t = (temperature - w.minTemperature) / (w.maxTemperature - w.minTemperature)
	}
	if err := writeFloat32LE(w.writers["temperatureColor"], t); err != nil {
		return err
	}
	w.bufferSizes["temperatureColor"] += floatBytes
	return nil
}

func (w *Writer) writeLayerHeightColor(layerHeight float32) error {
	t := float32(1)
	if w.maxLayerHeight > w.minLayerHeight {
		t = (layerHeight - w.minLayerHeight) / (w.maxLayerHeight - w.minLayerHeight)
	}
	if err := writeFloat32LE(w.writers["layerHeightColor"], t); err != nil {
		return err
	}
	w.bufferSizes["layerHeightColor"] += floatBytes
	return nil
}

func (w *Writer) updateLayerStartIndices() {
	w.state.layerStartIndices = append(w.state.layerStartIndices, w.getCurrentIndex())
}

func (w *Writer) LayerChange(z float32) error {
	// add to the list of Z heights
	w.state.layerHeights = append(w.state.layerHeights, z)
	// set starting indices for geometry this layer
	w.updateLayerStartIndices()
	return nil
}

func (w *Writer) SetExtrusionWidth(width float32) error {
	if w.state.travelling {
		w.state.travelBufferedExtrusionWidth = width
		return nil
	}
	if width == w.state.currentExtrusionWidth {
		return nil
	}
	if w.state.printLineBuffered {
		if err := w.outputPrintLine(); err != nil {
			return err
		}
		w.state.printLineBuffered = false
		w.state.lastLineWasPrint = true
	}
	w.state.currentExtrusionWidth = width
	return nil
}

func (w *Writer) SetLayerHeight(height float32) error {
	if w.state.travelling {
		w.state.travelBufferedLayerHeight = height
		return nil
	}
	if height == w.state.currentLayerHeight {
		return nil
	}
	if w.state.printLineBuffered {
		if err := w.outputPrintLine(); err != nil {
			return err
		}
		w.state.printLineBuffered = false
		w.state.lastLineWasPrint = true
	}
	w.state.currentLayerHeight = height
	return nil
}

func (w *Writer) SetTool(tool int) error {
	if w.state.travelling {
		w.state.travelBufferedTool = tool
		return nil
	}
	if tool == w.state.currentTool {
		return nil
	}
	if w.state.printLineBuffered {
		if err := w.outputPrintLine(); err != nil {
			return err
		}
		w.state.printLineBuffered = false
		w.state.lastLineWasPrint = true
	}
	w.state.currentTool = tool
	return nil
}

func (w *Writer) SetPathType(pathType PathType) error {
	if w.state.travelling {
		w.state.travelBufferedPathType = pathType
		return nil
	}
	if pathType == w.state.currentPathType {
		return nil
	}
	if w.state.printLineBuffered {
		if err := w.outputPrintLine(); err != nil {
			return err
		}
		w.state.printLineBuffered = false
		w.state.lastLineWasPrint = true
	}
	w.state.currentPathType = pathType
	return nil
}

func (w *Writer) SetFeedrate(feedrate float32) error {
	if feedrate == w.state.currentFeedrate {
		return nil
	}
	if w.state.printLineBuffered {
		if err := w.outputPrintLine(); err != nil {
			return err
		}
		w.state.printLineBuffered = false
		w.state.lastLineWasPrint = true
	}
	w.state.currentFeedrate = feedrate
	return nil
}

func (w *Writer) SetFanSpeed(pwmValue int) error {
	if pwmValue == w.state.currentFanSpeed {
		return nil
	}
	if w.state.printLineBuffered {
		if err := w.outputPrintLine(); err != nil {
			return err
		}
		w.state.printLineBuffered = false
		w.state.lastLineWasPrint = true
	}
	w.state.currentFanSpeed = pwmValue
	return nil
}

func (w *Writer) SetTemperature(temperature float32) error {
	if temperature == w.state.currentTemperature {
		return nil
	}
	if w.state.printLineBuffered {
		if err := w.outputPrintLine(); err != nil {
			return err
		}
		w.state.printLineBuffered = false
		w.state.lastLineWasPrint = true
	}
	w.state.currentTemperature = temperature
	return nil
}

func (w *Writer) GetCurrentPosition() (float32, float32, float32) {
	return w.state.currentX, w.state.currentY, w.state.currentZ
}

func (w *Writer) outputRetractPoint() error {
	if err := w.writeRetractPosition(w.state.currentX, w.state.currentY, w.state.currentZ); err != nil {
		return err
	}
	return w.writeIndexAtRetract(w.getCurrentIndex())
}

func (w *Writer) outputRestartPoint() error {
	if err := w.writeRestartPosition(w.state.currentX, w.state.currentY, w.state.currentZ); err != nil {
		return err
	}
	return w.writeIndexAtRestart(w.getCurrentIndex())
}

func (w *Writer) outputPingPoint() error {
	if err := w.writePingPosition(w.state.currentX, w.state.currentY, w.state.currentZ); err != nil {
		return err
	}
	return w.writeIndexAtPing(w.getCurrentIndex())
}

func (w *Writer) AddXYZTravelTo(x, y, z float32) error {
	return w.addXYZPrintLineTo(x, y, z, true, true)
}

func (w *Writer) AddXYTravelTo(x, y float32) error {
	return w.AddXYZTravelTo(x, y, w.state.currentZ)
}

func (w *Writer) AddRetract() error {
	// flush print line buffer if necessary
	if w.state.printLineBuffered {
		if err := w.outputPrintLine(); err != nil {
			return err
		}
		w.state.printLineBuffered = false
		w.state.lastLineWasPrint = true
	}
	return w.outputRetractPoint()
}

func (w *Writer) AddRetractAt(x, y, z float32, savePosition bool) error {
	// flush print line buffer if necessary
	if w.state.printLineBuffered {
		if err := w.outputPrintLine(); err != nil {
			return err
		}
		w.state.printLineBuffered = false
		w.state.lastLineWasPrint = true
	}
	// update history
	w.state.prevX = w.state.currentX
	w.state.prevY = w.state.currentY
	w.state.prevZ = w.state.currentZ
	w.state.currentX = x
	w.state.currentY = y
	w.state.currentZ = z
	w.state.lastLineWasPrint = false

	if err := w.outputRetractPoint(); err != nil {
		return err
	}

	if !savePosition {
		w.state.currentX = w.state.prevX
		w.state.currentY = w.state.prevY
		w.state.currentZ = w.state.prevZ
	}

	return nil
}

func (w *Writer) AddRestart() error {
	// flush print line buffer if necessary
	if w.state.printLineBuffered {
		if err := w.outputPrintLine(); err != nil {
			return err
		}
		w.state.printLineBuffered = false
		w.state.lastLineWasPrint = true
	}
	return w.outputRestartPoint()
}

func (w *Writer) AddRestartAt(x, y, z float32, savePosition bool) error {
	// flush print line buffer if necessary
	if w.state.printLineBuffered {
		if err := w.outputPrintLine(); err != nil {
			return err
		}
		w.state.printLineBuffered = false
		w.state.lastLineWasPrint = true
	}
	// update history
	w.state.prevX = w.state.currentX
	w.state.prevY = w.state.currentY
	w.state.prevZ = w.state.currentZ
	w.state.currentX = x
	w.state.currentY = y
	w.state.currentZ = z
	w.state.lastLineWasPrint = false

	if err := w.outputRestartPoint(); err != nil {
		return err
	}

	if !savePosition {
		w.state.currentX = w.state.prevX
		w.state.currentY = w.state.prevY
		w.state.currentZ = w.state.prevZ
	}

	return nil
}

func (w *Writer) AddPing() error {
	// flush print line buffer if necessary
	if w.state.printLineBuffered {
		if err := w.outputPrintLine(); err != nil {
			return err
		}
		w.state.printLineBuffered = false
		w.state.lastLineWasPrint = true
	}
	return w.outputPingPoint()
}

func (w *Writer) AddPingAt(x, y, z float32, savePosition bool) error {
	// flush print line buffer if necessary
	if w.state.printLineBuffered {
		if err := w.outputPrintLine(); err != nil {
			return err
		}
		w.state.printLineBuffered = false
		w.state.lastLineWasPrint = true
	}
	// update history
	w.state.prevX = w.state.currentX
	w.state.prevY = w.state.currentY
	w.state.prevZ = w.state.currentZ
	w.state.currentX = x
	w.state.currentY = y
	w.state.currentZ = z
	w.state.lastLineWasPrint = false

	if err := w.outputPingPoint(); err != nil {
		return err
	}

	if !savePosition {
		w.state.currentX = w.state.prevX
		w.state.currentY = w.state.prevY
		w.state.currentZ = w.state.prevZ
	}

	return nil
}

func (w *Writer) getLastIndex() uint32 {
	return (w.bufferSizes["position"] / (floatBytes * 3)) - 1
}

func (w *Writer) getCurrentIndex() uint32 {
	return w.bufferSizes["index"] / uint32Bytes
}

func (w *Writer) outputPrintLine() error {
	fromTool := 0
	t := float32(1)
	if w.state.transitionLineBuffered {
		fromTool = w.state.bufferedFromTool
		t = w.state.bufferedT
		w.state.transitionLineBuffered = false
	}

	//
	// position and normal
	//
	dirX := w.state.currentX - w.state.prevX
	dirY := w.state.currentY - w.state.prevY
	dirZ := w.state.currentZ - w.state.prevZ
	dirSize := float32(math.Sqrt(float64(dirX*dirX) + float64(dirY*dirY) + float64(dirZ*dirZ)))
	dirX /= dirSize
	dirY /= dirSize
	dirZ /= dirSize

	// starting vertex x2
	if err := w.writePosition(w.state.prevX, w.state.prevY, w.state.prevZ); err != nil {
		return err
	}
	if err := w.writePosition(w.state.prevX, w.state.prevY, w.state.prevZ); err != nil {
		return err
	}
	if err := w.writeNormal(dirX, dirY, dirZ); err != nil {
		return err
	}
	if err := w.writeNormal(dirX, dirY, dirZ); err != nil {
		return err
	}

	// if current segment is connected to previous segment, include corner triangles
	if w.state.lastLineWasPrint {
		lastIndex := w.getLastIndex()
		a := lastIndex - 3
		b := lastIndex - 2
		c := lastIndex - 1
		d := lastIndex - 0
		for _, index := range []uint32{a, b, c, c, b, d} {
			if err := w.writeIndex(index); err != nil {
				return err
			}
		}
	}

	// ending vertex x2
	if err := w.writePosition(w.state.currentX, w.state.currentY, w.state.currentZ); err != nil {
		return err
	}
	if err := w.writePosition(w.state.currentX, w.state.currentY, w.state.currentZ); err != nil {
		return err
	}
	if err := w.writeNormal(dirX, dirY, dirZ); err != nil {
		return err
	}
	if err := w.writeNormal(dirX, dirY, dirZ); err != nil {
		return err
	}

	lastIndex := w.getLastIndex()
	a := lastIndex - 3
	b := lastIndex - 2
	c := lastIndex - 1
	d := lastIndex - 0
	for _, index := range []uint32{a, b, c, c, b, d} {
		if err := w.writeIndex(index); err != nil {
			return err
		}
	}

	//
	// extrusion width
	//
	for i := 0; i < 4; i++ {
		if err := w.writeExtrusionWidth(w.state.currentExtrusionWidth); err != nil {
			return err
		}
	}

	//
	// layer height
	//
	for i := 0; i < 4; i++ {
		if err := w.writeLayerHeight(w.state.currentLayerHeight); err != nil {
			return err
		}
	}

	//
	// travel flag
	//
	for i := 0; i < 4; i++ {
		isTravel := w.state.travelling || w.state.currentPathType == PathTypeTravel
		if err := w.writeIsTravel(isTravel); err != nil {
			return err
		}
	}

	//
	// tool colors
	//
	for i := 0; i < 4; i++ {
		if err := w.writeToolColor(w.state.currentTool, fromTool, t); err != nil {
			return err
		}
	}

	//
	// path type colors
	//
	for i := 0; i < 4; i++ {
		if err := w.writePathTypeColor(w.state.currentPathType); err != nil {
			return err
		}
	}

	//
	// feedrate colors
	//
	for i := 0; i < 4; i++ {
		if err := w.writeFeedrateColor(w.state.currentFeedrate); err != nil {
			return err
		}
	}

	//
	// fan colors
	//
	for i := 0; i < 4; i++ {
		if err := w.writeFanSpeedColor(w.state.currentFanSpeed); err != nil {
			return err
		}
	}

	//
	// temperature colors
	//
	for i := 0; i < 4; i++ {
		if err := w.writeTemperatureColor(w.state.currentTemperature); err != nil {
			return err
		}
	}

	//
	// layer height colors
	//
	for i := 0; i < 4; i++ {
		if err := w.writeLayerHeightColor(w.state.currentLayerHeight); err != nil {
			return err
		}
	}

	return nil
}

func (w *Writer) addXYZPrintLineTo(x, y, z float32, savePosition, isTravel bool) error {
	zFloat := roundZ(z)

	if isTravel {
		if w.state.currentPathType != PathTypeTravel {
			// starting travel path after non-travel line
			w.state.travelBufferedPathType = w.state.currentPathType
			w.state.travelBufferedExtrusionWidth = w.state.currentExtrusionWidth
			w.state.travelBufferedLayerHeight = w.state.currentLayerHeight
			w.state.travelBufferedTool = w.state.currentTool
			if err := w.SetPathType(PathTypeTravel); err != nil {
				return err
			}
			if err := w.SetExtrusionWidth(travelExtrusionWidth); err != nil {
				return err
			}
			if err := w.SetLayerHeight(travelLayerHeight); err != nil {
				return err
			}
			if err := w.SetTool(travelTool); err != nil {
				return err
			}
			w.state.travelling = true
		}
	} else {
		if w.state.currentPathType == PathTypeTravel {
			// starting non-travel path after travel line
			w.state.travelling = false
			if err := w.SetPathType(w.state.travelBufferedPathType); err != nil {
				return err
			}
			if err := w.SetExtrusionWidth(w.state.travelBufferedExtrusionWidth); err != nil {
				return err
			}
			if err := w.SetLayerHeight(w.state.travelBufferedLayerHeight); err != nil {
				return err
			}
			if err := w.SetTool(w.state.travelBufferedTool); err != nil {
				return err
			}
		}
	}

	// handle print line buffering/merging
	if w.state.printLineBuffered {
		if directionallyCollinear(
			w.state.prevX, w.state.prevY, w.state.prevZ,
			w.state.currentX, w.state.currentY, w.state.currentZ,
			x, y, zFloat,
		) {
			// spoof history so that segments are merged
			w.state.currentX = w.state.prevX
			w.state.currentY = w.state.prevY
			w.state.currentZ = w.state.prevZ
		} else {
			if err := w.outputPrintLine(); err != nil {
				return err
			}
			w.state.lastLineWasPrint = true
		}
	}

	// update history
	w.state.prevX = w.state.currentX
	w.state.prevY = w.state.currentY
	w.state.prevZ = w.state.currentZ
	w.state.currentX = x
	w.state.currentY = y
	w.state.currentZ = zFloat
	w.state.printLineBuffered = true

	if math.Abs(float64(w.state.currentX-w.state.prevX)) < skipThreshold &&
		math.Abs(float64(w.state.currentY-w.state.prevY)) < skipThreshold &&
		math.Abs(float64(w.state.currentZ-w.state.prevZ)) < skipThreshold {
		// don't output exceedingly-small line segments
		w.state.lastLineWasPrint = false
		w.state.printLineBuffered = false
		return nil
	}

	if w.state.currentTool >= 0 {
		w.state.toolsSeen[w.state.currentTool] = true
	}
	w.state.pathTypesSeen[w.state.currentPathType] = true
	w.state.feedratesSeen[w.state.currentFeedrate] = true
	w.state.fanSpeedsSeen[w.state.currentFanSpeed] = true
	w.state.temperaturesSeen[w.state.currentTemperature] = true
	if w.state.currentPathType != PathTypeTravel && w.state.currentLayerHeight > 0 {
		w.state.layerHeightsSeen[w.state.currentLayerHeight] = true
	}

	if !savePosition {
		if err := w.outputPrintLine(); err != nil {
			return err
		}
		w.state.printLineBuffered = false
		w.state.lastLineWasPrint = false
		w.state.currentX = w.state.prevX
		w.state.currentY = w.state.prevY
		w.state.currentZ = w.state.prevZ
	}

	return nil
}

func (w *Writer) addXYPrintLineTo(x, y float32, savePosition, isTravel bool) error {
	return w.addXYZPrintLineTo(x, y, w.state.currentZ, savePosition, isTravel)
}

func (w *Writer) AddXYZPrintLineTo(x, y, z float32) error {
	return w.addXYZPrintLineTo(x, y, z, true, false)
}

func (w *Writer) AddXYPrintLineTo(x, y float32) error {
	return w.addXYZPrintLineTo(x, y, w.state.currentZ, true, false)
}

func (w *Writer) AddXYZTransitionLineTo(x, y, z float32, fromTool int, t float32) error {
	if err := w.addXYZPrintLineTo(x, y, z, true, false); err != nil {
		return err
	}
	w.state.bufferedFromTool = fromTool
	w.state.bufferedT = t
	w.state.transitionLineBuffered = true
	return nil
}

func (w *Writer) AddXYTransitionLineTo(x, y float32, fromTool int, t float32) error {
	if err := w.addXYPrintLineTo(x, y, true, false); err != nil {
		return err
	}
	w.state.bufferedFromTool = fromTool
	w.state.bufferedT = t
	w.state.transitionLineBuffered = true
	return nil
}

func (w *Writer) AddSideTransitionDangler() error {
	toZ := float32(math.Max(-20, float64(w.state.currentZ)-100))
	return w.addXYZPrintLineTo(w.state.currentX, w.state.currentY, toZ, false, false)
}
