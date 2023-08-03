package battlefield

import "github.com/farseeingnorthwest/playground/battlefield/v2/evaluation"

type Actor interface {
	Act(*Fighter, []*Fighter, Signal) *Action
	Forker
}

type BlindActor struct {
	proto Verb
	*evaluation.Bundle
}

func (a *BlindActor) Act(source *Fighter, targets []*Fighter, signal Signal) *Action {
	return &Action{
		Targets: targets,
		Verb:    a.proto.Fork(a.ForkWith(source.Warrior), signal),
	}
}

func (a *BlindActor) Fork(block *evaluation.Block, _ Signal) any {
	return &BlindActor{
		proto:  a.proto,
		Bundle: a.Bundle.Fork(block),
	}
}

type SelectiveActor struct {
	Selector
	Actor
}

func (a *SelectiveActor) Act(source *Fighter, targets []*Fighter, signal Signal) *Action {
	selected := a.Select(source, targets)
	if len(selected) == 0 {
		return nil
	}

	return a.Actor.Act(source, selected, signal)
}

func (a *SelectiveActor) Fork(block *evaluation.Block, signal Signal) any {
	return &SelectiveActor{
		Selector: a.Selector,
		Actor:    a.Actor.Fork(block, signal).(Actor),
	}
}

type Rng interface {
	Gen() float64 // [0, 1)
}

type ProbabilityActor struct {
	rng  Rng
	odds int
	Actor
}

func (a *ProbabilityActor) Act(source *Fighter, targets []*Fighter, signal Signal) *Action {
	if a.rng.Gen() > float64(a.odds)/100 {
		return nil
	}

	return a.Actor.Act(source, targets, signal)
}

func (a *ProbabilityActor) Fork(block *evaluation.Block, signal Signal) any {
	return &ProbabilityActor{
		rng:   a.rng,
		odds:  a.odds,
		Actor: a.Actor.Fork(block, signal).(Actor),
	}
}
