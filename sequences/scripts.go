package sequences

import (
    "../printerscript"
    "encoding/json"
    "io/ioutil"
    "strings"
)

type Scripts struct {
    Start string `json:"start"`
    End string `json:"end"`
    LayerChange string `json:"layerChange"`
    MaterialChange []string `json:"materialChange"`
    // todo: before side transitioning sequence
    // todo: side transition sequence
    // todo: after side transitioning sequence
}

type ParsedScripts struct {
    Start printerscript.ISequenceContext
    End printerscript.ISequenceContext
    LayerChange printerscript.ISequenceContext
    MaterialChange []printerscript.ISequenceContext
}

func trimDirective(script string) string {
    if strings.HasPrefix(script, "@printerscript 1.0") {
        return script[strings.IndexRune(script, '\n'):]
    }
    return script
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
        MaterialChange: make([]printerscript.ISequenceContext, 0, len(s.MaterialChange)),
    }
    if len(s.Start) > 0 {
        tree, err := printerscript.LexAndParse(trimDirective(s.Start))
        if err != nil {
            return parsed, err
        }
        parsed.Start = tree
    }
    if len(s.End) > 0 {
        tree, err := printerscript.LexAndParse(trimDirective(s.End))
        if err != nil {
            return parsed, err
        }
        parsed.End = tree
    }
    if len(s.LayerChange) > 0 {
        tree, err := printerscript.LexAndParse(trimDirective(s.LayerChange))
        if err != nil {
            return parsed, err
        }
        parsed.LayerChange = tree
    }
    for i, script := range s.MaterialChange {
        if len(script) > 0 {
            tree, err := printerscript.LexAndParse(trimDirective(script))
            if err != nil {
                return parsed, err
            }
            parsed.MaterialChange[i] = tree
        }
    }
    return parsed, nil
}
