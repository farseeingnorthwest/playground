package battlefield

import "sort"

type BattleField struct {
	warriors []Warrior
	reactors []Reactor
}

func NewBattleField(warriors []Warrior, reactors ...Reactor) *BattleField {
	return &BattleField{warriors, reactors}
}

func (b *BattleField) Warriors() []Warrior {
	return b.warriors
}

func (b *BattleField) React(signal ForkableSignal) {
	for _, r := range b.reactors {
		sig := signal.Fork(nil)
		r.React(sig, b)
		if s, ok := sig.(Renderer); ok {
			s.Render(b)
		}
	}

	warriors := make([]Warrior, len(b.warriors))
	copy(warriors, b.warriors)
	for i := 0; i < len(warriors); i++ {
		sort.Sort(&ByAxis{Speed, false, b, warriors[i:]})

		sig := signal.Fork(warriors[i])
		warriors[i].React(sig, b)
		if s, ok := sig.(Renderer); ok {
			s.Render(b)
		}
	}
}

func (b *BattleField) Run() {
	b.React(NewBattleStartSignal())
	for {
		b.React(NewRoundStartSignal())

		warriors := make([]Warrior, len(b.warriors))
		copy(warriors, b.warriors)
		for i := 0; i < len(warriors); i++ {
			sort.Sort(&ByAxis{Speed, false, b, warriors[i:]})

			sig := NewLaunchSignal(warriors[i])
			warriors[i].React(sig, b)
			sig.Render(b)

			healthy := Healthy.Select(b.warriors, nil, b)
			left := AbsoluteSideSelector(Left).Select(healthy, nil, b)
			if len(left) == 0 || len(healthy) == len(left) {
				return
			}
		}

		b.React(NewRoundEndSignal())
	}
}
