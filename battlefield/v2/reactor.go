package battlefield

type Reactor interface {
	React(Signal)
}

type Portfolio interface {
	Reactor

	Add(Reactor)
	Contains(any) bool
}
