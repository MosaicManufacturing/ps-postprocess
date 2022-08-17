package msf

import (
	"mosaicmfg.com/ps-postprocess/sequences"
	"testing"
)

// transition length >= 80 mm
// no infill transitions
// no variable-length transitions
// no tower brims needed
func Test_TowersBasic(t *testing.T) {
	palette := getTestPalette(80)
	printContent := `
;START_OF_PRINT
T0
G1 Z0
;LAYER_CHANGE
;Z:0.2
;HEIGHT:0.2
; printing object model id:0 copy 0
G1 X0 Y0 Z0.2 F1800 ; move to a real XYZ position
G1 E200 F2400 ; avoid testing logic specific to first piece length
G92 E0
T1
G1 E80 F2400
G92 E0
T2
G1 E80 F2400
G92 E0
T1
G1 E100 F2400
; stop printing object model id:0 copy 0
`
	expectedTransitions := []Transition{
		{
			Layer:            0,
			From:             0,
			To:               1,
			TotalExtrusion:   200,
			TransitionLength: 80,
			PurgeLength:      80,
			UsableInfill:     0,
		},
		{
			Layer:            0,
			From:             1,
			To:               2,
			TotalExtrusion:   280,
			TransitionLength: 80,
			PurgeLength:      80,
			UsableInfill:     0,
		},
		{
			Layer:            0,
			From:             2,
			To:               1,
			TotalExtrusion:   360,
			TransitionLength: 80,
			PurgeLength:      80,
			UsableInfill:     0,
		},
	}
	preflightResults := testTowerPreflight(t, &palette, printContent, expectedTransitions)
	testTowerOutput(t, &palette, printContent, &preflightResults, sequences.NewLocals())
}

// transition length < 80 mm, but enough extrusion on print to avoid piece length issues
// no infill transitions
// no variable-length transitions
// no tower brims needed
func Test_TowersShortTransitions(t *testing.T) {
	palette := getTestPalette(30)
	printContent := `
;START_OF_PRINT
T0
G1 Z0
;LAYER_CHANGE
;Z:0.2
;HEIGHT:0.2
; printing object model id:0 copy 0
G1 X0 Y0 Z0.2 F1800 ; move to a real XYZ position
G1 E200 F2400 ; avoid testing logic specific to first piece length
G92 E0
T1
G1 E80 F2400 ; enough to ensure no extra transition is needed
G92 E0
T2
G1 E90 F2400 ; enough to ensure no extra transition is needed
G92 E0
T1
G1 E100 F2400
; stop printing object model id:0 copy 0
`
	expectedTransitions := []Transition{
		{
			Layer:            0,
			From:             0,
			To:               1,
			TotalExtrusion:   200,
			TransitionLength: 30,
			PurgeLength:      30,
			UsableInfill:     0,
		},
		{
			Layer:            0,
			From:             1,
			To:               2,
			TotalExtrusion:   280,
			TransitionLength: 30,
			PurgeLength:      30,
			UsableInfill:     0,
		},
		{
			Layer:            0,
			From:             2,
			To:               1,
			TotalExtrusion:   370,
			TransitionLength: 30,
			PurgeLength:      30,
			UsableInfill:     0,
		},
	}
	preflightResults := testTowerPreflight(t, &palette, printContent, expectedTransitions)
	testTowerOutput(t, &palette, printContent, &preflightResults, sequences.NewLocals())
}

