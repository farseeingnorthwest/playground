package battlefield

import (
	"sort"
)

type Side uint8

const (
	Left Side = iota
	Right
)

type Signal uint8

const (
	Launch Signal = iota
)

type Attack struct {
	Objects []*Fighter
	Points  int
}

type Script struct {
	Subject *Fighter
	Attacks []Attack
}

type Warrior interface {
	Reactor

	Speed() int

	Functional() bool
	Render(Attack)
}

type Observer interface {
	Observe(*Script)
}

type Fighter struct {
	Warrior
	Side
	Position uint8
}

type bySpeed []*Fighter

func (f bySpeed) Len() int { return len(f) }
func (f bySpeed) Less(i, j int) bool {
	if f[i].Speed() != f[j].Speed() {
		return f[i].Speed() > f[j].Speed()
	}
	if f[i].Side != f[j].Side {
		return f[i].Side == Left
	}

	return f[i].Position < f[j].Position
}
func (f bySpeed) Swap(i, j int) { f[i], f[j] = f[j], f[i] }

func Fight(a, b []Warrior, observer Observer) {
	fighters := make([]*Fighter, len(a)+len(b))
	for i, f := range a {
		fighters[i] = &Fighter{f, Left, uint8(i)}
	}
	for i, f := range b {
		fighters[i+len(a)] = &Fighter{f, Right, uint8(i)}
	}

	for {
		n := 0

		sort.Sort(bySpeed(fighters))
		for _, f := range fighters {
			script := f.React(Launch, f, fighters)
			if script == nil {
				continue
			}

			n++
			observer.Observe(script)
			for _, attack := range script.Attacks {
				for _, obj := range attack.Objects {
					obj.Render(attack)
				}
			}
		}

		if n == 0 {
			break
		}
	}
}
