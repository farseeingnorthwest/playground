package battlefield

import "github.com/farseeingnorthwest/playground/battlefield/v2/evaluation"

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
	sig := NewEvaluationSignal(evaluation.Damage, w.Baseline.Damage(), nil)
	w.React(sig)

	return sig.Value()
}

func (w *Warrior) Defense() int {
	sig := NewEvaluationSignal(evaluation.Defense, w.Baseline.Defense(), nil)
	w.React(sig)

	return sig.Value()
}

func (w *Warrior) Speed() int {
	sig := NewEvaluationSignal(evaluation.Speed, w.Baseline.Speed(), nil)
	w.React(sig)

	return sig.Value()
}

func (w *Warrior) Health() (Ratio, int) {
	sig := NewEvaluationSignal(evaluation.Health, w.Baseline.Health(), nil)
	w.React(sig)

	return w.current, sig.Value()
}

func (w *Warrior) SetHealth(value Ratio) {
	w.current = value
}

func (w *Warrior) Component(axis evaluation.Axis) (value int) {
	switch axis {
	case evaluation.Damage:
		value = w.Damage()
	case evaluation.Defense:
		value = w.Defense()
	case evaluation.Health:
		r, m := w.Health()
		value = r.Current * m / r.Maximum
	case evaluation.HealthMax:
		_, value = w.Health()
	case evaluation.Speed:
		value = w.Speed()
	default:
		panic("bad axis")
	}

	return
}
