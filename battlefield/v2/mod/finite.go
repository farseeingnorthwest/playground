package mod

type Finite interface {
	Capacity() int
	Valid() bool
	WarmUp()
}

type FiniteMod struct {
	capacity int
}

func NewFiniteModifier(capacity int) *FiniteMod {
	return &FiniteMod{
		capacity: capacity,
	}
}

func (m *FiniteMod) Capacity() int {
	if m == nil {
		return 1
	}

	return m.capacity
}

func (m *FiniteMod) Valid() bool {
	return m.Capacity() > 0
}

func (m *FiniteMod) WarmUp() {
	if m == nil {
		return
	}

	m.capacity--
}

func (m *FiniteMod) Clone() any {
	if m == nil {
		return nil
	}

	return &FiniteMod{
		capacity: m.capacity,
	}
}
