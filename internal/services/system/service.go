package system

import (
	"log"
	"os"

	"github.com/HardDie/DeckBuilder/internal/config"
	dbSettings "github.com/HardDie/DeckBuilder/internal/db/settings"
	"github.com/HardDie/DeckBuilder/internal/dto"
	"github.com/HardDie/DeckBuilder/internal/entity"
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
func (s *system) GetSettings() (*entity.SettingInfo, error) {
	// Load default value
	settings := entity.NewSettings()

	// Try to read settings from file
	set, err := s.repositorySettings.Get()
	if err != nil {
		return nil, err
	}

	// If got no settings from file
	if set == nil {
		// Return default value
		return settings, nil
	}

	// Update default values
	settings.Lang = set.Lang
	settings.EnableBackShadow = set.EnableBackShadow
	settings.CardSize.ScaleX = set.CardSize.ScaleX
	settings.CardSize.ScaleY = set.CardSize.ScaleY
	settings.CardSize.ScaleZ = set.CardSize.ScaleZ
	return settings, nil
}
func (s *system) UpdateSettings(dtoObject *dto.UpdateSettingsDTO) (*entity.SettingInfo, error) {
	log.Println("Update settings")
	set, err := s.GetSettings()
	if err != nil {
		return nil, err
	}
	isUpdated := false
	switch dtoObject.Lang {
	case "en", "ru":
		if set.Lang != dtoObject.Lang {
			set.Lang = dtoObject.Lang
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
