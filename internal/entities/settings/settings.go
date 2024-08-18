package settings

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
