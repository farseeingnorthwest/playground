package battlefield

import "math"

const (
	Unlimited = math.MaxInt
)

type Finite interface {
	Capacity() int
	Valid() bool
	WarmUp()
}

type FiniteReactor struct {
	capacity int
}

func (c *FiniteReactor) Capacity() int {
	if c == nil {
		return Unlimited
	}

	return c.capacity
}

func (c *FiniteReactor) Valid() bool {
	return c.Capacity() > 0
}

func (c *FiniteReactor) WarmUp() {
	if c != nil && c.capacity > 0 {
		c.capacity--
	}
}

func (c *FiniteReactor) Fork() *FiniteReactor {
	if c == nil {
		return nil
	}

	return &FiniteReactor{
		capacity: c.capacity,
	}
}

type Periodic interface {
	Period() int
	Phase() int
	Free() bool
	CoolDown()
	WarmUp()
}

type PeriodicReactor struct {
	period int
	phase  int
}

func (c *PeriodicReactor) Period() int {
	return c.period
}

func (c *PeriodicReactor) Phase() int {
	return c.phase
}

func (c *PeriodicReactor) Free() bool {
	return c.Phase() == 0
}

func (c *PeriodicReactor) CoolDown() {
	if c.phase > 0 {
		c.phase--
	}
}

func (c *PeriodicReactor) WarmUp() {
	c.phase = c.period
}

type Stackable interface {
	StackLimit() int
}

type StackableReactor struct {
	stackLimit int
}

func (c *StackableReactor) StackLimit() int {
	return c.stackLimit
}

type Tagged interface {
	Tag() any
}

type TaggedReactor struct {
	tag any
}

func (c *TaggedReactor) Tag() any {
	return c.tag
}
