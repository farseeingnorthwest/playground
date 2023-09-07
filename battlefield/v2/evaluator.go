package battlefield

import "encoding/json"

var (
	_ Evaluator = ConstEvaluator(0)
	_ Evaluator = AxisEvaluator(0)
	_ Evaluator = BuffCounter{}
	_ Evaluator = LossEvaluator{}
	_ Evaluator = SelectCounter(nil)
	_ Evaluator = (*Adder)(nil)
	_ Evaluator = (*Multiplier)(nil)
	_ Evaluator = (*CustomEvaluator)(nil)
)

type Evaluator interface {
	Evaluate(Warrior, EvaluationContext) int
}

type EvaluationContext interface {
	Warriors() []Warrior
	React(RegularSignal)
}

type ConstEvaluator int

func (e ConstEvaluator) Evaluate(Warrior, EvaluationContext) int {
	return int(e)
}

func (e ConstEvaluator) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]int{
		"const": int(e),
	})
}

type AxisEvaluator Axis

func (e AxisEvaluator) Evaluate(warrior Warrior, ec EvaluationContext) int {
	return warrior.Component(Axis(e), ec)
}

func (e AxisEvaluator) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]string{"axis": Axis(e).String()})
}

type BuffCounter struct {
	tag any
}

func NewBuffCounter(tag any) BuffCounter {
	return BuffCounter{tag}
}

func (e BuffCounter) Evaluate(warrior Warrior, _ EvaluationContext) int {
	return len(warrior.Buffs(e.tag))
}

func (e BuffCounter) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]any{
		"count_buffs": e.tag,
	})
}

type LossEvaluator struct {
}

func (LossEvaluator) Evaluate(_ Warrior, ec EvaluationContext) int {
	return ec.(ActionContext).Action().Script().Loss()
}

func (LossEvaluator) MarshalJSON() ([]byte, error) {
	return json.Marshal("loss")
}

type SelectCounter []Selector

func NewSelectCounter(selectors ...Selector) SelectCounter {
	return selectors
}

func (e SelectCounter) Evaluate(warrior Warrior, ec EvaluationContext) int {
	signal := NewFreeSignal(warrior)
	warriors := ec.Warriors()
	for _, selector := range e {
		warriors = selector.Select(warriors, signal, ec)
	}

	return len(warriors)
}

func (e SelectCounter) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]any{
		"count": ([]Selector)(e),
	})
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

func (e *Adder) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]any{
		"add":       e.adder,
		"evaluator": e.evaluator,
	})
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

func (e *Multiplier) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]any{
		"mul":       e.multiplier,
		"evaluator": e.evaluator,
	})
}

type CustomEvaluator struct {
	evaluator func(Warrior, EvaluationContext) int
}

func NewCustomEvaluator(evaluator func(Warrior, EvaluationContext) int) *CustomEvaluator {
	return &CustomEvaluator{evaluator}
}

func (e *CustomEvaluator) Evaluate(warrior Warrior, ec EvaluationContext) int {
	return e.evaluator(warrior, ec)
}
