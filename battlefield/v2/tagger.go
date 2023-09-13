package battlefield

import (
	"encoding/json"
	"errors"
	"reflect"
)

var (
	_ Tagger = TagSet{}

	tagType   = make(map[string]reflect.Type)
	ErrBadTag = errors.New("bad tag")
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

func (t TagSet) Tags() []any {
	tags := make([]any, len(t))
	i := 0
	for tag := range t {
		tags[i] = tag
		i++
	}

	return tags
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
	return json.Marshal(t.Tags())
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

func queryTag(a any, proto any) (any, bool) {
	tagger, ok := a.(Tagger)
	if !ok {
		return proto, false
	}

	tag := tagger.Find(NewTypeMatcher(proto))
	if tag == nil {
		return proto, false
	}

	return tag, true
}

func QueryTag[T any](a any) (T, bool) {
	tag, ok := queryTag(a, *new(T))
	return tag.(T), ok
}

func QueryTagA[T any](a any) any {
	if tag, ok := QueryTag[T](a); ok {
		return tag
	}

	return nil
}

type Label string

func (l Label) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]string{
		"_kind": "label",
		"text":  string(l),
	})
}

func (l *Label) UnmarshalJSON(data []byte) error {
	var v struct{ Text string }
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	*l = Label(v.Text)
	return nil
}

type Priority int

func (p Priority) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]any{
		"_kind": "priority",
		"index": int(p),
	})
}

func (p *Priority) UnmarshalJSON(data []byte) error {
	var v struct{ Index int }
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	*p = Priority(v.Index)
	return nil
}

type ExclusionGroup uint8

func (g ExclusionGroup) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]any{
		"_kind": "exclusion_group",
		"index": uint8(g),
	})
}

func (g *ExclusionGroup) UnmarshalJSON(data []byte) error {
	var v struct{ Index uint8 }
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	*g = ExclusionGroup(v.Index)
	return nil
}

type StackingLimit struct {
	ID       string
	Capacity int
}

func NewStackingLimit(id string, capacity int) StackingLimit {
	return StackingLimit{id, capacity}
}

func (l StackingLimit) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]any{
		"_kind":    "stacking_limit",
		"id":       l.ID,
		"capacity": l.Capacity,
	})
}

type TagFile struct {
	Tag any
}

func (f *TagFile) UnmarshalJSON(data []byte) error {
	var s string
	if json.Unmarshal(data, &s) == nil {
		f.Tag = s
		return nil
	}

	var k kind
	if err := json.Unmarshal(data, &k); err != nil {
		return err
	}

	if t, ok := tagType[k.Kind]; ok {
		v := reflect.New(t)
		if err := json.Unmarshal(data, v.Interface()); err != nil {
			return err
		}

		f.Tag = v.Elem().Interface()
		return nil
	}

	return ErrBadTag
}

func RegisterTagType(kind string, proto any) {
	tagType[kind] = reflect.TypeOf(proto)
}

type kind struct {
	Kind string `json:"_kind"`
}

func init() {
	RegisterTagType("label", Label(""))
	RegisterTagType("priority", Priority(0))
	RegisterTagType("exclusion_group", ExclusionGroup(0))
	RegisterTagType("stacking_limit", StackingLimit{})
}
