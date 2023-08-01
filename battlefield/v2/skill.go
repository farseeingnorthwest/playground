package battlefield

import (
	"github.com/farseeingnorthwest/playground/battlefield/v2/evaluation"
	"github.com/farseeingnorthwest/playground/battlefield/v2/modifier"
)

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
					NewAttackProto(evaluation.Head),
					evaluation.NewBundleProto(
						evaluation.NewAxisEvaluator(
							evaluation.Damage,
							100,
						),
					),
				},
			},
		},
	}
	coordination     = "coordination"
	coordinationBuff = NewTaggedBuff(
		coordination,
		[]*EvaluationBuff{
			NewEvaluationBuff(evaluation.Damage, EvaluationMultiplier(108)),
			NewEvaluationBuff(evaluation.Defense, EvaluationMultiplier(108)),
		},
		TaggedCapacity(3),
	)
	healingBuff = &RoundStartReactor{
		FiniteModifier: modifier.NewFiniteModifier(3),
		actors: []Actor{
			SelectiveActor{
				CurrentSelector{},
				&BlindActor{
					NewHealProto(evaluation.Head),
					&evaluation.Bundle{},
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
						NewAttackProto(evaluation.Head),
						evaluation.NewBundleProto(
							evaluation.NewAxisEvaluator(evaluation.Damage, 90),
						),
					},
				},
				SelectiveActor{
					AndSelector{
						HealthSelector{},
						SideSelector{true},
					},
					&BlindActor{
						NewBuffProto(coordinationBuff, nil),
						&evaluation.Bundle{},
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
						&evaluation.Bundle{},
					},
				},
			},
		},
	}
)

type LaunchReactor struct {
	modifier.PeriodicModifier
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
		PeriodicModifier: a.PeriodicModifier,
		actors:           a.actors,
	}
}

type RoundStartReactor struct {
	*modifier.FiniteModifier
	modifier.PeriodicModifier
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
		a.FiniteModifier.WarmUp()
		a.PeriodicModifier.WarmUp()

	case *RoundEndSignal:
		a.CoolDown()
	}
}

func (a *RoundStartReactor) Fork(*evaluation.Block, Signal) Reactor {
	return &RoundStartReactor{
		FiniteModifier:   a.FiniteModifier.Clone().(*modifier.FiniteModifier),
		PeriodicModifier: a.PeriodicModifier,
		actors:           a.actors,
	}
}
