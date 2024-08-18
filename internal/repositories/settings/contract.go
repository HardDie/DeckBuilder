package settings

import (
	entitiesSettings "github.com/HardDie/DeckBuilder/internal/entities/settings"
)

type Settings interface {
	Get() (*entitiesSettings.Settings, error)
	Save(set *entitiesSettings.Settings) error
}
