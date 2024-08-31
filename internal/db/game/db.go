package game

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/HardDie/fsentry"
	"github.com/HardDie/fsentry/pkg/fsentry_error"
	"github.com/HardDie/fsentry/pkg/fsentry_types"

	entitiesGame "github.com/HardDie/DeckBuilder/internal/entities/game"
	er "github.com/HardDie/DeckBuilder/internal/errors"
	"github.com/HardDie/DeckBuilder/internal/logger"
	"github.com/HardDie/DeckBuilder/internal/utils"
)

type game struct {
	db        fsentry.IFSEntry
	gamesPath string
}

func New(db fsentry.IFSEntry) Game {
	return &game{
		db:        db,
		gamesPath: "games",
	}
}

func (d *game) Create(_ context.Context, name, description, image string) (*entitiesGame.Game, error) {
	info, err := d.db.CreateFolder(name, model{
		Description: fsentry_types.QS(description),
		Image:       fsentry_types.QS(image),
	}, d.gamesPath)
	if err != nil {
		if errors.Is(err, fsentry_error.ErrorExist) {
			return nil, er.GameExist
		} else if errors.Is(err, fsentry_error.ErrorBadName) {
			return nil, er.BadName
		} else {
			return nil, er.InternalError.AddMessage(err.Error())
		}
	}

	createdAt, updatedAt := d.convertCreateUpdate(info.CreatedAt, info.UpdatedAt)
	return &entitiesGame.Game{
		ID:          info.Id,
		Name:        info.Name.String(),
		Description: description,
		Image:       image,
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
	}, nil
}
func (d *game) Get(_ context.Context, name string) (*entitiesGame.Game, error) {
	info, err := d.db.GetFolder(name, d.gamesPath)
	if err != nil {
		if errors.Is(err, fsentry_error.ErrorNotExist) {
			return nil, er.GameNotExists.AddMessage(err.Error()).HTTP(http.StatusBadRequest)
		} else if errors.Is(err, fsentry_error.ErrorBadName) {
			return nil, er.BadName
		} else {
			return nil, er.InternalError.AddMessage(err.Error())
		}
	}

	var gInfo model
	err = json.Unmarshal(info.Data, &gInfo)
	if err != nil {
		return nil, er.InternalError.AddMessage(err.Error())
	}

	createdAt, updatedAt := d.convertCreateUpdate(info.CreatedAt, info.UpdatedAt)
	return &entitiesGame.Game{
		ID:          info.Id,
		Name:        info.Name.String(),
		Description: gInfo.Description.String(),
		Image:       gInfo.Image.String(),
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
	}, nil
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
	info, err := d.db.MoveFolder(oldName, newName, d.gamesPath)
	if err != nil {
		if errors.Is(err, fsentry_error.ErrorNotExist) {
			return nil, er.GameNotExists.AddMessage(err.Error()).HTTP(http.StatusBadRequest)
		} else if errors.Is(err, fsentry_error.ErrorBadName) {
			return nil, er.BadName
		} else {
			return nil, er.InternalError.AddMessage(err.Error())
		}
	}

	var gInfo model
	err = json.Unmarshal(info.Data, &gInfo)
	if err != nil {
		return nil, er.InternalError.AddMessage(err.Error())
	}

	createdAt, updatedAt := d.convertCreateUpdate(info.CreatedAt, info.UpdatedAt)
	return &entitiesGame.Game{
		ID:          info.Id,
		Name:        info.Name.String(),
		Description: gInfo.Description.String(),
		Image:       gInfo.Image.String(),
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
	}, nil
}
func (d *game) Update(_ context.Context, name, description, image string) (*entitiesGame.Game, error) {
	info, err := d.db.UpdateFolder(name, &model{
		Description: fsentry_types.QS(description),
		Image:       fsentry_types.QS(image),
	}, d.gamesPath)
	if err != nil {
		if errors.Is(err, fsentry_error.ErrorNotExist) {
			return nil, er.GameNotExists.AddMessage(err.Error()).HTTP(http.StatusBadRequest)
		} else if errors.Is(err, fsentry_error.ErrorBadName) {
			return nil, er.BadName
		} else {
			return nil, er.InternalError.AddMessage(err.Error())
		}
	}

	var gInfo model
	err = json.Unmarshal(info.Data, &gInfo)
	if err != nil {
		return nil, er.InternalError.AddMessage(err.Error())
	}

	createdAt, updatedAt := d.convertCreateUpdate(info.CreatedAt, info.UpdatedAt)
	return &entitiesGame.Game{
		ID:          info.Id,
		Name:        info.Name.String(),
		Description: gInfo.Description.String(),
		Image:       gInfo.Image.String(),
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
	}, nil
}
func (d *game) Delete(_ context.Context, name string) error {
	err := d.db.RemoveFolder(name, d.gamesPath)
	if err != nil {
		if errors.Is(err, fsentry_error.ErrorNotExist) {
			return er.GameNotExists.AddMessage(err.Error()).HTTP(http.StatusBadRequest)
		} else if errors.Is(err, fsentry_error.ErrorBadName) {
			return er.BadName
		} else {
			return er.InternalError.AddMessage(err.Error())
		}
	}
	return nil
}
func (d *game) Duplicate(_ context.Context, srcName, dstName string) (*entitiesGame.Game, error) {
	info, err := d.db.DuplicateFolder(srcName, dstName, d.gamesPath)
	if err != nil {
		if errors.Is(err, fsentry_error.ErrorNotExist) {
			return nil, er.GameNotExists.AddMessage(err.Error()).HTTP(http.StatusBadRequest)
		} else if errors.Is(err, fsentry_error.ErrorExist) {
			return nil, er.GameExist.AddMessage(err.Error()).HTTP(http.StatusBadRequest)
		} else if errors.Is(err, fsentry_error.ErrorBadName) {
			return nil, er.BadName
		} else {
			return nil, er.InternalError.AddMessage(err.Error())
		}
	}

	var gInfo model
	err = json.Unmarshal(info.Data, &gInfo)
	if err != nil {
		return nil, er.InternalError.AddMessage(err.Error())
	}

	createdAt, updatedAt := d.convertCreateUpdate(info.CreatedAt, info.UpdatedAt)
	return &entitiesGame.Game{
		ID:          info.Id,
		Name:        info.Name.String(),
		Description: gInfo.Description.String(),
		Image:       gInfo.Image.String(),
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
	}, nil
}
func (d *game) UpdateInfo(_ context.Context, name, newName string) error {
	return d.db.UpdateFolderNameWithoutTimestamp(name, newName, d.gamesPath)
}
func (d *game) ImageCreate(ctx context.Context, gameID string, data []byte) error {
	game, err := d.Get(ctx, gameID)
	if err != nil {
		return err
	}

	err = d.db.CreateBinary("image", data, d.gamesPath, game.ID)
	if err != nil {
		if errors.Is(err, fsentry_error.ErrorExist) {
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

	data, err := d.db.GetBinary("image", d.gamesPath, game.ID)
	if err != nil {
		if errors.Is(err, fsentry_error.ErrorNotExist) {
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
		if errors.Is(err, fsentry_error.ErrorNotExist) {
			return er.GameImageNotExists.AddMessage(err.Error())
		} else {
			return er.InternalError.AddMessage(err.Error())
		}
	}
	return nil
}

func (d *game) convertCreateUpdate(createdAt, updatedAt *time.Time) (time.Time, time.Time) {
	if createdAt == nil {
		createdAt = utils.Allocate(time.Now())
	}
	if updatedAt == nil {
		updatedAt = createdAt
	}
	return *createdAt, *updatedAt
}
