package battlefield

import (
	"errors"
)

const (
	Lt IntComparator = iota
	Le
	Eq
	Ge
	Gt
)

var (
	comparators = map[IntComparator]string{
		Lt: "<",
		Le: "<=",
		Eq: "=",
		Ge: ">=",
		Gt: ">",
	}

	ErrBadComparator = errors.New("bad comparator")
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

func (c IntComparator) String() string {
	return comparators[c]
}

func (c *IntComparator) UnmarshalText(text []byte) error {
	for i, name := range comparators {
		if string(text) == name {
			*c = IntComparator(i)
			return nil
		}
	}

	return ErrBadComparator
}
