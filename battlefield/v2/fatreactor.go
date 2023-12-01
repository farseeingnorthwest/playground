package battlefield

import (
	"encoding/json"

	"github.com/farseeingnorthwest/playground/battlefield/v2/functional"
)

var (
	_ ForkReactor = (*FatReactor)(nil)
)

type FatReactor struct {
	TagSet
	leading   *Leading
	cooling   *Cooling
	capacity  *Capacity
	responder *Responder
}

func NewFatReactor(options ...func(*FatReactor)) *FatReactor {
	r := &FatReactor{TagSet: NewTagSet()}
	for _, option := range options {
		option(r)
	}

	return r
}

func NewBuffReactor(axis Axis, bias bool, evaluator Evaluator, options ...func(*FatReactor)) *FatReactor {
	options = append(
		options,
		FatRespond(
			NewSignalTrigger(&EvaluationSignal{}),
			NewBuffer(axis, bias, evaluator),
		),
	)

	return NewFatReactor(options...)
}

func FatTags(tags ...any) func(*FatReactor) {
	return func(r *FatReactor) {
		r.TagSet = NewTagSet(tags...)
	}
}

func FatRespond(trigger Trigger, actors ...Actor) func(*FatReactor) {
	return func(r *FatReactor) {
		r.responder = &Responder{trigger, NewSequenceActor(actors...)}
	}
}

func FatLeading(trigger Trigger, count int) func(*FatReactor) {
	return func(r *FatReactor) {
		r.leading = &Leading{trigger, count}
	}
}

func FatCooling(trigger Trigger, count int) func(*FatReactor) {
	return func(r *FatReactor) {
		r.cooling = &Cooling{trigger, count, 0}
	}
}

func FatCapacity(trigger Trigger, count int) func(*FatReactor) {
	return func(r *FatReactor) {
		r.capacity = &Capacity{trigger, count}
	}
}

func (r *FatReactor) Update(options ...func(*FatReactor)) {
	for _, option := range options {
		option(r)
	}
}

func (r *FatReactor) React(signal Signal, ec EvaluationContext) {
	lc := &Lifecycle{}
	r.leading.React(signal, ec, lc)
	r.cooling.React(signal, ec, lc)
	r.capacity.React(signal, ec, lc)

	triggered := false
	ac := newActorContext(ec, r.capacity.Count())
	defer func() {
		var affairs LifecycleAffairs
		if triggered {
			n := r.capacity.Count() - ac.Capacity()
			if n < 1 {
				n = 1
			}
			r.capacity.Flush(lc, n)
			r.cooling.WarmUp(lc)
			affairs |= LifecycleTrigger
		}

		lc.Flush(signal, r, affairs, ec)
	}()

	if !r.leading.Ready() || !r.cooling.Ready() || !r.capacity.Ready() {
		return
	}

	if tagger, ok := signal.(Tagger); ok {
		if g := r.Find(NewTypeMatcher(ExclusionGroup(0))); g != nil {
			if tagger.Match(g) {
				return
			}

			defer func() {
				if triggered {
					tagger.Save(g)
				}
			}()
		}
	}

	if scripter, ok := signal.(Scripter); ok {
		scripter.Push(signal, r, ac.InstructionChannel())
		defer func() {
			if !triggered {
				scripter.Pop()
			}
		}()
	}

	triggered = r.responder != nil && r.responder.React(signal, ac)
}

func (r *FatReactor) Lifecycle() Lifecycle {
	lc := &Lifecycle{}
	if r.leading != nil {
		lc.SetLeading(r.leading.count)
	}
	if r.cooling != nil {
		lc.SetCooling(r.cooling.p, r.cooling.count)
	}
	if r.capacity != nil {
		lc.SetCapacity(r.capacity.count)
	}

	return *lc
}

func (r *FatReactor) Active() bool {
	return r.capacity.Ready()
}

func (r *FatReactor) Fork(evaluator Evaluator) any {
	res := r.responder
	if res != nil {
		res = res.Fork(evaluator).(*Responder)
	}

	return &FatReactor{
		r.TagSet,
		r.leading.Fork(),
		r.cooling.Fork(),
		r.capacity.Fork(),
		res,
	}
}

func (r *FatReactor) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Tags      TagSet     `json:"tags,omitempty"`
		Leading   *Leading   `json:"leading,omitempty"`
		Cooling   *Cooling   `json:"cooling,omitempty"`
		Capacity  *Capacity  `json:"capacity,omitempty"`
		Responder *Responder `json:"respond,omitempty"`
	}{
		r.TagSet,
		r.leading,
		r.cooling,
		r.capacity,
		r.responder,
	})
}

type FatReactorFile struct {
	*FatReactor
}

func (f *FatReactorFile) UnmarshalJSON(data []byte) error {
	var fr struct {
		Tags          []TagFile
		Leading       *Leading
		Cooling       *Cooling
		Capacity      *Capacity
		ResponderFile ResponderFile `json:"respond"`
	}
	if err := json.Unmarshal(data, &fr); err != nil {
		return err
	}

	f.FatReactor = &FatReactor{
		TagSet: NewTagSet(functional.Map(func(f TagFile) any {
			return f.Tag
		})(fr.Tags)...),
		leading:   fr.Leading,
		cooling:   fr.Cooling,
		capacity:  fr.Capacity,
		responder: fr.ResponderFile.Responder,
	}
	return nil
}

