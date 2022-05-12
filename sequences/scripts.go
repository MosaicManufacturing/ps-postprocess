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
	Start          printerscript.Tree
	End            printerscript.Tree
	LayerChange    printerscript.Tree
	MaterialChange []printerscript.Tree
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

func (s *Scripts) Parse() (ParsedScripts, error) {
	parsed := ParsedScripts{
		Start:          nil,
		End:            nil,
		LayerChange:    nil,
		MaterialChange: make([]printerscript.Tree, len(s.MaterialChange)),
	}
	if len(strings.TrimSpace(s.Start)) > 0 {
		tree, err := printerscript.LexAndParse(s.Start)
		if err != nil {
			return parsed, err
		}
		parsed.Start = tree
	}
	if len(strings.TrimSpace(s.End)) > 0 {
		tree, err := printerscript.LexAndParse(s.End)
		if err != nil {
			return parsed, err
		}
		parsed.End = tree
	}
	if len(strings.TrimSpace(s.LayerChange)) > 0 {
		tree, err := printerscript.LexAndParse(s.LayerChange)
		if err != nil {
			return parsed, err
		}
		parsed.LayerChange = tree
	}
	for i, script := range s.MaterialChange {
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
