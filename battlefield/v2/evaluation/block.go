package evaluation

var (
	Zero = &Block{}
)

type Block struct {
	value int
	prev  *Block
}

func (b *Block) Value() int {
	return b.value
}

func (b *Block) Prev() *Block {
	return b.prev
}

func (b *Block) Fork(value int) *Block {
	return &Block{
		value: value,
		prev:  b,
	}
}
