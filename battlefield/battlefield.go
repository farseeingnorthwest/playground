package battlefield

import (
	"sort"
)

type side int

const (
	Left side = iota
	Right
)

type Preparer interface {
	Prepare()
}

type Attacker interface {
	Attack() (damage int, critical bool)
	Velocity() int
}

type Sufferer interface {
	Suffer(attack int) (damage, overflow int)
	Health() int
}

type Warrior interface {
	Preparer
	Attacker
	Sufferer
}

type Attack struct {
	Attacker Warrior
	Sufferer Warrior
	Attack   int
	Critical bool
	Damage   int
	Overflow int
}

type fighter struct {
	Warrior
	seat int
	side
}

type byVelocity []fighter

func newByVelocity(warriors []Warrior, side side) byVelocity {
	fighters := make([]fighter, len(warriors))
	for i, warrior := range warriors {
		fighters[i] = fighter{warrior, i, side}
	}

	return fighters
}

func (fighters byVelocity) Len() int {
	return len(fighters)
}

func (fighters byVelocity) Less(i, j int) bool {
	return fighters[i].Velocity() > fighters[j].Velocity() ||
		(fighters[i].Velocity() == fighters[j].Velocity() &&
			fighters[i].side < fighters[j].side) ||
		(fighters[i].Velocity() == fighters[j].Velocity() &&
			fighters[i].side == fighters[j].side &&
			fighters[i].seat < fighters[j].seat)
}

func (fighters byVelocity) Swap(i, j int) {
	fighters[i], fighters[j] = fighters[j], fighters[i]
}

type Observer interface {
	Observe(attack Attack)
}

func Fight(a, b []Warrior, observer Observer) {
	fighters := append(newByVelocity(a, Left), newByVelocity(b, Right)...)
	sort.Sort(fighters)

	alive := []int{
		countIf(a, isAlive[Warrior]),
		countIf(b, isAlive[Warrior]),
	}
	for alive[Left] > 0 && alive[Right] > 0 {
		for _, fighter := range fighters {
			if fighter.Health() > 0 {
				fighter.Prepare()
			}
		}

		for _, attacker := range fighters {
			if attacker.Health() == 0 {
				continue
			}

			for _, sufferer := range fighters {
				if sufferer.Health() == 0 {
					continue
				}

				if attacker.side == sufferer.side {
					continue
				}

				attack, critical := attacker.Attack()
				damage, overflow := sufferer.Suffer(attack)
				if sufferer.Health() <= 0 {
					alive[sufferer.side]--
				}

				observer.Observe(Attack{
					attacker.Warrior,
					sufferer.Warrior,
					attack,
					critical,
					damage,
					overflow,
				})
				break
			}
		}
	}

	return
}

func countIf[T any](slice []T, predicate func(T) bool) (count int) {
	for _, element := range slice {
		if predicate(element) {
			count++
		}
	}

	return
}

func isAlive[T Sufferer](sufferer T) bool {
	return sufferer.Health() > 0
}
