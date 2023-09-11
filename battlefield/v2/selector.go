package battlefield

import (
	"encoding/json"
	"errors"
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
		"side": Side(s).String(),
	})
}

type SideSelector bool

func (s SideSelector) Select(inputs []Warrior, signal Signal, ec EvaluationContext) []Warrior {
	side := Side(bool(s) == bool(signal.Current().(Warrior).Side()))
	return AbsoluteSideSelector(side).Select(inputs, signal, ec)
}

func (s SideSelector) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]bool{
		"side": bool(s),
	})
}

type CurrentSelector struct {
}

func (CurrentSelector) Select(_ []Warrior, signal Signal, _ EvaluationContext) []Warrior {
	return []Warrior{signal.Current().(Warrior)}
}

func (CurrentSelector) MarshalJSON() ([]byte, error) {
	return json.Marshal("current")
}

type SourceSelector struct {
}

func (SourceSelector) Select(_ []Warrior, signal Signal, _ EvaluationContext) []Warrior {
	source, _ := signal.(ActionSignal).Action().Script().Source()
	return []Warrior{source.(Warrior)}
}

func (SourceSelector) MarshalJSON() ([]byte, error) {
	return json.Marshal("source")
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
		"sort": s.axis.String(),
		"asc":  s.asc,
	})
}

type ShuffleSelector struct {
	rng        Rng
	preference any
}

func NewShuffleSelector(rng Rng, preference any) ShuffleSelector {
	return ShuffleSelector{rng, preference}
}

func (s ShuffleSelector) Select(inputs []Warrior, _ Signal, _ EvaluationContext) (outputs []Warrior) {
	if len(inputs) < 2 {
		return inputs
	}

	outputs = make([]Warrior, len(inputs))
	copy(outputs, inputs)

	randoms := make([]int, len(inputs))
	for i := range randoms {
		randoms[i] = int(s.rng.Float64() * 1e6)
	}

	sort.Sort(&shuffle{s.preference, randoms, outputs})
	return
}

func (s ShuffleSelector) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]any{
		"shuffle": s.preference,
	})
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
	return json.Marshal(map[string]int{
		"take": int(s),
	})
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
	return json.Marshal(map[string]uint8{"counter_position": uint8(s)})
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
		"take_while": s.comparator.String(),
		"evaluator":  s.evaluator,
		"value":      s.value,
	})
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
	return json.Marshal([]Selector(s))
}

type SelectorFile struct {
	Selector Selector
}

func (f *SelectorFile) UnmarshalJSON(bytes []byte) error {
	var s string
	if err := json.Unmarshal(bytes, &s); err != nil {
		var e *json.UnmarshalTypeError
		if !errors.As(err, &e) {
			return err
		}
	} else {
		if selector, ok := map[string]Selector{
			"current": CurrentSelector{},
			"source":  SourceSelector{},
		}[s]; ok {
			f.Selector = selector
			return nil
		}

		return ErrBadSelector
	}

	var m map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &m); err != nil {
		var e *json.UnmarshalTypeError
		if !errors.As(err, &e) {
			return err
		}

		var fs []SelectorFile
		if err := json.Unmarshal(bytes, &fs); err != nil {
			return err
		}

		f.Selector = PipelineSelector(functional.Map(func(f SelectorFile) Selector {
			return f.Selector
		})(fs))
		return nil
	}

	if side, ok := m["side"]; ok {
		var s any
		if err := json.Unmarshal(side, &s); err != nil {
			return err
		}

		switch s := s.(type) {
		case bool:
			f.Selector = SideSelector(s)
			return nil
		case string:
			if side, ok := map[string]Side{
				"Left":  Left,
				"Right": Right,
			}[s]; ok {
				f.Selector = AbsoluteSideSelector(side)
				return nil
			}
		}

		return ErrBadSelector
	}

	if _, ok := m["sort"]; ok {
		var s struct {
			Sort Axis
			Asc  bool
		}
		if err := json.Unmarshal(bytes, &s); err != nil {
			return err
		}

		f.Selector = NewSortSelector(s.Sort, s.Asc)
		return nil
	}

	if preference, ok := m["shuffle"]; ok {
		var t TagFile
		if err := json.Unmarshal(preference, &t); err != nil {
			return err
		}

		f.Selector = NewShuffleSelector(DefaultRng, t.Tag)
		return nil
	}

	if take, ok := m["take"]; ok {
		var n int
		if err := json.Unmarshal(take, &n); err != nil {
			return err
		}

		f.Selector = FrontSelector(n)
		return nil
	}

	if counterPosition, ok := m["counter_position"]; ok {
		var r uint8
		if err := json.Unmarshal(counterPosition, &r); err != nil {
			return err
		}

		f.Selector = NewCounterPositionSelector(r)
		return nil
	}

	if _, ok := m["take_while"]; ok {
		var s struct {
			TakeWhile IntComparator `json:"take_while"`
			Evaluator EvaluatorFile
			Value     int
		}
		if err := json.Unmarshal(bytes, &s); err != nil {
			return err
		}

		f.Selector = NewWaterLevelSelector(s.TakeWhile, s.Evaluator.Evaluator, s.Value)
		return nil
	}

	return ErrBadSelector
}
