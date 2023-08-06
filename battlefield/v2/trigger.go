package battlefield

import "reflect"

type Trigger interface {
	Trigger(Signal) bool
}

type SignalTrigger struct {
	signal Signal
}

func NewSignalTrigger(signal Signal) *SignalTrigger {
	return &SignalTrigger{signal}
}

func (t *SignalTrigger) Trigger(signal Signal) bool {
	return reflect.TypeOf(t.signal) == reflect.TypeOf(signal)
}

type ActionTrigger interface {
	Trigger(Action, Signal) bool
}

type CurrentIsSourceTrigger struct {
}

func (CurrentIsSourceTrigger) Trigger(action Action, signal Signal) bool {
	a, _ := action.Script().Source()
	return a == signal.Current()
}

type CurrentIsTargetTrigger struct {
}

func (CurrentIsTargetTrigger) Trigger(action Action, signal Signal) bool {
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

func (t *ReactorTrigger) Trigger(action Action, _ Signal) bool {
	_, r := action.Script().Source()
	return r == t.reactor
}

type CriticalStrikeTrigger struct {
}

func (CriticalStrikeTrigger) Trigger(action Action, _ Signal) bool {
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

func (t *LossTrigger) Trigger(action Action, signal Signal) bool {
	a, ok := action.Verb().(*Attack)
	return ok && t.comparator.Compare(a.Loss(), t.evaluator.Evaluate(signal.Current().(Warrior)))
}

type FatTrigger struct {
	signalTrigger  SignalTrigger
	actionTriggers []ActionTrigger
}

func NewFatTrigger(signal Signal, triggers ...ActionTrigger) *FatTrigger {
	return &FatTrigger{SignalTrigger{signal}, triggers}
}

func (t *FatTrigger) Trigger(signal Signal) bool {
	if !t.signalTrigger.Trigger(signal) {
		return false
	}

	a := signal.(ActionSignal).Action()
	for _, trigger := range t.actionTriggers {
		if !trigger.Trigger(a, signal) {
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

func (t *OrTrigger) Trigger(signal Signal) bool {
	for _, trigger := range t.triggers {
		if trigger.Trigger(signal) {
			return true
		}
	}

	return false
}
