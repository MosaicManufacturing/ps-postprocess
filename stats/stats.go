package stats

import (
    "encoding/json"
    "math"
)

type BoundingBox struct {
    Min [3]float64 `json:"min"`
    Max [3]float64 `json:"max"`
}

type PrintSummary struct {
    Time int `json:"time"` // seconds
    Length []float32 `json:"length"` // mm
    Volume []float32 `json:"volume"` // mm3
    Splices int `json:"splices"`
    Pings int `json:"pings"`
    BoundingBox BoundingBox `json:"boundingBox"`
}

type jsonPrintSummary struct {
    PrintSummary
    TotalLength float32 `json:"totalLength"` // mm
    TotalVolume float32 `json:"totalVolume"` // mm3
    InputsUsed int `json:"inputsUsed"`
}

func NewPrintSummary(maxInputCount int) PrintSummary {
    return PrintSummary{
        Time:        0,
        Length:      make([]float32, maxInputCount),
        Volume:      make([]float32, maxInputCount),
        Splices:     0,
        Pings:       0,
        BoundingBox: BoundingBox{
            Min: [3]float64{math.Inf(1), math.Inf(1), math.Inf(1)},
            Max: [3]float64{math.Inf(-1), math.Inf(-1), math.Inf(-1)},
        },
    }
}

func (ps *PrintSummary) ToJSON() (string, error) {
    jps := jsonPrintSummary{
        PrintSummary: *ps,
        TotalLength:  0,
        TotalVolume:  0,
        InputsUsed:   0,
    }
    for i := 0; i < len(ps.Length); i++ {
        jps.TotalLength += ps.Length[i]
        jps.TotalVolume += ps.Volume[i]
        if ps.Length[i] > 0 {
            jps.InputsUsed++
        }
    }
    bytes, err := json.Marshal(jps)
    if err != nil {
        return "", err
    }
    return string(bytes), nil
}
