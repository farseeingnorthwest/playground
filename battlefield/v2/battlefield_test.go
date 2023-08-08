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
	rng.SetRng(NewSequence(0.1, 0.5))
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
				T0[3],
				T0[4],
				T0[5],
			),
			NewMyWarrior(
				MyBaseline{
					Damage:  12,
					Defense: 8,
					Health:  200,
				},
				Right,
				0,
				T0[0],
			),
		},
	)
	f.Run()
	// Output:
	// verb=buff reactor="+20% Damage" target.side=Left target.position=0 source.reactor="@BattleStart({$} +20% Damage)"
	// verb=buff reactor="[15] +2% CriticalOdds" target.side=Left target.position=0 source.reactor="@Launch([15] +2% CriticalOdds)"
	// verb=attack critical=true loss=98 overflow=0 source.side=Left source.position=0 source.reactor="@Launch({*} 480% Damage, {1} 520% Damage)" source.damage=57 target.side=Right target.position=0 target.defense=8 target.health.current=102 target.health.maximum=200
	// verb=attack critical=false loss=54 overflow=0 source.side=Left source.position=0 source.reactor="@Launch({*} 480% Damage, {1} 520% Damage)" source.damage=62 target.side=Right target.position=0 target.defense=8 target.health.current=48 target.health.maximum=200
	// verb=attack critical=false loss=7 overflow=0 source.side=Right source.position=0 source.reactor=NormalAttack source.damage=12 target.side=Left target.position=0 target.defense=5 target.health.current=93 target.health.maximum=100
	// verb=buff reactor="[15] +2% CriticalOdds" target.side=Left target.position=0 source.reactor="@Launch([15] +2% CriticalOdds)"
	// verb=attack critical=false loss=47 overflow=0 source.side=Left source.position=0 source.reactor="@Launch({1} 3 * 460% Damage)" source.damage=55 target.side=Right target.position=0 target.defense=8 target.health.current=1 target.health.maximum=200
	// verb=attack critical=false loss=47 overflow=46 source.side=Left source.position=0 source.reactor="@Launch({1} 3 * 460% Damage)" source.damage=55 target.side=Right target.position=0 target.defense=8 target.health.current=0 target.health.maximum=200
	// verb=attack critical=false loss=47 overflow=47 source.side=Left source.position=0 source.reactor="@Launch({1} 3 * 460% Damage)" source.damage=55 target.side=Right target.position=0 target.defense=8 target.health.current=0 target.health.maximum=200
}
