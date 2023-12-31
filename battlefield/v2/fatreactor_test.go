package battlefield_test

import . "github.com/farseeingnorthwest/playground/battlefield/v2"

const (
	Water Element = iota
	Fire
	Ice
	Wind
	Earth
	Thunder
	Dark
	Light
)

var (
	RngX          = &RngProxy{}
	SkillGroup    = ExclusionGroup(0)
	Shuffle       = NewShuffleSelector(RngX, Label("Taunt"))
	ElementTheory = NewTheoryActor(
		map[Element]map[Element]int{
			Water:   {Fire: 120, Thunder: 80},
			Fire:    {Ice: 120, Water: 80},
			Ice:     {Wind: 120, Fire: 80},
			Wind:    {Earth: 120, Ice: 80},
			Earth:   {Thunder: 120, Wind: 80},
			Thunder: {Water: 120, Earth: 80},
			Dark:    {Light: 120},
			Light:   {Dark: 120},
		},
	)
	Regular = []*FatReactor{
		NewFatReactor(
			FatTags(SkillGroup, Label("NormalAttack")),
			FatRespond(
				NewSignalTrigger(&LaunchSignal{}),
				NewSelectActor(
					NewVerbActor(&Attack{}, AxisEvaluator(Damage)),
					SideSelector(false),
					Healthy,
					Shuffle,
					FrontSelector(1),
				),
			),
		),
		NewFatReactor(
			FatTags(Priority(1000)),
			FatRespond(
				NewFatTrigger(
					&PreActionSignal{},
					CurrentIsSourceTrigger{},
					NewVerbTrigger[*Attack](),
				),
				NewProbabilityActor(
					RngX,
					AxisEvaluator(CriticalOdds),
					CriticalActor{},
				),
			),
		),
		NewFatReactor(
			FatTags(Priority(999)),
			FatRespond(
				NewFatTrigger(
					&PreActionSignal{},
					CurrentIsSourceTrigger{},
					NewVerbTrigger[*Attack](),
					CriticalStrikeTrigger{},
				),
				NewActionBuffer(
					AxisEvaluator(CriticalLoss),
					NewBuffer(Loss, false, nil),
				),
			),
		),
		NewFatReactor(
			FatRespond(
				NewFatTrigger(
					&PreActionSignal{},
					NewVerbTrigger[*Attack](),
				),
				NewActionBuffer(nil, ElementTheory),
			),
		),
	}

	Effect = map[string]*FatReactor{
		// 沉睡
		"Sleep": NewFatReactor(
			FatTags(SkillGroup, Priority(10), Label("Sleep")),
			FatRespond(
				NewSignalTrigger(&LaunchSignal{}),
				NewSequenceActor(),
			),
			FatCapacity(
				NewAnyTrigger(
					NewSignalTrigger(&RoundEndSignal{}),
					NewFatTrigger(
						&PostActionSignal{},
						VerbTrigger[*Attack]{},
						CurrentIsTargetTrigger{}),
				),
				1,
			),
		),

		// 暈眩
		"Stun": NewFatReactor(
			FatTags(SkillGroup, Priority(10), Label("Stun")),
			FatRespond(
				NewSignalTrigger(&LaunchSignal{}),
				NewSequenceActor(),
			),
			FatCapacity(
				NewSignalTrigger(&RoundEndSignal{}),
				1,
			),
		),

		// 多重屏障
		"Barrier": NewFatReactor(
			FatTags(Label("Barrier")),
			FatRespond(
				NewSignalTrigger(&PreLossSignal{}),
				NewLossStopper(NewMultiplier(10, AxisEvaluator(HealthMaximum)), true),
			),
			FatCapacity(nil, 1),
		),

		// 嘲諷
		"Taunt": NewFatReactor(
			FatTags(Label("Taunt")),
			FatCapacity(NewSignalTrigger(&RoundEndSignal{}), 2),
		),

		// 庇護
		"Sanctuary": NewFatReactor(
			FatTags(Label("Sanctuary"), Interest{}),
			FatRespond(
				NewSignalTrigger(&PreLossSignal{}),
				NewLossStopper(NewMultiplier(30, AxisEvaluator(HealthMaximum)), false),
			),
		),

		// 護盾
		"Shield": NewFatReactor(
			FatTags(Label("Shield")),
			FatRespond(
				NewSignalTrigger(&PreLossSignal{}),
				NewSelectActor(
					LossResister{},
					CurrentSelector{},
				),
			),
		),

		// 增益無效
		"BuffImmune": NewFatReactor(
			FatTags(Label("BuffImmune")),
			FatCapacity(NewSignalTrigger(&RoundEndSignal{}), 3),
			FatRespond(
				NewFatTrigger(
					&PreActionSignal{},
					CurrentIsTargetTrigger{},
					NewTagTrigger("Buff"),
				),
				NewSelectActor(ImmuneActor{}, CurrentSelector{}),
			),
		),

		// 控制效果免疫
		"ControlImmune": NewFatReactor(
			FatTags(Label("ControlImmune")),
			FatCapacity(NewSignalTrigger(&RoundEndSignal{}), 2),
			FatRespond(
				NewFatTrigger(
					&PreActionSignal{},
					CurrentIsTargetTrigger{},
					NewTagTrigger("Control"),
				),
				NewSelectActor(ImmuneActor{}, CurrentSelector{}),
			),
		),

		// 再生
		"Regeneration": NewFatReactor(
			FatTags(Label("Regeneration")),
			FatRespond(
				NewSignalTrigger(&RoundStartSignal{}),
				NewSelectActor(
					NewVerbActor(&Heal{}, NewMultiplier(7, AxisEvaluator(HealthMaximum))),
					CurrentSelector{},
				),
			),
			FatCapacity(nil, 1),
		),
	}

	Special = [][]*FatReactor{
		// ////////////////////////////////////////////////////////////
		// [0] 織田
		{
			// 對隨機 1 名敵人進行 3 次攻擊，每次造成攻擊力 460% 的傷害。
			NewFatReactor(
				FatTags(SkillGroup, Priority(1), Label("@Launch({1} 3 * 460% Damage)")),
				FatRespond(
					NewSignalTrigger(&LaunchSignal{}),
					NewSelectActor(
						NewRepeatActor(
							3,
							NewVerbActor(&Attack{}, NewMultiplier(460, AxisEvaluator(Damage))),
						),
						SideSelector(false),
						Healthy,
						Shuffle,
						FrontSelector(1),
					),
				),
				FatCooling(NewSignalTrigger(&RoundEndSignal{}), 5),
			),

			// 對全體敵人造成攻擊力 480% 的傷害。再對當前生命值最低的1名敵人造成攻擊力 520% 的傷害。
			NewFatReactor(
				FatTags(SkillGroup, Priority(2), Label("@Launch({*} 480% Damage; {1} 520% Damage)")),
				FatRespond(
					NewSignalTrigger(&LaunchSignal{}),
					NewSelectActor(
						NewVerbActor(&Attack{}, NewMultiplier(480, AxisEvaluator(Damage))),
						SideSelector(false),
						Healthy,
					),
					NewSelectActor(
						NewVerbActor(&Attack{}, NewMultiplier(520, AxisEvaluator(Damage))),
						SideSelector(false),
						Healthy,
						NewSortSelector(HealthPercent, true),
						FrontSelector(1),
					),
				),
				FatCooling(NewSignalTrigger(&RoundEndSignal{}), 4),
			),

			// 每次行動開始時，提升 2% 爆擊率(最高 15 層，無法被解除)。
			NewFatReactor(
				FatTags(Priority(3), Label("@Launch([15] +2% CriticalOdds)")),
				FatRespond(
					NewSignalTrigger(&LaunchSignal{}),
					NewSelectActor(
						NewVerbActor(
							NewBuff(false, nil, NewBuffReactor(
								CriticalOdds,
								true,
								ConstEvaluator(2),
								FatTags(NewStackingLimit(15), Label("[15] +2% CriticalOdds")))),
							nil,
						),
						CurrentSelector{},
					),
				),
			),

			// 提升 20% 攻擊力。(無法被解除)
			NewFatReactor(
				FatTags(Priority(4), Label("@BattleStart({$} +20% Damage)")),
				FatRespond(
					NewSignalTrigger(&BattleStartSignal{}),
					NewSelectActor(
						NewVerbActor(
							NewBuff(false, nil, NewBuffReactor(
								Damage,
								false,
								ConstEvaluator(120),
								FatTags(Label("+20% Damage")))),
							nil,
						),
						CurrentSelector{},
					),
				),
			),
		},

		// ////////////////////////////////////////////////////////////
		// [1] 豐臣
		{
			// 對隨機 2 名敵人造成攻擊力 420% 的傷害。並對目標附加 30% 被擊增傷(1 回合)&「沉睡」(1 回合)。
			NewFatReactor(
				FatTags(SkillGroup, Priority(1), Label("@Launch({2} 420% Damage, +30% Loss, Sleep)")),
				FatRespond(
					NewSignalTrigger(&LaunchSignal{}),
					NewSelectActor(
						NewSequenceActor(
							NewVerbActor(&Attack{}, NewMultiplier(420, AxisEvaluator(Damage))),
							NewVerbActor(
								// 30% 被擊增傷 (1回合)
								NewBuff(false, nil, NewBuffReactor(
									Loss,
									false,
									ConstEvaluator(130),
									FatTags(Label("+30% Loss")),
									FatCapacity(NewSignalTrigger(&RoundEndSignal{}), 1))),
								nil,
							),
							NewVerbActor(NewBuff(false, nil, Effect["Sleep"]), nil),
						),
						SideSelector(false),
						Healthy,
						Shuffle,
						FrontSelector(2),
					),
				),
				FatCooling(NewSignalTrigger(&RoundEndSignal{}), 5),
			),

			// 對攻擊力最高的 1 名敵人造成攻擊力 520% 的傷害。並有 70% 機率對目標附加「暈眩」(1 回合)。
			NewFatReactor(
				FatTags(SkillGroup, Priority(2), Label("@Launch({1} 520% Damage, P(70%, Stun))")),
				FatRespond(
					NewSignalTrigger(&LaunchSignal{}),
					NewSelectActor(
						NewSequenceActor(
							NewVerbActor(&Attack{}, NewMultiplier(520, AxisEvaluator(Damage))),
							NewProbabilityActor(RngX, ConstEvaluator(70), NewVerbActor(
								NewBuff(false, nil, Effect["Stun"]),
								nil,
							)),
						),
						SideSelector(false),
						Healthy,
						NewSortSelector(Damage, false),
						FrontSelector(1),
					),
				),
				FatCooling(NewSignalTrigger(&RoundEndSignal{}), 4),
			),

			// 自身附帶 2 種以上減益效果時觸發，解除所有減益效果。
			NewFatReactor(
				FatTags(Priority(2), Label("@PostAction({$/<Nerf> >= 2} Purge())")),
				FatRespond(
					NewAnyTrigger(
						NewSignalTrigger(&RoundStartSignal{}),
						NewSignalTrigger(&PostActionSignal{}),
					),
					NewSelectActor(
						NewVerbActor(NewPurge(RngX, "Nerf", 0), nil),
						CurrentSelector{},
						NewWaterLevelSelector(Ge, NewBuffCounter("Nerf"), 2),
					),
				),
				FatCooling(NewSignalTrigger(&RoundEndSignal{}), 4),
			),

			// 行動開始時獲得「多重屏障」(2 層，無法被解除)。
			NewFatReactor(
				FatTags(Priority(3), Label("@Launch({$} Barrier)")),
				FatRespond(
					NewSignalTrigger(&LaunchSignal{}),
					NewSelectActor(
						NewVerbActor(NewBuff(false, nil, Effect["Barrier"]), nil),
						CurrentSelector{},
					),
				),
				FatCooling(NewSignalTrigger(&RoundEndSignal{}), 4),
			),
		},

		// ////////////////////////////////////////////////////////////
		// [2] 上杉
		{
			// 對當前生命值百分比最高的 1 名敵人造成攻擊力 400% 的傷害；此攻擊必定爆擊。
			NewFatReactor(
				FatTags(SkillGroup, Priority(1), Label("@Launch({1^/%} 400% Critical Damage)")),
				FatRespond(
					NewSignalTrigger(&LaunchSignal{}),
					NewSelectActor(
						NewVerbActor(NewAttack(nil, true), NewMultiplier(400, AxisEvaluator(Damage))),
						SideSelector(false),
						Healthy,
						NewSortSelector(HealthPercent, false),
						FrontSelector(1),
					),
				),
				FatCooling(NewSignalTrigger(&RoundEndSignal{}), 3),
			),

			// 對當前生命值百分比最低的 1 名敵人造成攻擊力 400% 的傷害；此攻擊必定爆擊。
			NewFatReactor(
				FatTags(SkillGroup, Priority(2), Label("@Launch({1!/%} 400% Critical Damage)")),
				FatRespond(
					NewSignalTrigger(&LaunchSignal{}),
					NewSelectActor(
						NewVerbActor(NewAttack(nil, true), NewMultiplier(400, AxisEvaluator(Damage))),
						SideSelector(false),
						Healthy,
						NewSortSelector(HealthPercent, true),
						FrontSelector(1),
					),
				),
				FatCooling(NewSignalTrigger(&RoundEndSignal{}), 3),
			),

			// 提升 25% 最大生命值(無法被解除)。
			NewFatReactor(
				FatTags(Priority(3), Label("@BattleStart({$} +25% HealthMaximum)")),
				FatRespond(
					NewSignalTrigger(&BattleStartSignal{}),
					NewSelectActor(
						NewVerbActor(
							NewBuff(false, nil, NewBuffReactor(
								HealthMaximum,
								false,
								ConstEvaluator(125),
								FatTags(Label("+25% HealthMaximum")))),
							nil,
						),
						CurrentSelector{},
					),
				),
			),

			// 戰鬥開始時，每有一名友軍，全體友軍提升 2% 攻擊力。(無法被解除)
			NewFatReactor(
				FatTags(Priority(4), Label("@BattleStart({~} +2% Attack*)")),
				FatRespond(
					NewSignalTrigger(&BattleStartSignal{}),
					NewSelectActor(
						NewVerbActor(
							NewBuff(false, nil, NewBuffReactor(
								Damage,
								false,
								NewAdder(100, NewMultiplier(2, NewSelectCounter(
									SideSelector(true),
									Healthy,
								))),
								FatTags(Label("+2% Attack*")))),
							nil,
						),
						SideSelector(true),
					),
				),
			),
		},

		// ////////////////////////////////////////////////////////////
		// [3] 徳川
		{
			// 對攻擊力最高的敵人進行 4 次攻擊，每次造成攻擊力 340% 的傷害。
			NewFatReactor(
				FatTags(SkillGroup, Priority(1), Label("@Launch({1^/D} 4 * 340% Damage)")),
				FatRespond(
					NewSignalTrigger(&LaunchSignal{}),
					NewSelectActor(
						NewRepeatActor(
							4,
							NewVerbActor(&Attack{}, NewMultiplier(340, AxisEvaluator(Damage))),
						),
						SideSelector(false),
						Healthy,
						NewSortSelector(Damage, false),
						FrontSelector(1),
					),
				),
				FatCooling(NewSignalTrigger(&RoundEndSignal{}), 4),
			),

			// 對生命值百分比最低的 3 名友軍治療，恢復徳川攻擊力 250% 的生命值
			NewFatReactor(
				FatTags(SkillGroup, Priority(2), Label("@Launch({3!/%} 250% Damage+)")),
				FatRespond(
					NewSignalTrigger(&LaunchSignal{}),
					NewSelectActor(
						NewVerbActor(NewHeal(nil), NewMultiplier(250, AxisEvaluator(Damage))),
						SideSelector(true),
						Healthy,
						NewSortSelector(HealthPercent, true),
						FrontSelector(3),
					),
				),
				FatCooling(NewSignalTrigger(&RoundEndSignal{}), 4),
			),

			// 強化自身的普通攻擊(無法被解除)，普通攻擊爆擊時，提昇 5% 攻擊力(最高 3 層，3 回合)，並刷新層數的回合。
			NewFatReactor(
				FatTags(Priority(3), Label("@PreAction({$} +5% Attack*)")),
				FatRespond(
					NewFatTrigger(
						&PreActionSignal{},
						CurrentIsSourceTrigger{},
						NewVerbTrigger[*Attack](),
						NewReactorTrigger(Regular[0]),
						CriticalStrikeTrigger{},
					),
					NewSelectActor(
						NewVerbActor(
							NewBuff(false, nil, NewBuffReactor(
								Damage,
								false,
								ConstEvaluator(105),
								FatTags(NewStackingLimit(3), Label("+5% Attack*")),
								FatCapacity(NewSignalTrigger(&RoundEndSignal{}), 3))),
							nil,
						),
						CurrentSelector{},
					),
				),
			),

			// 自身的生命值百分比為 50% 以下時，獲得「庇護」(最大生命值 30%，無法被解除)。
			NewFatReactor(
				FatTags(Priority(4), Label("@PostAction({$/< 50%} Sanctuary)")),
				FatRespond(
					NewSignalTrigger(&PostActionSignal{}),
					NewSelectActor(
						NewVerbActor(NewBuff(false, nil, Effect["Sanctuary"]), nil),
						CurrentSelector{},
						Healthy,
						NewWaterLevelSelector(Lt, AxisEvaluator(HealthPercent), 50),
						NewWaterLevelSelector(Lt, NewBuffCounter(Label("Sanctuary")), 1),
					),
				),
			),
		},

		// ////////////////////////////////////////////////////////////
		// [4] 武田
		{
			// 對全體敵人造成攻擊力 300% 的傷害。
			NewFatReactor(
				FatTags(SkillGroup, Priority(1), Label("@Launch({*} 300% Damage)")),
				FatRespond(
					NewSignalTrigger(&LaunchSignal{}),
					NewSelectActor(
						NewVerbActor(&Attack{}, NewMultiplier(300, AxisEvaluator(Damage))),
						SideSelector(false),
						Healthy,
					),
				),
				FatCooling(NewSignalTrigger(&RoundEndSignal{}), 5),
			),

			// 對隨機 1 名敵人造成攻擊力 560% 的傷害，並使目標減少 25% 攻擊力(2 回合)。
			NewFatReactor(
				FatTags(SkillGroup, Priority(2), Label("@Launch({1} 560% Damage, -25% Damage)")),
				FatRespond(
					NewSignalTrigger(&LaunchSignal{}),
					NewSelectActor(
						NewSequenceActor(
							NewVerbActor(&Attack{}, NewMultiplier(560, AxisEvaluator(Damage))),
							NewVerbActor(
								NewBuff(false, nil, NewBuffReactor(
									Damage,
									false,
									ConstEvaluator(75),
									FatCapacity(NewSignalTrigger(&RoundEndSignal{}), 2),
									FatTags(Label("-25% Damage")),
								)),
								nil,
							),
						),
						SideSelector(false),
						Healthy,
						Shuffle,
						FrontSelector(1),
					),
				),
				FatCooling(NewSignalTrigger(&RoundEndSignal{}), 4),
			),

			// 行動開始時，若自身的生命值百分比為 60% 以上，獲得「嘲諷」(2 回合) & 減少 15% 攻擊力(2 回合)。
			NewFatReactor(
				FatTags(Priority(3), Label("@Launch({$/>= 60%}, -15% Damage, Taunt)")),
				FatRespond(
					NewSignalTrigger(&LaunchSignal{}),
					NewSelectActor(
						NewSequenceActor(
							NewVerbActor(NewBuff(false, nil, Effect["Taunt"]), nil),
							NewVerbActor(
								NewBuff(false, nil, NewBuffReactor(
									Damage,
									false,
									ConstEvaluator(85),
									FatCapacity(NewSignalTrigger(&RoundEndSignal{}), 2),
									FatTags(Label("-15% Damage")))),
								nil,
							),
						),
						CurrentSelector{},
						NewWaterLevelSelector(Ge, AxisEvaluator(HealthPercent), 60),
					),
				),
				FatCooling(NewSignalTrigger(&RoundEndSignal{}), 4),
			),

			// 提升 25% 最大生命值(無法被解除)。
			// [2][2]
		},

		// ////////////////////////////////////////////////////////////
		// [5] 梅花
		{
			// 對隨機 1 名敵人進行 3 次攻擊，每次造成攻擊力 550% 傷害。每次攻擊都有 50% 機率對目標附加「暈眩」(1 回合)
			NewFatReactor(
				FatTags(SkillGroup, Priority(1), Label("@Launch({1} 3 * 550% Damage, P(50%) Stun)")),
				FatRespond(
					NewSignalTrigger(&LaunchSignal{}),
					NewSelectActor(
						NewRepeatActor(
							3,
							NewVerbActor(&Attack{}, NewMultiplier(550, AxisEvaluator(Damage))),
							NewProbabilityActor(RngX, ConstEvaluator(50), NewVerbActor(
								NewBuff(false, nil, Effect["Stun"]),
								nil,
							)),
						),
						SideSelector(false),
						Healthy,
						Shuffle,
						FrontSelector(1),
					),
				),
				FatCooling(NewSignalTrigger(&RoundEndSignal{}), 4),
			),

			// 對隨機 3 名敵人造成攻擊力 510% 的傷害。
			NewFatReactor(
				FatTags(SkillGroup, Priority(2), Label("@Launch({3} 510% Damage)")),
				FatRespond(
					NewSignalTrigger(&LaunchSignal{}),
					NewSelectActor(
						NewVerbActor(&Attack{}, NewMultiplier(510, AxisEvaluator(Damage))),
						SideSelector(false),
						Healthy,
						Shuffle,
						FrontSelector(3),
					),
				),
				FatCooling(NewSignalTrigger(&RoundEndSignal{}), 4),
			),

			// 自身的生命值百分比為 50% 以下時，獲得 20% 被擊減傷。(無法被解除)
			NewFatReactor(
				FatTags(Priority(3), Label("@PostAction({$/< 50%}, -20% Loss)")),
				FatRespond(
					NewSignalTrigger(&PostActionSignal{}),
					NewSelectActor(
						NewVerbActor(
							NewBuff(false, nil, NewBuffReactor(
								Loss,
								false,
								ConstEvaluator(80),
								FatTags(Label("-20% Loss"), "PBCj8umTGgCwqpTWV8KtqP"),
							)),
							nil,
						),
						CurrentSelector{},
						Healthy,
						NewWaterLevelSelector(Lt, AxisEvaluator(HealthPercent), 50),
						NewWaterLevelSelector(Lt, NewBuffCounter("PBCj8umTGgCwqpTWV8KtqP"), 1),
					),
				),
			),

			// 提升 30% 攻擊力(無法被解除)。
			NewFatReactor(
				FatTags(Priority(4), Label("@BattleStart({$}, 30% Damage)")),
				FatRespond(
					NewSignalTrigger(&BattleStartSignal{}),
					NewSelectActor(
						NewVerbActor(
							NewBuff(false, nil, NewBuffReactor(
								Damage,
								false,
								ConstEvaluator(130),
								FatTags(Label("30% Damage")))),
							nil,
						),
						CurrentSelector{},
					),
				),
			),
		},

		// ////////////////////////////////////////////////////////////
		// [6] 鑽石
		{
			// 使全體友軍增加攻擊力，增幅為鑽石攻擊力 15% (2 回合)
			NewFatReactor(
				FatTags(SkillGroup, Priority(1), Label("@Launch({~}, 15% Damage*)")),
				FatRespond(
					NewSignalTrigger(&LaunchSignal{}),
					NewSelectActor(
						NewVerbActor(
							NewBuff(false, nil, NewBuffReactor(
								Damage,
								true,
								nil,
								FatCapacity(NewSignalTrigger(&RoundEndSignal{}), 2),
								FatTags(Label("15% Damage*")))),
							NewMultiplier(15, AxisEvaluator(Damage)),
						),
						SideSelector(true),
						Healthy,
					),
				),
				FatCooling(NewSignalTrigger(&RoundEndSignal{}), 4),
			),

			// 對隨機 1 名敵人進行 4 次攻擊，每次造成攻擊力 350% 的傷害；使全體友軍獲得總傷害 40% 的「護盾」。
			NewFatReactor(
				FatTags(SkillGroup, Priority(2), Label("@Launch({1} 4 * 350% Damage; {*} +40% Shield)")),
				FatRespond(
					NewSignalTrigger(&LaunchSignal{}),
					NewSelectActor(
						NewRepeatActor(4, NewVerbActor(&Attack{}, NewMultiplier(350, AxisEvaluator(Damage)))),
						SideSelector(false),
						Healthy,
						Shuffle,
						FrontSelector(1),
					),
					NewSelectActor(
						NewVerbActor(
							NewBuff(true, NewMultiplier(40, LossEvaluator{}), Effect["Shield"]),
							nil,
						),
						SideSelector(true),
						Healthy,
					),
				),
				FatCooling(NewSignalTrigger(&RoundEndSignal{}), 4),
			),

			// 提升20%最大生命值(無法被解除)。
			NewFatReactor(
				FatTags(Priority(2), Label("@BattleStart({$}, 20% HealthMaximum)")),
				FatRespond(
					NewSignalTrigger(&BattleStartSignal{}),
					NewSelectActor(
						NewVerbActor(
							NewBuff(false, nil, NewBuffReactor(
								HealthMaximum,
								false,
								ConstEvaluator(120),
								FatTags(Label("20% HealthMaximum")))),
							nil,
						),
						CurrentSelector{},
					),
				),
			),

			// 受到爆擊時，使攻擊的敵人減少 15% 攻擊力(1 回合)。
			NewFatReactor(
				FatTags(Priority(3), Label("@PostAction({&/C}, -15% Damage)")),
				FatRespond(
					NewFatTrigger(
						&PostActionSignal{},
						CurrentIsTargetTrigger{},
						CriticalStrikeTrigger{},
					),
					NewSelectActor(
						NewVerbActor(
							NewBuff(false, nil, NewBuffReactor(
								Damage,
								false,
								ConstEvaluator(85),
								FatCapacity(NewSignalTrigger(&RoundEndSignal{}), 1),
								FatTags(Label("-15% Damage")))),
							nil,
						),
						SourceSelector{},
					),
				),
			),
		},

		// ////////////////////////////////////////////////////////////
		// [7] 王牌
		{
			// 對面前的 1 名敵人進行 2 次攻擊，每次造成攻擊力 360% 的傷害。每次攻擊都有 15% 機率對目標附加「暈眩」(1 回合)。
			NewFatReactor(
				FatTags(SkillGroup, Priority(1), Label("@Launch({1} 2 * 360% Damage, 15% Stun)")),
				FatRespond(
					NewSignalTrigger(&LaunchSignal{}),
					NewSelectActor(
						NewRepeatActor(2,
							NewVerbActor(&Attack{}, NewMultiplier(360, AxisEvaluator(Damage))),
							NewProbabilityActor(RngX, ConstEvaluator(15), NewVerbActor(
								NewBuff(false, nil, Effect["Stun"]),
								nil,
							)),
						),
						CounterPositionSelector{},
						Healthy,
					),
				),
				FatCooling(NewSignalTrigger(&RoundEndSignal{}), 4),
			),

			// 對面前的1名敵人造成攻擊力 455% 的傷害。並對目標附加「增益無效」(3 回合)。
			NewFatReactor(
				FatTags(SkillGroup, Priority(2), Label("@Launch({1} 455% Damage, BuffImmune")),
				FatRespond(
					NewSignalTrigger(&LaunchSignal{}),
					NewSelectActor(
						NewSequenceActor(
							NewVerbActor(&Attack{}, NewMultiplier(455, AxisEvaluator(Damage))),
							NewVerbActor(NewBuff(false, nil, Effect["BuffImmune"]), nil),
						),
						CounterPositionSelector{},
						Healthy,
					),
				),
				FatCooling(NewSignalTrigger(&RoundEndSignal{}), 3),
			),

			// 戰鬥開始時，每有一名友軍，全體友軍提升 2% 攻擊力。(無法被解除)
			// [2][4]

			// 提升 20% 攻擊力(無法被解除)。
			// [0][4]
		},

		// ////////////////////////////////////////////////////////////
		// [8] 紅心
		{
			// 使全體友軍獲得紅心攻擊力 120% 的「護盾」(1 回合)&「再生」(7%，1 回合)。
			NewFatReactor(
				FatTags(SkillGroup, Priority(1), Label("@Launch({*} +120% Shield, 7% Regeneration)")),
				FatRespond(
					NewSignalTrigger(&LaunchSignal{}),
					NewSelectActor(
						NewSequenceActor(
							NewVerbActor(
								NewBuff(true, NewMultiplier(120, AxisEvaluator(Damage)), Effect["Shield"]),
								nil,
							),
							NewVerbActor(NewBuff(false, nil, Effect["Regeneration"]), nil),
						),
						SideSelector(true),
						Healthy,
					),
				),
				FatCooling(NewSignalTrigger(&RoundEndSignal{}), 4),
			),

			// 對生命值百分比最低的 1 名敵人造成攻擊力 340% 的傷害。再對生命值百分比最低的 1 名友軍治療，恢復總傷害 80% 的生命值。
			NewFatReactor(
				FatTags(SkillGroup, Priority(2), Label("@Launch({1!/%} 340% Damage, {~1!/%} 80% Damage+)")),
				FatRespond(
					NewSignalTrigger(&LaunchSignal{}),
					NewSelectActor(
						NewVerbActor(&Attack{}, NewMultiplier(340, AxisEvaluator(Damage))),
						SideSelector(false),
						Healthy,
						NewSortSelector(HealthPercent, true),
						FrontSelector(1),
					),
					NewSelectActor(
						NewVerbActor(NewHeal(NewMultiplier(80, LossEvaluator{})), nil),
						SideSelector(true),
						Healthy,
						NewSortSelector(HealthPercent, true),
						FrontSelector(1),
					),
				),
				FatCooling(NewSignalTrigger(&RoundEndSignal{}), 4),
			),

			// 行動開始時獲得「控制效果免疫」(2 回合，無法被解除)。
			NewFatReactor(
				FatTags(Priority(3), Label("@Launch({$} ControlImmune)")),
				FatRespond(
					NewSignalTrigger(&LaunchSignal{}),
					NewSelectActor(
						NewVerbActor(NewBuff(false, nil, Effect["ControlImmune"]), nil),
						CurrentSelector{},
					),
				),
				FatCooling(NewSignalTrigger(&RoundEndSignal{}), 5),
			),

			// 提升 15% 防禦力(無法被解除)。
			NewFatReactor(
				FatTags(Priority(4), Label("@BattleStart({$} +15% Defense)")),
				FatRespond(
					NewSignalTrigger(&BattleStartSignal{}),
					NewSelectActor(
						NewVerbActor(
							NewBuff(false, nil, NewBuffReactor(
								Defense,
								false,
								ConstEvaluator(115),
								FatTags(Label("+15% Defense"))),
							),
							nil,
						),
						CurrentSelector{},
					),
				),
			),
		},

		// ////////////////////////////////////////////////////////////
		// [9] 黑桃
		{
			// 對隨機 1 名敵人進行 3 次攻擊，每次造成攻擊力 350% 的傷害。並使目標減少 50% 防禦力(1 回合)。
			NewFatReactor(
				FatTags(SkillGroup, Priority(1), Label("@Launch({1} 3 * 350% Damage, -50% Defense)")),
				FatRespond(
					NewSignalTrigger(&LaunchSignal{}),
					NewSelectActor(
						NewSequenceActor(
							NewRepeatActor(3, NewVerbActor(&Attack{}, NewMultiplier(350, AxisEvaluator(Damage)))),
							NewVerbActor(
								NewBuff(false, nil, NewBuffReactor(
									Defense,
									false,
									ConstEvaluator(50),
									FatCapacity(NewSignalTrigger(&RoundEndSignal{}), 1),
									FatTags(Label("-50% Defense")))),
								nil,
							),
						),
						SideSelector(false),
						Healthy,
						Shuffle,
						FrontSelector(1),
					),
				),
				FatCooling(NewSignalTrigger(&RoundEndSignal{}), 4),
			),

			// 對生命值百分比最低的敵人進行 3 次攻擊，每次造成攻擊力 380% 的傷害。
			NewFatReactor(
				FatTags(SkillGroup, Priority(2), Label("@Launch({1!/%} 3 * 380% Damage)")),
				FatRespond(
					NewSignalTrigger(&LaunchSignal{}),
					NewSelectActor(
						NewRepeatActor(3, NewVerbActor(&Attack{}, NewMultiplier(380, AxisEvaluator(Damage)))),
						SideSelector(false),
						Healthy,
						NewSortSelector(HealthPercent, true),
						FrontSelector(1),
					),
				),
				FatCooling(NewSignalTrigger(&RoundEndSignal{}), 4),
			),

			// 第 4 次行動開始時，解除自己身上的隨機 2 種減益效果。
			NewFatReactor(
				FatTags(Priority(3), Label("@Launch[..4]({$} Remove 2 Nerf)")),
				FatRespond(
					NewSignalTrigger(&LaunchSignal{}),
					NewSelectActor(
						NewVerbActor(NewPurge(RngX, "Nerf", 2), nil),
						CurrentSelector{},
					),
				),
				FatLeading(NewSignalTrigger(&LaunchSignal{}), 4),
				FatCapacity(nil, 1),
			),

			// 提升20%攻擊力(無法被解除)。
			// [0][4]
		},
	}

	Scaffold = []*FatReactor{
		NewFatReactor(
			FatTags(SkillGroup, Priority(1), Label("#1")),
			FatRespond(
				NewSignalTrigger(&LaunchSignal{}),
				NewSelectActor(
					NewVerbActor(
						NewBuff(false, nil, NewFatReactor(
							FatTags(Label("Nerf #1"), "Nerf", "Control"),
						)),
						nil,
					),
					SideSelector(false),
				),
			),
			FatCapacity(nil, 1),
		),
		NewFatReactor(
			FatTags(SkillGroup, Priority(2), Label("#2")),
			FatRespond(
				NewSignalTrigger(&LaunchSignal{}),
				NewSelectActor(
					NewVerbActor(
						NewBuff(false, nil, NewFatReactor(
							FatTags(Label("Nerf #2"), "Nerf"),
						)),
						nil,
					),
					SideSelector(false),
				),
			),
			FatCapacity(nil, 1),
		),
		NewFatReactor(
			FatTags(Priority(3), Label("#3")),
			FatRespond(
				NewSignalTrigger(&LaunchSignal{}),
				NewSelectActor(
					NewVerbActor(
						NewBuff(false, nil, NewFatReactor(
							FatTags(Label("Buff #1"), "Buff"),
						)),
						nil,
					),
					CurrentSelector{},
				),
			),
			FatCapacity(nil, 1),
		),
	}
)

type Element int
