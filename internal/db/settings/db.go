package settings

import (
	"errors"

	"github.com/HardDie/fsentry"

	er "github.com/HardDie/DeckBuilder/internal/errors"
)

type settings struct {
	db *fsentry.DB
}

func New(db *fsentry.DB) Settings {
	return &settings{
		db: db,
	}
}

func (d *settings) Get() (*SettingInfo, error) {
	info, err := d.db.GetEntry[SettingInfo]("settings")
	if err != nil {
		if errors.Is(err, fsentry.ErrNotExist) {
			return nil, er.SettingsNotExists.AddMessage(err.Error())
		}
		return nil, er.InternalError.AddMessage(err.Error())
	}
	setting := info.Data
	return &setting, nil
}

func (d *settings) Set(data *SettingInfo) error {
	_, err := d.db.CreateEntry("settings", data)
	if err == nil {
		return nil
	}
	if !errors.Is(err, fsentry.ErrExist) {
		return err
	}
	_, err = d.db.UpdateEntry("settings", data)
	if err != nil {
		return er.InternalError.AddMessage(err.Error())
	}
	return nil
}
