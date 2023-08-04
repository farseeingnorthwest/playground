package battlefield

type Signal interface {
	Current() any
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

type scripter struct {
	scripts []Script
}

func (s *scripter) Render(b *BattleField) {
	for _, script := range s.scripts {
		script.Render(b)
	}
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

func (s *BattleStartSignal) Fork(current any) Signal {
	return &BattleStartSignal{current, scripter{}}
}

type LaunchSignal struct {
	current Warrior
	scripter
}

func NewLaunchSignal(current Warrior) *LaunchSignal {
	return &LaunchSignal{current, scripter{}}
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
	scripter
}

func NewLossSignal(target Warrior) *LossSignal {
	return &LossSignal{nil, target, scripter{}}
}

func (s *LossSignal) Current() any {
	return s.current
}

func (s *LossSignal) Target() Warrior {
	return s.target
}

func (s *LossSignal) Fork(current any) Signal {
	return &LossSignal{current, s.target, scripter{}}
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

func (s *PreActionSignal) Fork(current any) Signal {
	return &PreActionSignal{current, s.action, scripter{}}
}

type PostActionSignal struct {
	current any
	action  Action
	scripter
}

func NewPostActionSignal(action Action) *PostActionSignal {
	return &PostActionSignal{nil, action, scripter{}}
}

func (s *PostActionSignal) Current() any {
	return s.current
}

func (s *PostActionSignal) Action() Action {
	return s.action
}

func (s *PostActionSignal) Fork(current any) Signal {
	return &PostActionSignal{current, s.action, scripter{}}
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

func (s *RoundStartSignal) Fork(current any) Signal {
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

func (s *RoundEndSignal) Fork(current any) Signal {
	return &RoundEndSignal{current, scripter{}}
}
