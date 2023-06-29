package battlefield

type Attribute int

const (
	Attack Attribute = iota
	Critical
	Defense
	Health
	HealthCritical
	Velocity
	NumberOfAttributes
)

type Baseline struct {
	Attack   int
	Critical float64
	Defense  int
	Health   int
}

type Carrier struct {
	Attack   float64
	Critical float64
	Defense  float64
	Health   float64
	Velocity int
}

type Magnitude struct {
	Attack   float64
	Critical float64
	Defense  float64
	Health   float64
}

type Warrior struct {
	attack   int
	critical float64
	defense  int
	health   int
	velocity int
	buffers  []*bufferList

	criticalAttack bool
}

func NewWarrior(baseline *Baseline, carrier *Carrier, magnitude *Magnitude) *Warrior {
	buffers := make([]*bufferList, NumberOfAttributes)
	for i := range buffers {
		buffers[i] = &bufferList{}
	}
	buffers[HealthCritical].Append(healthCriticalBaseline{})

	return &Warrior{
		int(float64(baseline.Attack) * (carrier.Attack + magnitude.Attack)),
		baseline.Critical * (carrier.Critical + magnitude.Critical),
		int(float64(baseline.Defense) * (carrier.Defense + magnitude.Defense)),
		int(float64(baseline.Health) * (carrier.Health + magnitude.Health)),
		carrier.Velocity,
		buffers,
		false,
	}
}

func (w *Warrior) Attach(attribute Attribute, buffer Buffer) {
	w.buffers[attribute].Append(buffer)
}

func (w *Warrior) Prepare(r Randomizer) {
	for _, l := range w.buffers {
		l.Drain()
	}

	w.criticalAttack = r.Float64() < w.buffers[Critical].Buff(w.critical)
}

func (w *Warrior) Attack() (attack int, critical bool) {
	return int(w.buffers[Attack].Buff(float64(w.attack))), w.criticalAttack
}

func (w *Warrior) Velocity() int {
	return int(w.buffers[Velocity].Buff(float64(w.velocity)))
}

func (w *Warrior) Suffer(attack int, critical bool) (damage, overflow int) {
	d := float64(attack) - w.buffers[Defense].Buff(float64(w.defense))
	if d < 0 {
		return 0, 0
	}

	if critical {
		d = w.buffers[HealthCritical].Buff(d)
	}
	damage = int(w.buffers[Health].Buff(d))

	w.health -= damage
	if w.health < 0 {
		overflow = -w.health
		w.health = 0
	}

	return
}

func (w *Warrior) Health() int {
	return w.health
}
