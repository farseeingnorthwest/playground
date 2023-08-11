package battlefield_test

import "github.com/farseeingnorthwest/playground/battlefield/v2"

type RngProxy struct {
	rng battlefield.Rng
}

func (p *RngProxy) SetRng(rng battlefield.Rng) {
	p.rng = rng
}

func (p *RngProxy) Float64() float64 {
	return p.rng.Float64()
}

type Sequence struct {
	floats []float64
}

func NewSequence(floats ...float64) *Sequence {
	return &Sequence{floats}
}

func (s *Sequence) Float64() float64 {
	f := s.floats[0]
	if len(s.floats) > 1 {
		s.floats = s.floats[1:]
	}

	return f
}
