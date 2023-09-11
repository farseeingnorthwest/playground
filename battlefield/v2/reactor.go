package battlefield

import (
	"encoding/json"
	"errors"
	"fmt"

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
	_ Actor = TheoryActor[Side]{}
	_ Actor = (*ActionBuffer)(nil)

	ErrBadActor = errors.New("bad actor")
)

type Reactor interface {
	React(Signal, EvaluationContext)
}

type Forker interface {
	Fork(Evaluator) any
}

type ForkReactor interface {
	Reactor
	Forker
}

type Actor interface {
	Act(Signal, []Warrior, ActorContext) (trigger bool)
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

func (b Buffer) Act(signal Signal, _ []Warrior, ac ActorContext) bool {
	s := signal.(*EvaluationSignal)
	if b.axis != s.Axis() {
		return false
	}

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

	return true
}

func (b Buffer) Fork(evaluator Evaluator) any {
	if evaluator == nil {
		return b
	}

	return NewBuffer(b.axis, b.bias, evaluator)
}

func (b Buffer) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Axis      string    `json:"buff"`
		Bias      bool      `json:"bias,omitempty"`
		Evaluator Evaluator `json:"evaluator,omitempty"`
	}{
		b.axis.String(),
		b.bias,
		b.evaluator,
	})
}

type VerbActor struct {
	verb      Verb
	evaluator Evaluator
}

func NewVerbActor(verb Verb, evaluator Evaluator) VerbActor {
	return VerbActor{verb, evaluator}
}

func (a VerbActor) Act(signal Signal, targets []Warrior, ac ActorContext) bool {
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

	signal.(Scripter).Add(newAction(targets, a.verb.Fork(e).(Verb)))
	return true
}

func (a VerbActor) Fork(evaluator Evaluator) any {
	if evaluator == nil {
		return a
	}

	return NewVerbActor(a.verb, evaluator)
}

func (a VerbActor) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Verb      Verb      `json:"verb"`
		Evaluator Evaluator `json:"evaluator,omitempty"`
	}{
		a.verb,
		a.evaluator,
	})
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

func (a SelectActor) Act(signal Signal, warriors []Warrior, ac ActorContext) bool {
	warriors = a.selector.Select(warriors, signal, ac)
	if len(warriors) == 0 {
		return false
	}

	a.actor.Act(signal, warriors, ac)
	return true
}

func (a SelectActor) Fork(evaluator Evaluator) any {
	return NewSelectActor(a.actor.Fork(evaluator).(Actor), a.selector)
}

func (a SelectActor) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]any{
		"for": a.selector,
		"do":  a.actor,
	})
}

type ProbabilityActor struct {
	rng       Rng
	evaluator Evaluator
	actor     Actor
}

func NewProbabilityActor(rng Rng, evaluator Evaluator, actor Actor) ProbabilityActor {
	return ProbabilityActor{rng, evaluator, actor}
}

func (a ProbabilityActor) Act(signal Signal, warriors []Warrior, ec ActorContext) bool {
	if a.rng.Float64() < float64(a.evaluator.Evaluate(signal.Current().(Warrior), ec))/100 {
		a.actor.Act(signal, warriors, ec)
	}

	return true
}

func (a ProbabilityActor) Fork(evaluator Evaluator) any {
	return NewProbabilityActor(a.rng, a.evaluator, a.actor.Fork(evaluator).(Actor))
}

func (a ProbabilityActor) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]any{
		"probability": a.evaluator,
		"do":          a.actor,
	})
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

func (a SequenceActor) Act(signal Signal, warriors []Warrior, ac ActorContext) bool {
	for _, actor := range a {
		if !actor.Act(signal, warriors, ac) {
			return false
		}
	}

	return true
}

func (a SequenceActor) Fork(evaluator Evaluator) any {
	actors := make([]Actor, len(a))
	for i, actor := range a {
		actors[i] = actor.Fork(evaluator).(Actor)
	}

	return NewSequenceActor(actors...)
}

func (a SequenceActor) MarshalJSON() ([]byte, error) {
	if a == nil {
		return json.Marshal([]Actor{})
	}

	return json.Marshal([]Actor(a))
}

type RepeatActor struct {
	count int
	actor Actor
}

func NewRepeatActor(count int, actors ...Actor) RepeatActor {
	return RepeatActor{count, NewSequenceActor(actors...)}
}

