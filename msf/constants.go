package msf

import "math"

const EOL = "\r\n"

const (
    charLimitMSF1 = 20
    charLimitMSF2 = 32
)

var posInf = float32(math.Inf(1))
var negInf = float32(math.Inf(-1))

type Type string

const (
    TypeP1 Type = "palette"
    TypeP2 Type = "palette-2"
    TypeP3 Type = "palette-3"
)

type Model string

const (
    ModelP Model = "p"
    ModelPPlus Model = "p-plus"
    ModelP2 Model = "p2"
    ModelP2Pro Model = "p2-pro"
    ModelP2S Model = "p2s"
    ModelP2SPro Model = "p2s-pro"
    ModelP3 Model = "p3"
    ModelP3Pro Model = "p3-pro"
)

type TransitionMethod int
const (
    TransitionTower TransitionMethod = 1
    SideTransitions TransitionMethod = 2
)

const (
	MinSpliceLength = float32(80)
    MinFirstSpliceLengthP1 = float32(140)
    MinFirstSpliceLengthP2 = float32(100)
    MinFirstSpliceLengthP3 = float32(130)
)

const BowdenDefault = float32(150)

const CutterToScrollWheel = float32(760)

const PingExtrusionCounts = 600 // target extrusion between ping sequence pauses, in scroll wheel counts
const PingExtrusion = 20        // target extrusion between ping sequence pauses, mm
const Ping1PauseLength = 13000  // duration of first ping sequence pause, in ms
const Ping2PauseLength = 7000   // duration of second ping sequence pause, in ms
const PingMinSpacing = 350      // minimum distance (extrusion) between ping starts, in mm
