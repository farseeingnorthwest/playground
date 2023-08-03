package mod

type Stackable interface {
	StackLimit() int
}

type StackableMod struct {
	stackLimit int
}

func NewStackableModifier(stackLimit int) *StackableMod {
	return &StackableMod{
		stackLimit: stackLimit,
	}
}

func (m *StackableMod) StackLimit() int {
	return m.stackLimit
}

func (m *StackableMod) Clone() any {
	return &StackableMod{
		stackLimit: m.stackLimit,
	}
}
