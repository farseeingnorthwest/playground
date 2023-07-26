package battlefield

type Warrior interface {
	Portfolio

	Element() Element
	Attack() int
	Defense() int
	Speed() int

	Health() int
	SetHealth(int)
}
