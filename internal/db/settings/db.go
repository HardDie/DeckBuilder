package settings

import (
	"encoding/json"
	"errors"

	"github.com/HardDie/fsentry"
	"github.com/HardDie/fsentry/pkg/fsentry_error"

	er "github.com/HardDie/DeckBuilder/internal/errors"
)

type settings struct {
	db fsentry.IFSEntry
}

func New(db fsentry.IFSEntry) Settings {
	return &settings{
		db: db,
	}
}

func (d *settings) Get() (*SettingInfo, error) {
	info, err := d.db.GetEntry("settings")
	if err != nil {
		if errors.Is(err, fsentry_error.ErrorNotExist) {
			return nil, er.SettingsNotExists.AddMessage(err.Error())
		} else {
			return nil, er.InternalError.AddMessage(err.Error())
		}
	}
	setting := &SettingInfo{}

	err = json.Unmarshal(info.Data, setting)
	if err != nil {
		return nil, er.InternalError.AddMessage(err.Error())
	}

	return setting, nil
}
func (d *settings) Set(data *SettingInfo) error {
	err := d.db.CreateEntry("settings", data)
	if err == nil {
		return nil
	}
	if !errors.Is(err, fsentry_error.ErrorExist) {
		return err
	}
	err = d.db.UpdateEntry("settings", data)
	if err != nil {
		return er.InternalError.AddMessage(err.Error())
	}
	return nil
}
