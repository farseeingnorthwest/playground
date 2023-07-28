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
	Theory = ElementTheory{
		theory: map[Element]map[Element]TemporaryDamage{
			Water: {
				Fire:    TemporaryDamage{120, nil},
				Thunder: TemporaryDamage{80, nil},
			},
			Fire: {
				Ice:   TemporaryDamage{120, nil},
				Water: TemporaryDamage{80, nil},
			},
			Ice: {
				Wind: TemporaryDamage{120, nil},
				Fire: TemporaryDamage{80, nil},
			},
			Wind: {
				Earth: TemporaryDamage{120, nil},
				Ice:   TemporaryDamage{80, nil},
			},
			Earth: {
				Thunder: TemporaryDamage{120, nil},
				Wind:    TemporaryDamage{80, nil},
			},
			Thunder: {
				Water: TemporaryDamage{120, nil},
				Earth: TemporaryDamage{80, nil},
			},
			Dark: {
				Light: TemporaryDamage{120, nil},
			},
			Light: {
				Dark: TemporaryDamage{120, nil},
			},
		},
	}
)

type Element uint8

type ElementTheory struct {
	theory map[Element]map[Element]TemporaryDamage
}

func (t *ElementTheory) React(signal Signal) {
	sig, ok := signal.(*PreActionSignal)
	if !ok {
		return
	}
	_, ok = sig.Verb.(*Hit)
	if !ok {
		return
	}

	theory := t.theory[sig.Source.Element()]
	for _, object := range sig.Targets {
		if damage, ok := theory[object.Element()]; ok {
			sig.Add(&Action{
				Source:  sig.Source,
				Targets: []*Fighter{object},
				Verb:    NewBuffing(damage.Fork(sig.Action)),
			})
		}
	}
}
