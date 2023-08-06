package battlefield

type Leading struct {
	trigger Trigger
	count   int
}

func (l *Leading) React(signal Signal) {
	if l == nil {
		return
	}

	if l.trigger.Trigger(signal) && l.count > 0 {
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

func (c *Cooling) React(signal Signal) {
	if c == nil {
		return
	}

	if c.trigger.Trigger(signal) && c.p > 0 {
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

func (c *Capacity) React(signal Signal) {
	if c == nil {
		return
	}

	if c.trigger != nil && c.trigger.Trigger(signal) {
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

func (r *Responder) React(signal Signal, warriors []Warrior) {
	if r.trigger.Trigger(signal) {
		for _, actor := range r.actors {
			actor.Act(signal, warriors)
		}
	}
}

func (r *Responder) Fork(evaluator Evaluator) any {
	actors := make([]Actor, len(r.actors))
	for i, actor := range r.actors {
		actors[i] = actor.Fork(evaluator).(Actor)
	}

	return &Responder{r.trigger, actors}
}

type FatReactor struct {
	taggers    map[any]struct{}
	leading    *Leading
	cooling    *Cooling
	capacity   *Capacity
	responders []*Responder
}

func NewFatReactor(options ...func(*FatReactor)) *FatReactor {
	r := &FatReactor{taggers: make(map[any]struct{})}
	for _, option := range options {
		option(r)
	}

	return r
}

func FatTag(tags ...any) func(*FatReactor) {
	return func(r *FatReactor) {
		for _, tag := range tags {
			r.taggers[tag] = struct{}{}
		}
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

func (r *FatReactor) React(signal Signal, warriors []Warrior) {
	r.leading.React(signal)
	r.cooling.React(signal)
	r.capacity.React(signal)

	if !r.leading.Ready() || !r.cooling.Ready() || !r.capacity.Ready() {
		return
	}

	if scripter, ok := signal.(Scripter); ok {
		scripter.New(signal.Current(), r)
	}
	for _, responder := range r.responders {
		responder.React(signal, warriors)
	}

	r.cooling.WarmUp()
	r.capacity.Flush()
}

func (r *FatReactor) Tags() (tags []any) {
	for tag := range r.taggers {
		tags = append(tags, tag)
	}

	return
}

func (r *FatReactor) Match(tags ...any) bool {
	for _, tag := range tags {
		if _, ok := r.taggers[tag]; !ok {
			return false
		}
	}

	return true
}

func (r *FatReactor) Fork(evaluator Evaluator) any {
	responders := make([]*Responder, len(r.responders))
	for i, responder := range r.responders {
		responders[i] = responder.Fork(evaluator).(*Responder)
	}

	return &FatReactor{
		r.taggers,
		r.leading.Fork(),
		r.cooling.Fork(),
		r.capacity.Fork(),
		responders,
	}
}
