package battlefield

const (
	Left  Side = false
	Right      = true
)

const (
	Damage Axis = iota
	CriticalOdds
	CriticalLoss
	Defense
	Loss
	Health
	HealthPercent
	HealthMaximum
	Position
	Speed
)

type Axis uint8
type Side bool

type Ratio struct {
	Current int
	Maximum int
}

type Warrior interface {
	Portfolio
	Side() Side
	Health() Ratio
	SetHealth(Ratio)
	Component(Axis) int
}

type ByAxis struct {
	Axis
	Asc      bool
	Warriors []Warrior
}

func (a *ByAxis) Len() int {
	return len(a.Warriors)
}

func (a *ByAxis) Swap(i, j int) {
	a.Warriors[i], a.Warriors[j] = a.Warriors[j], a.Warriors[i]
}

func (a *ByAxis) Less(i, j int) bool {
	if a.Warriors[i].Component(a.Axis) != a.Warriors[j].Component(a.Axis) {
		return a.Warriors[i].Component(a.Axis) < a.Warriors[j].Component(a.Axis) == a.Asc
	}
	if a.Warriors[i].Component(Position) != a.Warriors[j].Component(Position) {
		return a.Warriors[i].Component(Position) < a.Warriors[j].Component(Position)
	}

	return a.Warriors[i].Side() == Left
}
