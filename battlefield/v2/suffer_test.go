package battlefield_test

import (
	"testing"

	. "github.com/farseeingnorthwest/playground/battlefield/v2"
)

func TestConstSufferer_Suffer(t *testing.T) {
	for _, test := range []struct {
		attack, defense, suffer int
	}{
		{attack: 10, defense: 10, suffer: 0},
		{attack: 10, defense: 5, suffer: 5},
		{attack: 10, defense: 15, suffer: 0},
	} {
		suffer := ConstSufferer{}.Suffer(test.attack, test.defense)
		if suffer != test.suffer {
			t.Errorf("ConstSufferer.Suffer(%d, %d) = %d, want %d", test.attack, test.defense, suffer, test.suffer)
		}
	}
}

func TestReciprocalSufferer_Suffer(t *testing.T) {
	for _, test := range []struct {
		scale, attack, defense, suffer int
	}{
		{scale: 1, attack: 10, defense: 10, suffer: 5},
		{scale: 1, attack: 10, defense: 5, suffer: 6},
		{scale: 1, attack: 10, defense: 15, suffer: 4},
		{scale: 2, attack: 10, defense: 10, suffer: 3},
		{scale: 2, attack: 10, defense: 5, suffer: 5},
	} {
		suffer := ReciprocalSufferer{Scale: test.scale}.Suffer(test.attack, test.defense)
		if suffer != test.suffer {
			t.Errorf("ReciprocalSufferer{%d}.Suffer(%d, %d) = %d, want %d", test.scale, test.attack, test.defense, suffer, test.suffer)
		}
	}
}
