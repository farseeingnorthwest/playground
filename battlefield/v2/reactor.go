package battlefield

type Reactor interface {
	React(Signal)
}

type Forker interface {
	Fork() interface{}
}

type Actor interface {
	Act(source *Fighter, targets []*Fighter) *Action
}

type Evaluator struct {
	Axis
	Percentage int
}

func (e *Evaluator) Evaluate(warrior *Warrior) int {
	var value int
	switch e.Axis {
	case Damage:
		value = warrior.Damage()
	case Defense:
		value = warrior.Defense()
	case Health:
		r, m := warrior.Health()
		value = r.Current * m / r.Maximum
	default:
		panic("bad axis")
	}

	return value * e.Percentage / 100
}

type Attacker struct {
	Evaluator
}

func (a *Attacker) Act(source *Fighter, targets []*Fighter) *Action {
	return &Action{
		Source:  source,
		Targets: targets,
		Verb:    NewAttack(a.Evaluate(source.Warrior)),
	}
}

type Healer struct {
	Evaluator
}

func (h *Healer) Act(source *Fighter, targets []*Fighter) *Action {
	return &Action{
		Source:  source,
		Targets: targets,
		Verb:    NewHeal(h.Evaluate(source.Warrior)),
	}
}

type Buffer struct {
	reactor interface {
		Reactor
		Forker
	}
}

func (b *Buffer) Act(source *Fighter, targets []*Fighter) *Action {
	return &Action{
		Source:  source,
		Targets: targets,
		Verb:    NewBuffing(b.reactor.Fork().(Reactor)),
	}
}

type SelectiveActor struct {
	Selector
	Actor
}

func (s SelectiveActor) Act(source *Fighter, targets []*Fighter) *Action {
	selected := s.Selector.Select(source, targets)
	if len(selected) == 0 {
		return nil
	}

	return s.Actor.Act(source, selected)
}

type Rng interface {
	Gen() float64 // [0, 1)
}

type Critical struct {
	rng  Rng
	odds int // percentage
	buff *ClearingBuff
}

func (c *Critical) React(signal Signal) {
	sig, ok := signal.(*PreActionSignal)
	if !ok {
		return
	}
	_, ok = sig.Verb.(*Attack)
	if !ok {
		return
	}

	if float64(c.odds)/100 <= c.rng.Gen() {
		return
	}

	sig.Append(&Action{
		Source:  sig.Source,
		Targets: sig.Targets,
		Verb:    NewBuffing(c.buff.Fork(sig.Action)),
	})
}
