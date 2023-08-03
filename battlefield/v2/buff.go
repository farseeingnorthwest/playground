package battlefield

import (
	"github.com/farseeingnorthwest/playground/battlefield/v2/evaluation"
	"github.com/farseeingnorthwest/playground/battlefield/v2/mod"
)

type EvaluationBuff struct {
	mod.TaggerMod
	*mod.FiniteMod

	axis       evaluation.Axis
	bias       int
	multiplier int
}

func NewEvaluationBuff(tag any, axis evaluation.Axis, options ...func(buff *EvaluationBuff)) *EvaluationBuff {
	buff := &EvaluationBuff{
		axis:       axis,
		bias:       0,
		multiplier: 100,
	}

	buff.SetTag(tag)
	for _, option := range options {
		option(buff)
	}

	return buff
}

func EvaluationBias(bias int) func(buff *EvaluationBuff) {
	return func(buff *EvaluationBuff) {
		buff.bias = bias
	}
}

func EvaluationMultiplier(multiplier int) func(buff *EvaluationBuff) {
	return func(buff *EvaluationBuff) {
		buff.multiplier = multiplier
	}
}

func EvaluationCapacity(capacity int) func(buff *EvaluationBuff) {
	return func(buff *EvaluationBuff) {
		buff.FiniteMod = mod.NewFiniteModifier(capacity)
	}
}

func (b *EvaluationBuff) React(signal Signal) {
	switch sig := signal.(type) {
	case *EvaluationSignal:
		if sig.Axis() != b.axis || sig.Action() != nil {
			return
		}

		sig.Map(func(points int) int {
			return points*b.multiplier/100 + b.bias
		})

	case *RoundEndSignal:
		b.WarmUp()
	}
}

func (b *EvaluationBuff) Fork(*evaluation.Block, Signal) any {
	return &EvaluationBuff{
		TaggerMod:  b.TaggerMod,
		FiniteMod:  b.FiniteMod.Clone().(*mod.FiniteMod),
		axis:       b.axis,
		bias:       b.bias,
		multiplier: b.multiplier,
	}
}

type ClearingBuff struct {
	mod.TaggerMod

	axis       evaluation.Axis
	bias       int
	multiplier int
	action     *Action
}

func NewClearingBuff(tag any, axis evaluation.Axis, action *Action, options ...func(buff *ClearingBuff)) *ClearingBuff {
	buff := &ClearingBuff{
		axis:       axis,
		bias:       0,
		multiplier: 100,

		action: action,
	}

	buff.SetTag(tag)
	for _, option := range options {
		option(buff)
	}

	return buff
}

func ClearingBias(bias int) func(buff *ClearingBuff) {
	return func(buff *ClearingBuff) {
		buff.bias = bias
	}
}

func ClearingMultiplier(multiplier int) func(buff *ClearingBuff) {
	return func(buff *ClearingBuff) {
		buff.multiplier = multiplier
	}
}

func (b *ClearingBuff) React(signal Signal) {
	switch sig := signal.(type) {
	case *EvaluationSignal:
		if sig.Axis() != b.axis || sig.Action() != b.action {
			return
		}

		sig.Map(func(points int) int {
			return points*b.multiplier/100 + b.bias
		})

	case *PostActionSignal:
		if sig.Action == b.action {
			b.WarmUp()
		}
	}
}

func (b *ClearingBuff) Capacity() int {
	if b.action == nil {
		return 0
	}

	return 1
}

func (b *ClearingBuff) WarmUp() {
	b.action = nil
}

func (b *ClearingBuff) Fork(_ *evaluation.Block, signal Signal) any {
	action := b.action
	if signal != nil {
		action = signal.(*PreActionSignal).Action
	}

	return &ClearingBuff{
		TaggerMod:  b.TaggerMod,
		axis:       b.axis,
		bias:       b.bias,
		multiplier: b.multiplier,
		action:     action,
	}
}

type CompoundBuff struct {
	mod.TaggerMod
	*mod.FiniteMod
	buffs []*EvaluationBuff
}

func NewCompoundBuff(tag any, buffs []*EvaluationBuff, options ...func(buff *CompoundBuff)) *CompoundBuff {
	buff := &CompoundBuff{
		buffs: buffs,
	}

	buff.SetTag(tag)
	for _, option := range options {
		option(buff)
	}

	return buff
}

func TaggedCapacity(capacity int) func(buff *CompoundBuff) {
	return func(buff *CompoundBuff) {
		buff.FiniteMod = mod.NewFiniteModifier(capacity)
	}
}

func (b *CompoundBuff) React(signal Signal) {
	switch sig := signal.(type) {
	case *EvaluationSignal:
		if sig.Action() != nil {
			return
		}

		for _, buff := range b.buffs {
			buff.React(sig)
		}

	case *RoundEndSignal:
		b.WarmUp()
	}
}

func (b *CompoundBuff) Fork(block *evaluation.Block, signal Signal) any {
	return &CompoundBuff{
		TaggerMod: b.TaggerMod,
		FiniteMod: b.FiniteMod.Clone().(*mod.FiniteMod),
		buffs:     b.buffs,
	}
}