func (a RepeatActor) Act(signal Signal, warriors []Warrior, ec ActorContext) bool {
	for i := 0; i < a.count; i++ {
		if !a.actor.Act(signal, warriors, ec) {
			return false
		}
	}

	return true
}

func (a RepeatActor) Fork(evaluator Evaluator) any {
	return NewRepeatActor(a.count, a.actor.Fork(evaluator).(Actor))
}

func (a RepeatActor) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]any{
		"repeat": a.count,
		"do":     a.actor,
	})
}

// CriticalActor set the `critical' flag and should be triggered on
// `PreAttackSignal' with high priority.
type CriticalActor struct {
}

func (CriticalActor) Act(signal Signal, _ []Warrior, _ ActorContext) bool {
	sig := signal.(ActionSignal)
	attack := sig.Action().Verb().(*Attack)
	attack.SetCritical(true)
	return true
}

func (CriticalActor) Fork(_ Evaluator) any {
	return CriticalActor{}
}

func (CriticalActor) MarshalJSON() ([]byte, error) {
	return json.Marshal("critical_strike")
}

type ImmuneActor struct {
}

func (ImmuneActor) Act(signal Signal, warriors []Warrior, _ ActorContext) bool {
	action := signal.(ActionSignal).Action()
	for _, w := range warriors {
		action.AddImmuneTarget(w)
	}

	return true
}

func (ImmuneActor) Fork(_ Evaluator) any {
	return ImmuneActor{}
}

func (ImmuneActor) MarshalJSON() ([]byte, error) {
	return json.Marshal("immune")
}

type LossStopper struct {
	evaluator Evaluator
	full      bool
}

func NewLossStopper(evaluator Evaluator, full bool) LossStopper {
	return LossStopper{evaluator, full}
}

func (s LossStopper) Act(signal Signal, _ []Warrior, ec ActorContext) bool {
	sig := signal.(*PreLossSignal)
	stopper := s.evaluator.Evaluate(sig.Current().(Warrior), ec)
	if sig.Loss() <= stopper {
		return false
	}

	if s.full {
		sig.SetLoss(0)
	} else {
		sig.SetLoss(stopper)
	}
	return true
}

func (s LossStopper) Fork(evaluator Evaluator) any {
	if evaluator == nil {
		return s
	}

	return NewLossStopper(evaluator, false)
}

func (s LossStopper) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]any{
		"stop_loss": s.full,
		"evaluator": s.evaluator,
	})
}

type LossResister struct {
}

func (LossResister) Act(signal Signal, _ []Warrior, ec ActorContext) bool {
	s := signal.(*PreLossSignal)
	r := min(s.Loss(), ec.Capacity())
	s.SetLoss(s.Loss() - r)
	ec.Flush(r)

	return true
}

func (LossResister) Fork(_ Evaluator) any {
	return LossResister{}
}

func (LossResister) MarshalJSON() ([]byte, error) {
	return json.Marshal("resist_loss")
}

type feature interface {
	comparable
	fmt.Stringer
}

type theory[T comparable] map[T]map[T]int

type TheoryActor[T feature] struct {
	theory theory[T]
}

func NewTheoryActor[T feature](theory map[T]map[T]int) TheoryActor[T] {
	return TheoryActor[T]{theory}
}

func (a TheoryActor[T]) Act(signal Signal, _ []Warrior, ac ActorContext) bool {
	scripter, _ := ac.(*actorContext).
		EvaluationContext.(ActionContext).Action().Script().Source()
	s, ok := QueryTag[T](scripter)
	if !ok {
		return false
	}
	theory, ok := a.theory[s]
	if !ok {
		return false
	}
	t, ok := QueryTag[T](signal.Current())
	if !ok {
		return false
	}
	m, ok := theory[t]
	if !ok {
		return false
	}

	sig := signal.(*EvaluationSignal)
	sig.Amend(func(v float64) float64 {
		return v * float64(m) / 100
	})

	return true
}

func (a TheoryActor[T]) Fork(_ Evaluator) any {
	return a
}

