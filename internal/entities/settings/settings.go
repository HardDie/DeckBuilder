package settings

import "math"

type CardSize struct {
	ScaleX float64
	ScaleY float64
	ScaleZ float64
}

type Settings struct {
	Lang             string
	EnableBackShadow bool
	CardSize         CardSize
}

func Default() Settings {
	return Settings{
		Lang:             "en",
		EnableBackShadow: false,
		CardSize: CardSize{
			ScaleX: 1,
			ScaleY: 1,
			ScaleZ: 1,
		},
	}
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
	s.CardSize.ScaleX = validScale(s.CardSize.ScaleX, def.CardSize.ScaleX)
	s.CardSize.ScaleY = validScale(s.CardSize.ScaleY, def.CardSize.ScaleY)
	s.CardSize.ScaleZ = validScale(s.CardSize.ScaleZ, def.CardSize.ScaleZ)
	return s
}

func validScale(v, def float64) float64 {
	if v > 0 && !math.IsInf(v, 1) {
		return v
	}
	return def
}
