package battlefield

import (
	"github.com/farseeingnorthwest/playground/battlefield/v2/evaluation"
)

var (
	NormalAttack = &LaunchReactor{
		NewModifiedReactor([]Actor{
			&SelectiveActor{
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
		}),
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
		NewModifiedReactor(
			[]Actor{
				&SelectiveActor{
					CurrentSelector{},
					&BlindActor{
						NewHealProto(evaluation.Head),
						&evaluation.Bundle{},
					},
				},
			}, Capacity(3),
		),
	}
	Active = []*LaunchReactor{
		{
			NewModifiedReactor([]Actor{
				&SelectiveActor{
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
				&SelectiveActor{
					AndSelector{
						HealthSelector{},
						SideSelector{true},
					},
					&BlindActor{
						NewBuffProto(coordinationBuff, nil),
						&evaluation.Bundle{},
					},
				},
			}),
		},
		{
			NewModifiedReactor([]Actor{
				&SelectiveActor{
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
			}),
		},
	}
)
