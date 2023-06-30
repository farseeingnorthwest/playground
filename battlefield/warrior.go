package battlefield

type Attribute int

const (
	Attack Attribute = iota
	Critical
	Defense
	Health
	HealthCritical
	Speed
	NumberOfAttributes
)

type Baseline struct {
	Attack   int
	Critical int // percentage
	Defense  int
	Health   int
	Speed    int
}

// Carrier percentage corrections
type Carrier struct {
	Attack   int
	Critical int
	Defense  int
	Health   int
	Speed    int
}

// Magnitude percentage corrections
type Magnitude struct {
	Attack   int
	Critical int
	Defense  int
	Health   int
	Speed    int
}

type Warrior struct {
	attack   int
	critical int // percentage
	defense  int
	health   int
	speed    int
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
		baseline.Attack * (carrier.Attack + magnitude.Attack) / 100,
		baseline.Critical * (carrier.Critical + magnitude.Critical) / 100,
		baseline.Defense * (carrier.Defense + magnitude.Defense) / 100,
		baseline.Health * (carrier.Health + magnitude.Health) / 100,
		baseline.Speed * (carrier.Speed + magnitude.Speed) / 100,
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

	w.criticalAttack = r.Float64()*100 < w.buffers[Critical].Buff(float64(w.critical))
}

func (w *Warrior) Attack() (attack int, critical bool) {
	return int(w.buffers[Attack].Buff(float64(w.attack))), w.criticalAttack
}

func (w *Warrior) Speed() int {
	return int(w.buffers[Speed].Buff(float64(w.speed)))
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
