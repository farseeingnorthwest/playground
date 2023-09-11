package battlefield_test

import (
	"encoding/json"
	"testing"

	. "github.com/farseeingnorthwest/playground/battlefield/v2"
	"github.com/stretchr/testify/assert"
)

func TestSelectorFile_UnmarshalJSON(t *testing.T) {
	for _, tt := range []struct {
		name     string
		selector Selector
	}{
		{"absolute side", AbsoluteSideSelector(Left)},
		{"side", SideSelector(false)},
		{"current", CurrentSelector{}},
		{"source", SourceSelector{}},
		{"sort", NewSortSelector(Damage, false)},
		{"shuffle", NewShuffleSelector(DefaultRng, Label("Taunt"))},
		{"front", FrontSelector(0)},
		{"counter position", CounterPositionSelector(0)},
		{"water level", NewWaterLevelSelector(
			Gt,
			AxisEvaluator(Health),
			0,
		)},
		{"pipeline", PipelineSelector{
			SideSelector(false),
			NewShuffleSelector(DefaultRng, Label("Taunt")),
			FrontSelector(1),
		}},
	} {
		t.Run(tt.name, func(t *testing.T) {
			bytes, err := json.Marshal(tt.selector)
			if err != nil {
				t.Fatal(err)
			}

			var f SelectorFile
			err = json.Unmarshal(bytes, &f)

			assert.NoError(t, err)
			assert.Equal(t, tt.selector, f.Selector)
		})
	}
}
