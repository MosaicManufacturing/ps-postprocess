package ultimaker

import (
    "../gcode"
    "fmt"
    "math"
    "time"
)

const EOL = "\r\n"

func getCurrentDateString() string {
    now := time.Now()
    return now.Format("2006-02-01")
}

type griffinOpts struct {
    firstTemperature float64
    firstBedTemperature float64
    materialVolumeUsed float64
    nozzleDiameter float64
    totalPrintTime int
    boundingBox gcode.BoundingBox
}

func getUltimakerGriffinHeader(opts griffinOpts) string {
    date := getCurrentDateString()
    initialExtTemp := int(math.Round(opts.firstTemperature))
    initialBedTemp := int(math.Round(opts.firstBedTemperature))
    materialVolume := float32(math.Round(opts.materialVolumeUsed))
    bbox := opts.boundingBox
    if bbox.Min[0] > bbox.Max[0] || bbox.Min[1] > bbox.Max[1] || bbox.Min[2] > bbox.Max[2] {
        // safety default
        bbox.Min = [3]float32{0, 0, 0}
        bbox.Max = [3]float32{10, 10, 10}
    }
    printTime := opts.totalPrintTime
    if printTime < 10 {
        printTime = 10
    }

    header := ";START_OF_HEADER" + EOL
    header += ";HEADER_VERSION:0.1" + EOL
    header += ";FLAVOR:Griffin" + EOL
    header += ";GENERATOR.NAME:Canvas" + EOL
    header += ";GENERATOR.VERSION:1.0.0" + EOL
    header += fmt.Sprintf(";GENERATOR.BUILD_DATE:%s%s", date, EOL)
    header += ";TARGET_MACHINE.NAME:Ultimaker S" + EOL
    header += fmt.Sprintf(";EXTRUDER_TRAIN.0.INITIAL_TEMPERATURE:%d%s", initialExtTemp, EOL)
    header += fmt.Sprintf(";EXTRUDER_TRAIN.0.MATERIAL.VOLUME_USED:%f%s", materialVolume, EOL)
    header += ";EXTRUDER_TRAIN.0.MATERIAL.GUID:506c9f0d-e3aa-4bd4-b2d2-23e2425b1aa9" + EOL
    header += fmt.Sprintf(";EXTRUDER_TRAIN.0.NOZZLE.DIAMETER:%.2f%s", opts.nozzleDiameter, EOL)
    header += fmt.Sprintf(";EXTRUDER_TRAIN.0.NOZZLE.NAME:AA %.2f%s", opts.nozzleDiameter, EOL)
    header += ";BUILD_PLATE.TYPE:glass" + EOL
    header += fmt.Sprintf(";BUILD_PLATE.INITIAL_TEMPERATURE:%d%s", initialBedTemp, EOL)
    header += fmt.Sprintf(";PRINT.TIME:%d%s", printTime, EOL)
    header += fmt.Sprintf(";PRINT.SIZE.MIN.X:%.3f%s", bbox.Min[0], EOL)
    header += fmt.Sprintf(";PRINT.SIZE.MIN.Y:%.3f%s", bbox.Min[1], EOL)
    header += fmt.Sprintf(";PRINT.SIZE.MIN.Z:%.3f%s", bbox.Min[2], EOL)
    header += fmt.Sprintf(";PRINT.SIZE.MAX.X:%.3f%s", bbox.Max[0], EOL)
    header += fmt.Sprintf(";PRINT.SIZE.MAX.Y:%.3f%s", bbox.Max[1], EOL)
    header += fmt.Sprintf(";PRINT.SIZE.MAX.Z:%.3f%s", bbox.Max[2], EOL)
    header += ";END_OF_HEADER" + EOL
    return header
}