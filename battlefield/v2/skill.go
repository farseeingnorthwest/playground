package battlefield

var (
	NormalAttack = &LaunchReactor{
		actors: []Actor{
			SelectiveActor{
				AndSelector{
					HealthSelector{},
					SideSelector{},
					RandomSelector{1},
				},
				&Attacker{Evaluator{Damage, 100}},
			},
		},
	}
	coordination     = "coordination"
	coordinationBuff = NewTaggedBuff(
		coordination,
		[]*EvaluationBuff{
			NewEvaluationBuff(Damage, EvaluationSlope(108)),
			NewEvaluationBuff(Defense, EvaluationSlope(108)),
		},
		TaggedCapacity(3),
	)
	healingBuff = &RoundStartReactor{
		FiniteReactor: &FiniteReactor{3},
		actors: []Actor{
			SelectiveActor{
				CurrentSelector{},
				&Healer{Evaluator{Damage, 150}}, // TODO:
			},
		},
	}
	Active = []*LaunchReactor{
		{
			actors: []Actor{
				SelectiveActor{
					AndSelector{
						HealthSelector{},
						SideSelector{},
					},
					&Attacker{Evaluator{Damage, 90}},
				},
				SelectiveActor{
					AndSelector{
						HealthSelector{},
						SideSelector{true},
					},
					&Buffer{coordinationBuff},
				},
			},
		},
		{
			actors: []Actor{
				SelectiveActor{
					AndSelector{
						HealthSelector{},
						SideSelector{true},
						RandomSelector{3},
					},
					&Buffer{healingBuff},
				},
			},
		},
	}
)

type LaunchReactor struct {
	PeriodicReactor
	actors []Actor
}

func (a *LaunchReactor) React(signal Signal) {
	switch sig := signal.(type) {
	case *LaunchSignal:
		if !a.Free() {
			return
		}

		var actions []*Action
		for _, actor := range a.actors {
			a := actor.Act(sig.Target, sig.Field.fighters)
			if a == nil {
				return
			}

			actions = append(actions, a)
		}

		sig.Append(actions...)
		sig.Launched = true
		a.WarmUp()

	case *RoundEndSignal:
		a.CoolDown()
	}
}

type RoundStartReactor struct {
	*FiniteReactor
	PeriodicReactor
	actors []Actor
}

func (a *RoundStartReactor) React(signal Signal) {
	switch sig := signal.(type) {
	case *RoundStartSignal:
		if !a.Free() {
			return
		}

		var actions []*Action
		for _, actor := range a.actors {
			a := actor.Act(sig.Current, sig.Field.fighters)
			if a == nil {
				return
			}

			actions = append(actions, a)
		}

		sig.Append(actions...)
		a.FiniteReactor.WarmUp()
		a.PeriodicReactor.WarmUp()

	case *RoundEndSignal:
		a.CoolDown()
	}
}

func (a *RoundStartReactor) Fork() any {
	return &RoundStartReactor{
		FiniteReactor:   a.FiniteReactor.Fork(),
		PeriodicReactor: a.PeriodicReactor,
		actors:          a.actors,
	}
}
