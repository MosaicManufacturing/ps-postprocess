package msf

import (
	"encoding/json"
	"io/ioutil"
	"math"
	"mosaicmfg.com/ps-postprocess/gcode"
	"mosaicmfg.com/ps-postprocess/printerscript"
	"strings"
)

type Material struct {
	ID         string `json:"id"`
	Index      int    `json:"index"`
	FilamentID int    `json:"filamentId"`
	Name       string `json:"name"`
	Color      string `json:"color"`
}

type SpliceSettings struct {
	IngoingID         string  `json:"ingoingId"`
	OutgoingID        string  `json:"outgoingId"`
	HeatFactor        float32 `json:"heatFactor"`
	CompressionFactor float32 `json:"compressionFactor"`
	CoolingFactor     float32 `json:"coolingFactor"`
	Reverse           bool    `json:"reverse"`
}

type Palette struct {
	// general
	Type           Type             `json:"type"`
	Model          Model            `json:"model"`
	MaterialMeta   []Material       `json:"materialMeta"`
	SpliceSettings []SpliceSettings `json:"spliceSettings"`

	// physical
	PrintExtruder    int     `json:"printExtruder"`
	FirmwarePurge    float32 `json:"firmwarePurge"`    // mm
	BowdenTubeLength float32 `json:"bowdenTubeLength"` // mm

	// slicer
	TravelSpeedXY float32 `json:"travelSpeedXY"` // mm/min
	TravelSpeedZ  float32 `json:"travelSpeedZ"`  // mm/min
	PrintBedMinX  float32 `json:"printBedMinX"`  // mm
	PrintBedMaxX  float32 `json:"printBedMaxX"`  // mm
	PrintBedMinY  float32 `json:"printBedMinY"`  // mm
	PrintBedMaxY  float32 `json:"printBedMaxY"`  // mm

	// transitions
	TransitionMethod    TransitionMethod `json:"TransitionMethod"`
	TransitionLengths   [][]float32      `json:"transitionLengths"` // mm
	TransitionTarget    float32          `json:"transitionTarget"`  // 0..100
	InfillTransitioning bool             `json:"infillTransitioning"`

	// transition tower generation
	TowerSize                 [2]float32 `json:"towerSize"`
	TowerPosition             [2]float32 `json:"towerPosition"`             // center of tower
	TowerMinDensity           float32    `json:"towerMinDensity"`           // 0..100
	TowerMinFirstLayerDensity float32    `json:"towerMinFirstLayerDensity"` // 0..100
	TowerMaxDensity           float32    `json:"towerMaxDensity"`           // 0..100
	TowerMinBrims             int        `json:"towerMinBrims"`
	TowerSpeed                []float32  `json:"towerSpeed"`               // mm/s
	FirstLayerTowerSpeed      []float32  `json:"firstLayerTowerSpeed"`     // mm/s
	TowerExtrusionWidth       float32    `json:"towerExtrusionWidth"`      // mm
	TowerExtrusionMultiplier  float32    `json:"towerExtrusionMultiplier"` // unitless
	TowerFirstLayerPerimeters bool       `json:"towerFirstLayerPerimeters"`
	InfillPerimeterOverlap    float32    `json:"infillPerimeterOverlap"` // 0..100
	RaftLayers                int        `json:"raftLayers"`
	RaftInflation             float32    `json:"raftInflation"`      // mm
	RaftExtrusionWidth        float32    `json:"raftExtrusionWidth"` // mm
	RaftStride                float32    `json:"raftStride"`         // mm
	UseFirmwareRetraction     bool       `json:"useFirmwareRetraction"`
	RetractDistance           []float32  `json:"retractDistance"` // mm
	RestartDistance           []float32  `json:"restartDistance"` // mm
	RetractFeedrate           []float32  `json:"retractFeedrate"` // mm/min
	RestartFeedrate           []float32  `json:"restartFeedrate"` // mm/min
	Wipe                      []bool     `json:"wipe"`
	ZLift                     []float32  `json:"zLift"`   // mm
	ZOffset                   float32    `json:"zOffset"` // mm

	// side transition scripting
	PreSideTransitionSequence  string `json:"preSideTransitionSequence"`
	SideTransitionSequence     string `json:"sideTransitionSequence"`
	PostSideTransitionSequence string `json:"postSideTransitionSequence"`
	PreSideTransitionScript    printerscript.Tree
	SideTransitionScript       printerscript.Tree
	PostSideTransitionScript   printerscript.Tree

	// side transitions
	SideTransitionJog        bool            `json:"sideTransitionJog"`
	SideTransitionPurgeSpeed float32         `json:"sideTransitionPurgeSpeed"` // mm/s
	SideTransitionMoveSpeed  float32         `json:"sideTransitionMoveSpeed"`  // mm/s
	SideTransitionX          float32         `json:"sideTransitionX"`          // mm
	SideTransitionY          float32         `json:"sideTransitionY"`          // mm
	SideTransitionEdge       gcode.Direction `json:"sideTransitionEdge"`
	SideTransitionEdgeOffset float32         `json:"sideTransitionEdgeOffset"` // mm

	// pings
	PingOffTowerDistance float32 `json:"pingOffTowerDistance"` // mm
	JogPauses            bool    `json:"jogPauses"`

	// P2/P3
	ClearBufferCommand string `json:"clearBufferCommand"`
	ConnectedMode      bool   `json:"connectedMode"`
	PrinterID          string `json:"printerId"`
	Filename           string `json:"filename"`

	// P1
	LoadingOffset     int     `json:"loadingOffset"`     // scroll wheel counts
	PrintValue        int     `json:"printValue"`        // scroll wheel counts
	CalibrationLength float32 `json:"calibrationLength"` // mm
}

