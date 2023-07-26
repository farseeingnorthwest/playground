package battlefield

type Reactor interface {
	React(Signal)
	Valid() bool
}

type NormalAttack struct {
	Selector
	Points int
}

func (a *NormalAttack) React(signal Signal) {
	launch, ok := signal.(*LaunchingSignal)
	if !ok {
		return
	}

	fighter := a.Select(launch.Subject, launch.Objects)
	if fighter == nil {
		return
	}

	launch.Add(&Action{
		Subject: launch.Subject,
		Objects: []Warrior{fighter},
		Verb:    &Attack{Points: a.Points},
	})
}

func (a *NormalAttack) Valid() bool {
	return true
}
