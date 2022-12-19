package sequences

import (
	"encoding/json"
	"io/ioutil"
)

type PreheatHints struct {
	Extruder    float32 `json:"extruder"`    // first extruder temperature used in the print
	ExtruderMax float32 `json:"extruderMax"` // highest extruder temperature used in the print
	Bed         float32 `json:"bed"`         // first bed temperature used in the print
	Chamber     float32 `json:"chamber"`     // first chamber temperature used in the print
}

func (p *PreheatHints) Save(path string) error {
	asJson, err := json.Marshal(p)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(path, asJson, 0644)
}
