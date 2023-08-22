package battlefield

import "sort"

var (
	_ EvaluationContext = (*BattleField)(nil)
)

type Renderer interface {
	Render(*BattleField)
}

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

func (b *BattleField) React(signal RegularSignal) {
	for _, r := range b.reactors {
		sig := signal.SetCurrent(nil)
		r.React(sig, b)
		if s, ok := sig.(Renderer); ok {
			s.Render(b)
		}
	}

	death := make(map[Warrior]struct{})
	if postActionSignal, ok := signal.(*PostActionSignal); ok {
		deaths := postActionSignal.Deaths()
		sort.Sort(&ByAxis{Speed, false, b, deaths})

		for _, w := range deaths {
			death[w] = struct{}{}
			sig := signal.SetCurrent(w)
			w.React(sig, b)
			if s, ok := sig.(Renderer); ok {
				s.Render(b)
			}
		}
	}

	warriors := make([]Warrior, len(b.warriors))
	copy(warriors, b.warriors)
	for i := 0; i < len(warriors); i++ {
		sort.Sort(&ByAxis{Speed, false, b, warriors[i:]})
		if _, ok := death[warriors[i]]; ok || warriors[i].Health().Current <= 0 {
			continue
		}

		sig := signal.SetCurrent(warriors[i])
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
			if warriors[i].Health().Current <= 0 {
				continue
			}

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
