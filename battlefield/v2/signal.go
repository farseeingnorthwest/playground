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

func (a *actions) Add(action *Action) {
	a.actions = append(a.actions, action)
}

func (a *actions) signalTrait() {}

type LaunchSignal struct {
	Target *Fighter
	Field  *BattleField
	actions
}

func NewLaunchingSignal(target *Fighter, field *BattleField) *LaunchSignal {
	return &LaunchSignal{
		Target: target,
		Field:  field,
	}
}

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
	Attack Axis = iota
	Defense
	Damage
	Health
	Speed
)

type EvaluationSignal struct {
	axis  Axis
	clear bool
	value int
}

func NewEvaluationSignal(axis Axis, clear bool, value int) *EvaluationSignal {
	return &EvaluationSignal{
		axis,
		clear,
		value,
	}
}

func (s *EvaluationSignal) Axis() Axis {
	return s.axis
}

func (s *EvaluationSignal) Clear() bool {
	return s.clear
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

func (s *EvaluationSignal) signalTrait() {}
