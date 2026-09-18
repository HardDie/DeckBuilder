package collection

import (
	"errors"
	"net/http"
	"time"

	"github.com/HardDie/fsentry"

	"github.com/HardDie/DeckBuilder/internal/config"
	entitiesCollection "github.com/HardDie/DeckBuilder/internal/entities/collection"
	er "github.com/HardDie/DeckBuilder/internal/errors"
	"github.com/HardDie/DeckBuilder/internal/images"
	"github.com/HardDie/DeckBuilder/internal/logger"
	"github.com/HardDie/DeckBuilder/internal/network"
)

type collection struct {
	cfg       *config.Config
	db        *fsentry.DB
	gamesPath string
}

func New(cfg *config.Config, db *fsentry.DB) Collection {
	return &collection{
		cfg:       cfg,
		db:        db,
		gamesPath: "games",
	}
}

func (r *collection) Create(gameID string, req CreateRequest) (*entitiesCollection.Collection, error) {
	c, err := r.create(gameID, req)
	if err != nil {
		return nil, err
	}

	if c.Image == "" && req.ImageFile == nil {
		return c, nil
	}

	if c.Image != "" {
		err = r.createImage(gameID, c.ID, c.Image)
		if err != nil {
			logger.Warn.Println("Unable to load image. The collection will be saved without an image.", err.Error())
		}
	} else if req.ImageFile != nil {
		err = r.createImageFromByte(gameID, c.ID, req.ImageFile)
		if err != nil {
			logger.Warn.Println("Invalid image. The collection will be saved without an image.", err.Error())
		}
	}

	return c, nil
}

func (r *collection) GetByID(gameID, collectionID string) (*entitiesCollection.Collection, error) {
	return r.get(gameID, collectionID)
}

func (r *collection) GetAll(gameID string) ([]*entitiesCollection.Collection, error) {
	return r.list(gameID)
}

func (r *collection) Update(gameID, collectionID string, req UpdateRequest) (*entitiesCollection.Collection, error) {
	oldCollection, err := r.get(gameID, collectionID)
	if err != nil {
		return nil, err
	}

	var newCollection *entitiesCollection.Collection
	if oldCollection.Name != req.Name {
		newCollection, err = r.move(gameID, oldCollection.Name, req.Name)
		if err != nil {
			return nil, err
		}
	}

	if oldCollection.Description != req.Description ||
		oldCollection.Image != req.Image ||
		req.ImageFile != nil {
		newCollection, err = r.update(gameID, updateRequest{
			Name:        req.Name,
			Description: req.Description,
			Image:       req.Image,
		})
		if err != nil {
			return nil, err
		}
	}

	if newCollection == nil {
		newCollection = oldCollection
	}

	if newCollection.Image == oldCollection.Image && req.ImageFile == nil {
		return newCollection, nil
	}

	if data, _, _ := r.GetImage(gameID, newCollection.ID); data != nil {
		err = r.imageDelete(gameID, collectionID)
		if err != nil {
			return nil, err
		}
	}

	if newCollection.Image == "" && req.ImageFile == nil {
		return newCollection, nil
	}

	if newCollection.Image != "" {
		err = r.createImage(gameID, newCollection.ID, newCollection.Image)
		if err != nil {
			logger.Warn.Println("Unable to load image. The collection will be saved without an image.", err.Error())
		}
	} else if req.ImageFile != nil {
		err = r.createImageFromByte(gameID, newCollection.ID, req.ImageFile)
		if err != nil {
			logger.Warn.Println("Invalid image. The collection will be saved without an image.", err.Error())
		}
	}

	return newCollection, nil
}

func (r *collection) DeleteByID(gameID, collectionID string) error {
	return r.delete(gameID, collectionID)
}

func (r *collection) GetImage(gameID, collectionID string) ([]byte, string, error) {
	data, err := r.imageGet(gameID, collectionID)
	if err != nil {
		return nil, "", err
	}

	imgType, err := images.ValidateImage(data)
	if err != nil {
		return nil, "", err
	}

	return data, imgType, nil
}

func (r *collection) createImage(gameID, collectionID, imageURL string) error {
	imageBytes, err := network.DownloadBytes(imageURL)
	if err != nil {
		return err
	}

	return r.createImageFromByte(gameID, collectionID, imageBytes)
}

func (r *collection) createImageFromByte(gameID, collectionID string, data []byte) error {
	_, err := images.ValidateImage(data)
	if err != nil {
		return err
	}

	return r.imageCreate(gameID, collectionID, data)
}

func (r *collection) getGame(gameID string) (string, error) {
	info, err := r.db.GetFolder[any](gameID, r.gamesPath)
	if err != nil {
		if errors.Is(err, fsentry.ErrNotExist) {
			return "", er.GameNotExists.AddMessage(err.Error()).HTTP(http.StatusBadRequest)
		} else if errors.Is(err, fsentry.ErrBadName) {
			return "", er.BadName
		} else {
			return "", er.InternalError.AddMessage(err.Error())
		}
	}
	return info.ID, nil
}

func (r *collection) create(gameID string, req CreateRequest) (*entitiesCollection.Collection, error) {
	gameIDResolved, err := r.getGame(gameID)
	if err != nil {
		return nil, err
	}

	info, err := r.db.CreateFolder(req.Name, model{
		Description: fsentry.QuotedString(req.Description),
		Image:       fsentry.QuotedString(req.Image),
	}, r.gamesPath, gameIDResolved)
	if err != nil {
		if errors.Is(err, fsentry.ErrExist) {
			return nil, er.CollectionExist
		} else if errors.Is(err, fsentry.ErrBadName) {
			return nil, er.BadName
		} else {
			return nil, er.InternalError.AddMessage(err.Error())
		}
	}

	return r.toEntity(info, gameID), nil
}

