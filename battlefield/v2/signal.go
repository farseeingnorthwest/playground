package battlefield

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

type ForkableSignal interface {
	Signal
	Fork(any) Signal
}

type Sourcer interface {
	Source() (any, Reactor)
}

type Renderer interface {
	Render(*BattleField)
}

type Scripter interface {
	Push(any, Reactor)
	Pop()
	Add(Action)
}

type myScripter struct {
	scripts []Script
}

func (s *myScripter) Push(any any, reactor Reactor) {
	s.scripts = append(s.scripts, NewMyScript(any, reactor))
}

func (s *myScripter) Pop() {
	s.scripts = s.scripts[:len(s.scripts)-1]
}

func (s *myScripter) Add(action Action) {
	s.scripts[len(s.scripts)-1].Add(action)
}

func (s *myScripter) Render(b *BattleField) {
	for _, script := range s.scripts {
		script.Render(b)
	}
}

type BattleStartSignal struct {
	current any
	myScripter
}

func NewBattleStartSignal() *BattleStartSignal {
	return &BattleStartSignal{nil, myScripter{}}
}

func (s *BattleStartSignal) Current() any {
	return s.current
}

func (s *BattleStartSignal) Fork(current any) Signal {
	return &BattleStartSignal{current, myScripter{}}
}

type LaunchSignal struct {
	TagSet
	current Warrior
	myScripter
}

func NewLaunchSignal(current Warrior) *LaunchSignal {
	return &LaunchSignal{NewTagSet(), current, myScripter{}}
}

func (s *LaunchSignal) Current() any {
	return s.current
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

type LossSignal struct {
	current any
	target  Warrior
	myScripter
}

func NewLossSignal(target Warrior) *LossSignal {
	return &LossSignal{nil, target, myScripter{}}
}

func (s *LossSignal) Current() any {
	return s.current
}

func (s *LossSignal) Target() Warrior {
	return s.target
}

func (s *LossSignal) Fork(current any) Signal {
	return &LossSignal{current, s.target, myScripter{}}
}

type ActionSignal interface {
	Action() Action
}

type PreActionSignal struct {
	current any
	action  Action
	myScripter
}

func NewPreActionSignal(action Action) *PreActionSignal {
	return &PreActionSignal{nil, action, myScripter{}}
}

func (s *PreActionSignal) Current() any {
	return s.current
}

func (s *PreActionSignal) Action() Action {
	return s.action
}

func (s *PreActionSignal) Fork(current any) Signal {
	return &PreActionSignal{current, s.action, myScripter{}}
}

type PostActionSignal struct {
	current any
	action  Action
	myScripter
}

func NewPostActionSignal(action Action) *PostActionSignal {
	return &PostActionSignal{nil, action, myScripter{}}
}

func (s *PostActionSignal) Current() any {
	return s.current
}

func (s *PostActionSignal) Action() Action {
	return s.action
}

func (s *PostActionSignal) Fork(current any) Signal {
	return &PostActionSignal{current, s.action, myScripter{}}
}

type RoundStartSignal struct {
	current any
	myScripter
}

func NewRoundStartSignal() *RoundStartSignal {
	return &RoundStartSignal{nil, myScripter{}}
}

func (s *RoundStartSignal) Current() any {
	return s.current
}

func (s *RoundStartSignal) Fork(current any) Signal {
	return &RoundStartSignal{current, myScripter{}}
}

type RoundEndSignal struct {
	current any
	myScripter
}

func NewRoundEndSignal() *RoundEndSignal {
	return &RoundEndSignal{nil, myScripter{}}
}

func (s *RoundEndSignal) Current() any {
	return s.current
}

func (s *RoundEndSignal) Fork(current any) Signal {
	return &RoundEndSignal{current, myScripter{}}
}

type EvaluationSignal struct {
	current any
	axis    Axis
	value   int
}

func NewEvaluationSignal(current any, axis Axis, value int) *EvaluationSignal {
	return &EvaluationSignal{current, axis, value}
}

func (s *EvaluationSignal) Current() any {
	return s.current
}

func (s *EvaluationSignal) Axis() Axis {
	return s.axis
}

func (s *EvaluationSignal) Value() int {
	return s.value
}

func (s *EvaluationSignal) SetValue(value int) {
	s.value = value
}

type LifecycleSignal struct {
	current   any
	scripter  any
	reactor   Reactor
	lifecycle *Lifecycle
}

func NewTempoSignal(scripter any, reactor Reactor, lifecycle *Lifecycle) *LifecycleSignal {
	return &LifecycleSignal{nil, scripter, reactor, lifecycle}
}

func (s *LifecycleSignal) Current() any {
	return s.current
}

func (s *LifecycleSignal) Source() (any, Reactor) {
	return s.scripter, s.reactor
}

func (s *LifecycleSignal) Lifecycle() *Lifecycle {
	return s.lifecycle
}

func (s *LifecycleSignal) Fork(current any) Signal {
	return &LifecycleSignal{current, s.scripter, s.reactor, s.lifecycle}
}
