package msf

import (
    "bufio"
    "io/ioutil"
    "mosaicmfg.com/ps-postprocess/gcode"
    "mosaicmfg.com/ps-postprocess/sequences"
    "path"
    "testing"
)

func getTestPalette(transitionLength float32) Palette {
    basicMaterial := Material{
        ID:         "1",
        Index:      1,
        FilamentID: 1,
        Name:       "Material",
        Color:      "000000",
    }
    basicSpliceSettings := SpliceSettings{
        IngoingID:         "1",
        OutgoingID:        "1",
        HeatFactor:        0,
        CompressionFactor: 0,
        CoolingFactor:     0,
        Reverse:           false,
    }
    materialMeta := make([]Material, 8)
    for i := 0; i < 8; i++ {
        materialMeta[i] = basicMaterial
    }
    spliceSettings := make([]SpliceSettings, 8 * 8)
    for i := 0; i < 8 * 8; i++ {
        spliceSettings[i] = basicSpliceSettings
    }
    transitionLengths := make([][]float32, 0, 8)
    for i := 0; i < 8; i++ {
        perIngoing := make([]float32, 8)
        for j := 0; j < 8; j++ {
            perIngoing[j] = transitionLength
        }
        transitionLengths = append(transitionLengths, perIngoing)
    }
    towerSpeeds := make([]float32, 8)
    for i := 0; i < 8; i++ {
        towerSpeeds[i] = 60
    }
    retractDistance := make([]float32, 8)
    retractFeedrate := make([]float32, 8)
    restartDistance := make([]float32, 8)
    restartFeedrate := make([]float32, 8)
    wipe := make([]bool, 8)
    zLift := make([]float32, 8)
    for i := 0; i < 8; i++ {
        retractDistance[i] = 1
        restartDistance[i] = 1
        retractFeedrate[i] = 30
        restartFeedrate[i] = 30
        wipe[i] = false
        zLift[i] = 0.5
    }

    palette := Palette{
        Type:                       TypeP3,
        Model:                      ModelP3Pro,
        MaterialMeta:               materialMeta,
        SpliceSettings:             spliceSettings,
        PrintExtruder:              0,
        FirmwarePurge:              0,
        BowdenTubeLength:           BowdenDefault,
        TravelSpeedXY:              0,
        TravelSpeedZ:               0,
        PrintBedMinX:               0,
        PrintBedMaxX:               100,
        PrintBedMinY:               0,
        PrintBedMaxY:               100,
        TransitionMethod:           CustomTower,
        TransitionLengths:          transitionLengths,
        TransitionTarget:           40,
        InfillTransitioning:        false,
        TowerSize:                  [2]float32{40, 60},
        TowerPosition:              [2]float32{50, 50},
        TowerMinDensity:            5,
        TowerMinFirstLayerDensity:  10,
        TowerMaxDensity:            100,
        TowerMinBrims:              0,
        TowerSpeed:                 towerSpeeds,
        TowerExtrusionWidth:        0.4,
        TowerExtrusionMultiplier:   100,
        TowerFirstLayerPerimeters:  false,
        InfillPerimeterOverlap:     20,
        RaftLayers:                 0,
        RaftInflation:              0,
        RaftExtrusionWidth:         0,
        RaftStride:                 0,
        UseFirmwareRetraction:      false,
        RetractDistance:            retractDistance,
        RestartDistance:            retractFeedrate,
        RetractFeedrate:            restartDistance,
        RestartFeedrate:            restartFeedrate,
        Wipe:                       wipe,
        ZLift:                      zLift,
        ZOffset:                    0,
        PreSideTransitionSequence:  "",
        SideTransitionSequence:     "",
        PostSideTransitionSequence: "",
        PreSideTransitionScript:    nil,
        SideTransitionScript:       nil,
        PostSideTransitionScript:   nil,
        SideTransitionJog:          false,
        SideTransitionPurgeSpeed:   0,
        SideTransitionMoveSpeed:    0,
        SideTransitionX:            0,
        SideTransitionY:            0,
        SideTransitionEdge:         0,
        SideTransitionEdgeOffset:   0,
        PingOffTowerDistance:       0,
        JogPauses:                  false,
        ClearBufferCommand:         "G4",
        ConnectedMode:              false,
        PrinterID:                  "1",
        Filename:                   "test",
        LoadingOffset:              0,
        PrintValue:                 0,
        CalibrationLength:          0,
    }
    return palette
}

