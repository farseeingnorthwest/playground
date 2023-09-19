package examples

import (
	"encoding/json"
	"errors"
)

const (
	Water Element = iota
	Fire
	Ice
	Wind
	Earth
	Thunder
	Dark
	Light
)

var (
	ErrBadElement = errors.New("bad element")
	ElementNames  = []string{"Water", "Fire", "Ice", "Wind", "Earth", "Thunder", "Dark", "Light"}
)

type Element uint8

func (e Element) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]string{
		"_kind": "element",
		"name":  ElementNames[e],
	})
}

func (e *Element) UnmarshalJSON(data []byte) error {
	var v struct{ Name string }
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	for i, name := range ElementNames {
		if v.Name == name {
			*e = Element(i)
			return nil
		}
	}

	return ErrBadElement
}
