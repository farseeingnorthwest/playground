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
	prepare, ok := signal.(*PreActionSignal)
	if !ok {
		return
	}
	_, ok = prepare.Verb.(*Attack)
	if !ok {
		return
	}

	if float64(c.odds)/100 <= c.rng.Gen() {
		return
	}

	prepare.Add(&Action{
		Subject: prepare.Subject,
		Objects: prepare.Objects,
		Verb:    &Buffing{Buff: c.damage.Fork()},
	})
}

func (c *Critical) Valid() bool {
	return true
}

type TemporaryDamage struct {
	factor int // percentage
	valid  bool
}

func (c *TemporaryDamage) React(signal Signal) {
	damage, ok := signal.(*DamageClearingSignal)
	if !ok {
		return
	}

	damage.Map(func(points int) int {
		return points * c.factor / 100
	})
	c.valid = false
}

func (c *TemporaryDamage) Valid() bool {
	return c.valid
}

func (c *TemporaryDamage) Fork() Buff {
	return &TemporaryDamage{
		factor: c.factor,
		valid:  true,
	}
}
