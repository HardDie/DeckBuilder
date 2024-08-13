package settings

import "github.com/HardDie/DeckBuilder/internal/entity"

type Settings interface {
	Get() (*entity.SettingInfo, error)
	Set(data *entity.SettingInfo) error
}
