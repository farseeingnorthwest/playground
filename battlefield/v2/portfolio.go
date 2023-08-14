package battlefield

import "container/list"

type Portfolio interface {
	Reactor

	Add(Reactor) Reactor
	Remove(Reactor)
	Buffs(tags ...any) []Reactor
}

type Priority int

type StackingLimit struct {
	reactors *list.List
	capacity int
}

func NewStackingLimit(capacity int) StackingLimit {
	return StackingLimit{list.New(), capacity}
}

func (p StackingLimit) Count() int {
	return p.reactors.Len()
}

func (p StackingLimit) Capacity() int {
	return p.capacity
}

func (p StackingLimit) Add(reactor Reactor) (overflow Reactor) {
	if p.reactors.Len() == p.capacity {
		e := p.reactors.Front()
		p.reactors.Remove(e)
		overflow = e.Value.(Reactor)
	}

	p.reactors.PushBack(reactor)
	return
}

func (p StackingLimit) Remove(reactor Reactor) {
	for e := p.reactors.Front(); e != nil; e = e.Next() {
		if e.Value == reactor {
			p.reactors.Remove(e)
			return
		}
	}
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
