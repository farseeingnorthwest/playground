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
				Fire:    TemporaryDamage{120, false},
				Thunder: TemporaryDamage{80, false},
			},
			Fire: {
				Ice:   TemporaryDamage{120, false},
				Water: TemporaryDamage{80, false},
			},
			Ice: {
				Wind: TemporaryDamage{120, false},
				Fire: TemporaryDamage{80, false},
			},
			Wind: {
				Earth: TemporaryDamage{120, false},
				Ice:   TemporaryDamage{80, false},
			},
			Earth: {
				Thunder: TemporaryDamage{120, false},
				Wind:    TemporaryDamage{80, false},
			},
			Thunder: {
				Water: TemporaryDamage{120, false},
				Earth: TemporaryDamage{80, false},
			},
			Dark: {
				Light: TemporaryDamage{120, false},
			},
			Light: {
				Dark: TemporaryDamage{120, false},
			},
		},
	}
)

type ElementTheory struct {
	theory map[Element]map[Element]TemporaryDamage
}

func (t *ElementTheory) React(signal Signal) {
	prepare, ok := signal.(*PreActionSignal)
	if !ok {
		return
	}
	_, ok = prepare.Verb.(*Attack)
	if !ok {
		return
	}

	theory := t.theory[prepare.Subject.Element()]
	for _, object := range prepare.Objects {
		if damage, ok := theory[object.Element()]; ok {
			prepare.Add(&Action{
				Subject: prepare.Subject,
				Objects: []Warrior{object},
				Verb:    &Buffing{Buff: damage.Fork()},
			})
		}
	}
}

func (t *ElementTheory) Valid() bool {
	return true
}
