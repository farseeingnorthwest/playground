package battlefield

import "reflect"

type Trigger interface {
	Trigger(Signal, EvaluationContext) bool
}

type SignalTrigger struct {
	signal Signal
}

func NewSignalTrigger(signal Signal) *SignalTrigger {
	return &SignalTrigger{signal}
}

func (t *SignalTrigger) Trigger(signal Signal, _ EvaluationContext) bool {
	return reflect.TypeOf(t.signal) == reflect.TypeOf(signal)
}

type ActionTrigger interface {
	Trigger(Action, Signal, EvaluationContext) bool
}

type CurrentIsSourceTrigger struct {
}

func (CurrentIsSourceTrigger) Trigger(action Action, signal Signal, _ EvaluationContext) bool {
	a, _ := action.Script().Source()
	return a == signal.Current()
}

type CurrentIsTargetTrigger struct {
}

func (CurrentIsTargetTrigger) Trigger(action Action, signal Signal, _ EvaluationContext) bool {
	for _, target := range action.Targets() {
		if target == signal.Current() {
			return true
		}
	}
	return false
}

type ReactorTrigger struct {
	reactor Reactor
}

func NewActionReactorTrigger(reactor Reactor) *ReactorTrigger {
	return &ReactorTrigger{reactor}
}

func (t *ReactorTrigger) Trigger(action Action, _ Signal, _ EvaluationContext) bool {
	_, r := action.Script().Source()
	return r == t.reactor
}

type VerbTrigger[T Verb] struct {
}

func NewVerbTrigger[T Verb]() VerbTrigger[T] {
	return VerbTrigger[T]{}
}

func (VerbTrigger[T]) Trigger(action Action, _ Signal, _ EvaluationContext) bool {
	_, ok := action.Verb().(T)
	return ok
}

type CriticalStrikeTrigger struct {
}

func (CriticalStrikeTrigger) Trigger(action Action, _ Signal, _ EvaluationContext) bool {
	a, ok := action.Verb().(*Attack)
	return ok && a.Critical()
}

type LossTrigger struct {
	comparator IntComparator
	evaluator  Evaluator
}

func NewLossTrigger(comparator IntComparator, evaluator Evaluator) *LossTrigger {
	return &LossTrigger{comparator, evaluator}
}

func (t *LossTrigger) Trigger(action Action, signal Signal, ec EvaluationContext) bool {
	sig := signal.(*LossSignal)
	a, ok := action.Verb().(*Attack)
	return ok && t.comparator.Compare(a.Loss()[sig.Target()], t.evaluator.Evaluate(signal.Current().(Warrior), ec))
}

type FatTrigger struct {
	signalTrigger  SignalTrigger
	actionTriggers []ActionTrigger
}

func NewFatTrigger(signal Signal, triggers ...ActionTrigger) *FatTrigger {
	return &FatTrigger{SignalTrigger{signal}, triggers}
}

func (t *FatTrigger) Trigger(signal Signal, ec EvaluationContext) bool {
	if !t.signalTrigger.Trigger(signal, ec) {
		return false
	}

	a := signal.(ActionSignal).Action()
	for _, trigger := range t.actionTriggers {
		if !trigger.Trigger(a, signal, ec) {
			return false
		}
	}

	return true
}

type OrTrigger struct {
	triggers []Trigger
}

func NewOrTrigger(triggers ...Trigger) *OrTrigger {
	return &OrTrigger{triggers}
}

func (t *OrTrigger) Trigger(signal Signal, ec EvaluationContext) bool {
	for _, trigger := range t.triggers {
		if trigger.Trigger(signal, ec) {
			return true
		}
	}

	return false
}
