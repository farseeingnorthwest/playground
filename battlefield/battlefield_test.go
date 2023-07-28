package battlefield

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFighterList_Drain(t *testing.T) {
	for i, tt := range []struct {
		fighters []Fighter
		expected []string
	}{
		{
			[]Fighter{
				newWarrior("A", 10, 10, 0, 10, false, false),
				newWarrior("B", 10, 10, 10, 10, false, false),
				newWarrior("C", 10, 10, 10, 10, false, false),
			},
			[]string{"C", "B"},
		},
		{
			[]Fighter{
				newWarrior("A", 10, 10, 10, 10, false, false),
				newWarrior("B", 10, 10, 10, 10, false, false),
				newWarrior("C", 10, 10, 0, 10, false, false),
			},
			[]string{"A", "B"},
		},
		{
			[]Fighter{
				newWarrior("A", 10, 10, 0, 10, false, false),
				newWarrior("B", 10, 10, 10, 10, false, false),
				newWarrior("C", 10, 10, 0, 10, false, false),
			},
			[]string{"B"},
		},
	} {
		t.Run(fmt.Sprintf("#%v", i), func(t *testing.T) {
			fl := newFighterList(tt.fighters, Right)
			fl = fl.Drain()

			assert.Equal(t, len(tt.expected), len(fl))
			for i, f := range fl {
				assert.Equal(t, tt.expected[i], f.Fighter.(*warrior).name)
			}
		})
	}
}

func TestFight(t *testing.T) {
	for i, tt := range []struct {
		randomizer Randomizer
		a, b       []Fighter
		expected   []string
	}{
		{
			NewSequence(0.5, 0.8, 0.8, 0.5, 0.2, 0.9, 0.3, 0.5, 0.8, 0.1, 0.8, 0.6),
			[]Fighter{
				newWarrior("A", 15, 5, 55, 5, false, false),
				newWarrior("B", 10, 8, 30, 9, false, false),
				newWarrior("C", 12, 9, 45, 9, true, false),
			},
			[]Fighter{
				newWarrior("D", 18, 12, 50, 9, true, false),
				newWarrior("E", 10, 5, 20, 10, false, false),
			},
			[]string{
				"E => B 🗡️(20/12/0) 18",
				"B => E 🗡️(20/15/0) 5",
				"C => E 🗡️(12/7/2) 0",
				"D => B 🗡️(18/10/0) 8",
				"A => D 🗡️(30/18/0) 32",
				"B => D 🗡️(10/0/0) 32",
				"C => D 🗡️(24/12/0) 20",
				"D => B 🗡️(36/28/20) 0",
				"A => D 🗡️(15/3/0) 17",
				"C => D 🗡️(12/0/0) 17",
				"D => C 🗡️(18/9/0) 36",
				"A => D 🗡️(30/18/1) 0",
			},
		},
		{
			NewSequence(0.1),
			[]Fighter{
				newWarrior("A", 15, 5, 55, 5, false, false),
				newWarrior("B", 10, 8, 30, 9, false, false),
				newWarrior("C", 12, 9, 45, 9, true, false),
			},
			[]Fighter{
				newWarrior("D", 18, 12, 50, 9, true, false),
				newWarrior("E", 10, 5, 20, 10, false, false),
			},
			[]string{
				"E => A 🗡️(20/15/0) 40",
				"B => D 🗡️(20/8/0) 42",
				"C => D 🗡️(12/0/0) 42",
				"D => A 🗡️(18/13/0) 27",
				"A => D 🗡️(30/18/0) 24",
				"E => A 🗡️(10/5/0) 22",
				"B => D 🗡️(10/0/0) 24",
				"C => D 🗡️(24/12/0) 12",
				"D => A 🗡️(36/31/9) 0",
				"E => C 🗡️(20/11/0) 34",
				"B => D 🗡️(20/8/0) 4",
				"C => D 🗡️(12/0/0) 4",
				"D => C 🗡️(18/9/0) 25",
				"E => C 🗡️(10/1/0) 24",
				"B => D 🗡️(10/0/0) 4",
				"C => D 🗡️(24/12/8) 0",
				"E => C 🗡️(20/11/0) 13",
				"B => E 🗡️(20/15/0) 5",
				"C => E 🗡️(12/7/2) 0",
			},
		},
		{
			NewSequence(0.1),
			[]Fighter{
				newWarrior("A", 15, 5, 55, 5, false, false),
				newWarrior("B", 10, 8, 30, 9, false, false),
				newWarrior("C", 12, 9, 45, 8, true, true),
			},
			[]Fighter{
				newWarrior("D", 18, 12, 50, 9, true, false),
				newWarrior("E", 10, 5, 20, 10, false, false),
			},
			[]string{
				"E => A 🗡️(20/15/0) 40",
				"B => D 🗡️(20/8/0) 42",
				"C => D 🗡️(12/0/0) 42",
				"D => A 🗡️(18/13/0) 27",
				"A => D 🗡️(30/18/0) 24",
				"C => D 🗡️(24/12/0) 12",
				"E => A 🗡️(10/5/0) 22",
				"B => D 🗡️(10/0/0) 12",
				"D => A 🗡️(36/31/9) 0",
				"C => D 🗡️(12/0/0) 12",
				"E => C 🗡️(20/11/0) 34",
				"B => D 🗡️(20/8/0) 4",
				"D => C 🗡️(18/9/0) 25",
				"C => D 🗡️(24/12/8) 0",
				"E => C 🗡️(10/1/0) 24",
				"B => E 🗡️(10/5/0) 15",
				"C => E 🗡️(12/7/0) 8",
				"E => C 🗡️(20/11/0) 13",
				"B => E 🗡️(20/15/7) 0",
			},
		},
	} {
		t.Run(fmt.Sprintf("#%v", i), func(t *testing.T) {
			observer := &journalObserver{}
			f := NewBattlefield(tt.randomizer)
			f.Fight(tt.a, tt.b, observer)

			assert.Equal(t, tt.expected, observer.attacks)
		})
	}
}

