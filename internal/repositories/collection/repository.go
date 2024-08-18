package collection

import (
	"context"

	"github.com/HardDie/DeckBuilder/internal/config"
	dbCollection "github.com/HardDie/DeckBuilder/internal/db/collection"
	"github.com/HardDie/DeckBuilder/internal/entity"
	"github.com/HardDie/DeckBuilder/internal/images"
	"github.com/HardDie/DeckBuilder/internal/logger"
	"github.com/HardDie/DeckBuilder/internal/network"
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

func (r *collection) Create(gameID string, req CreateRequest) (*entity.CollectionInfo, error) {
	c, err := r.collection.Create(context.Background(), gameID, req.Name, req.Description, req.Image)
	if err != nil {
		return nil, err
	}

	if c.Image == "" && req.ImageFile == nil {
		return c, nil
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

	return c, nil
}
func (r *collection) GetByID(gameID, collectionID string) (*entity.CollectionInfo, error) {
	_, resp, err := r.collection.Get(context.Background(), gameID, collectionID)
	return resp, err
}
func (r *collection) GetAll(gameID string) ([]*entity.CollectionInfo, error) {
	return r.collection.List(context.Background(), gameID)
}
func (r *collection) Update(gameID, collectionID string, req UpdateRequest) (*entity.CollectionInfo, error) {
	_, oldCollection, err := r.collection.Get(context.Background(), gameID, collectionID)
	if err != nil {
		return nil, err
	}

	var newCollection *entity.CollectionInfo
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
		return newCollection, nil
	}

	// If image exist, delete
	if data, _, _ := r.GetImage(gameID, newCollection.ID); data != nil {
		err = r.collection.ImageDelete(context.Background(), gameID, collectionID)
		if err != nil {
			return nil, err
		}
	}

	if newCollection.Image == "" && req.ImageFile == nil {
		return newCollection, nil
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

	return newCollection, nil
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
