package functional

type Option[T any] struct {
	value T
	ok    bool
}

func Some[T any](value T) Option[T] {
	return Option[T]{value, true}
}

func None[T any]() Option[T] {
	return Option[T]{}
}

func (o Option[T]) Value() T {
	return o.value
}

func (o Option[T]) Ok() bool {
	return o.ok
}

func (o Option[T]) UnwrapOr(value T) T {
	if o.ok {
		return o.value
	}

	return value
}

func (o Option[T]) Destruct() (T, bool) {
	return o.value, o.ok
}
