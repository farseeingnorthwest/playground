package battlefield_test

import (
	"log/slog"
	"os"

	. "github.com/farseeingnorthwest/playground/battlefield/v2"
)

func ExampleBattleField_Run_theory() {
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

	RngX.SetRng(NewSequence(0.5))
	f := NewBattleField(
		[]Warrior{
			NewMyWarrior(
				MyBaseline{
					Damage:  24,
					Defense: 5,
					Health:  100,
				},
				Left,
				0,
				WarriorSkills(Regular[0]),
				WarriorTags(Water),
			),
			NewMyWarrior(
				MyBaseline{
					Damage:  10,
					Defense: 4,
					Health:  40,
				},
				Right,
				0,
				WarriorSkills(Regular[0]),
				WarriorTags(Fire),
			),
			NewMyWarrior(
				MyBaseline{
					Damage:  10,
					Defense: 4,
					Health:  40,
				},
				Right,
				1,
				WarriorSkills(Regular[0]),
				WarriorTags(Thunder),
			),
		},
		Regular[3],
	)

	f.Run()
	// Output:
	// verb=attack critical=false loss=24 overflow=0 source.side=Left source.position=0 source.reactor=NormalAttack source.damage=24 target.side=Right target.position=0 target.defense=4 target.health.current=16 target.health.maximum=40
	// verb=attack critical=false loss=4 overflow=0 source.side=Right source.position=0 source.reactor=NormalAttack source.damage=10 target.side=Left target.position=0 target.defense=5 target.health.current=96 target.health.maximum=100
	// verb=attack critical=false loss=6 overflow=0 source.side=Right source.position=1 source.reactor=NormalAttack source.damage=10 target.side=Left target.position=0 target.defense=5 target.health.current=90 target.health.maximum=100
	// verb=attack critical=false loss=24 overflow=8 source.side=Left source.position=0 source.reactor=NormalAttack source.damage=24 target.side=Right target.position=0 target.defense=4 target.health.current=0 target.health.maximum=40
	// verb=attack critical=false loss=6 overflow=0 source.side=Right source.position=1 source.reactor=NormalAttack source.damage=10 target.side=Left target.position=0 target.defense=5 target.health.current=84 target.health.maximum=100
	// verb=attack critical=false loss=16 overflow=0 source.side=Left source.position=0 source.reactor=NormalAttack source.damage=24 target.side=Right target.position=1 target.defense=4 target.health.current=24 target.health.maximum=40
	// verb=attack critical=false loss=6 overflow=0 source.side=Right source.position=1 source.reactor=NormalAttack source.damage=10 target.side=Left target.position=0 target.defense=5 target.health.current=78 target.health.maximum=100
	// verb=attack critical=false loss=16 overflow=0 source.side=Left source.position=0 source.reactor=NormalAttack source.damage=24 target.side=Right target.position=1 target.defense=4 target.health.current=8 target.health.maximum=40
	// verb=attack critical=false loss=6 overflow=0 source.side=Right source.position=1 source.reactor=NormalAttack source.damage=10 target.side=Left target.position=0 target.defense=5 target.health.current=72 target.health.maximum=100
	// verb=attack critical=false loss=16 overflow=8 source.side=Left source.position=0 source.reactor=NormalAttack source.damage=24 target.side=Right target.position=1 target.defense=4 target.health.current=0 target.health.maximum=40
}

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

	RngX.SetRng(NewSequence(0.1, 0.5))
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
				WarriorSkills(
					Regular[0],
					Regular[1],
					Regular[2],
					Special[0][0].Fork(nil).(*FatReactor),
					Special[0][1].Fork(nil).(*FatReactor),
					Special[0][2].Fork(nil).(*FatReactor),
					Special[0][3].Fork(nil).(*FatReactor),
				),
			),
			NewMyWarrior(
				MyBaseline{
					Damage:  12,
					Defense: 8,
					Health:  200,
				},
				Right,
				0,
				WarriorSkills(
					Regular[0],
				),
			),
		},
	)

	f.Run()
	// Output:
	// verb=buff reactor="+20% Damage" target.side=Left target.position=0 source.reactor="@BattleStart({$} +20% Damage)"
	// signal=16 affairs=1 source.position=0 source.side=Left source.reactor="@Launch({*} 480% Damage; {1} 520% Damage)" lifecycle.leading=0 lifecycle.cooling="{Current:4 Maximum:4}" lifecycle.capacity=-1
	// stacking=1 verb=buff reactor="[15] +2% CriticalOdds" target.side=Left target.position=0 source.reactor="@Launch([15] +2% CriticalOdds)"
	// verb=attack critical=true loss=98 overflow=0 source.side=Left source.position=0 source.reactor="@Launch({*} 480% Damage; {1} 520% Damage)" source.damage=57 target.side=Right target.position=0 target.defense=8 target.health.current=102 target.health.maximum=200
	// verb=attack critical=false loss=54 overflow=0 source.side=Left source.position=0 source.reactor="@Launch({*} 480% Damage; {1} 520% Damage)" source.damage=62 target.side=Right target.position=0 target.defense=8 target.health.current=48 target.health.maximum=200
	// verb=attack critical=false loss=7 overflow=0 source.side=Right source.position=0 source.reactor=NormalAttack source.damage=12 target.side=Left target.position=0 target.defense=5 target.health.current=93 target.health.maximum=100
	// signal=74 affairs=0 source.position=0 source.side=Left source.reactor="@Launch({*} 480% Damage; {1} 520% Damage)" lifecycle.leading=0 lifecycle.cooling="{Current:3 Maximum:4}" lifecycle.capacity=-1
	// signal=85 affairs=1 source.position=0 source.side=Left source.reactor="@Launch({1} 3 * 460% Damage)" lifecycle.leading=0 lifecycle.cooling="{Current:5 Maximum:5}" lifecycle.capacity=-1
	// stacking=2 verb=buff reactor="[15] +2% CriticalOdds" target.side=Left target.position=0 source.reactor="@Launch([15] +2% CriticalOdds)"
	// verb=attack critical=false loss=47 overflow=0 source.side=Left source.position=0 source.reactor="@Launch({1} 3 * 460% Damage)" source.damage=55 target.side=Right target.position=0 target.defense=8 target.health.current=1 target.health.maximum=200
	// verb=attack critical=false loss=47 overflow=46 source.side=Left source.position=0 source.reactor="@Launch({1} 3 * 460% Damage)" source.damage=55 target.side=Right target.position=0 target.defense=8 target.health.current=0 target.health.maximum=200
	// targets=0 falseTargets=1 immuneTargets=0
}

