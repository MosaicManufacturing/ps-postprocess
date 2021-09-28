package gcode

import (
    "fmt"
    "math"
    "regexp"
    "strconv"
)

var timeEstimateRegexp = regexp.MustCompile("estimated printing time \\(normal mode\\) = (?:(\\d+)d ?)?(?:(\\d+)h ?)?(?:(\\d+)m ?)?(?:(\\d+)s ?)")

func ParseTimeString(str string) (float32, error) {
    matches := timeEstimateRegexp.FindStringSubmatch(str)
    timeTotal := 0
    if len(matches[1]) > 0 {
        // days
        days, err := strconv.ParseInt(matches[1], 10, 32)
        if err != nil {
            return 0, err
        }
        timeTotal += int(days) * 86400
    }
    if len(matches[2]) > 0 {
        // hours
        hours, err := strconv.ParseInt(matches[2], 10, 32)
        if err != nil {
            return 0, err
        }
        timeTotal += int(hours) * 3600
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

func GetTimeString(timeEstimate float32) string {
    seconds := int(math.Round(float64(timeEstimate)))
    days := seconds / 86400
    seconds -= days * 86400
    hours := seconds / 3600
    seconds -= hours * 3600
    minutes := seconds / 60
    seconds -= minutes * 60

    if days > 0 {
        return fmt.Sprintf("%dd %dh %dm %ds", days, hours, minutes, seconds)
    }
    if hours > 0 {
        return fmt.Sprintf("%dh %dm %ds", hours, minutes, seconds)
    }
    if minutes > 0 {
        return fmt.Sprintf("%dm %ds", minutes, seconds)
    }
    return fmt.Sprintf("%ds", seconds)
}
