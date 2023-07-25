package battlefield

type Signal interface {
}

type Launch struct {
	subject Warrior
	objects []Warrior
	actions []*Action
}

type Reactor interface {
	React(Signal)
	Valid() bool
}

type NormalAttack struct {
	Selector
	Points int
}

func (a *NormalAttack) React(signal Signal) {
	launch, ok := signal.(*Launch)
	if !ok {
		return
	}

	fighter := a.Select(launch.subject, launch.objects)
	if fighter == nil {
		return
	}

	launch.actions = append(launch.actions, &Action{
		Subject: launch.subject,
		Objects: []Warrior{fighter},
		Verb:    &Attack{Points: a.Points},
	})
}

func (a *NormalAttack) Valid() bool {
	return true
}