func TestFight_DeathSetup(t *testing.T) {
	warriors := []Fighter{
		newWarrior("A", 15, 5, 55, 5, false, false),
		newWarrior("B", 10, 8, 0, 9, false, false),
		newWarrior("C", 12, 9, 36, 8, true, true),
		newWarrior("D", 18, 12, 0, 9, true, false),
		newWarrior("E", 10, 5, 0, 10, false, false),
	}
	observer := &journalObserver{}
	f := NewBattlefield(DefaultRandomizer{})
	f.Fight(warriors[:3], warriors[3:], observer)

	assert.Empty(t, observer.attacks)
}

func BenchmarkFight(b *testing.B) {
	for i := 0; i < b.N; i++ {
		warriors := []Fighter{
			newWarrior("A", 15, 5, 55, 5, false, false),
			newWarrior("B", 10, 8, 30, 9, false, false),
			newWarrior("C", 12, 9, 45, 8, true, true),
			newWarrior("D", 18, 12, 50, 9, true, false),
			newWarrior("E", 10, 5, 20, 10, false, false),
		}
		observer := &dummyObserver{}
		f := NewBattlefield(DefaultRandomizer{})

		f.Fight(warriors[:3], warriors[3:], observer)
	}
}

func FuzzFight(f *testing.F) {
	f.Add(
		15, 5, 55, 5, false, false,
		10, 8, 30, 9, false, false,
		12, 9, 45, 8, true, true,
		18, 12, 50, 9, true, false,
		10, 5, 20, 10, false, false,
	)
	f.Fuzz(func(
		t *testing.T,
		aAttack, aDefense, aHealth, aSpeed int, aCritical bool, aSpeedUp bool,
		bAttack, bDefense, bHealth, bSpeed int, bCritical bool, bSpeedUp bool,
		cAttack, cDefense, cHealth, cSpeed int, cCritical bool, cSpeedUp bool,
		dAttack, dDefense, dHealth, dSpeed int, dCritical bool, dSpeedUp bool,
		eAttack, eDefense, eHealth, eSpeed int, eCritical bool, eSpeedUp bool,
	) {
		warriors := []Fighter{
			newFuzzWarrior("A", aAttack, aDefense, aHealth, aSpeed, aCritical, aSpeedUp),
			newFuzzWarrior("B", bAttack, bDefense, bHealth, bSpeed, bCritical, bSpeedUp),
			newFuzzWarrior("C", cAttack, cDefense, cHealth, cSpeed, cCritical, cSpeedUp),
			newFuzzWarrior("D", dAttack, dDefense, dHealth, dSpeed, dCritical, dSpeedUp),
			newFuzzWarrior("E", eAttack, eDefense, eHealth, eSpeed, eCritical, eSpeedUp),
		}
		observer := &observer{}
		f := NewBattlefield(DefaultRandomizer{})

		f.Fight(warriors[:3], warriors[3:], observer)
		var deaths []string
		for _, h := range observer.attacks {
			assert.NotContains(t, deaths, h.Attacker.(*warrior).name)
			assert.NotContains(t, deaths, h.Sufferer.(*warrior).name)
			assert.GreaterOrEqual(t, h.health, 0)
			if h.health > 0 {
				assert.Equal(t, h.Overflow, 0)
				continue
			}

			deaths = append(deaths, h.Sufferer.(*warrior).name)
			assert.GreaterOrEqual(t, h.Overflow, 0)
		}
	})
}

