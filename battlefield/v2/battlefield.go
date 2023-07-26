package battlefield

import (
	"sort"
)

type Side uint8

const (
	Left Side = iota
	Right
)

type Element uint8

type Warrior interface {
	Portfolio

	Element() Element
	Defense() int
	Health() int
	SetHealth(int)
	Speed() int
}

type Observer interface {
	Observe(*Action)
}

type Fighter struct {
	Warrior
	Side
	Position uint8
}

type bySpeed []Warrior

func (f bySpeed) Len() int { return len(f) }
func (f bySpeed) Less(i, j int) bool {
	if f[i].Speed() != f[j].Speed() {
		return f[i].Speed() > f[j].Speed()
	}
	a, b := f[i].(*Fighter), f[j].(*Fighter)
	if a.Side != b.Side {
		return a.Side == Left
	}

	return a.Position < a.Position
}
func (f bySpeed) Swap(i, j int) { f[i], f[j] = f[j], f[i] }

type BattleField struct {
	fighters []Warrior
	reactors []Reactor
	observer Observer
}

func NewBattleField(a, b []Warrior, reactors []Reactor, observer Observer) *BattleField {
	fighters := make([]Warrior, len(a)+len(b))
	for i, f := range a {
		fighters[i] = &Fighter{f, Left, uint8(i)}
	}
	for i, f := range b {
		fighters[i+len(a)] = &Fighter{f, Right, uint8(i)}
	}

	return &BattleField{
		fighters: fighters,
		reactors: reactors,
		observer: observer,
	}
}

func (b *BattleField) React(signal Signal) {
	for _, f := range b.fighters {
		if !f.Valid() {
			continue
		}

		f.React(signal)
	}
	for _, reactor := range b.reactors {
		if !reactor.Valid() {
			continue
		}

		reactor.React(signal)
	}
}

func (b *BattleField) Valid() bool {
	return true
}

func (b *BattleField) Fight() {
	for {
		n := 0

		sort.Sort(bySpeed(b.fighters))
		for _, f := range b.fighters {
			if f.Health() <= 0 {
				continue
			}

			signal := NewLaunchingSignal(f, b.fighters)
			f.React(signal)
			for _, a := range signal.Actions() {
				n++
				b.observer.Observe(a)
				a.Render(b)
			}
		}

		if n == 0 {
			break
		}
	}
}
