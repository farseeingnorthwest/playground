package battlefield_test

import (
	"encoding/json"
	"testing"

	. "github.com/farseeingnorthwest/playground/battlefield/v2"
	"github.com/stretchr/testify/assert"
)

func TestVerbFile_UnmarshalJSON(t *testing.T) {
	for _, tt := range []struct {
		name string
		verb Verb
	}{
		{"attack", NewAttack(nil, false)},
		{"heal", NewHeal(NewMultiplier(80, LossEvaluator{}))},
		{"buff", NewBuff(false, nil, Effect["Sleep"])},
		{"purge", NewPurge(DefaultRng, Label("Buff"), 0)},
	} {
		t.Run(tt.name, func(t *testing.T) {
			buf, err := json.Marshal(tt.verb)
			if err != nil {
				t.Fatal(err)
			}

			var f VerbFile
			err = json.Unmarshal(buf, &f)

			assert.NoError(t, err)
			assert.Equal(t, tt.verb, f.Verb)
		})
	}
}
