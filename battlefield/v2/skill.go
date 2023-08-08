package battlefield

var (
	rng   = &RngProxy{}
	Skill = ExclusionGroup(0)
	T0    = []Reactor{
		// 0
		NewFatReactor(
			FatTags(Skill, Label("NormalAttack")),
			FatRespond(
				NewSignalTrigger(&LaunchSignal{}),
				NewSelectActor(
					NewVerbActor(&Attack{}, AxisEvaluator(Damage)),
					Healthy,
					SideSelector(false),
					NewShuffleSelector(rng, nil),
					FrontSelector(1),
				),
			),
		),

		// 1
		NewFatReactor(
			FatRespond(
				NewFatTrigger(
					&PreActionSignal{},
					CurrentIsSourceTrigger{},
					NewVerbTrigger[*Attack](),
				),
				NewProbabilityActor(
					rng,
					AxisEvaluator(CriticalOdds),
					NewSequenceActor(
						CriticalActor{},
						NewActionBuffer(
							AxisEvaluator(CriticalLoss),
							NewBuffer(Loss, false, nil),
						),
					),
				),
			),
		),

		// ////////////////////////////////////////////////////////////
		// 織田

		// 對隨機1名敵人進行3次攻擊，每次造成攻擊力460%的傷害。
		// 2
		NewFatReactor(
			FatTags(Skill, Priority(1), Label("@Launch({1} 3 * 460% Damage)")),
			FatRespond(
				NewSignalTrigger(&LaunchSignal{}),
				NewSelectActor(
					NewSequenceActor(
						NewVerbActor(&Attack{}, NewMultiplier(AxisEvaluator(Damage), 460)),
						NewVerbActor(&Attack{}, NewMultiplier(AxisEvaluator(Damage), 460)),
						NewVerbActor(&Attack{}, NewMultiplier(AxisEvaluator(Damage), 460)),
					),
					Healthy,
					SideSelector(false),
					NewShuffleSelector(rng, nil),
					FrontSelector(1),
				),
			),
			FatCooling(NewSignalTrigger(&RoundEndSignal{}), 5),
		),

		// 對全體敵人造成攻擊力480%的傷害。再對當前生命值最低的1名敵人造成攻擊力520%的傷害。
		// 3
		NewFatReactor(
			FatTags(Skill, Priority(2), Label("@Launch({*} 480% Damage, {1} 520% Damage)")),
			FatRespond(
				NewSignalTrigger(&LaunchSignal{}),
				NewSelectActor(
					NewVerbActor(&Attack{}, NewMultiplier(AxisEvaluator(Damage), 480)),
					Healthy,
					SideSelector(false),
				),
				NewSelectActor(
					NewVerbActor(&Attack{}, NewMultiplier(AxisEvaluator(Damage), 520)),
					Healthy,
					SideSelector(false),
					NewSortSelector(HealthPercent, true),
					FrontSelector(1),
				),
			),
			FatCooling(NewSignalTrigger(&RoundEndSignal{}), 4),
		),

		// 每次行動開始時，提升2%爆擊率(最高15層，無法被解除)。
		// 4
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
		// 5
		NewFatReactor(
			FatTags(Priority(4), Label("@BattleStart({$} +20% Damage)")),
			FatRespond(
				NewSignalTrigger(&BattleStartSignal{}),
				NewSelectActor(
					NewVerbActor(
						NewBuff(
							nil,
							NewBuffReactor(Damage, false, ConstEvaluator(120), FatTags(
								Label("+20% Damage"))),
						),
						nil,
					),
					CurrentSelector{},
				),
			),
		),
	}
)
