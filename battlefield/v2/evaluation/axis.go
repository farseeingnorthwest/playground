package evaluation

const (
	Damage Axis = iota
	Defense
	Loss
	Healing
	Health
	HealthMax
	Speed
)

type Axis uint8
