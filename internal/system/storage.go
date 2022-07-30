package system

import (
	"tts_deck_build/internal/config"
	"tts_deck_build/internal/fs"
)

type SystemStorage struct {
	Config *config.Config
}

func NewSystemStorage(config *config.Config) *SystemStorage {
	return &SystemStorage{
		Config: config,
	}
}

func (s *SystemStorage) GetSettings() (*SettingInfo, error) {
	// Check if such an object exists
	isExist, err := fs.IsFileExist(s.Config.Settings())
	if err != nil || !isExist {
		return nil, err
	}

	// Read info from file
	return fs.OpenAndProcess(s.Config.Settings(), fs.JsonFromReader[SettingInfo])
}

func (s *SystemStorage) SaveSettings(set SettingInfo) error {
	return fs.CreateAndProcess[SettingInfo](s.Config.Settings(), set, fs.JsonToWriter[SettingInfo])
}
