package examples

import b "github.com/farseeingnorthwest/playground/battlefield/v2"

var (
	SkillGroup    = b.ExclusionGroup(0)
	Shuffle       = b.NewShuffleSelector(b.Label("Taunt"))
	ElementTheory = b.NewTheoryActor(
		map[any]map[any]int{
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
	Regular = []*b.FatReactor{
		b.NewFatReactor(
			b.FatTags(SkillGroup, b.Label("NormalAttack")),
			b.FatRespond(
				b.NewSignalTrigger(&b.LaunchSignal{}),
				b.NewSelectActor(
					b.NewVerbActor(b.NewAttack(nil, false), b.AxisEvaluator(b.Damage)),
					b.SideSelector(false),
					b.Healthy,
					Shuffle,
					b.FrontSelector(1),
				),
			),
		),
		b.NewFatReactor(
			b.FatTags(b.Priority(1000)),
			b.FatRespond(
				b.NewFatTrigger(
					&b.PreActionSignal{},
					b.CurrentIsSourceTrigger{},
					b.NewVerbTrigger[*b.Attack](),
				),
				b.NewProbabilityActor(
					b.AxisEvaluator(b.CriticalOdds),
					b.CriticalActor{},
				),
			),
		),
		b.NewFatReactor(
			b.FatTags(b.Priority(999)),
			b.FatRespond(
				b.NewFatTrigger(
					&b.PreActionSignal{},
					b.CurrentIsSourceTrigger{},
					b.NewVerbTrigger[*b.Attack](),
					b.CriticalStrikeTrigger{},
				),
				b.NewActionBuffer(
					b.AxisEvaluator(b.CriticalLoss),
					b.NewBuffer(b.Loss, false, nil),
				),
			),
		),
		b.NewFatReactor(
			b.FatRespond(
				b.NewFatTrigger(
					&b.PreActionSignal{},
					b.NewVerbTrigger[*b.Attack](),
				),
				b.NewActionBuffer(nil, ElementTheory),
			),
		),
	}

	Effect = map[string]*b.FatReactor{
		// 沉睡
		"Sleep": b.NewFatReactor(
			b.FatTags(SkillGroup, b.Priority(10), b.Label("Sleep")),
			b.FatRespond(
				b.NewSignalTrigger(&b.LaunchSignal{}),
			),
			b.FatCapacity(
				b.NewAnyTrigger(
					b.NewSignalTrigger(&b.RoundEndSignal{}),
					b.NewFatTrigger(
						&b.PostActionSignal{},
						b.VerbTrigger[*b.Attack]{},
						b.CurrentIsTargetTrigger{}),
				),
				1,
			),
		),

		// 暈眩
		"Stun": b.NewFatReactor(
			b.FatTags(SkillGroup, b.Priority(10), b.Label("Stun")),
			b.FatRespond(
				b.NewSignalTrigger(&b.LaunchSignal{}),
				b.NewSequenceActor(),
			),
			b.FatCapacity(
				b.NewSignalTrigger(&b.RoundEndSignal{}),
				1,
			),
		),

		// 多重屏障
		"Barrier": b.NewFatReactor(
			b.FatTags(b.Label("Barrier")),
			b.FatRespond(
				b.NewSignalTrigger(&b.PreLossSignal{}),
				b.NewLossStopper(b.NewMultiplier(10, b.AxisEvaluator(b.HealthMaximum)), true),
			),
			b.FatCapacity(nil, 1),
		),

		// 嘲諷
		"Taunt": b.NewFatReactor(
			b.FatTags(b.Label("Taunt")),
			b.FatCapacity(b.NewSignalTrigger(&b.RoundEndSignal{}), 2),
		),

		// 庇護
		"Sanctuary": b.NewFatReactor(
			b.FatTags(b.Label("Sanctuary"), b.Interest{}),
			b.FatRespond(
				b.NewSignalTrigger(&b.PreLossSignal{}),
				b.NewLossStopper(b.NewMultiplier(30, b.AxisEvaluator(b.HealthMaximum)), false),
			),
		),

		// 護盾
		"Shield": b.NewFatReactor(
			b.FatTags(b.Label("Shield")),
			b.FatRespond(
				b.NewSignalTrigger(&b.PreLossSignal{}),
				b.NewSelectActor(
					b.LossResister{},
					b.CurrentSelector{},
				),
			),
		),

		// 增益無效
		"BuffImmune": b.NewFatReactor(
			b.FatTags(b.Label("BuffImmune")),
			b.FatCapacity(b.NewSignalTrigger(&b.RoundEndSignal{}), 3),
			b.FatRespond(
				b.NewFatTrigger(
					&b.PreActionSignal{},
					b.CurrentIsTargetTrigger{},
					b.NewTagTrigger("Buff"),
				),
				b.NewSelectActor(b.ImmuneActor{}, b.CurrentSelector{}),
			),
		),

		// 控制效果免疫
		"ControlImmune": b.NewFatReactor(
			b.FatTags(b.Label("ControlImmune")),
			b.FatCapacity(b.NewSignalTrigger(&b.RoundEndSignal{}), 2),
			b.FatRespond(
				b.NewFatTrigger(
					&b.PreActionSignal{},
					b.CurrentIsTargetTrigger{},
					b.NewTagTrigger("Control"),
				),
				b.NewSelectActor(b.ImmuneActor{}, b.CurrentSelector{}),
			),
		),

		// 再生
		"Regeneration": b.NewFatReactor(
			b.FatTags(b.Label("Regeneration")),
			b.FatRespond(
				b.NewSignalTrigger(&b.RoundStartSignal{}),
				b.NewSelectActor(
					b.NewVerbActor(b.NewHeal(nil), b.NewMultiplier(7, b.AxisEvaluator(b.HealthMaximum))),
					b.CurrentSelector{},
				),
			),
			b.FatCapacity(nil, 1),
		),
	}

	Special = [][]*b.FatReactor{
		// ////////////////////////////////////////////////////////////
		// [0] 織田
		{
			// 對隨機 1 名敵人進行 3 次攻擊，每次造成攻擊力 460% 的傷害。
			b.NewFatReactor(
				b.FatTags(SkillGroup, b.Priority(1), b.Label("@Launch({1} 3 * 460% Damage)")),
				b.FatRespond(
					b.NewSignalTrigger(&b.LaunchSignal{}),
					b.NewSelectActor(
						b.NewRepeatActor(
							3,
							b.NewVerbActor(b.NewAttack(nil, false), b.NewMultiplier(460, b.AxisEvaluator(b.Damage))),
						),
						b.SideSelector(false),
						b.Healthy,
						Shuffle,
						b.FrontSelector(1),
					),
				),
				b.FatCooling(b.NewSignalTrigger(&b.RoundEndSignal{}), 5),
			),

			// 對全體敵人造成攻擊力 480% 的傷害。再對當前生命值最低的1名敵人造成攻擊力 520% 的傷害。
			b.NewFatReactor(
				b.FatTags(SkillGroup, b.Priority(2), b.Label("@Launch({*} 480% Damage; {1} 520% Damage)")),
				b.FatRespond(
					b.NewSignalTrigger(&b.LaunchSignal{}),
					b.NewSelectActor(
						b.NewVerbActor(b.NewAttack(nil, false), b.NewMultiplier(480, b.AxisEvaluator(b.Damage))),
						b.SideSelector(false),
						b.Healthy,
					),
					b.NewSelectActor(
						b.NewVerbActor(b.NewAttack(nil, false), b.NewMultiplier(520, b.AxisEvaluator(b.Damage))),
						b.SideSelector(false),
						b.Healthy,
						b.NewSortSelector(b.HealthPercent, true),
						b.FrontSelector(1),
					),
				),
				b.FatCooling(b.NewSignalTrigger(&b.RoundEndSignal{}), 4),
			),

			// 每次行動開始時，提升 2% 爆擊率(最高 15 層，無法被解除)。
			b.NewFatReactor(
				b.FatTags(b.Priority(3), b.Label("@Launch([15] +2% CriticalOdds)")),
				b.FatRespond(
					b.NewSignalTrigger(&b.LaunchSignal{}),
					b.NewSelectActor(
						b.NewVerbActor(
							b.NewBuff(false, nil, b.NewBuffReactor(
								b.CriticalOdds,
								true,
								b.ConstEvaluator(2),
								b.FatTags(
									b.NewStackingLimit("Qn7Vh9kuDNToXVG1kGFtzE", 15),
									b.Label("[15] +2% CriticalOdds"),
								),
							)),
							nil,
						),
						b.CurrentSelector{},
					),
				),
			),

			// 提升 20% 攻擊力。(無法被解除)
			b.NewFatReactor(
				b.FatTags(b.Priority(4), b.Label("@BattleStart({$} +20% Damage)")),
				b.FatRespond(
					b.NewSignalTrigger(&b.BattleStartSignal{}),
					b.NewSelectActor(
						b.NewVerbActor(
							b.NewBuff(false, nil, b.NewBuffReactor(
								b.Damage,
								false,
								b.ConstEvaluator(120),
								b.FatTags(b.Label("+20% Damage")))),
							nil,
						),
						b.CurrentSelector{},
					),
				),
			),
		},

		// ////////////////////////////////////////////////////////////
		// [1] 豐臣
		{
			// 對隨機 2 名敵人造成攻擊力 420% 的傷害。並對目標附加 30% 被擊增傷(1 回合)&「沉睡」(1 回合)。
			b.NewFatReactor(
				b.FatTags(SkillGroup, b.Priority(1), b.Label("@Launch({2} 420% Damage, +30% Loss, Sleep)")),
				b.FatRespond(
					b.NewSignalTrigger(&b.LaunchSignal{}),
					b.NewSelectActor(
						b.NewSequenceActor(
							b.NewVerbActor(b.NewAttack(nil, false), b.NewMultiplier(420, b.AxisEvaluator(b.Damage))),
							b.NewVerbActor(
								// 30% 被擊增傷 (1回合)
								b.NewBuff(false, nil, b.NewBuffReactor(
									b.Loss,
									false,
									b.ConstEvaluator(130),
									b.FatTags(b.Label("+30% Loss")),
									b.FatCapacity(b.NewSignalTrigger(&b.RoundEndSignal{}), 1))),
								nil,
							),
							b.NewVerbActor(b.NewBuff(false, nil, Effect["Sleep"]), nil),
						),
						b.SideSelector(false),
						b.Healthy,
						Shuffle,
						b.FrontSelector(2),
					),
				),
				b.FatCooling(b.NewSignalTrigger(&b.RoundEndSignal{}), 5),
			),

			// 對攻擊力最高的 1 名敵人造成攻擊力 520% 的傷害。並有 70% 機率對目標附加「暈眩」(1 回合)。
			b.NewFatReactor(
				b.FatTags(SkillGroup, b.Priority(2), b.Label("@Launch({1} 520% Damage, P(70%, Stun))")),
				b.FatRespond(
					b.NewSignalTrigger(&b.LaunchSignal{}),
					b.NewSelectActor(
						b.NewSequenceActor(
							b.NewVerbActor(b.NewAttack(nil, false), b.NewMultiplier(520, b.AxisEvaluator(b.Damage))),
							b.NewProbabilityActor(b.ConstEvaluator(70), b.NewVerbActor(
								b.NewBuff(false, nil, Effect["Stun"]),
								nil,
							)),
						),
						b.SideSelector(false),
						b.Healthy,
						b.NewSortSelector(b.Damage, false),
						b.FrontSelector(1),
					),
				),
				b.FatCooling(b.NewSignalTrigger(&b.RoundEndSignal{}), 4),
			),

			// 自身附帶 2 種以上減益效果時觸發，解除所有減益效果。
			b.NewFatReactor(
				b.FatTags(b.Priority(2), b.Label("@PostAction({$/<Nerf> >= 2} Purge())")),
				b.FatRespond(
					b.NewAnyTrigger(
						b.NewSignalTrigger(&b.RoundStartSignal{}),
						b.NewSignalTrigger(&b.PostActionSignal{}),
					),
					b.NewSelectActor(
						b.NewVerbActor(b.NewPurge("Nerf", 0), nil),
						b.CurrentSelector{},
						b.NewWaterLevelSelector(b.Ge, b.NewBuffCounter("Nerf"), 2),
					),
				),
				b.FatCooling(b.NewSignalTrigger(&b.RoundEndSignal{}), 4),
			),

			// 行動開始時獲得「多重屏障」(2 層，無法被解除)。
			b.NewFatReactor(
				b.FatTags(b.Priority(3), b.Label("@Launch({$} Barrier)")),
				b.FatRespond(
					b.NewSignalTrigger(&b.LaunchSignal{}),
					b.NewSelectActor(
						b.NewVerbActor(b.NewBuff(false, nil, Effect["Barrier"]), nil),
						b.CurrentSelector{},
					),
				),
				b.FatCooling(b.NewSignalTrigger(&b.RoundEndSignal{}), 4),
			),
		},

		// ////////////////////////////////////////////////////////////
		// [2] 上杉
		{
			// 對當前生命值百分比最高的 1 名敵人造成攻擊力 400% 的傷害；此攻擊必定爆擊。
			b.NewFatReactor(
				b.FatTags(SkillGroup, b.Priority(1), b.Label("@Launch({1^/%} 400% Critical Damage)")),
				b.FatRespond(
					b.NewSignalTrigger(&b.LaunchSignal{}),
					b.NewSelectActor(
						b.NewVerbActor(b.NewAttack(nil, true), b.NewMultiplier(400, b.AxisEvaluator(b.Damage))),
						b.SideSelector(false),
						b.Healthy,
						b.NewSortSelector(b.HealthPercent, false),
						b.FrontSelector(1),
					),
				),
				b.FatCooling(b.NewSignalTrigger(&b.RoundEndSignal{}), 3),
			),

			// 對當前生命值百分比最低的 1 名敵人造成攻擊力 400% 的傷害；此攻擊必定爆擊。
			b.NewFatReactor(
				b.FatTags(SkillGroup, b.Priority(2), b.Label("@Launch({1!/%} 400% Critical Damage)")),
				b.FatRespond(
					b.NewSignalTrigger(&b.LaunchSignal{}),
					b.NewSelectActor(
						b.NewVerbActor(b.NewAttack(nil, true), b.NewMultiplier(400, b.AxisEvaluator(b.Damage))),
						b.SideSelector(false),
						b.Healthy,
						b.NewSortSelector(b.HealthPercent, true),
						b.FrontSelector(1),
					),
				),
				b.FatCooling(b.NewSignalTrigger(&b.RoundEndSignal{}), 3),
			),

			// 提升 25% 最大生命值(無法被解除)。
			b.NewFatReactor(
				b.FatTags(b.Priority(3), b.Label("@BattleStart({$} +25% HealthMaximum)")),
				b.FatRespond(
					b.NewSignalTrigger(&b.BattleStartSignal{}),
					b.NewSelectActor(
						b.NewVerbActor(
							b.NewBuff(false, nil, b.NewBuffReactor(
								b.HealthMaximum,
								false,
								b.ConstEvaluator(125),
								b.FatTags(b.Label("+25% HealthMaximum")))),
							nil,
						),
						b.CurrentSelector{},
					),
				),
			),

			// 戰鬥開始時，每有一名友軍，全體友軍提升 2% 攻擊力。(無法被解除)
			b.NewFatReactor(
				b.FatTags(b.Priority(4), b.Label("@BattleStart({~} +2% Attack*)")),
				b.FatRespond(
					b.NewSignalTrigger(&b.BattleStartSignal{}),
					b.NewSelectActor(
						b.NewVerbActor(
							b.NewBuff(false, nil, b.NewBuffReactor(
								b.Damage,
								false,
								b.NewAdder(100, b.NewMultiplier(2, b.NewSelectCounter(
									b.SideSelector(true),
									b.Healthy,
								))),
								b.FatTags(b.Label("+2% Attack*")))),
							nil,
						),
						b.SideSelector(true),
					),
				),
			),
		},

		// ////////////////////////////////////////////////////////////
		// [3] 徳川
		{
			// 對攻擊力最高的敵人進行 4 次攻擊，每次造成攻擊力 340% 的傷害。
			b.NewFatReactor(
				b.FatTags(SkillGroup, b.Priority(1), b.Label("@Launch({1^/D} 4 * 340% Damage)")),
				b.FatRespond(
					b.NewSignalTrigger(&b.LaunchSignal{}),
					b.NewSelectActor(
						b.NewRepeatActor(
							4,
							b.NewVerbActor(b.NewAttack(nil, false), b.NewMultiplier(340, b.AxisEvaluator(b.Damage))),
						),
						b.SideSelector(false),
						b.Healthy,
						b.NewSortSelector(b.Damage, false),
						b.FrontSelector(1),
					),
				),
				b.FatCooling(b.NewSignalTrigger(&b.RoundEndSignal{}), 4),
			),

			// 對生命值百分比最低的 3 名友軍治療，恢復徳川攻擊力 250% 的生命值
			b.NewFatReactor(
				b.FatTags(SkillGroup, b.Priority(2), b.Label("@Launch({3!/%} 250% Damage+)")),
				b.FatRespond(
					b.NewSignalTrigger(&b.LaunchSignal{}),
					b.NewSelectActor(
						b.NewVerbActor(b.NewHeal(nil), b.NewMultiplier(250, b.AxisEvaluator(b.Damage))),
						b.SideSelector(true),
						b.Healthy,
						b.NewSortSelector(b.HealthPercent, true),
						b.FrontSelector(3),
					),
				),
				b.FatCooling(b.NewSignalTrigger(&b.RoundEndSignal{}), 4),
			),

			// 強化自身的普通攻擊(無法被解除)，普通攻擊爆擊時，提昇 5% 攻擊力(最高 3 層，3 回合)，並刷新層數的回合。
			b.NewFatReactor(
				b.FatTags(b.Priority(3), b.Label("@PreAction({$} +5% Attack*)")),
				b.FatRespond(
					b.NewFatTrigger(
						&b.PreActionSignal{},
						b.CurrentIsSourceTrigger{},
						b.NewVerbTrigger[*b.Attack](),
						b.NewReactorTrigger(b.Label("NormalAttack")),
						b.CriticalStrikeTrigger{},
					),
					b.NewSelectActor(
						b.NewVerbActor(
							b.NewBuff(false, nil, b.NewBuffReactor(
								b.Damage,
								false,
								b.ConstEvaluator(105),
								b.FatTags(
									b.NewStackingLimit("VUE3GbrQweSqdpa1DfoqYt", 3),
									b.Label("+5% Attack*"),
								),
								b.FatCapacity(b.NewSignalTrigger(&b.RoundEndSignal{}), 3))),
							nil,
						),
						b.CurrentSelector{},
					),
				),
			),

			// 自身的生命值百分比為 50% 以下時，獲得「庇護」(最大生命值 30%，無法被解除)。
			b.NewFatReactor(
				b.FatTags(b.Priority(4), b.Label("@PostAction({$/< 50%} Sanctuary)")),
				b.FatRespond(
					b.NewSignalTrigger(&b.PostActionSignal{}),
					b.NewSelectActor(
						b.NewVerbActor(b.NewBuff(false, nil, Effect["Sanctuary"]), nil),
						b.CurrentSelector{},
						b.Healthy,
						b.NewWaterLevelSelector(b.Lt, b.AxisEvaluator(b.HealthPercent), 50),
						b.NewWaterLevelSelector(b.Lt, b.NewBuffCounter(b.Label("Sanctuary")), 1),
					),
				),
			),
		},

		// ////////////////////////////////////////////////////////////
		// [4] 武田
		{
			// 對全體敵人造成攻擊力 300% 的傷害。
			b.NewFatReactor(
				b.FatTags(SkillGroup, b.Priority(1), b.Label("@Launch({*} 300% Damage)")),
				b.FatRespond(
					b.NewSignalTrigger(&b.LaunchSignal{}),
					b.NewSelectActor(
						b.NewVerbActor(b.NewAttack(nil, false), b.NewMultiplier(300, b.AxisEvaluator(b.Damage))),
						b.SideSelector(false),
						b.Healthy,
					),
				),
				b.FatCooling(b.NewSignalTrigger(&b.RoundEndSignal{}), 5),
			),

			// 對隨機 1 名敵人造成攻擊力 560% 的傷害，並使目標減少 25% 攻擊力(2 回合)。
			b.NewFatReactor(
				b.FatTags(SkillGroup, b.Priority(2), b.Label("@Launch({1} 560% Damage, -25% Damage)")),
				b.FatRespond(
					b.NewSignalTrigger(&b.LaunchSignal{}),
					b.NewSelectActor(
						b.NewSequenceActor(
							b.NewVerbActor(b.NewAttack(nil, false), b.NewMultiplier(560, b.AxisEvaluator(b.Damage))),
							b.NewVerbActor(
								b.NewBuff(false, nil, b.NewBuffReactor(
									b.Damage,
									false,
									b.ConstEvaluator(75),
									b.FatCapacity(b.NewSignalTrigger(&b.RoundEndSignal{}), 2),
									b.FatTags(b.Label("-25% Damage")),
								)),
								nil,
							),
						),
						b.SideSelector(false),
						b.Healthy,
						Shuffle,
						b.FrontSelector(1),
					),
				),
				b.FatCooling(b.NewSignalTrigger(&b.RoundEndSignal{}), 4),
			),

			// 行動開始時，若自身的生命值百分比為 60% 以上，獲得「嘲諷」(2 回合) & 減少 15% 攻擊力(2 回合)。
			b.NewFatReactor(
				b.FatTags(b.Priority(3), b.Label("@Launch({$/>= 60%}, -15% Damage, Taunt)")),
				b.FatRespond(
					b.NewSignalTrigger(&b.LaunchSignal{}),
					b.NewSelectActor(
						b.NewSequenceActor(
							b.NewVerbActor(b.NewBuff(false, nil, Effect["Taunt"]), nil),
							b.NewVerbActor(
								b.NewBuff(false, nil, b.NewBuffReactor(
									b.Damage,
									false,
									b.ConstEvaluator(85),
									b.FatCapacity(b.NewSignalTrigger(&b.RoundEndSignal{}), 2),
									b.FatTags(b.Label("-15% Damage")))),
								nil,
							),
						),
						b.CurrentSelector{},
						b.NewWaterLevelSelector(b.Ge, b.AxisEvaluator(b.HealthPercent), 60),
					),
				),
				b.FatCooling(b.NewSignalTrigger(&b.RoundEndSignal{}), 4),
			),

			// 提升 25% 最大生命值(無法被解除)。
			// [2][2]
		},

		// ////////////////////////////////////////////////////////////
		// [5] 梅花
		{
			// 對隨機 1 名敵人進行 3 次攻擊，每次造成攻擊力 550% 傷害。每次攻擊都有 50% 機率對目標附加「暈眩」(1 回合)
			b.NewFatReactor(
				b.FatTags(SkillGroup, b.Priority(1), b.Label("@Launch({1} 3 * 550% Damage, P(50%) Stun)")),
				b.FatRespond(
					b.NewSignalTrigger(&b.LaunchSignal{}),
					b.NewSelectActor(
						b.NewRepeatActor(
							3,
							b.NewVerbActor(b.NewAttack(nil, false), b.NewMultiplier(550, b.AxisEvaluator(b.Damage))),
							b.NewProbabilityActor(b.ConstEvaluator(50), b.NewVerbActor(
								b.NewBuff(false, nil, Effect["Stun"]),
								nil,
							)),
						),
						b.SideSelector(false),
						b.Healthy,
						Shuffle,
						b.FrontSelector(1),
					),
				),
				b.FatCooling(b.NewSignalTrigger(&b.RoundEndSignal{}), 4),
			),

			// 對隨機 3 名敵人造成攻擊力 510% 的傷害。
			b.NewFatReactor(
				b.FatTags(SkillGroup, b.Priority(2), b.Label("@Launch({3} 510% Damage)")),
				b.FatRespond(
					b.NewSignalTrigger(&b.LaunchSignal{}),
					b.NewSelectActor(
						b.NewVerbActor(b.NewAttack(nil, false), b.NewMultiplier(510, b.AxisEvaluator(b.Damage))),
						b.SideSelector(false),
						b.Healthy,
						Shuffle,
						b.FrontSelector(3),
					),
				),
				b.FatCooling(b.NewSignalTrigger(&b.RoundEndSignal{}), 4),
			),

			// 自身的生命值百分比為 50% 以下時，獲得 20% 被擊減傷。(無法被解除)
			b.NewFatReactor(
				b.FatTags(b.Priority(3), b.Label("@PostAction({$/< 50%}, -20% Loss)")),
				b.FatRespond(
					b.NewSignalTrigger(&b.PostActionSignal{}),
					b.NewSelectActor(
						b.NewVerbActor(
							b.NewBuff(false, nil, b.NewBuffReactor(
								b.Loss,
								false,
								b.ConstEvaluator(80),
								b.FatTags(b.Label("-20% Loss"), "PBCj8umTGgCwqpTWV8KtqP"),
							)),
							nil,
						),
						b.CurrentSelector{},
						b.Healthy,
						b.NewWaterLevelSelector(b.Lt, b.AxisEvaluator(b.HealthPercent), 50),
						b.NewWaterLevelSelector(b.Lt, b.NewBuffCounter("PBCj8umTGgCwqpTWV8KtqP"), 1),
					),
				),
			),

			// 提升 30% 攻擊力(無法被解除)。
			b.NewFatReactor(
				b.FatTags(b.Priority(4), b.Label("@BattleStart({$}, 30% Damage)")),
				b.FatRespond(
					b.NewSignalTrigger(&b.BattleStartSignal{}),
					b.NewSelectActor(
						b.NewVerbActor(
							b.NewBuff(false, nil, b.NewBuffReactor(
								b.Damage,
								false,
								b.ConstEvaluator(130),
								b.FatTags(b.Label("30% Damage")))),
							nil,
						),
						b.CurrentSelector{},
					),
				),
			),
		},

		// ////////////////////////////////////////////////////////////
		// [6] 鑽石
		{
			// 使全體友軍增加攻擊力，增幅為鑽石攻擊力 15% (2 回合)
			b.NewFatReactor(
				b.FatTags(SkillGroup, b.Priority(1), b.Label("@Launch({~}, 15% Damage*)")),
				b.FatRespond(
					b.NewSignalTrigger(&b.LaunchSignal{}),
					b.NewSelectActor(
						b.NewVerbActor(
							b.NewBuff(false, nil, b.NewBuffReactor(
								b.Damage,
								true,
								nil,
								b.FatCapacity(b.NewSignalTrigger(&b.RoundEndSignal{}), 2),
								b.FatTags(b.Label("15% Damage*")))),
							b.NewMultiplier(15, b.AxisEvaluator(b.Damage)),
						),
						b.SideSelector(true),
						b.Healthy,
					),
				),
				b.FatCooling(b.NewSignalTrigger(&b.RoundEndSignal{}), 4),
			),

			// 對隨機 1 名敵人進行 4 次攻擊，每次造成攻擊力 350% 的傷害；使全體友軍獲得總傷害 40% 的「護盾」。
			b.NewFatReactor(
				b.FatTags(SkillGroup, b.Priority(2), b.Label("@Launch({1} 4 * 350% Damage; {*} +40% Shield)")),
				b.FatRespond(
					b.NewSignalTrigger(&b.LaunchSignal{}),
					b.NewSelectActor(
						b.NewRepeatActor(4, b.NewVerbActor(b.NewAttack(nil, false), b.NewMultiplier(350, b.AxisEvaluator(b.Damage)))),
						b.SideSelector(false),
						b.Healthy,
						Shuffle,
						b.FrontSelector(1),
					),
					b.NewSelectActor(
						b.NewVerbActor(
							b.NewBuff(true, b.NewMultiplier(40, b.LossEvaluator{}), Effect["Shield"]),
							nil,
						),
						b.SideSelector(true),
						b.Healthy,
					),
				),
				b.FatCooling(b.NewSignalTrigger(&b.RoundEndSignal{}), 4),
			),

			// 提升20%最大生命值(無法被解除)。
			b.NewFatReactor(
				b.FatTags(b.Priority(2), b.Label("@BattleStart({$}, 20% HealthMaximum)")),
				b.FatRespond(
					b.NewSignalTrigger(&b.BattleStartSignal{}),
					b.NewSelectActor(
						b.NewVerbActor(
							b.NewBuff(false, nil, b.NewBuffReactor(
								b.HealthMaximum,
								false,
								b.ConstEvaluator(120),
								b.FatTags(b.Label("20% HealthMaximum")))),
							nil,
						),
						b.CurrentSelector{},
					),
				),
			),

			// 受到爆擊時，使攻擊的敵人減少 15% 攻擊力(1 回合)。
			b.NewFatReactor(
				b.FatTags(b.Priority(3), b.Label("@PostAction({&/C}, -15% Damage)")),
				b.FatRespond(
					b.NewFatTrigger(
						&b.PostActionSignal{},
						b.CurrentIsTargetTrigger{},
						b.CriticalStrikeTrigger{},
					),
					b.NewSelectActor(
						b.NewVerbActor(
							b.NewBuff(false, nil, b.NewBuffReactor(
								b.Damage,
								false,
								b.ConstEvaluator(85),
								b.FatCapacity(b.NewSignalTrigger(&b.RoundEndSignal{}), 1),
								b.FatTags(b.Label("-15% Damage")))),
							nil,
						),
						b.SourceSelector{},
					),
				),
			),
		},

		// ////////////////////////////////////////////////////////////
		// [7] 王牌
		{
			// 對面前的 1 名敵人進行 2 次攻擊，每次造成攻擊力 360% 的傷害。每次攻擊都有 15% 機率對目標附加「暈眩」(1 回合)。
			b.NewFatReactor(
				b.FatTags(SkillGroup, b.Priority(1), b.Label("@Launch({1} 2 * 360% Damage, 15% Stun)")),
				b.FatRespond(
					b.NewSignalTrigger(&b.LaunchSignal{}),
					b.NewSelectActor(
						b.NewRepeatActor(2,
							b.NewVerbActor(b.NewAttack(nil, false), b.NewMultiplier(360, b.AxisEvaluator(b.Damage))),
							b.NewProbabilityActor(b.ConstEvaluator(15), b.NewVerbActor(
								b.NewBuff(false, nil, Effect["Stun"]),
								nil,
							)),
						),
						b.NewCounterPositionSelector(0),
						b.Healthy,
					),
				),
				b.FatCooling(b.NewSignalTrigger(&b.RoundEndSignal{}), 4),
			),

			// 對面前的1名敵人造成攻擊力 455% 的傷害。並對目標附加「增益無效」(3 回合)。
			b.NewFatReactor(
				b.FatTags(SkillGroup, b.Priority(2), b.Label("@Launch({1} 455% Damage, BuffImmune")),
				b.FatRespond(
					b.NewSignalTrigger(&b.LaunchSignal{}),
					b.NewSelectActor(
						b.NewSequenceActor(
							b.NewVerbActor(b.NewAttack(nil, false), b.NewMultiplier(455, b.AxisEvaluator(b.Damage))),
							b.NewVerbActor(b.NewBuff(false, nil, Effect["BuffImmune"]), nil),
						),
						b.NewCounterPositionSelector(0),
						b.Healthy,
					),
				),
				b.FatCooling(b.NewSignalTrigger(&b.RoundEndSignal{}), 3),
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
			b.NewFatReactor(
				b.FatTags(SkillGroup, b.Priority(1), b.Label("@Launch({*} +120% Shield, 7% Regeneration)")),
				b.FatRespond(
					b.NewSignalTrigger(&b.LaunchSignal{}),
					b.NewSelectActor(
						b.NewSequenceActor(
							b.NewVerbActor(
								b.NewBuff(true, b.NewMultiplier(120, b.AxisEvaluator(b.Damage)), Effect["Shield"]),
								nil,
							),
							b.NewVerbActor(b.NewBuff(false, nil, Effect["Regeneration"]), nil),
						),
						b.SideSelector(true),
						b.Healthy,
					),
				),
				b.FatCooling(b.NewSignalTrigger(&b.RoundEndSignal{}), 4),
			),

			// 對生命值百分比最低的 1 名敵人造成攻擊力 340% 的傷害。再對生命值百分比最低的 1 名友軍治療，恢復總傷害 80% 的生命值。
			b.NewFatReactor(
				b.FatTags(SkillGroup, b.Priority(2), b.Label("@Launch({1!/%} 340% Damage, {~1!/%} 80% Damage+)")),
				b.FatRespond(
					b.NewSignalTrigger(&b.LaunchSignal{}),
					b.NewSelectActor(
						b.NewVerbActor(b.NewAttack(nil, false), b.NewMultiplier(340, b.AxisEvaluator(b.Damage))),
						b.SideSelector(false),
						b.Healthy,
						b.NewSortSelector(b.HealthPercent, true),
						b.FrontSelector(1),
					),
					b.NewSelectActor(
						b.NewVerbActor(b.NewHeal(b.NewMultiplier(80, b.LossEvaluator{})), nil),
						b.SideSelector(true),
						b.Healthy,
						b.NewSortSelector(b.HealthPercent, true),
						b.FrontSelector(1),
					),
				),
				b.FatCooling(b.NewSignalTrigger(&b.RoundEndSignal{}), 4),
			),

			// 行動開始時獲得「控制效果免疫」(2 回合，無法被解除)。
			b.NewFatReactor(
				b.FatTags(b.Priority(3), b.Label("@Launch({$} ControlImmune)")),
				b.FatRespond(
					b.NewSignalTrigger(&b.LaunchSignal{}),
					b.NewSelectActor(
						b.NewVerbActor(b.NewBuff(false, nil, Effect["ControlImmune"]), nil),
						b.CurrentSelector{},
					),
				),
				b.FatCooling(b.NewSignalTrigger(&b.RoundEndSignal{}), 5),
			),

			// 提升 15% 防禦力(無法被解除)。
			b.NewFatReactor(
				b.FatTags(b.Priority(4), b.Label("@BattleStart({$} +15% Defense)")),
				b.FatRespond(
					b.NewSignalTrigger(&b.BattleStartSignal{}),
					b.NewSelectActor(
						b.NewVerbActor(
							b.NewBuff(false, nil, b.NewBuffReactor(
								b.Defense,
								false,
								b.ConstEvaluator(115),
								b.FatTags(b.Label("+15% Defense"))),
							),
							nil,
						),
						b.CurrentSelector{},
					),
				),
			),
		},

		// ////////////////////////////////////////////////////////////
		// [9] 黑桃
		{
			// 對隨機 1 名敵人進行 3 次攻擊，每次造成攻擊力 350% 的傷害。並使目標減少 50% 防禦力(1 回合)。
			b.NewFatReactor(
				b.FatTags(SkillGroup, b.Priority(1), b.Label("@Launch({1} 3 * 350% Damage, -50% Defense)")),
				b.FatRespond(
					b.NewSignalTrigger(&b.LaunchSignal{}),
					b.NewSelectActor(
						b.NewSequenceActor(
							b.NewRepeatActor(3, b.NewVerbActor(b.NewAttack(nil, false), b.NewMultiplier(350, b.AxisEvaluator(b.Damage)))),
							b.NewVerbActor(
								b.NewBuff(false, nil, b.NewBuffReactor(
									b.Defense,
									false,
									b.ConstEvaluator(50),
									b.FatCapacity(b.NewSignalTrigger(&b.RoundEndSignal{}), 1),
									b.FatTags(b.Label("-50% Defense")))),
								nil,
							),
						),
						b.SideSelector(false),
						b.Healthy,
						Shuffle,
						b.FrontSelector(1),
					),
				),
				b.FatCooling(b.NewSignalTrigger(&b.RoundEndSignal{}), 4),
			),

			// 對生命值百分比最低的敵人進行 3 次攻擊，每次造成攻擊力 380% 的傷害。
			b.NewFatReactor(
				b.FatTags(SkillGroup, b.Priority(2), b.Label("@Launch({1!/%} 3 * 380% Damage)")),
				b.FatRespond(
					b.NewSignalTrigger(&b.LaunchSignal{}),
					b.NewSelectActor(
						b.NewRepeatActor(3, b.NewVerbActor(b.NewAttack(nil, false), b.NewMultiplier(380, b.AxisEvaluator(b.Damage)))),
						b.SideSelector(false),
						b.Healthy,
						b.NewSortSelector(b.HealthPercent, true),
						b.FrontSelector(1),
					),
				),
				b.FatCooling(b.NewSignalTrigger(&b.RoundEndSignal{}), 4),
			),

			// 第 4 次行動開始時，解除自己身上的隨機 2 種減益效果。
			b.NewFatReactor(
				b.FatTags(b.Priority(3), b.Label("@Launch[..4]({$} Remove 2 Nerf)")),
				b.FatRespond(
					b.NewSignalTrigger(&b.LaunchSignal{}),
					b.NewSelectActor(
						b.NewVerbActor(b.NewPurge("Nerf", 2), nil),
						b.CurrentSelector{},
					),
				),
				b.FatLeading(b.NewSignalTrigger(&b.LaunchSignal{}), 4),
				b.FatCapacity(nil, 1),
			),

			// 提升20%攻擊力(無法被解除)。
			// [0][4]
		},
	}

	Scaffold = []*b.FatReactor{
		b.NewFatReactor(
			b.FatTags(SkillGroup, b.Priority(1), b.Label("#1")),
			b.FatRespond(
				b.NewSignalTrigger(&b.LaunchSignal{}),
				b.NewSelectActor(
					b.NewVerbActor(
						b.NewBuff(false, nil, b.NewFatReactor(
							b.FatTags(b.Label("Nerf #1"), "Nerf", "Control"),
						)),
						nil,
					),
					b.SideSelector(false),
				),
			),
			b.FatCapacity(nil, 1),
		),
		b.NewFatReactor(
			b.FatTags(SkillGroup, b.Priority(2), b.Label("#2")),
			b.FatRespond(
				b.NewSignalTrigger(&b.LaunchSignal{}),
				b.NewSelectActor(
					b.NewVerbActor(
						b.NewBuff(false, nil, b.NewFatReactor(
							b.FatTags(b.Label("Nerf #2"), "Nerf"),
						)),
						nil,
					),
					b.SideSelector(false),
				),
			),
			b.FatCapacity(nil, 1),
		),
		b.NewFatReactor(
			b.FatTags(b.Priority(3), b.Label("#3")),
			b.FatRespond(
				b.NewSignalTrigger(&b.LaunchSignal{}),
				b.NewSelectActor(
					b.NewVerbActor(
						b.NewBuff(false, nil, b.NewFatReactor(
							b.FatTags(b.Label("Buff #1"), "Buff"),
						)),
						nil,
					),
					b.CurrentSelector{},
				),
			),
			b.FatCapacity(nil, 1),
		),
	}
)
