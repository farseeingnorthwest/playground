package battlefield

import (
	"encoding/json"
	"errors"
	"reflect"
	"sort"
)

var (
	_         Tagger = TagSet{}
	ErrBadTag        = errors.New("bad tag")
)

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

func (t TagSet) MarshalJSON() ([]byte, error) {
	tags := t.Tags()
	sort.Sort(byType(tags))

	return json.Marshal(tags)
}

type byType []any

func (b byType) Len() int {
	return len(b)
}

func (b byType) Less(i, j int) bool {
	return reflect.TypeOf(b[i]).Name() < reflect.TypeOf(b[j]).Name()
}

func (b byType) Swap(i, j int) {
	b[i], b[j] = b[j], b[i]
}

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

type Label string

func (l Label) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]string{"label": string(l)})
}

type Priority int

func (p Priority) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]int{"priority": int(p)})
}

type ExclusionGroup uint8

func (g ExclusionGroup) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]uint8{
		"group": uint8(g),
	})
}

func UnmarshalTag(data []byte) (any, error) {
	var tag map[string]any
	if err := json.Unmarshal(data, &tag); err != nil {
		return nil, err
	}

	if label, ok := tag["label"]; ok {
		if label, ok := label.(string); ok {
			return Label(label), nil
		}

		return nil, ErrBadTag
	}
	if priority, ok := tag["priority"]; ok {
		if priority, ok := priority.(int); ok {
			return Priority(priority), nil
		}

		return nil, ErrBadTag
	}
	if exclusionGroup, ok := tag["group"]; ok {
		if exclusionGroup, ok := exclusionGroup.(uint8); ok {
			return ExclusionGroup(exclusionGroup), nil
		}

		return nil, ErrBadTag
	}
	if stackingLimit, ok := tag["stacking"]; ok {
		if stackingLimit, ok := stackingLimit.(int); ok {
			return NewStackingLimit(stackingLimit), nil
		}

		return nil, ErrBadTag
	}

	return nil, ErrBadTag
}