func assertTransitionValuesMatch(t *testing.T, idx int, expected, actual Transition) {
    //fmt.Printf("transition %d\n", idx + 1)
    //fmt.Printf("  Layer: %d\n", actual.Layer)
    //fmt.Printf("  From: %d\n", actual.From)
    //fmt.Printf("  To: %d\n", actual.To)
    //fmt.Printf("  TotalExtrusion: %f\n", actual.TotalExtrusion)
    //fmt.Printf("  TransitionLength: %f\n", actual.TransitionLength)
    //fmt.Printf("  PurgeLength: %f\n", actual.PurgeLength)
    //fmt.Printf("  UsableInfill: %f\n", actual.UsableInfill)
    if actual.Layer != expected.Layer {
        t.Errorf("transition %d: expected Layer == %d, got %d", idx, expected.Layer, actual.Layer)
    }
    if actual.From != expected.From {
        t.Errorf("transition %d: expected From == %d, got %d", idx, expected.From, actual.From)
    }
    if actual.To != expected.To {
        t.Errorf("transition %d: expected To == %d, got %d", idx, expected.To, actual.To)
    }
    if actual.TotalExtrusion != expected.TotalExtrusion {
        t.Errorf("transition %d: expected TotalExtrusion == %f, got %f", idx, expected.TotalExtrusion, actual.TotalExtrusion)
    }
    if actual.TransitionLength != expected.TransitionLength {
        t.Errorf("transition %d: expected TransitionLength == %f, got %f", idx, expected.TransitionLength, actual.TransitionLength)
    }
    if actual.PurgeLength != expected.PurgeLength {
        t.Errorf("transition %d: expected PurgeLength == %f, got %f", idx, expected.PurgeLength, actual.PurgeLength)
    }
    if actual.UsableInfill != expected.UsableInfill {
        t.Errorf("transition %d: expected UsableInfill == %f, got %f", idx, expected.UsableInfill, actual.UsableInfill)
    }
}

func testTowerPreflight(t *testing.T, palette *Palette, printContent string, expectedTransitions []Transition) msfPreflight {
    gcodeLines := gcode.ParseLines(printContent)
    readerFn := func(callback gcode.LineCallback) error {
        for lineNumber, line := range gcodeLines {
            if err := callback(line, lineNumber); err != nil {
                return err
            }
        }
        return nil
    }
    results, err := _preflight(readerFn, palette)
    if err != nil {
        t.Fatal(err)
    }
    if len(results.transitions) != len(expectedTransitions) {
        t.Fatalf("expected %d transitions, got %d", len(expectedTransitions), len(results.transitions))
    }

    for idx, transition := range results.transitions {
        assertTransitionValuesMatch(t, idx, transition, expectedTransitions[idx])
    }

    return results
}

func testTowerOutput(t *testing.T, palette *Palette, printContent string, preflight *msfPreflight, locals sequences.Locals) {
    writer := bufio.NewWriter(ioutil.Discard)

    gcodeLines := gcode.ParseLines(printContent)
    readerFn := func(callback gcode.LineCallback) error {
        for lineNumber, line := range gcodeLines {
            if err := callback(line, lineNumber); err != nil {
                return err
            }
        }
        return nil
    }
    msfOut := NewMSF(palette)
    err := _paletteOutput(readerFn, writer, msfOut, palette, preflight, locals)
    if err != nil {
        t.Fatal(err)
    }
}

func testWithInputFiles(t *testing.T, folderName string) {
    printPath := path.Join("test-files", folderName, "print.gcode")
    palettePath := path.Join("test-files", folderName, "palette.json")
    localsPath := path.Join("test-files", folderName, "locals.json")
    perExtruderLocalsPath := path.Join("test-files", folderName, "perExtruderLocals.json")

    palette, err := LoadPaletteFromFile(palettePath)
    if err != nil {
        t.Fatal(err)
    }
    preflightResults, err := preflight(printPath, &palette)
    if err != nil {
        t.Fatal(err)
    }
    printContent, err := ioutil.ReadFile(printPath)
    locals := sequences.NewLocals()
    if err = locals.LoadGlobal(localsPath); err != nil {
        t.Fatal(err)
    }
    if err = locals.LoadPerExtruder(perExtruderLocalsPath); err != nil {
        t.Fatal(err)
    }
    testTowerOutput(t, &palette, string(printContent), &preflightResults, locals)
}