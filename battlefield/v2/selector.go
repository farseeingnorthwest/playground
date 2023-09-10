package battlefield

import (
	"encoding/json"
	"sort"
)

var (
	Healthy          = &WaterLevelSelector{Gt, AxisEvaluator(Health), 0}
	_       Selector = AbsoluteSideSelector(false)
	_       Selector = SideSelector(false)
	_       Selector = CurrentSelector{}
	_       Selector = SourceSelector{}
	_       Selector = SortSelector{}
	_       Selector = ShuffleSelector{}
	_       Selector = FrontSelector(0)
	_       Selector = CounterPositionSelector{}
	_       Selector = WaterLevelSelector{}
	_       Selector = PipelineSelector{}
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

type CounterPositionSelector struct {
}

func (CounterPositionSelector) Select(inputs []Warrior, signal Signal, _ EvaluationContext) []Warrior {
	current := signal.Current().(Warrior)
	for _, w := range inputs {
		if w != current && w.Component(Position, nil) == current.Component(Position, nil) {
			return []Warrior{w}
		}
	}

	return nil
}

func (CounterPositionSelector) MarshalJSON() ([]byte, error) {
	return json.Marshal("counter_position")
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
