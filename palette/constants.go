package palette

type Type string
type Model string

const EOL = "\r\n"

const (
    charLimitMSF1 = 20
    charLimitMSF2 = 32
)

const (
    TypeP1 Type = "palette"
    TypeP2 Type = "palette-2"
    TypeP3 Type = "palette-3"
)

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

const (
	MinSpliceLength = float32(80)
    MinFirstSpliceLengthP1 = float32(140)
    MinFirstSpliceLengthP2 = float32(100)
    MinFirstSpliceLengthP3 = float32(130)
)

const BowdenDefault = float32(150)

const CutterToScrollWheel = float32(760)

const PingExtrusionCounts = 600 // target extrusion between ping sequence pauses, in scroll wheel counts
const Ping1PauseLength = 13000  // ms; duration of first ping sequence pause
const Ping2PauseLength = 7000   // ms; duration of second ping sequence pause
const PingMinSpacing = 350      // mm; minimum distance (extrusion) between pings
