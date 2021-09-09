package ptp

import (
    "bufio"
    "errors"
    "fmt"
    "io/ioutil"
    "math"
    "os"
)

type state struct {
    lastLineWasPrint bool // if true, corner triangles will be created to join the lines
    printLineBuffered bool // if true, a print line from prevX/Y/Z to currentX/Y/Z needs to be output
    transitionLineBuffered bool // if true, use `bufferedFromTool` to create a gradient
    bufferedFromTool int // used for interpolation when outputting a buffered transition line
    bufferedT float32 // used for interpolation when outputting a buffered transition line
    travelLineBuffered bool // if true, a travel line from prevX/Y/Z to currentX/Y/Z needs to be output
    currentX float32
    currentY float32
    currentZ float32
    prevX float32
    prevY float32
    prevZ float32
    currentExtrusionWidth float32
    currentLayerHeight float32
    currentTool int
    currentPathType PathType
    currentFeedrate float32
    currentFanSpeed int
    currentTemperature float32
    zSeen map[float32]bool
    toolsSeen map[int]bool
    pathTypesSeen map[PathType]bool
    feedratesSeen map[float32]bool
    fanSpeedsSeen map[int]bool
    temperaturesSeen map[float32]bool
    layerHeightsSeen map[float32]bool
}

func getStartingState() state {
    return state{
        lastLineWasPrint:       false,
        printLineBuffered:      false,
        transitionLineBuffered: false,
        bufferedFromTool:       0,
        bufferedT:              0,
        travelLineBuffered:     false,
        currentX:               0,
        currentY:               0,
        currentZ:               0,
        prevX:                  0,
        prevY:                  0,
        prevZ:                  0,
        currentExtrusionWidth:  0,
        currentLayerHeight:     0,
        currentTool:            0,
        currentPathType:        PathTypeUnknown,
        currentFeedrate:        0,
        currentFanSpeed:        0,
        currentTemperature:     0,
        zSeen:                  make(map[float32]bool),
        toolsSeen:              make(map[int]bool),
        pathTypesSeen:          make(map[PathType]bool),
        feedratesSeen:          make(map[float32]bool),
        fanSpeedsSeen:          make(map[int]bool),
        temperaturesSeen:       make(map[float32]bool),
        layerHeightsSeen:       make(map[float32]bool),
    }
}

type Writer struct {
    version uint8
    paths map[string]string
    files map[string]*os.File
    writers map[string]*bufio.Writer
    bufferSizes map[string]int

    // bounds for interpolated color scales
    minFeedrate float32
    maxFeedrate float32
    minTemperature float32
    maxTemperature float32
    minLayerHeight float32
    maxLayerHeight float32

    brimIsSkirt bool // whether brim paths should be called skirts
    toolColors [][3]float32 // array of [r, g, b] floats in range 0..1
    state state
}

