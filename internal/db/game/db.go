package game

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/HardDie/fsentry"

	entitiesGame "github.com/HardDie/DeckBuilder/internal/entities/game"
	er "github.com/HardDie/DeckBuilder/internal/errors"
	"github.com/HardDie/DeckBuilder/internal/logger"
)

type game struct {
	db        *fsentry.DB
	gamesPath string
}

func New(db *fsentry.DB) Game {
	return &game{
		db:        db,
		gamesPath: "games",
	}
}

func (d *game) Create(_ context.Context, req CreateRequest) (*entitiesGame.Game, error) {
	info, err := d.db.CreateFolder(req.Name, model{
		Description: fsentry.QuotedString(req.Description),
		Image:       fsentry.QuotedString(req.Image),
	}, d.gamesPath)
	if err != nil {
		if errors.Is(err, fsentry.ErrExist) {
			return nil, er.GameExist
		} else if errors.Is(err, fsentry.ErrBadName) {
			return nil, er.BadName
		} else {
			return nil, er.InternalError.AddMessage(err.Error())
		}
	}

	return d.toEntity(info), nil
}

func (d *game) Get(_ context.Context, name string) (*entitiesGame.Game, error) {
	info, err := d.db.GetFolder[model](name, d.gamesPath)
	if err != nil {
		if errors.Is(err, fsentry.ErrNotExist) {
			return nil, er.GameNotExists.AddMessage(err.Error()).HTTP(http.StatusBadRequest)
		} else if errors.Is(err, fsentry.ErrBadName) {
			return nil, er.BadName
		} else {
			return nil, er.InternalError.AddMessage(err.Error())
		}
	}

	return d.toEntity(info), nil
}

func (d *game) List(ctx context.Context) ([]*entitiesGame.Game, error) {
	list, err := d.db.List(d.gamesPath)
	if err != nil {
		return nil, er.InternalError.AddMessage(err.Error())
	}

	var games []*entitiesGame.Game
	for _, folder := range list.Folders {
		game, err := d.Get(ctx, folder)
		if err != nil {
			logger.Error.Println(folder, err.Error())
			continue
		}
		if folder != game.ID {
			logger.Error.Println("Corrupted game folder:", folder)
			continue
		}
		games = append(games, game)
	}
	return games, nil
}

func (d *game) Move(_ context.Context, oldName, newName string) (*entitiesGame.Game, error) {
	info, err := d.db.MoveFolder[model](oldName, newName, d.gamesPath)
	if err != nil {
		if errors.Is(err, fsentry.ErrNotExist) {
			return nil, er.GameNotExists.AddMessage(err.Error()).HTTP(http.StatusBadRequest)
		} else if errors.Is(err, fsentry.ErrBadName) {
			return nil, er.BadName
		} else {
			return nil, er.InternalError.AddMessage(err.Error())
		}
	}

	return d.toEntity(info), nil
}

func (d *game) Update(_ context.Context, req UpdateRequest) (*entitiesGame.Game, error) {
	info, err := d.db.UpdateFolder(req.Name, model{
		Description: fsentry.QuotedString(req.Description),
		Image:       fsentry.QuotedString(req.Image),
	}, d.gamesPath)
	if err != nil {
		if errors.Is(err, fsentry.ErrNotExist) {
			return nil, er.GameNotExists.AddMessage(err.Error()).HTTP(http.StatusBadRequest)
		} else if errors.Is(err, fsentry.ErrBadName) {
			return nil, er.BadName
		} else {
			return nil, er.InternalError.AddMessage(err.Error())
		}
	}

	return d.toEntity(info), nil
}

func (d *game) Delete(_ context.Context, name string) error {
	err := d.db.RemoveFolder(name, d.gamesPath)
	if err != nil {
		if errors.Is(err, fsentry.ErrNotExist) {
			return er.GameNotExists.AddMessage(err.Error()).HTTP(http.StatusBadRequest)
		} else if errors.Is(err, fsentry.ErrBadName) {
			return er.BadName
		} else {
			return er.InternalError.AddMessage(err.Error())
		}
	}
	return nil
}

func (d *game) Duplicate(_ context.Context, srcName, dstName string) (*entitiesGame.Game, error) {
	info, err := d.db.DuplicateFolder[model](srcName, dstName, d.gamesPath)
	if err != nil {
		if errors.Is(err, fsentry.ErrNotExist) {
			return nil, er.GameNotExists.AddMessage(err.Error()).HTTP(http.StatusBadRequest)
		} else if errors.Is(err, fsentry.ErrExist) {
			return nil, er.GameExist.AddMessage(err.Error()).HTTP(http.StatusBadRequest)
		} else if errors.Is(err, fsentry.ErrBadName) {
			return nil, er.BadName
		} else {
			return nil, er.InternalError.AddMessage(err.Error())
		}
	}
	return d.toEntity(info), nil
}

func (d *game) UpdateInfo(_ context.Context, name, newName string) error {
	_, err := d.db.UpdateFolderNameWithoutTimestamp[model](name, newName, d.gamesPath)
	return err
}

func (d *game) ImageCreate(ctx context.Context, gameID string, data []byte) error {
	game, err := d.Get(ctx, gameID)
	if err != nil {
		return err
	}

	err = d.db.CreateBinary("image", data, d.gamesPath, game.ID)
	if err != nil {
		if errors.Is(err, fsentry.ErrExist) {
			return er.GameImageExist.AddMessage(err.Error())
		} else {
			return er.InternalError.AddMessage(err.Error())
		}
	}
	return nil
}

func (d *game) ImageGet(ctx context.Context, gameID string) ([]byte, error) {
	game, err := d.Get(ctx, gameID)
	if err != nil {
		return nil, err
	}

	data, err := d.db.GetBinary("image", nil, d.gamesPath, game.ID)
	if err != nil {
		if errors.Is(err, fsentry.ErrNotExist) {
			return nil, er.GameImageNotExists.AddMessage(err.Error())
		} else {
			return nil, er.InternalError.AddMessage(err.Error())
		}
	}
	return data, nil
}

func (d *game) ImageDelete(ctx context.Context, gameID string) error {
	game, err := d.Get(ctx, gameID)
	if err != nil {
		return err
	}

	err = d.db.RemoveBinary("image", d.gamesPath, game.ID)
	if err != nil {
		if errors.Is(err, fsentry.ErrNotExist) {
			return er.GameImageNotExists.AddMessage(err.Error())
		} else {
			return er.InternalError.AddMessage(err.Error())
		}
	}
	return nil
}

func (d *game) toEntity(info fsentry.FolderInfo[model]) *entitiesGame.Game {
	createdAt, updatedAt := d.convertCreateUpdate(info.CreatedAt, info.UpdatedAt)
	return &entitiesGame.Game{
		ID:          info.ID,
		Name:        info.Name,
		Description: info.Data.Description.String(),
		Image:       info.Data.Image.String(),
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
	}
}

func (d *game) convertCreateUpdate(createdAt, updatedAt time.Time) (time.Time, time.Time) {
	if createdAt.IsZero() {
		createdAt = time.Now()
	}
	if updatedAt.IsZero() {
		updatedAt = createdAt
	}
	return createdAt, updatedAt
}
