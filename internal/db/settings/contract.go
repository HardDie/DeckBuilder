package settings

type Settings interface {
	Get() (*SettingInfo, error)
	Set(data *SettingInfo) error
}

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
