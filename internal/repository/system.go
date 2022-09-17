package repository

import (
	"tts_deck_build/internal/config"
	"tts_deck_build/internal/entity"
	"tts_deck_build/internal/fs"
)

type ISystemRepository interface {
	GetSettings() (*entity.SettingInfo, error)
	SaveSettings(set entity.SettingInfo) error
}
type SystemRepository struct {
	cfg *config.Config
}

func NewSystemRepository(cfg *config.Config) *SystemRepository {
	return &SystemRepository{
		cfg: cfg,
	}
}

func (s *SystemRepository) GetSettings() (*entity.SettingInfo, error) {
	// Check if such an object exists
	isExist, err := fs.IsFileExist(s.cfg.Settings())
	if err != nil || !isExist {
		return nil, err
	}

	// Read info from file
	return fs.OpenAndProcess(s.cfg.Settings(), fs.JsonFromReader[entity.SettingInfo])
}
func (s *SystemRepository) SaveSettings(set entity.SettingInfo) error {
	return fs.CreateAndProcess[entity.SettingInfo](s.cfg.Settings(), set, fs.JsonToWriter[entity.SettingInfo])
}
