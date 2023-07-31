package battlefield

type Baseline interface {
	Element() Element
	Damage() int
	Defense() int
	Speed() int
	Health() int
}

type Ratio struct {
	Current int
	Maximum int
}

type Warrior struct {
	Baseline
	Portfolio

	current Ratio
}

func NewWarrior(b Baseline, portfolio Portfolio) *Warrior {
	return &Warrior{
		Baseline:  b,
		Portfolio: portfolio,
		current:   Ratio{b.Health(), b.Health()},
	}
}

func (w *Warrior) Damage() int {
	sig := NewEvaluationSignal(Damage, false, w.Baseline.Damage())
	w.React(sig)

	return sig.Value()
}

func (w *Warrior) Defense() int {
	sig := NewEvaluationSignal(Defense, false, w.Baseline.Defense())
	w.React(sig)

	return sig.Value()
}

func (w *Warrior) Speed() int {
	sig := NewEvaluationSignal(Speed, false, w.Baseline.Speed())
	w.React(sig)

	return sig.Value()
}

func (w *Warrior) Health() (Ratio, int) {
	sig := NewEvaluationSignal(Health, false, w.Baseline.Health())
	w.React(sig)

	return w.current, sig.Value()
}

func (w *Warrior) SetHealth(value Ratio) {
	w.current = value
}
