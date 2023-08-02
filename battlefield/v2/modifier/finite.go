package modifier

type Finite interface {
	Capacity() int
	Valid() bool
	WarmUp()
}

type FiniteModifier struct {
	capacity int
}

func NewFiniteModifier(capacity int) *FiniteModifier {
	return &FiniteModifier{
		capacity: capacity,
	}
}

func (m *FiniteModifier) Capacity() int {
	if m == nil {
		return 1
	}

	return m.capacity
}

func (m *FiniteModifier) Valid() bool {
	return m.Capacity() > 0
}

func (m *FiniteModifier) WarmUp() {
	if m == nil {
		return
	}

	m.capacity--
}

func (m *FiniteModifier) Clone() any {
	if m == nil {
		return nil
	}

	return &FiniteModifier{
		capacity: m.capacity,
	}
}
