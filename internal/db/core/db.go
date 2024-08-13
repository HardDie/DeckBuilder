package core

import (
	"errors"

	"github.com/HardDie/fsentry"
	"github.com/HardDie/fsentry/pkg/fsentry_error"

	er "github.com/HardDie/DeckBuilder/internal/errors"
)

type core struct {
	db        fsentry.IFSEntry
	gamesPath string
}

func New(db fsentry.IFSEntry) Core {
	return &core{
		db:        db,
		gamesPath: "games",
	}
}

func (d *core) Init() error {
	err := d.db.Init()
	if err != nil {
		return er.InternalError.AddMessage(err.Error())
	}
	_, err = d.db.CreateFolder(d.gamesPath, nil)
	if err != nil {
		if !errors.Is(err, fsentry_error.ErrorExist) {
			return er.InternalError.AddMessage(err.Error())
		}
	}
	return nil
}
func (d *core) Drop() error {
	err := d.db.Drop()
	if err != nil {
		return er.InternalError.AddMessage(err.Error())
	}
	return nil
}
