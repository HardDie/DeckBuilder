package settings

import "github.com/HardDie/DeckBuilder/internal/entity"

type Settings interface {
	Get() (*entity.SettingInfo, error)
	Save(set *entity.SettingInfo) error
}
