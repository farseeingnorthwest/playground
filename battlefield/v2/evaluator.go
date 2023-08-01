package battlefield

import "github.com/farseeingnorthwest/playground/battlefield/v2/evaluation"

type Evaluator interface {
	Evaluate(*Warrior, *evaluation.Block) int
}

type HeadEvaluator struct {
}

func (e *HeadEvaluator) Evaluate(_ *Warrior, block *evaluation.Block) int {
	return block.Value()
}

type PercentageEvaluator struct {
	Axis
	Percentage int
}

func (e *PercentageEvaluator) Evaluate(warrior *Warrior, block *evaluation.Block) int {
	var value int
	switch e.Axis {
	case Damage:
		value = warrior.Damage()
	case Defense:
		value = warrior.Defense()
	case Health:
		r, m := warrior.Health()
		value = r.Current * m / r.Maximum
	default:
		panic("bad axis")
	}

	return value * e.Percentage / 100
}

type EvalChain struct {
	Evaluator
	*evaluation.Block
}

func NewEvalChain(e Evaluator, block *evaluation.Block) *EvalChain {
	return &EvalChain{e, block}
}

func NewEvalChainProto(e Evaluator) *EvalChain {
	return NewEvalChain(e, nil)
}

func (e *EvalChain) ForkWith(warrior *Warrior) *evaluation.Block {
	if e.Evaluator == nil {
		return e.Block
	}

	return e.Fork(e.Evaluate(warrior, e.Block))
}