func ExampleBattleField_Run_special_0b() {
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

	RngX.SetRng(NewSequence(0.1, 0.5))
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
				WarriorSkills(
					Regular[0],
					Regular[1],
					Regular[2],
					Special[0][0],
					Special[0][1],
					Special[0][2],
					Special[0][3],
				),
			),
			NewMyWarrior(
				MyBaseline{
					Damage:  12,
					Defense: 8,
					Health:  150,
				},
				Right,
				0,
				WarriorSkills(
					Regular[0],
				),
			),
			NewMyWarrior(
				MyBaseline{
					Damage:  12,
					Defense: 2,
					Health:  160,
				},
				Right,
				1,
				WarriorSkills(
					Regular[0],
				),
			),
		},
	)

	f.Run()
	// Output:
	// verb=buff reactor="+20% Damage" target.side=Left target.position=0 source.reactor="@BattleStart({$} +20% Damage)"
	// signal=34 affairs=1 source.position=0 source.side=Left source.reactor="@Launch({*} 480% Damage; {1} 520% Damage)" lifecycle.leading=0 lifecycle.cooling="{Current:4 Maximum:4}" lifecycle.capacity=-1
	// stacking=3 verb=buff reactor="[15] +2% CriticalOdds" target.side=Left target.position=0 source.reactor="@Launch([15] +2% CriticalOdds)"
	// verb=attack critical=true loss=98 overflow=0 source.side=Left source.position=0 source.reactor="@Launch({*} 480% Damage; {1} 520% Damage)" source.damage=57 target.side=Right target.position=0 target.defense=8 target.health.current=52 target.health.maximum=150
	// verb=attack critical=true loss=110 overflow=0 source.side=Left source.position=0 source.reactor="@Launch({*} 480% Damage; {1} 520% Damage)" source.damage=57 target.side=Right target.position=1 target.defense=2 target.health.current=50 target.health.maximum=160
	// verb=attack critical=false loss=60 overflow=10 source.side=Left source.position=0 source.reactor="@Launch({*} 480% Damage; {1} 520% Damage)" source.damage=62 target.side=Right target.position=1 target.defense=2 target.health.current=0 target.health.maximum=160
	// verb=attack critical=false loss=7 overflow=0 source.side=Right source.position=0 source.reactor=NormalAttack source.damage=12 target.side=Left target.position=0 target.defense=5 target.health.current=93 target.health.maximum=100
	// signal=139 affairs=0 source.position=0 source.side=Left source.reactor="@Launch({*} 480% Damage; {1} 520% Damage)" lifecycle.leading=0 lifecycle.cooling="{Current:3 Maximum:4}" lifecycle.capacity=-1
	// signal=164 affairs=1 source.position=0 source.side=Left source.reactor="@Launch({1} 3 * 460% Damage)" lifecycle.leading=0 lifecycle.cooling="{Current:5 Maximum:5}" lifecycle.capacity=-1
	// stacking=4 verb=buff reactor="[15] +2% CriticalOdds" target.side=Left target.position=0 source.reactor="@Launch([15] +2% CriticalOdds)"
	// verb=attack critical=false loss=47 overflow=0 source.side=Left source.position=0 source.reactor="@Launch({1} 3 * 460% Damage)" source.damage=55 target.side=Right target.position=0 target.defense=8 target.health.current=5 target.health.maximum=150
	// verb=attack critical=false loss=47 overflow=42 source.side=Left source.position=0 source.reactor="@Launch({1} 3 * 460% Damage)" source.damage=55 target.side=Right target.position=0 target.defense=8 target.health.current=0 target.health.maximum=150
	// targets=0 falseTargets=1 immuneTargets=0
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

	RngX.SetRng(NewSequence(0.1, 0.9))
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
				WarriorSkills(
					Regular[0],
					Special[1][0],
					Special[1][1],
					Special[1][2],
					Special[1][3],
				),
			),
			NewMyWarrior(
				MyBaseline{
					Damage:  20,
					Defense: 8,
					Health:  84,
				},
				Right,
				0,
				WarriorSkills(
					Regular[0],
					Scaffold[0].Fork(nil).(*FatReactor),
					Scaffold[1].Fork(nil).(*FatReactor),
				),
			),
		},
	)

	f.Run()
	// Output:
	// signal=9 affairs=1 source.position=0 source.side=Left source.reactor="@Launch({$} Barrier)" lifecycle.leading=0 lifecycle.cooling="{Current:4 Maximum:4}" lifecycle.capacity=-1
	// signal=9 affairs=1 source.position=0 source.side=Left source.reactor="@Launch({1} 520% Damage, P(70%, Stun))" lifecycle.leading=0 lifecycle.cooling="{Current:4 Maximum:4}" lifecycle.capacity=-1
	// verb=buff reactor=Barrier target.side=Left target.position=0 source.reactor="@Launch({$} Barrier)"
	// verb=attack critical=false loss=44 overflow=0 source.side=Left source.position=0 source.reactor="@Launch({1} 520% Damage, P(70%, Stun))" source.damage=52 target.side=Right target.position=0 target.defense=8 target.health.current=40 target.health.maximum=84
	// verb=buff reactor=Stun target.side=Right target.position=0 source.reactor="@Launch({1} 520% Damage, P(70%, Stun))"
	// signal=48 affairs=0 source.position=0 source.side=Left source.reactor="@Launch({$} Barrier)" lifecycle.leading=0 lifecycle.cooling="{Current:3 Maximum:4}" lifecycle.capacity=-1
	// signal=48 affairs=0 source.position=0 source.side=Left source.reactor="@Launch({1} 520% Damage, P(70%, Stun))" lifecycle.leading=0 lifecycle.cooling="{Current:3 Maximum:4}" lifecycle.capacity=-1
	// signal=48 affairs=0 source.position=0 source.side=Right source.reactor=Stun lifecycle.leading=0 lifecycle.cooling="{Current:0 Maximum:0}" lifecycle.capacity=0
	// signal=65 affairs=1 source.position=0 source.side=Left source.reactor="@Launch({2} 420% Damage, +30% Loss, Sleep)" lifecycle.leading=0 lifecycle.cooling="{Current:5 Maximum:5}" lifecycle.capacity=-1
	// verb=attack critical=false loss=34 overflow=0 source.side=Left source.position=0 source.reactor="@Launch({2} 420% Damage, +30% Loss, Sleep)" source.damage=42 target.side=Right target.position=0 target.defense=8 target.health.current=6 target.health.maximum=84
	// verb=buff reactor="+30% Loss" target.side=Right target.position=0 source.reactor="@Launch({2} 420% Damage, +30% Loss, Sleep)"
	// verb=buff reactor=Sleep target.side=Right target.position=0 source.reactor="@Launch({2} 420% Damage, +30% Loss, Sleep)"
	// signal=101 affairs=0 source.position=0 source.side=Left source.reactor="@Launch({$} Barrier)" lifecycle.leading=0 lifecycle.cooling="{Current:2 Maximum:4}" lifecycle.capacity=-1
	// signal=101 affairs=0 source.position=0 source.side=Left source.reactor="@Launch({1} 520% Damage, P(70%, Stun))" lifecycle.leading=0 lifecycle.cooling="{Current:2 Maximum:4}" lifecycle.capacity=-1
	// signal=101 affairs=0 source.position=0 source.side=Left source.reactor="@Launch({2} 420% Damage, +30% Loss, Sleep)" lifecycle.leading=0 lifecycle.cooling="{Current:4 Maximum:5}" lifecycle.capacity=-1
	// signal=101 affairs=0 source.position=0 source.side=Right source.reactor=Sleep lifecycle.leading=0 lifecycle.cooling="{Current:0 Maximum:0}" lifecycle.capacity=0
	// signal=101 affairs=0 source.position=0 source.side=Right source.reactor="+30% Loss" lifecycle.leading=0 lifecycle.cooling="{Current:0 Maximum:0}" lifecycle.capacity=0
	// verb=attack critical=false loss=2 overflow=0 source.side=Left source.position=0 source.reactor=NormalAttack source.damage=10 target.side=Right target.position=0 target.defense=8 target.health.current=4 target.health.maximum=84
	// signal=140 affairs=1 source.position=0 source.side=Right source.reactor=#2 lifecycle.leading=0 lifecycle.cooling="{Current:0 Maximum:0}" lifecycle.capacity=0
	// verb=buff reactor="Nerf #2" target.side=Left target.position=0 source.reactor=#2
	// signal=153 affairs=0 source.position=0 source.side=Left source.reactor="@Launch({$} Barrier)" lifecycle.leading=0 lifecycle.cooling="{Current:1 Maximum:4}" lifecycle.capacity=-1
	// signal=153 affairs=0 source.position=0 source.side=Left source.reactor="@Launch({1} 520% Damage, P(70%, Stun))" lifecycle.leading=0 lifecycle.cooling="{Current:1 Maximum:4}" lifecycle.capacity=-1
	// signal=153 affairs=0 source.position=0 source.side=Left source.reactor="@Launch({2} 420% Damage, +30% Loss, Sleep)" lifecycle.leading=0 lifecycle.cooling="{Current:3 Maximum:5}" lifecycle.capacity=-1
	// verb=attack critical=false loss=2 overflow=0 source.side=Left source.position=0 source.reactor=NormalAttack source.damage=10 target.side=Right target.position=0 target.defense=8 target.health.current=2 target.health.maximum=84
	// signal=186 affairs=1 source.position=0 source.side=Right source.reactor=#1 lifecycle.leading=0 lifecycle.cooling="{Current:0 Maximum:0}" lifecycle.capacity=0
	// verb=buff reactor="Nerf #1" target.side=Left target.position=0 source.reactor=#1
	// signal=194 affairs=1 source.position=0 source.side=Left source.reactor="@PostAction({$/<Nerf> >= 2} Purge())" lifecycle.leading=0 lifecycle.cooling="{Current:4 Maximum:4}" lifecycle.capacity=-1
	// verb=purge reactors="[Nerf #2 Nerf #1]"
	// signal=209 affairs=0 source.position=0 source.side=Left source.reactor="@Launch({$} Barrier)" lifecycle.leading=0 lifecycle.cooling="{Current:0 Maximum:4}" lifecycle.capacity=-1
	// signal=209 affairs=0 source.position=0 source.side=Left source.reactor="@Launch({1} 520% Damage, P(70%, Stun))" lifecycle.leading=0 lifecycle.cooling="{Current:0 Maximum:4}" lifecycle.capacity=-1
	// signal=209 affairs=0 source.position=0 source.side=Left source.reactor="@PostAction({$/<Nerf> >= 2} Purge())" lifecycle.leading=0 lifecycle.cooling="{Current:3 Maximum:4}" lifecycle.capacity=-1
	// signal=209 affairs=0 source.position=0 source.side=Left source.reactor="@Launch({2} 420% Damage, +30% Loss, Sleep)" lifecycle.leading=0 lifecycle.cooling="{Current:2 Maximum:5}" lifecycle.capacity=-1
	// signal=229 affairs=1 source.position=0 source.side=Left source.reactor="@Launch({$} Barrier)" lifecycle.leading=0 lifecycle.cooling="{Current:4 Maximum:4}" lifecycle.capacity=-1
	// signal=229 affairs=1 source.position=0 source.side=Left source.reactor="@Launch({1} 520% Damage, P(70%, Stun))" lifecycle.leading=0 lifecycle.cooling="{Current:4 Maximum:4}" lifecycle.capacity=-1
	// verb=buff reactor=Barrier target.side=Left target.position=0 source.reactor="@Launch({$} Barrier)"
	// verb=attack critical=false loss=44 overflow=42 source.side=Left source.position=0 source.reactor="@Launch({1} 520% Damage, P(70%, Stun))" source.damage=52 target.side=Right target.position=0 target.defense=8 target.health.current=0 target.health.maximum=84
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

	RngX.SetRng(NewSequence(0.1, 0.9))
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
				WarriorSkills(
					Regular[0],
					Regular[1],
					Regular[2],
					Special[2][0],
					Special[2][1],
					Special[2][2],
					Special[2][3],
				),
			),
			NewMyWarrior(
				MyBaseline{
					Damage:  12,
					Defense: 8,
					Health:  130,
				},
				Right,
				0,
				WarriorSkills(Regular[0]),
			),
		},
	)

	f.Run()
	// Output:
	// verb=buff reactor="+2% Attack*" target.side=Left target.position=0 source.reactor="@BattleStart({~} +2% Attack*)"
	// verb=buff reactor="+25% HealthMaximum" target.side=Left target.position=0 source.reactor="@BattleStart({$} +25% HealthMaximum)"
	// signal=23 affairs=1 source.position=0 source.side=Left source.reactor="@Launch({1!/%} 400% Critical Damage)" lifecycle.leading=0 lifecycle.cooling="{Current:3 Maximum:3}" lifecycle.capacity=-1
	// verb=attack critical=true loss=64 overflow=0 source.side=Left source.position=0 source.reactor="@Launch({1!/%} 400% Critical Damage)" source.damage=40 target.side=Right target.position=0 target.defense=8 target.health.current=66 target.health.maximum=130
	// verb=attack critical=false loss=7 overflow=0 source.side=Right source.position=0 source.reactor=NormalAttack source.damage=12 target.side=Left target.position=0 target.defense=5 target.health.current=118 target.health.maximum=125
	// signal=62 affairs=0 source.position=0 source.side=Left source.reactor="@Launch({1!/%} 400% Critical Damage)" lifecycle.leading=0 lifecycle.cooling="{Current:2 Maximum:3}" lifecycle.capacity=-1
	// signal=73 affairs=1 source.position=0 source.side=Left source.reactor="@Launch({1^/%} 400% Critical Damage)" lifecycle.leading=0 lifecycle.cooling="{Current:3 Maximum:3}" lifecycle.capacity=-1
	// verb=attack critical=true loss=64 overflow=0 source.side=Left source.position=0 source.reactor="@Launch({1^/%} 400% Critical Damage)" source.damage=40 target.side=Right target.position=0 target.defense=8 target.health.current=2 target.health.maximum=130
	// verb=attack critical=false loss=7 overflow=0 source.side=Right source.position=0 source.reactor=NormalAttack source.damage=12 target.side=Left target.position=0 target.defense=5 target.health.current=111 target.health.maximum=125
	// signal=112 affairs=0 source.position=0 source.side=Left source.reactor="@Launch({1!/%} 400% Critical Damage)" lifecycle.leading=0 lifecycle.cooling="{Current:1 Maximum:3}" lifecycle.capacity=-1
	// signal=112 affairs=0 source.position=0 source.side=Left source.reactor="@Launch({1^/%} 400% Critical Damage)" lifecycle.leading=0 lifecycle.cooling="{Current:2 Maximum:3}" lifecycle.capacity=-1
	// verb=attack critical=false loss=2 overflow=0 source.side=Left source.position=0 source.reactor=NormalAttack source.damage=10 target.side=Right target.position=0 target.defense=8 target.health.current=0 target.health.maximum=130
}

