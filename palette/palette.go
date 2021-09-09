package palette

import (
    "encoding/json"
    "io/ioutil"
    "math"
)

type Material struct {
    ID string `json:"id"`
    Index int `json:"index"`
    Name string `json:"name"`
    Color string `json:"color"`
}

type SpliceSettings struct {
    IngoingID string `json:"ingoingId"`
    OutgoingID string `json:"outgoingId"`
    HeatFactor float32 `json:"heatFactor"`
    CompressionFactor float32 `json:"compressionFactor"`
    CoolingFactor float32 `json:"coolingFactor"`
    Reverse bool `json:"reverse"`
}

type Palette struct {
    // general settings
    Type Type `json:"type"`
    Model Model `json:"model"`
    Makerbot5thGen bool `json:"makerbot5thGen"`
    MaterialMeta []Material `json:"materialMeta"`
    SpliceSettings []SpliceSettings `json:"spliceSettings"`

    // output settings
    PrintExtruder int `json:"printExtruder"`
    FirmwarePurge float32 `json:"firmwarePurge"` // mm
    BowdenTubeLength float32 `json:"bowdenTubeLength"` // mm
    TransitionLengths [][]float32 `json:"transitionLengths"` // mm
    TransitionTarget float32 `json:"transitionTarget"` // 0..100

    // P2/P3 settings
    PrinterID string `json:"printerId"`
    ConnectedMode bool `json:"connectedMode"`

    // P1 settings
    LoadingOffset int `json:"loadingOffset"`
    PrintValue int `json:"printValue"`
    CalibrationLength float32 `json:"calibrationLength"`
}

func LoadFromFile(path string) (Palette, error) {
    palette := Palette{}
    bytes, err := ioutil.ReadFile(path)
    if err != nil {
        return palette, err
    }
    if err := json.Unmarshal(bytes, &palette); err != nil {
        return palette, err
    }
    return palette, nil
}

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
