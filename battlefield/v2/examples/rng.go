package examples

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
