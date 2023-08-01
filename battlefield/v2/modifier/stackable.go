package modifier

type Stackable interface {
	StackLimit() int
}

type StackableModifier struct {
	stackLimit int
}

func NewStackableModifier(stackLimit int) *StackableModifier {
	return &StackableModifier{
		stackLimit: stackLimit,
	}
}

func (m *StackableModifier) StackLimit() int {
	return m.stackLimit
}

func (m *StackableModifier) Clone() any {
	return &StackableModifier{
		stackLimit: m.stackLimit,
	}
}
