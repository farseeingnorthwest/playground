package battlefield

import "math/rand"

type JustRng struct{}

func (JustRng) Float64() float64 {
	return rand.Float64()
}
