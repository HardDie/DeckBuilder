package settings

import (
	"errors"
	"fmt"

	"github.com/HardDie/fsentry"

	"github.com/HardDie/DeckBuilder/internal/config"
	entitiesSettings "github.com/HardDie/DeckBuilder/internal/entities/settings"
	"github.com/HardDie/DeckBuilder/internal/utils"
)

type settings struct {
	cfg *config.Config
	db  *fsentry.DB
}

func New(cfg *config.Config, db *fsentry.DB) Settings {
	return &settings{
		cfg: cfg,
		db:  db,
	}
}

// Get returns the stored settings, or the defaults when none are stored.
func (r *settings) Get() (*entitiesSettings.Settings, error) {
	info, err := r.db.GetEntry[model]("settings")
	if errors.Is(err, fsentry.ErrNotExist) {
		return utils.Allocate(entitiesSettings.Default()), nil
	}
	if err != nil {
		return nil, fmt.Errorf("read settings: %w", err)
	}
	resp := info.Data
	return &entitiesSettings.Settings{
		Lang:             resp.Lang,
		EnableBackShadow: resp.EnableBackShadow,
		CardScale:        resp.CardSize.ScaleX,
		LogLevel:         resp.LogLevel,
	}, nil
}

func (r *settings) Save(req *entitiesSettings.Settings) error {
	return r.set(&model{
		Lang:             req.Lang,
		EnableBackShadow: req.EnableBackShadow,
		CardSize: cardSize{
			ScaleX: req.CardScale,
			ScaleY: 1,
			ScaleZ: req.CardScale,
		},
		LogLevel: req.LogLevel,
	})
}

func (r *settings) set(data *model) error {
	_, err := r.db.CreateEntry("settings", data)
	if err == nil {
		return nil
	}
	if !errors.Is(err, fsentry.ErrExist) {
		return fmt.Errorf("save settings: %w", err)
	}
	_, err = r.db.UpdateEntry("settings", data)
	if err != nil {
		return fmt.Errorf("save settings: %w", err)
	}
	return nil
}
