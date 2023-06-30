package battlefield

import "math/rand"

type Randomizer interface {
	Float64() float64
}

type DefaultRandomizer struct{}

func (DefaultRandomizer) Float64() float64 {
	return rand.Float64()
}

type Sequence struct {
	values []float64
	i      int
}

func NewSequence(values ...float64) *Sequence {
	return &Sequence{values, 0}
}

func (c *Sequence) Float64() (v float64) {
	v = c.values[c.i]
	c.i++
	if c.i >= len(c.values) {
		c.i = 0
	}

	return
}
