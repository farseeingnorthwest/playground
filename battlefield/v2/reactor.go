package battlefield

import (
	"encoding/json"
	"errors"
	"reflect"

	"github.com/farseeingnorthwest/playground/battlefield/v2/functional"
)

var (
	_ Actor = Buffer{}
	_ Actor = VerbActor{}
	_ Actor = SelectActor{}
	_ Actor = ProbabilityActor{}
	_ Actor = SequenceActor{}
	_ Actor = RepeatActor{}
	_ Actor = CriticalActor{}
	_ Actor = ImmuneActor{}
	_ Actor = LossStopper{}
	_ Actor = LossResister{}
	_ Actor = TheoryActor{}
	_ Actor = (*ActionBuffer)(nil)

	ErrBadActor = errors.New("bad actor")
)

type Reactor interface {
	React(Signal, EvaluationContext)
	Active() bool
}

type Forker interface {
	Fork(Evaluator) any
}

type ForkReactor interface {
	Reactor
	Forker
}

type Actor interface {
	Act(Signal, []Warrior, ActorContext)
	Forker
}

type Buffer struct {
	axis      Axis
	bias      bool
	evaluator Evaluator
}

func NewBuffer(axis Axis, bias bool, evaluator Evaluator) Buffer {
	return Buffer{axis, bias, evaluator}
}

func (b Buffer) Destruct() (Axis, bool, Evaluator) {
	return b.axis, b.bias, b.evaluator
}

func (b Buffer) Act(signal Signal, _ []Warrior, ac ActorContext) {
	s := signal.(*EvaluationSignal)
	if b.axis != s.Axis() {
		ac.SetTrigger(false)
		return
	}

	ac.SetTrigger(true)
	var current Warrior
	if warrior, ok := signal.Current().(Warrior); ok {
		current = warrior
	}

	if b.bias {
		s.Amend(func(v float64) float64 {
			return v + float64(b.evaluator.Evaluate(current, ac))
		})
	} else {
		s.Amend(func(v float64) float64 {
			return v * float64(b.evaluator.Evaluate(current, ac)) / 100
		})
	}
}

func (b Buffer) Fork(evaluator Evaluator) any {
	if evaluator == nil {
		return b
	}

	return NewBuffer(b.axis, b.bias, evaluator)
}

func (b Buffer) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Kind      string    `json:"_kind"`
		Axis      string    `json:"axis"`
		Bias      bool      `json:"bias,omitempty"`
		Evaluator Evaluator `json:"evaluator,omitempty"`
	}{
		"buffer",
		b.axis.String(),
		b.bias,
		b.evaluator,
	})
}

func (b *Buffer) UnmarshalJSON(data []byte) error {
	var v struct {
		Axis      Axis
		Bias      bool
		Evaluator EvaluatorFile
	}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	*b = NewBuffer(v.Axis, v.Bias, v.Evaluator.Evaluator)
	return nil
}

type VerbActor struct {
	verb      Verb
	evaluator Evaluator
}

func NewVerbActor(verb Verb, evaluator Evaluator) VerbActor {
	return VerbActor{verb, evaluator}
}

func (a VerbActor) Act(signal Signal, targets []Warrior, ac ActorContext) {
	ac.SetTrigger(true)
	e := a.evaluator
	if e != nil {
		var current Warrior
		if warrior, ok := signal.Current().(Warrior); ok {
			current = warrior
		}
		e = NewCustomEvaluator(func(Warrior, EvaluationContext) int {
			return a.evaluator.Evaluate(current, ac)
		})
	}

	ac.Queue(newAction(ac.Next(), targets, a.verb.Fork(e).(Verb)))
}

func (a VerbActor) Fork(evaluator Evaluator) any {
	if evaluator == nil {
		return a
	}

	return NewVerbActor(a.verb, evaluator)
}

func (a VerbActor) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Kind      string    `json:"_kind"`
		Verb      Verb      `json:"verb"`
		Evaluator Evaluator `json:"evaluator,omitempty"`
	}{
		"verb",
		a.verb,
		a.evaluator,
	})
}

func (a *VerbActor) UnmarshalJSON(data []byte) error {
	var v struct {
		Verb      VerbFile
		Evaluator EvaluatorFile
	}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	*a = NewVerbActor(v.Verb.Verb, v.Evaluator.Evaluator)
	return nil
}

type SelectActor struct {
	actor    Actor
	selector Selector
}

func NewSelectActor(actor Actor, selectors ...Selector) SelectActor {
	if len(selectors) == 1 {
		return SelectActor{actor, selectors[0]}
	}

	return SelectActor{actor, PipelineSelector(selectors)}
}

