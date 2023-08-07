package battlefield

type Portfolio interface {
	Reactor

	Add(Reactor)
	Remove(Reactor)
	Buffs(tags ...any) []Reactor
}

type FatPortfolio struct {
	reactors []Reactor
}

func NewFatPortfolio() *FatPortfolio {
	return &FatPortfolio{}
}

func (p *FatPortfolio) React(signal Signal, warriors []Warrior) {
	for _, r := range p.reactors {
		r.React(signal, warriors)
	}
}

func (p *FatPortfolio) Add(reactor Reactor) {
	p.reactors = append(p.reactors, reactor)
}

func (p *FatPortfolio) Remove(reactor Reactor) {
	for i, r := range p.reactors {
		if r == reactor {
			p.reactors = append(p.reactors[:i], p.reactors[i+1:]...)
			return
		}
	}
}

func (p *FatPortfolio) Buffs(tags ...any) (buffs []Reactor) {
	if len(tags) == 0 {
		return p.reactors
	}

	for _, r := range p.reactors {
		if tagger, ok := r.(Tagger); ok && tagger.Match(tags...) {
			buffs = append(buffs, r)
		}
	}

	return
}
