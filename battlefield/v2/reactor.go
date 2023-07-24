package battlefield

type Reactor interface {
	React(Signal, *Fighter, []*Fighter) *Script
}

type NormalAttack struct {
	Selector
	Points int
}

func (a *NormalAttack) React(signal Signal, i *Fighter, fighters []*Fighter) *Script {
	if signal != Launch || !i.Functional() {
		return nil
	}

	fighter := a.Select(i, fighters)
	if fighter == nil {
		return nil
	}

	return &Script{
		Subject: i,
		Attacks: []Attack{
			{Objects: []*Fighter{fighter}, Points: a.Points},
		},
	}
}
