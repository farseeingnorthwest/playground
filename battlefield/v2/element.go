package battlefield

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
	up     = NewClearingBuff(Loss, nil, ClearingSlope(120))
	down   = NewClearingBuff(Loss, nil, ClearingSlope(80))
	Theory = ElementTheory{
		theory: map[Element]map[Element]*ClearingBuff{
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
	theory map[Element]map[Element]*ClearingBuff
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

	theory := t.theory[sig.Source.Element()]
	for _, object := range sig.Targets {
		if damage, ok := theory[object.Element()]; ok {
			sig.Append(&Action{
				Source:  sig.Source,
				Targets: []*Fighter{object},
				Verb:    NewBuffing(damage.Fork(sig.Action)),
			})
		}
	}
}
