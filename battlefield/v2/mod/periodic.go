package mod

type Periodic interface {
	Period() int
	Phase() int
	Free() bool
	CoolDown()
	WarmUp()
}

type PeriodicMod struct {
	period int
	phase  int
}

func (m *PeriodicMod) Period() int {
	return m.period
}

func (m *PeriodicMod) SetPeriod(period int) {
	m.period = period
}

func (m *PeriodicMod) Phase() int {
	return m.phase
}

func (m *PeriodicMod) SetPhase(phase int) {
	m.phase = phase
}

func (m *PeriodicMod) Free() bool {
	return m.Phase() == 0
}

func (m *PeriodicMod) CoolDown() {
	if m.phase > 0 {
		m.phase--
	}
}

func (m *PeriodicMod) WarmUp() {
	m.phase = m.period
}
