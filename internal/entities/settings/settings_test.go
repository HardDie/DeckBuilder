package settings

import (
	"math"
	"testing"
)

func TestNormalize(t *testing.T) {
	tests := []struct {
		name string
		in   Settings
		want Settings
	}{
		{"valid_kept", Settings{Lang: "ru", EnableBackShadow: true, CardScale: 1.25},
			Settings{Lang: "ru", EnableBackShadow: true, CardScale: 1.25}},
		{"all_missing", Settings{}, Default()},
		{"unknown_lang", Settings{Lang: "fr", CardScale: 1}, Settings{Lang: "en", CardScale: 1}},
		{"zero_scale", Settings{Lang: "en", CardScale: 0}, Default()},
		{"negative_scale", Settings{Lang: "en", CardScale: -2}, Default()},
		{"nan_scale", Settings{Lang: "en", CardScale: math.NaN()}, Default()},
		{"inf_scale", Settings{Lang: "en", CardScale: math.Inf(1)}, Settings{Lang: "en", CardScale: MaxCardScale}},
		{"too_small", Settings{Lang: "en", CardScale: 0.01}, Settings{Lang: "en", CardScale: MinCardScale}},
		{"too_big", Settings{Lang: "en", CardScale: 25}, Settings{Lang: "en", CardScale: MaxCardScale}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.in.Normalize(); got != tt.want {
				t.Fatalf("got %+v, want %+v", got, tt.want)
			}
		})
	}
}

func TestValidCardScale(t *testing.T) {
	for v, want := range map[float64]bool{
		0.1: true, 1: true, 10: true,
		0: false, 0.09: false, 10.01: false, -1: false, math.NaN(): false, math.Inf(1): false,
	} {
		if got := ValidCardScale(v); got != want {
			t.Errorf("ValidCardScale(%v) = %v, want %v", v, got, want)
		}
	}
}
