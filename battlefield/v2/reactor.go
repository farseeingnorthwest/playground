package battlefield

type Reactor interface {
	React(Signal)
}

type Validator interface {
	Validate() bool
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

	targets := a.Select(launch.Target, launch.Field.Warriors())
	if len(targets) == 0 {
		return
	}

	launch.Add(&Action{
		Source:  launch.Target,
		Targets: targets,
		Verb:    NewAttack(a.Points),
	})
}
