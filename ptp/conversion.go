package ptp

import (
	"errors"
	"strconv"
	"strings"
)

func parseToolColors(serialized string) ([][3]float32, error) {
	perToolVals := strings.Split(serialized, "|")
	toolColors := make([][3]float32, 0, len(perToolVals))
	for _, colors := range perToolVals {
		rgbParts := strings.Split(colors, ",")
		if len(rgbParts) != 3 {
			return nil, errors.New("expected 3 components for RGB value")
		}
		r, rErr := strconv.ParseFloat(rgbParts[0], 32)
		if rErr != nil {
			return nil, rErr
		}
		g, gErr := strconv.ParseFloat(rgbParts[1], 32)
		if gErr != nil {
			return nil, gErr
		}
		b, bErr := strconv.ParseFloat(rgbParts[2], 32)
		if bErr != nil {
			return nil, bErr
		}
		thisToolColors := [3]float32{float32(r), float32(g), float32(b)}
		toolColors = append(toolColors, thisToolColors)
	}
	return toolColors, nil
}

func convertPathType(hint string) PathType {
	switch hint {
	case "Perimeter":
		return PathTypeInnerPerimeter
	case "External perimeter":
		fallthrough
	case "Overhang perimeter":
		return PathTypeOuterPerimeter
	case "Internal infill":
		return PathTypeInfill
	case "Solid infill":
		fallthrough
	case "Top solid infill":
		return PathTypeSolidLayer
	case "Bridge infill":
		return PathTypeBridge
	case "Ironing":
		return PathTypeIroning
	case "Gap fill":
		return PathTypeGapFill
	case "Skirt":
		fallthrough
	case "Skirt/Brim":
		return PathTypeBrim
	case "Support material":
		return PathTypeSupport
	case "Support material interface":
		return PathTypeSupportInterface
	case "Wipe tower":
		return PathTypeTransition
	case "Side transition":
		return PathTypeTransition
	case "Custom":
		return PathTypeStartSequence
	}
	return PathTypeUnknown
}
