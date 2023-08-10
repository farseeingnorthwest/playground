package battlefield

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

type Label string

type Capacitor interface {
	Capacity() int
	Flush(int)
}

type PlainCapacitor struct {
	capacity int
}

func NewPlainCapacitor(capacity int) *PlainCapacitor {
	return &PlainCapacitor{capacity}
}

func (c *PlainCapacitor) Capacity() int {
	return c.capacity
}

func (c *PlainCapacitor) Flush(n int) {
	c.capacity -= n
}

type ActorContext interface {
	EvaluationContext
	Capacitor
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

func NewBuffer(axis Axis, bias bool, evaluator Evaluator) *Buffer {
	return &Buffer{axis, bias, evaluator}
}

func (a *Buffer) Act(signal Signal, _ []Warrior, ec ActorContext) bool {
	s := signal.(*EvaluationSignal)
	if a.axis != s.Axis() {
		return false
	}

	var current Warrior
	if warrior, ok := signal.Current().(Warrior); ok {
		current = warrior
	}

	if a.bias {
		s.Amend(func(v float64) float64 {
			return v + float64(a.evaluator.Evaluate(current, ec))
		})
	} else {
		s.Amend(func(v float64) float64 {
			return v * float64(a.evaluator.Evaluate(current, ec)) / 100
		})
	}

	return true
}

func (a *Buffer) Fork(evaluator Evaluator) any {
	if evaluator == nil {
		return a
	}

	return &Buffer{a.axis, a.bias, evaluator}
}

type VerbActor struct {
	verb      Verb
	evaluator Evaluator
}

func NewVerbActor(verb Verb, evaluator Evaluator) *VerbActor {
	return &VerbActor{verb, evaluator}
}

func (a *VerbActor) Act(signal Signal, targets []Warrior, ec ActorContext) bool {
	e := a.evaluator
	if e != nil {
		var current Warrior
		if warrior, ok := signal.Current().(Warrior); ok {
			current = warrior
		}
		e = NewCustomEvaluator(func(warrior Warrior, context EvaluationContext) int {
			return a.evaluator.Evaluate(current, ec)
		})
	}

	signal.(Scripter).Add(NewMyAction(targets, a.verb.Fork(e).(Verb)))
	return true
}

func (a *VerbActor) Fork(evaluator Evaluator) any {
	if evaluator == nil {
		return a
	}

	return NewVerbActor(a.verb, evaluator)
}

type SelectActor struct {
	actor     Actor
	selectors []Selector
}

func NewSelectActor(actor Actor, selectors ...Selector) *SelectActor {
	return &SelectActor{actor, selectors}
}

func (a *SelectActor) Act(signal Signal, warriors []Warrior, ec ActorContext) bool {
	for _, selector := range a.selectors {
		warriors = selector.Select(warriors, signal, ec)
		if len(warriors) == 0 {
			return false
		}
	}

	a.actor.Act(signal, warriors, ec)
	return true
}

func (a *SelectActor) Fork(evaluator Evaluator) any {
	return NewSelectActor(a.actor.Fork(evaluator).(Actor), a.selectors...)
}

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

type ActionBuffer struct {
	evaluator Evaluator
	buffer    *Buffer
}

func NewActionBuffer(evaluator Evaluator, buffer *Buffer) *ActionBuffer {
	return &ActionBuffer{evaluator, buffer}
}

func (b *ActionBuffer) Act(signal Signal, _ []Warrior, ec ActorContext) bool {
	sig := signal.(ActionSignal)
	e := b.evaluator
	if e != nil {
		var current Warrior
		if warrior, ok := signal.Current().(Warrior); ok {
			current = warrior
		}
		e = ConstEvaluator(e.Evaluate(current, ec))
	}

	sig.Action().Add(
		NewFatReactor(FatRespond(
			NewSignalTrigger(&EvaluationSignal{}),
			b.buffer.Fork(e).(*Buffer),
		)),
	)
	return true
}

func (b *ActionBuffer) Fork(evaluator Evaluator) any {
	if evaluator == nil {
		return b
	}

	return NewActionBuffer(evaluator, b.buffer)
}

type ProbabilityActor struct {
	rng       Rng
	evaluator Evaluator
	actor     Actor
}

func NewProbabilityActor(rng Rng, evaluator Evaluator, actor Actor) *ProbabilityActor {
	return &ProbabilityActor{rng, evaluator, actor}
}

func (a *ProbabilityActor) Act(signal Signal, warriors []Warrior, ec ActorContext) bool {
	if a.rng.Float64() < float64(a.evaluator.Evaluate(signal.Current().(Warrior), ec))/100 {
		a.actor.Act(signal, warriors, ec)
	}

	return true
}

func (a *ProbabilityActor) Fork(evaluator Evaluator) any {
	return NewProbabilityActor(a.rng, a.evaluator, a.actor.Fork(evaluator).(Actor))
}

type SequenceActor struct {
	actors []Actor
}

func NewSequenceActor(actors ...Actor) *SequenceActor {
	return &SequenceActor{actors}
}

func (a *SequenceActor) Act(signal Signal, warriors []Warrior, ec ActorContext) bool {
	for _, actor := range a.actors {
		if !actor.Act(signal, warriors, ec) {
			return false
		}
	}

	return true
}

func (a *SequenceActor) Fork(evaluator Evaluator) any {
	actors := make([]Actor, len(a.actors))
	for i, actor := range a.actors {
		actors[i] = actor.Fork(evaluator).(Actor)
	}

	return NewSequenceActor(actors...)
}

type LossStopper struct {
	evaluator Evaluator
	zero      bool
}

func NewLossStopper(evaluator Evaluator, zero bool) *LossStopper {
	return &LossStopper{evaluator, zero}
}

func (s *LossStopper) Act(signal Signal, _ []Warrior, ec ActorContext) bool {
	sig := signal.(*PreLossSignal)
	stopper := s.evaluator.Evaluate(sig.Current().(Warrior), ec)
	if sig.Loss() <= stopper {
		return false
	}

	if s.zero {
		sig.SetLoss(0)
	} else {
		sig.SetLoss(stopper)
	}
	return true
}

func (s *LossStopper) Fork(evaluator Evaluator) any {
	if evaluator == nil {
		return s
	}

	return NewLossStopper(evaluator, false)
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

type RepeatActor struct {
	count int
	actor Actor
}

func NewRepeatActor(count int, actors ...Actor) *RepeatActor {
	return &RepeatActor{count, NewSequenceActor(actors...)}
}

func (a *RepeatActor) Act(signal Signal, warriors []Warrior, ec ActorContext) bool {
	for i := 0; i < a.count; i++ {
		if !a.actor.Act(signal, warriors, ec) {
			return false
		}
	}

	return true
}

func (a *RepeatActor) Fork(evaluator Evaluator) any {
	return NewRepeatActor(a.count, a.actor.Fork(evaluator).(Actor))
}
