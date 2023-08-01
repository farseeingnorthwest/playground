package battlefield

import (
	"github.com/farseeingnorthwest/playground/battlefield/v2/evaluation"
	"github.com/farseeingnorthwest/playground/battlefield/v2/modifier"
)

type EvaluationBuff struct {
	axis       evaluation.Axis
	bias       int
	multiplier int

	*modifier.FiniteModifier
}

func NewEvaluationBuff(axis evaluation.Axis, options ...func(buff *EvaluationBuff)) *EvaluationBuff {
	buff := &EvaluationBuff{
		axis:       axis,
		bias:       0,
		multiplier: 100,
	}

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
		buff.FiniteModifier = modifier.NewFiniteModifier(capacity)
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
		axis:       b.axis,
		bias:       b.bias,
		multiplier: b.multiplier,

		FiniteModifier: b.FiniteModifier.Clone().(*modifier.FiniteModifier),
	}
}

type ClearingBuff struct {
	axis       evaluation.Axis
	bias       int
	multiplier int

	action *Action
}

func NewClearingBuff(axis evaluation.Axis, action *Action, options ...func(buff *ClearingBuff)) *ClearingBuff {
	buff := &ClearingBuff{
		axis:       axis,
		bias:       0,
		multiplier: 100,

		action: action,
	}

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

func (b *ClearingBuff) Fork(_ *evaluation.Block, signal Signal) Reactor {
	action := b.action
	if signal != nil {
		action = signal.(*PreActionSignal).Action
	}

	return &ClearingBuff{
		axis:       b.axis,
		bias:       b.bias,
		multiplier: b.multiplier,
		action:     action,
	}
}

type TaggedBuff struct {
	*modifier.TaggedModifier
	*modifier.FiniteModifier
	buffs []*EvaluationBuff
}

func NewTaggedBuff(tag any, buffs []*EvaluationBuff, options ...func(buff *TaggedBuff)) *TaggedBuff {
	return &TaggedBuff{
		TaggedModifier: modifier.NewTaggedModifier(tag),
		buffs:          buffs,
	}
}

func TaggedCapacity(capacity int) func(buff *TaggedBuff) {
	return func(buff *TaggedBuff) {
		buff.FiniteModifier = modifier.NewFiniteModifier(capacity)
	}
}

func (b *TaggedBuff) React(signal Signal) {
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

func (b *TaggedBuff) Fork(block *evaluation.Block, signal Signal) Reactor {
	return &TaggedBuff{
		TaggedModifier: b.TaggedModifier,
		FiniteModifier: b.FiniteModifier.Clone().(*modifier.FiniteModifier),
		buffs:          b.buffs,
	}
}
