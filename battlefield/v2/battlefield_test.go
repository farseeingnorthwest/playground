package battlefield

import (
	"bytes"
	"fmt"
	"reflect"
	"testing"

	"github.com/farseeingnorthwest/playground/battlefield/v2/evaluation"

	"github.com/stretchr/testify/assert"
)

func TestBattlefield_Fight(t *testing.T) {
	ob := &observer{}
	b := NewBattleField(
		[]*Warrior{
			NewWarrior(
				&baseline{
					element: Water,
					damage:  10,
					defense: 5,
					speed:   10,
					health:  20,
				},
				&FatPortfolio{
					[]Reactor{
						NormalAttack,
						&PreAttackReactor{
							NewModifiedReactor([]Actor{
								&ProbabilityActor{
									rng:  &mockRng{0.001},
									odds: 10,
									Actor: &BlindActor{
										proto: NewBuffProto(
											NewClearingBuff(evaluation.Loss, nil, ClearingMultiplier(200)),
											nil,
										),
									},
								},
							}),
						},
					},
				},
			),
		},
		[]*Warrior{
			NewWarrior(
				&baseline{
					element: Fire,
					damage:  15,
					defense: 5,
					speed:   9,
					health:  22,
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
			"<Buff> L#0 -> [R#0]",
			"<Buff> L#0 -> [R#0]",
			"<Attack> L#0 -> [R#0]",
			"<Buff> R#0 -> [L#0]",
			"<Attack> R#0 -> [L#0]",
			"<Buff> L#0 -> [R#0]",
			"<Attack> L#0 -> [R#0]",
			"<Buff> R#0 -> [L#0]",
			"<Attack> R#0 -> [L#0]",
			"<Buff> L#0 -> [R#0]",
			"<Attack> L#0 -> [R#0]",
		},
		ob.scripts,
	)
}

type baseline struct {
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
	scripts []string
}

func (ob *observer) React(s Signal) {
	switch sig := s.(type) {
	case *PostActionSignal:
		ob.scripts = append(ob.scripts, fmt.Sprintf(
			"<%s> %s -> [%s]",
			reflect.TypeOf(sig.Verb).Elem().Name(),
			positions([]*Fighter{sig.Source}).String(),
			positions(sig.Targets).String(),
		))
	}
}

type positions []*Fighter

func (p positions) String() string {
	var b bytes.Buffer

	for i, f := range p {
		if i > 0 {
			b.WriteString(", ")
		}
		b.WriteString(fmt.Sprintf("%s#%d", []string{"L", "R"}[f.Side], f.Position))
	}

	return b.String()
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
