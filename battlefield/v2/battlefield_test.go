package battlefield

import (
	"bytes"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBattlefield_Fight(t *testing.T) {
	ob := &mockObserver{}
	Fight(
		[]Warrior{
			&mockFighter{
				NormalAttack: NormalAttack{&RandomSelector{}, 10},
				speed:        10,
				health:       20,
			},
		},
		[]Warrior{
			&mockFighter{
				NormalAttack: NormalAttack{&RandomSelector{}, 9},
				speed:        9,
				health:       22,
			},
		},
		ob,
	)

	assert.Equal(
		t,
		[]string{
			"L#0 -> {[R#0] / 10}",
			"R#0 -> {[L#0] / 9}",
			"L#0 -> {[R#0] / 10}",
			"R#0 -> {[L#0] / 9}",
			"L#0 -> {[R#0] / 10}",
		},
		ob.scripts,
	)
}

type mockFighter struct {
	NormalAttack

	speed  int
	health int
}

func (f *mockFighter) Speed() int {
	return f.speed
}

func (f *mockFighter) Functional() bool {
	return f.health > 0
}

func (f *mockFighter) Render(attack Attack) {
	f.health -= attack.Points
}

type mockObserver struct {
	scripts []string
}

func (o *mockObserver) Observe(script *Script) {
	o.scripts = append(o.scripts, tr(script))
}

func tr(script *Script) string {
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

	id(script.Subject)
	p(" -> {")
	for i, attack := range script.Attacks {
		comma(i)
		p("[")
		for j, object := range attack.Objects {
			comma(j)
			id(object)
		}
		p("] / %d", attack.Points)
	}
	p("}")

	return b.String()
}
