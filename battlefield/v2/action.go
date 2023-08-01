package battlefield

import "github.com/farseeingnorthwest/playground/battlefield/v2/evaluation"

type Action struct {
	Source  *Fighter
	Targets []*Fighter
	Verb
}

func (a *Action) Render(f *BattleField) {
	pre := NewPreActionSignal(a)
	f.React(pre)
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
	Evaluator
	chain *evaluation.Block
}

func NewAttackProto(e Evaluator) *Attack {
	return &Attack{e, nil}
}

func (a *Attack) Render(target, source *Warrior, action *Action) {
	damage := NewEvaluationSignal(Damage, a.Evaluate(target, a.chain), action)
	source.React(damage)
	defense := NewEvaluationSignal(Defense, target.Defense(), action)
	target.React(defense)

	loss := NewEvaluationSignal(Loss, damage.Value()-defense.Value(), action)
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
	return &Attack{a.Evaluator, chain}
}

type Heal struct {
	Evaluator
	chain *evaluation.Block
}

func NewHealProto(e Evaluator) *Heal {
	return &Heal{e, nil}
}

func (h *Heal) Render(target, _ *Warrior, action *Action) {
	heal := NewEvaluationSignal(Healing, h.Evaluate(target, h.chain), action)
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
	return &Heal{h.Evaluator, chain}
}

type Buff struct {
	reactor Reactor
	*EvalChain
}

func NewBuffProto(reactor Reactor, e Evaluator) *Buff {
	return &Buff{
		reactor,
		NewEvalChainProto(e),
	}
}

func (h *Buff) Render(target, _ *Warrior, _ *Action) {
	target.Append(h.reactor.Fork(h.ForkWith(target), nil))
}

func (h *Buff) Fork(chain *evaluation.Block, signal Signal) Verb {
	return &Buff{
		h.reactor.Fork(nil, signal),
		NewEvalChain(h.Evaluator, chain),
	}
}

type Purging struct {
}

func NewPurging() *Purging {
	return &Purging{}
}

func (p *Purging) Render(target, _ *Warrior, _ *Action) {
	// TODO:
}
