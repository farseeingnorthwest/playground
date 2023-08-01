package modifier

type Tagged interface {
	Tag() any
}

type TaggedModifier struct {
	tag any
}

func NewTaggedModifier(tag any) *TaggedModifier {
	return &TaggedModifier{
		tag: tag,
	}
}

func (m *TaggedModifier) Tag() any {
	return m.tag
}

func (m *TaggedModifier) Clone() any {
	return &TaggedModifier{
		tag: m.tag,
	}
}
