package battlefield

type Scripter interface {
	Push(Signal, Reactor, chan Action)
	Pop()
}

type scripter struct {
	scripts []Script
}

func (s *scripter) Push(signal Signal, reactor Reactor, aChan chan Action) {
	s.scripts = append(s.scripts, newScript(signal, reactor, aChan))
}

func (s *scripter) Pop() {
	s.scripts = s.scripts[:len(s.scripts)-1]
}

func (s *scripter) Render(b *BattleField) {
	for _, script := range s.scripts {
		script.Render(b)
	}
}
