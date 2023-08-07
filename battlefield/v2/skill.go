package battlefield

var (
	rng   = &RngProxy{}
	Skill = ExclusionGroup(0)
	T0    = []Reactor{
		// 0
		NewFatReactor(
			FatTags(Skill),
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
				),
				NewProbabilityActor(
					rng,
					AxisEvaluator(CriticalOdds),
					NewSequenceActor(
						CriticalActor{},
						NewActionBuffer(
							AxisEvaluator(CriticalLoss),
							NewBuffer(Loss, nil),
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
			FatTags(Skill, Priority(1)),
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
	}
)
