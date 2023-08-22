package functional

func First[A any, B any](a A, _ B) A {
	return a
}

func Second[A any, B any](_ A, b B) B {
	return b
}
