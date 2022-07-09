package settings

import (
	"tts_deck_build/internal/config"
)

type SettingService struct {
	storage *SettingStorage
}

func NewService() *SettingService {
	return &SettingService{
		storage: NewSettingStorage(config.GetConfig()),
	}
}

func (s *SettingService) Get() (*SettingInfo, error) {
	// Load default value
	settings := NewSettings()

	// Try to read settings from file
	set, err := s.storage.Get()
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

func (s *SettingService) Update(dto *UpdateSettingsDTO) (*SettingInfo, error) {
	set, err := s.Get()
	if err != nil {
		return nil, err
	}
	isUpdated := false
	switch dto.Lang {
	case "en", "ru":
		set.Lang = dto.Lang
		isUpdated = true
	}
	if isUpdated {
		err = s.storage.Save(*set)
		if err != nil {
			return nil, err
		}
	}
	return set, nil
}
