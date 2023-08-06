package battlefield

type Script interface {
	Renderer
	Source() (any, Reactor)
	Add(Action)
}

type MyScript struct {
	scripter any
	reactor  Reactor
	actions  []Action
}

func NewMyScript(scripter any, reactor Reactor) *MyScript {
	return &MyScript{scripter, reactor, nil}
}

func (s *MyScript) Source() (any, Reactor) {
	return s.scripter, s.reactor
}

func (s *MyScript) Add(action Action) {
	s.actions = append(s.actions, action)
	action.SetScript(s)
}

func (s *MyScript) Render(b *BattleField) {
	for _, action := range s.actions {
		action.Render(b)
	}
}

type Action interface {
	Portfolio
	Renderer
	Script() Script
	SetScript(Script)
	Targets() []Warrior
	Verb() Verb
}

type MyAction struct {
	*MyPortfolio
	script  Script
	targets []Warrior
	verb    Verb
}

func NewMyAction(targets []Warrior, verb Verb) *MyAction {
	return &MyAction{NewMyPortfolio(), nil, targets, verb}
}

func (a *MyAction) Script() Script {
	return a.script
}

func (a *MyAction) SetScript(script Script) {
	a.script = script
}

func (a *MyAction) Targets() []Warrior {
	return a.targets
}

func (a *MyAction) Verb() Verb {
	return a.verb
}

func (a *MyAction) Render(b *BattleField) {
	b.React(NewPreActionSignal(a))

	for _, target := range a.targets {
		a.verb.Render(target, a)
	}

	b.React(NewPostActionSignal(a))
}

type Verb interface {
	Render(target Warrior, action Action)
	Forker
}

type Attack struct {
	evaluator Evaluator
	critical  bool
	loss      int
}

func NewAttack(evaluator Evaluator, critical bool) *Attack {
	return &Attack{evaluator, critical, 0}
}

func (a *Attack) Critical() bool {
	return a.critical
}

func (a *Attack) SetCritical(critical bool) {
	a.critical = critical
}

func (a *Attack) Loss() int {
	return a.loss
}

func (a *Attack) Fork(evaluator Evaluator) any {
	if evaluator == nil {
		return a
	}

	return &Attack{evaluator, a.critical, 0}
}

func (a *Attack) Render(target Warrior, action Action) {
	damage := a.evaluator.Evaluate(target)
	defense := target.Component(Defense)
	a.loss = damage - defense
	if a.loss < 0 {
		a.loss = 0
	}

	e := NewEvaluationSignal(target, Loss, a.loss)
	action.React(e, nil)
	loss := NewPreLossSignal(target, e.Value())
	target.React(loss, nil)

	a.loss = loss.Loss()
	r := target.Health()
	m := target.Component(HealthMaximum)
	c := r.Current*m/r.Maximum - a.loss
	if c < 0 {
		a.loss += c
		c = 0
	}
	target.SetHealth(Ratio{c, m})
}

type Heal struct {
	evaluator Evaluator
	rise      int
}

func NewHeal(evaluator Evaluator) *Heal {
	return &Heal{evaluator, 0}
}

func (h *Heal) Rise() int {
	return h.rise
}

func (h *Heal) Fork(evaluator Evaluator) any {
	if evaluator == nil {
		return h
	}

	return &Heal{evaluator, 0}
}

func (h *Heal) Render(target Warrior, action Action) {
	r := target.Health()
	m := target.Component(HealthMaximum)
	c := r.Current * m / r.Maximum
	h.rise = h.evaluator.Evaluate(target)

	c += h.rise
	if c > m {
		h.rise -= c - m
		c = m
	}
	target.SetHealth(Ratio{c, m})
}

type Buff struct {
	evaluator Evaluator
	reactor   ForkReactor
	perAction bool
}

func NewBuff(evaluator Evaluator, reactor ForkReactor, perAction bool) *Buff {
	return &Buff{evaluator, reactor, perAction}
}

func (b *Buff) Fork(evaluator Evaluator) any {
	if evaluator == nil {
		return b
	}

	return &Buff{evaluator, b.reactor, b.perAction}
}

func (b *Buff) Render(target Warrior, action Action) {
	reactor := b.reactor
	if b.evaluator != nil {
		e := ConstEvaluator(b.evaluator.Evaluate(target))
		reactor = reactor.Fork(e).(ForkReactor)
	}

	if b.perAction {
		action.Add(reactor)
	} else {
		target.Add(reactor)
	}
}

type Purge struct {
	rng   Rng
	tag   any
	count int
}

func NewPurge(rng Rng, tag any, count int) *Purge {
	return &Purge{rng, tag, count}
}

func (p *Purge) Fork(Evaluator) any {
	return p
}

func (p *Purge) Render(target Warrior, action Action) {
	buffs := target.Buffs(p.tag)
	m, n := len(buffs), p.count
	if m > n && n > 0 {
		for ; n > 0; n-- {
			i := int(p.rng.Float64() * float64(m))
			m--
			buffs[i], buffs[m] = buffs[m], buffs[i]
		}

		m--
		buffs = buffs[m:]
	}

	for _, buff := range buffs {
		target.Remove(buff)
	}
}
