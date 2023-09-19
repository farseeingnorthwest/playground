package battlefield_test

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"

	. "github.com/farseeingnorthwest/playground/battlefield/v2"
	. "github.com/farseeingnorthwest/playground/battlefield/v2/examples"
)

func TestFatReactor_UnmarshalJSON(t *testing.T) {
	var rng Rng
	rng, DefaultRng = DefaultRng, RngX
	defer (func() {
		DefaultRng = rng
	})()

	regress := func(v any) func(*testing.T) {
		return func(t *testing.T) {
			data, err := json.Marshal(v)
			if err != nil {
				t.Fatal(err)
			}

			var f FatReactorFile
			if err := json.Unmarshal(data, &f); err != nil {
				t.Fatal(err)
			}

			opts := cmp.Options{
				cmp.Exporter(func(reflect.Type) bool {
					return true
				}),
			}
			if !cmp.Equal(v, f.FatReactor, opts...) {
				t.Error(cmp.Diff(v, f.FatReactor, opts...))
			}
		}
	}

	for i, v := range Regular {
		t.Run(fmt.Sprintf("regular #%d", i), regress(v))
	}

	for k, v := range Effect {
		t.Run(k, regress(v))
	}

	for i, v := range Special {
		for j, vv := range v {
			t.Run(fmt.Sprintf("special #%d-%d", i, j), regress(vv))
		}
	}
}

func init() {
	RegisterTagType("element", Element(0))
}
