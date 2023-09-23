package battlefield

import (
	"encoding/json"
	"errors"
	"reflect"

	"github.com/farseeingnorthwest/playground/battlefield/v2/functional"
)

var (
	_ Evaluator = ConstEvaluator(0)
	_ Evaluator = AxisEvaluator(0)
	_ Evaluator = BuffCounter{}
	_ Evaluator = LossEvaluator{}
	_ Evaluator = SelectCounter(nil)
	_ Evaluator = Adder{}
	_ Evaluator = Multiplier{}
	_ Evaluator = CustomEvaluator{}

	ErrBadEvaluator = errors.New("bad evaluator")
)

type Evaluator interface {
	Evaluate(Warrior, EvaluationContext) int
}

type EvaluationContext interface {
	Warriors() []Warrior
	React(RegularSignal)
	Next() int
}

type ConstEvaluator int

func (e ConstEvaluator) Evaluate(Warrior, EvaluationContext) int {
	return int(e)
}

func (e ConstEvaluator) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]any{
		"_kind": "const",
		"value": int(e),
	})
}

func (e *ConstEvaluator) UnmarshalJSON(bytes []byte) error {
	var v struct{ Value int }
	if err := json.Unmarshal(bytes, &v); err != nil {
		return err
	}

	*e = ConstEvaluator(v.Value)
	return nil
}

type AxisEvaluator Axis

func (e AxisEvaluator) Evaluate(warrior Warrior, ec EvaluationContext) int {
	return warrior.Component(Axis(e), ec)
}

func (e AxisEvaluator) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]string{
		"_kind": "axis",
		"axis":  Axis(e).String(),
	})
}

func (e *AxisEvaluator) UnmarshalJSON(bytes []byte) error {
	var v struct{ Axis Axis }
	if err := json.Unmarshal(bytes, &v); err != nil {
		return err
	}

	*e = AxisEvaluator(v.Axis)
	return nil
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
		"_kind": "buff_count",
		"tag":   e.tag,
	})
}

func (e *BuffCounter) UnmarshalJSON(bytes []byte) error {
	var v struct{ Tag TagFile }
	if err := json.Unmarshal(bytes, &v); err != nil {
		return err
	}

	e.tag = v.Tag.Tag
	return nil
}

type LossEvaluator struct {
}

func (LossEvaluator) Evaluate(_ Warrior, ec EvaluationContext) int {
	return ec.(ActionContext).Action().Script().Loss()
}

func (LossEvaluator) MarshalJSON() ([]byte, error) {
	return json.Marshal(kind{"loss"})
}

type SelectCounter PipelineSelector

func NewSelectCounter(selectors ...Selector) SelectCounter {
	return selectors
}

func (e SelectCounter) Evaluate(warrior Warrior, ec EvaluationContext) int {
	signal := NewFreeSignal(ec.Next(), warrior)
	warriors := ec.Warriors()
	for _, selector := range e {
		warriors = selector.Select(warriors, signal, ec)
	}

	return len(warriors)
}

func (e SelectCounter) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]any{
		"_kind":     "select_counter",
		"selectors": ([]Selector)(e),
	})
}

func (e *SelectCounter) UnmarshalJSON(bytes []byte) error {
	var v struct{ Selectors []SelectorFile }
	if err := json.Unmarshal(bytes, &v); err != nil {
		return err
	}

	*e = NewSelectCounter(functional.Map(func(f SelectorFile) Selector {
		return f.Selector
	})(v.Selectors)...)
	return nil
}

type Adder struct {
	adder     int
	evaluator Evaluator
}

func NewAdder(adder int, evaluator Evaluator) Adder {
	return Adder{adder, evaluator}
}

func (e Adder) Evaluate(warrior Warrior, ec EvaluationContext) int {
	return e.adder + e.evaluator.Evaluate(warrior, ec)
}

func (e Adder) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]any{
		"_kind":     "adder",
		"adder":     e.adder,
		"evaluator": e.evaluator,
	})
}

func (e *Adder) UnmarshalJSON(bytes []byte) error {
	var v struct {
		Adder     int
		Evaluator EvaluatorFile
	}
	if err := json.Unmarshal(bytes, &v); err != nil {
		return err
	}

	*e = Adder{v.Adder, v.Evaluator.Evaluator}
	return nil
}

type Multiplier struct {
	multiplier int
	evaluator  Evaluator
}

func NewMultiplier(multiplier int, evaluator Evaluator) Multiplier {
	return Multiplier{multiplier, evaluator}
}

func (e Multiplier) Evaluate(warrior Warrior, ec EvaluationContext) int {
	return e.multiplier * e.evaluator.Evaluate(warrior, ec) / 100
}

func (e Multiplier) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]any{
		"_kind":      "multiplier",
		"multiplier": e.multiplier,
		"evaluator":  e.evaluator,
	})
}

func (e *Multiplier) UnmarshalJSON(bytes []byte) error {
	var v struct {
		Multiplier int
		Evaluator  EvaluatorFile
	}
	if err := json.Unmarshal(bytes, &v); err != nil {
		return err
	}

	*e = Multiplier{v.Multiplier, v.Evaluator.Evaluator}
	return nil
}

type CustomEvaluator struct {
	evaluator func(Warrior, EvaluationContext) int
}

func NewCustomEvaluator(evaluator func(Warrior, EvaluationContext) int) CustomEvaluator {
	return CustomEvaluator{evaluator}
}

func (e CustomEvaluator) Evaluate(warrior Warrior, ec EvaluationContext) int {
	return e.evaluator(warrior, ec)
}

type EvaluatorFile struct {
	Evaluator Evaluator
}

func (f *EvaluatorFile) UnmarshalJSON(bytes []byte) error {
	var k kind
	if err := json.Unmarshal(bytes, &k); err != nil {
		return err
	}

	if evaluator, ok := evaluatorType[k.Kind]; ok {
		v := reflect.New(evaluator)
		if err := json.Unmarshal(bytes, v.Interface()); err != nil {
			return err
		}

		f.Evaluator = v.Elem().Interface().(Evaluator)
		return nil
	}

	return ErrBadEvaluator
}

var evaluatorType = map[string]reflect.Type{
	"const":          reflect.TypeOf(ConstEvaluator(0)),
	"axis":           reflect.TypeOf(AxisEvaluator(0)),
	"buff_count":     reflect.TypeOf(BuffCounter{}),
	"loss":           reflect.TypeOf(LossEvaluator{}),
	"select_counter": reflect.TypeOf(SelectCounter(nil)),
	"adder":          reflect.TypeOf(Adder{}),
	"multiplier":     reflect.TypeOf(Multiplier{}),
}
