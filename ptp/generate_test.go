package ptp

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func testFanSpeed(t *testing.T, name string, gcodeContent string, expectedFanSpeed []legendEntry) {
	t.Run(name, func(t *testing.T) {
		// create a temporary directory for test files
		tempDir, err := os.MkdirTemp("", "ptp_test")
		if err != nil {
			t.Fatal(err)
		}
		defer os.RemoveAll(tempDir)

		gcodePath := filepath.Join(tempDir, "test.gcode")
		err = os.WriteFile(gcodePath, []byte(gcodeContent), 0644)
		if err != nil {
			t.Fatal(err)
		}

		// set up output paths
		outPath := filepath.Join(tempDir, "output.ptp")
		legendPath := outPath + ".legend"

		// arguments for GenerateToolpath
		argv := []string{
			gcodePath,     // input file
			outPath,       // output file
			"0.4",         // initialExtrusionWidth
			"0.2",         // initialLayerHeight
			"0.0",         // zOffset
			"false",       // brimIsSkirt
			"1.0,0.0,0.0", // toolColors (red for tool 0)
		}

		// call GenerateToolpath
		GenerateToolpath(argv)

		// read the legend file
		legendData, err := os.ReadFile(legendPath)
		if err != nil {
			t.Fatal(err)
		}

		// parse the legend JSON
		var legend ptpLegend
		err = json.Unmarshal(legendData, &legend)
		if err != nil {
			t.Fatal(err)
		}

		if len(legend.FanSpeed) != len(expectedFanSpeed) {
			t.Errorf("Expected %d legend entries, got %d", len(expectedFanSpeed), len(legend.FanSpeed))
			return
		}

		for i, expected := range expectedFanSpeed {
			actual := legend.FanSpeed[i]
			if actual.Label != expected.Label || actual.Color != expected.Color {
				t.Errorf("Entry %d: expected [%s, %s], got [%s, %s]", i, expected.Label, expected.Color, actual.Label, actual.Color)
			}
		}
	})
}

