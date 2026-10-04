package dto

type Settings struct {
	Lang             string  `json:"lang"`
	EnableBackShadow bool    `json:"enable_back_shadow"`
	CardScale        float64 `json:"card_scale"`
}
