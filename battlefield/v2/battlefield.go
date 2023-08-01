package battlefield

import (
	"github.com/farseeingnorthwest/playground/battlefield/v2/modifier"
	"sort"
)

type Side uint8

const (
	Left Side = iota
	Right
)

type Observer interface {
	Observe(*Action)
}

type Fighter struct {
	*Warrior
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

	return f[i].Position < f[i].Position
}
func (f bySpeed) Swap(i, j int) { f[i], f[j] = f[j], f[i] }

type BattleField struct {
	fighters []*Fighter
	reactors []Reactor
	observer Observer
}

func NewBattleField(a, b []*Warrior, reactors []Reactor, observer Observer) *BattleField {
	fighters := make([]*Fighter, len(a)+len(b))
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

func (b *BattleField) Warriors() []*Fighter {
	return b.fighters
}

func (b *BattleField) React(signal Signal) {
	for _, f := range b.fighters {
		f.React(signal)
	}
	for _, reactor := range b.reactors {
		if r, ok := reactor.(modifier.Finite); ok && !r.Valid() {
			continue
		}
		if r, ok := reactor.(modifier.Periodic); ok && !r.Free() {
			continue
		}

		reactor.React(signal)
	}
}

func (b *BattleField) Fight() {
	for {
		n := 0

		sorted := false
		for i := 0; i < len(b.fighters); i++ {
			if !sorted {
				sort.Sort(bySpeed(b.fighters[i:]))
				sorted = true
			}
			if b.fighters[i].current.Current <= 0 {
				continue
			}

			sorted = false
			signal := NewLaunchSignal(b.fighters[i], b)
			b.fighters[i].React(signal)
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
