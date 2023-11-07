package battlefield

import "sort"

const (
	gaugeMul = 1_000
	gaugeMax = 1_000_000
)

var (
	_ EvaluationContext = (*BattleField)(nil)
)

type Renderer interface {
	Render(EvaluationContext)
}

type BattleField struct {
	Rng
	Sufferer
	warriors []Warrior
	reactors []Reactor
	deadline int
	sequence int
}

type Option func(*BattleField)

func NewBattleField(rng Rng, warriors []Warrior, options ...Option) *BattleField {
	f := &BattleField{
		Rng:      rng,
		Sufferer: ConstSufferer{},
		warriors: warriors,
		deadline: 1_000_000,
	}
	for _, opt := range options {
		opt(f)
	}

	return f
}

func FieldSufferer(s Sufferer) Option {
	return func(b *BattleField) {
		b.Sufferer = s
	}
}

func FieldReactor(r Reactor) Option {
	return func(b *BattleField) {
		b.reactors = append(b.reactors, r)
	}
}

func Deadline(deadline int) Option {
	return func(b *BattleField) {
		b.deadline = deadline
	}
}

func (b *BattleField) Warriors() []Warrior {
	return b.warriors
}

func (b *BattleField) React(signal RegularSignal) {
	b.react(signal, speedSorter{b}, b)
}

func (b *BattleField) react(signal RegularSignal, sorter warriorSorter, ec EvaluationContext) {
	for _, r := range b.reactors {
		sig := signal.SetCurrent(nil)
		r.React(sig, ec)
		if s, ok := sig.(Renderer); ok {
			s.Render(ec)
		}
	}

	death := make(map[Warrior]struct{})
	if postActionSignal, ok := signal.(*PostActionSignal); ok {
		deaths := postActionSignal.Deaths()
		sorter.Sort(deaths)

		for _, w := range deaths {
			death[w] = struct{}{}
			sig := signal.SetCurrent(w)
			w.React(sig, ec)
			if s, ok := sig.(Renderer); ok {
				s.Render(ec)
			}
		}
	}

	warriors := make([]Warrior, len(b.warriors))
	copy(warriors, b.warriors)
	for i := 0; i < len(warriors); i++ {
		sorter.Sort(warriors[i:])
		if _, ok := death[warriors[i]]; ok || warriors[i].Health().Current <= 0 {
			continue
		}

		sig := signal.SetCurrent(warriors[i])
		warriors[i].React(sig, ec)
		if s, ok := sig.(Renderer); ok {
			s.Render(ec)
		}
	}
}

func (b *BattleField) Next() int {
	b.sequence += 1
	return b.sequence
}

func (b *BattleField) Run() {
	b.React(NewBattleStartSignal(b.Next()))
	for b.sequence < b.deadline {
		b.React(NewRoundStartSignal(b.Next()))

		warriors := make([]Warrior, len(b.warriors))
		copy(warriors, b.warriors)
		for i := 0; i < len(warriors); i++ {
			sort.Sort(&ByAxis{Speed, false, b, warriors[i:]})
			if warriors[i].Health().Current <= 0 {
				continue
			}

			sig := NewLaunchSignal(b.Next(), warriors[i])
			warriors[i].React(sig, b)
			sig.Render(b)

			if b.end() {
				return
			}
		}

		b.React(NewRoundEndSignal(b.Next()))
	}
}

func (b *BattleField) end() bool {
	healthy := Healthy.Select(b.warriors, nil, b)
	left := AbsoluteSideSelector(Left).Select(healthy, nil, b)

	return len(left) == 0 || len(healthy) == len(left)
}

type ATBattleField struct {
	*BattleField
	progress map[Warrior]int
	elapsed  int
}

func NewATBattleField(rng Rng, warriors []Warrior, options ...Option) *ATBattleField {
	progress := make(map[Warrior]int)
	for _, w := range warriors {
		progress[w] = 0
	}

	return &ATBattleField{
		NewBattleField(rng, warriors, options...),
		progress,
		0,
	}
}

