package battlefield

type Health struct {
	Current int
	Maximum int
}

type Interest interface {
	Target() *Warrior
}

type interest struct {
	target *Warrior
}

func (i *interest) Target() *Warrior {
	return i.target
}

type AttackInterest struct {
	interest

	Damage   int
	Defense  int
	Loss     int
	Overflow int
	Health
}

type HealingInterest struct {
	interest

	Healing  int
	Overflow int
	Health
}

type BuffingInterest struct {
	interest

	Buff Reactor
}

type PurgingInterest struct {
	interest
}
