package repository

import (
	"errors"

	"github.com/HardDie/DeckBuilder/internal/config"
	dbSettings "github.com/HardDie/DeckBuilder/internal/db/settings"
	"github.com/HardDie/DeckBuilder/internal/entity"
	errors2 "github.com/HardDie/DeckBuilder/internal/errors"
)

type ISystemRepository interface {
	GetSettings() (*entity.SettingInfo, error)
	SaveSettings(set *entity.SettingInfo) error
}
type SystemRepository struct {
	cfg      *config.Config
	settings dbSettings.Settings
}

func NewSystemRepository(
	cfg *config.Config,
	settings dbSettings.Settings,
) *SystemRepository {
	return &SystemRepository{
		cfg:      cfg,
		settings: settings,
	}
}

func (s *SystemRepository) GetSettings() (*entity.SettingInfo, error) {
	resp, err := s.settings.Get()
	if err != nil {
		if errors.Is(err, errors2.SettingsNotExists) {
			return entity.NewSettings(), nil
		} else {
			return nil, err
		}
	}
	return resp, nil
}
func (s *SystemRepository) SaveSettings(set *entity.SettingInfo) error {
	return s.settings.Set(set)
}
