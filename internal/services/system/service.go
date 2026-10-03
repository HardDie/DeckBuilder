package system

import (
	entitiesSettings "github.com/HardDie/DeckBuilder/internal/entities/settings"
	"github.com/HardDie/DeckBuilder/internal/logger"
	repositoriesSettings "github.com/HardDie/DeckBuilder/internal/repositories/settings"
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
func (s *system) GetSettings() (*entitiesSettings.Settings, error) {
	return s.repositorySettings.Get()
}
func (s *system) UpdateSettings(req UpdateSettingsRequest) (*entitiesSettings.Settings, error) {
	logger.Info.Println("Update settings")
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
