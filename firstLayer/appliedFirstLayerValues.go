package firstLayer

import (
	"encoding/json"
	"io/ioutil"
)

type FirstLayer struct {
	ZOffset        float32 `json:"zOffset"`
	BedTemperature float32 `json:"bedTemperature"`
}

func (p *FirstLayer) Save(path string) error {
	asJson, err := json.Marshal(p)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(path, asJson, 0644)
}
