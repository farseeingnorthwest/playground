package battlefield

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBattlefield_Fight(t *testing.T) {
	ob := &mockObserver{}
	b := NewBattleField(
		[]Warrior{
			&mockFighter{
				FatPortfolio: FatPortfolio{
					[]Reactor{
						&NormalAttack{&RandomSelector{}, 10},
						&Critical{
							rng:    &mockRng{.001},
							odds:   10,
							damage: &TemporaryDamage{200, false},
						},
					},
				},
				element: Water,
				defense: 5,
				health:  20,
				speed:   10,
			},
		},
		[]Warrior{
			&mockFighter{
				FatPortfolio: FatPortfolio{[]Reactor{
					&NormalAttack{&RandomSelector{}, 15},
				}},
				element: Fire,
				defense: 5,
				health:  22,
				speed:   9,
			},
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

type mockFighter struct {
	FatPortfolio

	element Element
	defense int
	health  int
	speed   int
}

func (f *mockFighter) Element() Element {
	return f.element
}

func (f *mockFighter) Defense() int {
	return f.defense
}

func (f *mockFighter) Health() int {
	return f.health
}

func (f *mockFighter) SetHealth(health int) {
	f.health = health
}

func (f *mockFighter) Speed() int {
	return f.speed
}

func (f *mockFighter) Functional() bool {
	return f.health > 0
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

	id(script.Subject.(*Fighter))
	p(" -> {[")
	for j, object := range script.Objects {
		comma(j)
		id(object.(*Fighter))
	}
	p("] / %d}", script.Verb.(*Attack).Points)

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
