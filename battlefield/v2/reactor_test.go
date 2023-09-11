package battlefield_test

import (
	"encoding/json"
	"testing"

	. "github.com/farseeingnorthwest/playground/battlefield/v2"
	"github.com/stretchr/testify/assert"
)

func TestActorFile_UnmarshalJSON(t *testing.T) {
	for _, tt := range []struct {
		name  string
		actor Actor
	}{
		{"buffer", NewBuffer(Damage, false, ConstEvaluator(110))},
		{"verb", NewVerbActor(NewAttack(nil, false), AxisEvaluator(Damage))},
		{"select", NewSelectActor(NewVerbActor(NewAttack(nil, false), AxisEvaluator(Damage)), SideSelector(false))},
		{"probability", NewProbabilityActor(
			DefaultRng,
			AxisEvaluator(CriticalOdds),
			NewVerbActor(NewAttack(nil, false), AxisEvaluator(Damage)),
		)},
		{"sequence", NewSequenceActor(
			NewVerbActor(NewAttack(nil, false), AxisEvaluator(Damage)),
			NewVerbActor(NewAttack(nil, false), AxisEvaluator(Damage)),
		)},
		{"repeat", NewRepeatActor(3, NewVerbActor(NewAttack(nil, false), AxisEvaluator(Damage)))},
		{"critical", CriticalActor{}},
		{"immune", ImmuneActor{}},
		{"loss stopper", NewLossStopper(NewMultiplier(10, AxisEvaluator(HealthMaximum)), true)},
		{"loss resister", LossResister{}},
		{"theory", ElementTheory},
	} {
		t.Run(tt.name, func(t *testing.T) {
			bytes, err := json.Marshal(tt.actor)
			if err != nil {
				t.Fatal(err)
			}

			var f ActorFile[Element]
			err = json.Unmarshal(bytes, &f)

			assert.NoError(t, err)
			assert.Equal(t, tt.actor, f.Actor)
		})
	}
}
