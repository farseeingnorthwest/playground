package battlefield

import "github.com/farseeingnorthwest/playground/battlefield/v2/evaluation"

type Reactor interface {
	React(Signal)
	Fork(*evaluation.Block, Signal) Reactor
}

type Actor interface {
	Act(*Fighter, []*Fighter) *Action
}

type BlindActor struct {
	proto Verb
	*evaluation.Bundle
}

func (a *BlindActor) Act(source *Fighter, targets []*Fighter) *Action {
	return &Action{
		Source:  source,
		Targets: targets,
		Verb:    a.proto.Fork(a.ForkWith(source.Warrior), nil),
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

type ProbabilityAttackReactor struct {
	rng   Rng
	odds  int // percentage
	proto Verb
}

func (c *ProbabilityAttackReactor) React(signal Signal) {
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
		Verb:    c.proto.Fork(nil, signal),
	})
}

func (c *ProbabilityAttackReactor) Fork(_ *evaluation.Block, _ Signal) Reactor {
	return &ProbabilityAttackReactor{
		rng:   c.rng,
		odds:  c.odds,
		proto: c.proto,
	}
}
