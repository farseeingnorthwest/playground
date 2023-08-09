package battlefield

type Leading struct {
	trigger Trigger
	count   int
}

func (l *Leading) React(signal Signal, ec EvaluationContext) {
	if l == nil {
		return
	}

	if l.trigger.Trigger(signal, ec) && l.count > 0 {
		l.count--
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

func (c *Cooling) React(signal Signal, ec EvaluationContext) {
	if c == nil {
		return
	}

	if c.trigger.Trigger(signal, ec) && c.p > 0 {
		c.p--
	}
}

func (c *Cooling) WarmUp() {
	if c == nil {
		return
	}

	c.p = c.count
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

func (c *Capacity) React(signal Signal, ec EvaluationContext) {
	if c == nil {
		return
	}

	if c.trigger != nil && c.trigger.Trigger(signal, ec) {
		c.count--
	}
}

func (c *Capacity) Flush() {
	if c == nil {
		return
	}

	if c.trigger == nil {
		c.count--
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
	actors  []Actor
}

func (r *Responder) React(signal Signal, ec EvaluationContext) (trigger bool) {
	trigger = r.trigger.Trigger(signal, ec)
	if trigger {
		for _, actor := range r.actors {
			actor.Act(signal, ec.Warriors(), ec)
		}
	}

	return
}

func (r *Responder) Fork(evaluator Evaluator) any {
	actors := make([]Actor, len(r.actors))
	for i, actor := range r.actors {
		actors[i] = actor.Fork(evaluator).(Actor)
	}

	return &Responder{r.trigger, actors}
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
		r.responders = append(r.responders, &Responder{trigger, actors})
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

func (r *FatReactor) React(signal Signal, ec EvaluationContext) {
	r.leading.React(signal, ec)
	r.cooling.React(signal, ec)
	r.capacity.React(signal, ec)

	if !r.leading.Ready() || !r.cooling.Ready() || !r.capacity.Ready() {
		return
	}

	trigger := false
	defer func() {
		if trigger {
			r.capacity.Flush()
			r.cooling.WarmUp()
		}
	}()

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
		if responder.React(signal, ec) {
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
