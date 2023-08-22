package battlefield

const (
	Lt IntComparator = iota
	Le
	Eq
	Ge
	Gt
)

type IntComparator uint8

func (c IntComparator) Compare(a, b int) bool {
	switch c {
	case Lt:
		return a < b
	case Le:
		return a <= b
	case Eq:
		return a == b
	case Ge:
		return a >= b
	case Gt:
		return a > b

	default:
		panic("bad comparator")
	}
}