func ExampleBattleField_Run_special_3() {
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

	RngX.SetRng(NewSequence(0.5, 0.5, 0.5, 0.5, 0.01, 0.9))
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
				WarriorSkills(
					Regular[0],
					Regular[1],
					Regular[2],
					Special[3][0],
					Special[3][1],
					Special[3][2],
					Special[3][3],
				),
			),
			NewMyWarrior(
				MyBaseline{
					Damage:  60,
					Defense: 8,
					Health:  110,
				},
				Right,
				0,
				WarriorSkills(Regular[0]),
			),
		},
	)

	f.Run()
	// Output:
	// signal=9 affairs=1 source.position=0 source.side=Left source.reactor="@Launch({3!/%} 250% Damage+)" lifecycle.leading=0 lifecycle.cooling="{Current:4 Maximum:4}" lifecycle.capacity=-1
	// verb=heal rise=25 overflow=25 source.side=Left source.position=0 source.reactor="@Launch({3!/%} 250% Damage+)" target.side=Left target.position=0 target.health.current=100 target.health.maximum=100
	// verb=attack critical=false loss=55 overflow=0 source.side=Right source.position=0 source.reactor=NormalAttack source.damage=60 target.side=Left target.position=0 target.defense=5 target.health.current=45 target.health.maximum=100
	// verb=buff reactor=Sanctuary target.side=Left target.position=0 source.reactor="@PostAction({$/< 50%} Sanctuary)"
	// signal=51 affairs=0 source.position=0 source.side=Left source.reactor="@Launch({3!/%} 250% Damage+)" lifecycle.leading=0 lifecycle.cooling="{Current:3 Maximum:4}" lifecycle.capacity=-1
	// signal=62 affairs=1 source.position=0 source.side=Left source.reactor="@Launch({1^/D} 4 * 340% Damage)" lifecycle.leading=0 lifecycle.cooling="{Current:4 Maximum:4}" lifecycle.capacity=-1
	// verb=attack critical=false loss=26 overflow=0 source.side=Left source.position=0 source.reactor="@Launch({1^/D} 4 * 340% Damage)" source.damage=34 target.side=Right target.position=0 target.defense=8 target.health.current=84 target.health.maximum=110
	// verb=attack critical=false loss=26 overflow=0 source.side=Left source.position=0 source.reactor="@Launch({1^/D} 4 * 340% Damage)" source.damage=34 target.side=Right target.position=0 target.defense=8 target.health.current=58 target.health.maximum=110
	// verb=attack critical=false loss=26 overflow=0 source.side=Left source.position=0 source.reactor="@Launch({1^/D} 4 * 340% Damage)" source.damage=34 target.side=Right target.position=0 target.defense=8 target.health.current=32 target.health.maximum=110
	// verb=attack critical=false loss=26 overflow=0 source.side=Left source.position=0 source.reactor="@Launch({1^/D} 4 * 340% Damage)" source.damage=34 target.side=Right target.position=0 target.defense=8 target.health.current=6 target.health.maximum=110
	// signal=134 affairs=1 source.position=0 source.side=Left source.reactor=Sanctuary lifecycle.leading=0 lifecycle.cooling="{Current:0 Maximum:0}" lifecycle.capacity=-1
	// verb=attack critical=false loss=30 overflow=0 source.side=Right source.position=0 source.reactor=NormalAttack source.damage=60 target.side=Left target.position=0 target.defense=5 target.health.current=15 target.health.maximum=100
	// signal=146 affairs=0 source.position=0 source.side=Left source.reactor="@Launch({3!/%} 250% Damage+)" lifecycle.leading=0 lifecycle.cooling="{Current:2 Maximum:4}" lifecycle.capacity=-1
	// signal=146 affairs=0 source.position=0 source.side=Left source.reactor="@Launch({1^/D} 4 * 340% Damage)" lifecycle.leading=0 lifecycle.cooling="{Current:3 Maximum:4}" lifecycle.capacity=-1
	// stacking=1 verb=buff reactor="+5% Attack*" target.side=Left target.position=0 source.reactor="@PreAction({$} +5% Attack*)"
	// verb=attack critical=true loss=4 overflow=0 source.side=Left source.position=0 source.reactor=NormalAttack source.damage=10 target.side=Right target.position=0 target.defense=8 target.health.current=2 target.health.maximum=110
	// signal=196 affairs=1 source.position=0 source.side=Left source.reactor=Sanctuary lifecycle.leading=0 lifecycle.cooling="{Current:0 Maximum:0}" lifecycle.capacity=-1
	// verb=attack critical=false loss=30 overflow=15 source.side=Right source.position=0 source.reactor=NormalAttack source.damage=60 target.side=Left target.position=0 target.defense=5 target.health.current=0 target.health.maximum=100
}

