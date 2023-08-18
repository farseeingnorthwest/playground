package battlefield

import "log/slog"

type Script interface {
	Renderer
	Source() (any, Reactor)
	Add(Action)
	Loss() int
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

func (s *MyScript) Loss() int {
	loss := 0
	for _, action := range s.actions {
		if a, ok := action.Verb().(*Attack); ok {
			for _, lo := range a.Loss() {
				loss += lo
			}
		}
	}

	return loss
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
	ImmuneTargets() map[Warrior]struct{}
	FalseTargets() []Warrior
	AddImmuneTarget(Warrior)
	Verb() Verb
}

type MyAction struct {
	*FatPortfolio
	script        Script
	targets       []Warrior
	immuneTargets map[Warrior]struct{}
	falseTargets  []Warrior
	verb          Verb
}

func NewMyAction(targets []Warrior, verb Verb) *MyAction {
	return &MyAction{NewFatPortfolio(), nil, targets, make(map[Warrior]struct{}), nil, verb}
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

func (a *MyAction) ImmuneTargets() map[Warrior]struct{} {
	return a.immuneTargets
}

func (a *MyAction) AddImmuneTarget(target Warrior) {
	a.immuneTargets[target] = struct{}{}
}

func (a *MyAction) FalseTargets() []Warrior {
	return a.falseTargets
}

func (a *MyAction) Verb() Verb {
	return a.verb
}

func (a *MyAction) Render(b *BattleField) {
	b.React(NewPreActionSignal(a))

	i, j := 0, len(a.targets)
	for i < j {
		target := a.targets[i]
		if _, ok := a.immuneTargets[target]; ok {
			slog.Debug(
				"render",
				slog.Group("immune",
					slog.Any("side", target.(Warrior).Side()),
					slog.Int("position", target.(Warrior).Position())),
			)
			i++
			continue
		}

		if a.verb.Render(target, NewMyActionContext(a, b)) {
			i++
			continue
		}
		if j--; i < j {
			a.targets[i], a.targets[j] = a.targets[j], a.targets[i]
		}
	}

	a.falseTargets = a.targets[j:]
	a.targets = a.targets[:j]
	if len(a.falseTargets) > 0 || len(a.immuneTargets) > 0 {
		slog.Debug("render",
			slog.Int("targets", len(a.targets)),
			slog.Int("falseTargets", len(a.falseTargets)),
			slog.Int("immuneTargets", len(a.immuneTargets)),
		)
	}

	b.React(NewPostActionSignal(a))
}

type MyActionContext struct {
	EvaluationContext
	action Action
}

func NewMyActionContext(action Action, ac EvaluationContext) *MyActionContext {
	return &MyActionContext{ac, action}
}

func (ac *MyActionContext) Action() Action {
	return ac.action
}

type Verb interface {
	Render(target Warrior, ac ActionContext) bool
	Forker
}

type Attack struct {
	evaluator Evaluator
	critical  bool
	loss      map[Warrior]int
}

func NewAttack(evaluator Evaluator, critical bool) *Attack {
	return &Attack{evaluator, critical, make(map[Warrior]int)}
}

func (a *Attack) Critical() bool {
	return a.critical
}

func (a *Attack) SetCritical(critical bool) {
	a.critical = critical
}

func (a *Attack) Loss() map[Warrior]int {
	return a.loss
}

func (a *Attack) Fork(evaluator Evaluator) any {
	if evaluator == nil {
		evaluator = a.evaluator
	}

	return &Attack{evaluator, a.critical, make(map[Warrior]int)}
}

func (a *Attack) Render(target Warrior, ac ActionContext) bool {
	if target.Health().Current <= 0 {
		return false
	}

	damage := a.evaluator.Evaluate(target, ac)
	defense := target.Component(Defense, ac)
	t := damage - defense
	if t < 0 {
		t = 0
	}

	e := NewEvaluationSignal(target, Loss, t)
	ac.Action().React(e, ac)
	target.React(e, ac)
	loss := NewPreLossSignal(target, e.Value())
	target.React(loss, ac)

	r := target.Health()
	m := target.Component(HealthMaximum, ac)
	c := r.Current*m/r.Maximum - loss.Loss()
	overflow := 0
	if c < 0 {
		overflow = -c
		c = 0
	}
	target.SetHealth(Ratio{c, m})
	a.loss[target] = loss.Loss() - overflow
	source, reactor := ac.Action().Script().Source()
	slog.Debug(
		"render",
		slog.String("verb", "attack"),
		slog.Bool("critical", a.critical),
		slog.Int("loss", loss.Loss()),
		slog.Int("overflow", overflow),
		slog.Group("source",
			slog.Any("side", source.(Warrior).Side()),
			slog.Int("position", source.(Warrior).Position()),
			slog.Any("reactor", QueryTagA[Label](reactor)),
			slog.Int("damage", damage)),
		slog.Group("target",
			slog.Any("side", target.Side()),
			slog.Int("position", target.Position()),
			slog.Int("defense", defense),
			slog.Group("health",
				slog.Int("current", c),
				slog.Int("maximum", m))),
	)

	return true
}

type Heal struct {
	evaluator Evaluator
	rise      map[Warrior]int
}

func NewHeal(evaluator Evaluator) *Heal {
	return &Heal{evaluator, make(map[Warrior]int)}
}

func (h *Heal) Rise() map[Warrior]int {
	return h.rise
}

func (h *Heal) Fork(evaluator Evaluator) any {
	if evaluator == nil {
		evaluator = h.evaluator
	}

	return &Heal{evaluator, make(map[Warrior]int)}
}

func (h *Heal) Render(target Warrior, ac ActionContext) bool {
	if target.Health().Current <= 0 {
		return false
	}

	r := target.Health()
	m := target.Component(HealthMaximum, ac)
	c := r.Current * m / r.Maximum
	rise := h.evaluator.Evaluate(target, ac)

	c += rise
	overflow := 0
	if c > m {
		overflow = c - m
		c = m
	}
	target.SetHealth(Ratio{c, m})
	h.rise[target] = rise - overflow
	source, reactor := ac.Action().Script().Source()
	slog.Debug(
		"render",
		slog.String("verb", "heal"),
		slog.Int("rise", rise),
		slog.Int("overflow", overflow),
		slog.Group("source",
			slog.Any("side", source.(Warrior).Side()),
			slog.Int("position", source.(Warrior).Position()),
			slog.Any("reactor", QueryTagA[Label](reactor))),
		slog.Group("target",
			slog.Any("side", target.Side()),
			slog.Int("position", target.Position()),
			slog.Group("health",
				slog.Int("current", c),
				slog.Int("maximum", m))),
	)
	return true
}

type Buff struct {
	capacity  bool
	evaluator Evaluator
	reactor   ForkReactor
}

func NewBuff(capacity bool, evaluator Evaluator, reactor ForkReactor) *Buff {
	return &Buff{capacity && evaluator != nil, evaluator, reactor}
}

func (b *Buff) Reactor() ForkReactor {
	return b.reactor
}

func (b *Buff) Fork(evaluator Evaluator) any {
	if evaluator == nil {
		return b
	}

	return &Buff{b.capacity, evaluator, b.reactor}
}

func (b *Buff) Render(target Warrior, ac ActionContext) bool {
	logger := slog.With()
	if target.Health().Current <= 0 {
		return false
	}

	var reactor Reactor
	e := b.evaluator
	if e != nil {
		e = ConstEvaluator(e.Evaluate(target, ac))
	}
	if !b.capacity {
		reactor = b.reactor.Fork(e).(Reactor)
	} else {
		reactor = b.reactor.Fork(nil).(ForkReactor)

		c := int(e.(ConstEvaluator))
		reactor.(*FatReactor).Amend(FatCapacity(nil, c))
		logger = logger.With(slog.Int("capacity", c))
	}

	if overflow := target.Add(reactor); overflow != nil {
		slog.Debug(
			"render",
			slog.String("verb", "buff/overflow"),
			slog.Any("reactor", QueryTagA[Label](reactor)),
			slog.Group("target",
				slog.Any("side", target.Side()),
				slog.Int("position", target.Position()),
			),
		)
		ac.React(NewLifecycleSignal(target, overflow, nil))
	}
	if stacking, ok := QueryTag[StackingLimit](reactor); ok {
		logger = logger.With("stacking", stacking.Count())
	}
	logger.Debug(
		"render",
		slog.String("verb", "buff"),
		slog.Any("reactor", QueryTagA[Label](reactor)),
		slog.Group("target",
			slog.Any("side", target.Side()),
			slog.Int("position", target.Position()),
		),
		slog.Group("source",
			slog.Any("reactor", QueryTagA[Label](Second(ac.Action().Script().Source())))),
	)
	return true
}

type Purge struct {
	rng      Rng
	tag      any
	count    int
	reactors []Reactor
}

func NewPurge(rng Rng, tag any, count int) *Purge {
	return &Purge{rng, tag, count, nil}
}

func (p *Purge) Reactors() []Reactor {
	return p.reactors
}

func (p *Purge) Fork(Evaluator) any {
	return p
}

func (p *Purge) Render(target Warrior, ac ActionContext) bool {
	if target.Health().Current <= 0 {
		return false
	}

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

	tags := make([]any, len(buffs))
	for i, buff := range buffs {
		target.Remove(buff)
		tags[i] = QueryTagA[Label](buff)
	}
	p.reactors = buffs
	slog.Debug("render", slog.String("verb", "purge"), slog.Any("reactors", tags))
	return true
}
