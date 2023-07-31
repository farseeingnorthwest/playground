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
		a.Verb.Render(object.Warrior, a.Source.Warrior, a)
	}

	post := NewPostActionSignal(a)
	f.React(post)
	for _, action := range post.Actions() {
		action.Render(f)
	}
}

type Verb interface {
	Render(*Warrior, *Warrior, *Action)
}

type Attack struct {
	points int
}

func NewAttack(points int) *Attack {
	return &Attack{
		points: points,
	}
}

func (a *Attack) Render(target, source *Warrior, action *Action) {
	damage := NewEvaluationSignal(Damage, a.points, action)
	source.React(damage)
	defense := NewEvaluationSignal(Defense, target.Defense(), action)
	target.React(defense)

	loss := NewEvaluationSignal(Loss, damage.Value()-defense.Value(), action)
	target.React(loss)
	if loss.Value() < 0 {
		loss.SetValue(0)
	}

	r, m := target.Health()
	c := r.Current * m / r.Maximum
	c -= loss.Value()
	if c < 0 {
		c = 0
	}

	target.current = Ratio{c, m}
}

type Heal struct {
	points int
}

func NewHeal(points int) *Heal {
	return &Heal{
		points: points,
	}
}

func (h *Heal) Render(target, _ *Warrior, action *Action) {
	heal := NewEvaluationSignal(Healing, h.points, action)
	target.React(heal)
	if heal.Value() < 0 {
		heal.SetValue(0)
	}

	r, m := target.Health()
	c := r.Current * m / r.Maximum
	c += heal.Value()
	if c > m {
		c = m
	}

	target.current = Ratio{c, m}
}

type Buffing struct {
	buff Reactor
}

func NewBuffing(buff Reactor) *Buffing {
	return &Buffing{
		buff: buff,
	}
}

func (h *Buffing) Render(target, _ *Warrior, _ *Action) {
	target.Append(h.buff)
}

type Purging struct {
}

func NewPurging() *Purging {
	return &Purging{}
}

func (p *Purging) Render(target, _ *Warrior, _ *Action) {
	// TODO:
}
