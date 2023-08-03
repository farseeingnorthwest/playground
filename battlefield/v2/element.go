package battlefield

import (
	"github.com/farseeingnorthwest/playground/battlefield/v2/evaluation"
	"github.com/farseeingnorthwest/playground/battlefield/v2/mod"
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
	up     = NewBuffProto(NewClearingBuff("元素提高伤害", evaluation.Loss, nil, ClearingMultiplier(120)), nil)
	down   = NewBuffProto(NewClearingBuff("元素降低伤害", evaluation.Loss, nil, ClearingMultiplier(80)), nil)
	Theory = ElementTheory{
		TaggerMod: mod.NewTaggerMod("元素"),
		theory: map[Element]map[Element]Verb{
			Water: {
				Fire:    up,
				Thunder: down,
			},
			Fire: {
				Ice:   up,
				Water: down,
			},
			Ice: {
				Wind: up,
				Fire: down,
			},
			Wind: {
				Earth: up,
				Ice:   down,
			},
			Earth: {
				Thunder: up,
				Wind:    down,
			},
			Thunder: {
				Water: up,
				Earth: down,
			},
			Dark: {
				Light: up,
			},
			Light: {
				Dark: up,
			},
		},
	}
)

type Element uint8

type ElementTheory struct {
	mod.TaggerMod
	theory map[Element]map[Element]Verb
}

func (t *ElementTheory) React(signal Signal) {
	sig, ok := signal.(*PreActionSignal)
	if !ok {
		return
	}
	_, ok = sig.Verb.(*Attack)
	if !ok {
		return
	}

	theory := t.theory[sig.Script.Current.Element()]
	for _, object := range sig.Targets {
		if damage, ok := theory[object.Element()]; ok {
			sig.Append(NewScript(
				nil,
				t,
				&Action{
					Targets: []*Fighter{object},
					Verb:    damage.Fork(nil, signal),
				},
			))
		}
	}
}
