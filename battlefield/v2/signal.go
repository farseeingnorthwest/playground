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

type LaunchingSignal struct {
	Subject Warrior
	Objects []Warrior
	actions
}

func NewLaunchingSignal(subject Warrior, objects []Warrior) *LaunchingSignal {
	return &LaunchingSignal{
		Subject: subject,
		Objects: objects,
	}
}

type PreActionSignal struct {
	*Action
	actions
}

func NewPreActionSignal(action *Action) *PreActionSignal {
	return &PreActionSignal{
		Action: action,
	}
}

type clearingSignal struct {
	value int
}

func (s *clearingSignal) Value() int {
	return s.value
}

func (s *clearingSignal) SetValue(value int) {
	s.value = value
}

func (s *clearingSignal) Map(fn ...func(int) int) {
	for _, f := range fn {
		s.value = f(s.value)
	}
}

func (s *clearingSignal) signalTrait() {}

type AttackClearingSignal struct {
	clearingSignal
}

func NewAttackClearingSignal(value int) *AttackClearingSignal {
	return &AttackClearingSignal{
		clearingSignal: clearingSignal{
			value: value,
		},
	}
}

type DefenseClearingSignal struct {
	clearingSignal
}

func NewDefenseClearingSignal(value int) *DefenseClearingSignal {
	return &DefenseClearingSignal{
		clearingSignal: clearingSignal{
			value: value,
		},
	}
}

type DamageClearingSignal struct {
	clearingSignal
}

func NewDamageClearingSignal(value int) *DamageClearingSignal {
	return &DamageClearingSignal{
		clearingSignal: clearingSignal{
			value: value,
		},
	}
}
