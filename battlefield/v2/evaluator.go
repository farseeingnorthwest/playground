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

type EvaluationContext interface {
	Warriors() []Warrior
	React(ForkableSignal)
}

type Evaluator interface {
	Evaluate(Warrior, EvaluationContext) int
}

type ConstEvaluator int

func (e ConstEvaluator) Evaluate(Warrior, EvaluationContext) int {
	return int(e)
}

type AxisEvaluator Axis

func (e AxisEvaluator) Evaluate(warrior Warrior, ec EvaluationContext) int {
	return warrior.Component(Axis(e), ec)
}

type BuffCounter struct {
	tag any
}

func (e BuffCounter) Evaluate(warrior Warrior, _ EvaluationContext) int {
	return len(warrior.Buffs(e.tag))
}

type Adder struct {
	adder     int
	evaluator Evaluator
}

func NewAdder(adder int, evaluator Evaluator) *Adder {
	return &Adder{adder, evaluator}
}

func (e *Adder) Evaluate(warrior Warrior, ec EvaluationContext) int {
	return e.adder + e.evaluator.Evaluate(warrior, ec)
}

type Multiplier struct {
	multiplier int
	evaluator  Evaluator
}

func NewMultiplier(multiplier int, evaluator Evaluator) *Multiplier {
	return &Multiplier{multiplier, evaluator}
}

func (e *Multiplier) Evaluate(warrior Warrior, ec EvaluationContext) int {
	return e.multiplier * e.evaluator.Evaluate(warrior, ec) / 100
}

type SelectCounter struct {
	selectors []Selector
}

func NewSelectCounter(selectors ...Selector) *SelectCounter {
	return &SelectCounter{selectors}
}

func (e *SelectCounter) Evaluate(warrior Warrior, ec EvaluationContext) int {
	signal := NewFreeSignal(warrior)
	warriors := ec.Warriors()
	for _, selector := range e.selectors {
		warriors = selector.Select(warriors, signal, ec)
	}

	return len(warriors)
}
