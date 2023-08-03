package battlefield

import (
	"github.com/farseeingnorthwest/playground/battlefield/v2/evaluation"
)

var (
	Layer0 = []ForkableReactor{
		// 0
		NewCompoundBuff(
			coordination,
			[]*EvaluationBuff{
				NewEvaluationBuff("提升攻击力 8%", evaluation.Damage, EvaluationMultiplier(108)),
				NewEvaluationBuff("提升防御力 8%", evaluation.Defense, EvaluationMultiplier(108)),
			},
			TaggedCapacity(3),
		),
		// 1
		&RoundStartReactor{NewModifiedReactor(
			"回合開始時治療",
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
		NewEvaluationBuff("提升防御力 20%", evaluation.Defense, EvaluationMultiplier(120), EvaluationCapacity(2)),
		// 3
		NewEvaluationBuff("降低防御力 20%", evaluation.Defense, EvaluationMultiplier(80), EvaluationCapacity(2)),
		// 4
		NewEvaluationBuff("降低防御力 5%", evaluation.Defense, EvaluationMultiplier(95)),
		// 5
		NewEvaluationBuff("降低防御力 100%", evaluation.Defense, EvaluationMultiplier(0)),
		// 6
		NewEvaluationBuff("降低生命值上限 50%", evaluation.HealthMax, EvaluationMultiplier(50)),
		// 7
		NewEvaluationBuff("提升攻击力 5%", evaluation.Damage, EvaluationMultiplier(105)),
		// 8
		NewEvaluationBuff("提升攻击力 30%", evaluation.Damage, EvaluationMultiplier(130)),
		// 9
		NewEvaluationBuff("提升防御力 10%", evaluation.Defense, EvaluationMultiplier(110)),
		// 10
		NewEvaluationBuff("提升生命值上限 10%", evaluation.HealthMax, EvaluationMultiplier(110)),
	}
	Active = []*LaunchReactor{
		{NewModifiedReactor("普通攻击", []Actor{
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
		{NewModifiedReactor("攻擊敵方全體 90% 自身傷害，並給予全隊共鬥效果，持續三回合", []Actor{
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
		{NewModifiedReactor("賦予隨機三位友方回合開始時治療攻擊力 150% 治療效果，持續三回合", []Actor{
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
		{NewModifiedReactor("提升全隊防禦力 20%，持續 2 回合", []Actor{
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
		{NewModifiedReactor("降低敵方全體防禦 20%，持續 2 回合", []Actor{
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
		{NewModifiedReactor("對自己造成 25% 自身傷害，對敵方全體造成 500% 自身傷害", []Actor{
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
		{NewModifiedReactor("對自己造成 50% 自身傷害，對敵方隨機一位造成 1000% 自身傷害", []Actor{
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
		{NewModifiedReactor("對敵方一位造成 175% 自身傷害", []Actor{
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
		{NewModifiedReactor("攻擊敵方全體 70% 傷害", []Actor{
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
		&RoundStartReactor{NewModifiedReactor("每 1 回合降低敵方防禦最高的單位防禦 5%，可疊加，不可被清除", []Actor{
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
		&RoundStartReactor{NewModifiedReactor("每回合提升傷害 5%，可疊加，不可被清除", []Actor{
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
	return &PreAttackReactor{NewModifiedReactor("暴击", []Actor{
		&ProbabilityActor{
			rng,
			odds,
			&BlindActor{
				NewBuffProto(
					NewClearingBuff(
						"暴击伤害",
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
