package battlefield

import (
	"math/rand"
	"sort"

	"github.com/farseeingnorthwest/playground/battlefield/v2/evaluation"
)

type Selector interface {
	Select(*Fighter, []*Fighter) []*Fighter
}

type AndSelector []Selector

func (s AndSelector) Select(i *Fighter, fighters []*Fighter) []*Fighter {
	for _, selector := range s {
		fighters = selector.Select(i, fighters)
	}

	return fighters
}

type CopySelector struct{}

func (s CopySelector) Select(_ *Fighter, fighters []*Fighter) []*Fighter {
	a := make([]*Fighter, len(fighters))
	copy(a, fighters)
	return a
}

type HealthSelector struct {
}

func (s HealthSelector) Select(_ *Fighter, fighters []*Fighter) []*Fighter {
	var a []*Fighter
	for _, f := range fighters {
		if f.current.Current > 0 {
			a = append(a, f)
		}
	}

	return a
}

type SideSelector struct {
	Own bool
}

func (s SideSelector) Select(i *Fighter, fighters []*Fighter) []*Fighter {
	var a []*Fighter
	for _, f := range fighters {
		if f.Side == i.Side == s.Own {
			a = append(a, f)
		}
	}

	return a
}

type FrontSelector struct {
	Count int
}

func (s FrontSelector) Select(_ *Fighter, fighters []*Fighter) []*Fighter {
	if s.Count < len(fighters) {
		return fighters[:s.Count]
	}

	return fighters
}

type RandomSelector struct {
	Count int
}

func (s RandomSelector) Select(_ *Fighter, fighters []*Fighter) []*Fighter {
	for i := range fighters {
		if s.Count <= i {
			return fighters[:i]
		}

		j := i + rand.Intn(len(fighters)-i)
		fighters[i], fighters[j] = fighters[j], fighters[i]
	}

	return fighters
}

type AxisSelector struct {
	evaluation.Axis
	Asc bool
}

type byAxis struct {
	evaluation.Axis
	Asc     bool
	Fighter []*Fighter
}

func (a byAxis) Len() int { return len(a.Fighter) }
func (a byAxis) Less(i, j int) bool {
	return a.Fighter[i].Component(a.Axis) < a.Fighter[j].Component(a.Axis) == a.Asc
}
func (a byAxis) Swap(i, j int) { a.Fighter[i], a.Fighter[j] = a.Fighter[j], a.Fighter[i] }

func (s AxisSelector) Select(_ *Fighter, fighters []*Fighter) []*Fighter {
	sort.Sort(byAxis{s.Axis, s.Asc, fighters})

	return fighters
}

type TagSelector struct {
	Tag any
}

func (s TagSelector) Select(_ *Fighter, fighters []*Fighter) []*Fighter {
	var a []*Fighter
	for _, f := range fighters {
		if f.Contains(s.Tag) {
			a = append(a, f)
		}
	}

	return a
}

type CurrentSelector struct{}

func (s CurrentSelector) Select(current *Fighter, _ []*Fighter) []*Fighter {
	return []*Fighter{current}
}