func TestFanSpeedEdgeCases(t *testing.T) {
	// expected colors
	offColor := floatsToHex(fanColorMin[0], fanColorMin[1], fanColorMin[2])
	onColor := floatsToHex(fanColorMax[0], fanColorMax[1], fanColorMax[2])

	// test case 1: No fan commands
	testFanSpeed(t, "NoFanCommands", `; Test G-code with no fan commands
G90 ; use absolute coordinates
M83 ; extruder relative mode
M104 S200 ; set extruder temp
M109 S200 ; wait for extruder temp
G1 Z0.2 F720
G1 X10 Y10 Z0.2 F1800
;TYPE:Internal infill
G1 X20 Y20 Z0.2 F1800 E1.0 ; print move
G1 X30 Y30 Z0.2 F1800 E2.0 ; print move
;LAYER_CHANGE
;Z:0.4
;HEIGHT:0.2
G1 Z0.4 F720
G1 X10 Y10 Z0.4 F1800
;TYPE:Internal infill
G1 X20 Y20 Z0.4 F1800 E3.0 ; print move
G1 X30 Y30 Z0.4 F1800 E4.0 ; print move
`, []legendEntry{{Label: "Off", Color: offColor}})

	// test case 2: Only fan off (M107)
	testFanSpeed(t, "OnlyFanOff", `; Test G-code with only fan off
G90
M83
M104 S200
M109 S200
G1 Z0.2 F720
G1 X10 Y10 Z0.2 F1800
;TYPE:Internal infill
G1 X20 Y20 Z0.2 F1800 E1.0
M107 ; fan off
G1 X30 Y30 Z0.2 F1800 E2.0
;HEIGHT:0.2
`, []legendEntry{{Label: "Off", Color: offColor}})

	// test case 3: Only fan on (M106 S255)
	testFanSpeed(t, "OnlyFanOn", `; Test G-code with only fan on
G90
M83
M104 S200
M109 S200
G1 Z0.2 F720
G1 X10 Y10 Z0.2 F1800
;TYPE:Internal infill
G1 X20 Y20 Z0.2 F1800 E1.0
M106 S255 ; fan on
G1 X30 Y30 Z0.2 F1800 E2.0
;HEIGHT:0.2
`, []legendEntry{{Label: "Off", Color: offColor}, {Label: "On", Color: onColor}})

	// test case 4: Mix of off and on
	testFanSpeed(t, "FanOffAndOn", `; Test G-code with fan off and on
G90
M83
M104 S200
M109 S200
G1 Z0.2 F720
G1 X10 Y10 Z0.2 F1800
;TYPE:Internal infill
G1 X20 Y20 Z0.2 F1800 E1.0
M106 S0 ; fan off
G1 X25 Y25 Z0.2 F1800 E1.5
M106 S255 ; fan on
G1 X30 Y30 Z0.2 F1800 E2.0
;HEIGHT:0.2
`, []legendEntry{{Label: "Off", Color: offColor}, {Label: "On", Color: onColor}})

	// test case 5: Single percentage (M106 S128)
	testFanSpeed(t, "SinglePercentage", `; Test G-code with single fan percentage
G90
M83
M104 S200
M109 S200
G1 Z0.2 F720
G1 X10 Y10 Z0.2 F1800
;TYPE:Internal infill
G1 X20 Y20 Z0.2 F1800 E1.0
M106 S128 ; 50% fan
G1 X30 Y30 Z0.2 F1800 E2.0
;HEIGHT:0.2
`, []legendEntry{{
		Label: "50%",
		Color: onColor,
	}})

	// test case 6: Many fan speeds (should trigger gradated legend with legendSteps)
	testFanSpeed(t, "ManyFanSpeeds", `; Test G-code with many different fan speeds
G90
M83
M104 S200
M109 S200
G1 Z0.2 F720
G1 X10 Y10 Z0.2 F1800
;TYPE:Internal infill
G1 X10 Y10 Z0.2 F1800 E1.0
M106 S25 ; 10%
G1 X15 Y15 Z0.2 F1800 E1.5
M106 S51 ; 20%
G1 X20 Y20 Z0.2 F1800 E2.0
M106 S76 ; 30%
G1 X25 Y25 Z0.2 F1800 E2.5
M106 S102 ; 40%
G1 X30 Y30 Z0.2 F1800 E3.0
M106 S127 ; 50%
G1 X35 Y35 Z0.2 F1800 E3.5
M106 S153 ; 60%
G1 X40 Y40 Z0.2 F1800 E4.0
M106 S178 ; 70%
G1 X45 Y45 Z0.2 F1800 E4.5
M106 S204 ; 80%
G1 X50 Y50 Z0.2 F1800 E5.0
M106 S229 ; 90%
G1 X55 Y55 Z0.2 F1800 E5.5
;HEIGHT:0.2
`, []legendEntry{
		{Label: "90%", Color: onColor},
		{Label: "77%", Color: floatsToHex(lerp(fanColorMax[0], fanColorMin[0], 1.0/6.0), lerp(fanColorMax[1], fanColorMin[1], 1.0/6.0), lerp(fanColorMax[2], fanColorMin[2], 1.0/6.0))},
		{Label: "63%", Color: floatsToHex(lerp(fanColorMax[0], fanColorMin[0], 2.0/6.0), lerp(fanColorMax[1], fanColorMin[1], 2.0/6.0), lerp(fanColorMax[2], fanColorMin[2], 2.0/6.0))},
		{Label: "50%", Color: floatsToHex(lerp(fanColorMax[0], fanColorMin[0], 3.0/6.0), lerp(fanColorMax[1], fanColorMin[1], 3.0/6.0), lerp(fanColorMax[2], fanColorMin[2], 3.0/6.0))},
		{Label: "37%", Color: floatsToHex(lerp(fanColorMax[0], fanColorMin[0], 4.0/6.0), lerp(fanColorMax[1], fanColorMin[1], 4.0/6.0), lerp(fanColorMax[2], fanColorMin[2], 4.0/6.0))},
		{Label: "23%", Color: floatsToHex(lerp(fanColorMax[0], fanColorMin[0], 5.0/6.0), lerp(fanColorMax[1], fanColorMin[1], 5.0/6.0), lerp(fanColorMax[2], fanColorMin[2], 5.0/6.0))},
		{Label: "10%", Color: offColor},
	})
}

