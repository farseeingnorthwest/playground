package battlefield

type Action struct {
	Subject Warrior
	Objects []Warrior
	Verb
}

func (a *Action) Render(f *BattleField) {
	signal := NewPreActionSignal(a)
	f.React(signal)
	for _, action := range signal.Actions() {
		action.Render(f)
	}

	for _, object := range a.Objects {
		a.Verb.Render(object, a.Subject)
	}

	// TODO:
}

type Verb interface {
	Render(Warrior, Warrior)
}

type Attack struct {
	Points int
}

func (a *Attack) Render(object, subject Warrior) {
	attack, defense := NewAttackClearingSignal(a.Points), NewDefenseClearingSignal(object.Defense())
	subject.React(attack)
	object.React(defense)

	damage := NewDamageClearingSignal(attack.Value() - defense.Value())
	object.React(damage)
	if damage.Value() < 0 {
		damage.SetValue(0)
	}

	current := object.Health()
	current -= damage.Value()
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
