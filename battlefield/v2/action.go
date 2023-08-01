package battlefield

import "github.com/farseeingnorthwest/playground/battlefield/v2/evaluation"

type Action struct {
	Source  *Fighter
	Targets []*Fighter
	Verb
}

func (a *Action) Render(f *BattleField) {
	pre := NewPreActionSignal(a)
	f.React(pre) // TODO: render actions per warrior
	for _, action := range pre.Actions() {
		action.Render(f)
	}

	for _, object := range a.Targets {
		a.Verb.Render(object.Warrior, a.Source.Warrior, a)
	}

	post := NewPostActionSignal(a)
	f.React(post)
	for _, action := range post.Actions() {
		action.Render(f)
	}
}

type Verb interface {
	Render(target, source *Warrior, action *Action)
	Fork(*evaluation.Block, Signal) Verb
}

type Attack struct {
	*evaluation.Bundle
}

func NewAttackProto(e evaluation.Evaluator) *Attack {
	return &Attack{evaluation.NewBundleProto(e)}
}

func (a *Attack) Render(target, source *Warrior, action *Action) {
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
	if c < 0 {
		c = 0
	}

	target.current = Ratio{c, m}
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

func (h *Heal) Render(target, _ *Warrior, action *Action) {
	heal := NewEvaluationSignal(evaluation.Healing, h.Evaluate(target), action)
	target.React(heal)
	if heal.Value() < 0 {
		heal.SetValue(0)
	}

	r, m := target.Health()
	c := r.Current * m / r.Maximum
	c += heal.Value()
	if c > m {
		c = m
	}

	target.current = Ratio{c, m}
}

func (h *Heal) Fork(chain *evaluation.Block, _ Signal) Verb {
	return &Heal{h.Bundle.Fork(chain)}
}

type Buff struct {
	reactor Reactor
	*evaluation.Bundle
}

func NewBuffProto(reactor Reactor, e evaluation.Evaluator) *Buff {
	return &Buff{
		reactor,
		evaluation.NewBundleProto(e),
	}
}

func (b *Buff) Render(target, _ *Warrior, _ *Action) {
	target.Append(b.reactor.Fork(b.ForkWith(target), nil))
}

func (b *Buff) Fork(chain *evaluation.Block, signal Signal) Verb {
	return &Buff{
		b.reactor.Fork(nil, signal),
		b.Bundle.Fork(chain),
	}
}

type Purge struct {
}

func NewPurgingProto() *Purge {
	return &Purge{}
}

func (*Purge) Render(target, _ *Warrior, _ *Action) {
}

func (*Purge) Fork(_ *evaluation.Block, _ Signal) Verb {
	return &Purge{}
}
