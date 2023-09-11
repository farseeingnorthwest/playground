package battlefield_test

import (
	"encoding/json"
	"testing"

	. "github.com/farseeingnorthwest/playground/battlefield/v2"
	"github.com/stretchr/testify/assert"
)

func TestTagFile_UnmarshalJSON(t *testing.T) {
	for _, tt := range []struct {
		name  string
		value any
	}{
		{"label", Label("foo")},
		{"priority", Priority(10)},
		{"exclusion group", ExclusionGroup(1)},
		{"stacking limit", NewStackingLimit(1)},
	} {
		t.Run(tt.name, func(t *testing.T) {
			bytes, err := json.Marshal(tt.value)
			if err != nil {
				t.Fatal(err)
			}

			var f TagFile
			err = json.Unmarshal(bytes, &f)

			assert.NoError(t, err)
			assert.Equal(t, tt.value, f.Tag)
		})
	}
}