func ExampleBattleField_Run_special_3b() {
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

	RngX.SetRng(NewSequence(0.01))
	f := NewBattleField(
		[]Warrior{
			NewMyWarrior(
				MyBaseline{
					Damage:       20,
					Defense:      5,
					CriticalOdds: 10,
					CriticalLoss: 200,
					Health:       100,
				},
				Left,
				0,
				WarriorSkills(
					Regular[0],
					Regular[1],
					Special[3][2],
				),
			),
			NewMyWarrior(
				MyBaseline{
					Damage:  22,
					Defense: 8,
					Health:  60,
				},
				Right,
				0,
				WarriorSkills(Regular[0]),
			),
		},
	)

	f.Run()
	// Output:
	// stacking=2 verb=buff reactor="+5% Attack*" target.side=Left target.position=0 source.reactor="@PreAction({$} +5% Attack*)"
	// verb=attack critical=true loss=13 overflow=0 source.side=Left source.position=0 source.reactor=NormalAttack source.damage=21 target.side=Right target.position=0 target.defense=8 target.health.current=47 target.health.maximum=60
	// verb=attack critical=false loss=17 overflow=0 source.side=Right source.position=0 source.reactor=NormalAttack source.damage=22 target.side=Left target.position=0 target.defense=5 target.health.current=83 target.health.maximum=100
	// signal=49 affairs=0 source.position=0 source.side=Left source.reactor="+5% Attack*" lifecycle.leading=0 lifecycle.cooling="{Current:0 Maximum:0}" lifecycle.capacity=2
	// stacking=3 verb=buff reactor="+5% Attack*" target.side=Left target.position=0 source.reactor="@PreAction({$} +5% Attack*)"
	// verb=attack critical=true loss=14 overflow=0 source.side=Left source.position=0 source.reactor=NormalAttack source.damage=22 target.side=Right target.position=0 target.defense=8 target.health.current=33 target.health.maximum=60
	// verb=attack critical=false loss=17 overflow=0 source.side=Right source.position=0 source.reactor=NormalAttack source.damage=22 target.side=Left target.position=0 target.defense=5 target.health.current=66 target.health.maximum=100
	// signal=100 affairs=0 source.position=0 source.side=Left source.reactor="+5% Attack*" lifecycle.leading=0 lifecycle.cooling="{Current:0 Maximum:0}" lifecycle.capacity=1
	// signal=100 affairs=0 source.position=0 source.side=Left source.reactor="+5% Attack*" lifecycle.leading=0 lifecycle.cooling="{Current:0 Maximum:0}" lifecycle.capacity=2
	// verb=buff/overflow reactor="+5% Attack*" target.side=Left target.position=0
	// stacking=3 verb=buff reactor="+5% Attack*" target.side=Left target.position=0 source.reactor="@PreAction({$} +5% Attack*)"
	// verb=attack critical=true loss=15 overflow=0 source.side=Left source.position=0 source.reactor=NormalAttack source.damage=23 target.side=Right target.position=0 target.defense=8 target.health.current=18 target.health.maximum=60
	// verb=attack critical=false loss=17 overflow=0 source.side=Right source.position=0 source.reactor=NormalAttack source.damage=22 target.side=Left target.position=0 target.defense=5 target.health.current=49 target.health.maximum=100
	// signal=157 affairs=0 source.position=0 source.side=Left source.reactor="+5% Attack*" lifecycle.leading=0 lifecycle.cooling="{Current:0 Maximum:0}" lifecycle.capacity=0
	// signal=157 affairs=0 source.position=0 source.side=Left source.reactor="+5% Attack*" lifecycle.leading=0 lifecycle.cooling="{Current:0 Maximum:0}" lifecycle.capacity=1
	// signal=157 affairs=0 source.position=0 source.side=Left source.reactor="+5% Attack*" lifecycle.leading=0 lifecycle.cooling="{Current:0 Maximum:0}" lifecycle.capacity=2
	// verb=buff/overflow reactor="+5% Attack*" target.side=Left target.position=0
	// stacking=3 verb=buff reactor="+5% Attack*" target.side=Left target.position=0 source.reactor="@PreAction({$} +5% Attack*)"
	// verb=attack critical=true loss=15 overflow=0 source.side=Left source.position=0 source.reactor=NormalAttack source.damage=23 target.side=Right target.position=0 target.defense=8 target.health.current=3 target.health.maximum=60
	// verb=attack critical=false loss=17 overflow=0 source.side=Right source.position=0 source.reactor=NormalAttack source.damage=22 target.side=Left target.position=0 target.defense=5 target.health.current=32 target.health.maximum=100
	// signal=217 affairs=0 source.position=0 source.side=Left source.reactor="+5% Attack*" lifecycle.leading=0 lifecycle.cooling="{Current:0 Maximum:0}" lifecycle.capacity=0
	// signal=217 affairs=0 source.position=0 source.side=Left source.reactor="+5% Attack*" lifecycle.leading=0 lifecycle.cooling="{Current:0 Maximum:0}" lifecycle.capacity=1
	// signal=217 affairs=0 source.position=0 source.side=Left source.reactor="+5% Attack*" lifecycle.leading=0 lifecycle.cooling="{Current:0 Maximum:0}" lifecycle.capacity=2
	// verb=buff/overflow reactor="+5% Attack*" target.side=Left target.position=0
	// stacking=3 verb=buff reactor="+5% Attack*" target.side=Left target.position=0 source.reactor="@PreAction({$} +5% Attack*)"
	// verb=attack critical=true loss=15 overflow=12 source.side=Left source.position=0 source.reactor=NormalAttack source.damage=23 target.side=Right target.position=0 target.defense=8 target.health.current=0 target.health.maximum=60
}

func ExampleBattleField_Run_special_4() {
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

	RngX.SetRng(NewSequence(0.5))
	f := NewBattleField(
		[]Warrior{
			NewMyWarrior(
				MyBaseline{
					Damage:       12,
					Defense:      5,
					CriticalOdds: 10,
					CriticalLoss: 200,
					Health:       100,
				},
				Left,
				0,
				WarriorSkills(
					Regular[0],
					Regular[1],
					Special[4][0],
					Special[4][1],
					Special[4][2],
				),
			),
			NewMyWarrior(
				MyBaseline{
					Damage:  20,
					Defense: 8,
					Health:  72,
				},
				Right,
				0,
				WarriorSkills(Regular[0]),
			),
		},
	)

	f.Run()
	// Output:
	// signal=9 affairs=1 source.position=0 source.side=Left source.reactor="@Launch({$/>= 60%}, -15% Damage, Taunt)" lifecycle.leading=0 lifecycle.cooling="{Current:4 Maximum:4}" lifecycle.capacity=-1
	// signal=9 affairs=1 source.position=0 source.side=Left source.reactor="@Launch({1} 560% Damage, -25% Damage)" lifecycle.leading=0 lifecycle.cooling="{Current:4 Maximum:4}" lifecycle.capacity=-1
	// verb=buff reactor=Taunt target.side=Left target.position=0 source.reactor="@Launch({$/>= 60%}, -15% Damage, Taunt)"
	// verb=buff reactor="-15% Damage" target.side=Left target.position=0 source.reactor="@Launch({$/>= 60%}, -15% Damage, Taunt)"
	// verb=attack critical=false loss=48 overflow=0 source.side=Left source.position=0 source.reactor="@Launch({1} 560% Damage, -25% Damage)" source.damage=56 target.side=Right target.position=0 target.defense=8 target.health.current=24 target.health.maximum=72
	// verb=buff reactor="-25% Damage" target.side=Right target.position=0 source.reactor="@Launch({1} 560% Damage, -25% Damage)"
	// verb=attack critical=false loss=10 overflow=0 source.side=Right source.position=0 source.reactor=NormalAttack source.damage=15 target.side=Left target.position=0 target.defense=5 target.health.current=90 target.health.maximum=100
	// signal=69 affairs=0 source.position=0 source.side=Left source.reactor="@Launch({$/>= 60%}, -15% Damage, Taunt)" lifecycle.leading=0 lifecycle.cooling="{Current:3 Maximum:4}" lifecycle.capacity=-1
	// signal=69 affairs=0 source.position=0 source.side=Left source.reactor="@Launch({1} 560% Damage, -25% Damage)" lifecycle.leading=0 lifecycle.cooling="{Current:3 Maximum:4}" lifecycle.capacity=-1
	// signal=69 affairs=0 source.position=0 source.side=Left source.reactor=Taunt lifecycle.leading=0 lifecycle.cooling="{Current:0 Maximum:0}" lifecycle.capacity=1
	// signal=69 affairs=0 source.position=0 source.side=Left source.reactor="-15% Damage" lifecycle.leading=0 lifecycle.cooling="{Current:0 Maximum:0}" lifecycle.capacity=1
	// signal=69 affairs=0 source.position=0 source.side=Right source.reactor="-25% Damage" lifecycle.leading=0 lifecycle.cooling="{Current:0 Maximum:0}" lifecycle.capacity=1
	// signal=92 affairs=1 source.position=0 source.side=Left source.reactor="@Launch({*} 300% Damage)" lifecycle.leading=0 lifecycle.cooling="{Current:5 Maximum:5}" lifecycle.capacity=-1
	// verb=attack critical=false loss=22 overflow=0 source.side=Left source.position=0 source.reactor="@Launch({*} 300% Damage)" source.damage=30 target.side=Right target.position=0 target.defense=8 target.health.current=2 target.health.maximum=72
	// verb=attack critical=false loss=10 overflow=0 source.side=Right source.position=0 source.reactor=NormalAttack source.damage=15 target.side=Left target.position=0 target.defense=5 target.health.current=80 target.health.maximum=100
	// signal=128 affairs=0 source.position=0 source.side=Left source.reactor="@Launch({$/>= 60%}, -15% Damage, Taunt)" lifecycle.leading=0 lifecycle.cooling="{Current:2 Maximum:4}" lifecycle.capacity=-1
	// signal=128 affairs=0 source.position=0 source.side=Left source.reactor="@Launch({1} 560% Damage, -25% Damage)" lifecycle.leading=0 lifecycle.cooling="{Current:2 Maximum:4}" lifecycle.capacity=-1
	// signal=128 affairs=0 source.position=0 source.side=Left source.reactor="@Launch({*} 300% Damage)" lifecycle.leading=0 lifecycle.cooling="{Current:4 Maximum:5}" lifecycle.capacity=-1
	// signal=128 affairs=0 source.position=0 source.side=Left source.reactor=Taunt lifecycle.leading=0 lifecycle.cooling="{Current:0 Maximum:0}" lifecycle.capacity=0
	// signal=128 affairs=0 source.position=0 source.side=Left source.reactor="-15% Damage" lifecycle.leading=0 lifecycle.cooling="{Current:0 Maximum:0}" lifecycle.capacity=0
	// signal=128 affairs=0 source.position=0 source.side=Right source.reactor="-25% Damage" lifecycle.leading=0 lifecycle.cooling="{Current:0 Maximum:0}" lifecycle.capacity=0
	// verb=attack critical=false loss=4 overflow=2 source.side=Left source.position=0 source.reactor=NormalAttack source.damage=12 target.side=Right target.position=0 target.defense=8 target.health.current=0 target.health.maximum=72
}

