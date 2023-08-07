package battlefield

import (
	"os"

	"golang.org/x/exp/slog"
)

func ExampleBattleField_Run() {
	h := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if len(groups) == 0 {
				switch a.Key {
				case slog.TimeKey, slog.LevelKey, slog.MessageKey:
					return slog.Attr{}
				}
			}

			return a
		},
	})
	slog.SetDefault(slog.New(h))
	rng.SetRng(NewSequence(0.5, 0.001, 0.5))
	f := NewBattleField(
		[]Warrior{
			NewMyWarrior(
				MyBaseline{
					Damage:       10,
					Defense:      5,
					CriticalOdds: 10,
					CriticalLoss: 200,
					Health:       100,
				},
				Left,
				0,
				T0[0],
				T0[1],
				T0[2],
			),
			NewMyWarrior(
				MyBaseline{
					Damage:  12,
					Defense: 8,
					Health:  120,
				},
				Right,
				0,
				T0[0],
			),
		},
	)
	f.Run()
	// Output:
	// verb=attack critical=true loss=4 overflow=0 source.side=Left source.position=0 source.damage=10 target.side=Right target.position=0 target.defense=8 target.health.current=18 target.health.maximum=22
	// verb=attack critical=false loss=7 overflow=0 source.side=Right source.position=0 source.damage=12 target.side=Left target.position=0 target.defense=5 target.health.current=13 target.health.maximum=20
	// verb=attack critical=false loss=2 overflow=0 source.side=Left source.position=0 source.damage=10 target.side=Right target.position=0 target.defense=8 target.health.current=16 target.health.maximum=22
	// verb=attack critical=false loss=7 overflow=0 source.side=Right source.position=0 source.damage=12 target.side=Left target.position=0 target.defense=5 target.health.current=6 target.health.maximum=20
	// verb=attack critical=false loss=2 overflow=0 source.side=Left source.position=0 source.damage=10 target.side=Right target.position=0 target.defense=8 target.health.current=14 target.health.maximum=22
	// verb=attack critical=false loss=7 overflow=1 source.side=Right source.position=0 source.damage=12 target.side=Left target.position=0 target.defense=5 target.health.current=0 target.health.maximum=20
}
