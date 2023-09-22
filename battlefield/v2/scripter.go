package battlefield

type Scripter interface {
	Push(Signal, Reactor)
	Pop()
	Add(Action)
}

type scripter struct {
	scripts []Script
}

func (s *scripter) Push(signal Signal, reactor Reactor) {
	s.scripts = append(s.scripts, newScript(signal, reactor))
}

func (s *scripter) Pop() {
	s.scripts = s.scripts[:len(s.scripts)-1]
}

func (s *scripter) Add(action Action) {
	s.scripts[len(s.scripts)-1].Add(action)
}

func (s *scripter) Render(b *BattleField) {
	for _, script := range s.scripts {
		script.Render(b)
	}
}
