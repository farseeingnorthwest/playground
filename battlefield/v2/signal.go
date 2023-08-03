package battlefield

import "github.com/farseeingnorthwest/playground/battlefield/v2/evaluation"

type Signal interface {
	Current() any
}

type scripts struct {
	scripts []*Script
}

func (a *scripts) Scripts() []*Script {
	return a.scripts
}

func (a *scripts) Append(scripts ...*Script) {
	a.scripts = append(a.scripts, scripts...)
}

type LaunchSignal struct {
	Target   *Fighter
	Field    *BattleField
	Launched bool

	scripts
}

func NewLaunchSignal(target *Fighter, field *BattleField) *LaunchSignal {
	return &LaunchSignal{
		Target: target,
		Field:  field,
	}
}

func (s *LaunchSignal) Current() any {
	return s.Target
}

type RoundStartSignal struct {
	current *Fighter
	Field   *BattleField
	scripts
}

func (s *RoundStartSignal) Current() any {
	return s.current
}

type RoundEndSignal struct {
	current any
}

func (s *RoundEndSignal) Current() any {
	return s.current
}

type PreScriptSignal struct {
	current any
	*Script
}

func (s *PreScriptSignal) Current() any {
	return s.current
}

func (s *PreScriptSignal) Fork(current any) any {
	return &PreScriptSignal{current, s.Script}
}

func (s *PreScriptSignal) Scripts() []*Script {
	return nil
}

type PostScriptSignal struct {
	current any
	*Script
}

func (s *PostScriptSignal) Current() any {
	return s.current
}

func (s *PostScriptSignal) Fork(current any) any {
	return &PostScriptSignal{current, s.Script}
}

func (s *PostScriptSignal) Scripts() []*Script {
	return nil
}

type ActionSignal interface {
	Signal
	Fork(current any) any
	Scripts() []*Script
}

type actionSignal struct {
	current any
	*Action
	scripts
}

func (s *actionSignal) Current() any {
	return s.current
}

func (s *actionSignal) Fork(current any) any {
	return &actionSignal{
		current: current,
		Action:  s.Action,
	}
}

func (s *actionSignal) Scripts() []*Script {
	return s.scripts.scripts
}

type PreActionSignal struct {
	*actionSignal
}

func NewPreActionSignal(action *Action) *PreActionSignal {
	return &PreActionSignal{
		actionSignal: &actionSignal{
			Action: action,
		},
	}
}

func (s *PreActionSignal) Fork(current any) any {
	return &PreActionSignal{s.actionSignal.Fork(current).(*actionSignal)}
}

type PostActionSignal struct {
	*actionSignal
}

func NewPostActionSignal(action *Action) *PostActionSignal {
	return &PostActionSignal{
		actionSignal: &actionSignal{
			Action: action,
		},
	}
}

func (s *PostActionSignal) Fork(current any) any {
	return &PostActionSignal{s.actionSignal.Fork(current).(*actionSignal)}
}

type EvaluationSignal struct {
	axis   evaluation.Axis
	value  int
	action *Action
}

func NewEvaluationSignal(axis evaluation.Axis, value int, action *Action) *EvaluationSignal {
	return &EvaluationSignal{
		axis,
		value,
		action,
	}
}

func (s *EvaluationSignal) Current() any {
	return nil
}

func (s *EvaluationSignal) Axis() evaluation.Axis {
	return s.axis
}

func (s *EvaluationSignal) Value() int {
	return s.value
}

func (s *EvaluationSignal) SetValue(value int) {
	s.value = value
}

func (s *EvaluationSignal) Map(fn ...func(int) int) {
	for _, f := range fn {
		s.value = f(s.value)
	}
}

func (s *EvaluationSignal) Action() *Action {
	return s.action
}
