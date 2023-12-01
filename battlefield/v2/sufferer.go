package battlefield

import "math/big"

type Sufferer interface {
	Suffer(attack, defense int) int
}

type ConstSufferer struct{}

func (s ConstSufferer) Suffer(attack, defense int) int {
	suffer := attack - defense
	if suffer < 0 {
		return 0
	}

	return suffer
}

type ReciprocalSufferer struct {
	Scale int
}

// Suffer returns the loss of the defender.
//
//	attack * attack / (defense + scale * defense)
func (s ReciprocalSufferer) Suffer(attack, defense int) int {
	a := big.NewInt(int64(attack))
	suffer := new(big.Int).Div(
		new(big.Int).Mul(a, a),
		new(big.Int).Add(a, new(big.Int).Mul(
			big.NewInt(int64(s.Scale)),
			big.NewInt(int64(defense)),
		)),
	)

	return int(suffer.Int64())
}
