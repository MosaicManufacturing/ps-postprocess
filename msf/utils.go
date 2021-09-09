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
