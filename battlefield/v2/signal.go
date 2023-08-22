package battlefield

var (
	_ Signal        = (*FreeSignal)(nil)
	_ Signal        = (*EvaluationSignal)(nil)
	_ Signal        = (*PreLossSignal)(nil)
	_ LocalSignal   = (*LaunchSignal)(nil)
	_ ScriptSignal  = (*BattleStartSignal)(nil)
	_ ScriptSignal  = (*RoundStartSignal)(nil)
	_ ScriptSignal  = (*RoundEndSignal)(nil)
	_ ActionSignal  = (*PreActionSignal)(nil)
	_ ActionSignal  = (*PostActionSignal)(nil)
	_ RegularSignal = (*LifecycleSignal)(nil)
)

type Signal interface {
	Current() any
}

type FreeSignal struct {
	current any
}

func NewFreeSignal(current any) *FreeSignal {
	return &FreeSignal{current}
}

func (s *FreeSignal) Current() any {
	return s.current
}

type EvaluationSignal struct {
	current any
	axis    Axis
	value   float64
}

func NewEvaluationSignal(current any, axis Axis, value int) *EvaluationSignal {
	return &EvaluationSignal{current, axis, float64(value)}
}

func (s *EvaluationSignal) Current() any {
	return s.current
}

func (s *EvaluationSignal) Axis() Axis {
	return s.axis
}

func (s *EvaluationSignal) Value() int {
	return int(s.value)
}

func (s *EvaluationSignal) Amend(f func(float64) float64) {
	s.value = f(s.value)
}

type PreLossSignal struct {
	current Warrior
	loss    int
}

func NewPreLossSignal(current Warrior, loss int) *PreLossSignal {
	return &PreLossSignal{current, loss}
}

func (s *PreLossSignal) Current() any {
	return s.current
}

func (s *PreLossSignal) Loss() int {
	return s.loss
}

func (s *PreLossSignal) SetLoss(loss int) {
	s.loss = loss
}

type LocalSignal interface {
	Signal
	Scripter
}

type LaunchSignal struct {
	TagSet
	current Warrior
	scripter
}

func NewLaunchSignal(current Warrior) *LaunchSignal {
	return &LaunchSignal{NewTagSet(), current, scripter{}}
}

func (s *LaunchSignal) Current() any {
	return s.current
}

type RegularSignal interface {
	Signal
	SetCurrent(any) Signal
}

type ScriptSignal interface {
	RegularSignal
	Scripter
}

type BattleStartSignal struct {
	current any
	scripter
}

func NewBattleStartSignal() *BattleStartSignal {
	return &BattleStartSignal{nil, scripter{}}
}

func (s *BattleStartSignal) Current() any {
	return s.current
}

func (s *BattleStartSignal) SetCurrent(current any) Signal {
	return &BattleStartSignal{current, scripter{}}
}

type RoundStartSignal struct {
	current any
	scripter
}

func NewRoundStartSignal() *RoundStartSignal {
	return &RoundStartSignal{nil, scripter{}}
}

func (s *RoundStartSignal) Current() any {
	return s.current
}

func (s *RoundStartSignal) SetCurrent(current any) Signal {
	return &RoundStartSignal{current, scripter{}}
}

type RoundEndSignal struct {
	current any
	scripter
}

func NewRoundEndSignal() *RoundEndSignal {
	return &RoundEndSignal{nil, scripter{}}
}

func (s *RoundEndSignal) Current() any {
	return s.current
}

func (s *RoundEndSignal) SetCurrent(current any) Signal {
	return &RoundEndSignal{current, scripter{}}
}

type ActionSignal interface {
	ScriptSignal
	Action() Action
}

type PreActionSignal struct {
	current any
	action  Action
	scripter
}

func NewPreActionSignal(action Action) *PreActionSignal {
	return &PreActionSignal{nil, action, scripter{}}
}

func (s *PreActionSignal) Current() any {
	return s.current
}

func (s *PreActionSignal) Action() Action {
	return s.action
}

func (s *PreActionSignal) SetCurrent(current any) Signal {
	return &PreActionSignal{current, s.action, scripter{}}
}

type PostActionSignal struct {
	current any
	action  Action
	deaths  []Warrior
	scripter
}

func NewPostActionSignal(action Action, deaths []Warrior) *PostActionSignal {
	return &PostActionSignal{nil, action, deaths, scripter{}}
}

func (s *PostActionSignal) Current() any {
	return s.current
}

func (s *PostActionSignal) Action() Action {
	return s.action
}

func (s *PostActionSignal) Deaths() []Warrior {
	return s.deaths
}

func (s *PostActionSignal) SetCurrent(current any) Signal {
	return &PostActionSignal{current, s.action, s.deaths, scripter{}}
}

type LifecycleSignal struct {
	current   any
	scripter  any
	reactor   Reactor
	lifecycle *Lifecycle // nil if stacking limit overflow
}

func NewLifecycleSignal(scripter any, reactor Reactor, lifecycle *Lifecycle) *LifecycleSignal {
	return &LifecycleSignal{nil, scripter, reactor, lifecycle}
}

func (s *LifecycleSignal) Current() any {
	return s.current
}

func (s *LifecycleSignal) SetCurrent(current any) Signal {
	return &LifecycleSignal{current, s.scripter, s.reactor, s.lifecycle}
}

func (s *LifecycleSignal) Source() (any, Reactor) {
	return s.scripter, s.reactor
}

func (s *LifecycleSignal) Lifecycle() *Lifecycle {
	return s.lifecycle
}
