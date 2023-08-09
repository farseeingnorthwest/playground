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

func (s Side) String() string {
	if s == Left {
		return "Left"
	}

	return "Right"
}

type Ratio struct {
	Current int
	Maximum int
}

type Warrior interface {
	Portfolio
	Side() Side
	Position() int
	Health() Ratio
	SetHealth(Ratio)
	Component(Axis, EvaluationContext) int
}

type ByAxis struct {
	Axis
	Asc      bool
	Context  EvaluationContext
	Warriors []Warrior
}

func (a *ByAxis) Len() int {
	return len(a.Warriors)
}

func (a *ByAxis) Swap(i, j int) {
	a.Warriors[i], a.Warriors[j] = a.Warriors[j], a.Warriors[i]
}

func (a *ByAxis) Less(i, j int) bool {
	if a.Warriors[i].Component(a.Axis, a.Context) != a.Warriors[j].Component(a.Axis, a.Context) {
		return a.Warriors[i].Component(a.Axis, a.Context) < a.Warriors[j].Component(a.Axis, a.Context) == a.Asc
	}
	if a.Warriors[i].Component(Position, nil) != a.Warriors[j].Component(Position, nil) {
		return a.Warriors[i].Component(Position, nil) < a.Warriors[j].Component(Position, nil)
	}

	return a.Warriors[i].Side() == Left
}

type Baseline interface {
	Component(Axis) int
}

type MyBaseline struct {
	Damage       int
	CriticalOdds int
	CriticalLoss int
	Defense      int
	Health       int
	Speed        int
}

func (b MyBaseline) Component(axis Axis) int {
	switch axis {
	case Damage:
		return b.Damage
	case CriticalOdds:
		return b.CriticalOdds
	case CriticalLoss:
		return b.CriticalLoss
	case Defense:
		return b.Defense
	case HealthMaximum:
		return b.Health
	case Speed:
		return b.Speed

	default:
		panic("unknown axis")
	}
}

type MyWarrior struct {
	*FatPortfolio
	baseline Baseline
	side     Side
	position int
	health   Ratio
}

func NewMyWarrior(baseline Baseline, side Side, position int, reactors ...Reactor) *MyWarrior {
	portfolio := NewFatPortfolio()
	for _, r := range reactors {
		portfolio.Add(r)
	}

	return &MyWarrior{
		FatPortfolio: portfolio,
		baseline:     baseline,
		side:         side,
		position:     position,
		health: Ratio{
			baseline.Component(HealthMaximum),
			baseline.Component(HealthMaximum),
		},
	}
}

func (w *MyWarrior) Baseline() Baseline {
	return w.baseline
}

func (w *MyWarrior) Side() Side {
	return w.side
}

func (w *MyWarrior) Position() int {
	return w.position
}

func (w *MyWarrior) Health() Ratio {
	return w.health
}

func (w *MyWarrior) SetHealth(health Ratio) {
	w.health = health
}

func (w *MyWarrior) Component(axis Axis, ec EvaluationContext) int {
	switch axis {
	case Position:
		return w.position

	case Health:
		m := w.Component(HealthMaximum, ec)
		return w.health.Current * m / w.health.Maximum

	case HealthPercent:
		return w.health.Current * 100 / w.health.Maximum

	default:
		signal := NewEvaluationSignal(w, axis, w.baseline.Component(axis))
		w.React(signal, ec)
		return signal.Value()
	}
}
