package sequences

import (
	"encoding/json"
	"io/ioutil"
	"mosaicmfg.com/ps-postprocess/printerscript"
	"strings"
)

type Scripts struct {
	Start          string   `json:"start"`
	End            string   `json:"end"`
	LayerChange    string   `json:"layerChange"`
	MaterialChange []string `json:"materialChange"`
}

type ParsedScripts struct {
	Start                        printerscript.Tree
	End                          printerscript.Tree
	LayerChange                  printerscript.Tree
	MaterialChange               []printerscript.Tree
	CoolingModuleSpeedPercentage []int
	EnableCoolingModuleAtLayer   []int
}

func LoadScripts(jsonPath string) (Scripts, error) {
	scripts := Scripts{
		MaterialChange: make([]string, 0),
	}
	data, err := ioutil.ReadFile(jsonPath)
	if err != nil {
		return scripts, err
	}
	err = json.Unmarshal(data, &scripts)
	return scripts, err
}

func (s *Scripts) Parse(locals Locals) (ParsedScripts, error) {
	parsed := ParsedScripts{
		Start:                        nil,
		End:                          nil,
		LayerChange:                  nil,
		MaterialChange:               make([]printerscript.Tree, len(s.MaterialChange)),
		CoolingModuleSpeedPercentage: make([]int, 0),
		EnableCoolingModuleAtLayer:   make([]int, 0),
	}

	coolingModuleSpeedPercentage := locals.PerExtruder["coolingModuleSpeedPercentage"]
	intValuesSpeedPercentage := make([]int, len(coolingModuleSpeedPercentage))
	for i, value := range coolingModuleSpeedPercentage {
		intValuesSpeedPercentage[i] = int(value)
	}

	enableCoolingModuleAtLayer := locals.PerExtruder["enableCoolingModuleAtLayer"]
	intValuesEnableAtLayer := make([]int, len(enableCoolingModuleAtLayer))
	for i, value := range enableCoolingModuleAtLayer {
		intValuesEnableAtLayer[i] = int(value)
	}

	parsed.CoolingModuleSpeedPercentage = intValuesSpeedPercentage
	parsed.EnableCoolingModuleAtLayer = intValuesEnableAtLayer

	s.Start = printerscript.Normalize(s.Start)
	if len(strings.TrimSpace(s.Start)) > 0 {
		tree, err := printerscript.LexAndParse(s.Start)
		if err != nil {
			return parsed, err
		}
		parsed.Start = tree
	}
	s.End = printerscript.Normalize(s.End)
	if len(strings.TrimSpace(s.End)) > 0 {
		tree, err := printerscript.LexAndParse(s.End)
		if err != nil {
			return parsed, err
		}
		parsed.End = tree
	}
	s.LayerChange = printerscript.Normalize(s.LayerChange)
	if len(strings.TrimSpace(s.LayerChange)) > 0 {
		tree, err := printerscript.LexAndParse(s.LayerChange)
		if err != nil {
			return parsed, err
		}
		parsed.LayerChange = tree
	}
	for i, script := range s.MaterialChange {
		script = printerscript.Normalize(script)
		s.MaterialChange[i] = script
		if len(strings.TrimSpace(script)) > 0 {
			tree, err := printerscript.LexAndParse(script)
			if err != nil {
				return parsed, err
			}
			parsed.MaterialChange[i] = tree
		}
	}
	return parsed, nil
}
