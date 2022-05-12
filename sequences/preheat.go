package sequences

import (
	"encoding/json"
	"io/ioutil"
)

type PreheatHints struct {
	Extruder float32 `json:"extruder"`
	Bed      float32 `json:"bed"`
	Chamber  float32 `json:"chamber"`
}

func (p *PreheatHints) Save(path string) error {
	asJson, err := json.Marshal(p)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(path, asJson, 0644)
}