func (r *collection) get(gameID, name string) (*entitiesCollection.Collection, error) {
	gameIDResolved, err := r.getGame(gameID)
	if err != nil {
		return nil, err
	}

	info, err := r.db.GetFolder[model](name, r.gamesPath, gameIDResolved)
	if err != nil {
		if errors.Is(err, fsentry.ErrNotExist) {
			return nil, er.CollectionNotExists.AddMessage(err.Error()).HTTP(http.StatusBadRequest)
		} else if errors.Is(err, fsentry.ErrBadName) {
			return nil, er.BadName
		} else {
			return nil, er.InternalError.AddMessage(err.Error())
		}
	}

	return r.toEntity(info, gameID), nil
}

func (r *collection) list(gameID string) ([]*entitiesCollection.Collection, error) {
	gameIDResolved, err := r.getGame(gameID)
	if err != nil {
		return nil, err
	}

	list, err := r.db.List(r.gamesPath, gameIDResolved)
	if err != nil {
		return nil, er.InternalError.AddMessage(err.Error())
	}

	var collections []*entitiesCollection.Collection
	for _, folder := range list.Folders {
		collection, err := r.get(gameIDResolved, folder)
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

func (r *collection) move(gameID, oldName, newName string) (*entitiesCollection.Collection, error) {
	gameIDResolved, err := r.getGame(gameID)
	if err != nil {
		return nil, err
	}

	info, err := r.db.MoveFolder[model](oldName, newName, r.gamesPath, gameIDResolved)
	if err != nil {
		if errors.Is(err, fsentry.ErrNotExist) {
			return nil, er.CollectionNotExists.AddMessage(err.Error()).HTTP(http.StatusBadRequest)
		} else if errors.Is(err, fsentry.ErrBadName) {
			return nil, er.BadName
		} else {
			return nil, er.InternalError.AddMessage(err.Error())
		}
	}

	return r.toEntity(info, gameID), nil
}

type updateRequest struct {
	Name        string
	Description string
	Image       string
}

func (r *collection) update(gameID string, req updateRequest) (*entitiesCollection.Collection, error) {
	gameIDResolved, err := r.getGame(gameID)
	if err != nil {
		return nil, err
	}

	info, err := r.db.UpdateFolder(req.Name, model{
		Description: fsentry.QuotedString(req.Description),
		Image:       fsentry.QuotedString(req.Image),
	}, r.gamesPath, gameIDResolved)
	if err != nil {
		if errors.Is(err, fsentry.ErrNotExist) {
			return nil, er.CollectionNotExists.AddMessage(err.Error()).HTTP(http.StatusBadRequest)
		} else if errors.Is(err, fsentry.ErrBadName) {
			return nil, er.BadName
		} else {
			return nil, er.InternalError.AddMessage(err.Error())
		}
	}

	return r.toEntity(info, gameID), nil
}

func (r *collection) delete(gameID, name string) error {
	gameIDResolved, err := r.getGame(gameID)
	if err != nil {
		return err
	}

	err = r.db.RemoveFolder(name, r.gamesPath, gameIDResolved)
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

func (r *collection) imageCreate(gameID, collectionID string, data []byte) error {
	collection, err := r.get(gameID, collectionID)
	if err != nil {
		return err
	}

	err = r.db.CreateBinary("image", data, r.gamesPath, gameID, collection.ID)
	if err != nil {
		if errors.Is(err, fsentry.ErrExist) {
			return er.CollectionImageExist.AddMessage(err.Error())
		} else {
			return er.InternalError.AddMessage(err.Error())
		}
	}
	return nil
}

func (r *collection) imageGet(gameID, collectionID string) ([]byte, error) {
	collection, err := r.get(gameID, collectionID)
	if err != nil {
		return nil, err
	}

	data, err := r.db.GetBinary("image", nil, r.gamesPath, gameID, collection.ID)
	if err != nil {
		if errors.Is(err, fsentry.ErrNotExist) {
			return nil, er.CollectionImageNotExists.AddMessage(err.Error())
		} else {
			return nil, er.InternalError.AddMessage(err.Error())
		}
	}
	return data, nil
}

func (r *collection) imageDelete(gameID, collectionID string) error {
	collection, err := r.get(gameID, collectionID)
	if err != nil {
		return err
	}

	err = r.db.RemoveBinary("image", r.gamesPath, gameID, collection.ID)
	if err != nil {
		if errors.Is(err, fsentry.ErrNotExist) {
			return er.CollectionImageNotExists.AddMessage(err.Error())
		} else {
			return er.InternalError.AddMessage(err.Error())
		}
	}
	return nil
}

func (r *collection) toEntity(info fsentry.FolderInfo[model], gameID string) *entitiesCollection.Collection {
	createdAt, updatedAt := r.convertCreateUpdate(info.CreatedAt, info.UpdatedAt)
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

func (r *collection) convertCreateUpdate(createdAt, updatedAt time.Time) (time.Time, time.Time) {
	if createdAt.IsZero() {
		createdAt = time.Now()
	}
	if updatedAt.IsZero() {
		updatedAt = createdAt
	}
	return createdAt, updatedAt
}
