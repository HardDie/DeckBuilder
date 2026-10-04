package settings

import (
	"log/slog"
	"math"
	"testing"
)

func TestNormalize(t *testing.T) {
	tests := []struct {
		name string
		in   Settings
		want Settings
	}{
		{"valid_kept", Settings{Lang: "ru", EnableBackShadow: true, CardScale: 1.25, LogLevel: LogLevelDebug},
			Settings{Lang: "ru", EnableBackShadow: true, CardScale: 1.25, LogLevel: LogLevelDebug}},
		{"all_missing", Settings{}, Default()},
		{"unknown_lang", Settings{Lang: "fr", CardScale: 1, LogLevel: LogLevelInfo}, Settings{Lang: "en", CardScale: 1, LogLevel: LogLevelInfo}},
		{"zero_scale", Settings{Lang: "en", CardScale: 0, LogLevel: LogLevelInfo}, Default()},
		{"negative_scale", Settings{Lang: "en", CardScale: -2, LogLevel: LogLevelInfo}, Default()},
		{"nan_scale", Settings{Lang: "en", CardScale: math.NaN(), LogLevel: LogLevelInfo}, Default()},
		{"inf_scale", Settings{Lang: "en", CardScale: math.Inf(1), LogLevel: LogLevelInfo}, Settings{Lang: "en", CardScale: MaxCardScale, LogLevel: LogLevelInfo}},
		{"too_small", Settings{Lang: "en", CardScale: 0.01, LogLevel: LogLevelInfo}, Settings{Lang: "en", CardScale: MinCardScale, LogLevel: LogLevelInfo}},
		{"too_big", Settings{Lang: "en", CardScale: 25, LogLevel: LogLevelInfo}, Settings{Lang: "en", CardScale: MaxCardScale, LogLevel: LogLevelInfo}},
		{"missing_log_level", Settings{Lang: "en", CardScale: 1}, Default()},
		{"unknown_log_level", Settings{Lang: "en", CardScale: 1, LogLevel: "trace"}, Default()},
		{"warn_kept", Settings{Lang: "en", CardScale: 1, LogLevel: LogLevelWarn}, Settings{Lang: "en", CardScale: 1, LogLevel: LogLevelWarn}},
		{"error_kept", Settings{Lang: "en", CardScale: 1, LogLevel: LogLevelError}, Settings{Lang: "en", CardScale: 1, LogLevel: LogLevelError}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.in.Normalize(); got != tt.want {
				t.Fatalf("got %+v, want %+v", got, tt.want)
			}
		})
	}
}

func TestSlogLevel(t *testing.T) {
	for in, want := range map[string]slog.Level{
		LogLevelDebug: slog.LevelDebug,
		LogLevelInfo:  slog.LevelInfo,
		LogLevelWarn:  slog.LevelWarn,
		LogLevelError: slog.LevelError,
		"":            slog.LevelInfo,
		"trace":       slog.LevelInfo,
	} {
		if got := (Settings{LogLevel: in}).SlogLevel(); got != want {
			t.Errorf("SlogLevel(%q) = %v, want %v", in, got, want)
		}
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
