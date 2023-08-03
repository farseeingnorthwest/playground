package battlefield

import (
	"github.com/farseeingnorthwest/playground/battlefield/v2/evaluation"
	"github.com/farseeingnorthwest/playground/battlefield/v2/mod"
)

type RawReactor interface {
	React(Signal)
}

type Reactor interface {
	RawReactor
	mod.Tagger
}

type Forker interface {
	Fork(*evaluation.Block, Signal) any
}

type ForkableReactor interface {
	Reactor
	Forker
}

type ModifiedReactor struct {
	mod.TaggerMod
	*mod.FiniteMod
	mod.PeriodicMod
	actors []Actor
}

func NewModifiedReactor(tag any, actors []Actor, options ...func(*ModifiedReactor)) *ModifiedReactor {
	m := &ModifiedReactor{
		actors: actors,
	}

	m.SetTag(tag)
	for _, option := range options {
		option(m)
	}

	return m
}

func Capacity(capacity int) func(*ModifiedReactor) {
	return func(m *ModifiedReactor) {
		m.FiniteMod = mod.NewFiniteModifier(capacity)
	}
}

func Period(period int) func(*ModifiedReactor) {
	return func(m *ModifiedReactor) {
		m.PeriodicMod.SetPeriod(period)
	}
}

func Phase(phase int) func(*ModifiedReactor) {
	return func(m *ModifiedReactor) {
		m.PeriodicMod.SetPhase(phase)
	}
}

func (m *ModifiedReactor) Fork(block *evaluation.Block, signal Signal) *ModifiedReactor {
	actors := make([]Actor, len(m.actors))
	for i, actor := range m.actors {
		actors[i] = actor.Fork(block, signal).(Actor)
	}

	return &ModifiedReactor{
		TaggerMod:   m.TaggerMod,
		FiniteMod:   m.FiniteMod.Clone().(*mod.FiniteMod),
		PeriodicMod: m.PeriodicMod,
		actors:      actors,
	}
}

func (m *ModifiedReactor) WarmUp() {
	m.FiniteMod.WarmUp()
	m.PeriodicMod.WarmUp()
}

func (m *ModifiedReactor) act(source *Fighter, targets []*Fighter, signal Signal) (actions []*Action) {
	if !m.Free() {
		return
	}

	for _, actor := range m.actors {
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
