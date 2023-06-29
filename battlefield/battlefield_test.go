package battlefield

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFight(t *testing.T) {
	warriors := []Warrior{
		newWarrior("A", 15, 5, 55, 5, false),
		newWarrior("B", 10, 8, 30, 9, false),
		newWarrior("C", 12, 9, 45, 9, true),
		newWarrior("D", 18, 12, 50, 9, true),
		newWarrior("E", 10, 5, 20, 10, false),
	}
	observer := &journalObserver{}
	Fight(warriors[:3], warriors[3:], observer)

	expected := []string{
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
	}
	assert.Equal(t, expected, observer.attacks)
}

func TestFight_DeathSetup(t *testing.T) {
	warriors := []Warrior{
		newWarrior("A", 15, 5, 55, 5, false),
		newWarrior("B", 10, 8, 0, 9, false),
		newWarrior("C", 12, 9, 36, 9, true),
		newWarrior("D", 18, 12, 0, 9, true),
		newWarrior("E", 10, 5, 0, 10, false),
	}
	observer := &journalObserver{}
	Fight(warriors[:3], warriors[3:], observer)

	assert.Empty(t, observer.attacks)
}

func BenchmarkFight(b *testing.B) {
	for i := 0; i < b.N; i++ {
		warriors := []Warrior{
			newWarrior("A", 15, 5, 55, 5, false),
			newWarrior("B", 10, 8, 30, 9, false),
			newWarrior("C", 12, 9, 45, 9, true),
			newWarrior("D", 18, 12, 50, 9, true),
			newWarrior("E", 10, 5, 20, 10, false),
		}
		observer := &dummyObserver{}

		Fight(warriors[:3], warriors[3:], observer)
	}
}

func FuzzFight(f *testing.F) {
	f.Add(
		15, 5, 55, 5, false,
		10, 8, 30, 9, false,
		12, 9, 45, 9, true,
		18, 12, 50, 9, true,
		10, 5, 20, 10, false,
	)
	f.Fuzz(func(
		t *testing.T,
		aAttack, aDefense, aHealth, aVelocity int, aCritical bool,
		bAttack, bDefense, bHealth, bVelocity int, bCritical bool,
		cAttack, cDefense, cHealth, cVelocity int, cCritical bool,
		dAttack, dDefense, dHealth, dVelocity int, dCritical bool,
		eAttack, eDefense, eHealth, eVelocity int, eCritical bool,
	) {
		warriors := []Warrior{
			newFuzzWarrior("A", aAttack, aDefense, aHealth, aVelocity, aCritical),
			newFuzzWarrior("B", bAttack, bDefense, bHealth, bVelocity, bCritical),
			newFuzzWarrior("C", cAttack, cDefense, cHealth, cVelocity, cCritical),
			newFuzzWarrior("D", dAttack, dDefense, dHealth, dVelocity, dCritical),
			newFuzzWarrior("E", eAttack, eDefense, eHealth, eVelocity, eCritical),
		}
		observer := &observer{}

		Fight(warriors[:3], warriors[3:], observer)
		var deaths []string
		for _, h := range observer.attacks {
			assert.NotContains(t, deaths, h.Attacker.(*warrior).name)
			assert.NotContains(t, deaths, h.Sufferer.(*warrior).name)
			assert.GreaterOrEqual(t, h.health, 0)
			if h.health > 0 {
				assert.Equal(t, h.Overflow, 0)
			} else {
				deaths = append(deaths, h.Sufferer.(*warrior).name)
				assert.GreaterOrEqual(t, h.Overflow, 0)
			}
		}
	})
}

type warrior struct {
	name     string
	attack   int
	defense  int
	health   int
	velocity int
	critical bool
}

func newWarrior(name string, attack, defense, health, velocity int, critical bool) *warrior {
	return &warrior{name, attack, defense, health, velocity, critical}
}

func newFuzzWarrior(name string, attack, defense, health, velocity int, critical bool) *warrior {
	attack, defense, health, velocity = abs(attack), abs(defense), abs(health), abs(velocity)
	if attack < 3 {
		attack = 3
	}
	if attack <= defense {
		defense = attack - 1
	}

	return &warrior{name, attack, defense, health, velocity, critical}
}

func abs(i int) int {
	if i < 0 {
		return -i
	}
	return i
}

func (w *warrior) Prepare() {
	w.critical = !w.critical
}

func (w *warrior) Attack() (int, bool) {
	if w.critical {
		return w.attack * 2, w.critical
	}

	return w.attack, w.critical
}

func (w *warrior) Velocity() int {
	return w.velocity
}

func (w *warrior) Suffer(attack int) (damage, overflow int) {
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

func (o *journalObserver) Observe(attack Attack) {
	o.attacks = append(o.attacks, fmt.Sprintf(
		"%v => %v 🗡️(%v/%v/%v) %v",
		attack.Attacker.(*warrior).name,
		attack.Sufferer.(*warrior).name,
		attack.Attack,
		attack.Damage,
		attack.Overflow,
		attack.Sufferer.(Warrior).Health(),
	))
}

type dummyObserver struct {
}

func (d *dummyObserver) Observe(Attack) {
}

type withHealth struct {
	Attack
	health int
}

type observer struct {
	attacks []withHealth
}

func (o *observer) Observe(attack Attack) {
	o.attacks = append(o.attacks, withHealth{attack, attack.Sufferer.Health()})
}
