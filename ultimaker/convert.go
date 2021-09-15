package ultimaker

import (
    "../gcode"
    "io"
    "log"
    "os"
    "strconv"
)

func AddHeader(argv []string) {
    argc := len(argv)

    if argc != 8 {
        log.Fatalln("expected 8 command-line arguments")
    }
    inpath := argv[0]
    outpath := argv[1]
    firstTemperature, err := strconv.ParseFloat(argv[2], 32)
    if err != nil {
        log.Fatalln(err)
    }
    firstBedTemperature, err := strconv.ParseFloat(argv[3], 32)
    if err != nil {
        log.Fatalln(err)
    }
    materialVolumeUsed, err := strconv.ParseFloat(argv[4], 32)
    if err != nil {
        log.Fatalln(err)
    }
    nozzleDiameter, err := strconv.ParseFloat(argv[5], 32)
    if err != nil {
        log.Fatalln(err)
    }
    totalPrintTime, err := strconv.ParseInt(argv[6], 10, 32)
    if err != nil {
        log.Fatalln(err)
    }
    boundingBox, err := gcode.UnserializeBoundingBox(argv[7])
    if err != nil {
        log.Fatalln(err)
    }

    opts := griffinOpts{
        firstTemperature:    firstTemperature,
        firstBedTemperature: firstBedTemperature,
        materialVolumeUsed:  materialVolumeUsed,
        nozzleDiameter:      nozzleDiameter,
        totalPrintTime:      int(totalPrintTime),
        boundingBox:         boundingBox,
    }

    header := getUltimakerGriffinHeader(opts)

    outfile, err := os.Create(outpath)
    if err != nil {
        log.Fatalln(err)
    }

    // write header first
    if _, err := outfile.WriteString(header); err != nil {
        log.Fatalln(err)
    }

    // concat entire infile
    infile, err := os.Open(inpath)
    if err != nil {
        log.Fatalln(err)
    }
    if _, err := io.Copy(outfile, infile); err != nil {
        log.Fatalln(err)
    }
    if err := infile.Close(); err != nil {
        log.Fatalln(err)
    }

    // finalize and close temporary file
    if err := outfile.Sync(); err != nil {
        log.Fatalln(err)
    }
    if err := outfile.Close(); err != nil {
        log.Fatalln(err)
    }
}
