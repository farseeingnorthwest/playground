package battlefield

import (
	"encoding/json"
	"errors"

	"github.com/farseeingnorthwest/playground/battlefield/v2/functional"
)

var (
	_ Evaluator = ConstEvaluator(0)
	_ Evaluator = AxisEvaluator(0)
	_ Evaluator = BuffCounter{}
	_ Evaluator = LossEvaluator{}
	_ Evaluator = SelectCounter(nil)
	_ Evaluator = (*Adder)(nil)
	_ Evaluator = (*Multiplier)(nil)
	_ Evaluator = (*CustomEvaluator)(nil)

	ErrBadEvaluator = errors.New("bad evaluator")
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
		"count": e.tag,
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

type SelectCounter PipelineSelector

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
		"count_if": ([]Selector)(e),
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

type EvaluatorFile struct {
	Evaluator Evaluator
}

func (f *EvaluatorFile) UnmarshalJSON(bytes []byte) error {
	var s string
	if err := json.Unmarshal(bytes, &s); err != nil {
		var e *json.UnmarshalTypeError
		if !errors.As(err, &e) {
			return err
		}
	} else if s == "" {
		return nil
	} else {
		if e, ok := map[string]Evaluator{
			"loss": LossEvaluator{},
		}[s]; ok {
			f.Evaluator = e
			return nil
		}

		return ErrBadEvaluator
	}

	var m map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &m); err != nil {
		return err
	}

	if c, ok := m["const"]; ok {
		var n int
		if err := json.Unmarshal(c, &n); err != nil {
			return err
		}

		f.Evaluator = ConstEvaluator(n)
		return nil
	}

	if a, ok := m["axis"]; ok {
		var axis Axis
		if err := json.Unmarshal(a, &axis); err != nil {
			return err
		}

		f.Evaluator = AxisEvaluator(axis)
		return nil
	}

	if count, ok := m["count"]; ok {
		var t TagFile
		if err := json.Unmarshal(count, &t); err != nil {
			return err
		}

		f.Evaluator = NewBuffCounter(t.Tag)
		return nil
	}

	if countIf, ok := m["count_if"]; ok {
		var fs []SelectorFile
		if err := json.Unmarshal(countIf, &fs); err != nil {
			return err
		}

		f.Evaluator = NewSelectCounter(functional.Map(func(f SelectorFile) Selector {
			return f.Selector
		})(fs)...)
		return nil
	}

	if _, ok := m["add"]; ok {
		var s struct {
			Add       int
			Evaluator EvaluatorFile
		}
		if err := json.Unmarshal(bytes, &s); err != nil {
			return err
		}

		f.Evaluator = NewAdder(s.Add, s.Evaluator.Evaluator)
		return nil
	}

	if _, ok := m["mul"]; ok {
		var s struct {
			Mul       int
			Evaluator EvaluatorFile
		}
		if err := json.Unmarshal(bytes, &s); err != nil {
			return err
		}

		f.Evaluator = NewMultiplier(s.Mul, s.Evaluator.Evaluator)
		return nil
	}

	return ErrBadEvaluator
}
