package gcode

import (
    "fmt"
    "math"
    "strconv"
    "strings"
)

const delimiterRow = "|"
const delimiterCol = ","

var posInf = float32(math.Inf(1))
var negInf = float32(math.Inf(-1))

type BoundingBox struct {
    Min [3]float32
    Max [3]float32
}

func NewBoundingBox() BoundingBox {
    return BoundingBox{
        Min: [3]float32{posInf, posInf, posInf},
        Max: [3]float32{negInf, negInf, negInf},
    }
}

func (b *BoundingBox) ExpandX(x float32) {
    if x < b.Min[0] { b.Min[0] = x }
    if x > b.Max[0] { b.Max[0] = x }
}

func (b *BoundingBox) ExpandY(y float32) {
    if y < b.Min[1] { b.Min[1] = y }
    if y > b.Max[1] { b.Max[1] = y }
}

func (b *BoundingBox) ExpandZ(z float32) {
    if z < b.Min[2] { b.Min[2] = z }
    if z > b.Max[2] { b.Max[2] = z }
}

func (b *BoundingBox) Serialize() string {
    serializedMin := fmt.Sprintf("%.6e%s%.6e%s%.6e", b.Min[0], delimiterCol, b.Min[1], delimiterCol, b.Min[2])
    serializedMax := fmt.Sprintf("%.6e%s%.6e%s%.6e", b.Max[0], delimiterCol, b.Max[1], delimiterCol, b.Max[2])
    return fmt.Sprintf("%s%s%s", serializedMin, delimiterRow, serializedMax)
}

func UnserializeBoundingBox(str string) (BoundingBox, error) {
    bbox := NewBoundingBox()
    lines := strings.Split(str, delimiterRow)
    if len(lines) != 2 {
        return bbox, fmt.Errorf("expected 2 rows in serialized BoundingBox, found %d", len(lines))
    }
    for i := 0; i < 2; i++ {
        line := strings.Split(lines[i], delimiterCol)
        if len(line) != 3 {
            return bbox, fmt.Errorf("expected 3 columns in serialized BoundingBox, found %d", len(line))
        }
        for j := 0; j < 3; j++ {
            value, err := strconv.ParseFloat(line[j], 32)
            if err != nil {
                return bbox, err
            }
            if i == 0 {
                bbox.Min[j] = float32(value)
            } else {
                bbox.Max[j] = float32(value)
            }
        }
    }
    return bbox, nil
}
