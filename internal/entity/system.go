package entity

type SettingInfo struct {
	Lang string `json:"lang"`
}

func NewSettings() *SettingInfo {
	return &SettingInfo{
		Lang: "en",
	}
}

type Status struct {
	Type     string  `json:"type"`
	Message  string  `json:"message"`
	Progress float32 `json:"progress"`
}
