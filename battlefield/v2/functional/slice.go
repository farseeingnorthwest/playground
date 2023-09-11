package functional

func Map[A any, B any](f func(A) B) func([]A) []B {
	return func(as []A) []B {
		if len(as) == 0 {
			return nil
		}

		bs := make([]B, len(as))
		for i, a := range as {
			bs[i] = f(a)
		}

		return bs
	}
}