func (a TheoryActor[T]) MarshalJSON() ([]byte, error) {
	t := make(theory[string])
	for k, v := range a.theory {
		t[k.String()] = make(map[string]int)
		for kk, vv := range v {
			t[k.String()][kk.String()] = vv
		}
	}

	return json.Marshal(map[string]theory[string]{
		"theory": t,
	})
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

func (b ActionBuffer) Act(signal Signal, _ []Warrior, ac ActorContext) bool {
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
	return true
}

func (b ActionBuffer) Fork(evaluator Evaluator) any {
	if evaluator == nil {
		return b
	}

	return NewActionBuffer(evaluator, b.buffer)
}

func (b ActionBuffer) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Buffer    Actor     `json:"buff_a"`
		Evaluator Evaluator `json:"evaluator,omitempty"`
	}{
		b.buffer,
		b.evaluator,
	})
}

type ActorContext interface {
	EvaluationContext
	Capacitor
}

type actorContext struct {
	EvaluationContext
	*capacitor
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

type ActorFile[T feature] struct {
	Actor Actor
}

func (f *ActorFile[T]) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		var e *json.UnmarshalTypeError
		if !errors.As(err, &e) {
			return err
		}
	} else {
		if actor, ok := map[string]Actor{
			"critical_strike": CriticalActor{},
			"immune":          ImmuneActor{},
			"resist_loss":     LossResister{},
		}[s]; ok {
			f.Actor = actor
			return nil
		}

		return ErrBadActor
	}

	var m map[string]json.RawMessage
	if err := json.Unmarshal(data, &m); err != nil {
		var e *json.UnmarshalTypeError
		if !errors.As(err, &e) {
			return err
		}

		var fs []ActorFile[T]
		if err := json.Unmarshal(data, &fs); err != nil {
			return err
		}

		f.Actor = NewSequenceActor(functional.Map(func(f ActorFile[T]) Actor {
			return f.Actor
		})(fs)...)
		return nil
	}

	if _, ok := m["buff"]; ok {
		var b struct {
			Axis      Axis `json:"buff"`
			Bias      bool
			Evaluator EvaluatorFile
		}
		if err := json.Unmarshal(data, &b); err != nil {
			return err
		}

		f.Actor = NewBuffer(b.Axis, b.Bias, b.Evaluator.Evaluator)
		return nil
	}

	if _, ok := m["verb"]; ok {
		var v struct {
			Verb      VerbFile[T]
			Evaluator EvaluatorFile
		}
		if err := json.Unmarshal(data, &v); err != nil {
			return err
		}

		f.Actor = NewVerbActor(v.Verb.Verb, v.Evaluator.Evaluator)
		return nil
	}

	if _, ok := m["stop_loss"]; ok {
		var s struct {
			Full      bool `json:"stop_loss"`
			Evaluator EvaluatorFile
		}
		if err := json.Unmarshal(data, &s); err != nil {
			return err
		}

		f.Actor = NewLossStopper(s.Evaluator.Evaluator, s.Full)
		return nil
	}

	if thr, ok := m["theory"]; ok {
		var t theory[T]
		if err := json.Unmarshal(thr, &t); err != nil {
			return err
		}

		f.Actor = NewTheoryActor(t)
		return nil
	}

	if _, ok := m["buff_a"]; ok {
		var b struct {
			Buffer    ActorFile[T] `json:"buff_a"`
			Evaluator EvaluatorFile
		}
		if err := json.Unmarshal(data, &b); err != nil {
			return err
		}

		f.Actor = NewActionBuffer(b.Evaluator.Evaluator, b.Buffer.Actor)
		return nil
	}

	data2, ok := m["do"]
	if !ok {
		return ErrBadActor
	}

	var f2 ActorFile[T]
	if err := json.Unmarshal(data2, &f2); err != nil {
		return err
	}

	if data, ok := m["for"]; ok {
		var s SelectorFile
		if err := json.Unmarshal(data, &s); err != nil {
			return err
		}

		f.Actor = NewSelectActor(f2.Actor, s.Selector)
		return nil
	}

	if data, ok := m["probability"]; ok {
		var e EvaluatorFile
		if err := json.Unmarshal(data, &e); err != nil {
			return err
		}

		f.Actor = NewProbabilityActor(DefaultRng, e.Evaluator, f2.Actor)
		return nil
	}

	if data, ok := m["repeat"]; ok {
		var n int
		if err := json.Unmarshal(data, &n); err != nil {
			return err
		}

		f.Actor = NewRepeatActor(n, f2.Actor)
		return nil
	}

	return ErrBadActor
}
