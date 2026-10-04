package settings

import "math"

// Card scale limits for the settings dialog. TTS itself has no documented limit.
const (
	MinCardScale = 0.1
	MaxCardScale = 10
)

type Settings struct {
	Lang             string
	EnableBackShadow bool
	// CardScale is the size of spawned cards and decks. 1 is the TTS default.
	// TTS scales a card in width and length only, so it goes to X and Z; Y stays 1.
	CardScale float64
}

func Default() Settings {
	return Settings{
		Lang:             "en",
		EnableBackShadow: false,
		CardScale:        1,
	}
}

// ValidCardScale reports whether v is a card scale the settings dialog accepts.
func ValidCardScale(v float64) bool {
	return v >= MinCardScale && v <= MaxCardScale
}

// Normalize replaces values that cannot be valid with their defaults.
// A field missing from an older or hand-edited settings file reads as zero;
// a card scale of 0 would make TTS spawn invisible cards.
func (s Settings) Normalize() Settings {
	def := Default()
	switch s.Lang {
	case "en", "ru":
	default:
		s.Lang = def.Lang
	}
	switch {
	case math.IsNaN(s.CardScale) || s.CardScale <= 0:
		s.CardScale = def.CardScale
	case s.CardScale < MinCardScale:
		s.CardScale = MinCardScale
	case s.CardScale > MaxCardScale:
		s.CardScale = MaxCardScale
	}
	return s
}