func ExampleBattleField_Run_special_5() {
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

	RngX.SetRng(NewSequence(0.5, 0.1, 0.5))
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
				WarriorSkills(
					Regular[0],
					Special[5][0],
					Special[5][1],
					Special[5][2],
					Special[5][3],
				),
			),
			NewMyWarrior(
				MyBaseline{
					Damage:  60,
					Defense: 8,
					Health:  255,
				},
				Right,
				0,
				WarriorSkills(Regular[0]),
			),
		},
	)

	f.Run()
	// Output:
	// verb=buff reactor="30% Damage" target.side=Left target.position=0 source.reactor="@BattleStart({$}, 30% Damage)"
	// signal=17 affairs=1 source.position=0 source.side=Left source.reactor="@Launch({3} 510% Damage)" lifecycle.leading=0 lifecycle.cooling="{Current:4 Maximum:4}" lifecycle.capacity=-1
	// verb=attack critical=false loss=58 overflow=0 source.side=Left source.position=0 source.reactor="@Launch({3} 510% Damage)" source.damage=66 target.side=Right target.position=0 target.defense=8 target.health.current=197 target.health.maximum=255
	// verb=attack critical=false loss=55 overflow=0 source.side=Right source.position=0 source.reactor=NormalAttack source.damage=60 target.side=Left target.position=0 target.defense=5 target.health.current=45 target.health.maximum=100
	// verb=buff reactor="-20% Loss" target.side=Left target.position=0 source.reactor="@PostAction({$/< 50%}, -20% Loss)"
	// signal=62 affairs=0 source.position=0 source.side=Left source.reactor="@Launch({3} 510% Damage)" lifecycle.leading=0 lifecycle.cooling="{Current:3 Maximum:4}" lifecycle.capacity=-1
	// signal=73 affairs=1 source.position=0 source.side=Left source.reactor="@Launch({1} 3 * 550% Damage, P(50%) Stun)" lifecycle.leading=0 lifecycle.cooling="{Current:4 Maximum:4}" lifecycle.capacity=-1
	// verb=attack critical=false loss=63 overflow=0 source.side=Left source.position=0 source.reactor="@Launch({1} 3 * 550% Damage, P(50%) Stun)" source.damage=71 target.side=Right target.position=0 target.defense=8 target.health.current=134 target.health.maximum=255
	// verb=attack critical=false loss=63 overflow=0 source.side=Left source.position=0 source.reactor="@Launch({1} 3 * 550% Damage, P(50%) Stun)" source.damage=71 target.side=Right target.position=0 target.defense=8 target.health.current=71 target.health.maximum=255
	// verb=buff reactor=Stun target.side=Right target.position=0 source.reactor="@Launch({1} 3 * 550% Damage, P(50%) Stun)"
	// verb=attack critical=false loss=63 overflow=0 source.side=Left source.position=0 source.reactor="@Launch({1} 3 * 550% Damage, P(50%) Stun)" source.damage=71 target.side=Right target.position=0 target.defense=8 target.health.current=8 target.health.maximum=255
	// signal=130 affairs=0 source.position=0 source.side=Left source.reactor="@Launch({3} 510% Damage)" lifecycle.leading=0 lifecycle.cooling="{Current:2 Maximum:4}" lifecycle.capacity=-1
	// signal=130 affairs=0 source.position=0 source.side=Left source.reactor="@Launch({1} 3 * 550% Damage, P(50%) Stun)" lifecycle.leading=0 lifecycle.cooling="{Current:3 Maximum:4}" lifecycle.capacity=-1
	// signal=130 affairs=0 source.position=0 source.side=Right source.reactor=Stun lifecycle.leading=0 lifecycle.cooling="{Current:0 Maximum:0}" lifecycle.capacity=0
	// verb=attack critical=false loss=5 overflow=0 source.side=Left source.position=0 source.reactor=NormalAttack source.damage=13 target.side=Right target.position=0 target.defense=8 target.health.current=3 target.health.maximum=255
	// verb=attack critical=false loss=44 overflow=0 source.side=Right source.position=0 source.reactor=NormalAttack source.damage=60 target.side=Left target.position=0 target.defense=5 target.health.current=1 target.health.maximum=100
	// signal=181 affairs=0 source.position=0 source.side=Left source.reactor="@Launch({3} 510% Damage)" lifecycle.leading=0 lifecycle.cooling="{Current:1 Maximum:4}" lifecycle.capacity=-1
	// signal=181 affairs=0 source.position=0 source.side=Left source.reactor="@Launch({1} 3 * 550% Damage, P(50%) Stun)" lifecycle.leading=0 lifecycle.cooling="{Current:2 Maximum:4}" lifecycle.capacity=-1
	// verb=attack critical=false loss=5 overflow=2 source.side=Left source.position=0 source.reactor=NormalAttack source.damage=13 target.side=Right target.position=0 target.defense=8 target.health.current=0 target.health.maximum=255
}

