package system

import (
	"github.com/HardDie/DeckBuilder/internal/dto"
	"github.com/HardDie/DeckBuilder/internal/entity"
)

type System interface {
	Quit()
	GetSettings() (*entity.SettingInfo, error)
	UpdateSettings(dtoObject *dto.UpdateSettingsDTO) (*entity.SettingInfo, error)
}
