package system

import (
	entitiesSettings "github.com/HardDie/DeckBuilder/internal/entities/settings"
	repositoriesSettings "github.com/HardDie/DeckBuilder/internal/repositories/settings"
	"log/slog"
)

type system struct {
	repositorySettings repositoriesSettings.Settings
}

func New(repositorySettings repositoriesSettings.Settings) System {
	return &system{
		repositorySettings: repositorySettings,
	}
}

// GetSettings returns the stored settings, or the defaults when none are stored.
// Missing or invalid values read as their defaults (see Settings.Normalize).
func (s *system) GetSettings() (*entitiesSettings.Settings, error) {
	set, err := s.repositorySettings.Get()
	if err != nil {
		return nil, err
	}
	normalized := set.Normalize()
	return &normalized, nil
}
func (s *system) UpdateSettings(req UpdateSettingsRequest) (*entitiesSettings.Settings, error) {
	slog.Info("update settings", "lang", req.Lang)
	set, err := s.GetSettings()
	if err != nil {
		return nil, err
	}
	isUpdated := false
	switch req.Lang {
	case "en", "ru":
		if set.Lang != req.Lang {
			set.Lang = req.Lang
			isUpdated = true
		}
	}
	if isUpdated {
		err = s.repositorySettings.Save(set)
		if err != nil {
			return nil, err
		}
	}
	return set, nil
}
