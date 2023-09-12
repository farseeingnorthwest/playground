package battlefield

import (
	"encoding/json"
	"errors"
	"reflect"
	"sort"
)

var (
	_ Tagger = TagSet{}

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

type StackingLimit struct {
	ID       string `json:"stacking"`
	Capacity int
}

func NewStackingLimit(id string, capacity int) StackingLimit {
	return StackingLimit{id, capacity}
}

type TagFile struct {
	Tag any
}

func (f *TagFile) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		var e *json.UnmarshalTypeError
		if !errors.As(err, &e) {
			return err
		}
	} else {
		f.Tag = s
		return nil
	}

	var tag map[string]any
	if err := json.Unmarshal(data, &tag); err != nil {
		return err
	}

	if label, ok := tag["label"]; ok {
		if label, ok := label.(string); ok {
			f.Tag = Label(label)
			return nil
		}
	} else if priority, ok := tag["priority"]; ok {
		if priority, ok := priority.(float64); ok {
			f.Tag = Priority(priority)
			return nil
		}
	} else if exclusionGroup, ok := tag["group"]; ok {
		if exclusionGroup, ok := exclusionGroup.(float64); ok {
			f.Tag = ExclusionGroup(exclusionGroup)
			return nil
		}
	} else if _, ok := tag["stacking"]; ok {
		var stacking StackingLimit
		if err := json.Unmarshal(data, &stacking); err != nil {
			return err
		}

		f.Tag = stacking
		return nil
	}

	return ErrBadTag
}
