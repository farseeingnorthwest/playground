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

func (c *Lifecycle) Flush(signal Signal, reactor Reactor, affairs LifecycleAffairs, ec EvaluationContext) {
	if signal.Current() == nil {
		return
	}

	if c.Leading.Ok() || c.Cooling.Ok() || c.Capacity.Ok() || (affairs != 0 && QueryTagA[Interest](reactor) != nil) {
		slog.Debug(
			"flush",
			slog.Int("signal", signal.ID()),
			slog.Int("affairs", int(affairs)),
			slog.Group("source",
				slog.Int("position", signal.Current().(Warrior).Position()),
				slog.Any("side", signal.Current().(Warrior).Side()),
				slog.Any("reactor", QueryTagA[Label](reactor))),
			slog.Group("lifecycle",
				slog.Int("leading", c.Leading.Value()),
				slog.Any("cooling", c.Cooling.Value()),
				slog.Int("capacity", c.Capacity.UnwrapOr(-1))),
		)
		ec.React(NewLifecycleSignal(ec.Next(), signal, signal.Current(), reactor, c, affairs))
	}
}
