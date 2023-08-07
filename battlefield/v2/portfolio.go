package battlefield

import "container/list"

type Portfolio interface {
	Reactor

	Add(Reactor)
	Remove(Reactor)
	Buffs(tags ...any) []Reactor
}

type Priority int

func priorityTag(a any) any {
	tagger, ok := a.(Tagger)
	if !ok {
		return nil
	}

	return tagger.Find(NewTypeMatcher(Priority(0)))
}

type FatPortfolio struct {
	reactors *list.List
}

func NewFatPortfolio() *FatPortfolio {
	return &FatPortfolio{list.New()}
}

func (p *FatPortfolio) React(signal Signal, warriors []Warrior) {
	for e := p.reactors.Front(); e != nil; e = e.Next() {
		e.Value.(Reactor).React(signal, warriors)
	}
}

func (p *FatPortfolio) Add(reactor Reactor) {
	if pr := priorityTag(reactor); pr != nil {
		for e := p.reactors.Front(); e != nil; e = e.Next() {
			pr2 := priorityTag(e.Value)
			if pr2 == nil || pr.(Priority) > pr2.(Priority) {
				p.reactors.InsertBefore(reactor, e)
				return
			}
		}
	}

	p.reactors.PushBack(reactor)
}

func (p *FatPortfolio) Remove(reactor Reactor) {
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
