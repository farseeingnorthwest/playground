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
	Prepare(Randomizer)
}

type Attacker interface {
	Attack() (attack int, critical bool)
	Velocity() int
}

type Sufferer interface {
	Suffer(attack int, critical bool) (damage, overflow int)
	Health() int
}

type Fighter interface {
	Preparer
	Attacker
	Sufferer
}

type Action struct {
	Attacker Fighter
	Sufferer Fighter
	Attack   int
	Critical bool
	Damage   int
	Overflow int
}

type fighter struct {
	Fighter
	seat int
	side
}

type byVelocity []fighter

func newByVelocity(warriors []Fighter, side side) byVelocity {
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
	Observe(attack Action)
}

func Fight(a, b []Fighter, observer Observer, randomizer Randomizer) {
	fighters := append(newByVelocity(a, Left), newByVelocity(b, Right)...)
	sort.Sort(fighters)

	alive := []int{
		countIf(a, isAlive[Fighter]),
		countIf(b, isAlive[Fighter]),
	}
	for alive[Left] > 0 && alive[Right] > 0 {
		for _, fighter := range fighters {
			if fighter.Health() > 0 {
				fighter.Prepare(randomizer)
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
				damage, overflow := sufferer.Suffer(attack, critical)
				if sufferer.Health() <= 0 {
					alive[sufferer.side]--
				}

				observer.Observe(Action{
					attacker.Fighter,
					sufferer.Fighter,
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