func ExampleBattleField_Run_special_6() {
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

	RngX.SetRng(NewSequence(0.01, 0.5))
	f := NewBattleField(
		[]Warrior{
			NewMyWarrior(
				MyBaseline{
					Damage:  20,
					Defense: 5,
					Health:  100,
				},
				Left,
				0,
				WarriorSkills(
					Regular[0],
					Special[6][0],
					Special[6][1],
					Special[6][2],
					Special[6][3],
				),
			),
			NewMyWarrior(
				MyBaseline{
					Damage:       40,
					Defense:      8,
					Health:       270,
					CriticalOdds: 10,
					CriticalLoss: 200,
				},
				Right,
				0,
				WarriorSkills(
					Regular[0],
					Regular[1],
					Regular[2],
				),
			),
		},
	)

	f.Run()
	// Output:
	// verb=buff reactor="20% HealthMaximum" target.side=Left target.position=0 source.reactor="@BattleStart({$}, 20% HealthMaximum)"
	// signal=16 affairs=1 source.position=0 source.side=Left source.reactor="@Launch({1} 4 * 350% Damage; {*} +40% Shield)" lifecycle.leading=0 lifecycle.cooling="{Current:4 Maximum:4}" lifecycle.capacity=-1
	// verb=attack critical=false loss=62 overflow=0 source.side=Left source.position=0 source.reactor="@Launch({1} 4 * 350% Damage; {*} +40% Shield)" source.damage=70 target.side=Right target.position=0 target.defense=8 target.health.current=208 target.health.maximum=270
	// verb=attack critical=false loss=62 overflow=0 source.side=Left source.position=0 source.reactor="@Launch({1} 4 * 350% Damage; {*} +40% Shield)" source.damage=70 target.side=Right target.position=0 target.defense=8 target.health.current=146 target.health.maximum=270
	// verb=attack critical=false loss=62 overflow=0 source.side=Left source.position=0 source.reactor="@Launch({1} 4 * 350% Damage; {*} +40% Shield)" source.damage=70 target.side=Right target.position=0 target.defense=8 target.health.current=84 target.health.maximum=270
	// verb=attack critical=false loss=62 overflow=0 source.side=Left source.position=0 source.reactor="@Launch({1} 4 * 350% Damage; {*} +40% Shield)" source.damage=70 target.side=Right target.position=0 target.defense=8 target.health.current=22 target.health.maximum=270
	// capacity=99 verb=buff reactor=Shield target.side=Left target.position=0 source.reactor="@Launch({1} 4 * 350% Damage; {*} +40% Shield)"
	// signal=90 affairs=1 source.position=0 source.side=Left source.reactor=Shield lifecycle.leading=0 lifecycle.cooling="{Current:0 Maximum:0}" lifecycle.capacity=29
	// verb=attack critical=true loss=0 overflow=0 source.side=Right source.position=0 source.reactor=NormalAttack source.damage=40 target.side=Left target.position=0 target.defense=5 target.health.current=120 target.health.maximum=120
	// verb=buff reactor="-15% Damage" target.side=Right target.position=0 source.reactor="@PostAction({&/C}, -15% Damage)"
	// signal=107 affairs=0 source.position=0 source.side=Left source.reactor="@Launch({1} 4 * 350% Damage; {*} +40% Shield)" lifecycle.leading=0 lifecycle.cooling="{Current:3 Maximum:4}" lifecycle.capacity=-1
	// signal=107 affairs=0 source.position=0 source.side=Right source.reactor="-15% Damage" lifecycle.leading=0 lifecycle.cooling="{Current:0 Maximum:0}" lifecycle.capacity=0
	// signal=121 affairs=1 source.position=0 source.side=Left source.reactor="@Launch({~}, 15% Damage*)" lifecycle.leading=0 lifecycle.cooling="{Current:4 Maximum:4}" lifecycle.capacity=-1
	// verb=buff reactor="15% Damage*" target.side=Left target.position=0 source.reactor="@Launch({~}, 15% Damage*)"
	// signal=146 affairs=1 source.position=0 source.side=Left source.reactor=Shield lifecycle.leading=0 lifecycle.cooling="{Current:0 Maximum:0}" lifecycle.capacity=0
	// verb=attack critical=false loss=6 overflow=0 source.side=Right source.position=0 source.reactor=NormalAttack source.damage=40 target.side=Left target.position=0 target.defense=5 target.health.current=114 target.health.maximum=120
	// signal=156 affairs=0 source.position=0 source.side=Left source.reactor="@Launch({1} 4 * 350% Damage; {*} +40% Shield)" lifecycle.leading=0 lifecycle.cooling="{Current:2 Maximum:4}" lifecycle.capacity=-1
	// signal=156 affairs=0 source.position=0 source.side=Left source.reactor="@Launch({~}, 15% Damage*)" lifecycle.leading=0 lifecycle.cooling="{Current:3 Maximum:4}" lifecycle.capacity=-1
	// signal=156 affairs=0 source.position=0 source.side=Left source.reactor="15% Damage*" lifecycle.leading=0 lifecycle.cooling="{Current:0 Maximum:0}" lifecycle.capacity=1
	// verb=attack critical=false loss=15 overflow=0 source.side=Left source.position=0 source.reactor=NormalAttack source.damage=23 target.side=Right target.position=0 target.defense=8 target.health.current=7 target.health.maximum=270
	// verb=attack critical=false loss=35 overflow=0 source.side=Right source.position=0 source.reactor=NormalAttack source.damage=40 target.side=Left target.position=0 target.defense=5 target.health.current=79 target.health.maximum=120
	// signal=206 affairs=0 source.position=0 source.side=Left source.reactor="@Launch({1} 4 * 350% Damage; {*} +40% Shield)" lifecycle.leading=0 lifecycle.cooling="{Current:1 Maximum:4}" lifecycle.capacity=-1
	// signal=206 affairs=0 source.position=0 source.side=Left source.reactor="@Launch({~}, 15% Damage*)" lifecycle.leading=0 lifecycle.cooling="{Current:2 Maximum:4}" lifecycle.capacity=-1
	// signal=206 affairs=0 source.position=0 source.side=Left source.reactor="15% Damage*" lifecycle.leading=0 lifecycle.cooling="{Current:0 Maximum:0}" lifecycle.capacity=0
	// verb=attack critical=false loss=12 overflow=5 source.side=Left source.position=0 source.reactor=NormalAttack source.damage=20 target.side=Right target.position=0 target.defense=8 target.health.current=0 target.health.maximum=270
}

func ExampleBattleField_Run_special_7() {
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

	RngX.SetRng(NewSequence(0.5, 0.1))
	f := NewBattleField(
		[]Warrior{
			NewMyWarrior(
				MyBaseline{
					Damage:  20,
					Defense: 5,
					Health:  100,
				},
				Left,
				0,
				WarriorSkills(
					Regular[0],
					Special[7][0],
					Special[7][1],
				),
			),
			NewMyWarrior(
				MyBaseline{
					Damage:       40,
					Defense:      8,
					Health:       270,
					CriticalOdds: 10,
					CriticalLoss: 200,
				},
				Right,
				0,
				WarriorSkills(
					Regular[0],
					Scaffold[2].Fork(nil).(*FatReactor),
				),
			),
		},
	)

	f.Run()
	// Output:
	// signal=9 affairs=1 source.position=0 source.side=Left source.reactor="@Launch({1} 455% Damage, BuffImmune" lifecycle.leading=0 lifecycle.cooling="{Current:3 Maximum:3}" lifecycle.capacity=-1
	// verb=attack critical=false loss=83 overflow=0 source.side=Left source.position=0 source.reactor="@Launch({1} 455% Damage, BuffImmune" source.damage=91 target.side=Right target.position=0 target.defense=8 target.health.current=187 target.health.maximum=270
	// verb=buff reactor=BuffImmune target.side=Right target.position=0 source.reactor="@Launch({1} 455% Damage, BuffImmune"
	// signal=35 affairs=1 source.position=0 source.side=Right source.reactor=#3 lifecycle.leading=0 lifecycle.cooling="{Current:0 Maximum:0}" lifecycle.capacity=0
	// immune.side=Right immune.position=0
	// targets=1 falseTargets=0 immuneTargets=1
	// verb=attack critical=false loss=35 overflow=0 source.side=Right source.position=0 source.reactor=NormalAttack source.damage=40 target.side=Left target.position=0 target.defense=5 target.health.current=65 target.health.maximum=100
	// signal=61 affairs=0 source.position=0 source.side=Left source.reactor="@Launch({1} 455% Damage, BuffImmune" lifecycle.leading=0 lifecycle.cooling="{Current:2 Maximum:3}" lifecycle.capacity=-1
	// signal=61 affairs=0 source.position=0 source.side=Right source.reactor=BuffImmune lifecycle.leading=0 lifecycle.cooling="{Current:0 Maximum:0}" lifecycle.capacity=2
	// signal=75 affairs=1 source.position=0 source.side=Left source.reactor="@Launch({1} 2 * 360% Damage, 15% Stun)" lifecycle.leading=0 lifecycle.cooling="{Current:4 Maximum:4}" lifecycle.capacity=-1
	// verb=attack critical=false loss=64 overflow=0 source.side=Left source.position=0 source.reactor="@Launch({1} 2 * 360% Damage, 15% Stun)" source.damage=72 target.side=Right target.position=0 target.defense=8 target.health.current=123 target.health.maximum=270
	// verb=attack critical=false loss=64 overflow=0 source.side=Left source.position=0 source.reactor="@Launch({1} 2 * 360% Damage, 15% Stun)" source.damage=72 target.side=Right target.position=0 target.defense=8 target.health.current=59 target.health.maximum=270
	// verb=buff reactor=Stun target.side=Right target.position=0 source.reactor="@Launch({1} 2 * 360% Damage, 15% Stun)"
	// signal=116 affairs=0 source.position=0 source.side=Left source.reactor="@Launch({1} 455% Damage, BuffImmune" lifecycle.leading=0 lifecycle.cooling="{Current:1 Maximum:3}" lifecycle.capacity=-1
	// signal=116 affairs=0 source.position=0 source.side=Left source.reactor="@Launch({1} 2 * 360% Damage, 15% Stun)" lifecycle.leading=0 lifecycle.cooling="{Current:3 Maximum:4}" lifecycle.capacity=-1
	// signal=116 affairs=0 source.position=0 source.side=Right source.reactor=Stun lifecycle.leading=0 lifecycle.cooling="{Current:0 Maximum:0}" lifecycle.capacity=0
	// signal=116 affairs=0 source.position=0 source.side=Right source.reactor=BuffImmune lifecycle.leading=0 lifecycle.cooling="{Current:0 Maximum:0}" lifecycle.capacity=1
	// verb=attack critical=false loss=12 overflow=0 source.side=Left source.position=0 source.reactor=NormalAttack source.damage=20 target.side=Right target.position=0 target.defense=8 target.health.current=47 target.health.maximum=270
	// verb=attack critical=false loss=35 overflow=0 source.side=Right source.position=0 source.reactor=NormalAttack source.damage=40 target.side=Left target.position=0 target.defense=5 target.health.current=30 target.health.maximum=100
	// signal=168 affairs=0 source.position=0 source.side=Left source.reactor="@Launch({1} 455% Damage, BuffImmune" lifecycle.leading=0 lifecycle.cooling="{Current:0 Maximum:3}" lifecycle.capacity=-1
	// signal=168 affairs=0 source.position=0 source.side=Left source.reactor="@Launch({1} 2 * 360% Damage, 15% Stun)" lifecycle.leading=0 lifecycle.cooling="{Current:2 Maximum:4}" lifecycle.capacity=-1
	// signal=168 affairs=0 source.position=0 source.side=Right source.reactor=BuffImmune lifecycle.leading=0 lifecycle.cooling="{Current:0 Maximum:0}" lifecycle.capacity=0
	// signal=185 affairs=1 source.position=0 source.side=Left source.reactor="@Launch({1} 455% Damage, BuffImmune" lifecycle.leading=0 lifecycle.cooling="{Current:3 Maximum:3}" lifecycle.capacity=-1
	// verb=attack critical=false loss=83 overflow=36 source.side=Left source.position=0 source.reactor="@Launch({1} 455% Damage, BuffImmune" source.damage=91 target.side=Right target.position=0 target.defense=8 target.health.current=0 target.health.maximum=270
	// targets=0 falseTargets=1 immuneTargets=0
}

