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
	Type     string
	Message  string
	Progress float32
}
