package collection

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/HardDie/fsentry"

	dbGame "github.com/HardDie/DeckBuilder/internal/db/game"
	entitiesCollection "github.com/HardDie/DeckBuilder/internal/entities/collection"
	er "github.com/HardDie/DeckBuilder/internal/errors"
	"github.com/HardDie/DeckBuilder/internal/logger"
)

type collection struct {
	db        *fsentry.DB
	gamesPath string

	game dbGame.Game
}

func New(db *fsentry.DB, game dbGame.Game) Collection {
	return &collection{
		db:        db,
		gamesPath: "games",

		game: game,
	}
}

func (d *collection) Create(ctx context.Context, req CreateRequest) (*entitiesCollection.Collection, error) {
	game, err := d.game.Get(ctx, req.GameID)
	if err != nil {
		return nil, err
	}

	info, err := d.db.CreateFolder(req.Name, model{
		Description: fsentry.QuotedString(req.Description),
		Image:       fsentry.QuotedString(req.Image),
	}, d.gamesPath, game.ID)
	if err != nil {
		if errors.Is(err, fsentry.ErrExist) {
			return nil, er.CollectionExist
		} else if errors.Is(err, fsentry.ErrBadName) {
			return nil, er.BadName
		} else {
			return nil, er.InternalError.AddMessage(err.Error())
		}
	}

	return d.toEntity(info, req.GameID), nil
}

func (d *collection) Get(ctx context.Context, gameID, name string) (*entitiesCollection.Collection, error) {
	game, err := d.game.Get(ctx, gameID)
	if err != nil {
		return nil, err
	}

	info, err := d.db.GetFolder[model](name, d.gamesPath, game.ID)
	if err != nil {
		if errors.Is(err, fsentry.ErrNotExist) {
			return nil, er.CollectionNotExists.AddMessage(err.Error()).HTTP(http.StatusBadRequest)
		} else if errors.Is(err, fsentry.ErrBadName) {
			return nil, er.BadName
		} else {
			return nil, er.InternalError.AddMessage(err.Error())
		}
	}

	return d.toEntity(info, gameID), nil
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

	info, err := d.db.MoveFolder[model](oldName, newName, d.gamesPath, game.ID)
	if err != nil {
		if errors.Is(err, fsentry.ErrNotExist) {
			return nil, er.CollectionNotExists.AddMessage(err.Error()).HTTP(http.StatusBadRequest)
		} else if errors.Is(err, fsentry.ErrBadName) {
			return nil, er.BadName
		} else {
			return nil, er.InternalError.AddMessage(err.Error())
		}
	}

	return d.toEntity(info, gameID), nil
}

func (d *collection) Update(ctx context.Context, req UpdateRequest) (*entitiesCollection.Collection, error) {
	game, err := d.game.Get(ctx, req.GameID)
	if err != nil {
		return nil, err
	}

	info, err := d.db.UpdateFolder(req.Name, model{
		Description: fsentry.QuotedString(req.Description),
		Image:       fsentry.QuotedString(req.Image),
	}, d.gamesPath, game.ID)
	if err != nil {
		if errors.Is(err, fsentry.ErrNotExist) {
			return nil, er.CollectionNotExists.AddMessage(err.Error()).HTTP(http.StatusBadRequest)
		} else if errors.Is(err, fsentry.ErrBadName) {
			return nil, er.BadName
		} else {
			return nil, er.InternalError.AddMessage(err.Error())
		}
	}

	return d.toEntity(info, req.GameID), nil
}

func (d *collection) Delete(ctx context.Context, gameID, name string) error {
	game, err := d.game.Get(ctx, gameID)
	if err != nil {
		return err
	}

	err = d.db.RemoveFolder(name, d.gamesPath, game.ID)
	if err != nil {
		if errors.Is(err, fsentry.ErrNotExist) {
			return er.CollectionNotExists.AddMessage(err.Error()).HTTP(http.StatusBadRequest)
		} else if errors.Is(err, fsentry.ErrBadName) {
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
		if errors.Is(err, fsentry.ErrExist) {
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

	data, err := d.db.GetBinary("image", nil, d.gamesPath, gameID, collection.ID)
	if err != nil {
		if errors.Is(err, fsentry.ErrNotExist) {
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
		if errors.Is(err, fsentry.ErrNotExist) {
			return er.CollectionImageNotExists.AddMessage(err.Error())
		} else {
			return er.InternalError.AddMessage(err.Error())
		}
	}
	return nil
}

func (d *collection) toEntity(info fsentry.FolderInfo[model], gameID string) *entitiesCollection.Collection {
	createdAt, updatedAt := d.convertCreateUpdate(info.CreatedAt, info.UpdatedAt)
	return &entitiesCollection.Collection{
		ID:          info.ID,
		Name:        info.Name,
		Description: info.Data.Description.String(),
		Image:       info.Data.Image.String(),
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
		GameID:      gameID,
	}
}

func (d *collection) convertCreateUpdate(createdAt, updatedAt time.Time) (time.Time, time.Time) {
	if createdAt.IsZero() {
		createdAt = time.Now()
	}
	if updatedAt.IsZero() {
		updatedAt = createdAt
	}
	return createdAt, updatedAt
}
