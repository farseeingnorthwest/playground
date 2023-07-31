package battlefield

type EvaluationBuff struct {
	axis  Axis
	bias  int
	slope int

	*FiniteReactor
}

func NewEvaluationBuff(axis Axis, options ...func(buff *EvaluationBuff)) *EvaluationBuff {
	buff := &EvaluationBuff{
		axis:  axis,
		bias:  0,
		slope: 100,
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

func EvaluationSlope(slope int) func(buff *EvaluationBuff) {
	return func(buff *EvaluationBuff) {
		buff.slope = slope
	}
}

func EvaluationCapacity(capacity int) func(buff *EvaluationBuff) {
	return func(buff *EvaluationBuff) {
		buff.FiniteReactor = &FiniteReactor{capacity}
	}
}

func (b *EvaluationBuff) React(signal Signal) {
	switch sig := signal.(type) {
	case *EvaluationSignal:
		if sig.Axis() != b.axis || sig.Action() != nil {
			return
		}

		sig.Map(func(points int) int {
			return points*b.slope/100 + b.bias
		})

	case *RoundEndSignal:
		b.WarmUp()
	}
}

func (b *EvaluationBuff) Fork() any {
	return &EvaluationBuff{
		axis:  b.axis,
		bias:  b.bias,
		slope: b.slope,

		FiniteReactor: b.FiniteReactor.Fork(),
	}
}

type ClearingBuff struct {
	axis   Axis
	bias   int
	slope  int
	action *Action
}

func NewClearingBuff(axis Axis, action *Action, options ...func(buff *ClearingBuff)) *ClearingBuff {
	buff := &ClearingBuff{
		axis:   axis,
		bias:   0,
		slope:  100,
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

func ClearingSlope(slope int) func(buff *ClearingBuff) {
	return func(buff *ClearingBuff) {
		buff.slope = slope
	}
}

func (b *ClearingBuff) React(signal Signal) {
	switch sig := signal.(type) {
	case *EvaluationSignal:
		if sig.Axis() != b.axis || sig.Action() != b.action {
			return
		}

		sig.Map(func(points int) int {
			return points*b.slope/100 + b.bias
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

func (b *ClearingBuff) Fork(a *Action) *ClearingBuff {
	return &ClearingBuff{
		axis:   b.axis,
		bias:   b.bias,
		slope:  b.slope,
		action: a,
	}
}

type TaggedBuff struct {
	TaggedReactor
	*FiniteReactor
	buffs []*EvaluationBuff
}

func NewTaggedBuff(tag any, buffs []*EvaluationBuff, options ...func(buff *TaggedBuff)) *TaggedBuff {
	return &TaggedBuff{
		TaggedReactor: TaggedReactor{tag},
		buffs:         buffs,
	}
}

func TaggedCapacity(capacity int) func(buff *TaggedBuff) {
	return func(buff *TaggedBuff) {
		buff.FiniteReactor = &FiniteReactor{capacity}
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

func (b *TaggedBuff) Fork() any {
	return &TaggedBuff{
		TaggedReactor: b.TaggedReactor,
		FiniteReactor: b.FiniteReactor.Fork(),
		buffs:         b.buffs,
	}
}
