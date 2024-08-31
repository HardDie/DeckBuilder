package collection

import (
	"context"
	"time"

	"github.com/HardDie/DeckBuilder/internal/config"
	dbCollection "github.com/HardDie/DeckBuilder/internal/db/collection"
	entitiesCollection "github.com/HardDie/DeckBuilder/internal/entities/collection"
	"github.com/HardDie/DeckBuilder/internal/images"
	"github.com/HardDie/DeckBuilder/internal/logger"
	"github.com/HardDie/DeckBuilder/internal/network"
	"github.com/HardDie/DeckBuilder/internal/utils"
)

type collection struct {
	cfg        *config.Config
	collection dbCollection.Collection
}

func New(cfg *config.Config, c dbCollection.Collection) Collection {
	return &collection{
		cfg:        cfg,
		collection: c,
	}
}

func (r *collection) Create(gameID string, req CreateRequest) (*entitiesCollection.Collection, error) {
	c, err := r.collection.Create(context.Background(), gameID, req.Name, req.Description, req.Image)
	if err != nil {
		return nil, err
	}

	if c.Image == "" && req.ImageFile == nil {
		return r.oldEntityToNew(c), nil
	}

	if c.Image != "" {
		// Download image
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

	return r.oldEntityToNew(c), nil
}
func (r *collection) GetByID(gameID, collectionID string) (*entitiesCollection.Collection, error) {
	_, resp, err := r.collection.Get(context.Background(), gameID, collectionID)
	return r.oldEntityToNew(resp), err
}
func (r *collection) GetAll(gameID string) ([]*entitiesCollection.Collection, error) {
	items, err := r.collection.List(context.Background(), gameID)
	if err != nil {
		return nil, err
	}
	res := make([]*entitiesCollection.Collection, 0, len(items))
	for _, item := range items {
		res = append(res, r.oldEntityToNew(item))
	}
	return res, nil
}
func (r *collection) Update(gameID, collectionID string, req UpdateRequest) (*entitiesCollection.Collection, error) {
	_, oldCollection, err := r.collection.Get(context.Background(), gameID, collectionID)
	if err != nil {
		return nil, err
	}

	var newCollection *dbCollection.CollectionInfo
	if oldCollection.Name != req.Name {
		// Rename folder
		newCollection, err = r.collection.Move(context.Background(), gameID, oldCollection.Name, req.Name)
		if err != nil {
			return nil, err
		}
	}

	if oldCollection.Description != req.Description ||
		oldCollection.Image != req.Image ||
		req.ImageFile != nil {
		// Update data
		newCollection, err = r.collection.Update(context.Background(), gameID, req.Name, req.Description, req.Image)
		if err != nil {
			return nil, err
		}
	}

	if newCollection == nil {
		// If nothing has changed
		newCollection = oldCollection
	}

	// If the image has not been changed
	if newCollection.Image == oldCollection.Image && req.ImageFile == nil {
		return r.oldEntityToNew(newCollection), nil
	}

	// If image exist, delete
	if data, _, _ := r.GetImage(gameID, newCollection.ID); data != nil {
		err = r.collection.ImageDelete(context.Background(), gameID, collectionID)
		if err != nil {
			return nil, err
		}
	}

	if newCollection.Image == "" && req.ImageFile == nil {
		return r.oldEntityToNew(newCollection), nil
	}

	if newCollection.Image != "" {
		// Download image
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

	return r.oldEntityToNew(newCollection), nil
}
func (r *collection) DeleteByID(gameID, collectionID string) error {
	return r.collection.Delete(context.Background(), gameID, collectionID)
}
func (r *collection) GetImage(gameID, collectionID string) ([]byte, string, error) {
	data, err := r.collection.ImageGet(context.Background(), gameID, collectionID)
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
	// Download image
	imageBytes, err := network.DownloadBytes(imageURL)
	if err != nil {
		return err
	}

	return r.createImageFromByte(gameID, collectionID, imageBytes)
}
func (r *collection) createImageFromByte(gameID, collectionID string, data []byte) error {
	// Validate image
	_, err := images.ValidateImage(data)
	if err != nil {
		return err
	}

	// Write image to file
	return r.collection.ImageCreate(context.Background(), gameID, collectionID, data)
}

func (r *collection) oldEntityToNew(g *dbCollection.CollectionInfo) *entitiesCollection.Collection {
	if g == nil {
		return nil
	}
	createdAt, updatedAt := r.convertCreateUpdate(g.CreatedAt, g.UpdatedAt)
	return &entitiesCollection.Collection{
		ID:          g.ID,
		Name:        g.Name,
		Description: g.Description,
		Image:       g.Image,
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,

		GameID: g.GameID,
	}
}
func (r *collection) convertCreateUpdate(createdAt, updatedAt *time.Time) (time.Time, time.Time) {
	if createdAt == nil {
		createdAt = utils.Allocate(time.Now())
	}
	if updatedAt == nil {
		updatedAt = createdAt
	}
	return *createdAt, *updatedAt
}
