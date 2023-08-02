package battlefield

import (
	"github.com/farseeingnorthwest/playground/battlefield/v2/evaluation"
)

var (
	Layer0 = []ForkableReactor{
		// 0
		NewTaggedBuff(
			coordination,
			[]*EvaluationBuff{
				NewEvaluationBuff(evaluation.Damage, EvaluationMultiplier(108)),
				NewEvaluationBuff(evaluation.Defense, EvaluationMultiplier(108)),
			},
			TaggedCapacity(3),
		),
		// 1
		&RoundStartReactor{NewModifiedReactor(
			[]Actor{
				&SelectiveActor{
					CurrentSelector{},
					&BlindActor{
						NewHealProto(evaluation.Head),
						nil,
					},
				},
			},
			Capacity(3),
		)},
		// 2
		NewEvaluationBuff(evaluation.Defense, EvaluationMultiplier(120), EvaluationCapacity(2)),
		// 3
		NewEvaluationBuff(evaluation.Defense, EvaluationMultiplier(80), EvaluationCapacity(2)),
		// 4
		NewEvaluationBuff(evaluation.Defense, EvaluationMultiplier(95)),
		// 5
		NewEvaluationBuff(evaluation.Defense, EvaluationMultiplier(0)),
		// 6
		NewEvaluationBuff(evaluation.HealthMax, EvaluationMultiplier(50)),
		// 7
		NewEvaluationBuff(evaluation.Damage, EvaluationMultiplier(105)),
		// 8
		NewEvaluationBuff(evaluation.Damage, EvaluationMultiplier(130)),
		// 9
		NewEvaluationBuff(evaluation.Defense, EvaluationMultiplier(110)),
		// 10
		NewEvaluationBuff(evaluation.HealthMax, EvaluationMultiplier(110)),
	}
	Active = []*LaunchReactor{
		{NewModifiedReactor([]Actor{
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
		})},
		{NewModifiedReactor([]Actor{
			&SelectiveActor{
				AndSelector{
					HealthSelector{},
					SideSelector{},
				},
				&BlindActor{
					NewAttackProto(evaluation.Head),
					evaluation.NewBundleProto(
						evaluation.NewAxisEvaluator(
							evaluation.Damage,
							90,
						),
					),
				},
			},
			&SelectiveActor{
				AndSelector{
					HealthSelector{},
					SideSelector{true},
				},
				&BlindActor{
					NewBuffProto(Layer0[0], nil),
					nil,
				},
			},
		}, Period(4))},
		{NewModifiedReactor([]Actor{
			&SelectiveActor{
				AndSelector{
					HealthSelector{},
					SideSelector{true},
					RandomSelector{3},
				},
				&BlindActor{
					NewBuffProto(Layer0[1], nil),
					evaluation.NewBundleProto(
						evaluation.NewAxisEvaluator(
							evaluation.Damage,
							150,
						),
					),
				},
			},
		}, Period(4))},
		{NewModifiedReactor([]Actor{
			&SelectiveActor{
				AndSelector{
					HealthSelector{},
					SideSelector{true},
				},
				&BlindActor{
					NewBuffProto(Layer0[2], nil),
					nil,
				},
			},
		}, Period(4))},
		{NewModifiedReactor([]Actor{
			&SelectiveActor{
				AndSelector{
					HealthSelector{},
					SideSelector{true},
				},
				&BlindActor{
					NewBuffProto(Layer0[2], nil),
					nil,
				},
			},
		}, Period(4))},
		{NewModifiedReactor([]Actor{
			&SelectiveActor{
				AndSelector{
					HealthSelector{},
					SideSelector{},
				},
				&BlindActor{
					NewBuffProto(Layer0[3], nil),
					nil,
				},
			},
		}, Period(4))},
		{NewModifiedReactor([]Actor{
			&SelectiveActor{
				CurrentSelector{},
				&BlindActor{
					NewAttackProto(evaluation.Head),
					evaluation.NewBundleProto(
						evaluation.NewAxisEvaluator(
							evaluation.Damage,
							25,
						),
					),
				},
			},
			&SelectiveActor{
				AndSelector{
					HealthSelector{},
					SideSelector{},
				},
				&BlindActor{
					NewAttackProto(evaluation.Head),
					evaluation.NewBundleProto(
						evaluation.NewAxisEvaluator(
							evaluation.Damage,
							500,
						),
					),
				},
			},
		}, Period(4))},
		{NewModifiedReactor([]Actor{
			&SelectiveActor{
				CurrentSelector{},
				&BlindActor{
					NewAttackProto(evaluation.Head),
					evaluation.NewBundleProto(
						evaluation.NewAxisEvaluator(
							evaluation.Damage,
							50,
						),
					),
				},
			},
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
							1000,
						),
					),
				},
			},
		}, Period(4))},
		{NewModifiedReactor([]Actor{
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
							175,
						),
					),
				},
			},
		})},
		{NewModifiedReactor([]Actor{
			&SelectiveActor{
				AndSelector{
					HealthSelector{},
					SideSelector{},
				},
				&BlindActor{
					NewAttackProto(evaluation.Head),
					evaluation.NewBundleProto(
						evaluation.NewAxisEvaluator(
							evaluation.Damage,
							70,
						),
					),
				},
			},
		})},
	}
	Passive = []Reactor{
		&RoundStartReactor{NewModifiedReactor([]Actor{
			&SelectiveActor{
				AndSelector{
					HealthSelector{},
					SideSelector{},
					AxisSelector{evaluation.Defense, false},
				},
				&BlindActor{
					NewBuffProto(Layer0[4], nil),
					nil,
				},
			},
		})},
		&RoundStartReactor{NewModifiedReactor([]Actor{
			&SelectiveActor{
				CurrentSelector{},
				&BlindActor{
					NewBuffProto(Layer0[7], nil),
					nil,
				},
			},
		})},
	}

	NormalAttack = Active[0]
	coordination = "coordination"
)

func NewCriticalAttack(rng Rng, odds int, multiplier int) Reactor {
	return &PreAttackReactor{NewModifiedReactor([]Actor{
		&ProbabilityActor{
			rng,
			odds,
			&BlindActor{
				NewBuffProto(
					NewClearingBuff(
						evaluation.Loss,
						nil,
						ClearingMultiplier(multiplier),
					),
					nil,
				),
				nil,
			},
		},
	})}
}
