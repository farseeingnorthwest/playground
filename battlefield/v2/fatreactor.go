package battlefield

import "log/slog"

type Lifecycle struct {
	Leading  Option[int]
	Cooling  Option[Ratio]
	Capacity Option[int]
}

func (c *Lifecycle) SetLeading(count int) {
	c.Leading = Some(count)
}

func (c *Lifecycle) SetCooling(current int, maximum int) {
	c.Cooling = Some(Ratio{current, maximum})
}

func (c *Lifecycle) SetCapacity(count int) {
	c.Capacity = Some(count)
}

func (c *Lifecycle) Flush(current any, reactor Reactor, ec EvaluationContext) {
	if c.Leading.Ok() || c.Cooling.Ok() || c.Capacity.Ok() {
		slog.Debug(
			"flush",
			slog.Group("source",
				slog.Int("position", current.(Warrior).Position()),
				slog.Any("side", current.(Warrior).Side()),
				slog.Any("reactor", QueryTagA[Label](reactor))),
			slog.Group("lifecycle",
				slog.Int("leading", c.Leading.Value()),
				slog.Any("cooling", c.Cooling.Value()),
				slog.Int("capacity", c.Capacity.UnwrapOr(-1))),
		)
		ec.React(NewLifecycleSignal(current, reactor, c))
	}
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

type Responder struct {
	trigger Trigger
	actor   Actor
}

func (r *Responder) React(signal Signal, ec ActorContext) bool {
	if !r.trigger.Trigger(signal, ec) {
		return false
	}

	return r.actor.Act(signal, ec.Warriors(), ec)
}

func (r *Responder) Fork(evaluator Evaluator) any {
	return &Responder{r.trigger, r.actor.Fork(evaluator).(Actor)}
}

type ExclusionGroup uint8

type FatReactor struct {
	TagSet
	leading    *Leading
	cooling    *Cooling
	capacity   *Capacity
	responders []*Responder
}

func NewFatReactor(options ...func(*FatReactor)) *FatReactor {
	r := &FatReactor{}
	for _, option := range options {
		option(r)
	}

	return r
}

func FatTags(tags ...any) func(*FatReactor) {
	return func(r *FatReactor) {
		r.TagSet = NewTagSet(tags...)
	}
}

func FatRespond(trigger Trigger, actors ...Actor) func(*FatReactor) {
	return func(r *FatReactor) {
		r.responders = append(r.responders, &Responder{trigger, NewSequenceActor(actors...)})
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

type PlainActorContext struct {
	EvaluationContext
	*PlainCapacitor
}

func (r *FatReactor) React(signal Signal, ec EvaluationContext) {
	lc := &Lifecycle{}
	r.leading.React(signal, ec, lc)
	r.cooling.React(signal, ec, lc)
	r.capacity.React(signal, ec, lc)

	trigger := false
	ac := &PlainActorContext{
		ec,
		NewPlainCapacitor(r.capacity.Count()),
	}
	defer func() {
		if trigger {
			n := r.capacity.Count() - ac.Capacity()
			if n < 1 {
				n = 1
			}
			r.capacity.Flush(lc, n)
			r.cooling.WarmUp(lc)
		}

		lc.Flush(signal.Current(), r, ec)
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
				if trigger {
					tagger.Save(g)
				}
			}()
		}
	}

	if scripter, ok := signal.(Scripter); ok {
		scripter.Push(signal.Current(), r)
		defer func() {
			if !trigger {
				scripter.Pop()
			}
		}()
	}
	for _, responder := range r.responders {
		if responder.React(signal, ac) {
			trigger = true
		}
	}
}

func (r *FatReactor) Fork(evaluator Evaluator) any {
	responders := make([]*Responder, len(r.responders))
	for i, responder := range r.responders {
		responders[i] = responder.Fork(evaluator).(*Responder)
	}

	return &FatReactor{
		r.TagSet,
		r.leading.Fork(),
		r.cooling.Fork(),
		r.capacity.Fork(),
		responders,
	}
}

func (r *FatReactor) Amend(options ...func(*FatReactor)) {
	for _, option := range options {
		option(r)
	}
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
