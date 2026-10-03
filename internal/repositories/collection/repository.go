package collection

import (
	"errors"

	"github.com/HardDie/fsentry"

	"github.com/HardDie/DeckBuilder/internal/config"
	entitiesCollection "github.com/HardDie/DeckBuilder/internal/entities/collection"
	er "github.com/HardDie/DeckBuilder/internal/errors"
	"github.com/HardDie/DeckBuilder/internal/images"
	"github.com/HardDie/DeckBuilder/internal/logger"
	"github.com/HardDie/DeckBuilder/internal/repositories"
	"github.com/HardDie/DeckBuilder/internal/utils"
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
	change, imageErr := repositories.ResolveImage("", req.Image, req.ImageFile)
	req.Image = change.URL

	item, err := r.create(gameID, req)
	if err != nil {
		return nil, err
	}
	if change.Data != nil {
		if err = r.imageCreate(gameID, item.ID, change.Data); err != nil {
			imageErr = err
		}
	}
	item.ImageError = imageErr
	return item, nil
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
	// Download and validate before any write, so a bad image keeps the old one.
	change, imageErr := repositories.ResolveImage(oldCollection.Image, req.Image, req.ImageFile)

	newCollection := oldCollection
	if oldCollection.Name != req.Name {
		newCollection, err = r.move(gameID, oldCollection.Name, req.Name)
		if err != nil {
			return nil, err
		}
	}

	if oldCollection.Description != req.Description ||
		oldCollection.Image != change.URL ||
		change.Data != nil {
		newCollection, err = r.update(gameID, updateRequest{
			Name:        req.Name,
			Description: req.Description,
			Image:       change.URL,
		})
		if err != nil {
			return nil, err
		}
	}

	if change.Data != nil || change.Clear {
		err = r.imageDelete(gameID, newCollection.ID)
		if err != nil && !errors.Is(err, er.CollectionImageNotExists) {
			return nil, err
		}
	}
	if change.Data != nil {
		if err = r.imageCreate(gameID, newCollection.ID, change.Data); err != nil {
			imageErr = err
		}
	}
	newCollection.ImageError = imageErr
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

	imgType, err := images.ImageType(data)
	if err != nil {
		return nil, "", err
	}

	return data, imgType, nil
}

func (r *collection) create(gameID string, req CreateRequest) (*entitiesCollection.Collection, error) {
	info, err := r.db.CreateFolder(req.Name, model{
		Description: fsentry.QuotedString(req.Description),
		Image:       fsentry.QuotedString(req.Image),
	}, r.gamesPath, gameID)
	if err != nil {
		if mapped := er.MissingAncestor(err); mapped != nil {
			return nil, mapped
		}
		return nil, repositories.MapFsentry(err, repositories.FsentrySentinels{
			Exist: er.CollectionExist,
		})
	}

	return r.toEntity(info, gameID), nil
}

func (r *collection) get(gameID, name string) (*entitiesCollection.Collection, error) {
	info, err := r.db.GetFolder[model](name, r.gamesPath, gameID)
	if err != nil {
		if mapped := er.MissingAncestor(err); mapped != nil {
			return nil, mapped
		}
		return nil, repositories.MapFsentry(err, repositories.FsentrySentinels{
			NotExist: er.CollectionNotExists,
			Message:  true,
		})
	}

	return r.toEntity(info, gameID), nil
}

func (r *collection) list(gameID string) ([]*entitiesCollection.Collection, error) {
	list, err := r.db.List(r.gamesPath, gameID)
	if err != nil {
		if mapped := er.MissingAncestor(err); mapped != nil {
			return nil, mapped
		}
		return nil, er.InternalError.AddMessage(err.Error())
	}

	var collections []*entitiesCollection.Collection
	for _, folder := range list.Folders {
		collection, err := r.get(gameID, folder)
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
	info, err := r.db.MoveFolder[model](oldName, newName, r.gamesPath, gameID)
	if err != nil {
		if mapped := er.MissingAncestor(err); mapped != nil {
			return nil, mapped
		}
		return nil, repositories.MapFsentry(err, repositories.FsentrySentinels{
			NotExist: er.CollectionNotExists,
			Message:  true,
		})
	}

	return r.toEntity(info, gameID), nil
}

type updateRequest struct {
	Name        string
	Description string
	Image       string
}

func (r *collection) update(gameID string, req updateRequest) (*entitiesCollection.Collection, error) {
	info, err := r.db.UpdateFolder(req.Name, model{
		Description: fsentry.QuotedString(req.Description),
		Image:       fsentry.QuotedString(req.Image),
	}, r.gamesPath, gameID)
	if err != nil {
		if mapped := er.MissingAncestor(err); mapped != nil {
			return nil, mapped
		}
		return nil, repositories.MapFsentry(err, repositories.FsentrySentinels{
			NotExist: er.CollectionNotExists,
			Message:  true,
		})
	}

	return r.toEntity(info, gameID), nil
}

func (r *collection) delete(gameID, name string) error {
	err := r.db.RemoveFolder(name, r.gamesPath, gameID)
	if err != nil {
		if mapped := er.MissingAncestor(err); mapped != nil {
			return mapped
		}
		return repositories.MapFsentry(err, repositories.FsentrySentinels{
			NotExist: er.CollectionNotExists,
			Message:  true,
		})
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
		if mapped := er.MissingAncestor(err); mapped != nil {
			return mapped
		}
		return repositories.MapFsentry(err, repositories.FsentrySentinels{
			Exist:   er.CollectionImageExist,
			Message: true,
		})
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
		if mapped := er.MissingAncestor(err); mapped != nil {
			return nil, mapped
		}
		return nil, repositories.MapFsentry(err, repositories.FsentrySentinels{
			NotExist: er.CollectionImageNotExists,
			Message:  true,
		})
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
		if mapped := er.MissingAncestor(err); mapped != nil {
			return mapped
		}
		return repositories.MapFsentry(err, repositories.FsentrySentinels{
			NotExist: er.CollectionImageNotExists,
			Message:  true,
		})
	}
	return nil
}

func (r *collection) toEntity(info fsentry.FolderInfo[model], gameID string) *entitiesCollection.Collection {
	createdAt, updatedAt := utils.NormalizeTimestamps(info.CreatedAt, info.UpdatedAt)
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
