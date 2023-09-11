package battlefield

import (
	"encoding/json"
	"errors"
	"reflect"

	"github.com/farseeingnorthwest/playground/battlefield/v2/functional"
)

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

	ErrBadTrigger = errors.New("bad trigger")
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

func (t SignalTrigger) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]string{
		"signal": t.signal.Name(),
	})
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

func (t AnyTrigger) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.triggers)
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

func (t *FatTrigger) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]any{
		"signal": t.signalTrigger.signal.Name(),
		"if":     t.actionTriggers,
	})
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

func (CurrentIsSourceTrigger) MarshalJSON() ([]byte, error) {
	return json.Marshal("current_is_source")
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

func (CurrentIsTargetTrigger) MarshalJSON() ([]byte, error) {
	return json.Marshal("current_is_target")
}

type ReactorTrigger struct {
	tag any
}

func NewReactorTrigger(tag any) ReactorTrigger {
	return ReactorTrigger{tag}
}

func (t ReactorTrigger) Trigger(action Action, _ Signal, _ EvaluationContext) bool {
	_, r := action.Script().Source()
	tagger, ok := r.(Tagger)
	if !ok {
		return false
	}

	return tagger.Match(t.tag)
}

func (t ReactorTrigger) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]any{
		"source": t.tag,
	})
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

func (VerbTrigger[T]) MarshalJSON() ([]byte, error) {
	var verb T

	return json.Marshal(map[string]string{
		"verb": verb.Name(),
	})
}

type CriticalStrikeTrigger struct {
}

func (CriticalStrikeTrigger) Trigger(action Action, _ Signal, _ EvaluationContext) bool {
	a, ok := action.Verb().(*Attack)
	return ok && a.Critical()
}

func (CriticalStrikeTrigger) MarshalJSON() ([]byte, error) {
	return json.Marshal("critical_strike")
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

func (t TagTrigger) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]any{
		"tag": t.tag,
	})
}

type TriggerFile struct {
	Trigger Trigger
}

func (f *TriggerFile) UnmarshalJSON(bytes []byte) error {
	var m map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &m); err != nil {
		var e *json.UnmarshalTypeError
		if !errors.As(err, &e) {
			return err
		}

		var fs []TriggerFile
		if err := json.Unmarshal(bytes, &fs); err != nil {
			return err
		}

		f.Trigger = NewAnyTrigger(
			functional.Map(func(f TriggerFile) Trigger {
				return f.Trigger
			})(fs)...,
		)
		return nil
	}

	sig, ok := m["signal"]
	if !ok {
		return ErrBadTrigger
	}
	var s string
	if err := json.Unmarshal(sig, &s); err != nil {
		return err
	}
	signal, ok := map[string]Signal{
		"evaluation":   &EvaluationSignal{},
		"pre_loss":     &PreLossSignal{},
		"launch":       &LaunchSignal{},
		"battle_start": &BattleStartSignal{},
		"round_start":  &RoundStartSignal{},
		"round_end":    &RoundEndSignal{},
		"pre_action":   &PreActionSignal{},
		"post_action":  &PostActionSignal{},
		"lifecycle":    &LifecycleSignal{},
	}[s]
	if !ok {
		return ErrBadTrigger
	}

	actionTriggers, ok := m["if"]
	if !ok {
		f.Trigger = NewSignalTrigger(signal)
		return nil
	}

	var fs []ActionTriggerFile
	if err := json.Unmarshal(actionTriggers, &fs); err != nil {
		return err
	}

	f.Trigger = NewFatTrigger(
		signal,
		functional.Map(func(f ActionTriggerFile) ActionTrigger {
			return f.Trigger
		})(fs)...,
	)
	return nil
}

type ActionTriggerFile struct {
	Trigger ActionTrigger
}

func (f *ActionTriggerFile) UnmarshalJSON(bytes []byte) error {
	var s string
	if err := json.Unmarshal(bytes, &s); err != nil {
		var e *json.UnmarshalTypeError
		if !errors.As(err, &e) {
			return err
		}
	} else {
		if trigger, ok := map[string]ActionTrigger{
			"current_is_source": CurrentIsSourceTrigger{},
			"current_is_target": CurrentIsTargetTrigger{},
			"critical_strike":   CriticalStrikeTrigger{},
		}[s]; ok {
			f.Trigger = trigger
			return nil
		}

		return ErrBadTrigger
	}

	var m map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &m); err != nil {
		var e *json.UnmarshalTypeError
		if !errors.As(err, &e) {
			return err
		}
	}
	if source, ok := m["source"]; ok {
		var t TagFile
		if err := json.Unmarshal(source, &t); err != nil {
			return err
		}

		f.Trigger = NewReactorTrigger(t.Tag)
		return nil
	} else if verb, ok := m["verb"]; ok {
		var v string
		if err := json.Unmarshal(verb, &v); err != nil {
			return err
		}

		if t, ok := map[string]ActionTrigger{
			"attack": NewVerbTrigger[*Attack](),
			"heal":   NewVerbTrigger[*Heal](),
			"buff":   NewVerbTrigger[*Buff](),
			"purge":  NewVerbTrigger[*Purge](),
		}[v]; ok {
			f.Trigger = t
			return nil
		}
	} else if tag, ok := m["tag"]; ok {
		var t TagFile
		if err := json.Unmarshal(tag, &t); err != nil {
			return err
		}

		f.Trigger = NewTagTrigger(t.Tag)
		return nil
	}

	return ErrBadTrigger
}