func testTemperature(t *testing.T, name string, gcodeContent string, expectedTemperature []legendEntry) {
	t.Run(name, func(t *testing.T) {
		// create a temporary directory for test files
		tempDir, err := os.MkdirTemp("", "ptp_test")
		if err != nil {
			t.Fatal(err)
		}
		defer os.RemoveAll(tempDir)

		gcodePath := filepath.Join(tempDir, "test.gcode")
		err = os.WriteFile(gcodePath, []byte(gcodeContent), 0644)
		if err != nil {
			t.Fatal(err)
		}

		// set up output paths
		outPath := filepath.Join(tempDir, "output.ptp")
		legendPath := outPath + ".legend"

		// arguments for GenerateToolpath
		argv := []string{
			gcodePath,     // input file
			outPath,       // output file
			"0.4",         // initialExtrusionWidth
			"0.2",         // initialLayerHeight
			"0.0",         // zOffset
			"false",       // brimIsSkirt
			"1.0,0.0,0.0", // toolColors (red for tool 0)
		}

		// call GenerateToolpath
		GenerateToolpath(argv)

		// read the legend file
		legendData, err := os.ReadFile(legendPath)
		if err != nil {
			t.Fatal(err)
		}

		// parse the legend JSON
		var legend ptpLegend
		err = json.Unmarshal(legendData, &legend)
		if err != nil {
			t.Fatal(err)
		}

		if len(legend.Temperature) != len(expectedTemperature) {
			t.Errorf("Expected %d legend entries, got %d", len(expectedTemperature), len(legend.Temperature))
			return
		}

		for i, expected := range expectedTemperature {
			actual := legend.Temperature[i]
			if actual.Label != expected.Label || actual.Color != expected.Color {
				t.Errorf("Entry %d: expected [%s, %s], got [%s, %s]", i, expected.Label, expected.Color, actual.Label, actual.Color)
			}
		}
	})
}

