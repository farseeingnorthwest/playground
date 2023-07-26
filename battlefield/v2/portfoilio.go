package battlefield

type Portfolio interface {
	Reactor

	Add(Reactor)
}

type FatPortfolio struct {
	reactors []Reactor
}

func (p *FatPortfolio) React(signal Signal) {
	for _, buff := range p.reactors {
		if v, ok := buff.(Validator); ok && !v.Validate() {
			continue
		}

		buff.React(signal)
	}
}

func (p *FatPortfolio) Add(reactor Reactor) {
	p.reactors = append(p.reactors, reactor)
}
