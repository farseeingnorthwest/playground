package battlefield

import (
	"os"

	"golang.org/x/exp/slog"
)

func ExampleBattleField_Run_special_0() {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
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
	})))

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
				Regular[0],
				Regular[1],
				Regular[2],
				Special[0][0],
				Special[0][1],
				Special[0][2],
				Special[0][3],
			),
			NewMyWarrior(
				MyBaseline{
					Damage:  12,
					Defense: 8,
					Health:  200,
				},
				Right,
				0,
				Regular[0],
			),
		},
	)
	f.Run()
	// Output:
	// verb=buff reactor="+20% Damage" target.side=Left target.position=0 source.reactor="@BattleStart({$} +20% Damage)"
	// verb=buff reactor="[15] +2% CriticalOdds" target.side=Left target.position=0 source.reactor="@Launch([15] +2% CriticalOdds)"
	// verb=attack critical=true loss=98 overflow=0 source.side=Left source.position=0 source.reactor="@Launch({*} 480% Damage; {1} 520% Damage)" source.damage=57 target.side=Right target.position=0 target.defense=8 target.health.current=102 target.health.maximum=200
	// verb=attack critical=false loss=54 overflow=0 source.side=Left source.position=0 source.reactor="@Launch({*} 480% Damage; {1} 520% Damage)" source.damage=62 target.side=Right target.position=0 target.defense=8 target.health.current=48 target.health.maximum=200
	// verb=attack critical=false loss=7 overflow=0 source.side=Right source.position=0 source.reactor=NormalAttack source.damage=12 target.side=Left target.position=0 target.defense=5 target.health.current=93 target.health.maximum=100
	// verb=buff reactor="[15] +2% CriticalOdds" target.side=Left target.position=0 source.reactor="@Launch([15] +2% CriticalOdds)"
	// verb=attack critical=false loss=47 overflow=0 source.side=Left source.position=0 source.reactor="@Launch({1} 3 * 460% Damage)" source.damage=55 target.side=Right target.position=0 target.defense=8 target.health.current=1 target.health.maximum=200
	// verb=attack critical=false loss=47 overflow=46 source.side=Left source.position=0 source.reactor="@Launch({1} 3 * 460% Damage)" source.damage=55 target.side=Right target.position=0 target.defense=8 target.health.current=0 target.health.maximum=200
	// verb=attack critical=false loss=47 overflow=47 source.side=Left source.position=0 source.reactor="@Launch({1} 3 * 460% Damage)" source.damage=55 target.side=Right target.position=0 target.defense=8 target.health.current=0 target.health.maximum=200
}

func ExampleBattleField_Run_special_1() {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
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
	})))

	rng.SetRng(NewSequence(0.1, 0.9))
	f := NewBattleField(
		[]Warrior{
			NewMyWarrior(
				MyBaseline{
					Damage:  10,
					Defense: 5,
					Health:  100,
				},
				Left,
				0,
				Regular[0],
				Special[1][0],
				Special[1][1],
			),
			NewMyWarrior(
				MyBaseline{
					Damage:  12,
					Defense: 8,
					Health:  84,
				},
				Right,
				0,
				Regular[0],
			),
		},
	)
	f.Run()
	// Output:
	// verb=attack critical=false loss=44 overflow=0 source.side=Left source.position=0 source.reactor="@Launch({1} 520% Damage, P(70%, Dizzy))" source.damage=52 target.side=Right target.position=0 target.defense=8 target.health.current=40 target.health.maximum=84
	// verb=buff reactor=Dizzy target.side=Right target.position=0 source.reactor="@Launch({1} 520% Damage, P(70%, Dizzy))"
	// verb=attack critical=false loss=34 overflow=0 source.side=Left source.position=0 source.reactor="@Launch({2} 420% Damage, +30% Loss, Sleeping)" source.damage=42 target.side=Right target.position=0 target.defense=8 target.health.current=6 target.health.maximum=84
	// verb=buff reactor="+30% Loss" target.side=Right target.position=0 source.reactor="@Launch({2} 420% Damage, +30% Loss, Sleeping)"
	// verb=buff reactor=Sleeping target.side=Right target.position=0 source.reactor="@Launch({2} 420% Damage, +30% Loss, Sleeping)"
	// verb=attack critical=false loss=2 overflow=0 source.side=Left source.position=0 source.reactor=NormalAttack source.damage=10 target.side=Right target.position=0 target.defense=8 target.health.current=4 target.health.maximum=84
	// verb=attack critical=false loss=7 overflow=0 source.side=Right source.position=0 source.reactor=NormalAttack source.damage=12 target.side=Left target.position=0 target.defense=5 target.health.current=93 target.health.maximum=100
	// verb=attack critical=false loss=2 overflow=0 source.side=Left source.position=0 source.reactor=NormalAttack source.damage=10 target.side=Right target.position=0 target.defense=8 target.health.current=2 target.health.maximum=84
	// verb=attack critical=false loss=7 overflow=0 source.side=Right source.position=0 source.reactor=NormalAttack source.damage=12 target.side=Left target.position=0 target.defense=5 target.health.current=86 target.health.maximum=100
	// verb=attack critical=false loss=44 overflow=42 source.side=Left source.position=0 source.reactor="@Launch({1} 520% Damage, P(70%, Dizzy))" source.damage=52 target.side=Right target.position=0 target.defense=8 target.health.current=0 target.health.maximum=84
}

func ExampleBattleField_Run_special_2() {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
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
	})))

	rng.SetRng(NewSequence(0.1, 0.9))
	f := NewBattleField(
		[]Warrior{
			NewMyWarrior(
				MyBaseline{
					Damage:       10,
					Defense:      5,
					CriticalLoss: 200,
					Health:       100,
				},
				Left,
				0,
				Regular[0],
				Regular[1],
				Regular[2],
				Special[2][0],
				Special[2][1],
				Special[2][2],
				Special[2][3],
			),
			NewMyWarrior(
				MyBaseline{
					Damage:  12,
					Defense: 8,
					Health:  130,
				},
				Right,
				0,
				Regular[0],
			),
		},
	)
	f.Run()
	// Output:
	// verb=buff reactor="+2% Attack*" target.side=Left target.position=0 source.reactor="@BattleStart({~} +2% Attack*)"
	// verb=buff reactor="+25% HealthMaximum" target.side=Left target.position=0 source.reactor="@BattleStart({$} +25% HealthMaximum)"
	// verb=attack critical=true loss=64 overflow=0 source.side=Left source.position=0 source.reactor="@Launch({1!/%} 400% Critical Damage)" source.damage=40 target.side=Right target.position=0 target.defense=8 target.health.current=66 target.health.maximum=130
	// verb=attack critical=false loss=7 overflow=0 source.side=Right source.position=0 source.reactor=NormalAttack source.damage=12 target.side=Left target.position=0 target.defense=5 target.health.current=118 target.health.maximum=125
	// verb=attack critical=true loss=64 overflow=0 source.side=Left source.position=0 source.reactor="@Launch({1^/%} 400% Critical Damage)" source.damage=40 target.side=Right target.position=0 target.defense=8 target.health.current=2 target.health.maximum=130
	// verb=attack critical=false loss=7 overflow=0 source.side=Right source.position=0 source.reactor=NormalAttack source.damage=12 target.side=Left target.position=0 target.defense=5 target.health.current=111 target.health.maximum=125
	// verb=attack critical=false loss=2 overflow=0 source.side=Left source.position=0 source.reactor=NormalAttack source.damage=10 target.side=Right target.position=0 target.defense=8 target.health.current=0 target.health.maximum=130
}
