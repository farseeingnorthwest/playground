package battlefield

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewWarrior(t *testing.T) {
	warrior := NewWarrior(
		&Baseline{20, .125, 10, 100},
		&Carrier{1.1, 1, 1.1, 1.5, 9},
		&Magnitude{.1, 1, .1, .5},
	)

	assert.Equal(t, 24, warrior.attack)
	assert.Equal(t, .25, warrior.critical)
	assert.Equal(t, 12, warrior.defense)
	assert.Equal(t, 200, warrior.health)
	assert.Equal(t, 9, warrior.velocity)
	assert.Equal(t, 0, warrior.buffers[Attack].Len())
	assert.Equal(t, 0, warrior.buffers[Critical].Len())
	assert.Equal(t, 0, warrior.buffers[Defense].Len())
	assert.Equal(t, 0, warrior.buffers[Health].Len())
	assert.Equal(t, 1, warrior.buffers[HealthCritical].Len())
	assert.Equal(t, 0, warrior.buffers[Velocity].Len())
	assert.Equal(t, false, warrior.criticalAttack)
}

func TestWarrior_Attach(t *testing.T) {
	for i, tt := range []struct {
		a Attribute
		n int
	}{
		{Health, 1},
		{HealthCritical, 2},
	} {
		t.Run(fmt.Sprintf("#%v", i), func(t *testing.T) {
			warrior := NewWarrior(&Baseline{}, &Carrier{}, &Magnitude{})
			buffer := &volatileBuffer{}

			warrior.Attach(tt.a, buffer)

			assert.Equal(t, tt.n, warrior.buffers[tt.a].Len())
			p := (*bufferNode)(warrior.buffers[tt.a])
			for i := 0; i < tt.n; i++ {
				p = p.next
			}
			assert.Equal(t, buffer, p.Buffer)
		})
	}
}

func TestWarrior_Prepare(t *testing.T) {
	warrior := NewWarrior(&Baseline{}, &Carrier{}, &Magnitude{})
	warrior.Attach(Critical, &volatileBuffer{r: 2, v: .5})
	warrior.Attach(HealthCritical, &volatileBuffer{})

	warrior.Prepare(&fixedRandomizer{.55})
	assert.Equal(t, 1, warrior.buffers[Critical].Len())
	assert.Equal(t, 1, warrior.buffers[HealthCritical].Len())
	assert.Equal(t, false, warrior.criticalAttack)

	warrior.Prepare(&fixedRandomizer{.45})
	assert.Equal(t, 1, warrior.buffers[Critical].Len())
	assert.Equal(t, true, warrior.criticalAttack)

	warrior.Prepare(&fixedRandomizer{.05})
	assert.Equal(t, 0, warrior.buffers[Critical].Len())
	assert.Equal(t, false, warrior.criticalAttack)
}

func TestWarrior_Attack(t *testing.T) {
	for _, tt := range []struct {
		title    string
		attack   int
		critical bool
	}{
		{"normal", 5, false},
		{"critical", 10, true},
	} {
		t.Run(tt.title, func(t *testing.T) {
			warrior := NewWarrior(&Baseline{}, &Carrier{}, &Magnitude{})
			warrior.Attach(Attack, &volatileBuffer{r: 2, v: float64(tt.attack)})
			warrior.criticalAttack = tt.critical

			attack, critical := warrior.Attack()

			assert.Equal(t, tt.attack, attack)
			assert.Equal(t, tt.critical, critical)
		})
	}
}

func TestWarrior_Velocity(t *testing.T) {
	warrior := NewWarrior(&Baseline{}, &Carrier{}, &Magnitude{})
	warrior.Attach(Velocity, &volatileBuffer{r: 2, v: 10})

	assert.Equal(t, 10, warrior.Velocity())
}

func TestWarrior_Suffer(t *testing.T) {
	for _, tt := range []struct {
		title    string
		attack   int
		critical bool
		health   int
		damage   int
		overflow int
	}{
		{"normal", 5, false, 100, 6, 0},
		{"critical", 10, true, 100, 12, 0},
		{"normal overflow", 5, false, 2, 6, 4},
		{"critical overflow", 10, true, 10, 12, 2},
		{"cannot break the defense", 2, false, 100, 0, 0},
		{"critical cannot break the defense", 2, true, 100, 0, 0},
	} {
		t.Run(tt.title, func(t *testing.T) {
			warrior := NewWarrior(&Baseline{}, &Carrier{}, &Magnitude{})
			warrior.Attach(Defense, &volatileBuffer{r: 2, v: 2})
			warrior.Attach(Health, &buffer{2})
			warrior.Attach(HealthCritical, &buffer{.5})

			warrior.health = tt.health
			damage, overflow := warrior.Suffer(tt.attack, tt.critical)

			assert.Equal(t, tt.damage, damage)
			assert.Equal(t, tt.overflow, overflow)
		})
	}
}

type buffer struct {
	f float64
}

func (b *buffer) Drain() int {
	return 1
}

func (b *buffer) Buff(v float64) float64 {
	return v * b.f
}

type fixedRandomizer struct {
	value float64
}

func (c *fixedRandomizer) Float64() float64 {
	return c.value
}
