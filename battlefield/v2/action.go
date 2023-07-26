package battlefield

type Action struct {
	Source  *Fighter
	Targets []*Fighter
	Verb
}

func (a *Action) Render(f *BattleField) {
	pre := NewPreActionSignal(a)
	f.React(pre)
	for _, action := range pre.Actions() {
		action.Render(f)
	}

	for _, object := range a.Targets {
		a.Verb.Render(object, a.Source)
	}

	post := NewPostActionSignal(a)
	f.React(post)
	for _, action := range post.Actions() {
		action.Render(f)
	}
}

type Verb interface {
	Render(Warrior, Warrior)
}

type Attack struct {
	points int
}

func NewAttack(points int) *Attack {
	return &Attack{
		points: points,
	}
}

func (a *Attack) Render(target, source Warrior) {
	attack, defense := NewAttackClearingSignal(a.points), NewDefenseClearingSignal(target.Defense())
	source.React(attack)
	target.React(defense)

	damage := NewDamageClearingSignal(attack.Value() - defense.Value())
	target.React(damage)
	if damage.Value() < 0 {
		damage.SetValue(0)
	}

	current := target.Health()
	current -= damage.Value()
	if current < 0 {
		current = 0
	}

	target.SetHealth(current)
}

type Healing struct {
	points int
}

func NewHealing(points int) *Healing {
	return &Healing{
		points: points,
	}
}

func (h *Healing) Render(_, _ Warrior) {
	// TODO:
}

type Buffing struct {
	buff Buff
}

func NewBuffing(buff Buff) *Buffing {
	return &Buffing{
		buff: buff,
	}
}

func (h *Buffing) Render(target, _ Warrior) {
	target.Add(h.buff)
}

type Purging struct {
}

func NewPurging() *Purging {
	return &Purging{}
}

func (p *Purging) Render(target, _ Warrior) {
	// TODO:
}
