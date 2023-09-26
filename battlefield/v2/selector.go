package battlefield

import (
	"encoding/json"
	"errors"
	"reflect"
	"sort"

	"github.com/farseeingnorthwest/playground/battlefield/v2/functional"
)

var (
	Healthy          = WaterLevelSelector{Gt, AxisEvaluator(Health), 0}
	_       Selector = AbsoluteSideSelector(false)
	_       Selector = SideSelector(false)
	_       Selector = CurrentSelector{}
	_       Selector = SourceSelector{}
	_       Selector = SortSelector{}
	_       Selector = ShuffleSelector{}
	_       Selector = FrontSelector(0)
	_       Selector = CounterPositionSelector(0)
	_       Selector = WaterLevelSelector{}
	_       Selector = PipelineSelector{}

	ErrBadSelector = errors.New("bad selector")
)

type Selector interface {
	Select([]Warrior, Signal, EvaluationContext) []Warrior
}

type AbsoluteSideSelector Side

func (s AbsoluteSideSelector) Select(inputs []Warrior, _ Signal, _ EvaluationContext) (outputs []Warrior) {
	for _, w := range inputs {
		if w.Side() == Side(s) {
			outputs = append(outputs, w)
		}
	}

	return
}

func (s AbsoluteSideSelector) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]string{
		"_kind": "absolute_side",
		"side":  Side(s).String(),
	})
}

func (s *AbsoluteSideSelector) UnmarshalJSON(bytes []byte) error {
	var v struct{ Side Side }
	if err := json.Unmarshal(bytes, &v); err != nil {
		return err
	}

	*s = AbsoluteSideSelector(v.Side)
	return nil
}

type SideSelector bool

func (s SideSelector) Select(inputs []Warrior, signal Signal, ec EvaluationContext) []Warrior {
	side := Side(bool(s) == bool(signal.Current().(Warrior).Side()))
	return AbsoluteSideSelector(side).Select(inputs, signal, ec)
}

func (s SideSelector) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]any{
		"_kind": "side",
		"side":  bool(s),
	})
}

func (s *SideSelector) UnmarshalJSON(bytes []byte) error {
	var v struct{ Side bool }
	if err := json.Unmarshal(bytes, &v); err != nil {
		return err
	}

	*s = SideSelector(v.Side)
	return nil
}

type CurrentSelector struct {
}

func (CurrentSelector) Select(_ []Warrior, signal Signal, _ EvaluationContext) []Warrior {
	return []Warrior{signal.Current().(Warrior)}
}

func (CurrentSelector) MarshalJSON() ([]byte, error) {
	return json.Marshal(kind{"current"})
}

type SourceSelector struct {
}

func (SourceSelector) Select(_ []Warrior, signal Signal, _ EvaluationContext) []Warrior {
	_, source, _ := signal.(ActionSignal).Action().Script().Source()
	return []Warrior{source.(Warrior)}
}

func (SourceSelector) MarshalJSON() ([]byte, error) {
	return json.Marshal(kind{"source"})
}

type SortSelector struct {
	axis Axis
	asc  bool
}

func NewSortSelector(axis Axis, asc bool) SortSelector {
	return SortSelector{axis, asc}
}

func (s SortSelector) Select(inputs []Warrior, _ Signal, ec EvaluationContext) (outputs []Warrior) {
	outputs = make([]Warrior, len(inputs))
	copy(outputs, inputs)

	sort.Sort(&ByAxis{s.axis, s.asc, ec, outputs})
	return
}

func (s SortSelector) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]any{
		"_kind": "sort",
		"axis":  s.axis.String(),
		"asc":   s.asc,
	})
}

func (s *SortSelector) UnmarshalJSON(bytes []byte) error {
	var v struct {
		Axis Axis
		Asc  bool
	}
	if err := json.Unmarshal(bytes, &v); err != nil {
		return err
	}

	*s = NewSortSelector(v.Axis, v.Asc)
	return nil
}

type ShuffleSelector struct {
	preference any
}

func NewShuffleSelector(preference any) ShuffleSelector {
	return ShuffleSelector{preference}
}

func (s ShuffleSelector) Select(inputs []Warrior, _ Signal, ec EvaluationContext) (outputs []Warrior) {
	if len(inputs) < 2 {
		return inputs
	}

	outputs = make([]Warrior, len(inputs))
	copy(outputs, inputs)

	randoms := make([]int, len(inputs))
	for i := range randoms {
		randoms[i] = int(ec.Float64() * 1e6)
	}

	sort.Sort(&shuffle{s.preference, randoms, outputs})
	return
}

func (s ShuffleSelector) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]any{
		"_kind":      "shuffle",
		"preference": s.preference,
	})
}

func (s *ShuffleSelector) UnmarshalJSON(bytes []byte) error {
	var v struct{ Preference TagFile }
	if err := json.Unmarshal(bytes, &v); err != nil {
		return err
	}

	*s = NewShuffleSelector(v.Preference.Tag)
	return nil
}

type shuffle struct {
	preference any
	randoms    []int
	warriors   []Warrior
}

func (s *shuffle) Len() int {
	return len(s.warriors)
}

func (s *shuffle) Swap(i, j int) {
	s.warriors[i], s.warriors[j] = s.warriors[j], s.warriors[i]
	s.randoms[i], s.randoms[j] = s.randoms[j], s.randoms[i]
}

