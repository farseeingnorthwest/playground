package battlefield

import (
	"math/rand"
	"sort"
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

type SequenceSelector struct {
	Evaluate func(*Fighter) int
	Asc      bool
}

type byEvaluate struct {
	Evaluate func(*Fighter) int
	Asc      bool
	Fighter  []*Fighter
}

func (e byEvaluate) Len() int { return len(e.Fighter) }
func (e byEvaluate) Less(i, j int) bool {
	return e.Evaluate(e.Fighter[i]) < e.Evaluate(e.Fighter[j]) == e.Asc
}
func (e byEvaluate) Swap(i, j int) { e.Fighter[i], e.Fighter[j] = e.Fighter[j], e.Fighter[i] }

func (s SequenceSelector) Select(_ *Fighter, fighters []*Fighter) []*Fighter {
	sort.Sort(byEvaluate{s.Evaluate, s.Asc, fighters})

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
