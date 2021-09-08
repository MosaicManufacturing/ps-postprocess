package palette

import "math"

type Material struct {
    ID string
    Index int
    Name string
    Color string
}

type SpliceSettings struct {
    IngoingID string
    OutgoingID string
    HeatFactor float32
    CompressionFactor float32
    CoolingFactor float32
    Reverse bool
}

type Palette struct {
    // general settings
    Type Type
    Model Model
    Makerbot5thGen bool
    MaterialMeta []Material
    SpliceSettings []SpliceSettings

    // output settings
    PrintExtruder int
    FirmwarePurge float32 // mm
    BowdenTubeLength float32 // mm
    TransitionTarget float32 // 0..100

    // P2/P3 settings
    PrinterID string
    ConnectedMode bool

    // P1 settings
    LoadingOffset int
    PrintValue int
    CalibrationLength float32
}

// todo: some way to load from disk or JSON?

func (p Palette) GetInputCount() int {
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
        return "mcfx"
    }
    return "mcf"
}

func (p Palette) GetSpliceCore() string {
    switch p.Type {
    case TypeP1: return "P"
    case TypeP2:
        switch p.Model {
        case ModelP2: return "SC"
        case ModelP2Pro: return "SCP"
        case ModelP2S: return "SCS"
        case ModelP2SPro: return "SCSP"
        }
    case TypeP3:
        switch p.Model {
        case ModelP3: return "P3-SC"
        case ModelP3Pro: return "P3-SCP"
        }
    }
    return ""
}

func (p Palette) GetFirstSpliceMinLength() float32 {
    if p.Type == TypeP1 {
        return MinFirstSpliceLengthP1
    }
    if p.Type == TypeP2 {
        return MinFirstSpliceLengthP2
    }
    return MinFirstSpliceLengthP3
}

func (p Palette) GetPulsesPerMM() float32 {
    if p.PrintValue == 0 || p.CalibrationLength == 0 {
        return 0
    }
    ppm := float64(p.PrintValue) / float64(p.CalibrationLength + p.FirmwarePurge)
    return float32(math.Max(20, math.Min(40, ppm)))
}

func (p Palette) GetEffectiveLoadingOffset() float32 {
    ppm := p.GetPulsesPerMM()
    if ppm == 0 {
        return 0
    }
    return (float32(p.LoadingOffset) / ppm) + CutterToScrollWheel
}
