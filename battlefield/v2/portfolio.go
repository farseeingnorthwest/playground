package battlefield

import "container/list"

type Portfolio interface {
	Reactor

	Add(Reactor)
	Remove(Reactor)
	Buffs(tags ...any) []Reactor
}

type Priority int

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
