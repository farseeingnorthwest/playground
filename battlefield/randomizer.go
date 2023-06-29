package battlefield

import "math/rand"

type Randomizer interface {
	Float64() float64
}

type DefaultRandomizer struct{}

func (DefaultRandomizer) Float64() float64 {
	return rand.Float64()
}
