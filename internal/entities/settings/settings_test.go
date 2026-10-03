package settings

import (
	"math"
	"testing"
)

func TestNormalize(t *testing.T) {
	scale := func(x, y, z float64) CardSize { return CardSize{ScaleX: x, ScaleY: y, ScaleZ: z} }
	tests := []struct {
		name string
		in   Settings
		want Settings
	}{
		{"valid_kept", Settings{Lang: "ru", EnableBackShadow: true, CardSize: scale(1.25, 1.5, 1.75)},
			Settings{Lang: "ru", EnableBackShadow: true, CardSize: scale(1.25, 1.5, 1.75)}},
		{"all_missing", Settings{}, Default()},
		{"unknown_lang", Settings{Lang: "fr", CardSize: scale(1, 1, 1)}, Settings{Lang: "en", CardSize: scale(1, 1, 1)}},
		{"zero_scale", Settings{Lang: "en", CardSize: scale(0, 0, 0)}, Default()},
		{"negative_scale", Settings{Lang: "en", CardSize: scale(-2, 1, 1)}, Default()},
		{"nan_and_inf", Settings{Lang: "en", CardSize: scale(math.NaN(), math.Inf(1), 1)}, Default()},
		{"one_bad_axis", Settings{Lang: "en", CardSize: scale(2, 0, 3)}, Settings{Lang: "en", CardSize: scale(2, 1, 3)}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.in.Normalize(); got != tt.want {
				t.Fatalf("got %+v, want %+v", got, tt.want)
			}
		})
	}
}
