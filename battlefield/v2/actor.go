package battlefield

import "github.com/farseeingnorthwest/playground/battlefield/v2/evaluation"

type Actor interface {
	Act(*Fighter, []*Fighter, Signal) *Action
}

type BlindActor struct {
	proto Verb
	*evaluation.Bundle
}

func (a *BlindActor) Act(source *Fighter, targets []*Fighter, signal Signal) *Action {
	return &Action{
		Source:  source,
		Targets: targets,
		Verb:    a.proto.Fork(a.ForkWith(source.Warrior), signal),
	}
}

type SelectiveActor struct {
	Selector
	Actor
}

func (s SelectiveActor) Act(source *Fighter, targets []*Fighter, signal Signal) *Action {
	selected := s.Selector.Select(source, targets)
	if len(selected) == 0 {
		return nil
	}

	return s.Actor.Act(source, selected, signal)
}

type Rng interface {
	Gen() float64 // [0, 1)
}

type ProbabilityActor struct {
	rng  Rng
	odds int
	Actor
}

func (p ProbabilityActor) Act(source *Fighter, targets []*Fighter, signal Signal) *Action {
	if p.rng.Gen() > float64(p.odds)/100 {
		return nil
	}

	return p.Actor.Act(source, targets, signal)
}
