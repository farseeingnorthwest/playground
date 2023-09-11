package battlefield_test

import (
	"encoding/json"
	"testing"

	. "github.com/farseeingnorthwest/playground/battlefield/v2"
	"github.com/stretchr/testify/assert"
)

func TestEvaluatorFile_UnmarshalJSON(t *testing.T) {
	for _, tt := range []struct {
		name      string
		evaluator Evaluator
	}{
		{"const", ConstEvaluator(113)},
		{"axis", AxisEvaluator(Loss)},
		{"buff counter", NewBuffCounter(Label("foo"))},
		{"loss", LossEvaluator{}},
		{"select counter", NewSelectCounter(SideSelector(false))},
		{"add", NewAdder(10, AxisEvaluator(CriticalOdds))},
		{"mul", NewMultiplier(110, AxisEvaluator(HealthMaximum))},
	} {
		t.Run(tt.name, func(t *testing.T) {
			bytes, err := json.Marshal(tt.evaluator)
			if err != nil {
				t.Fatal(err)
			}

			var f EvaluatorFile
			err = json.Unmarshal(bytes, &f)

			assert.NoError(t, err)
			assert.Equal(t, tt.evaluator, f.Evaluator)
		})
	}
}
