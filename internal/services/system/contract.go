package system

import (
	"github.com/HardDie/DeckBuilder/internal/entity"
)

type System interface {
	Quit()
	GetSettings() (*entity.SettingInfo, error)
	UpdateSettings(req UpdateSettingsRequest) (*entity.SettingInfo, error)
}

type UpdateSettingsRequest struct {
	Lang string
}
