package sequences

import (
	"encoding/json"
	"io/ioutil"
)

type Locals struct {
	Global      map[string]float64
	PerExtruder map[string][]float64
}

func NewLocals() Locals {
	return Locals{
		Global:      make(map[string]float64),
		PerExtruder: make(map[string][]float64),
	}
}

func (l *Locals) LoadGlobal(jsonPath string) error {
	data, err := ioutil.ReadFile(jsonPath)
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, &l.Global)
	return err
}

func (l *Locals) LoadPerExtruder(jsonPath string) error {
	data, err := ioutil.ReadFile(jsonPath)
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, &l.PerExtruder)
	return err
}

func (l *Locals) Prepare(extruder int, extras map[string]float64) map[string]float64 {
	locals := make(map[string]float64)
	// "global" locals are constant for the entire print
	for k, v := range l.Global {
		locals[k] = v
	}
	// go from "per-extruder" locals to a concrete value
	for k, v := range l.PerExtruder {
		locals[k] = v[extruder]
	}
	// add on extra locals last
	for k, v := range extras {
		locals[k] = v
	}
	return locals
}
