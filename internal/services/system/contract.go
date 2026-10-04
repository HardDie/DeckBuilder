package system

import (
	entitiesSettings "github.com/HardDie/DeckBuilder/internal/entities/settings"
)

type System interface {
	GetSettings() (*entitiesSettings.Settings, error)
	UpdateSettings(req UpdateSettingsRequest) (*entitiesSettings.Settings, error)
}

// UpdateSettingsRequest replaces the settings the dialog edits.
// An empty or unknown Lang keeps the stored language.
// An empty or unknown LogLevel keeps the stored log level.
type UpdateSettingsRequest struct {
	Lang             string
	EnableBackShadow bool
	CardScale        float64
	LogLevel         string
}
