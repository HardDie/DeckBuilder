package settings

import (
	"errors"

	"github.com/HardDie/fsentry"

	"github.com/HardDie/DeckBuilder/internal/config"
	entitiesSettings "github.com/HardDie/DeckBuilder/internal/entities/settings"
	er "github.com/HardDie/DeckBuilder/internal/errors"
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

func (r *settings) Get() (*entitiesSettings.Settings, error) {
	resp, err := r.get()
	if err != nil {
		if errors.Is(err, er.SettingsNotExists) {
			return utils.Allocate(entitiesSettings.Default()), nil
		}
		return nil, err
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
	return r.set(&model{
		Lang:             req.Lang,
		EnableBackShadow: req.EnableBackShadow,
		CardSize: cardSize{
			ScaleX: req.CardSize.ScaleX,
			ScaleY: req.CardSize.ScaleY,
			ScaleZ: req.CardSize.ScaleZ,
		},
	})
}

func (r *settings) get() (*model, error) {
	info, err := r.db.GetEntry[model]("settings")
	if err != nil {
		if errors.Is(err, fsentry.ErrNotExist) {
			return nil, er.SettingsNotExists.AddMessage(err.Error())
		}
		return nil, er.InternalError.AddMessage(err.Error())
	}
	setting := info.Data
	return &setting, nil
}

func (r *settings) set(data *model) error {
	_, err := r.db.CreateEntry("settings", data)
	if err == nil {
		return nil
	}
	if !errors.Is(err, fsentry.ErrExist) {
		return err
	}
	_, err = r.db.UpdateEntry("settings", data)
	if err != nil {
		return er.InternalError.AddMessage(err.Error())
	}
	return nil
}
