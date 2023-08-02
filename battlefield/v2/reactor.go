package battlefield

import (
	"github.com/farseeingnorthwest/playground/battlefield/v2/evaluation"
	"github.com/farseeingnorthwest/playground/battlefield/v2/modifier"
)

type Reactor interface {
	React(Signal)
}

type Forker interface {
	Fork(*evaluation.Block, Signal) any
}

type ForkableReactor interface {
	Reactor
	Forker
}

type ModifiedReactor struct {
	*modifier.FiniteModifier
	*modifier.PeriodicModifier
	actors []Actor
}

func NewModifiedReactor(actors []Actor, options ...func(*ModifiedReactor)) *ModifiedReactor {
	a := &ModifiedReactor{
		PeriodicModifier: &modifier.PeriodicModifier{},
		actors:           actors,
	}
	for _, option := range options {
		option(a)
	}

	return a
}

func Capacity(capacity int) func(*ModifiedReactor) {
	return func(a *ModifiedReactor) {
		a.FiniteModifier = modifier.NewFiniteModifier(capacity)
	}
}

func Period(period int) func(*ModifiedReactor) {
	return func(a *ModifiedReactor) {
		a.PeriodicModifier.SetPeriod(period)
	}
}

func Phase(phase int) func(*ModifiedReactor) {
	return func(a *ModifiedReactor) {
		a.PeriodicModifier.SetPhase(phase)
	}
}

func (a *ModifiedReactor) Fork(block *evaluation.Block, signal Signal) *ModifiedReactor {
	actors := make([]Actor, len(a.actors))
	for i, actor := range a.actors {
		actors[i] = actor.Fork(block, signal).(Actor)
	}

	return &ModifiedReactor{
		FiniteModifier:   a.FiniteModifier.Clone().(*modifier.FiniteModifier),
		PeriodicModifier: a.PeriodicModifier.Clone().(*modifier.PeriodicModifier),
		actors:           actors,
	}
}

func (a *ModifiedReactor) WarmUp() {
	a.FiniteModifier.WarmUp()
	a.PeriodicModifier.WarmUp()
}

func (a *ModifiedReactor) act(source *Fighter, targets []*Fighter, signal Signal) (actions []*Action) {
	if !a.Free() {
		return
	}

	for _, actor := range a.actors {
		a := actor.Act(source, targets, signal)
		if a == nil {
			return nil
		}

		actions = append(actions, a)
	}

	return
}

type LaunchReactor struct {
	*ModifiedReactor
}

func (a *LaunchReactor) React(signal Signal) {
	switch sig := signal.(type) {
	case *LaunchSignal:
		actions := a.act(sig.Target, sig.Field.fighters, sig)
		if actions == nil {
			return
		}

		sig.Append(actions...)
		sig.Launched = true
		a.WarmUp()

	case *RoundEndSignal:
		a.CoolDown()
	}
}

func (a *LaunchReactor) Fork(block *evaluation.Block, signal Signal) Reactor {
	return &LaunchReactor{a.ModifiedReactor.Fork(block, signal)}
}

type RoundStartReactor struct {
	*ModifiedReactor
}

func (a *RoundStartReactor) React(signal Signal) {
	switch sig := signal.(type) {
	case *RoundStartSignal:
		actions := a.act(sig.current, sig.Field.fighters, sig)
		if actions == nil {
			return
		}

		sig.Append(actions...)
		a.WarmUp()
		a.CoolDown()
	}
}

func (a *RoundStartReactor) Fork(block *evaluation.Block, signal Signal) any {
	return &RoundStartReactor{a.ModifiedReactor.Fork(block, signal)}
}

type PreAttackReactor struct {
	*ModifiedReactor
}

func (c *PreAttackReactor) React(signal Signal) {
	switch sig := signal.(type) {
	case *PreActionSignal:
		_, ok := sig.Verb.(*Attack)
		if !ok {
			return
		}
		if sig.Current() != sig.Source {
			return
		}

		actions := c.act(sig.Source, sig.Targets, sig)
		if actions == nil {
			return
		}

		sig.Append(actions...)

	case *RoundEndSignal:
		c.WarmUp()
		c.CoolDown()
	}
}

func (c *PreAttackReactor) Fork(block *evaluation.Block, signal Signal) any {
	return &PreAttackReactor{c.ModifiedReactor.Fork(block, signal)}
}
