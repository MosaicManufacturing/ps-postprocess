package sequences

import (
    "regexp"
    "strconv"
)

const startPlaceholder = ";*/*/*/*/* START SEQUENCE */*/*/*/*"
const endPlaceholder = ";*/*/*/*/* END SEQUENCE */*/*/*/*"
const layerChangePrefix = ";*/*/*/*/* LAYER CHANGE SEQUENCE ("
const materialChangePrefix = ";*/*/*/*/* MATERIAL CHANGE SEQUENCE ("

var layerChangeRegexp = regexp.MustCompile(";\\*/\\*/\\*/\\*/\\* LAYER CHANGE SEQUENCE \\((\\d+), (\\d+(?:\\.\\d+)?)\\) \\*/\\*/\\*/\\*/\\*")
var materialChangeRegexp = regexp.MustCompile(";\\*/\\*/\\*/\\*/\\* MATERIAL CHANGE SEQUENCE \\((\\d+)\\) \\*/\\*/\\*/\\*/\\*")

func parseLayerChangePlaceholder(placeholder string) (layer int, layerZ float64, err error) {
    matches := layerChangeRegexp.FindStringSubmatch(placeholder)
    if len(matches[1]) > 0 {
        // layer number
        parsed, err := strconv.ParseInt(matches[1], 10, 32)
        if err != nil {
            return 0, 0, err
        }
        layer = int(parsed)
    }
    if len(matches[2]) > 0 {
        // layer Z
        layerZ, err = strconv.ParseFloat(matches[2], 64)
    }
    return
}

func parseMaterialChangePlaceholder(placeholder string) (toTool int, err error) {
    if err != nil {
        return 0, err
    }
    matches := materialChangeRegexp.FindStringSubmatch(placeholder)
    if len(matches[1]) > 0 {
        // layer number
        parsed, err := strconv.ParseInt(matches[1], 10, 32)
        if err != nil {
            return 0, err
        }
        toTool = int(parsed)
    }
    return
}
