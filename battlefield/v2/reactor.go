package battlefield

type Reactor interface {
	React(Signal, []Warrior)
}

type Forker interface {
	Fork(Evaluator) any
}

type Tagger interface {
	Tags() []any
	Match(...any) bool
}

type ForkReactor interface {
	Reactor
	Forker
}

type Actor interface {
	Act(Signal, []Warrior)
	Forker
}

type Buffer struct {
	axis       Axis
	multiplier Evaluator
}

func NewBuffer(axis Axis, multiplier Evaluator) *Buffer {
	return &Buffer{axis, multiplier}
}

func (a *Buffer) Act(signal Signal, _ []Warrior) {
	s := signal.(*EvaluationSignal)
	var current Warrior
	if warrior, ok := signal.Current().(Warrior); ok {
		current = warrior
	}
	s.SetValue(a.multiplier.Evaluate(current) * s.Value() / 100)
}

func (a *Buffer) Fork(evaluator Evaluator) any {
	return NewBuffer(a.axis, evaluator)
}

type VerbActor struct {
	verb      Verb
	evaluator Evaluator
}

func NewVerbActor(verb Verb, evaluator Evaluator) *VerbActor {
	return &VerbActor{verb, evaluator}
}

func (a *VerbActor) Act(signal Signal, targets []Warrior) {
	e := a.evaluator
	if e != nil {
		var current Warrior
		if warrior, ok := signal.Current().(Warrior); ok {
			current = warrior
		}
		e = ConstEvaluator(e.Evaluate(current))
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

func (a *SelectActor) Act(signal Signal, warriors []Warrior) {
	for _, selector := range a.selectors {
		warriors = selector.Select(warriors, signal)
		if len(warriors) == 0 {
			return
		}
	}

	a.actor.Act(signal, warriors)
}

func (a *SelectActor) Fork(evaluator Evaluator) any {
	return NewSelectActor(a.actor.Fork(evaluator).(Actor), a.selectors...)
}

type CriticalActor struct {
}

func (CriticalActor) Act(signal Signal, _ []Warrior) {
	sig := signal.(ActionSignal)
	attack := sig.Action().Verb().(*Attack)
	attack.SetCritical(true)
}

func (CriticalActor) Fork(_ Evaluator) any {
	return CriticalActor{}
}

type Rng interface {
	Float64() float64
}

type ProbabilityActor struct {
	rng       Rng
	evaluator Evaluator
	actor     Actor
}

func NewProbabilityActor(rng Rng, evaluator Evaluator, actor Actor) *ProbabilityActor {
	return &ProbabilityActor{rng, evaluator, actor}
}

func (a *ProbabilityActor) Act(signal Signal, warriors []Warrior) {
	if int(a.rng.Float64()*100) < a.evaluator.Evaluate(signal.Current().(Warrior)) {
		a.actor.Act(signal, warriors)
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

func (a *SequenceActor) Act(signal Signal, warriors []Warrior) {
	for _, actor := range a.actors {
		actor.Act(signal, warriors)
	}
}

func (a *SequenceActor) Fork(evaluator Evaluator) any {
	actors := make([]Actor, len(a.actors))
	for i, actor := range a.actors {
		actors[i] = actor.Fork(evaluator).(Actor)
	}

	return NewSequenceActor(actors...)
}
