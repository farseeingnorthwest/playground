package battlefield

var (
	rng = &RngProxy{}
	T0  = []Reactor{
		NewFatReactor(
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
	}
)
