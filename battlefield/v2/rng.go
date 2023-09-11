package battlefield

import "math/rand"

type Rng interface {
	Float64() float64
}

var DefaultRng Rng = JustRng{}

type JustRng struct{}

func (JustRng) Float64() float64 {
	return rand.Float64()
}