func (b *ATBattleField) Run() {
	b.React(NewBattleStartSignal(b.Next()))
	sort.Sort(&ByAxis{Speed, false, b, b.warriors})
	for _, w := range b.warriors {
		b.React(NewATRoundStartSignal(b.Next(), w))
	}
	var current Warrior
	for b.sequence < b.deadline && !b.end() {
		if current != nil {
			b.React(NewATRoundEndSignal(b.Next(), current))
			if b.progress[current] >= gaugeMax {
				b.progress[current] = 0
				b.React(NewATRoundStartSignal(b.Next(), current))
			}
		}

		// find the next warrior to launch
		forecasts := make([]forecast, 0, len(b.warriors))
		for w, p := range b.progress {
			if w.Health().Current <= 0 {
				b.progress[w] = 0
				continue
			}

			speed := w.Component(Speed, b)
			forecasts = append(forecasts, forecast{
				(gaugeMax - p + speed - 1) / speed,
				speed,
				w,
			})
		}
		if len(forecasts) == 0 {
			return
		}

		sort.Slice(forecasts, func(i, j int) bool {
			if forecasts[i].duration != forecasts[j].duration {
				return forecasts[i].duration < forecasts[j].duration
			}

			if forecasts[i].speed != forecasts[j].speed {
				return forecasts[i].speed > forecasts[j].speed
			}

			p, q := forecasts[i].warrior.Component(Position, b), forecasts[j].warrior.Component(Position, b)
			if p != q {
				return p < q
			}

			return forecasts[i].warrior.Side() == Left
		})
		current = forecasts[0].warrior
		t := forecasts[0].duration

		// update the progresses before launch
		for _, f := range forecasts {
			p := b.progress[f.warrior]
			p += f.speed * t
			if p >= gaugeMax {
				p = gaugeMax
			}
			b.progress[f.warrior] = p
		}
		b.elapsed += t

		sig := NewLaunchSignal(b.Next(), current)
		current.React(sig, b)
		sig.Render(b)
	}
}

func (b *ATBattleField) React(signal RegularSignal) {
	b.react(signal, progressSorter{b}, b)
}

func (b *ATBattleField) Progress(w Warrior) int {
	return b.progress[w] / gaugeMul
}

func (b *ATBattleField) Distance(w Warrior, progress int) {
	p := b.progress[w]
	p += progress * gaugeMul
	if p < 0 {
		p = 0
	} else if p >= gaugeMax {
		p = gaugeMax
	}

	b.progress[w] = p
}

func (b *ATBattleField) Interval(w Warrior, milli int) {
	v := w.Component(Speed, b)
	t := milli * gaugeMul / 1000

	b.Distance(w, v*t)
}

func (b *ATBattleField) Milli() int {
	return b.elapsed * 1000 / gaugeMul
}

type forecast struct {
	duration int
	speed    int
	warrior  Warrior
}

type warriorSorter interface {
	Sort([]Warrior)
}

type speedSorter struct {
	ec EvaluationContext
}

func (s speedSorter) Sort(warriors []Warrior) {
	sort.Sort(&ByAxis{Speed, false, s.ec, warriors})
}

type progressSorter struct {
	f *ATBattleField
}

func (s progressSorter) Sort(warriors []Warrior) {
	sort.Slice(warriors, func(i, j int) bool {
		if s.f.progress[warriors[i]] != s.f.progress[warriors[j]] {
			return s.f.progress[warriors[i]] > s.f.progress[warriors[j]]
		}

		p, q := warriors[i].Component(Speed, s.f), warriors[j].Component(Speed, s.f)
		if p != q {
			return p > q
		}

		p, q = warriors[i].Component(Position, s.f), warriors[j].Component(Position, s.f)
		if p != q {
			return p < q
		}

		return warriors[i].Side() == Left
	})
}
