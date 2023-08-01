package modifier

type Periodic interface {
	Period() int
	Phase() int
	Free() bool
	CoolDown()
	WarmUp()
}

type PeriodicModifier struct {
	period int
	phase  int
}

func (m *PeriodicModifier) Period() int {
	return m.period
}

func (m *PeriodicModifier) Phase() int {
	return m.phase
}

func (m *PeriodicModifier) Free() bool {
	return m.Phase() == 0
}

func (m *PeriodicModifier) CoolDown() {
	if m.phase > 0 {
		m.phase--
	}
}

func (m *PeriodicModifier) WarmUp() {
	m.phase = m.period
}

func (m *PeriodicModifier) Clone() any {
	return &PeriodicModifier{
		period: m.period,
		phase:  m.phase,
	}
}
