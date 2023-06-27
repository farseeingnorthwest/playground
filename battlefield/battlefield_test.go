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
	observer := &observer{}
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

type observer struct {
	attacks []string
}

func (o *observer) Observe(attack Attack) {
	o.attacks = append(o.attacks, fmt.Sprintf("%v => %v 🗡️(%v/%v/%v) %v",
		attack.Attacker.(*warrior).name,
		attack.Sufferer.(*warrior).name,
		attack.Attack,
		attack.Damage,
		attack.Overflow,
		attack.Sufferer.(Warrior).Health(),
	))
}
