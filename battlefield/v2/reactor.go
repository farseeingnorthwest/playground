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

type Actor interface {
	Act(Signal, []Warrior, EvaluationContext)
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

func (a *Buffer) Act(signal Signal, _ []Warrior, ec EvaluationContext) {
	s := signal.(*EvaluationSignal)
	if a.axis != s.Axis() {
		return
	}

	var current Warrior
	if warrior, ok := signal.Current().(Warrior); ok {
		current = warrior
	}

	if a.bias {
		s.SetValue(a.evaluator.Evaluate(current, ec) + s.Value())
	} else {
		s.SetValue(a.evaluator.Evaluate(current, ec) * s.Value() / 100)
	}
}

func (a *Buffer) Fork(evaluator Evaluator) any {
	return &Buffer{a.axis, a.bias, evaluator}
}

type VerbActor struct {
	verb      Verb
	evaluator Evaluator
}

func NewVerbActor(verb Verb, evaluator Evaluator) *VerbActor {
	return &VerbActor{verb, evaluator}
}

func (a *VerbActor) Act(signal Signal, targets []Warrior, ec EvaluationContext) {
	e := a.evaluator
	if e != nil {
		var current Warrior
		if warrior, ok := signal.Current().(Warrior); ok {
			current = warrior
		}
		e = ConstEvaluator(e.Evaluate(current, ec))
	}

	signal.(Scripter).Add(NewMyAction(targets, a.verb.Fork(e).(Verb)))
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

func (a *SelectActor) Act(signal Signal, warriors []Warrior, ec EvaluationContext) {
	for _, selector := range a.selectors {
		warriors = selector.Select(warriors, signal, ec)
		if len(warriors) == 0 {
			return
		}
	}

	a.actor.Act(signal, warriors, ec)
}

func (a *SelectActor) Fork(evaluator Evaluator) any {
	return NewSelectActor(a.actor.Fork(evaluator).(Actor), a.selectors...)
}

type CriticalActor struct {
}

func (CriticalActor) Act(signal Signal, _ []Warrior, _ EvaluationContext) {
	sig := signal.(ActionSignal)
	attack := sig.Action().Verb().(*Attack)
	attack.SetCritical(true)
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

func (b *ActionBuffer) Act(signal Signal, _ []Warrior, ec EvaluationContext) {
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

func (a *ProbabilityActor) Act(signal Signal, warriors []Warrior, ec EvaluationContext) {
	if a.rng.Float64() < float64(a.evaluator.Evaluate(signal.Current().(Warrior), ec))/100 {
		a.actor.Act(signal, warriors, ec)
	}
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

func (a *SequenceActor) Act(signal Signal, warriors []Warrior, ec EvaluationContext) {
	for _, actor := range a.actors {
		actor.Act(signal, warriors, ec)
	}
}

func (a *SequenceActor) Fork(evaluator Evaluator) any {
	actors := make([]Actor, len(a.actors))
	for i, actor := range a.actors {
		actors[i] = actor.Fork(evaluator).(Actor)
	}

	return NewSequenceActor(actors...)
}
