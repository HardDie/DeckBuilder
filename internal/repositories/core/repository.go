package core

import (
	"errors"

	"github.com/HardDie/fsentry"

	er "github.com/HardDie/DeckBuilder/internal/errors"
)

type core struct {
	db        *fsentry.DB
	gamesPath string
}

func New(db *fsentry.DB) Core {
	return &core{
		db:        db,
		gamesPath: "games",
	}
}

func (r *core) Init() error {
	err := r.db.Init()
	if err != nil {
		return er.InternalError.AddMessage(err.Error())
	}
	_, err = r.db.CreateFolder[any](r.gamesPath, nil)
	if err != nil {
		if !errors.Is(err, fsentry.ErrExist) {
			return er.InternalError.AddMessage(err.Error())
		}
	}
	return nil
}

func (r *core) Drop() error {
	err := r.db.Drop()
	if err != nil {
		return er.InternalError.AddMessage(err.Error())
	}
	return nil
}
