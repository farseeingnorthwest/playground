package battlefield

import "math/rand"

var (
	rng = rand.New(rand.NewSource(0))
	T0  = []Reactor{
		NewFatReactor(
			FatRespond(
				NewSignalTrigger(&LaunchSignal{}),
				NewVerbActor(&Attack{}, AxisEvaluator(Damage)),
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
