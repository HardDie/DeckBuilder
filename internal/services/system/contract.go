package system

import (
	entitiesSettings "github.com/HardDie/DeckBuilder/internal/entities/settings"
)

type System interface {
	Quit()
	GetSettings() (*entitiesSettings.Settings, error)
	UpdateSettings(req UpdateSettingsRequest) (*entitiesSettings.Settings, error)
}

type UpdateSettingsRequest struct {
	Lang string
}