func ExampleBattleField_Run_special_8() {
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

	RngX.SetRng(NewSequence(0.5, 0.1))
	f := NewBattleField(
		[]Warrior{
			NewMyWarrior(
				MyBaseline{
					Damage:  20,
					Defense: 5,
					Health:  60,
				},
				Left,
				0,
				WarriorSkills(
					Regular[0],
					Special[8][0],
					Special[8][1],
					Special[8][2],
					Special[8][3],
				),
			),
			NewMyWarrior(
				MyBaseline{
					Damage:  40,
					Defense: 8,
					Health:  270,
				},
				Right,
				0,
				WarriorSkills(
					Regular[0],
					Scaffold[0].Fork(nil).(*FatReactor),
				),
			),
		},
	)

	f.Run()
	// Output:
	// verb=buff reactor="+15% Defense" target.side=Left target.position=0 source.reactor="@BattleStart({$} +15% Defense)"
	// signal=16 affairs=1 source.position=0 source.side=Left source.reactor="@Launch({$} ControlImmune)" lifecycle.leading=0 lifecycle.cooling="{Current:5 Maximum:5}" lifecycle.capacity=-1
	// signal=16 affairs=1 source.position=0 source.side=Left source.reactor="@Launch({1!/%} 340% Damage, {~1!/%} 80% Damage+)" lifecycle.leading=0 lifecycle.cooling="{Current:4 Maximum:4}" lifecycle.capacity=-1
	// verb=buff reactor=ControlImmune target.side=Left target.position=0 source.reactor="@Launch({$} ControlImmune)"
	// verb=attack critical=false loss=60 overflow=0 source.side=Left source.position=0 source.reactor="@Launch({1!/%} 340% Damage, {~1!/%} 80% Damage+)" source.damage=68 target.side=Right target.position=0 target.defense=8 target.health.current=210 target.health.maximum=270
	// verb=heal rise=48 overflow=48 source.side=Left source.position=0 source.reactor="@Launch({1!/%} 340% Damage, {~1!/%} 80% Damage+)" target.side=Left target.position=0 target.health.current=60 target.health.maximum=60
	// signal=54 affairs=1 source.position=0 source.side=Right source.reactor=#1 lifecycle.leading=0 lifecycle.cooling="{Current:0 Maximum:0}" lifecycle.capacity=0
	// immune.side=Left immune.position=0
	// targets=1 falseTargets=0 immuneTargets=1
	// signal=67 affairs=0 source.position=0 source.side=Left source.reactor="@Launch({$} ControlImmune)" lifecycle.leading=0 lifecycle.cooling="{Current:4 Maximum:5}" lifecycle.capacity=-1
	// signal=67 affairs=0 source.position=0 source.side=Left source.reactor="@Launch({1!/%} 340% Damage, {~1!/%} 80% Damage+)" lifecycle.leading=0 lifecycle.cooling="{Current:3 Maximum:4}" lifecycle.capacity=-1
	// signal=67 affairs=0 source.position=0 source.side=Left source.reactor=ControlImmune lifecycle.leading=0 lifecycle.cooling="{Current:0 Maximum:0}" lifecycle.capacity=1
	// signal=84 affairs=1 source.position=0 source.side=Left source.reactor="@Launch({*} +120% Shield, 7% Regeneration)" lifecycle.leading=0 lifecycle.cooling="{Current:4 Maximum:4}" lifecycle.capacity=-1
	// capacity=24 verb=buff reactor=Shield target.side=Left target.position=0 source.reactor="@Launch({*} +120% Shield, 7% Regeneration)"
	// verb=buff reactor=Regeneration target.side=Left target.position=0 source.reactor="@Launch({*} +120% Shield, 7% Regeneration)"
	// signal=115 affairs=1 source.position=0 source.side=Left source.reactor=Shield lifecycle.leading=0 lifecycle.cooling="{Current:0 Maximum:0}" lifecycle.capacity=0
	// verb=attack critical=false loss=11 overflow=0 source.side=Right source.position=0 source.reactor=NormalAttack source.damage=40 target.side=Left target.position=0 target.defense=5 target.health.current=49 target.health.maximum=60
	// signal=125 affairs=0 source.position=0 source.side=Left source.reactor="@Launch({$} ControlImmune)" lifecycle.leading=0 lifecycle.cooling="{Current:3 Maximum:5}" lifecycle.capacity=-1
	// signal=125 affairs=0 source.position=0 source.side=Left source.reactor="@Launch({1!/%} 340% Damage, {~1!/%} 80% Damage+)" lifecycle.leading=0 lifecycle.cooling="{Current:2 Maximum:4}" lifecycle.capacity=-1
	// signal=125 affairs=0 source.position=0 source.side=Left source.reactor="@Launch({*} +120% Shield, 7% Regeneration)" lifecycle.leading=0 lifecycle.cooling="{Current:3 Maximum:4}" lifecycle.capacity=-1
	// signal=125 affairs=0 source.position=0 source.side=Left source.reactor=ControlImmune lifecycle.leading=0 lifecycle.cooling="{Current:0 Maximum:0}" lifecycle.capacity=0
	// signal=140 affairs=1 source.position=0 source.side=Left source.reactor=Regeneration lifecycle.leading=0 lifecycle.cooling="{Current:0 Maximum:0}" lifecycle.capacity=0
	// verb=heal rise=4 overflow=0 source.side=Left source.position=0 source.reactor=Regeneration target.side=Left target.position=0 target.health.current=53 target.health.maximum=60
	// verb=attack critical=false loss=12 overflow=0 source.side=Left source.position=0 source.reactor=NormalAttack source.damage=20 target.side=Right target.position=0 target.defense=8 target.health.current=198 target.health.maximum=270
	// verb=attack critical=false loss=35 overflow=0 source.side=Right source.position=0 source.reactor=NormalAttack source.damage=40 target.side=Left target.position=0 target.defense=5 target.health.current=18 target.health.maximum=60
	// signal=189 affairs=0 source.position=0 source.side=Left source.reactor="@Launch({$} ControlImmune)" lifecycle.leading=0 lifecycle.cooling="{Current:2 Maximum:5}" lifecycle.capacity=-1
	// signal=189 affairs=0 source.position=0 source.side=Left source.reactor="@Launch({1!/%} 340% Damage, {~1!/%} 80% Damage+)" lifecycle.leading=0 lifecycle.cooling="{Current:1 Maximum:4}" lifecycle.capacity=-1
	// signal=189 affairs=0 source.position=0 source.side=Left source.reactor="@Launch({*} +120% Shield, 7% Regeneration)" lifecycle.leading=0 lifecycle.cooling="{Current:2 Maximum:4}" lifecycle.capacity=-1
	// verb=attack critical=false loss=12 overflow=0 source.side=Left source.position=0 source.reactor=NormalAttack source.damage=20 target.side=Right target.position=0 target.defense=8 target.health.current=186 target.health.maximum=270
	// verb=attack critical=false loss=35 overflow=17 source.side=Right source.position=0 source.reactor=NormalAttack source.damage=40 target.side=Left target.position=0 target.defense=5 target.health.current=0 target.health.maximum=60
}