func TestTemperatureEdgeCases(t *testing.T) {
	// expected colors
	minColor := floatsToHex(temperatureColorMin[0], temperatureColorMin[1], temperatureColorMin[2])
	maxColor := floatsToHex(temperatureColorMax[0], temperatureColorMax[1], temperatureColorMax[2])

	// test case 1: Single temperature
	testTemperature(t, "SingleTemperature", `; Test G-code with single temperature
G90
M83
M104 S200
M109 S200
G1 Z0.2 F720
G1 X10 Y10 Z0.2 F1800
;TYPE:Internal infill
G1 X20 Y20 Z0.2 F1800 E1.0
G1 X30 Y30 Z0.2 F1800 E2.0
;HEIGHT:0.2
`, []legendEntry{{
		Label: "200 °C",
		Color: maxColor,
	}})

	// test case 2: Two different temperatures
	testTemperature(t, "TwoTemperatures", `; Test G-code with two temperatures
G90
M83
M104 S200
M109 S200
G1 Z0.2 F720
G1 X10 Y10 Z0.2 F1800
;TYPE:Internal infill
G1 X20 Y20 Z0.2 F1800 E1.0
M104 S220
M109 S220
G1 X30 Y30 Z0.2 F1800 E2.0
;HEIGHT:0.2
`, []legendEntry{
		{Label: "220 °C", Color: maxColor},
		{Label: "200 °C", Color: minColor},
	})

	// test case 3: Three different temperatures
	testTemperature(t, "ThreeTemperatures", `; Test G-code with three temperatures
G90
M83
M104 S200
G1 Z0.2 F720
G1 X10 Y10 Z0.2 F1800
;TYPE:Internal infill
G1 X20 Y20 Z0.2 F1800 E1.0
M109 S210
G1 X25 Y25 Z0.2 F1800 E1.5
M104 S220
G1 X30 Y30 Z0.2 F1800 E2.0
;HEIGHT:0.2
`, []legendEntry{
		{Label: "220 °C", Color: maxColor},
		{Label: "210 °C", Color: floatsToHex(lerp(temperatureColorMax[0], temperatureColorMin[0], 0.5), lerp(temperatureColorMax[1], temperatureColorMin[1], 0.5), lerp(temperatureColorMax[2], temperatureColorMin[2], 0.5))},
		{Label: "200 °C", Color: minColor},
	})

	// test case 4: Eight different temperatures (should trigger gradated legend with legendSteps)
	testTemperature(t, "EightTemperatures", `; Test G-code with eight temperatures
G90
M83
M104 S200
M109 S200
G1 Z0.2 F720
G1 X10 Y10 Z0.2 F1800
;TYPE:Internal infill
G1 X10 Y10 Z0.2 F1800 E1.0
M104 S210
M109 S210
G1 X15 Y15 Z0.2 F1800 E1.5
M104 S220
M109 S220
G1 X20 Y20 Z0.2 F1800 E2.0
M104 S230
M109 S230
G1 X25 Y25 Z0.2 F1800 E2.5
M104 S240
M109 S240
G1 X30 Y30 Z0.2 F1800 E3.0
M104 S250
M109 S250
G1 X35 Y35 Z0.2 F1800 E3.5
M104 S260
M109 S260
G1 X40 Y40 Z0.2 F1800 E4.0
M104 S270
M109 S270
G1 X45 Y45 Z0.2 F1800 E4.5
M104 S280
M109 S280
G1 X50 Y50 Z0.2 F1800 E5.0
M104 S290
M109 S290
G1 X55 Y55 Z0.2 F1800 E5.5
;HEIGHT:0.2
`, []legendEntry{
		{Label: "290 °C", Color: maxColor},
		{Label: "275 °C", Color: floatsToHex(lerp(temperatureColorMax[0], temperatureColorMin[0], 1.0/6.0), lerp(temperatureColorMax[1], temperatureColorMin[1], 1.0/6.0), lerp(temperatureColorMax[2], temperatureColorMin[2], 1.0/6.0))},
		{Label: "260 °C", Color: floatsToHex(lerp(temperatureColorMax[0], temperatureColorMin[0], 2.0/6.0), lerp(temperatureColorMax[1], temperatureColorMin[1], 2.0/6.0), lerp(temperatureColorMax[2], temperatureColorMin[2], 2.0/6.0))},
		{Label: "245 °C", Color: floatsToHex(lerp(temperatureColorMax[0], temperatureColorMin[0], 3.0/6.0), lerp(temperatureColorMax[1], temperatureColorMin[1], 3.0/6.0), lerp(temperatureColorMax[2], temperatureColorMin[2], 3.0/6.0))},
		{Label: "230 °C", Color: floatsToHex(lerp(temperatureColorMax[0], temperatureColorMin[0], 4.0/6.0), lerp(temperatureColorMax[1], temperatureColorMin[1], 4.0/6.0), lerp(temperatureColorMax[2], temperatureColorMin[2], 4.0/6.0))},
		{Label: "215 °C", Color: floatsToHex(lerp(temperatureColorMax[0], temperatureColorMin[0], 5.0/6.0), lerp(temperatureColorMax[1], temperatureColorMin[1], 5.0/6.0), lerp(temperatureColorMax[2], temperatureColorMin[2], 5.0/6.0))},
		{Label: "200 °C", Color: minColor},
	})
}

