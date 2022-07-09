package settings

import (
	"tts_deck_build/internal/config"
	"tts_deck_build/internal/fs"
)

type SettingStorage struct {
	Config *config.Config
}

func NewSettingStorage(config *config.Config) *SettingStorage {
	return &SettingStorage{
		Config: config,
	}
}

func (s *SettingStorage) Get() (*SettingInfo, error) {
	// Check if such an object exists
	isExist, err := fs.IsFileExist(s.Config.Settings())
	if err != nil || !isExist {
		return nil, err
	}

	// Read info from file
	return fs.OpenAndProcess(s.Config.Settings(), fs.JsonFromReader[SettingInfo])
}

func (s *SettingStorage) Save(set SettingInfo) error {
	return fs.CreateAndProcess[SettingInfo](s.Config.Settings(), set, fs.JsonToWriter[SettingInfo])
}
