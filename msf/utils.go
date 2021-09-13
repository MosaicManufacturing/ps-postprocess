package msf

import (
    "bufio"
    "encoding/binary"
    "fmt"
    "io"
    "io/ioutil"
    "math"
    "os"
    "path"
    "regexp"
    "strconv"
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

func writeLine(writer *bufio.Writer, line string) error {
    _, err := writer.WriteString(line + EOL)
    return err
}

func writeLines(writer *bufio.Writer, lines string) error {
    _, err := writer.WriteString(lines)
    return err
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

func parseTimeEstimate(str string) (float32, error) {
    r, err := regexp.Compile("estimated printing time \\(normal mode\\) = (?:(\\d+)d ?)?(?:(\\d+)h ?)?(?:(\\d+)m ?)?(?:(\\d+)s ?)")
    if err != nil {
        return 0, err
    }
    matches := r.FindStringSubmatch(str)
    timeTotal := 0
    if len(matches[1]) > 0 {
        // days
        days, err := strconv.ParseInt(matches[1], 10, 32)
        if err != nil {
            return 0, err
        }
        timeTotal += int(days) * 60 * 60 * 24
    }
    if len(matches[2]) > 0 {
        // hours
        hours, err := strconv.ParseInt(matches[2], 10, 32)
        if err != nil {
            return 0, err
        }
        timeTotal += int(hours) * 60 * 60
    }
    if len(matches[3]) > 0 {
        // minutes
        minutes, err := strconv.ParseInt(matches[3], 10, 32)
        if err != nil {
            return 0, err
        }
        timeTotal += int(minutes) * 60
    }
    if len(matches[4]) > 0 {
        // seconds
        seconds, err := strconv.ParseInt(matches[4], 10, 32)
        if err != nil {
            return 0, err
        }
        timeTotal += int(seconds)
    }
    return float32(timeTotal), nil
}

func prependFile(filepath, content string) error {
    // create a temporary file
    tempfile, err := ioutil.TempFile(path.Dir(filepath), "")
    if err != nil {
        return err
    }

    // write prepended content first
    if _, err := tempfile.WriteString(content); err != nil {
        return err
    }

    // now append content of original file
    reader, err := os.Open(filepath)
    if err != nil {
        return err
    }
    if _, err := io.Copy(tempfile, reader); err != nil {
        return err
    }
    if err := reader.Close(); err != nil {
        return err
    }

    // finalize and close temporary file
    if err := tempfile.Sync(); err != nil {
        return err
    }
    if err := tempfile.Close(); err != nil {
        return err
    }

    // overwrite original file with temporary one
    return os.Rename(tempfile.Name(), filepath)
}
