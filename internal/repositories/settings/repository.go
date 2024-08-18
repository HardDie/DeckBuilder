package settings

import (
	"errors"

	"github.com/HardDie/DeckBuilder/internal/config"
	dbSettings "github.com/HardDie/DeckBuilder/internal/db/settings"
	entitiesSettings "github.com/HardDie/DeckBuilder/internal/entities/settings"
	errors2 "github.com/HardDie/DeckBuilder/internal/errors"
	"github.com/HardDie/DeckBuilder/internal/utils"
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

func (r *settings) Get() (*entitiesSettings.Settings, error) {
	resp, err := r.dbSettings.Get()
	if err != nil {
		if errors.Is(err, errors2.SettingsNotExists) {
			return utils.Allocate(entitiesSettings.Default()), nil
		} else {
			return nil, err
		}
	}
	return &entitiesSettings.Settings{
		Lang:             resp.Lang,
		EnableBackShadow: resp.EnableBackShadow,
		CardSize: entitiesSettings.CardSize{
			ScaleX: resp.CardSize.ScaleX,
			ScaleY: resp.CardSize.ScaleY,
			ScaleZ: resp.CardSize.ScaleZ,
		},
	}, nil
}
func (r *settings) Save(req *entitiesSettings.Settings) error {
	return r.dbSettings.Set(&dbSettings.SettingInfo{
		Lang:             req.Lang,
		EnableBackShadow: req.EnableBackShadow,
		CardSize: dbSettings.CardSize{
			ScaleX: req.CardSize.ScaleX,
			ScaleY: req.CardSize.ScaleY,
			ScaleZ: req.CardSize.ScaleZ,
		},
	})
}
