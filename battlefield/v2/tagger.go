package battlefield

import "reflect"

type Matcher interface {
	Match(any) bool
}

type TypeMatcher struct {
	typ any
}

func NewTypeMatcher(proto any) *TypeMatcher {
	return &TypeMatcher{reflect.TypeOf(proto)}
}

func (m *TypeMatcher) Match(any any) bool {
	return m.typ == reflect.TypeOf(any)
}

type Tagger interface {
	Tags() []any
	Match(...any) bool
	Find(Matcher) any
	Save(any)
}

type TagSet map[any]struct{}

func NewTagSet(tags ...any) TagSet {
	t := TagSet(make(map[any]struct{}))
	for _, tag := range tags {
		t.Save(tag)
	}

	return t
}

func (t TagSet) Tags() (tags []any) {
	for tag := range t {
		tags = append(tags, tag)
	}

	return
}

func (t TagSet) Match(tags ...any) bool {
	for _, tag := range tags {
		if _, ok := t[tag]; !ok {
			return false
		}
	}

	return true
}

func (t TagSet) Find(matcher Matcher) any {
	for tag := range t {
		if matcher.Match(tag) {
			return tag
		}
	}

	return nil
}

func (t TagSet) Save(tag any) {
	t[tag] = struct{}{}
}

func QueryTag[T any](a any) (T, bool) {
	proto := new(T)
	tagger, ok := a.(Tagger)
	if !ok {
		return *proto, false
	}

	tag := tagger.Find(NewTypeMatcher(*proto))
	if tag == nil {
		return *proto, false
	}

	return tag.(T), true
}

func QueryTagA[T any](a any) any {
	if tag, ok := QueryTag[T](a); ok {
		return tag
	}

	return nil
}

func First[A any, B any](a A, _ B) A {
	return a
}

func Second[A any, B any](a A, b B) B {
	return b
}
