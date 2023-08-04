package battlefield

import "sort"

type Selector interface {
	Select([]Warrior, Signal) []Warrior
}

type SideSelector struct {
	own bool
}

func (s *SideSelector) Select(inputs []Warrior, signal Signal) (outputs []Warrior) {
	for _, w := range inputs {
		if w.Side() == signal.Current().(Warrior).Side() == s.own {
			outputs = append(outputs, w)
		}
	}

	return
}

type CurrentSelector struct {
}

func (s *CurrentSelector) Select(_ []Warrior, signal Signal) []Warrior {
	return []Warrior{signal.Current().(Warrior)}
}

type AxisSelector struct {
	axis Axis
	asc  bool
}

func (s *AxisSelector) Select(inputs []Warrior, _ Signal) (outputs []Warrior) {
	outputs = make([]Warrior, len(inputs))
	copy(outputs, inputs)

	sort.Sort(&ByAxis{s.axis, s.asc, outputs})
	return
}

type CounterSelector struct {
}

func (s *CounterSelector) Select(inputs []Warrior, signal Signal) []Warrior {
	current := signal.Current().(Warrior)
	for _, w := range inputs {
		if w != current && w.Component(Position) == current.Component(Position) {
			return []Warrior{w}
		}
	}

	return nil
}

type ShuffleSelector struct {
	randInt    func() int
	preference any
}

func (s *ShuffleSelector) Select(inputs []Warrior, _ Signal) (outputs []Warrior) {
	outputs = make([]Warrior, len(inputs))
	copy(outputs, inputs)

	randoms := make([]int, len(inputs))
	for i := range randoms {
		randoms[i] = s.randInt()
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
	if s.warriors[i].Contains(s.preference) != s.warriors[j].Contains(s.preference) {
		return s.warriors[i].Contains(s.preference)
	}

	return s.randoms[i] < s.randoms[j]
}

type FrontSelector struct {
	count int
}

func (s *FrontSelector) Select(inputs []Warrior, _ Signal) (outputs []Warrior) {
	if len(inputs) <= s.count {
		return inputs
	}

	outputs = make([]Warrior, s.count)
	copy(outputs, inputs[:s.count])
	return
}

type SourceSelector struct {
}

func (s *SourceSelector) Select(_ []Warrior, signal Signal) []Warrior {
	source, _ := signal.(Sourcer).Source()
	return []Warrior{source.(Warrior)}
}
