package battlefield

type Rng interface {
	Gen() float64 // [0, 1)
}

type Buff interface {
	Reactor
}

type Critical struct {
	rng    Rng
	odds   int // percentage
	damage *TemporaryDamage
}

func (c *Critical) React(signal Signal) {
	sig, ok := signal.(*PreActionSignal)
	if !ok {
		return
	}
	_, ok = sig.Verb.(*Attack)
	if !ok {
		return
	}

	if float64(c.odds)/100 <= c.rng.Gen() {
		return
	}

	sig.Add(&Action{
		Source:  sig.Source,
		Targets: sig.Targets,
		Verb:    NewBuffing(c.damage.Fork(sig.Action)),
	})
}

type TemporaryDamage struct {
	factor int // percentage
	action *Action
}

func (c *TemporaryDamage) React(signal Signal) {
	switch sig := signal.(type) {
	case *DamageClearingSignal:
		sig.Map(func(points int) int {
			return points * c.factor / 100
		})

	case *PostActionSignal:
		if sig.Action == c.action {
			c.action = nil
		}
	}
}

func (c *TemporaryDamage) Validate() bool {
	return c.action != nil
}

func (c *TemporaryDamage) Fork(a *Action) Buff {
	return &TemporaryDamage{
		factor: c.factor,
		action: a,
	}
}
