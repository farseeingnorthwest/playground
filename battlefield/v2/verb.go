package battlefield

type Prepare struct {
	*Action
	Actions []*Action
}

type Action struct {
	Subject Warrior
	Objects []Warrior
	Verb
}

func (a *Action) Render(warriors []Warrior) {
	signal := Prepare{Action: a}
	for _, warrior := range warriors {
		warrior.React(&signal)
	}
	for _, action := range signal.Actions {
		action.Render(warriors)
	}

	a.Verb.Render(a.Objects[0], a.Subject)

	// TODO:
}

type Verb interface {
	Render(Warrior, Warrior)
}

type Attack struct {
	Points int
}

type AttackSignal struct {
	Points int
}

type DefenseSignal struct {
	Points int
}

type DamageSignal struct {
	Points int
}

func (a *Attack) Render(object, subject Warrior) {
	attack, defense := AttackSignal{a.Points}, DefenseSignal{object.Defense()}
	subject.React(&attack)
	object.React(&defense)

	damage := DamageSignal{attack.Points - defense.Points}
	object.React(&damage)
	if damage.Points < 0 {
		damage.Points = 0
	}

	current := object.Health()
	current -= damage.Points
	if current < 0 {
		current = 0
	}

	object.SetHealth(current)
}

type Healing struct {
	Points int
}

func (h *Healing) Render(w, _ Warrior) {
	// TODO:
}

type Buffing struct {
	Buff
}

func (h *Buffing) Render(w, _ Warrior) {
	w.Add(h.Buff)
}
