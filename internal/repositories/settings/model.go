package settings

type cardSize struct {
	ScaleX float64 `json:"scaleX"`
	ScaleY float64 `json:"scaleY"`
	ScaleZ float64 `json:"scaleZ"`
}

type model struct {
	Lang             string   `json:"lang"`
	EnableBackShadow bool     `json:"enable_back_shadow"`
	CardSize         cardSize `json:"card_size"`
}
