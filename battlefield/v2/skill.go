package battlefield

import "math/rand"

var (
	rng = rand.New(rand.NewSource(0))
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
						NewVerbActor(
							NewBuff(
								nil,
								NewFatReactor(
									FatRespond(
										NewSignalTrigger(&EvaluationSignal{}),
										NewBuffer(Loss, nil),
									),
								),
								true,
							),
							AxisEvaluator(CriticalLoss),
						),
					),
				),
			),
		),
	}
)