func ExampleBattleField_Run_special_9() {
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

	RngX.SetRng(NewSequence(0.5, 0.1))
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
				WarriorSkills(
					Regular[0],
					Special[9][0],
					Special[9][1],
					Special[9][2],
				),
			),
			NewMyWarrior(
				MyBaseline{
					Damage:  20,
					Defense: 8,
					Health:  200,
				},
				Right,
				0,
				WarriorSkills(
					Regular[0],
					Scaffold[0].Fork(nil).(*FatReactor),
					Scaffold[1].Fork(nil).(*FatReactor),
				),
			),
		},
	)

	f.Run()
	// Output:
	// signal=9 affairs=0 source.position=0 source.side=Left source.reactor="@Launch[..4]({$} Remove 2 Nerf)" lifecycle.leading=3 lifecycle.cooling="{Current:0 Maximum:0}" lifecycle.capacity=-1
	// signal=9 affairs=1 source.position=0 source.side=Left source.reactor="@Launch({1!/%} 3 * 380% Damage)" lifecycle.leading=0 lifecycle.cooling="{Current:4 Maximum:4}" lifecycle.capacity=-1
	// verb=attack critical=false loss=30 overflow=0 source.side=Left source.position=0 source.reactor="@Launch({1!/%} 3 * 380% Damage)" source.damage=38 target.side=Right target.position=0 target.defense=8 target.health.current=170 target.health.maximum=200
	// verb=attack critical=false loss=30 overflow=0 source.side=Left source.position=0 source.reactor="@Launch({1!/%} 3 * 380% Damage)" source.damage=38 target.side=Right target.position=0 target.defense=8 target.health.current=140 target.health.maximum=200
	// verb=attack critical=false loss=30 overflow=0 source.side=Left source.position=0 source.reactor="@Launch({1!/%} 3 * 380% Damage)" source.damage=38 target.side=Right target.position=0 target.defense=8 target.health.current=110 target.health.maximum=200
	// signal=55 affairs=1 source.position=0 source.side=Right source.reactor=#2 lifecycle.leading=0 lifecycle.cooling="{Current:0 Maximum:0}" lifecycle.capacity=0
	// verb=buff reactor="Nerf #2" target.side=Left target.position=0 source.reactor=#2
	// signal=68 affairs=0 source.position=0 source.side=Left source.reactor="@Launch({1!/%} 3 * 380% Damage)" lifecycle.leading=0 lifecycle.cooling="{Current:3 Maximum:4}" lifecycle.capacity=-1
	// signal=79 affairs=0 source.position=0 source.side=Left source.reactor="@Launch[..4]({$} Remove 2 Nerf)" lifecycle.leading=2 lifecycle.cooling="{Current:0 Maximum:0}" lifecycle.capacity=-1
	// signal=79 affairs=1 source.position=0 source.side=Left source.reactor="@Launch({1} 3 * 350% Damage, -50% Defense)" lifecycle.leading=0 lifecycle.cooling="{Current:4 Maximum:4}" lifecycle.capacity=-1
	// verb=attack critical=false loss=27 overflow=0 source.side=Left source.position=0 source.reactor="@Launch({1} 3 * 350% Damage, -50% Defense)" source.damage=35 target.side=Right target.position=0 target.defense=8 target.health.current=83 target.health.maximum=200
	// verb=attack critical=false loss=27 overflow=0 source.side=Left source.position=0 source.reactor="@Launch({1} 3 * 350% Damage, -50% Defense)" source.damage=35 target.side=Right target.position=0 target.defense=8 target.health.current=56 target.health.maximum=200
	// verb=attack critical=false loss=27 overflow=0 source.side=Left source.position=0 source.reactor="@Launch({1} 3 * 350% Damage, -50% Defense)" source.damage=35 target.side=Right target.position=0 target.defense=8 target.health.current=29 target.health.maximum=200
	// verb=buff reactor="-50% Defense" target.side=Right target.position=0 source.reactor="@Launch({1} 3 * 350% Damage, -50% Defense)"
	// signal=132 affairs=1 source.position=0 source.side=Right source.reactor=#1 lifecycle.leading=0 lifecycle.cooling="{Current:0 Maximum:0}" lifecycle.capacity=0
	// verb=buff reactor="Nerf #1" target.side=Left target.position=0 source.reactor=#1
	// signal=145 affairs=0 source.position=0 source.side=Left source.reactor="@Launch({1!/%} 3 * 380% Damage)" lifecycle.leading=0 lifecycle.cooling="{Current:2 Maximum:4}" lifecycle.capacity=-1
	// signal=145 affairs=0 source.position=0 source.side=Left source.reactor="@Launch({1} 3 * 350% Damage, -50% Defense)" lifecycle.leading=0 lifecycle.cooling="{Current:3 Maximum:4}" lifecycle.capacity=-1
	// signal=145 affairs=0 source.position=0 source.side=Right source.reactor="-50% Defense" lifecycle.leading=0 lifecycle.cooling="{Current:0 Maximum:0}" lifecycle.capacity=0
	// signal=162 affairs=0 source.position=0 source.side=Left source.reactor="@Launch[..4]({$} Remove 2 Nerf)" lifecycle.leading=1 lifecycle.cooling="{Current:0 Maximum:0}" lifecycle.capacity=-1
	// verb=attack critical=false loss=2 overflow=0 source.side=Left source.position=0 source.reactor=NormalAttack source.damage=10 target.side=Right target.position=0 target.defense=8 target.health.current=27 target.health.maximum=200
	// verb=attack critical=false loss=15 overflow=0 source.side=Right source.position=0 source.reactor=NormalAttack source.damage=20 target.side=Left target.position=0 target.defense=5 target.health.current=85 target.health.maximum=100
	// signal=197 affairs=0 source.position=0 source.side=Left source.reactor="@Launch({1!/%} 3 * 380% Damage)" lifecycle.leading=0 lifecycle.cooling="{Current:1 Maximum:4}" lifecycle.capacity=-1
	// signal=197 affairs=0 source.position=0 source.side=Left source.reactor="@Launch({1} 3 * 350% Damage, -50% Defense)" lifecycle.leading=0 lifecycle.cooling="{Current:2 Maximum:4}" lifecycle.capacity=-1
	// signal=211 affairs=1 source.position=0 source.side=Left source.reactor="@Launch[..4]({$} Remove 2 Nerf)" lifecycle.leading=0 lifecycle.cooling="{Current:0 Maximum:0}" lifecycle.capacity=0
	// verb=purge reactors="[Nerf #2 Nerf #1]"
	// verb=attack critical=false loss=2 overflow=0 source.side=Left source.position=0 source.reactor=NormalAttack source.damage=10 target.side=Right target.position=0 target.defense=8 target.health.current=25 target.health.maximum=200
	// verb=attack critical=false loss=15 overflow=0 source.side=Right source.position=0 source.reactor=NormalAttack source.damage=20 target.side=Left target.position=0 target.defense=5 target.health.current=70 target.health.maximum=100
	// signal=253 affairs=0 source.position=0 source.side=Left source.reactor="@Launch({1!/%} 3 * 380% Damage)" lifecycle.leading=0 lifecycle.cooling="{Current:0 Maximum:4}" lifecycle.capacity=-1
	// signal=253 affairs=0 source.position=0 source.side=Left source.reactor="@Launch({1} 3 * 350% Damage, -50% Defense)" lifecycle.leading=0 lifecycle.cooling="{Current:1 Maximum:4}" lifecycle.capacity=-1
	// signal=267 affairs=1 source.position=0 source.side=Left source.reactor="@Launch({1!/%} 3 * 380% Damage)" lifecycle.leading=0 lifecycle.cooling="{Current:4 Maximum:4}" lifecycle.capacity=-1
	// verb=attack critical=false loss=30 overflow=5 source.side=Left source.position=0 source.reactor="@Launch({1!/%} 3 * 380% Damage)" source.damage=38 target.side=Right target.position=0 target.defense=8 target.health.current=0 target.health.maximum=200
	// targets=0 falseTargets=1 immuneTargets=0
	// targets=0 falseTargets=1 immuneTargets=0
}
