package settings

import (
	"errors"

	"github.com/HardDie/DeckBuilder/internal/config"
	dbSettings "github.com/HardDie/DeckBuilder/internal/db/settings"
	"github.com/HardDie/DeckBuilder/internal/entity"
	errors2 "github.com/HardDie/DeckBuilder/internal/errors"
)

type settings struct {
	cfg        *config.Config
	dbSettings dbSettings.Settings
}

func New(cfg *config.Config, dbSettings dbSettings.Settings) Settings {
	return &settings{
		cfg:        cfg,
		dbSettings: dbSettings,
	}
}

func (r *settings) Get() (*entity.SettingInfo, error) {
	resp, err := r.dbSettings.Get()
	if err != nil {
		if errors.Is(err, errors2.SettingsNotExists) {
			return entity.NewSettings(), nil
		} else {
			return nil, err
		}
	}
	return resp, nil
}
func (r *settings) Save(set *entity.SettingInfo) error {
	return r.dbSettings.Set(set)
}