func testFeedrate(t *testing.T, name string, gcodeContent string, expectedFeedrate []legendEntry) {
	t.Run(name, func(t *testing.T) {
		// create a temporary directory for test files
		tempDir, err := os.MkdirTemp("", "ptp_test")
		if err != nil {
			t.Fatal(err)
		}
		defer os.RemoveAll(tempDir)

		gcodePath := filepath.Join(tempDir, "test.gcode")
		err = os.WriteFile(gcodePath, []byte(gcodeContent), 0644)
		if err != nil {
			t.Fatal(err)
		}

		// set up output paths
		outPath := filepath.Join(tempDir, "output.ptp")
		legendPath := outPath + ".legend"

		// arguments for GenerateToolpath
		argv := []string{
			gcodePath,     // input file
			outPath,       // output file
			"0.4",         // initialExtrusionWidth
			"0.2",         // initialLayerHeight
			"0.0",         // zOffset
			"false",       // brimIsSkirt
			"1.0,0.0,0.0", // toolColors (red for tool 0)
		}

		// call GenerateToolpath
		GenerateToolpath(argv)

		// read the legend file
		legendData, err := os.ReadFile(legendPath)
		if err != nil {
			t.Fatal(err)
		}

		// parse the legend JSON
		var legend ptpLegend
		err = json.Unmarshal(legendData, &legend)
		if err != nil {
			t.Fatal(err)
		}

		if len(legend.Feedrate) != len(expectedFeedrate) {
			t.Errorf("Expected %d legend entries, got %d", len(expectedFeedrate), len(legend.Feedrate))
			return
		}

		for i, expected := range expectedFeedrate {
			actual := legend.Feedrate[i]
			if actual.Label != expected.Label || actual.Color != expected.Color {
				t.Errorf("Entry %d: expected [%s, %s], got [%s, %s]", i, expected.Label, expected.Color, actual.Label, actual.Color)
			}
		}
	})
}

func TestFeedrateEdgeCases(t *testing.T) {
	// expected colors
	minColor := floatsToHex(feedrateColorMin[0], feedrateColorMin[1], feedrateColorMin[2])
	maxColor := floatsToHex(feedrateColorMax[0], feedrateColorMax[1], feedrateColorMax[2])

	// test case 1: Single feedrate
	testFeedrate(t, "SingleFeedrate", `; Test G-code with single feedrate
G90
M83
M104 S200
M109 S200
G1 Z0.2 F720
G1 X10 Y10 Z0.2 F1800
;TYPE:Internal infill
G1 X20 Y20 Z0.2 F1800 E1.0
G1 X30 Y30 Z0.2 F1800 E2.0
;HEIGHT:0.2
`, []legendEntry{{
		Label: "1800 mm/min",
		Color: maxColor,
	}})

	// test case 2: Two different feedrates
	testFeedrate(t, "TwoFeedrates", `; Test G-code with two feedrates
G90
M83
M104 S200
M109 S200
G1 Z0.2 F720
G1 X10 Y10 Z0.2 F1800
;TYPE:Internal infill
G1 X20 Y20 Z0.2 F1800 E1.0
G1 X30 Y30 Z0.2 F2400 E2.0
;HEIGHT:0.2
`, []legendEntry{
		{Label: "2400 mm/min", Color: maxColor},
		{Label: "1800 mm/min", Color: minColor},
	})

	// test case 3: Three different feedrates
	testFeedrate(t, "ThreeFeedrates", `; Test G-code with three feedrates
G90
M83
M104 S200
M109 S200
G1 Z0.2 F720
G1 X10 Y10 Z0.2 F1800
;TYPE:Internal infill
G1 X20 Y20 Z0.2 F1800 E1.0
G1 X25 Y25 Z0.2 F2100 E1.5
G1 X30 Y30 Z0.2 F2400 E2.0
;HEIGHT:0.2
`, []legendEntry{
		{Label: "2400 mm/min", Color: maxColor},
		{Label: "2100 mm/min", Color: floatsToHex(lerp(feedrateColorMin[0], feedrateColorMax[0], 0.5), lerp(feedrateColorMin[1], feedrateColorMax[1], 0.5), lerp(feedrateColorMin[2], feedrateColorMax[2], 0.5))},
		{Label: "1800 mm/min", Color: minColor},
	})

	// test case 4: Eight different feedrates (should trigger gradated legend with legendSteps)
	testFeedrate(t, "EightFeedrates", `; Test G-code with eight feedrates
G90
M83
M104 S200
M109 S200
G1 Z0.2 F720
G1 X10 Y10 Z0.2 F1800
;TYPE:Internal infill
G1 X10 Y10 Z0.2 F1800 E1.0
G1 X15 Y15 Z0.2 F1950 E1.5
G1 X20 Y20 Z0.2 F2100 E2.0
G1 X25 Y25 Z0.2 F2250 E2.5
G1 X30 Y30 Z0.2 F2400 E3.0
G1 X35 Y35 Z0.2 F2550 E3.5
G1 X40 Y40 Z0.2 F2700 E4.0
G1 X45 Y45 Z0.2 F2850 E4.5
;HEIGHT:0.2
`, []legendEntry{
		{Label: "2850 mm/min", Color: maxColor},
		{Label: "2675 mm/min", Color: floatsToHex(lerp(feedrateColorMax[0], feedrateColorMin[0], 1.0/6.0), lerp(feedrateColorMax[1], feedrateColorMin[1], 1.0/6.0), lerp(feedrateColorMax[2], feedrateColorMin[2], 1.0/6.0))},
		{Label: "2500 mm/min", Color: floatsToHex(lerp(feedrateColorMax[0], feedrateColorMin[0], 2.0/6.0), lerp(feedrateColorMax[1], feedrateColorMin[1], 2.0/6.0), lerp(feedrateColorMax[2], feedrateColorMin[2], 2.0/6.0))},
		{Label: "2325 mm/min", Color: floatsToHex(lerp(feedrateColorMax[0], feedrateColorMin[0], 3.0/6.0), lerp(feedrateColorMax[1], feedrateColorMin[1], 3.0/6.0), lerp(feedrateColorMax[2], feedrateColorMin[2], 3.0/6.0))},
		{Label: "2150 mm/min", Color: floatsToHex(lerp(feedrateColorMax[0], feedrateColorMin[0], 4.0/6.0), lerp(feedrateColorMax[1], feedrateColorMin[1], 4.0/6.0), lerp(feedrateColorMax[2], feedrateColorMin[2], 4.0/6.0))},
		{Label: "1975 mm/min", Color: floatsToHex(lerp(feedrateColorMax[0], feedrateColorMin[0], 5.0/6.0), lerp(feedrateColorMax[1], feedrateColorMin[1], 5.0/6.0), lerp(feedrateColorMax[2], feedrateColorMin[2], 5.0/6.0))},
		{Label: "1800 mm/min", Color: minColor},
	})
}

