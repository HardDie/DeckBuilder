package system

import (
	"log"

	entitiesSettings "github.com/HardDie/DeckBuilder/internal/entities/settings"
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

func (s *system) GetSettings() (*entitiesSettings.Settings, error) {
	settings := entitiesSettings.Default()

	set, err := s.repositorySettings.Get()
	if err != nil {
		return nil, err
	}

	if set == nil {
		return &settings, nil
	}

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
