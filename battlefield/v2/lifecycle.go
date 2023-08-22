package battlefield

import (
	"log/slog"

	"github.com/farseeingnorthwest/playground/battlefield/v2/functional"
)

type Lifecycle struct {
	Leading  functional.Option[int]
	Cooling  functional.Option[Ratio]
	Capacity functional.Option[int]
}

func (c *Lifecycle) SetLeading(count int) {
	c.Leading = functional.Some(count)
}

func (c *Lifecycle) SetCooling(current int, maximum int) {
	c.Cooling = functional.Some(Ratio{current, maximum})
}

func (c *Lifecycle) SetCapacity(count int) {
	c.Capacity = functional.Some(count)
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
