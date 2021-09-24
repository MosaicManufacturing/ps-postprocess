package msf

import (
    "../sequences"
    "fmt"
    "log"
    "os"
)

// explaining outpath and msfpath:
// - P1:            outpath == *.msf.gcode,  msfpath == *.msf
// - P2 accessory:  outpath == *.maf.gcode,  msfpath == *.maf
// - P2 connected:  outpath == *.mcf.gcode,  [no msfpath]
// - P3 accessory:  outpath == *.gcode,      msfpath == *.json
// - P3 connected:  outpath == *.gcode,      msfpath == *.json

func ConvertForPalette(argv []string) {
    argc := len(argv)

    if argc < 6 {
        log.Fatalln("expected 6 command-line arguments")
    }
    inpath := argv[0] // unmodified G-code file
    outpath := argv[1] // modified G-code file
    msfpath := argv[2] // supplementary MSF file, if applicable
    palettepath := argv[3] // serialized Palette data
    localsPath := argv[4] // JSON-stringified locals
    perExtruderLocalsPath := argv[5] // JSON-stringified locals

    palette, err := LoadFromFile(palettepath)
    if err != nil {
        log.Fatalln(err)
    }

    locals := sequences.NewLocals()
    if err := locals.LoadGlobal(localsPath); err != nil {
        log.Fatalln(err)
    }
    if err := locals.LoadPerExtruder(perExtruderLocalsPath); err != nil {
        log.Fatalln(err)
    }

    // preflight: run through the G-code once to determine all necessary
    // information for performing modifications

    // - drives used
    // - splice lengths -- check early if any splices will be too short
    // - number of pings
    // - bounding box
    preflightResults, err := preflight(inpath, &palette)
    if err != nil {
        log.Fatalln(err)
    }
    if preflightResults.totalDrivesUsed() <= 1 {
        fmt.Println("NO_PALETTE")
        os.Exit(0)
    }

    // output: run through the G-code once and apply modifications
    // using information determined in preflight

    // - start of print O commands
    // - add initial toolchange to Palette extruder
    // - remove toolchange commands
    // - accessory pings (two pauses with precise-ish amount of E between them)
    // - connected pings
    // - print summary in footer
    err = paletteOutput(inpath, outpath, msfpath, &palette, &preflightResults, locals)
    if err != nil {
        log.Fatalln(err)
    }
}