func (a SelectActor) Act(signal Signal, warriors []Warrior, ac ActorContext) {
	warriors = a.selector.Select(warriors, signal, ac)
	if len(warriors) == 0 {
		ac.SetTrigger(false)
		return
	}

	ac.SetTrigger(true)
	a.actor.Act(signal, warriors, ac)
}

func (a SelectActor) Fork(evaluator Evaluator) any {
	return NewSelectActor(a.actor.Fork(evaluator).(Actor), a.selector)
}

func (a SelectActor) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]any{
		"_kind":    "select",
		"selector": a.selector,
		"do":       a.actor,
	})
}

func (a *SelectActor) UnmarshalJSON(data []byte) error {
	var v struct {
		Selector SelectorFile
		Actor    ActorFile `json:"do"`
	}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	*a = NewSelectActor(v.Actor.Actor, v.Selector.Selector)
	return nil
}

type ProbabilityActor struct {
	rng       Rng
	evaluator Evaluator
	actor     Actor
}

func NewProbabilityActor(rng Rng, evaluator Evaluator, actor Actor) ProbabilityActor {
	return ProbabilityActor{rng, evaluator, actor}
}

func (a ProbabilityActor) Act(signal Signal, warriors []Warrior, ac ActorContext) {
	if float64(a.evaluator.Evaluate(signal.Current().(Warrior), ac))/100 <= a.rng.Float64() {
		ac.SetTrigger(false)
		return
	}

	ac.SetTrigger(true)
	a.actor.Act(signal, warriors, ac)
}

func (a ProbabilityActor) Fork(evaluator Evaluator) any {
	return NewProbabilityActor(a.rng, a.evaluator, a.actor.Fork(evaluator).(Actor))
}

func (a ProbabilityActor) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]any{
		"_kind":       "probability",
		"probability": a.evaluator,
		"do":          a.actor,
	})
}

func (a *ProbabilityActor) UnmarshalJSON(data []byte) error {
	var v struct {
		Probability EvaluatorFile
		Actor       ActorFile `json:"do"`
	}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	*a = NewProbabilityActor(DefaultRng, v.Probability.Evaluator, v.Actor.Actor)
	return nil
}

type SequenceActor []Actor

func NewSequenceActor(actors ...Actor) SequenceActor {
	switch len(actors) {
	case 0:
		return nil
	case 1:
		if actor, ok := actors[0].(SequenceActor); ok {
			return actor
		}
	}

	return SequenceActor(actors)
}

func (a SequenceActor) Act(signal Signal, warriors []Warrior, ac ActorContext) {
	for _, actor := range a {
		actor.Act(signal, warriors, ac)
		if !ac.Trigger() {
			break
		}
	}
	ac.SetTrigger(true)
}

func (a SequenceActor) Fork(evaluator Evaluator) any {
	actors := make([]Actor, len(a))
	for i, actor := range a {
		actors[i] = actor.Fork(evaluator).(Actor)
	}

	return NewSequenceActor(actors...)
}

func (a SequenceActor) MarshalJSON() ([]byte, error) {
	actors := []Actor(a)
	if actors == nil {
		actors = []Actor{}
	}

	return json.Marshal(map[string]any{
		"_kind": "sequence",
		"do":    actors,
	})
}

func (a *SequenceActor) UnmarshalJSON(data []byte) error {
	var v struct {
		Actors []ActorFile `json:"do"`
	}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	*a = NewSequenceActor(functional.Map(func(f ActorFile) Actor {
		return f.Actor
	})(v.Actors)...)
	return nil
}

type RepeatActor struct {
	count int
	actor Actor
}

func NewRepeatActor(count int, actors ...Actor) RepeatActor {
	return RepeatActor{count, NewSequenceActor(actors...)}
}

func (a RepeatActor) Act(signal Signal, warriors []Warrior, ac ActorContext) {
	for i := 0; i < a.count; i++ {
		a.actor.Act(signal, warriors, ac)
		if !ac.Trigger() {
			break
		}
	}
	ac.SetTrigger(true)
}

func (a RepeatActor) Fork(evaluator Evaluator) any {
	return NewRepeatActor(a.count, a.actor.Fork(evaluator).(Actor))
}

func (a RepeatActor) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]any{
		"_kind": "repeat",
		"count": a.count,
		"do":    a.actor,
	})
}

func (a *RepeatActor) UnmarshalJSON(data []byte) error {
	var v struct {
		Count int
		Actor ActorFile `json:"do"`
	}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	*a = NewRepeatActor(v.Count, v.Actor.Actor)
	return nil
}

// CriticalActor set the `critical' flag and should be triggered on
// `PreAttackSignal' with high priority.
type CriticalActor struct {
}

func (CriticalActor) Act(signal Signal, _ []Warrior, ac ActorContext) {
	ac.SetTrigger(true)
	sig := signal.(ActionSignal)
	attack := sig.Action().Verb().(*Attack)
	attack.SetCritical(true)
}

