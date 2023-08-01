package evaluation

var (
	Head = &HeadEvaluator{}
)

type vector interface {
	Component(Axis) int
}

type Evaluator interface {
	Evaluate(vector, *Block) int
}

type HeadEvaluator struct {
}

func (e *HeadEvaluator) Evaluate(_ vector, block *Block) int {
	return block.Value()
}

type AxisEvaluator struct {
	Axis
	percentage int
}

func NewAxisEvaluator(axis Axis, percentage int) *AxisEvaluator {
	return &AxisEvaluator{axis, percentage}
}

func (e *AxisEvaluator) Evaluate(v vector, _ *Block) int {
	return v.Component(e.Axis) * e.percentage / 100
}

type Bundle struct {
	e Evaluator
	b *Block
}

func NewBundleProto(e Evaluator) *Bundle {
	return &Bundle{e, nil}
}

func (b *Bundle) Block() *Block {
	return b.b
}

func (b *Bundle) Evaluate(v vector) int {
	if b.e == nil {
		panic("evaluator is not set")
	}

	return b.e.Evaluate(v, b.b)
}

func (b *Bundle) Fork(block *Block) *Bundle {
	return &Bundle{b.e, block}
}

func (b *Bundle) ForkWith(v vector) *Block {
	if b.e == nil {
		return b.b
	}

	return b.b.Fork(b.Evaluate(v))
}
