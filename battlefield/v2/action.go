package battlefield

import "log/slog"

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
	*FatPortfolio
	script  Script
	targets []Warrior
	verb    Verb
}

func NewMyAction(targets []Warrior, verb Verb) *MyAction {
	return &MyAction{NewFatPortfolio(), nil, targets, verb}
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
		a.verb.Render(target, a, b)
	}

	b.React(NewPostActionSignal(a))
}

type Verb interface {
	Render(target Warrior, action Action, ec EvaluationContext)
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

func (a *Attack) Render(target Warrior, action Action, ec EvaluationContext) {
	damage := a.evaluator.Evaluate(target, ec)
	defense := target.Component(Defense, ec)
	t := damage - defense
	if t < 0 {
		t = 0
	}

	e := NewEvaluationSignal(target, Loss, t)
	action.React(e, ec)
	target.React(e, ec)
	loss := NewPreLossSignal(target, e.Value())
	target.React(loss, ec)

	r := target.Health()
	m := target.Component(HealthMaximum, ec)
	c := r.Current*m/r.Maximum - loss.Loss()
	overflow := 0
	if c < 0 {
		overflow = -c
		c = 0
	}
	target.SetHealth(Ratio{c, m})
	a.loss = loss.Loss() - overflow
	source, _ := action.Script().Source()
	slog.Debug(
		"render",
		slog.String("verb", "attack"),
		slog.Bool("critical", a.critical),
		slog.Int("loss", loss.Loss()),
		slog.Int("overflow", overflow),
		slog.Group("source",
			slog.Any("side", source.(Warrior).Side()),
			slog.Int("position", source.(Warrior).Position()),
			slog.Any("reactor", QueryTagA[Label](Second(action.Script().Source()))),
			slog.Int("damage", damage)),
		slog.Group("target",
			slog.Any("side", target.Side()),
			slog.Int("position", target.Position()),
			slog.Int("defense", defense),
			slog.Group("health",
				slog.Int("current", c),
				slog.Int("maximum", m))),
	)

	ec.React(NewLossSignal(target, action))
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

func (h *Heal) Render(target Warrior, action Action, ec EvaluationContext) {
	r := target.Health()
	m := target.Component(HealthMaximum, ec)
	c := r.Current * m / r.Maximum
	rise := h.evaluator.Evaluate(target, ec)

	c += rise
	overflow := 0
	if c > m {
		overflow = c - m
		c = m
	}
	target.SetHealth(Ratio{c, m})
	h.rise = rise - overflow
	source, _ := action.Script().Source()
	slog.Debug(
		"render",
		slog.String("verb", "heal"),
		slog.Int("rise", rise),
		slog.Int("overflow", overflow),
		slog.Group("source",
			slog.Any("side", source.(Warrior).Side()),
			slog.Int("position", source.(Warrior).Position()),
			slog.Any("reactor", QueryTagA[Label](Second(action.Script().Source())))),
		slog.Group("target",
			slog.Any("side", target.Side()),
			slog.Int("position", target.Position()),
			slog.Group("health",
				slog.Int("current", c),
				slog.Int("maximum", m))),
	)
}

type Buff struct {
	evaluator Evaluator
	reactor   ForkReactor
}

func NewBuff(evaluator Evaluator, reactor ForkReactor) *Buff {
	return &Buff{evaluator, reactor}
}

func (b *Buff) Fork(evaluator Evaluator) any {
	if evaluator == nil {
		return b
	}

	return &Buff{evaluator, b.reactor}
}

func (b *Buff) Render(target Warrior, action Action, ec EvaluationContext) {
	reactor := b.reactor
	if b.evaluator != nil {
		e := ConstEvaluator(b.evaluator.Evaluate(target, ec))
		reactor = reactor.Fork(e).(ForkReactor)
	}

	target.Add(reactor)
	slog.Debug("render",
		slog.String("verb", "buff"),
		slog.Any("reactor", QueryTagA[Label](reactor)),
		slog.Group("target",
			slog.Any("side", target.Side()),
			slog.Int("position", target.Position()),
		),
		slog.Group("source",
			slog.Any("reactor", QueryTagA[Label](Second(action.Script().Source())))),
	)
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

func (p *Purge) Render(target Warrior, _ Action, _ EvaluationContext) {
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
