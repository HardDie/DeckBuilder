package collection

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/HardDie/fsentry"
	"github.com/HardDie/fsentry/pkg/fsentry_error"
	"github.com/HardDie/fsentry/pkg/fsentry_types"

	dbCommon "github.com/HardDie/DeckBuilder/internal/db/common"
	dbGame "github.com/HardDie/DeckBuilder/internal/db/game"
	entitiesCollection "github.com/HardDie/DeckBuilder/internal/entities/collection"
	er "github.com/HardDie/DeckBuilder/internal/errors"
	"github.com/HardDie/DeckBuilder/internal/logger"
	"github.com/HardDie/DeckBuilder/internal/utils"
)

type collection struct {
	db        fsentry.IFSEntry
	gamesPath string

	game dbGame.Game
}

func New(db fsentry.IFSEntry, game dbGame.Game) Collection {
	return &collection{
		db:        db,
		gamesPath: "games",

		game: game,
	}
}

func (d *collection) Create(ctx context.Context, gameID, name, description, image string) (*entitiesCollection.Collection, error) {
	game, err := d.game.Get(ctx, gameID)
	if err != nil {
		return nil, err
	}

	info, err := d.db.CreateFolder(name, &model{
		Description: fsentry_types.QS(description),
		Image:       fsentry_types.QS(image),
	}, d.gamesPath, game.ID)
	if err != nil {
		if errors.Is(err, fsentry_error.ErrorExist) {
			return nil, er.CollectionExist
		} else if errors.Is(err, fsentry_error.ErrorBadName) {
			return nil, er.BadName
		} else {
			return nil, er.InternalError.AddMessage(err.Error())
		}
	}

	createdAt, updatedAt := d.convertCreateUpdate(info.CreatedAt, info.UpdatedAt)
	return &entitiesCollection.Collection{
		ID:          info.Id,
		Name:        info.Name.String(),
		Description: description,
		Image:       image,
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,

		GameID: gameID,
	}, nil
}
func (d *collection) Get(ctx context.Context, gameID, name string) (*entitiesCollection.Collection, error) {
	game, err := d.game.Get(ctx, gameID)
	if err != nil {
		return nil, err
	}

	info, err := d.db.GetFolder(name, d.gamesPath, game.ID)
	if err != nil {
		if errors.Is(err, fsentry_error.ErrorNotExist) {
			return nil, er.CollectionNotExists.AddMessage(err.Error()).HTTP(http.StatusBadRequest)
		} else if errors.Is(err, fsentry_error.ErrorBadName) {
			return nil, er.BadName
		} else {
			return nil, er.InternalError.AddMessage(err.Error())
		}
	}

	var cInfo dbCommon.Info
	err = json.Unmarshal(info.Data, &cInfo)
	if err != nil {
		return nil, er.InternalError.AddMessage(err.Error())
	}

	createdAt, updatedAt := d.convertCreateUpdate(info.CreatedAt, info.UpdatedAt)
	return &entitiesCollection.Collection{
		ID:          info.Id,
		Name:        info.Name.String(),
		Description: cInfo.Description.String(),
		Image:       cInfo.Image.String(),
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,

		GameID: gameID,
	}, nil
}
func (d *collection) List(ctx context.Context, gameID string) ([]*entitiesCollection.Collection, error) {
	game, err := d.game.Get(ctx, gameID)
	if err != nil {
		return nil, err
	}

	list, err := d.db.List(d.gamesPath, game.ID)
	if err != nil {
		return nil, er.InternalError.AddMessage(err.Error())
	}

	var collections []*entitiesCollection.Collection
	for _, folder := range list.Folders {
		collection, err := d.Get(ctx, game.ID, folder)
		if err != nil {
			logger.Error.Println(folder, err.Error())
			continue
		}
		if folder != collection.ID {
			logger.Error.Println("Corrupted collection folder:", folder)
			continue
		}
		collections = append(collections, collection)
	}
	return collections, nil
}
func (d *collection) Move(ctx context.Context, gameID, oldName, newName string) (*entitiesCollection.Collection, error) {
	game, err := d.game.Get(ctx, gameID)
	if err != nil {
		return nil, err
	}

	info, err := d.db.MoveFolder(oldName, newName, d.gamesPath, game.ID)
	if err != nil {
		if errors.Is(err, fsentry_error.ErrorNotExist) {
			return nil, er.CollectionNotExists.AddMessage(err.Error()).HTTP(http.StatusBadRequest)
		} else if errors.Is(err, fsentry_error.ErrorBadName) {
			return nil, er.BadName
		} else {
			return nil, er.InternalError.AddMessage(err.Error())
		}
	}

	var cInfo dbCommon.Info
	err = json.Unmarshal(info.Data, &cInfo)
	if err != nil {
		return nil, er.InternalError.AddMessage(err.Error())
	}

	createdAt, updatedAt := d.convertCreateUpdate(info.CreatedAt, info.UpdatedAt)
	return &entitiesCollection.Collection{
		ID:          info.Id,
		Name:        info.Name.String(),
		Description: cInfo.Description.String(),
		Image:       cInfo.Image.String(),
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,

		GameID: gameID,
	}, nil
}
func (d *collection) Update(ctx context.Context, gameID, name, description, image string) (*entitiesCollection.Collection, error) {
	game, err := d.game.Get(ctx, gameID)
	if err != nil {
		return nil, err
	}

	info, err := d.db.UpdateFolder(name, &dbCommon.Info{
		Description: fsentry_types.QS(description),
		Image:       fsentry_types.QS(image),
	}, d.gamesPath, game.ID)
	if err != nil {
		if errors.Is(err, fsentry_error.ErrorNotExist) {
			return nil, er.CollectionNotExists.AddMessage(err.Error()).HTTP(http.StatusBadRequest)
		} else if errors.Is(err, fsentry_error.ErrorBadName) {
			return nil, er.BadName
		} else {
			return nil, er.InternalError.AddMessage(err.Error())
		}
	}

	var cInfo dbCommon.Info
	err = json.Unmarshal(info.Data, &cInfo)
	if err != nil {
		return nil, er.InternalError.AddMessage(err.Error())
	}

	createdAt, updatedAt := d.convertCreateUpdate(info.CreatedAt, info.UpdatedAt)
	return &entitiesCollection.Collection{
		ID:          info.Id,
		Name:        info.Name.String(),
		Description: cInfo.Description.String(),
		Image:       cInfo.Image.String(),
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,

		GameID: gameID,
	}, nil
}
func (d *collection) Delete(ctx context.Context, gameID, name string) error {
	game, err := d.game.Get(ctx, gameID)
	if err != nil {
		return err
	}

	err = d.db.RemoveFolder(name, d.gamesPath, game.ID)
	if err != nil {
		if errors.Is(err, fsentry_error.ErrorNotExist) {
			return er.CollectionNotExists.AddMessage(err.Error()).HTTP(http.StatusBadRequest)
		} else if errors.Is(err, fsentry_error.ErrorBadName) {
			return er.BadName
		} else {
			return er.InternalError.AddMessage(err.Error())
		}
	}
	return nil
}
func (d *collection) ImageCreate(ctx context.Context, gameID, collectionID string, data []byte) error {
	collection, err := d.Get(ctx, gameID, collectionID)
	if err != nil {
		return err
	}

	err = d.db.CreateBinary("image", data, d.gamesPath, gameID, collection.ID)
	if err != nil {
		if errors.Is(err, fsentry_error.ErrorExist) {
			return er.CollectionImageExist.AddMessage(err.Error())
		} else {
			return er.InternalError.AddMessage(err.Error())
		}
	}
	return nil
}
func (d *collection) ImageGet(ctx context.Context, gameID, collectionID string) ([]byte, error) {
	collection, err := d.Get(ctx, gameID, collectionID)
	if err != nil {
		return nil, err
	}

	data, err := d.db.GetBinary("image", d.gamesPath, gameID, collection.ID)
	if err != nil {
		if errors.Is(err, fsentry_error.ErrorNotExist) {
			return nil, er.CollectionImageNotExists.AddMessage(err.Error())
		} else {
			return nil, er.InternalError.AddMessage(err.Error())
		}
	}
	return data, nil
}
func (d *collection) ImageDelete(ctx context.Context, gameID, collectionID string) error {
	collection, err := d.Get(ctx, gameID, collectionID)
	if err != nil {
		return err
	}

	err = d.db.RemoveBinary("image", d.gamesPath, gameID, collection.ID)
	if err != nil {
		if errors.Is(err, fsentry_error.ErrorNotExist) {
			return er.CollectionImageNotExists.AddMessage(err.Error())
		} else {
			return er.InternalError.AddMessage(err.Error())
		}
	}
	return nil
}

func (d *collection) convertCreateUpdate(createdAt, updatedAt *time.Time) (time.Time, time.Time) {
	if createdAt == nil {
		createdAt = utils.Allocate(time.Now())
	}
	if updatedAt == nil {
		updatedAt = createdAt
	}
	return *createdAt, *updatedAt
}
