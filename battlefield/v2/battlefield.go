package battlefield

import "sort"

type BattleField struct {
	warriors []Warrior
	reactors []Reactor
}

func NewBattleField(warriors []Warrior, reactors ...Reactor) *BattleField {
	return &BattleField{warriors, reactors}
}

func (b *BattleField) React(signal ForkableSignal) {
	for _, r := range b.reactors {
		sig := signal.Fork(nil)
		r.React(sig)
		if s, ok := sig.(Renderer); ok {
			s.Render(b)
		}
	}

	warriors := make([]Warrior, len(b.warriors))
	copy(warriors, b.warriors)
	for i := 0; i < len(warriors); i++ {
		sort.Sort(&ByAxis{Speed, false, warriors[i:]})

		sig := signal.Fork(warriors[i])
		warriors[i].React(sig)
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
			sort.Sort(&ByAxis{Speed, false, warriors[i:]})

			sig := NewLaunchSignal(warriors[i])
			warriors[i].React(sig)
			sig.Render(b)

			// TODO: check if the battle is over
		}

		b.React(NewRoundEndSignal())
	}
}
