package battlefield

type Portfolio interface {
	Reactor

	Append(Reactor)
	Contains(any) bool
}

type FatPortfolio struct {
	reactors []Reactor
}

func NewFatPortfolio() *FatPortfolio {
	return &FatPortfolio{}
}

func (p *FatPortfolio) React(signal Signal) {
	for _, reactor := range p.reactors {
		if r, ok := reactor.(Finite); ok && !r.Valid() {
			continue
		}
		if r, ok := reactor.(Periodic); ok && !r.Free() {
			continue
		}

		reactor.React(signal)
	}
}

func (p *FatPortfolio) Append(reactor Reactor) {
	p.reactors = append(p.reactors, reactor)
}

func (p *FatPortfolio) Contains(tag any) bool {
	for _, reactor := range p.reactors {
		if r, ok := reactor.(Tagged); ok && r.Tag() == tag {
			return true
		}
	}

	return false
}
