package entity

type CardSize struct {
	ScaleX float64 `json:"scaleX"`
	ScaleY float64 `json:"scaleY"`
	ScaleZ float64 `json:"scaleZ"`
}

type SettingInfo struct {
	Lang             string   `json:"lang"`
	EnableBackShadow bool     `json:"enable_back_shadow"`
	CardSize         CardSize `json:"card_size"`
}

func NewSettings() *SettingInfo {
	return &SettingInfo{
		Lang:             "en",
		EnableBackShadow: false,
		CardSize: CardSize{
			ScaleX: 1,
			ScaleY: 1,
			ScaleZ: 1,
		},
	}
}

type Status struct {
	Type     string  `json:"type"`
	Message  string  `json:"message"`
	Progress float32 `json:"progress"`
	Status   string  `json:"status"`
}