func testLayerHeight(t *testing.T, name string, gcodeContent string, expectedLayerHeight []legendEntry) {
	t.Run(name, func(t *testing.T) {
		// create a temporary directory for test files
		tempDir, err := os.MkdirTemp("", "ptp_test")
		if err != nil {
			t.Fatal(err)
		}
		defer os.RemoveAll(tempDir)

		gcodePath := filepath.Join(tempDir, "test.gcode")
		err = os.WriteFile(gcodePath, []byte(gcodeContent), 0644)
		if err != nil {
			t.Fatal(err)
		}

		// set up output paths
		outPath := filepath.Join(tempDir, "output.ptp")
		legendPath := outPath + ".legend"

		// arguments for GenerateToolpath
		argv := []string{
			gcodePath,     // input file
			outPath,       // output file
			"0.4",         // initialExtrusionWidth
			"0.2",         // initialLayerHeight
			"0.0",         // zOffset
			"false",       // brimIsSkirt
			"1.0,0.0,0.0", // toolColors (red for tool 0)
		}

		// call GenerateToolpath
		GenerateToolpath(argv)

		// read the legend file
		legendData, err := os.ReadFile(legendPath)
		if err != nil {
			t.Fatal(err)
		}

		// parse the legend JSON
		var legend ptpLegend
		err = json.Unmarshal(legendData, &legend)
		if err != nil {
			t.Fatal(err)
		}

		if len(legend.LayerHeight) != len(expectedLayerHeight) {
			t.Errorf("Expected %d legend entries, got %d", len(expectedLayerHeight), len(legend.LayerHeight))
			return
		}

		for i, expected := range expectedLayerHeight {
			actual := legend.LayerHeight[i]
			if actual.Label != expected.Label || actual.Color != expected.Color {
				t.Errorf("Entry %d: expected [%s, %s], got [%s, %s]", i, expected.Label, expected.Color, actual.Label, actual.Color)
			}
		}
	})
}

