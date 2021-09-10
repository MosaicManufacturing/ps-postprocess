package msf

import (
    "encoding/binary"
    "fmt"
    "math"
    "strings"
)

func intToHexString(value uint, minHexDigits int) string {
    return fmt.Sprintf("%0*x", minHexDigits, value)
}

func int16ToHexString(value int16) string {
    return intToHexString(uint(uint16(value)), 4)
}

func floatToHexString(value float32) string {
    bits := math.Float32bits(value)
    buf := make([]byte, 4)
    binary.BigEndian.PutUint32(buf, bits)
    asUint := uint(binary.BigEndian.Uint32(buf))
    return intToHexString(asUint, 8)
}

func replaceSpaces(input string) string {
    return strings.ReplaceAll(input, " ", "_")
}

func truncate(input string, length int) string {
    if len(input) <= length {
        return input
    }
    return input[:length]
}

func msfVersionToO21(major, minor uint) string {
    versionNumber := (major * 10) + minor
    return fmt.Sprintf("O21 D%s%s", intToHexString(versionNumber, 4), EOL)
}

func getLineLength(x1, y1, x2, y2 float32) float32 {
    dx := float64(x2 - x1)
    dy := float64(y2 - y1)
    return float32(math.Sqrt(dx * dx + dy * dy))
}

func lerp(minVal, maxVal, t float32) float32 {
    boundedT := float32(math.Max(0, math.Min(1, float64(t))))
    return ((1 - boundedT) * minVal) + (t * maxVal)
}
