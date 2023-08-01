package battlefield

import (
	"bytes"
	"fmt"
	"github.com/farseeingnorthwest/playground/battlefield/v2/evaluation"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBattlefield_Fight(t *testing.T) {
	ob := &mockObserver{}
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
						&ProbabilityAttackReactor{
							rng:  &mockRng{.001},
							odds: 10,
							proto: NewBuffProto(
								NewClearingBuff(evaluation.Loss, nil, ClearingSlope(200)),
								nil,
							),
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
		[]Reactor{
			&Theory,
		},
		ob,
	)
	b.Fight()

	assert.Equal(
		t,
		[]string{
			"L#0 -> {[R#0] / 10}",
			"R#0 -> {[L#0] / 15}",
			"L#0 -> {[R#0] / 10}",
			"R#0 -> {[L#0] / 15}",
			"L#0 -> {[R#0] / 10}",
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

type mockObserver struct {
	scripts []string
}

func (o *mockObserver) Observe(script *Action) {
	o.scripts = append(o.scripts, tr(script))
}

func tr(script *Action) string {
	var b bytes.Buffer

	p := func(format string, a ...any) {
		_, _ = fmt.Fprintf(&b, format, a...)
	}
	id := func(f *Fighter) {
		p("%s#%d", []string{"L", "R"}[f.Side], f.Position)
	}
	comma := func(i int) {
		if i > 0 {
			p(", ")
		}
	}

	id(script.Source)
	p(" -> {[")
	for j, object := range script.Targets {
		comma(j)
		id(object)
	}
	p("] / %d}", script.Verb.(*Attack).Block().Value())

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
