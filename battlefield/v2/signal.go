package battlefield

const (
	LifecycleTrigger LifecycleAffairs = 1 << iota
	LifecycleOverflow
)

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
	ID() int
	Current() any
}

type FreeSignal struct {
	id      int
	current any
}

func NewFreeSignal(id int, current any) *FreeSignal {
	return &FreeSignal{id, current}
}

func (s *FreeSignal) ID() int {
	return s.id
}

func (s *FreeSignal) Current() any {
	return s.current
}

type EvaluationSignal struct {
	id      int
	current any
	axis    Axis
	value   float64
}

func NewEvaluationSignal(id int, current any, axis Axis, value int) *EvaluationSignal {
	return &EvaluationSignal{id, current, axis, float64(value)}
}

func (s *EvaluationSignal) ID() int {
	return s.id
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
	id      int
	current Warrior
	action  Action
	loss    int
}

func NewPreLossSignal(id int, current Warrior, action Action, loss int) *PreLossSignal {
	return &PreLossSignal{id, current, action, loss}
}

func (s *PreLossSignal) ID() int {
	return s.id
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
	id int
	TagSet
	current Warrior
	scripter
}

func NewLaunchSignal(id int, current Warrior) *LaunchSignal {
	return &LaunchSignal{id, NewTagSet(), current, scripter{}}
}

func (s *LaunchSignal) ID() int {
	return s.id
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
	id      int
	current any
	scripter
}

func NewBattleStartSignal(id int) *BattleStartSignal {
	return &BattleStartSignal{id, nil, scripter{}}
}

func (s *BattleStartSignal) ID() int {
	return s.id
}

func (s *BattleStartSignal) Current() any {
	return s.current
}

func (s *BattleStartSignal) SetCurrent(current any) Signal {
	return &BattleStartSignal{s.id, current, scripter{}}
}

type RoundStartSignal struct {
	id      int
	current any
	scripter
}

func NewRoundStartSignal(id int) *RoundStartSignal {
	return &RoundStartSignal{id, nil, scripter{}}
}

func (s *RoundStartSignal) ID() int {
	return s.id
}

func (s *RoundStartSignal) Current() any {
	return s.current
}

func (s *RoundStartSignal) SetCurrent(current any) Signal {
	return &RoundStartSignal{s.id, current, scripter{}}
}

type RoundEndSignal struct {
	id      int
	current any
	scripter
}

func NewRoundEndSignal(id int) *RoundEndSignal {
	return &RoundEndSignal{id, nil, scripter{}}
}

func (s *RoundEndSignal) ID() int {
	return s.id
}

func (s *RoundEndSignal) Current() any {
	return s.current
}

func (s *RoundEndSignal) SetCurrent(current any) Signal {
	return &RoundEndSignal{s.id, current, scripter{}}
}

type ActionSignal interface {
	ScriptSignal
	Action() Action
}

type PreActionSignal struct {
	id      int
	current any
	action  Action
	scripter
}

func NewPreActionSignal(id int, action Action) *PreActionSignal {
	return &PreActionSignal{id, nil, action, scripter{}}
}

func (s *PreActionSignal) ID() int {
	return s.id
}

func (s *PreActionSignal) Current() any {
	return s.current
}

func (s *PreActionSignal) Action() Action {
	return s.action
}

func (s *PreActionSignal) SetCurrent(current any) Signal {
	return &PreActionSignal{s.id, current, s.action, scripter{}}
}

type PostActionSignal struct {
	id      int
	current any
	action  Action
	deaths  []Warrior
	scripter
}

func NewPostActionSignal(id int, action Action, deaths []Warrior) *PostActionSignal {
	return &PostActionSignal{id, nil, action, deaths, scripter{}}
}

func (s *PostActionSignal) ID() int {
	return s.id
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
	return &PostActionSignal{s.id, current, s.action, s.deaths, scripter{}}
}

type LifecycleSignal struct {
	id        int
	current   any
	parent    Signal
	scripter  any
	reactor   Reactor
	lifecycle *Lifecycle // nil if stacking limit overflow
	affairs   LifecycleAffairs
}

type LifecycleAffairs int

func NewLifecycleSignal(id int, parent Signal, scripter any, reactor Reactor, lifecycle *Lifecycle, affairs LifecycleAffairs) *LifecycleSignal {
	return &LifecycleSignal{id, nil, parent, scripter, reactor, lifecycle, affairs}
}

func (s *LifecycleSignal) ID() int {
	return s.id
}

func (s *LifecycleSignal) Current() any {
	return s.current
}

func (s *LifecycleSignal) SetCurrent(current any) Signal {
	return &LifecycleSignal{s.id, current, s.parent, s.scripter, s.reactor, s.lifecycle, s.affairs}
}

func (s *LifecycleSignal) Parent() Signal {
	return s.parent
}

func (s *LifecycleSignal) Source() (any, Reactor) {
	return s.scripter, s.reactor
}

func (s *LifecycleSignal) Lifecycle() *Lifecycle {
	return s.lifecycle
}

func (s *LifecycleSignal) Affairs() LifecycleAffairs {
	return s.affairs
}
