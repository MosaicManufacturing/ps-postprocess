package sequences

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
)

func LoadLocals(jsonPath string) (map[string]float64, error) {
    locals := make(map[string]float64)
    data, err := ioutil.ReadFile(jsonPath)
    if err != nil {
        return locals, err
    }
    fmt.Println(string(data))
    err = json.Unmarshal(data, &locals)
    return locals, err
}

func MergeLocals(a, b map[string]float64) map[string]float64 {
    locals := make(map[string]float64)
    for k, v := range a {
        locals[k] = v
    }
    // do b second, allowing it to overwrite values from a
    for k, v := range b {
        locals[k] = v
    }
    return locals
}
