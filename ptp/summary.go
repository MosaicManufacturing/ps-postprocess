package ptp

import (
	"encoding/json"
	"io/ioutil"
)

type Point struct {
	X float32
	Y float32
	Z float32
}
type BoundingBox struct {
	Min Point
	Max Point
}
type Summary struct {
	BoundingBox BoundingBox
}

func (p *Summary) Save(path string) error {
	// create a temporary structure for JSON marshaling
	type TempBoundingBox struct {
		Min []float32 `json:"min"`
		Max []float32 `json:"max"`
	}
	type TempSummary struct {
		BoundingBox TempBoundingBox `json:"boundingBox"`
	}

	// convert Points to slices
	tempSummary := TempSummary{
		BoundingBox: TempBoundingBox{
			Min: []float32{p.BoundingBox.Min.X, p.BoundingBox.Min.Y, p.BoundingBox.Min.Z},
			Max: []float32{p.BoundingBox.Max.X, p.BoundingBox.Max.Y, p.BoundingBox.Max.Z},
		},
	}

	// marshal the temporary structure to JSON
	asJson, err := json.Marshal(tempSummary)
	if err != nil {
		return err
	}

	// write the JSON to a file
	return ioutil.WriteFile(path, asJson, 0644)
}
