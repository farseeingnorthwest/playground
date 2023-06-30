package battlefield

import "sort"

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
	Speed() int
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

type fighterInPosition struct {
	Fighter
	seat int
	side
}

type fighterList []fighterInPosition

func newFighterList(fighters []Fighter, side side) fighterList {
	fips := make([]fighterInPosition, len(fighters))
	for i, f := range fighters {
		fips[i] = fighterInPosition{f, i, side}
	}

	return fips
}

func (f fighterList) Drain() fighterList {
	i, n := 0, len(f)
	for i < n {
		if f[i].Health() > 0 {
			i++
			continue
		}

		n--
		f[i] = f[n]
	}

	return f[:n]
}

type bySpeed []fighterInPosition

func (f bySpeed) Len() int {
	return len(f)
}

func (f bySpeed) Less(i, j int) bool {
	return f[i].Speed() > f[j].Speed() ||
		(f[i].Speed() == f[j].Speed() &&
			f[i].side < f[j].side) ||
		(f[i].Speed() == f[j].Speed() &&
			f[i].side == f[j].side &&
			f[i].seat < f[j].seat)
}

func (f bySpeed) Swap(i, j int) {
	f[i], f[j] = f[j], f[i]
}

type Observer interface {
	Observe(attack Action)
}

type Battlefield struct {
	randomizer Randomizer
}

func NewBattlefield(randomizer Randomizer) *Battlefield {
	return &Battlefield{randomizer}
}

func (f *Battlefield) Fight(a, b []Fighter, observer Observer) {
	sides := []fighterList{
		newFighterList(a, Left).Drain(),
		newFighterList(b, Right).Drain(),
	}
	if len(sides[Left]) == 0 || len(sides[Right]) == 0 {
		return
	}

	var fighters fighterList
	for _, side := range sides {
		fighters = append(fighters, side...)
	}
	for {
		for _, fighter := range fighters {
			fighter.Prepare(f.randomizer)
		}

		sort.Sort(bySpeed(fighters))
		for _, attacker := range fighters {
			if attacker.Health() == 0 {
				continue
			}

			sufferers := sides[Left+Right-attacker.side]
			sufferer := sufferers[int(f.randomizer.Float64()*float64(len(sufferers)))]
			attack, critical := attacker.Attack()
			damage, overflow := sufferer.Suffer(attack, critical)
			observer.Observe(Action{
				attacker.Fighter,
				sufferer.Fighter,
				attack,
				critical,
				damage,
				overflow,
			})

			if sufferer.Health() <= 0 {
				sides[sufferer.side] = sides[sufferer.side].Drain()
				if len(sides[sufferer.side]) == 0 {
					return
				}
			}
		}

		fighters = fighters.Drain()
	}
}
