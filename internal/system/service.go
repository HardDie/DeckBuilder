package system

import (
	"os"

	"github.com/HardDie/DeckBuilder/internal/config"
	"github.com/HardDie/DeckBuilder/internal/dto"
	"github.com/HardDie/DeckBuilder/internal/entity"
	"github.com/HardDie/DeckBuilder/internal/repository"
)

type SystemService struct {
	rep repository.ISystemRepository
}

func NewService(cfg *config.Config) *SystemService {
	return &SystemService{
		rep: repository.NewSystemRepository(cfg),
	}
}

func (s *SystemService) Quit() {
	os.Exit(0)
}
func (s *SystemService) GetSettings() (*entity.SettingInfo, error) {
	// Load default value
	settings := entity.NewSettings()

	// Try to read settings from file
	set, err := s.rep.GetSettings()
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
func (s *SystemService) UpdateSettings(dtoObject *dto.UpdateSettingsDTO) (*entity.SettingInfo, error) {
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
		err = s.rep.SaveSettings(*set)
		if err != nil {
			return nil, err
		}
	}
	return set, nil
}
