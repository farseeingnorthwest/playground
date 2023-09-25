package battlefield

type Scripter interface {
	Push(Signal, Reactor, chan Instruction)
	Pop()
}

type scripter struct {
	scripts []Script
}

func (s *scripter) Push(signal Signal, reactor Reactor, ich chan Instruction) {
	s.scripts = append(s.scripts, newScript(signal, reactor, ich))
}

func (s *scripter) Pop() {
	s.scripts = s.scripts[:len(s.scripts)-1]
}

func (s *scripter) Render(b *BattleField) {
	for _, script := range s.scripts {
		script.Render(b)
	}
}
