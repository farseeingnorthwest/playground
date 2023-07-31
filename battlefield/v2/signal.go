package battlefield

type Signal interface {
	signalTrait()
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

func (a *actions) signalTrait() {}

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

type RoundStartSignal struct {
	Current *Fighter
	Field   *BattleField
	actions
}

func (*RoundStartSignal) signalTrait() {}

type RoundEndSignal struct{}

func (*RoundEndSignal) signalTrait() {}

type actionSignal struct {
	*Action
	actions
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

type Axis uint8

const (
	Damage Axis = iota
	Defense
	Loss
	Healing
	Health
	Speed
)

type EvaluationSignal struct {
	axis   Axis
	value  int
	action *Action
}

func NewEvaluationSignal(axis Axis, value int, action *Action) *EvaluationSignal {
	return &EvaluationSignal{
		axis,
		value,
		action,
	}
}

func (s *EvaluationSignal) Axis() Axis {
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

func (s *EvaluationSignal) signalTrait() {}
