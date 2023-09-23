package battlefield

import (
	"encoding/json"
	"errors"
	"log/slog"
	"reflect"
)

var (
	_ Script        = (*script)(nil)
	_ Action        = (*action)(nil)
	_ Verb          = (*Attack)(nil)
	_ Verb          = (*Heal)(nil)
	_ Verb          = (*Buff)(nil)
	_ Verb          = (*Purge)(nil)
	_ ActionContext = (*actionContext)(nil)

	ErrBadVerb = errors.New("bad verb")
)

type Script interface {
	Renderer
	Source() (Signal, any, Reactor)
	Add(Action)
	Loss() int
}

type script struct {
	signal   Signal
	scripter any
	reactor  Reactor
	actions  []Action
}

func newScript(signal Signal, reactor Reactor) *script {
	return &script{signal, signal.Current(), reactor, nil}
}

func (s *script) Source() (Signal, any, Reactor) {
	return s.signal, s.scripter, s.reactor
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
	ID() int
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
	id            int
	script        Script
	targets       []Warrior
	falseTargets  []Warrior
	immuneTargets map[Warrior]struct{}
	verb          Verb
	*FatPortfolio
}

func newAction(id int, targets []Warrior, verb Verb) *action {
	return &action{id, nil, targets, nil, make(map[Warrior]struct{}), verb, NewFatPortfolio()}
}

func (a *action) ID() int {
	return a.id
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
	b.React(NewPreActionSignal(b.Next(), a))

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

	b.React(NewPostActionSignal(b.Next(), a, deaths))
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

	e := NewEvaluationSignal(ac.Next(), target, Loss, t)
	ac.Action().React(e, ac)
	target.React(e, ac)
	loss := NewPreLossSignal(ac.Next(), target, ac.Action(), e.Value())
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
	_, source, reactor := ac.Action().Script().Source()
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
	return json.Marshal(struct {
		Verb      string    `json:"_verb"`
		Evaluator Evaluator `json:"evaluator,omitempty"`
		Critical  bool      `json:"critical,omitempty"`
	}{
		"attack",
		a.evaluator,
		a.critical,
	})
}

func (a *Attack) UnmarshalJSON(data []byte) error {
	var attack struct {
		Evaluator EvaluatorFile
		Critical  bool
	}
	if err := json.Unmarshal(data, &attack); err != nil {
		return err
	}

	*a = Attack{
		attack.Evaluator.Evaluator,
		attack.Critical,
		make(map[Warrior]int),
	}
	return nil
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
	_, source, reactor := ac.Action().Script().Source()
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
	return json.Marshal(struct {
		Verb      string    `json:"_verb"`
		Evaluator Evaluator `json:"evaluator,omitempty"`
	}{
		"heal",
		h.evaluator,
	})
}

func (h *Heal) UnmarshalJSON(data []byte) error {
	var heal struct {
		Evaluator EvaluatorFile
	}
	if err := json.Unmarshal(data, &heal); err != nil {
		return err
	}

	*h = Heal{
		heal.Evaluator.Evaluator,
		make(map[Warrior]int),
	}
	return nil
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
		signal, _, _ := ac.Action().Script().Source()
		ac.React(NewLifecycleSignal(ac.Next(), signal, target, overflow, nil, LifecycleOverflow))
	}
	if lm, ok := QueryTag[StackingLimit](reactor); ok {
		logger = logger.With("stacking", target.Stacking(lm).Count())
	}
	_, _, source := ac.Action().Script().Source()
	logger.Debug(
		"render",
		slog.String("verb", "buff"),
		slog.Any("reactor", QueryTagA[Label](reactor)),
		slog.Group("target",
			slog.Any("side", target.Side()),
			slog.Int("position", target.Position()),
		),
		slog.Group("source",
			slog.Any("reactor", QueryTagA[Label](source))),
	)
	return true
}

func (b *Buff) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Verb      string    `json:"_verb"`
		Reactor   Reactor   `json:"reactor"`
		Evaluator Evaluator `json:"evaluator,omitempty"`
		Capacity  bool      `json:"capacity,omitempty"`
	}{
		"buff",
		b.reactor,
		b.evaluator,
		b.capacity,
	})
}

func (b *Buff) UnmarshalJSON(data []byte) error {
	var buff struct {
		Reactor   FatReactorFile
		Evaluator EvaluatorFile
		Capacity  bool
	}
	if err := json.Unmarshal(data, &buff); err != nil {
		return err
	}

	*b = Buff{
		buff.Capacity,
		buff.Evaluator.Evaluator,
		buff.Reactor.FatReactor,
	}
	return nil
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
	return json.Marshal(struct {
		Verb  string `json:"_verb"`
		Tag   any    `json:"tag"`
		Count int    `json:"count,omitempty"`
	}{
		"purge",
		p.tag,
		p.count,
	})
}

func (p *Purge) UnmarshalJSON(data []byte) error {
	var purge struct {
		Tag   TagFile
		Count int
	}
	if err := json.Unmarshal(data, &purge); err != nil {
		return err
	}

	*p = Purge{
		DefaultRng,
		purge.Tag.Tag,
		purge.Count,
		nil,
	}
	return nil
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

type VerbFile struct {
	Verb Verb
}

func (f *VerbFile) UnmarshalJSON(data []byte) error {
	var v struct {
		Verb string `json:"_verb"`
	}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	if verb, ok := verbType[v.Verb]; ok {
		f.Verb = reflect.New(verb).Interface().(Verb)
		return json.Unmarshal(data, f.Verb)
	}

	return ErrBadVerb
}

var verbType = map[string]reflect.Type{
	"attack": reflect.TypeOf(Attack{}),
	"heal":   reflect.TypeOf(Heal{}),
	"buff":   reflect.TypeOf(Buff{}),
	"purge":  reflect.TypeOf(Purge{}),
}
