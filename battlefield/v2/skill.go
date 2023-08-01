package battlefield

import "github.com/farseeingnorthwest/playground/battlefield/v2/evaluation"

var (
	NormalAttack = &LaunchReactor{
		actors: []Actor{
			SelectiveActor{
				AndSelector{
					HealthSelector{},
					SideSelector{},
					RandomSelector{1},
				},
				&BlindActor{
					NewAttackProto(Head),
					EvalChain{
						&PercentageEvaluator{
							Damage,
							100,
						},
						nil,
					},
				},
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
				&BlindActor{
					NewHealProto(Head),
					EvalChain{},
				},
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
					&BlindActor{
						NewAttackProto(Head),
						EvalChain{
							&PercentageEvaluator{Damage, 90},
							nil,
						},
					},
				},
				SelectiveActor{
					AndSelector{
						HealthSelector{},
						SideSelector{true},
					},
					&BlindActor{
						NewBuffProto(coordinationBuff, nil),
						EvalChain{},
					},
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
					&BlindActor{
						NewBuffProto(healingBuff, nil),
						EvalChain{},
					},
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

func (a *LaunchReactor) Fork(*evaluation.Block, Signal) Reactor {
	return &LaunchReactor{
		PeriodicReactor: a.PeriodicReactor,
		actors:          a.actors,
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

func (a *RoundStartReactor) Fork(*evaluation.Block, Signal) Reactor {
	return &RoundStartReactor{
		FiniteReactor:   a.FiniteReactor.Fork(),
		PeriodicReactor: a.PeriodicReactor,
		actors:          a.actors,
	}
}
