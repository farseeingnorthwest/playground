package battlefield

type Script interface {
	Renderer
	Source() (any, Reactor)
}

type Action interface {
	Renderer
	Script() Script
	Targets() []Warrior
	Verb() Verb
}

type Verb interface {
	Render(target Warrior, action Action)
}

type Attack struct {
}

func (a *Attack) Render(target Warrior, action Action) {
}
