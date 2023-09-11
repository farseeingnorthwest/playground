package battlefield_test

import (
	"encoding/json"
	"testing"

	. "github.com/farseeingnorthwest/playground/battlefield/v2"
	"github.com/stretchr/testify/assert"
)

func TestTriggerFile_UnmarshalJSON(t *testing.T) {
	for _, tt := range []struct {
		name    string
		trigger Trigger
	}{
		{"signal", NewSignalTrigger(&LaunchSignal{})},
		{"any", NewAnyTrigger(
			NewSignalTrigger(&LaunchSignal{}),
			NewSignalTrigger(&RoundEndSignal{}),
		)},
		{"fat", NewFatTrigger(
			&LaunchSignal{},
			CurrentIsTargetTrigger{},
		)},
	} {
		t.Run(tt.name, func(t *testing.T) {
			bytes, err := json.Marshal(tt.trigger)
			if err != nil {
				t.Fatal(err)
			}

			var f TriggerFile
			err = json.Unmarshal(bytes, &f)

			assert.NoError(t, err)
			assert.Equal(t, tt.trigger, f.Trigger)
		})
	}
}

func TestActionTriggerFile_UnmarshalJSON(t *testing.T) {
	for _, tt := range []struct {
		name    string
		trigger ActionTrigger
	}{
		{"current is source", CurrentIsSourceTrigger{}},
		{"current is target", CurrentIsTargetTrigger{}},
		{"source label", NewReactorTrigger(Label("foo"))},
		{"verb", NewVerbTrigger[*Attack]()},
		{"critical strike", CriticalStrikeTrigger{}},
		{"tag", NewTagTrigger(Label("foo"))},
	} {
		t.Run(tt.name, func(t *testing.T) {
			bytes, err := json.Marshal(tt.trigger)
			if err != nil {
				t.Fatal(err)
			}

			var f ActionTriggerFile
			err = json.Unmarshal(bytes, &f)

			assert.NoError(t, err)
			assert.Equal(t, tt.trigger, f.Trigger)
		})
	}
}