func (CriticalActor) Fork(_ Evaluator) any {
	return CriticalActor{}
}

func (CriticalActor) MarshalJSON() ([]byte, error) {
	return json.Marshal(kind{"critical"})
}

type ImmuneActor struct {
}

func (ImmuneActor) Act(signal Signal, warriors []Warrior, ac ActorContext) {
	ac.SetTrigger(true)
	action := signal.(ActionSignal).Action()
	for _, w := range warriors {
		action.AddImmuneTarget(w)
	}
}

func (ImmuneActor) Fork(_ Evaluator) any {
	return ImmuneActor{}
}

func (ImmuneActor) MarshalJSON() ([]byte, error) {
	return json.Marshal(kind{"immune"})
}

type LossStopper struct {
	evaluator Evaluator
	full      bool
}

func NewLossStopper(evaluator Evaluator, full bool) LossStopper {
	return LossStopper{evaluator, full}
}

func (s LossStopper) Act(signal Signal, _ []Warrior, ac ActorContext) {
	sig := signal.(*PreLossSignal)
	stopper := s.evaluator.Evaluate(sig.Current().(Warrior), ac)
	if sig.Loss() <= stopper {
		ac.SetTrigger(false)
		return
	}

	ac.SetTrigger(true)
	if s.full {
		sig.SetLoss(0)
	} else {
		sig.SetLoss(stopper)
	}
}

func (s LossStopper) Fork(evaluator Evaluator) any {
	if evaluator == nil {
		return s
	}

	return NewLossStopper(evaluator, false)
}

func (s LossStopper) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]any{
		"_kind":     "loss_stopper",
		"full":      s.full,
		"evaluator": s.evaluator,
	})
}

func (s *LossStopper) UnmarshalJSON(data []byte) error {
	var v struct {
		Full      bool
		Evaluator EvaluatorFile
	}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	*s = NewLossStopper(v.Evaluator.Evaluator, v.Full)
	return nil
}

type LossResister struct {
}

func (LossResister) Act(signal Signal, _ []Warrior, ac ActorContext) {
	ac.SetTrigger(true)
	s := signal.(*PreLossSignal)
	r := min(s.Loss(), ac.Capacity())
	s.SetLoss(s.Loss() - r)
	ac.Flush(r)
}

func (LossResister) Fork(_ Evaluator) any {
	return LossResister{}
}

func (LossResister) MarshalJSON() ([]byte, error) {
	return json.Marshal(kind{"loss_resister"})
}

type TheoryActor map[any]map[any]int

func NewTheoryActor(theory map[any]map[any]int) TheoryActor {
	return TheoryActor(theory)
}

func (a TheoryActor) Act(signal Signal, _ []Warrior, ac ActorContext) {
	var proto any
	for k := range a {
		proto = k
		break
	}

	_, scripter, _ := ac.(*actorContext).
		EvaluationContext.(ActionContext).Action().Script().Source()
	s, ok := queryTag(scripter, proto)
	if !ok {
		ac.SetTrigger(false)
		return
	}
	theory, ok := a[s]
	if !ok {
		ac.SetTrigger(false)
		return
	}
	t, ok := queryTag(signal.Current(), proto)
	if !ok {
		ac.SetTrigger(false)
		return
	}
	m, ok := theory[t]
	if !ok {
		ac.SetTrigger(false)
		return
	}

	ac.SetTrigger(true)
	sig := signal.(*EvaluationSignal)
	sig.Amend(func(v float64) float64 {
		return v * float64(m) / 100
	})
}

func (a TheoryActor) Fork(_ Evaluator) any {
	return a
}

func (a TheoryActor) MarshalJSON() ([]byte, error) {
	var clauses []any
	for k, v := range a {
		for kk, vv := range v {
			clauses = append(clauses, struct {
				Source any `json:"s"`
				Target any `json:"t"`
				Value  int `json:"v"`
			}{k, kk, vv})
		}
	}

	return json.Marshal(map[string]any{
		"_kind":   "theory",
		"clauses": clauses,
	})
}

func (a *TheoryActor) UnmarshalJSON(data []byte) error {
	var v struct {
		Clauses []struct {
			Source TagFile `json:"s"`
			Target TagFile `json:"t"`
			Value  int     `json:"v"`
		}
	}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	theory := make(map[any]map[any]int)
	for _, c := range v.Clauses {
		if _, ok := theory[c.Source.Tag]; !ok {
			theory[c.Source.Tag] = make(map[any]int)
		}

		theory[c.Source.Tag][c.Target.Tag] = c.Value
	}

	*a = NewTheoryActor(theory)
	return nil
}

// ActionBuffer is a special actor that buffs the action, mostly used
// to amplify the damage.
type ActionBuffer struct {
	evaluator Evaluator
	buffer    Actor
}

