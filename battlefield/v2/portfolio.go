package battlefield

import (
	"container/list"
)

var (
	_ Portfolio = (*FatPortfolio)(nil)
)

type Portfolio interface {
	Reactor

	Add(Reactor) Reactor
	Remove(Reactor)
	Buffs(tags ...any) []Reactor
	Stacking(StackingLimit) *Stack
}

type FatPortfolio struct {
	reactors *list.List
	stacking map[StackingLimit]*Stack
}

func NewFatPortfolio() *FatPortfolio {
	return &FatPortfolio{list.New(), make(map[StackingLimit]*Stack)}
}

func (p *FatPortfolio) React(signal Signal, ec EvaluationContext) {
	for e := p.reactors.Front(); e != nil; e = e.Next() {
		e.Value.(Reactor).React(signal, ec)
	}
}

func (p *FatPortfolio) Active() bool {
	return true
}

func (p *FatPortfolio) Add(reactor Reactor) (overflow Reactor) {
	if lm, ok := QueryTag[StackingLimit](reactor); ok {
		if _, ok := p.stacking[lm]; !ok {
			p.stacking[lm] = NewStack(lm.Capacity)
		}
		if overflow = p.stacking[lm].Add(reactor); overflow != nil {
			p.remove(overflow)
		}
	}

	pr, ok := QueryTag[Priority](reactor)
	if !ok {
		pr = 0
	}
	for e := p.reactors.Front(); e != nil; e = e.Next() {
		pr2, ok := QueryTag[Priority](e.Value)
		if !ok {
			pr2 = 0
		}
		if pr > pr2 {
			p.reactors.InsertBefore(reactor, e)
			return
		}
	}

	p.reactors.PushBack(reactor)
	return
}

func (p *FatPortfolio) Remove(reactor Reactor) {
	if lm, ok := QueryTag[StackingLimit](reactor); ok {
		p.stacking[lm].Remove(reactor)
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
		r := e.Value.(Reactor)
		if !r.Active() {
			continue
		}
		if len(tags) == 0 {
			buffs = append(buffs, r)
			continue
		}
		if tagger, ok := r.(Tagger); ok && tagger.Match(tags...) {
			buffs = append(buffs, r)
		}
	}

	return
}

func (p *FatPortfolio) Stacking(lm StackingLimit) *Stack {
	return p.stacking[lm]
}

type Stack struct {
	reactors *list.List
	capacity int
}

func NewStack(capacity int) *Stack {
	return &Stack{list.New(), capacity}
}

func (l *Stack) Count() int {
	return l.reactors.Len()
}

func (l *Stack) Capacity() int {
	return l.capacity
}

func (l *Stack) Add(reactor Reactor) (overflow Reactor) {
	if l.reactors.Len() == l.capacity {
		e := l.reactors.Front()
		l.reactors.Remove(e)
		overflow = e.Value.(Reactor)
	}

	l.reactors.PushBack(reactor)
	return
}

func (l *Stack) Remove(reactor Reactor) {
	for e := l.reactors.Front(); e != nil; e = e.Next() {
		if e.Value == reactor {
			l.reactors.Remove(e)
			return
		}
	}
}