type warrior struct {
	name     string
	attack   int
	defense  int
	health   int
	speed    int
	critical bool
	speedUp  bool
}

func newWarrior(name string, attack, defense, health, speed int, critical, speedUp bool) *warrior {
	return &warrior{name, attack, defense, health, speed, critical, speedUp}
}

func newFuzzWarrior(name string, attack, defense, health, speed int, critical, speedUp bool) *warrior {
	attack, defense, health, speed = abs(attack), abs(defense), abs(health), abs(speed)
	if attack < 3 {
		attack = 3
	}
	if attack <= defense {
		defense = attack - 1
	}

	return &warrior{name, attack, defense, health, speed, critical, speedUp}
}

func abs(i int) int {
	if i < 0 {
		return -i
	}
	return i
}

func (w *warrior) Prepare(Randomizer) {
	w.critical = !w.critical
	if w.speedUp {
		w.speed++
	}
}

func (w *warrior) Attack() (int, bool) {
	if w.critical {
		return w.attack * 2, w.critical
	}

	return w.attack, w.critical
}

func (w *warrior) Speed() int {
	return w.speed
}

func (w *warrior) Suffer(attack int, _ bool) (damage, overflow int) {
	damage = attack - w.defense
	if damage < 0 {
		damage = 0
	}
	w.health -= damage
	if w.health < 0 {
		overflow = -w.health
		w.health = 0
	}
	return
}

func (w *warrior) Health() int {
	return w.health
}

type journalObserver struct {
	attacks []string
}

func (o *journalObserver) Observe(attack Action) {
	o.attacks = append(o.attacks, fmt.Sprintf(
		"%v => %v 🗡️(%v/%v/%v) %v",
		attack.Attacker.(*warrior).name,
		attack.Sufferer.(*warrior).name,
		attack.Attack,
		attack.Damage,
		attack.Overflow,
		attack.Sufferer.(Fighter).Health(),
	))
}

type dummyObserver struct {
}

func (d *dummyObserver) Observe(Action) {
}

type withHealth struct {
	Action
	health int
}

type observer struct {
	attacks []withHealth
}

func (o *observer) Observe(attack Action) {
	o.attacks = append(o.attacks, withHealth{attack, attack.Sufferer.Health()})
}
