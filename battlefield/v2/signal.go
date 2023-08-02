package battlefield

import "github.com/farseeingnorthwest/playground/battlefield/v2/evaluation"

type Signal interface {
	Current() any
}

type actions struct {
	actions []*Action
}

func (a *actions) Actions() []*Action {
	return a.actions
}

func (a *actions) Append(actions ...*Action) {
	a.actions = append(a.actions, actions...)
}

type LaunchSignal struct {
	Target   *Fighter
	Field    *BattleField
	Launched bool

	actions
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
	actions
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

type actionSignal struct {
	current any
	*Action
	actions
}

func (s *actionSignal) Current() any {
	return s.current
}

type PreActionSignal struct {
	actionSignal
}

func NewPreActionSignal(action *Action) *PreActionSignal {
	return &PreActionSignal{
		actionSignal: actionSignal{
			Action: action,
		},
	}
}

type PostActionSignal struct {
	actionSignal
}

func NewPostActionSignal(action *Action) *PostActionSignal {
	return &PostActionSignal{
		actionSignal: actionSignal{
			Action: action,
		},
	}
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