func NewWriter(outpath string, brimIsSkirt bool, toolColors [][3]float32) Writer {
    return Writer{
        version:        ptpVersion,
        paths:          map[string]string{
            "main":             outpath,
            "legend":           fmt.Sprintf("%s.%s", outpath, "legend"),
            "normal":           fmt.Sprintf("%s.%s", outpath, "normal"),
            "index":            fmt.Sprintf("%s.%s", outpath, "index"),
            "extrusionWidth":   fmt.Sprintf("%s.%s", outpath, "extrusionWidth"),
            "layerHeight":      fmt.Sprintf("%s.%s", outpath, "layerHeight"),
            "travelPosition":   fmt.Sprintf("%s.%s", outpath, "travelPosition"),
            "retractPosition":  fmt.Sprintf("%s.%s", outpath, "retractPosition"),
            "restartPosition":  fmt.Sprintf("%s.%s", outpath, "restartPosition"),
            "toolColor":        fmt.Sprintf("%s.%s", outpath, "toolColor"),
            "pathTypeColor":    fmt.Sprintf("%s.%s", outpath, "pathTypeColor"),
            "feedrateColor":    fmt.Sprintf("%s.%s", outpath, "feedrateColor"),
            "fanSpeedColor":    fmt.Sprintf("%s.%s", outpath, "fanSpeedColor"),
            "temperatureColor": fmt.Sprintf("%s.%s", outpath, "temperatureColor"),
            "layerHeightColor": fmt.Sprintf("%s.%s", outpath, "layerHeightColor"),
        },
        files:          map[string]*os.File{
            "main":             nil,
            "normal":           nil,
            "index":            nil,
            "extrusionWidth":   nil,
            "layerHeight":      nil,
            "travelPosition":   nil,
            "retractPosition":  nil,
            "restartPosition":  nil,
            "toolColor":        nil,
            "pathTypeColor":    nil,
            "feedrateColor":    nil,
            "fanSpeedColor":    nil,
            "temperatureColor": nil,
            "layerHeightColor": nil,
        },
        writers:        map[string]*bufio.Writer{
            "main":             nil,
            "normal":           nil,
            "index":            nil,
            "extrusionWidth":   nil,
            "layerHeight":      nil,
            "travelPosition":   nil,
            "retractPosition":  nil,
            "restartPosition":  nil,
            "toolColor":        nil,
            "pathTypeColor":    nil,
            "feedrateColor":    nil,
            "fanSpeedColor":    nil,
            "temperatureColor": nil,
            "layerHeightColor": nil,
        },
        bufferSizes:    map[string]int{
            "position":         0,
            "normal":           0,
            "index":            0,
            "extrusionWidth":   0,
            "layerHeight":      0,
            "travelPosition":   0,
            "retractPosition":  0,
            "restartPosition":  0,
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
        state:          getStartingState(),
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

    openForWrite(w, "main")
    openForWrite(w, "normal")
    openForWrite(w, "index")
    openForWrite(w, "extrusionWidth")
    openForWrite(w, "layerHeight")
    openForWrite(w, "travelPosition")
    openForWrite(w, "retractPosition")
    openForWrite(w, "restartPosition")
    openForWrite(w, "toolColor")
    openForWrite(w, "pathTypeColor")
    openForWrite(w, "feedrateColor")
    openForWrite(w, "fanSpeedColor")
    openForWrite(w, "temperatureColor")
    openForWrite(w, "layerHeightColor")
    return w.writeHeader()
}

func (w *Writer) writeHeader() error {
    buf := make([]byte, headerSize)
    buf[0] = w.version // only first byte of header is used
    _, err := w.writers["main"].Write(buf)
    return err
}

func (w *Writer) Finalize() error {
    // flush any remaining buffers
    if w.state.printLineBuffered {
        w.outputPrintLine()
        w.state.printLineBuffered = false
    } else if w.state.travelLineBuffered {
        w.outputTravelLine()
        w.state.travelLineBuffered = false
    }

    // close the temp files
    flushAndClose(w, "normal")
    flushAndClose(w, "index")
    flushAndClose(w, "extrusionWidth")
    flushAndClose(w, "layerHeight")
    flushAndClose(w, "travelPosition")
    flushAndClose(w, "retractPosition")
    flushAndClose(w, "restartPosition")
    flushAndClose(w, "toolColor")
    flushAndClose(w, "pathTypeColor")
    flushAndClose(w, "feedrateColor")
    flushAndClose(w, "fanSpeedColor")
    flushAndClose(w, "temperatureColor")
    flushAndClose(w, "layerHeightColor")

    // concatenate the files
    concatOntoWriter(w, "main", "normal")
    concatOntoWriter(w, "main", "index")
    concatOntoWriter(w, "main", "extrusionWidth")
    concatOntoWriter(w, "main", "layerHeight")

    // write legend and commit main file
    w.saveLegend()
    flushAndClose(w, "main")
    return nil
}

func (w *Writer) saveLegend() error {
    legend, err := w.getLegend()
    if err != nil {
        return err
    }
    return ioutil.WriteFile(w.paths["legend"], legend, 0644)
}

func (w *Writer) writePosition(x, y, z float32) {
    writeFloat32LE(w.writers["main"], x)
    writeFloat32LE(w.writers["main"], y)
    writeFloat32LE(w.writers["main"], z)
    w.bufferSizes["position"] += floatBytes * 3
}

func (w *Writer) writeNormal(x, y, z float32) {
    writeFloat32LE(w.writers["normal"], x)
    writeFloat32LE(w.writers["normal"], y)
    writeFloat32LE(w.writers["normal"], z)
    w.bufferSizes["normal"] += floatBytes * 3
}

func (w *Writer) writeIndex(idx uint32) {
    writeUint32LE(w.writers["index"], idx)
    w.bufferSizes["index"] += uint32Bytes
}

func (w *Writer) writeExtrusionWidth(width float32) {
    writeFloat32LE(w.writers["extrusionWidth"], width)
    w.bufferSizes["extrusionWidth"] += floatBytes
}

func (w *Writer) writeLayerHeight(height float32) {
    writeFloat32LE(w.writers["layerHeight"], height)
    w.bufferSizes["layerHeight"] += floatBytes
}

func (w *Writer) writeTravelPosition(x, y, z float32) {
    writeFloat32LE(w.writers["travelPosition"], x)
    writeFloat32LE(w.writers["travelPosition"], y)
    writeFloat32LE(w.writers["travelPosition"], z)
    w.bufferSizes["travelPosition"] += floatBytes * 3
}

func (w *Writer) writeRetractPosition(x, y, z float32) {
    writeFloat32LE(w.writers["retractPosition"], x)
    writeFloat32LE(w.writers["retractPosition"], y)
    writeFloat32LE(w.writers["retractPosition"], z)
    w.bufferSizes["retractPosition"] += floatBytes * 3
}

func (w *Writer) writeRestartPosition(x, y, z float32) {
    writeFloat32LE(w.writers["restartPosition"], x)
    writeFloat32LE(w.writers["restartPosition"], y)
    writeFloat32LE(w.writers["restartPosition"], z)
    w.bufferSizes["restartPosition"] += floatBytes * 3
}

func (w *Writer) writeToolColor(toTool, fromTool int, t float32) {
    var r, g, b float32
    if t >= 1 {
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
    writeFloat32LE(w.writers["toolColor"], r)
    writeFloat32LE(w.writers["toolColor"], g)
    writeFloat32LE(w.writers["toolColor"], b)
    w.bufferSizes["toolColor"] += floatBytes * 3
}

func (w *Writer) writePathTypeColor(pathType PathType) {
    writeFloat32LE(w.writers["pathTypeColor"], pathTypeColors[pathType][0])
    writeFloat32LE(w.writers["pathTypeColor"], pathTypeColors[pathType][1])
    writeFloat32LE(w.writers["pathTypeColor"], pathTypeColors[pathType][2])
    w.bufferSizes["pathTypeColor"] += floatBytes * 3
}

func (w *Writer) writeFeedrateColor(feedrate float32) {
    t := float32(1)
    if w.maxFeedrate > w.minFeedrate {
        t = (feedrate - w.minFeedrate) / (w.maxFeedrate - w.minFeedrate)
    }
    writeFloat32LE(w.writers["feedrateColor"], t)
    w.bufferSizes["feedrateColor"] += floatBytes
}

func (w *Writer) writeFanSpeedColor(pwmValue int) {
    t := float32(pwmValue) / 255
    writeFloat32LE(w.writers["fanSpeedColor"], t)
    w.bufferSizes["fanSpeedColor"] += floatBytes
}

func (w *Writer) writeTemperatureColor(temperature float32) {
    t := float32(1)
    if w.maxTemperature > w.minTemperature {
        t = (temperature - w.minTemperature) / (w.maxTemperature - w.minTemperature)
    }
    writeFloat32LE(w.writers["temperatureColor"], t)
    w.bufferSizes["temperatureColor"] += floatBytes
}

func (w *Writer) writeLayerHeightColor(layerHeight float32) {
    t := float32(1)
    if w.maxLayerHeight > w.minLayerHeight {
        t = (layerHeight - w.minLayerHeight) / (w.maxLayerHeight - w.minLayerHeight)
    }
    writeFloat32LE(w.writers["layerHeightColor"], t)
    w.bufferSizes["layerHeightColor"] += floatBytes
}

func (w *Writer) SetExtrusionWidth(width float32) {
    if width == w.state.currentExtrusionWidth {
        return
    }
    if w.state.printLineBuffered {
        w.outputPrintLine()
        w.state.printLineBuffered = false
        w.state.lastLineWasPrint = true
    }
    w.state.currentExtrusionWidth = width
}

func (w *Writer) SetLayerHeight(height float32) {
    if height == w.state.currentLayerHeight {
        return
    }
    if w.state.printLineBuffered {
        w.outputPrintLine()
        w.state.printLineBuffered = false
        w.state.lastLineWasPrint = true
    }
    w.state.currentLayerHeight = height
}

func (w *Writer) SetTool(tool int) {
    if tool == w.state.currentTool {
        return
    }
    if w.state.printLineBuffered {
        w.outputPrintLine()
        w.state.printLineBuffered = false
        w.state.lastLineWasPrint = true
    }
    w.state.currentTool = tool
}

func (w *Writer) SetPathType(pathType PathType) {
    if pathType == w.state.currentPathType {
        return
    }
    if w.state.printLineBuffered {
        w.outputPrintLine()
        w.state.printLineBuffered = false
        w.state.lastLineWasPrint = true
    }
    w.state.currentPathType = pathType
}

func (w *Writer) SetFeedrate(feedrate float32) {
    if feedrate == w.state.currentFeedrate {
        return
    }
    if w.state.printLineBuffered {
        w.outputPrintLine()
        w.state.printLineBuffered = false
        w.state.lastLineWasPrint = true
    }
    w.state.currentFeedrate = feedrate
}

func (w *Writer) SetFanSpeed(pwmValue int) {
    if pwmValue == w.state.currentFanSpeed {
        return
    }
    if w.state.printLineBuffered {
        w.outputPrintLine()
        w.state.printLineBuffered = false
        w.state.lastLineWasPrint = true
    }
    w.state.currentFanSpeed = pwmValue
}

func (w *Writer) SetTemperature(temperature float32) {
    if temperature == w.state.currentTemperature {
        return
    }
    if w.state.printLineBuffered {
        w.outputPrintLine()
        w.state.printLineBuffered = false
        w.state.lastLineWasPrint = true
    }
    w.state.currentTemperature = temperature
}

func (w *Writer) GetCurrentPosition() (float32, float32, float32) {
    return w.state.currentX, w.state.currentY, w.state.currentZ
}

func (w *Writer) outputTravelLine() {
    w.writeTravelPosition(w.state.prevX, w.state.prevY, w.state.prevZ)
    w.writeTravelPosition(w.state.currentX, w.state.currentY, w.state.currentZ)
}

func (w *Writer) outputRetractPoint() {
    w.writeRetractPosition(w.state.currentX, w.state.currentY, w.state.currentZ)
}

func (w *Writer) outputRestartPoint() {
    w.writeRestartPosition(w.state.currentX, w.state.currentY, w.state.currentZ)
}

func (w *Writer) AddXYZTravelTo(x, y, z float32) {
    // flush print line buffer if necessary
    if w.state.printLineBuffered {
        w.outputPrintLine()
        w.state.printLineBuffered = false
        w.state.lastLineWasPrint = true
    }
    // handle travel line buffering/merging
    if w.state.travelLineBuffered {
        isMergeable := directionallyCollinear(
            w.state.prevX, w.state.prevY, w.state.prevZ,
            w.state.currentX, w.state.currentY, w.state.currentZ,
            x, y, z,
        )
        if isMergeable {
        // spoof history
        w.state.currentX = w.state.prevX
        w.state.currentY = w.state.prevY
        w.state.currentZ = w.state.prevZ
        } else {
            w.outputTravelLine()
        }
    }
    // update history
    w.state.prevX = w.state.currentX
    w.state.prevY = w.state.currentY
    w.state.prevZ = w.state.currentZ
    w.state.currentX = x
    w.state.currentY = y
    w.state.currentZ = z
    w.state.lastLineWasPrint = false
    w.state.travelLineBuffered = true

    if math.Abs(float64(w.state.currentX - w.state.prevX)) < skipThreshold &&
        math.Abs(float64(w.state.currentY - w.state.prevY)) < skipThreshold &&
        math.Abs(float64(w.state.currentZ - w.state.prevZ)) < skipThreshold {
        // don't output exceedingly-small line segments
        w.state.travelLineBuffered = false
        return
    }
}

func (w *Writer) AddXYTravelTo(x, y float32) {
    w.AddXYZTravelTo(x, y, w.state.currentZ)
}

func (w *Writer) AddRetract() {
    // flush print line buffer if necessary
    if w.state.printLineBuffered {
        w.outputPrintLine()
        w.state.printLineBuffered = false
        w.state.lastLineWasPrint = true
    }
    // flush travel line buffer if necessary
    if w.state.travelLineBuffered {
        w.outputTravelLine()
        w.state.travelLineBuffered = false
        w.state.lastLineWasPrint = false
    }
    w.outputRetractPoint()
}

func (w *Writer) AddRestart() {
    // flush print line buffer if necessary
    if w.state.printLineBuffered {
        w.outputPrintLine()
        w.state.printLineBuffered = false
        w.state.lastLineWasPrint = true
    }
    // flush travel line buffer if necessary
    if w.state.travelLineBuffered {
        w.outputTravelLine()
        w.state.travelLineBuffered = false
        w.state.lastLineWasPrint = false
    }
    w.outputRestartPoint()
}

func (w *Writer) AddRetractAt(x, y, z float32, savePosition bool) {
    // flush print line buffer if necessary
    if w.state.printLineBuffered {
        w.outputPrintLine()
        w.state.printLineBuffered = false
        w.state.lastLineWasPrint = true
    }
    // flush travel line buffer if necessary
    if w.state.travelLineBuffered {
        w.outputTravelLine()
        w.state.travelLineBuffered = false
        w.state.lastLineWasPrint = false
    }
    // update history
    w.state.prevX = w.state.currentX
    w.state.prevY = w.state.currentY
    w.state.prevZ = w.state.currentZ
    w.state.currentX = x
    w.state.currentY = y
    w.state.currentZ = z
    w.state.lastLineWasPrint = false

    w.outputRetractPoint()

    if !savePosition {
        w.state.currentX = w.state.prevX
        w.state.currentY = w.state.prevY
        w.state.currentZ = w.state.prevZ
    }
}

func (w *Writer) AddRestartAt(x, y, z float32, savePosition bool) {
    // flush print line buffer if necessary
    if w.state.printLineBuffered {
        w.outputPrintLine()
        w.state.printLineBuffered = false
        w.state.lastLineWasPrint = true
    }
    // flush travel line buffer if necessary
    if w.state.travelLineBuffered {
        w.outputTravelLine()
        w.state.travelLineBuffered = false
        w.state.lastLineWasPrint = false
    }
    // update history
    w.state.prevX = w.state.currentX
    w.state.prevY = w.state.currentY
    w.state.prevZ = w.state.currentZ
    w.state.currentX = x
    w.state.currentY = y
    w.state.currentZ = z
    w.state.lastLineWasPrint = false

    w.outputRestartPoint()

    if !savePosition {
        w.state.currentX = w.state.prevX
        w.state.currentY = w.state.prevY
        w.state.currentZ = w.state.prevZ
    }
}

func (w *Writer) outputPrintLine() {
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
    dirSize := float32(math.Sqrt(float64(dirX * dirX) + float64(dirY * dirY) + float64(dirZ * dirZ)))
    dirX /= dirSize
    dirY /= dirSize
    dirZ /= dirSize

    // starting vertex x2
    w.writePosition(w.state.prevX, w.state.prevY, w.state.prevZ)
    w.writePosition(w.state.prevX, w.state.prevY, w.state.prevZ)
    w.writeNormal(dirX, dirY, dirZ)
    w.writeNormal(dirX, dirY, dirZ)

    // if current segment is connected to previous segment, include corner triangles
    if w.state.lastLineWasPrint {
        lastIndex := uint32((w.bufferSizes["position"] / (floatBytes * 3)) - 1)
        a := lastIndex - 3
        b := lastIndex - 2
        c := lastIndex - 1
        d := lastIndex - 0
        w.writeIndex(a); w.writeIndex(b); w.writeIndex(c)
        w.writeIndex(c); w.writeIndex(b); w.writeIndex(d)
    }

    // ending vertex x2
    w.writePosition(w.state.currentX, w.state.currentY, w.state.currentZ)
    w.writePosition(w.state.currentX, w.state.currentY, w.state.currentZ)
    w.writeNormal(dirX, dirY, dirZ)
    w.writeNormal(dirX, dirY, dirZ)

    lastIndex := uint32((w.bufferSizes["position"] / (floatBytes * 3)) - 1)
    a := lastIndex - 3
    b := lastIndex - 2
    c := lastIndex - 1
    d := lastIndex - 0
    w.writeIndex(a); w.writeIndex(b); w.writeIndex(c)
    w.writeIndex(c); w.writeIndex(b); w.writeIndex(d)

    //
    // extrusion width
    //
    w.writeExtrusionWidth(w.state.currentExtrusionWidth)
    w.writeExtrusionWidth(w.state.currentExtrusionWidth)
    w.writeExtrusionWidth(w.state.currentExtrusionWidth)
    w.writeExtrusionWidth(w.state.currentExtrusionWidth)

    //
    // layer height
    //
    w.writeLayerHeight(w.state.currentLayerHeight)
    w.writeLayerHeight(w.state.currentLayerHeight)
    w.writeLayerHeight(w.state.currentLayerHeight)
    w.writeLayerHeight(w.state.currentLayerHeight)

    //
    // tool colors
    //
    w.writeToolColor(w.state.currentTool, fromTool, t)
    w.writeToolColor(w.state.currentTool, fromTool, t)
    w.writeToolColor(w.state.currentTool, fromTool, t)
    w.writeToolColor(w.state.currentTool, fromTool, t)

    //
    // path type colors
    //
    w.writePathTypeColor(w.state.currentPathType)
    w.writePathTypeColor(w.state.currentPathType)
    w.writePathTypeColor(w.state.currentPathType)
    w.writePathTypeColor(w.state.currentPathType)

    //
    // feedrate colors
    //
    w.writeFeedrateColor(w.state.currentFeedrate)
    w.writeFeedrateColor(w.state.currentFeedrate)
    w.writeFeedrateColor(w.state.currentFeedrate)
    w.writeFeedrateColor(w.state.currentFeedrate)

    //
    // fan colors
    //
    w.writeFanSpeedColor(w.state.currentFanSpeed)
    w.writeFanSpeedColor(w.state.currentFanSpeed)
    w.writeFanSpeedColor(w.state.currentFanSpeed)
    w.writeFanSpeedColor(w.state.currentFanSpeed)

    //
    // temperature colors
    //
    w.writeTemperatureColor(w.state.currentTemperature)
    w.writeTemperatureColor(w.state.currentTemperature)
    w.writeTemperatureColor(w.state.currentTemperature)
    w.writeTemperatureColor(w.state.currentTemperature)

    //
    // layer height colors
    //
    w.writeLayerHeightColor(w.state.currentLayerHeight)
    w.writeLayerHeightColor(w.state.currentLayerHeight)
    w.writeLayerHeightColor(w.state.currentLayerHeight)
    w.writeLayerHeightColor(w.state.currentLayerHeight)
}

func (w *Writer) addXYZPrintLineTo(x, y, z float32, fromTool int, t float32, savePosition bool) {
    zFloat := roundZ(z)
    // flush travel line buffer if necessary
    if w.state.travelLineBuffered {
        w.outputTravelLine()
        w.state.travelLineBuffered = false
    }
    // handle print line buffering/merging
    if w.state.printLineBuffered {
        isMergeable := directionallyCollinear(
            w.state.prevX, w.state.prevY, w.state.prevZ,
            w.state.currentX, w.state.currentY, w.state.currentZ,
            x, y, zFloat,
        )
        if isMergeable {
            // spoof history so that segments are merged
            w.state.currentX = w.state.prevX
            w.state.currentY = w.state.prevY
            w.state.currentZ = w.state.prevZ
        } else {
            w.outputPrintLine()
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

    if math.Abs(float64(w.state.currentX - w.state.prevX)) < skipThreshold &&
        math.Abs(float64(w.state.currentY - w.state.prevY)) < skipThreshold &&
        math.Abs(float64(w.state.currentZ - w.state.prevZ)) < skipThreshold {
        // don't output exceedingly-small line segments
        w.state.lastLineWasPrint = false
        w.state.printLineBuffered = false
        return
    }

    w.state.zSeen[zFloat] = true
    w.state.toolsSeen[w.state.currentTool] = true
    w.state.pathTypesSeen[w.state.currentPathType] = true
    w.state.feedratesSeen[w.state.currentFeedrate] = true
    w.state.fanSpeedsSeen[w.state.currentFanSpeed] = true
    w.state.temperaturesSeen[w.state.currentTemperature] = true
    if w.state.currentLayerHeight > 0 {
        w.state.layerHeightsSeen[w.state.currentLayerHeight] = true
    }

    if !savePosition {
        w.outputPrintLine()
        w.state.printLineBuffered = false
        w.state.lastLineWasPrint = false
        w.state.currentX = w.state.prevX
        w.state.currentY = w.state.prevY
        w.state.currentZ = w.state.prevZ
    }
}

func (w *Writer) addXYPrintLineTo(x, y float32, fromTool int, t float32, savePosition bool) {
    w.addXYZPrintLineTo(x, y, w.state.currentZ, fromTool, t, savePosition)
}

func (w *Writer) AddXYZPrintLineTo(x, y, z float32) {
    w.addXYZPrintLineTo(x, y, z, 0, 1, true)
}

func (w *Writer) AddXYPrintLineTo(x, y float32) {
    w.addXYZPrintLineTo(x, y, w.state.currentZ, 0, 1, true)
}

func (w *Writer) AddXYZTransitionLineTo(x, y, z float32, fromTool int, t float32) {
    w.addXYZPrintLineTo(x, y, z, fromTool, t, true)
    w.state.bufferedFromTool = fromTool
    w.state.bufferedT = t
    w.state.transitionLineBuffered = true
}

func (w *Writer) AddXYTransitionLineTo(x, y float32, fromTool int, t float32) {
    w.addXYPrintLineTo(x, y, fromTool, t, true)
    w.state.bufferedFromTool = fromTool
    w.state.bufferedT = t
    w.state.transitionLineBuffered = true
}

func (w *Writer) AddSideTransitionDangler() {
    toZ := float32(math.Max(-20, float64(w.state.currentZ) - 100))
    w.addXYZPrintLineTo(w.state.currentX, w.state.currentY, toZ, 0, 1, false)
}