type Leading struct {
	trigger Trigger
	count   int
}

func (l *Leading) React(signal Signal, ec EvaluationContext, lc *Lifecycle) {
	if l == nil {
		return
	}

	if l.trigger.Trigger(signal, ec) && l.count > 0 {
		l.count--
		lc.SetLeading(l.count)
	}
}

func (l *Leading) Ready() bool {
	return l == nil || l.count == 0
}

func (l *Leading) Fork() *Leading {
	if l == nil {
		return nil
	}

	return &Leading{l.trigger, l.count}
}

func (l *Leading) MarshalJSON() ([]byte, error) {
	return json.Marshal(lc{
		l.count,
		l.trigger,
	})
}

func (l *Leading) UnmarshalJSON(data []byte) error {
	var lc lcf
	if err := json.Unmarshal(data, &lc); err != nil {
		return err
	}

	l.count = lc.Count
	l.trigger = lc.Trigger.Trigger
	return nil
}

type Cooling struct {
	trigger Trigger
	count   int
	p       int
}

func (c *Cooling) React(signal Signal, ec EvaluationContext, lc *Lifecycle) {
	if c == nil {
		return
	}

	if c.trigger.Trigger(signal, ec) && c.p > 0 {
		c.p--
		lc.SetCooling(c.p, c.count)
	}
}

func (c *Cooling) WarmUp(lc *Lifecycle) {
	if c == nil {
		return
	}

	c.p = c.count
	lc.SetCooling(c.p, c.count)
}

func (c *Cooling) Ready() bool {
	return c == nil || c.p == 0
}

func (c *Cooling) Fork() *Cooling {
	if c == nil {
		return nil
	}

	return &Cooling{c.trigger, c.count, 0}
}

func (c *Cooling) MarshalJSON() ([]byte, error) {
	return json.Marshal(lc{
		c.count,
		c.trigger,
	})
}

func (c *Cooling) UnmarshalJSON(data []byte) error {
	var lc lcf
	if err := json.Unmarshal(data, &lc); err != nil {
		return err
	}

	c.count = lc.Count
	c.trigger = lc.Trigger.Trigger
	c.p = 0
	return nil
}

type Capacity struct {
	trigger Trigger
	count   int
}

func (c *Capacity) React(signal Signal, ec EvaluationContext, lc *Lifecycle) {
	if c == nil {
		return
	}

	if c.trigger != nil && c.trigger.Trigger(signal, ec) && c.count > 0 {
		c.count--
		lc.SetCapacity(c.count)
	}
}

func (c *Capacity) Count() int {
	if c == nil {
		return 0
	}

	return c.count
}

func (c *Capacity) Flush(lc *Lifecycle, n int) {
	if c == nil {
		return
	}

	if c.trigger == nil && c.count > 0 {
		c.count -= n
		lc.SetCapacity(c.count)
	}
}

func (c *Capacity) Ready() bool {
	return c == nil || c.count > 0
}

func (c *Capacity) Fork() *Capacity {
	if c == nil {
		return nil
	}

	return &Capacity{c.trigger, c.count}
}

func (c *Capacity) MarshalJSON() ([]byte, error) {
	return json.Marshal(lc{
		c.count,
		c.trigger,
	})
}

func (c *Capacity) UnmarshalJSON(data []byte) error {
	var lc lcf
	if err := json.Unmarshal(data, &lc); err != nil {
		return err
	}

	c.count = lc.Count
	c.trigger = lc.Trigger.Trigger
	return nil
}

type lc struct {
	Count   int `json:"count"`
	Trigger `json:"when,omitempty"`
}

type lcf struct {
	Count   int
	Trigger TriggerFile `json:"when,omitempty"`
}

type Responder struct {
	trigger Trigger
	actor   SequenceActor
}

func (r *Responder) React(signal Signal, ac ActorContext) bool {
	if !r.trigger.Trigger(signal, ac) {
		return false
	}

	go func() {
		r.actor.Act(signal, ac.Warriors(), ac)
		ac.Resolve(true)
	}()

	return ac.WaitTriggered()
}

func (r *Responder) Fork(evaluator Evaluator) any {
	return &Responder{r.trigger, r.actor.Fork(evaluator).(SequenceActor)}
}

func (r *Responder) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Trigger `json:"when"`
		Actor   `json:"then"`
	}{
		r.trigger,
		r.actor,
	})
}

type ResponderFile struct {
	*Responder
}

func (f *ResponderFile) UnmarshalJSON(data []byte) error {
	var s struct {
		TriggerFile `json:"when"`
		ActorFile   `json:"then"`
	}
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	f.Responder = &Responder{s.TriggerFile.Trigger, s.ActorFile.Actor.(SequenceActor)}
	return nil
}
