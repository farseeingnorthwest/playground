package battlefield

var (
	rng        = &RngProxy{}
	SkillGroup = ExclusionGroup(0)
	Regular    = []Reactor{
		NewFatReactor(
			FatTags(SkillGroup, Label("NormalAttack")),
			FatRespond(
				NewSignalTrigger(&LaunchSignal{}),
				NewSelectActor(
					NewVerbActor(&Attack{}, AxisEvaluator(Damage)),
					SideSelector(false),
					Healthy,
					NewShuffleSelector(rng, nil),
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
					rng,
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
	}

	Special = [][]Reactor{
		// ////////////////////////////////////////////////////////////
		// [0] 織田
		{
			// 對隨機1名敵人進行3次攻擊，每次造成攻擊力460%的傷害。
			NewFatReactor(
				FatTags(SkillGroup, Priority(1), Label("@Launch({1} 3 * 460% Damage)")),
				FatRespond(
					NewSignalTrigger(&LaunchSignal{}),
					NewSelectActor(
						NewSequenceActor(
							NewVerbActor(&Attack{}, NewMultiplier(AxisEvaluator(Damage), 460)),
							NewVerbActor(&Attack{}, NewMultiplier(AxisEvaluator(Damage), 460)),
							NewVerbActor(&Attack{}, NewMultiplier(AxisEvaluator(Damage), 460)),
						),
						SideSelector(false),
						Healthy,
						NewShuffleSelector(rng, nil),
						FrontSelector(1),
					),
				),
				FatCooling(NewSignalTrigger(&RoundEndSignal{}), 5),
			),

			// 對全體敵人造成攻擊力480%的傷害。再對當前生命值最低的1名敵人造成攻擊力520%的傷害。
			NewFatReactor(
				FatTags(SkillGroup, Priority(2), Label("@Launch({*} 480% Damage; {1} 520% Damage)")),
				FatRespond(
					NewSignalTrigger(&LaunchSignal{}),
					NewSelectActor(
						NewVerbActor(&Attack{}, NewMultiplier(AxisEvaluator(Damage), 480)),
						SideSelector(false),
						Healthy,
					),
					NewSelectActor(
						NewVerbActor(&Attack{}, NewMultiplier(AxisEvaluator(Damage), 520)),
						SideSelector(false),
						Healthy,
						NewSortSelector(HealthPercent, true),
						FrontSelector(1),
					),
				),
				FatCooling(NewSignalTrigger(&RoundEndSignal{}), 4),
			),

			// 每次行動開始時，提升2%爆擊率(最高15層，無法被解除)。
			NewFatReactor(
				FatTags(Priority(3), Label("@Launch([15] +2% CriticalOdds)")),
				FatRespond(
					NewSignalTrigger(&LaunchSignal{}),
					NewSelectActor(
						NewVerbActor(
							NewBuff(
								nil,
								NewBuffReactor(CriticalOdds, true, ConstEvaluator(2), FatTags(
									Label("[15] +2% CriticalOdds"))),
							),
							nil,
						),
						CurrentSelector{},
					),
				),
			),

			// 提升20%攻擊力。(無法被解除)
			NewFatReactor(
				FatTags(Priority(4), Label("@BattleStart({$} +20% Damage)")),
				FatRespond(
					NewSignalTrigger(&BattleStartSignal{}),
					NewSelectActor(
						NewVerbActor(
							NewBuff(nil, NewBuffReactor(
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
			// 對隨機2名敵人造成攻擊力420%的傷害。並對目標附加30%被擊增傷(1回合)&「沉睡」(1回合)。
			NewFatReactor(
				FatTags(SkillGroup, Priority(1), Label("@Launch({2} 420% Damage, +30% Loss, Sleeping)")),
				FatRespond(
					NewSignalTrigger(&LaunchSignal{}),
					NewSelectActor(
						NewSequenceActor(
							NewVerbActor(&Attack{}, NewMultiplier(AxisEvaluator(Damage), 420)),
							NewVerbActor(
								// 30% 被擊增傷 (1回合)
								NewBuff(nil, NewBuffReactor(
									Loss,
									false,
									ConstEvaluator(130),
									FatTags(Label("+30% Loss")),
									FatCapacity(NewSignalTrigger(&RoundEndSignal{}), 1))),
								nil,
							),
							// 「沉睡」 (1回合)
							NewVerbActor(
								NewBuff(nil, NewFatReactor(
									FatTags(SkillGroup, Priority(10), Label("Sleeping")),
									FatRespond(
										NewSignalTrigger(&LaunchSignal{}),
										NewSequenceActor(),
									),
									FatCapacity(
										NewOrTrigger(
											NewSignalTrigger(&RoundEndSignal{}),
											NewFatTrigger(
												&PostActionSignal{},
												VerbTrigger[*Attack]{},
												CurrentIsTargetTrigger{}),
										),
										1,
									),
								)),
								nil,
							),
						),
						SideSelector(false),
						Healthy,
						NewShuffleSelector(rng, nil),
						FrontSelector(2),
					),
				),
				FatCooling(NewSignalTrigger(&RoundEndSignal{}), 5),
			),

			// 對攻擊力最高的1名敵人造成攻擊力520%的傷害。並有70%機率對目標附加「暈眩」(1回合)。
			NewFatReactor(
				FatTags(SkillGroup, Priority(2), Label("@Launch({1} 520% Damage, P(70%, Dizzy))")),
				FatRespond(
					NewSignalTrigger(&LaunchSignal{}),
					NewSelectActor(
						NewSequenceActor(
							NewVerbActor(&Attack{}, NewMultiplier(AxisEvaluator(Damage), 520)),
							// 「暈眩」 (1回合)
							NewProbabilityActor(rng, ConstEvaluator(70), NewVerbActor(
								NewBuff(nil, NewFatReactor(
									FatTags(SkillGroup, Priority(10), Label("Dizzy")),
									FatRespond(
										NewSignalTrigger(&LaunchSignal{}),
										NewSequenceActor(),
									),
									FatCapacity(
										NewSignalTrigger(&RoundEndSignal{}),
										1,
									),
								)),
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

			// 自身附帶2種以上減益效果時觸發，解除所有減益效果。

			// 行動開始時獲得「多重屏障」(2層，無法被解除)。
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
						NewVerbActor(NewAttack(nil, true), NewMultiplier(AxisEvaluator(Damage), 400)),
						SideSelector(false),
						Healthy,
						NewSortSelector(HealthPercent, false),
						FrontSelector(1),
					),
				),
				FatCooling(NewSignalTrigger(&RoundEndSignal{}), 3),
			),

			// 對當前生命值百分比最低的 1 名敵人造成攻擊力400%的傷害；此攻擊必定爆擊。
			NewFatReactor(
				FatTags(SkillGroup, Priority(2), Label("@Launch({1!/%} 400% Critical Damage)")),
				FatRespond(
					NewSignalTrigger(&LaunchSignal{}),
					NewSelectActor(
						NewVerbActor(NewAttack(nil, true), NewMultiplier(AxisEvaluator(Damage), 400)),
						SideSelector(false),
						Healthy,
						NewSortSelector(HealthPercent, true),
						FrontSelector(1),
					),
				),
				FatCooling(NewSignalTrigger(&RoundEndSignal{}), 3),
			),

			// 提升25%最大生命值(無法被解除)。
			NewFatReactor(
				FatTags(SkillGroup, Priority(3), Label("@BattleStart({$} +25% HealthMaximum)")),
				FatRespond(
					NewSignalTrigger(&BattleStartSignal{}),
					NewSelectActor(
						NewVerbActor(
							NewBuff(nil, NewBuffReactor(
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

			// 戰鬥開始時，每有一名友軍，全體友軍提升2%攻擊力。(無法被解除)
			NewFatReactor(
				FatTags(SkillGroup, Priority(4), Label("@BattleStart({~} +2% Attack*)")),
				FatRespond(
					NewSignalTrigger(&BattleStartSignal{}),
					NewSelectActor(
						NewVerbActor(
							NewBuff(nil, NewBuffReactor(
								Damage,
								false,
								ConstEvaluator(102), // TODO:
								FatTags(Label("+2% Attack*")))),
							nil,
						),
						SideSelector(true),
					),
				),
			),
		},
	}
)
