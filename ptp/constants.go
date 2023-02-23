package ptp

const ptpVersion = uint8(6)

const (
	floatBytes  = 4
	uint8Bytes  = 1
	uint32Bytes = 4
)

// segments with all dimension deltas smaller than this will be skipped
const skipThreshold = 0.01

// tolerance used by collinearity-checking functions
const collinearityEpsilon = 10e-5

// in-file header only contains version information
// (buffer offsets and sizes are now stored in the legend file)
const headerSize = uint32(4)

var (
	colorWhite      = [3]float32{0xff / 255.0, 0xff / 255.0, 0xff / 255.0} // #ffffff
	colorGreen      = [3]float32{0x6b / 255.0, 0xa7 / 255.0, 0x31 / 255.0} // #6ba731
	colorDarkGrey   = [3]float32{0x32 / 255.0, 0x29 / 255.0, 0x2f / 255.0} // #32292f
	colorYellow     = [3]float32{0xf7 / 255.0, 0xb5 / 255.0, 0x38 / 255.0} // #f7b538
	colorPurple     = [3]float32{0x3d / 255.0, 0x31 / 255.0, 0x5b / 255.0} // #3d315b
	colorLilac      = [3]float32{0x97 / 255.0, 0x89 / 255.0, 0xba / 255.0} // #9789ba
	colorLightGreen = [3]float32{0x71 / 255.0, 0xb5 / 255.0, 0x98 / 255.0} // #71b598
	colorTeal       = [3]float32{0x3b / 255.0, 0x8e / 255.0, 0xa5 / 255.0} // #3b8ea5
	colorRed        = [3]float32{0xdb / 255.0, 0x32 / 255.0, 0x4d / 255.0} // #db324d
	colorPink       = [3]float32{0xd6 / 255.0, 0x7a / 255.0, 0x89 / 255.0} // #d67a89
	colorOrange     = [3]float32{0xd5 / 255.0, 0x57 / 255.0, 0x3b / 255.0} // #d5573b
	colorSky        = [3]float32{0xd4 / 255.0, 0xde / 255.0, 0xff / 255.0} // #d4deff
	colorLightGrey  = [3]float32{0xd1 / 255.0, 0xd1 / 255.0, 0xd1 / 255.0} // #d1d1d1
)

const (
	travelExtrusionWidth = 0.08
	travelLayerHeight    = 0.08
	travelTool           = -1
)

var travelColor = [3]float32{0x99 / 255.0, 0x99 / 255.0, 0x99 / 255.0}

type PathType int

const (
	PathTypeUnknown PathType = iota
	PathTypeTravel
	PathTypeSequence
	PathTypeRaft
	PathTypeBrim
	PathTypeSupport
	PathTypeSupportInterface
	PathTypeInnerPerimeter
	PathTypeOuterPerimeter
	PathTypeSolidLayer
	PathTypeInfill
	PathTypeGapFill
	PathTypeBridge
	PathTypeIroning
	PathTypeTransition
	pathTypeCount
)

var pathTypeNames = map[PathType]string{
	PathTypeUnknown:          "Unknown",
	PathTypeTravel:           "Travel",
	PathTypeSequence:         "User Sequence",
	PathTypeRaft:             "Raft",
	PathTypeBrim:             "Skirt/Brim",
	PathTypeSupport:          "Support",
	PathTypeSupportInterface: "Support Interface",
	PathTypeInnerPerimeter:   "Inner Perimeter",
	PathTypeOuterPerimeter:   "Outer Perimeter",
	PathTypeSolidLayer:       "Solid Layer",
	PathTypeInfill:           "Infill",
	PathTypeGapFill:          "Gap Fill",
	PathTypeBridge:           "Bridge",
	PathTypeIroning:          "Ironing",
	PathTypeTransition:       "Transition",
}

var pathTypeColors = map[PathType][3]float32{
	PathTypeUnknown:          colorWhite,
	PathTypeTravel:           travelColor,
	PathTypeSequence:         colorDarkGrey,
	PathTypeRaft:             colorLilac,
	PathTypeBrim:             colorSky,
	PathTypeSupport:          colorPurple,
	PathTypeSupportInterface: colorLilac,
	PathTypeInnerPerimeter:   colorLightGreen,
	PathTypeOuterPerimeter:   colorTeal,
	PathTypeSolidLayer:       colorRed,
	PathTypeInfill:           colorYellow,
	PathTypeGapFill:          colorOrange,
	PathTypeBridge:           colorSky,
	PathTypeIroning:          colorPink,
	PathTypeTransition:       colorLightGrey,
}

var pathTypeColorStrings = map[PathType]string{
	PathTypeUnknown:          "#ffffff",
	PathTypeTravel:           "#999999",
	PathTypeSequence:         "#32292f",
	PathTypeRaft:             "#9789ba",
	PathTypeBrim:             "#d4deff",
	PathTypeSupport:          "#3d315b",
	PathTypeSupportInterface: "#9789ba",
	PathTypeInnerPerimeter:   "#71b598",
	PathTypeOuterPerimeter:   "#3b8ea5",
	PathTypeSolidLayer:       "#db324d",
	PathTypeInfill:           "#f7b538",
	PathTypeGapFill:          "#d5573b",
	PathTypeBridge:           "#bcd4de",
	PathTypeIroning:          "#d67a89",
	PathTypeTransition:       "#d1d1d1",
}

var feedrateColorMin = colorRed
var feedrateColorMax = colorTeal

var fanColorMin = colorRed
var fanColorMax = colorGreen

var temperatureColorMin = colorTeal
var temperatureColorMax = colorRed

var layerHeightColorMin = colorTeal
var layerHeightColorMax = colorOrange