// todo: explanatory comment
func Test_TowersExtraPurge(t *testing.T) {
	palette := getTestPalette(30)
	printContent := `
;START_OF_PRINT
T0
G1 Z0
;LAYER_CHANGE
;Z:0.2
;HEIGHT:0.2
; printing object model id:0 copy 0
G1 X0 Y0 Z0.2 F1800 ; move to a real XYZ position
G1 E200 F2400 ; avoid testing logic specific to first piece length
G92 E0
T1
G1 E80 F2400 ; enough to ensure no extra transition is needed
G92 E0
T2
G1 E10 F2400 ; will need extra transition
G92 E0
T1
G1 E10 F2400 ; will need extra transition
G92 E0
T2
G1 E10 F2400 ; will need extra transition
G92 E0
T1
G1 E10 F2400 ; will need extra transition
G92 E0
T2
G1 E10 F2400 ; will need extra transition
G92 E0
T1
G1 E10 F2400 ; will need extra transition
G92 E0
T2
G1 E50 F2400
; stop printing object model id:0 copy 0
`
	expectedTransitions := []Transition{
		{
			Layer:            0,
			From:             0,
			To:               1,
			TotalExtrusion:   200,
			TransitionLength: 30,
			PurgeLength:      30,
			UsableInfill:     0,
		},
		{
			Layer:            0,
			From:             1,
			To:               2,
			TotalExtrusion:   280,
			TransitionLength: 30,
			PurgeLength:      30,
			UsableInfill:     0,
		},
		{
			Layer:            0,
			From:             2,
			To:               1,
			TotalExtrusion:   290,
			TransitionLength: 30,
			PurgeLength:      70,
			UsableInfill:     0,
		},
		{
			Layer:            0,
			From:             1,
			To:               2,
			TotalExtrusion:   300,
			TransitionLength: 30,
			PurgeLength:      70,
			UsableInfill:     0,
		},
		{
			Layer:            0,
			From:             2,
			To:               1,
			TotalExtrusion:   310,
			TransitionLength: 30,
			PurgeLength:      70,
			UsableInfill:     0,
		},
		{
			Layer:            0,
			From:             1,
			To:               2,
			TotalExtrusion:   320,
			TransitionLength: 30,
			PurgeLength:      70,
			UsableInfill:     0,
		},
		{
			Layer:            0,
			From:             2,
			To:               1,
			TotalExtrusion:   330,
			TransitionLength: 30,
			PurgeLength:      70,
			UsableInfill:     0,
		},
		{
			Layer:            0,
			From:             1,
			To:               2,
			TotalExtrusion:   340,
			TransitionLength: 30,
			PurgeLength:      70,
			UsableInfill:     0,
		},
	}
	preflightResults := testTowerPreflight(t, &palette, printContent, expectedTransitions)
	testTowerOutput(t, &palette, printContent, &preflightResults, sequences.NewLocals())
}

// todo: explanatory comment
func Test_TowersMultipleExtraPurges(t *testing.T) {
	palette := getTestPalette(30)
	printContent := `
;START_OF_PRINT
T0
G1 Z0
;LAYER_CHANGE
;Z:0.2
;HEIGHT:0.2
; printing object model id:0 copy 0
G1 X0 Y0 Z0.2 F1800 ; move to a real XYZ position
G1 E200 F2400 ; avoid testing logic specific to first piece length
G92 E0
T1
G1 E10 F2400 ; will need extra transition
G92 E0
T2
G1 E10 F2400 ; will need extra transition
G92 E0
T1
G1 E50 F2400
; stop printing object model id:0 copy 0
`
	expectedTransitions := []Transition{
		{
			Layer:            0,
			From:             0,
			To:               1,
			TotalExtrusion:   200,
			TransitionLength: 30,
			PurgeLength:      30,
			UsableInfill:     0,
		},
		{
			Layer:            0,
			From:             1,
			To:               2,
			TotalExtrusion:   210,
			TransitionLength: 30,
			PurgeLength:      70,
			UsableInfill:     0,
		},
		{
			Layer:            0,
			From:             2,
			To:               1,
			TotalExtrusion:   220,
			TransitionLength: 30,
			PurgeLength:      70,
			UsableInfill:     0,
		},
	}
	preflightResults := testTowerPreflight(t, &palette, printContent, expectedTransitions)
	testTowerOutput(t, &palette, printContent, &preflightResults, sequences.NewLocals())
}

func Test_4ColorTowerShortTransitions(t *testing.T) {
	testWithInputFiles(t, "1")
}

func Test_4ColorTowerShortTransitionsWithInfill(t *testing.T) {
	testWithInputFiles(t, "2")
}

func Test_8ColorTowerShortTransitions(t *testing.T) {
	testWithInputFiles(t, "3")
}

func Test_8ColorTowerShortTransitionsWithInfill(t *testing.T) {
	testWithInputFiles(t, "4")
}

func Test_8ColorTowerVariableTransitions(t *testing.T) {
	testWithInputFiles(t, "5")
}

// TODO: add tests for the following cases (including combinations):
//   - first piece length handling (auto-brims)
//   - inclusion of sparse layers
//   - infill transitions enabled
//   - variable transition lengths