func NewActionBuffer(evaluator Evaluator, buffer Actor) ActionBuffer {
	return ActionBuffer{evaluator, buffer}
}

func (b ActionBuffer) Act(signal Signal, _ []Warrior, ac ActorContext) {
	ac.SetTrigger(true)
	sig := signal.(ActionSignal)
	e := b.evaluator
	if e != nil {
		var current Warrior
		if warrior, ok := signal.Current().(Warrior); ok {
			current = warrior
		}
		e = ConstEvaluator(e.Evaluate(current, ac))
	}

	sig.Action().Add(
		NewFatReactor(FatRespond(
			NewSignalTrigger(&EvaluationSignal{}),
			b.buffer.Fork(e).(Actor),
		)),
	)
}

func (b ActionBuffer) Fork(evaluator Evaluator) any {
	if evaluator == nil {
		return b
	}

	return NewActionBuffer(evaluator, b.buffer)
}

func (b ActionBuffer) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Kind      string    `json:"_kind"`
		Buffer    Actor     `json:"buffer"`
		Evaluator Evaluator `json:"evaluator,omitempty"`
	}{
		"action_buffer",
		b.buffer,
		b.evaluator,
	})
}

func (b *ActionBuffer) UnmarshalJSON(data []byte) error {
	var v struct {
		Buffer    ActorFile `json:"buffer"`
		Evaluator EvaluatorFile
	}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	*b = NewActionBuffer(v.Evaluator.Evaluator, v.Buffer.Actor)
	return nil
}

type ActorContext interface {
	EvaluationContext
	Capacitor
	Go(func())
	Queue(Action)
	Done()
	Drain()
	Await() bool
	Trigger() bool
	SetTrigger(bool)
}

type Instruction struct {
	action Action
	done   chan struct{}
}

type actorContext struct {
	EvaluationContext
	*capacitor
	trigger bool
	fx      chan func()
	ich     chan Instruction
	done    chan struct{}
}

func newActorContext(ec EvaluationContext, capacity int) *actorContext {
	return &actorContext{
		ec,
		newCapacitor(capacity),
		false,
		nil,
		make(chan Instruction),
		make(chan struct{}),
	}
}

func (a *actorContext) Go(fn func()) {
	if a.fx == nil {
		a.fx = make(chan func(), 16)
		go func() {
			for f := range a.fx {
				f()
			}
		}()
	}

	a.fx <- fn
}

func (a *actorContext) Queue(action Action) {
	done := make(chan struct{})
	a.ich <- Instruction{action, done}
	<-done
}

func (a *actorContext) Done() {
	if a.fx == nil {
		close(a.ich)
		return
	}

	a.Go(func() { close(a.ich) })
	close(a.fx)
}

func (a *actorContext) Drain() {
	for _ = range a.ich {
	}
}

func (a *actorContext) Await() bool {
	<-a.done
	return a.trigger
}

func (a *actorContext) Trigger() bool {
	return a.trigger
}

func (a *actorContext) SetTrigger(trigger bool) {
	if a.done != nil {
		a.trigger = trigger
		a.done <- struct{}{}
		a.done = nil
	}
}

type Capacitor interface {
	Capacity() int
	Flush(int)
}

type capacitor struct {
	capacity int
}

func newCapacitor(capacity int) *capacitor {
	return &capacitor{capacity}
}

func (c *capacitor) Capacity() int {
	return c.capacity
}

func (c *capacitor) Flush(n int) {
	c.capacity -= n
}

type ActorFile struct {
	Actor Actor
}

func (f *ActorFile) UnmarshalJSON(data []byte) error {
	var k kind
	if err := json.Unmarshal(data, &k); err != nil {
		return err
	}

	if actor, ok := actorType[k.Kind]; ok {
		v := reflect.New(actor)
		if err := json.Unmarshal(data, v.Interface()); err != nil {
			return err
		}

		f.Actor = v.Elem().Interface().(Actor)
		return nil
	}

	return ErrBadActor
}

var actorType = map[string]reflect.Type{
	"buffer":        reflect.TypeOf(Buffer{}),
	"verb":          reflect.TypeOf(VerbActor{}),
	"select":        reflect.TypeOf(SelectActor{}),
	"probability":   reflect.TypeOf(ProbabilityActor{}),
	"sequence":      reflect.TypeOf(SequenceActor{}),
	"repeat":        reflect.TypeOf(RepeatActor{}),
	"critical":      reflect.TypeOf(CriticalActor{}),
	"immune":        reflect.TypeOf(ImmuneActor{}),
	"loss_stopper":  reflect.TypeOf(LossStopper{}),
	"loss_resister": reflect.TypeOf(LossResister{}),
	"theory":        reflect.TypeOf(TheoryActor{}),
	"action_buffer": reflect.TypeOf(ActionBuffer{}),
}
