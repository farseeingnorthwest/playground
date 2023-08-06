package battlefield

const (
	Lt IntComparator = iota
	Le
	Eq
	Ge
	Gt
)

type IntComparator uint8

func (c IntComparator) Compare(a, b int) bool {
	switch c {
	case Lt:
		return a < b
	case Le:
		return a <= b
	case Eq:
		return a == b
	case Ge:
		return a >= b
	case Gt:
		return a > b

	default:
		panic("bad comparator")
	}
}

type Evaluator interface {
	Evaluate(Warrior) int
}

type ConstEvaluator int

func (e ConstEvaluator) Evaluate(Warrior) int {
	return int(e)
}

type AxisEvaluator Axis

func (e AxisEvaluator) Evaluate(warrior Warrior) int {
	return warrior.Component(Axis(e))
}

type BuffCounter struct {
	tag any
}

func (e BuffCounter) Evaluate(warrior Warrior) int {
	return len(warrior.Buffs(e.tag))
}

type Multiplier struct {
	evaluator  Evaluator
	multiplier int
}

func NewMultiplier(evaluator Evaluator, multiplier int) *Multiplier {
	return &Multiplier{evaluator, multiplier}
}

func (e *Multiplier) Evaluate(warrior Warrior) int {
	return e.evaluator.Evaluate(warrior) * e.multiplier / 100
}
