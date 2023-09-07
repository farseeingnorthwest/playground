package battlefield

import (
	"encoding/json"
	"log/slog"

	"github.com/farseeingnorthwest/playground/battlefield/v2/functional"
)

var (
	_ Script        = (*script)(nil)
	_ Action        = (*action)(nil)
	_ Verb          = (*Attack)(nil)
	_ Verb          = (*Heal)(nil)
	_ Verb          = (*Buff)(nil)
	_ Verb          = (*Purge)(nil)
	_ ActionContext = (*actionContext)(nil)
)

type Script interface {
	Renderer
	Source() (any, Reactor)
	Add(Action)
	Loss() int
}

type script struct {
	scripter any
	reactor  Reactor
	actions  []Action
}

func newScript(scripter any, reactor Reactor) *script {
	return &script{scripter, reactor, nil}
}

func (s *script) Source() (any, Reactor) {
	return s.scripter, s.reactor
}

func (s *script) Add(action Action) {
	s.actions = append(s.actions, action)
	action.SetScript(s)
}

func (s *script) Loss() int {
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

func (s *script) Render(b *BattleField) {
	for _, action := range s.actions {
		action.Render(b)
	}
}

type Action interface {
	Script() Script
	SetScript(Script)
	Targets() []Warrior
	FalseTargets() []Warrior
	ImmuneTargets() map[Warrior]struct{}
	AddImmuneTarget(Warrior)
	Verb() Verb
	Portfolio
	Renderer
}

type action struct {
	*FatPortfolio
	script        Script
	targets       []Warrior
	falseTargets  []Warrior
	immuneTargets map[Warrior]struct{}
	verb          Verb
}

func newAction(targets []Warrior, verb Verb) *action {
	return &action{NewFatPortfolio(), nil, targets, nil, make(map[Warrior]struct{}), verb}
}

func (a *action) Script() Script {
	return a.script
}

func (a *action) SetScript(script Script) {
	a.script = script
}

func (a *action) Targets() []Warrior {
	return a.targets
}

func (a *action) ImmuneTargets() map[Warrior]struct{} {
	return a.immuneTargets
}

func (a *action) AddImmuneTarget(target Warrior) {
	a.immuneTargets[target] = struct{}{}
}

func (a *action) FalseTargets() []Warrior {
	return a.falseTargets
}

func (a *action) Verb() Verb {
	return a.verb
}

func (a *action) Render(b *BattleField) {
	b.React(NewPreActionSignal(a))

	i, j := 0, len(a.targets)
	var deaths []Warrior
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

		if a.verb.Render(target, newActionContext(a, b)) {
			if target.Health().Current <= 0 {
				deaths = append(deaths, target)
			}

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

	b.React(NewPostActionSignal(a, deaths))
}

type Verb interface {
	Name() string
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

func (*Attack) Name() string {
	return "attack"
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

func (a *Attack) MarshalJSON() ([]byte, error) {
	return json.Marshal(attack{a.evaluator, a.critical})
}

type attack struct {
	Evaluator Evaluator `json:"attack"`
	Critical  bool      `json:"critical,omitempty"`
}

type Heal struct {
	evaluator Evaluator
	rise      map[Warrior]int
}

func NewHeal(evaluator Evaluator) *Heal {
	return &Heal{evaluator, make(map[Warrior]int)}
}

func (*Heal) Name() string {
	return "heal"
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

func (h *Heal) MarshalJSON() ([]byte, error) {
	return json.Marshal(heal{h.evaluator})
}

type heal struct {
	Evaluator Evaluator `json:"heal"`
}

type Buff struct {
	capacity  bool
	evaluator Evaluator
	reactor   ForkReactor
}

func NewBuff(capacity bool, evaluator Evaluator, reactor ForkReactor) *Buff {
	return &Buff{capacity && evaluator != nil, evaluator, reactor}
}

func (*Buff) Name() string {
	return "buff"
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
		reactor.(*FatReactor).Update(FatCapacity(nil, c))
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
			slog.Any("reactor", QueryTagA[Label](functional.Second(ac.Action().Script().Source())))),
	)
	return true
}

func (b *Buff) MarshalJSON() ([]byte, error) {
	return json.Marshal(buff{b.reactor, b.evaluator, b.capacity})
}

type buff struct {
	Reactor   Reactor   `json:"buff"`
	Evaluator Evaluator `json:"evaluator,omitempty"`
	Capacity  bool      `json:"capacity,omitempty"`
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

func (*Purge) Name() string {
	return "purge"
}

func (p *Purge) Reactors() []Reactor {
	return p.reactors
}

func (p *Purge) Fork(Evaluator) any {
	return p
}

func (p *Purge) Render(target Warrior, _ ActionContext) bool {
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

func (p *Purge) MarshalJSON() ([]byte, error) {
	return json.Marshal(purge{p.tag, p.count})
}

type purge struct {
	Tag   any `json:"purge"`
	Count int `json:"count,omitempty"`
}

type ActionContext interface {
	EvaluationContext
	Action() Action
}

type actionContext struct {
	EvaluationContext
	action Action
}

func newActionContext(action Action, ac EvaluationContext) *actionContext {
	return &actionContext{ac, action}
}

func (c *actionContext) Action() Action {
	return c.action
}
