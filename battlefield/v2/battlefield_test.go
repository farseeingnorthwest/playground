package battlefield

import (
	"fmt"
	"github.com/farseeingnorthwest/playground/battlefield/v2/mod"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBattlefield_Fight(t *testing.T) {
	ob := &observer{}
	ob.SetTag("Ob")
	b := NewBattleField(
		[]*Warrior{
			NewWarrior(
				&baseline{
					TaggerMod: mod.NewTaggerMod("Alice"),
					element:   Water,
					damage:    10,
					defense:   5,
					speed:     10,
					health:    20,
				},
				&FatPortfolio{
					[]Reactor{
						NormalAttack,
						NewCriticalAttack(&mockRng{0.001}, 10, 200),
					},
				},
			),
		},
		[]*Warrior{
			NewWarrior(
				&baseline{
					TaggerMod: mod.NewTaggerMod("Bob"),
					element:   Fire,
					damage:    15,
					defense:   5,
					speed:     9,
					health:    22,
				},
				&FatPortfolio{[]Reactor{
					NormalAttack,
				}},
			),
		},
		ob,
		&Theory,
	)
	b.Fight()

	assert.Equal(
		t,
		[]string{
			"<Buff> Alice -> [Bob]",
			"<Buff> Alice -> [Bob]",
			"<Attack> Alice -> [Bob]",
			"<Buff> Bob -> [Alice]",
			"<Attack> Bob -> [Alice]",
			"<Buff> Alice -> [Bob]",
			"<Attack> Alice -> [Bob]",
			"<Buff> Bob -> [Alice]",
			"<Attack> Bob -> [Alice]",
			"<Buff> Alice -> [Bob]",
			"<Attack> Alice -> [Bob]",
		},
		ob.scripts,
	)
}

type baseline struct {
	mod.TaggerMod
	element Element
	damage  int
	defense int
	speed   int
	health  int
}

func (f *baseline) Element() Element {
	return f.element
}

func (f *baseline) Damage() int {
	return f.damage
}

func (f *baseline) Defense() int {
	return f.defense
}

func (f *baseline) Health() int {
	return f.health
}

func (f *baseline) Speed() int {
	return f.speed
}

type observer struct {
	mod.TaggerMod
	scripts []string
}

func (ob *observer) React(s Signal) {
	switch sig := s.(type) {
	case *PostActionSignal:
		ob.scripts = append(ob.scripts, fmt.Sprintf(
			"<%s> %v -> [%v]",
			reflect.TypeOf(sig.Verb).Elem().Name(),
			sig.Source.Baseline.(mod.Tagger).Tag(),
			sig.Targets[0].Baseline.(mod.Tagger).Tag(),
		))
	}
}

type mockRng struct {
	initial float64
}

func (r *mockRng) Gen() (f float64) {
	f = 1 - 1e-6
	if 0 <= r.initial && r.initial < 1 {
		f = r.initial
	}

	r.initial = -1
	return
}
