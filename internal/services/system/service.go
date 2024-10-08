package system

import (
	"log"
	"os"

	"github.com/HardDie/DeckBuilder/internal/config"
	dbSettings "github.com/HardDie/DeckBuilder/internal/db/settings"
	entitiesSettings "github.com/HardDie/DeckBuilder/internal/entities/settings"
	repositoriesSettings "github.com/HardDie/DeckBuilder/internal/repositories/settings"
)

type system struct {
	repositorySettings repositoriesSettings.Settings
}

func New(cfg *config.Config, settings dbSettings.Settings) System {
	return &system{
		repositorySettings: repositoriesSettings.New(cfg, settings),
	}
}

func (s *system) Quit() {
	os.Exit(0)
}
func (s *system) GetSettings() (*entitiesSettings.Settings, error) {
	// Load default value
	settings := entitiesSettings.Default()

	// Try to read settings from file
	set, err := s.repositorySettings.Get()
	if err != nil {
		return nil, err
	}

	// If got no settings from file
	if set == nil {
		// Return default value
		return &settings, nil
	}

	// Update default values
	settings.Lang = set.Lang
	settings.EnableBackShadow = set.EnableBackShadow
	settings.CardSize.ScaleX = set.CardSize.ScaleX
	settings.CardSize.ScaleY = set.CardSize.ScaleY
	settings.CardSize.ScaleZ = set.CardSize.ScaleZ
	return &settings, nil
}
func (s *system) UpdateSettings(req UpdateSettingsRequest) (*entitiesSettings.Settings, error) {
	log.Println("Update settings")
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
