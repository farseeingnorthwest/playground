package mod

type Tagger interface {
	Tag() any
	SetTag(any)
}

type TaggerMod struct {
	tag any
}

func NewTaggerMod(tag any) TaggerMod {
	return TaggerMod{tag: tag}
}

func (m *TaggerMod) Tag() any {
	return m.tag
}

func (m *TaggerMod) SetTag(tag any) {
	m.tag = tag
}
