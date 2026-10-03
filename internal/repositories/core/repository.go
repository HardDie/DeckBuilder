package core

import (
	"errors"
	"fmt"

	"github.com/HardDie/fsentry"
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
		return fmt.Errorf("open data folder: %w", err)
	}
	_, err = r.db.CreateFolder[any](r.gamesPath, nil)
	if err != nil {
		if !errors.Is(err, fsentry.ErrExist) {
			return fmt.Errorf("create games folder: %w", err)
		}
	}
	return nil
}

func (r *core) Drop() error {
	err := r.db.Drop()
	if err != nil {
		return fmt.Errorf("drop data folder: %w", err)
	}
	return nil
}
