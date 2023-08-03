package battlefield

import (
	"encoding/json"

	"github.com/farseeingnorthwest/playground/battlefield/v2/mod"
)

type ScriptJournal struct {
	Current *Fighter
	Source  Reactor
	Actions []*ActionJournal
}

func (s *ScriptJournal) MarshalJSON() ([]byte, error) {
	return json.Marshal(marshalScript(s))
}

type ActionJournal struct {
	Interests   []Interest
	PreScripts  []*ScriptJournal
	PostScripts []*ScriptJournal
	post        bool
}

func (a *ActionJournal) MarshalJSON() ([]byte, error) {
	return json.Marshal(marshalAction(a))
}

type Observer struct {
	mod.TaggerMod
	Scripts     []*ScriptJournal
	scriptStack []*ScriptJournal
}

func (o *Observer) React(s Signal) {
	switch sig := s.(type) {
	case *PreScriptSignal:
		o.push(&ScriptJournal{
			Current: sig.Script.Current,
			Source:  sig.Script.Source,
		})

	case *PreActionSignal:
		o.put(&ActionJournal{})

	case *PostActionSignal:
		o.log(sig.Interests)

	case *PostScriptSignal:
		o.pop()
	}
}

func (o *Observer) push(s *ScriptJournal) {
	o.scriptStack = append(o.scriptStack, s)
}

func (o *Observer) pop() {
	latest := o.current()
	o.scriptStack[len(o.scriptStack)-1] = nil
	o.scriptStack = o.scriptStack[:len(o.scriptStack)-1]
	if len(o.scriptStack) == 0 {
		o.Scripts = append(o.Scripts, latest)
		return
	}

	current := o.current()
	a := current.Actions[len(current.Actions)-1]
	if a.post {
		a.PostScripts = append(a.PostScripts, latest)
	} else {
		a.PreScripts = append(a.PreScripts, latest)
	}
}

func (o *Observer) put(a *ActionJournal) {
	current := o.current()
	current.Actions = append(current.Actions, a)
}

func (o *Observer) log(i []Interest) {
	current := o.current()
	a := current.Actions[len(current.Actions)-1]
	a.Interests = i
	a.post = true
}

func (o *Observer) current() *ScriptJournal {
	return o.scriptStack[len(o.scriptStack)-1]
}

func fmap[T, U any](f func(T) U) func([]T) []U {
	return func(a []T) []U {
		b := make([]U, len(a))
		for i, x := range a {
			b[i] = f(x)
		}
		return b
	}
}

func marshalScript(s *ScriptJournal) map[string]any {
	var current any
	if s.Current != nil {
		current = s.Current.Baseline.(mod.Tagger).Tag()
	}

	return map[string]any{
		"current": current,
		"sources": s.Source.Tag(),
		"actions": fmap(marshalAction)(s.Actions),
	}
}

func marshalAction(a *ActionJournal) map[string]any {
	return map[string]any{
		"interests":    fmap(marshalInterest)(a.Interests),
		"pre_scripts":  fmap(marshalScript)(a.PreScripts),
		"post_scripts": fmap(marshalScript)(a.PostScripts),
	}
}

func marshalInterest(i Interest) map[string]any {
	target := i.Target().Baseline.(mod.Tagger).Tag()
	switch i := i.(type) {
	case *AttackInterest:
		return map[string]any{
			"verb":     "attack",
			"target":   target,
			"damage":   i.Damage,
			"defense":  i.Defense,
			"loss":     i.Loss,
			"overflow": i.Overflow,
			"health":   i.Health,
		}

	case *HealingInterest:
		return map[string]any{
			"verb":     "heal",
			"target":   target,
			"healing":  i.Healing,
			"overflow": i.Overflow,
			"health":   i.Health,
		}

	case *BuffingInterest:
		return map[string]any{
			"verb":   "buff",
			"target": target,
			"buff":   i.Buff.Tag(),
		}

	case *PurgingInterest:
		return map[string]any{
			"verb":   "purge",
			"target": target,
		}

	default:
		panic("bad interest type")
	}
}
