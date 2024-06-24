package firstLayer

import (
	"encoding/json"
	"io/ioutil"
)

type FirstLayerStyleSettings struct {
	ZOffsetPerExt  []float32 `json:"zOffsetPerExt"`
	BedTemperature []float32 `json:"bedTemperature"`
}

func LoadFirstLayerStylesFromFile(path string) (FirstLayerStyleSettings, error) {
	firstLayerStyleSettings := FirstLayerStyleSettings{}
	bytes, err := ioutil.ReadFile(path)

	if err != nil {
		return firstLayerStyleSettings, err
	}
	if err := json.Unmarshal(bytes, &firstLayerStyleSettings); err != nil {
		return firstLayerStyleSettings, err
	}
	return firstLayerStyleSettings, nil
}
