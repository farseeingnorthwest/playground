package battlefield

import "reflect"

var (
	_ Trigger       = SignalTrigger{}
	_ Trigger       = AnyTrigger{}
	_ Trigger       = (*FatTrigger)(nil)
	_ ActionTrigger = CurrentIsSourceTrigger{}
	_ ActionTrigger = CurrentIsTargetTrigger{}
	_ ActionTrigger = ReactorTrigger{}
	_ ActionTrigger = VerbTrigger[*Attack]{}
	_ ActionTrigger = CriticalStrikeTrigger{}
	_ ActionTrigger = TagTrigger{}
)

type Trigger interface {
	Trigger(Signal, EvaluationContext) bool
}

type SignalTrigger struct {
	signal Signal
}

func NewSignalTrigger(signal Signal) SignalTrigger {
	return SignalTrigger{signal}
}

func (t SignalTrigger) Trigger(signal Signal, _ EvaluationContext) bool {
	return reflect.TypeOf(t.signal) == reflect.TypeOf(signal)
}

type AnyTrigger struct {
	triggers []Trigger
}

func NewAnyTrigger(triggers ...Trigger) AnyTrigger {
	return AnyTrigger{triggers}
}

func (t AnyTrigger) Trigger(signal Signal, ec EvaluationContext) bool {
	for _, trigger := range t.triggers {
		if trigger.Trigger(signal, ec) {
			return true
		}
	}

	return false
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

type ActionTrigger interface {
	Trigger(Action, Signal, EvaluationContext) bool
}

type CurrentIsSourceTrigger struct {
}

func (CurrentIsSourceTrigger) Trigger(action Action, signal Signal, _ EvaluationContext) bool {
	_, a, _ := action.Script().Source()
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

func NewReactorTrigger(reactor Reactor) ReactorTrigger {
	return ReactorTrigger{reactor}
}

func (t ReactorTrigger) Trigger(action Action, _ Signal, _ EvaluationContext) bool {
	_, _, r := action.Script().Source()
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

type TagTrigger struct {
	tag any
}

func NewTagTrigger(tag any) TagTrigger {
	return TagTrigger{tag}
}

func (t TagTrigger) Trigger(action Action, _ Signal, _ EvaluationContext) bool {
	buff, ok := action.Verb().(*Buff)
	if !ok {
		return false
	}
	tagger, ok := buff.Reactor().(Tagger)
	if !ok {
		return false
	}

	return tagger.Match(t.tag)
}
