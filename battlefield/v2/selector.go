package battlefield

import "sort"

var (
	Healthy = &WaterLevelSelector{Gt, AxisEvaluator(Health), 0}
)

type Selector interface {
	Select([]Warrior, Signal) []Warrior
}

type AbsoluteSideSelector Side

func (s AbsoluteSideSelector) Select(inputs []Warrior, _ Signal) (outputs []Warrior) {
	for _, w := range inputs {
		if w.Side() == Side(s) {
			outputs = append(outputs, w)
		}
	}

	return
}

type SideSelector bool

func (s SideSelector) Select(inputs []Warrior, signal Signal) []Warrior {
	side := Side(bool(s) == bool(signal.Current().(Warrior).Side()))
	return AbsoluteSideSelector(side).Select(inputs, signal)
}

type CurrentSelector struct {
}

func (s *CurrentSelector) Select(_ []Warrior, signal Signal) []Warrior {
	return []Warrior{signal.Current().(Warrior)}
}

type SourceSelector struct {
}

func (s *SourceSelector) Select(_ []Warrior, signal Signal) []Warrior {
	source, _ := signal.(Sourcer).Source()
	return []Warrior{source.(Warrior)}
}

type SortSelector struct {
	axis Axis
	asc  bool
}

func (s *SortSelector) Select(inputs []Warrior, _ Signal) (outputs []Warrior) {
	outputs = make([]Warrior, len(inputs))
	copy(outputs, inputs)

	sort.Sort(&ByAxis{s.axis, s.asc, outputs})
	return
}

type ShuffleSelector struct {
	rng        Rng
	preference any
}

func NewShuffleSelector(rng Rng, preference any) *ShuffleSelector {
	return &ShuffleSelector{rng, preference}
}

func (s *ShuffleSelector) Select(inputs []Warrior, _ Signal) (outputs []Warrior) {
	outputs = make([]Warrior, len(inputs))
	copy(outputs, inputs)

	randoms := make([]int, len(inputs))
	for i := range randoms {
		randoms[i] = int(s.rng.Float64() * 1e6)
	}

	sort.Sort(&shuffle{s.preference, randoms, outputs})
	return
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

func (s FrontSelector) Select(inputs []Warrior, _ Signal) (outputs []Warrior) {
	if len(inputs) <= int(s) {
		return inputs
	}

	outputs = make([]Warrior, int(s))
	copy(outputs, inputs[:int(s)])
	return
}

type WaterLevelSelector struct {
	comparator IntComparator
	evaluator  Evaluator
	value      int
}

func (s *WaterLevelSelector) Select(inputs []Warrior, _ Signal) (outputs []Warrior) {
	for _, w := range inputs {
		if s.comparator.Compare(s.evaluator.Evaluate(w), s.value) {
			outputs = append(outputs, w)
		}
	}

	return
}

type CounterPositionSelector struct {
}

func (s *CounterPositionSelector) Select(inputs []Warrior, signal Signal) []Warrior {
	current := signal.Current().(Warrior)
	for _, w := range inputs {
		if w != current && w.Component(Position) == current.Component(Position) {
			return []Warrior{w}
		}
	}

	return nil
}