func LoadPaletteFromFile(path string) (Palette, error) {
	palette := Palette{
		BowdenTubeLength: BowdenDefault,
	}
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return palette, err
	}
	if err := json.Unmarshal(bytes, &palette); err != nil {
		return palette, err
	}

	// apply defaults
	if palette.BowdenTubeLength < BowdenDefault {
		palette.BowdenTubeLength = BowdenDefault
	}

	// lex and parse scripts just once now, and re-use the parse trees when evaluating
	palette.PreSideTransitionSequence = printerscript.Normalize(palette.PreSideTransitionSequence)
	if len(strings.TrimSpace(palette.PreSideTransitionSequence)) > 0 {
		tree, err := printerscript.LexAndParse(palette.PreSideTransitionSequence)
		if err != nil {
			return palette, err
		}
		palette.PreSideTransitionScript = tree
	}
	palette.SideTransitionSequence = printerscript.Normalize(palette.SideTransitionSequence)
	if len(strings.TrimSpace(palette.SideTransitionSequence)) > 0 {
		tree, err := printerscript.LexAndParse(palette.SideTransitionSequence)
		if err != nil {
			return palette, err
		}
		palette.SideTransitionScript = tree
	}
	palette.PostSideTransitionSequence = printerscript.Normalize(palette.PostSideTransitionSequence)
	if len(strings.TrimSpace(palette.PostSideTransitionSequence)) > 0 {
		tree, err := printerscript.LexAndParse(palette.PostSideTransitionSequence)
		if err != nil {
			return palette, err
		}
		palette.PostSideTransitionScript = tree
	}

	return palette, nil
}

func (p Palette) ProductName() string {
	if p.Type == TypeElement {
		return "Element"
	}
	return "Palette"
}

func (p Palette) SupportsPings() bool {
	return p.Type != TypeElement
}

func (p Palette) GetInputCount() int {
	if p.Type == TypeElement {
		return 8
	}
	if p.Type == TypeP3 && p.Model == ModelP3Pro {
		return 8
	}
	return 4
}

func (p Palette) GetAccessoryModeExtension() string {
	switch p.Type {
	case TypeP2:
		return "maf"
	case TypeP3:
		return "mafx"
	}
	return "msf"
}

func (p Palette) GetConnectedModeExtension() string {
	switch p.Type {
	case TypeP3:
		fallthrough
	case TypeElement:
		return "mcfx"
	}
	return "mcf"
}

func (p Palette) GetSpliceCore() string {
	switch p.Type {
	case TypeP1:
		return "P"
	case TypeP2:
		switch p.Model {
		case ModelP2:
			return "SC"
		case ModelP2Pro:
			return "SCP"
		case ModelP2S:
			return "SCS"
		case ModelP2SPro:
			return "SCSP"
		}
	case TypeP3:
		switch p.Model {
		case ModelP3:
			return "P3-SC"
		case ModelP3Pro:
			return "P3-SCP"
		}
	}
	return ""
}

func (p Palette) GetFirstSpliceMinLength() float32 {
	switch p.Type {
	case TypeElement:
		return MinFirstSpliceLengthElement
	case TypeP1:
		return MinFirstSpliceLengthP1
	case TypeP2:
		return MinFirstSpliceLengthP2
	case TypeP3:
		return MinFirstSpliceLengthP3
	}
	return 0
}

func (p Palette) GetSpliceMinLength() float32 {
	if p.Type == TypeElement {
		return MinSpliceLengthElement
	}
	return MinSpliceLength
}

func (p Palette) GetPulsesPerMM() float32 {
	if p.PrintValue == 0 || p.CalibrationLength == 0 {
		return 0
	}
	ppm := float64(p.PrintValue) / float64(p.CalibrationLength+p.FirmwarePurge)
	return float32(math.Max(20, math.Min(40, ppm)))
}

func (p Palette) GetPingExtrusion() float32 {
	if !p.SupportsPings() {
		return 0
	}
	if p.Type == TypeP1 {
		return PingExtrusionCounts / p.GetPulsesPerMM()
	}
	return PingExtrusion
}

func (p Palette) GetEffectiveLoadingOffset() float32 {
	ppm := p.GetPulsesPerMM()
	if ppm == 0 {
		return 0
	}
	return (float32(p.LoadingOffset) / ppm) + CutterToScrollWheel
}

func (p Palette) GetTransitionLength(toTool, fromTool int) float32 {
	return p.TransitionLengths[toTool][fromTool]
}

func (p Palette) GetTowerPrintSpeed(toTool, fromTool, layerIndex int) float32 {
	// use the slower of the two material settings for this transition
	toSpeed := p.TowerSpeed[toTool]
	fromSpeed := p.TowerSpeed[fromTool]
	if layerIndex == 0 {
		toSpeed = p.FirstLayerTowerSpeed[toTool]
		fromSpeed = p.FirstLayerTowerSpeed[fromTool]
	}
	if toSpeed < fromSpeed {
		return toSpeed * 60
	}
	return fromSpeed * 60
}
