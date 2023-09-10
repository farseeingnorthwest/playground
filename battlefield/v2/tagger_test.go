package battlefield_test

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUnmarshalTag(t *testing.T) {
	for _, tt := range []struct {
		name  string
		value any
	}{} {
		t.Run(tt.name, func(t *testing.T) {
			bytes, err := json.Marshal(tt.value)
			if err != nil {
				t.Fatal(err)
			}

			var tag any
			err = json.Unmarshal(bytes, tag)

			assert.NoError(t, err)
			assert.Equal(t, tt.value, tag)
		})
	}
}
