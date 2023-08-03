package battlefield

import "github.com/farseeingnorthwest/playground/battlefield/v2/evaluation"

type Script struct {
	Current *Fighter
	Source  Reactor
	Actions []*Action
}

func NewScript(current *Fighter, source Reactor, actions ...*Action) *Script {
	s := &Script{current, source, actions}
	for _, a := range actions {
		a.Script = s
	}

	return s
}

func (s *Script) Render(f *BattleField) {
	f.React(&PreScriptSignal{Script: s})
	for _, action := range s.Actions {
		action.Render(f)
	}
	f.React(&PostScriptSignal{Script: s})
}

type Action struct {
	Script  *Script
	Targets []*Fighter
	Verb
	Interests []Interest
}

func (a *Action) Render(f *BattleField) {
	pre := NewPreActionSignal(a)
	f.React(pre)

	var current *Warrior
	if a.Script.Current != nil {
		current = a.Script.Current.Warrior
	}
	for _, object := range a.Targets {
		a.Interests = append(a.Interests, a.Verb.Render(object.Warrior, current, a))
	}

	post := NewPostActionSignal(a)
	f.React(post)
}

type Verb interface {
	Render(target, source *Warrior, action *Action) Interest
	Fork(*evaluation.Block, Signal) Verb
}

type Attack struct {
	*evaluation.Bundle
}

func NewAttackProto(e evaluation.Evaluator) *Attack {
	return &Attack{evaluation.NewBundleProto(e)}
}

func (a *Attack) Render(target, source *Warrior, action *Action) Interest {
	damage := NewEvaluationSignal(evaluation.Damage, a.Evaluate(target), action)
	source.React(damage)
	defense := NewEvaluationSignal(evaluation.Defense, target.Defense(), action)
	target.React(defense)

	loss := NewEvaluationSignal(evaluation.Loss, damage.Value()-defense.Value(), action)
	target.React(loss)
	if loss.Value() < 0 {
		loss.SetValue(0)
	}

	r, m := target.Health()
	c := r.Current * m / r.Maximum
	c -= loss.Value()
	overflow := 0
	if c < 0 {
		overflow = -c
		c = 0
	}
	target.SetHealth(Ratio{c, m})

	return &AttackInterest{
		interest: interest{target},
		Damage:   damage.Value(),
		Defense:  defense.Value(),
		Loss:     loss.Value(),
		Overflow: overflow,
		Health: Health{
			Current: c,
			Maximum: m,
		},
	}
}

func (a *Attack) Fork(chain *evaluation.Block, _ Signal) Verb {
	return &Attack{a.Bundle.Fork(chain)}
}

type Heal struct {
	*evaluation.Bundle
}

func NewHealProto(e evaluation.Evaluator) *Heal {
	return &Heal{evaluation.NewBundleProto(e)}
}

func (h *Heal) Render(target, _ *Warrior, action *Action) Interest {
	healing := NewEvaluationSignal(evaluation.Healing, h.Evaluate(target), action)
	target.React(healing)
	if healing.Value() < 0 {
		healing.SetValue(0)
	}

	r, m := target.Health()
	c := r.Current * m / r.Maximum
	c += healing.Value()
	overflow := 0
	if c > m {
		overflow = c - m
		c = m
	}
	target.SetHealth(Ratio{c, m})

	return &HealingInterest{
		interest: interest{target},
		Healing:  healing.Value(),
		Overflow: overflow,
		Health: Health{
			Current: c,
			Maximum: m,
		},
	}
}

func (h *Heal) Fork(chain *evaluation.Block, _ Signal) Verb {
	return &Heal{h.Bundle.Fork(chain)}
}

type Buff struct {
	reactor ForkableReactor
	*evaluation.Bundle
}

func NewBuffProto(reactor ForkableReactor, e evaluation.Evaluator) *Buff {
	return &Buff{
		reactor,
		evaluation.NewBundleProto(e),
	}
}

func (b *Buff) Render(target, _ *Warrior, _ *Action) Interest {
	reactor := b.reactor.Fork(b.ForkWith(target), nil).(Reactor)
	target.Append(reactor)

	return &BuffingInterest{
		interest: interest{target},
		Buff:     reactor,
	}
}

func (b *Buff) Fork(chain *evaluation.Block, signal Signal) Verb {
	return &Buff{
		b.reactor.Fork(nil, signal).(ForkableReactor),
		b.Bundle.Fork(chain),
	}
}

type Purge struct {
}

func NewPurgingProto() *Purge {
	return &Purge{}
}

func (*Purge) Render(target, _ *Warrior, _ *Action) Interest {
	return &PurgingInterest{
		interest: interest{target},
	}
}

func (*Purge) Fork(_ *evaluation.Block, _ Signal) Verb {
	return &Purge{}
}
