package battlefield

import (
	"container/list"
	"encoding/json"
)

var (
	_ Portfolio = (*FatPortfolio)(nil)
)

type Priority int

func (p Priority) MarshalJSON() ([]byte, error) {
	return json.Marshal(pr{int(p)})
}

type pr struct {
	Priority int `json:"priority"`
}

type Portfolio interface {
	Reactor

	Add(Reactor) Reactor
	Remove(Reactor)
	Buffs(tags ...any) []Reactor
}

type FatPortfolio struct {
	reactors *list.List
}

func NewFatPortfolio() *FatPortfolio {
	return &FatPortfolio{list.New()}
}

func (p *FatPortfolio) React(signal Signal, ec EvaluationContext) {
	for e := p.reactors.Front(); e != nil; e = e.Next() {
		e.Value.(Reactor).React(signal, ec)
	}
}

func (p *FatPortfolio) Add(reactor Reactor) (overflow Reactor) {
	if stacking, ok := QueryTag[StackingLimit](reactor); ok {
		if overflow = stacking.Add(reactor); overflow != nil {
			p.remove(overflow)
		}
	}

	if pr, ok := QueryTag[Priority](reactor); ok {
		for e := p.reactors.Front(); e != nil; e = e.Next() {
			pr2, ok := QueryTag[Priority](e.Value)
			if !ok || pr > pr2 {
				p.reactors.InsertBefore(reactor, e)
				return
			}
		}
	}

	p.reactors.PushBack(reactor)
	return
}

func (p *FatPortfolio) Remove(reactor Reactor) {
	if stacking, ok := QueryTag[StackingLimit](reactor); ok {
		stacking.Remove(reactor)
	}

	p.remove(reactor)
}

func (p *FatPortfolio) remove(reactor Reactor) {
	for e := p.reactors.Front(); e != nil; e = e.Next() {
		if e.Value == reactor {
			p.reactors.Remove(e)
			return
		}
	}
}

func (p *FatPortfolio) Buffs(tags ...any) (buffs []Reactor) {
	for e := p.reactors.Front(); e != nil; e = e.Next() {
		if len(tags) == 0 {
			buffs = append(buffs, e.Value.(Reactor))
		} else if tagger, ok := e.Value.(Tagger); ok && tagger.Match(tags...) {
			buffs = append(buffs, e.Value.(Reactor))
		}
	}

	return
}

type StackingLimit struct {
	reactors *list.List
	capacity int
}

func NewStackingLimit(capacity int) StackingLimit {
	return StackingLimit{list.New(), capacity}
}

func (l StackingLimit) Count() int {
	return l.reactors.Len()
}

func (l StackingLimit) Capacity() int {
	return l.capacity
}

func (l StackingLimit) Add(reactor Reactor) (overflow Reactor) {
	if l.reactors.Len() == l.capacity {
		e := l.reactors.Front()
		l.reactors.Remove(e)
		overflow = e.Value.(Reactor)
	}

	l.reactors.PushBack(reactor)
	return
}

func (l StackingLimit) Remove(reactor Reactor) {
	for e := l.reactors.Front(); e != nil; e = e.Next() {
		if e.Value == reactor {
			l.reactors.Remove(e)
			return
		}
	}
}
