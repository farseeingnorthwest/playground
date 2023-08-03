package battlefield

import (
	"encoding/json"
	"testing"

	"github.com/farseeingnorthwest/playground/battlefield/v2/mod"

	"github.com/stretchr/testify/assert"
)

func TestBattlefield_Fight(t *testing.T) {
	ob := &Observer{}
	ob.SetTag("Observer")
	b := NewBattleField(
		[]*Warrior{
			NewWarrior(
				&baseline{
					TaggerMod: mod.NewTaggerMod("Alice"),
					element:   Water,
					damage:    10,
					defense:   5,
					speed:     10,
					health:    20,
				},
				&FatPortfolio{
					[]Reactor{
						NormalAttack,
						NewCriticalAttack(&mockRng{0.001}, 10, 200),
					},
				},
			),
		},
		[]*Warrior{
			NewWarrior(
				&baseline{
					TaggerMod: mod.NewTaggerMod("Bob"),
					element:   Fire,
					damage:    15,
					defense:   5,
					speed:     9,
					health:    22,
				},
				&FatPortfolio{[]Reactor{
					NormalAttack,
				}},
			),
		},
		ob,
		&Theory,
	)
	b.Fight()

	scripts, e := json.Marshal(ob.Scripts)
	assert.NoError(t, e)
	assert.JSONEq(
		t,
		`[
{"actions":[{"interests":[{"damage":10,"defense":5,"health":{"Current":10,"Maximum":22},"loss":12,"overflow":0,"target":"Bob","verb":"attack"}],"post_scripts":[],"pre_scripts":[{"actions":[{"interests":[{"buff":"元素提高伤害","target":"Bob","verb":"buff"}],"post_scripts":[],"pre_scripts":[]}],"current":null,"sources":"元素"},{"actions":[{"interests":[{"buff":"暴击提升伤害","target":"Bob","verb":"buff"}],"post_scripts":[],"pre_scripts":[]}],"current":"Alice","sources":"暴击"}]}],"current":"Alice","sources":"普通攻击"},
{"actions":[{"interests":[{"damage":15,"defense":5,"health":{"Current":12,"Maximum":20},"loss":8,"overflow":0,"target":"Alice","verb":"attack"}],"post_scripts":[],"pre_scripts":[{"actions":[{"interests":[{"buff":"元素降低伤害","target":"Alice","verb":"buff"}],"post_scripts":[],"pre_scripts":[]}],"current":null,"sources":"元素"}]}],"current":"Bob","sources":"普通攻击"},
{"actions":[{"interests":[{"damage":10,"defense":5,"health":{"Current":4,"Maximum":22},"loss":6,"overflow":0,"target":"Bob","verb":"attack"}],"post_scripts":[],"pre_scripts":[{"actions":[{"interests":[{"buff":"元素提高伤害","target":"Bob","verb":"buff"}],"post_scripts":[],"pre_scripts":[]}],"current":null,"sources":"元素"}]}],"current":"Alice","sources":"普通攻击"},
{"actions":[{"interests":[{"damage":15,"defense":5,"health":{"Current":4,"Maximum":20},"loss":8,"overflow":0,"target":"Alice","verb":"attack"}],"post_scripts":[],"pre_scripts":[{"actions":[{"interests":[{"buff":"元素降低伤害","target":"Alice","verb":"buff"}],"post_scripts":[],"pre_scripts":[]}],"current":null,"sources":"元素"}]}],"current":"Bob","sources":"普通攻击"},
{"actions":[{"interests":[{"damage":10,"defense":5,"health":{"Current":0,"Maximum":22},"loss":6,"overflow":2,"target":"Bob","verb":"attack"}],"post_scripts":[],"pre_scripts":[{"actions":[{"interests":[{"buff":"元素提高伤害","target":"Bob","verb":"buff"}],"post_scripts":[],"pre_scripts":[]}],"current":null,"sources":"元素"}]}],"current":"Alice","sources":"普通攻击"}
]`,
		string(scripts),
	)
}

type baseline struct {
	mod.TaggerMod
	element Element
	damage  int
	defense int
	speed   int
	health  int
}

func (f *baseline) Element() Element {
	return f.element
}

func (f *baseline) Damage() int {
	return f.damage
}

func (f *baseline) Defense() int {
	return f.defense
}

func (f *baseline) Health() int {
	return f.health
}

func (f *baseline) Speed() int {
	return f.speed
}

type mockRng struct {
	initial float64
}

func (r *mockRng) Gen() (f float64) {
	f = 1 - 1e-6
	if 0 <= r.initial && r.initial < 1 {
		f = r.initial
	}

	r.initial = -1
	return
}