func TestLayerHeightEdgeCases(t *testing.T) {
	// expected colors
	minColor := floatsToHex(layerHeightColorMin[0], layerHeightColorMin[1], layerHeightColorMin[2])
	maxColor := floatsToHex(layerHeightColorMax[0], layerHeightColorMax[1], layerHeightColorMax[2])

	// test case 1: Single layer height
	testLayerHeight(t, "SingleLayerHeight", `; Test G-code with single layer height
G90
M83
M104 S200
M109 S200
G1 Z0.2 F720
G1 X10 Y10 Z0.2 F1800
;TYPE:Internal infill
G1 X20 Y20 Z0.2 F1800 E1.0
G1 X30 Y30 Z0.2 F1800 E2.0
;LAYER_CHANGE
;Z:0.4
;HEIGHT:0.2
G1 Z0.4 F720
G1 X10 Y10 Z0.4 F1800
;TYPE:Internal infill
G1 X20 Y20 Z0.4 F1800 E3.0
G1 X30 Y30 Z0.4 F1800 E4.0
`, []legendEntry{{
		Label: "0.2 mm",
		Color: maxColor,
	}})

	// test case 2: Two different layer heights
	testLayerHeight(t, "TwoLayerHeights", `; Test G-code with two layer heights
G90
M83
M104 S200
M109 S200
G1 Z0.0 
G1 X10 Y10 Z0.0 F1800
;TYPE:Internal infill
G1 X20 Y20 Z0.0 F1800 E0.5
;LAYER_CHANGE
;Z:0.2
;HEIGHT:0.2
G1 Z0.2 F720
G1 X10 Y10 Z0.2 F1800
;TYPE:Internal infill
G1 X20 Y20 Z0.2 F1800 E1.0
;LAYER_CHANGE
;Z:0.35
;HEIGHT:0.15
G1 Z0.35 F720
G1 X10 Y10 Z0.35 F1800
;TYPE:Internal infill
G1 X20 Y20 Z0.35 F1800 E2.0
`, []legendEntry{
		{Label: "0.2 mm", Color: maxColor},
		{Label: "0.15 mm", Color: minColor},
	})

	// test case 3: Three different layer heights
	testLayerHeight(t, "ThreeLayerHeights", `; Test G-code with three layer heights
G90
M83
M104 S200
M109 S200
G1 Z0.2 F720
G1 X10 Y10 Z0.2 F1800
;TYPE:Internal infill
G1 X20 Y20 Z0.2 F1800 E1.0
;LAYER_CHANGE
;Z:0.35
;HEIGHT:0.15
G1 Z0.35 F720
G1 X10 Y10 Z0.35 F1800
;TYPE:Internal infill
G1 X20 Y20 Z0.35 F1800 E1.5
;LAYER_CHANGE
;Z:0.45
;HEIGHT:0.10
G1 Z0.45 F720
G1 X10 Y10 Z0.45 F1800
;TYPE:Internal infill
G1 X20 Y20 Z0.45 F1800 E2.0
;LAYER_CHANGE
;Z:0.65
;HEIGHT:0.2
G1 Z0.65 F720
G1 X10 Y10 Z0.65 F1800
;TYPE:Internal infill
G1 X20 Y20 Z0.65 F1800 E2.5
`, []legendEntry{
		{Label: "0.2 mm", Color: maxColor},
		{Label: "0.15 mm", Color: floatsToHex(lerp(layerHeightColorMin[0], layerHeightColorMax[0], 0.5), lerp(layerHeightColorMin[1], layerHeightColorMax[1], 0.5), lerp(layerHeightColorMin[2], layerHeightColorMax[2], 0.5))},
		{Label: "0.1 mm", Color: minColor},
	})

	// test case 4: Eight different layer heights (should trigger gradated legend with legendSteps)
	testLayerHeight(t, "EightLayerHeights", `; Test G-code with eight layer heights
G90
M83
M104 S200
M109 S200
G1 Z0.2 F720
G1 X10 Y10 Z0.2 F1800
;TYPE:Internal infill
G1 X10 Y10 Z0.2 F1800 E1.0
;LAYER_CHANGE
;Z:0.3
;HEIGHT:0.1
G1 Z0.3 F720
G1 X10 Y10 Z0.3 F1800
;TYPE:Internal infill
G1 X15 Y15 Z0.3 F1800 E1.5
;LAYER_CHANGE
;Z:0.38
;HEIGHT:0.08
G1 Z0.38 F720
G1 X10 Y10 Z0.38 F1800
;TYPE:Internal infill
G1 X20 Y20 Z0.38 F1800 E2.0
;LAYER_CHANGE
;Z:0.44
;HEIGHT:0.06
G1 Z0.44 F720
G1 X10 Y10 Z0.44 F1800
;TYPE:Internal infill
G1 X25 Y25 Z0.44 F1800 E2.5
;LAYER_CHANGE
;Z:0.48
;HEIGHT:0.04
G1 Z0.48 F720
G1 X10 Y10 Z0.48 F1800
;TYPE:Internal infill
G1 X30 Y30 Z0.48 F1800 E3.0
;LAYER_CHANGE
;Z:0.5
;HEIGHT:0.02
G1 Z0.5 F720
G1 X10 Y10 Z0.5 F1800
;TYPE:Internal infill
G1 X35 Y35 Z0.5 F1800 E3.5
;LAYER_CHANGE
;Z:0.51
;HEIGHT:0.01
G1 Z0.51 F720
G1 X10 Y10 Z0.51 F1800
;TYPE:Internal infill
G1 X40 Y40 Z0.51 F1800 E4.0
;LAYER_CHANGE
;Z:0.71
;HEIGHT:0.2
G1 Z0.71 F720
G1 X10 Y10 Z0.71 F1800
;TYPE:Internal infill
G1 X45 Y45 Z0.71 F1800 E4.5
`, []legendEntry{
		{Label: "0.2 mm", Color: maxColor},
		{Label: "0.168 mm", Color: floatsToHex(lerp(layerHeightColorMax[0], layerHeightColorMin[0], 1.0/6.0), lerp(layerHeightColorMax[1], layerHeightColorMin[1], 1.0/6.0), lerp(layerHeightColorMax[2], layerHeightColorMin[2], 1.0/6.0))},
		{Label: "0.136 mm", Color: floatsToHex(lerp(layerHeightColorMax[0], layerHeightColorMin[0], 2.0/6.0), lerp(layerHeightColorMax[1], layerHeightColorMin[1], 2.0/6.0), lerp(layerHeightColorMax[2], layerHeightColorMin[2], 2.0/6.0))},
		{Label: "0.104 mm", Color: floatsToHex(lerp(layerHeightColorMax[0], layerHeightColorMin[0], 3.0/6.0), lerp(layerHeightColorMax[1], layerHeightColorMin[1], 3.0/6.0), lerp(layerHeightColorMax[2], layerHeightColorMin[2], 3.0/6.0))},
		{Label: "0.072 mm", Color: floatsToHex(lerp(layerHeightColorMax[0], layerHeightColorMin[0], 4.0/6.0), lerp(layerHeightColorMax[1], layerHeightColorMin[1], 4.0/6.0), lerp(layerHeightColorMax[2], layerHeightColorMin[2], 4.0/6.0))},
		{Label: "0.04 mm", Color: floatsToHex(lerp(layerHeightColorMax[0], layerHeightColorMin[0], 5.0/6.0), lerp(layerHeightColorMax[1], layerHeightColorMin[1], 5.0/6.0), lerp(layerHeightColorMax[2], layerHeightColorMin[2], 5.0/6.0))},
		{Label: "0.01 mm", Color: minColor},
	})
}
