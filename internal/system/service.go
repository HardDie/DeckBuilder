package system

import (
	"os"

	"tts_deck_build/internal/config"
	"tts_deck_build/internal/dto"
)

type SystemService struct {
	storage *SystemStorage
}

func NewService() *SystemService {
	return &SystemService{
		storage: NewSystemStorage(config.GetConfig()),
	}
}

func (s *SystemService) Quit() {
	os.Exit(0)
}

func (s *SystemService) GetSettings() (*SettingInfo, error) {
	// Load default value
	settings := NewSettings()

	// Try to read settings from file
	set, err := s.storage.GetSettings()
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
	return settings, nil
}

func (s *SystemService) UpdateSettings(dtoObject *dto.UpdateSettingsDTO) (*SettingInfo, error) {
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
		err = s.storage.SaveSettings(*set)
		if err != nil {
			return nil, err
		}
	}
	return set, nil
}
