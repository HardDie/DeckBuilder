package entity

type SettingInfo struct {
	Lang string `json:"lang"`
}

func NewSettings() *SettingInfo {
	return &SettingInfo{
		Lang: "en",
	}
}
