package battlefield

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
	_ Actor = TheoryActor[any]{}
	_ Actor = (*ActionBuffer)(nil)
)

type Label string

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

type SelectActor struct {
	actor     Actor
	selectors []Selector
}

func NewSelectActor(actor Actor, selectors ...Selector) SelectActor {
	return SelectActor{actor, selectors}
}

func (a SelectActor) Act(signal Signal, warriors []Warrior, ac ActorContext) {
	for _, selector := range a.selectors {
		warriors = selector.Select(warriors, signal, ac)
		if len(warriors) == 0 {
			ac.SetTrigger(false)
			return
		}
	}

	ac.SetTrigger(true)
	a.actor.Act(signal, warriors, ac)
}

func (a SelectActor) Fork(evaluator Evaluator) any {
	return NewSelectActor(a.actor.Fork(evaluator).(Actor), a.selectors...)
}

type Rng interface {
	Float64() float64
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

type SequenceActor struct {
	actors []Actor
}

func NewSequenceActor(actors ...Actor) SequenceActor {
	return SequenceActor{actors}
}

func (a SequenceActor) Act(signal Signal, warriors []Warrior, ac ActorContext) {
	for _, actor := range a.actors {
		actor.Act(signal, warriors, ac)
		if !ac.Trigger() {
			break
		}
	}
	ac.SetTrigger(true)
}

func (a SequenceActor) Fork(evaluator Evaluator) any {
	actors := make([]Actor, len(a.actors))
	for i, actor := range a.actors {
		actors[i] = actor.Fork(evaluator).(Actor)
	}

	return NewSequenceActor(actors...)
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

type TheoryActor[T comparable] struct {
	theory map[T]map[T]int
}

func NewTheoryActor[T comparable](theory map[T]map[T]int) TheoryActor[T] {
	return TheoryActor[T]{theory}
}

func (a TheoryActor[T]) Act(signal Signal, _ []Warrior, ac ActorContext) {
	_, scripter, _ := ac.(*actorContext).
		EvaluationContext.(ActionContext).Action().Script().Source()
	s, ok := QueryTag[T](scripter)
	if !ok {
		ac.SetTrigger(false)
		return
	}
	theory, ok := a.theory[s]
	if !ok {
		ac.SetTrigger(false)
		return
	}
	t, ok := QueryTag[T](signal.Current())
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

func (a TheoryActor[T]) Fork(_ Evaluator) any {
	return a
}

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

type ActorContext interface {
	EvaluationContext
	Capacitor
	Queue(Action)
	Done()
	Await() bool
	Trigger() bool
	SetTrigger(bool)
}

type actorContext struct {
	EvaluationContext
	*capacitor
	trigger   bool
	aChan     chan Action
	semaphore chan struct{}
}

func newActorContext(ec EvaluationContext, capacity int) *actorContext {
	return &actorContext{ec, newCapacitor(capacity), false, make(chan Action), make(chan struct{})}
}

func (a *actorContext) Queue(action Action) {
	a.aChan <- action
}

func (a *actorContext) Done() {
	close(a.aChan)
}

func (a *actorContext) Await() bool {
	<-a.semaphore
	return a.trigger
}

func (a *actorContext) Trigger() bool {
	return a.trigger
}

func (a *actorContext) SetTrigger(trigger bool) {
	if a.semaphore != nil {
		a.trigger = trigger
		a.semaphore <- struct{}{}
		a.semaphore = nil
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