func (s *shuffle) Less(i, j int) bool {
	if s.preference != nil {
		p, q := len(s.warriors[i].Buffs(s.preference)) > 0, len(s.warriors[j].Buffs(s.preference)) > 0
		if p != q {
			return p
		}
	}

	return s.randoms[i] < s.randoms[j]
}

type FrontSelector int

func (s FrontSelector) Select(inputs []Warrior, _ Signal, _ EvaluationContext) (outputs []Warrior) {
	if len(inputs) <= int(s) {
		return inputs
	}

	outputs = make([]Warrior, s)
	copy(outputs, inputs[:s])
	return
}

func (s FrontSelector) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]any{
		"_kind": "front",
		"count": int(s),
	})
}

func (s *FrontSelector) UnmarshalJSON(bytes []byte) error {
	var v struct{ Count int }
	if err := json.Unmarshal(bytes, &v); err != nil {
		return err
	}

	*s = FrontSelector(v.Count)
	return nil
}

type CounterPositionSelector uint8

func NewCounterPositionSelector(r uint8) CounterPositionSelector {
	return CounterPositionSelector(r)
}

func (s CounterPositionSelector) Select(inputs []Warrior, signal Signal, _ EvaluationContext) (warriors []Warrior) {
	current := signal.Current().(Warrior)
	for _, w := range inputs {
		if w.Side() == current.Side() {
			continue
		}
		d := w.Component(Position, nil) - current.Component(Position, nil)
		if d < 0 {
			d = -d
		}
		if d <= int(s) {
			warriors = append(warriors, w)
		}
	}

	return
}

func (s CounterPositionSelector) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]any{
		"_kind": "counter_position",
		"r":     uint8(s),
	})
}

func (s *CounterPositionSelector) UnmarshalJSON(bytes []byte) error {
	var v struct{ R uint8 }
	if err := json.Unmarshal(bytes, &v); err != nil {
		return err
	}

	*s = CounterPositionSelector(v.R)
	return nil
}

type WaterLevelSelector struct {
	comparator IntComparator
	evaluator  Evaluator
	value      int
}

func NewWaterLevelSelector(comparator IntComparator, evaluator Evaluator, value int) WaterLevelSelector {
	return WaterLevelSelector{comparator, evaluator, value}
}

func (s WaterLevelSelector) Select(inputs []Warrior, _ Signal, ec EvaluationContext) (outputs []Warrior) {
	for _, w := range inputs {
		if s.comparator.Compare(s.evaluator.Evaluate(w, ec), s.value) {
			outputs = append(outputs, w)
		}
	}

	return
}

func (s WaterLevelSelector) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]any{
		"_kind":      "water_level",
		"comparator": s.comparator.String(),
		"evaluator":  s.evaluator,
		"value":      s.value,
	})
}

func (s *WaterLevelSelector) UnmarshalJSON(bytes []byte) error {
	var v struct {
		Comparator IntComparator
		Evaluator  EvaluatorFile
		Value      int
	}
	if err := json.Unmarshal(bytes, &v); err != nil {
		return err
	}

	*s = NewWaterLevelSelector(v.Comparator, v.Evaluator.Evaluator, v.Value)
	return nil
}

type PipelineSelector []Selector

func (s PipelineSelector) Select(inputs []Warrior, signal Signal, ec EvaluationContext) (outputs []Warrior) {
	outputs = inputs
	for _, selector := range s {
		outputs = selector.Select(outputs, signal, ec)
	}

	return
}

func (s PipelineSelector) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]any{
		"_kind":     "pipeline",
		"selectors": []Selector(s),
	})
}

func (s *PipelineSelector) UnmarshalJSON(bytes []byte) error {
	var v struct{ Selectors []SelectorFile }
	if err := json.Unmarshal(bytes, &v); err != nil {
		return err
	}

	*s = PipelineSelector(functional.Map(func(f SelectorFile) Selector {
		return f.Selector
	})(v.Selectors))
	return nil
}

type SelectorFile struct {
	Selector Selector
}

func (f *SelectorFile) UnmarshalJSON(bytes []byte) error {
	var k kind
	if err := json.Unmarshal(bytes, &k); err != nil {
		return err
	}

	if selector, ok := selectorType[k.Kind]; ok {
		v := reflect.New(selector)
		if err := json.Unmarshal(bytes, v.Interface()); err != nil {
			return err
		}

		f.Selector = v.Elem().Interface().(Selector)
		return nil
	}

	return ErrBadSelector
}

var selectorType = map[string]reflect.Type{
	"absolute_side":    reflect.TypeOf(AbsoluteSideSelector(false)),
	"side":             reflect.TypeOf(SideSelector(false)),
	"current":          reflect.TypeOf(CurrentSelector{}),
	"source":           reflect.TypeOf(SourceSelector{}),
	"sort":             reflect.TypeOf(SortSelector{}),
	"shuffle":          reflect.TypeOf(ShuffleSelector{}),
	"front":            reflect.TypeOf(FrontSelector(0)),
	"counter_position": reflect.TypeOf(CounterPositionSelector(0)),
	"water_level":      reflect.TypeOf(WaterLevelSelector{}),
	"pipeline":         reflect.TypeOf(PipelineSelector{}),
}
