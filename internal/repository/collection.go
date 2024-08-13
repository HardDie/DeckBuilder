package repository

import (
	"context"

	"github.com/HardDie/DeckBuilder/internal/config"
	dbCollection "github.com/HardDie/DeckBuilder/internal/db/collection"
	"github.com/HardDie/DeckBuilder/internal/dto"
	"github.com/HardDie/DeckBuilder/internal/entity"
	"github.com/HardDie/DeckBuilder/internal/images"
	"github.com/HardDie/DeckBuilder/internal/logger"
	"github.com/HardDie/DeckBuilder/internal/network"
)

type ICollectionRepository interface {
	Create(gameID string, req *dto.CreateCollectionDTO) (*entity.CollectionInfo, error)
	GetByID(gameID, collectionID string) (*entity.CollectionInfo, error)
	GetAll(gameID string) ([]*entity.CollectionInfo, error)
	Update(gameID, collectionID string, dtoObject *dto.UpdateCollectionDTO) (*entity.CollectionInfo, error)
	DeleteByID(gameID, collectionID string) error
	GetImage(gameID, collectionID string) ([]byte, string, error)
}
type CollectionRepository struct {
	cfg        *config.Config
	collection dbCollection.Collection
}

func NewCollectionRepository(cfg *config.Config, collection dbCollection.Collection) *CollectionRepository {
	return &CollectionRepository{
		cfg:        cfg,
		collection: collection,
	}
}

func (s *CollectionRepository) Create(gameID string, dtoObject *dto.CreateCollectionDTO) (*entity.CollectionInfo, error) {
	collection, err := s.collection.Create(context.Background(), gameID, dtoObject.Name, dtoObject.Description, dtoObject.Image)
	if err != nil {
		return nil, err
	}

	if collection.Image == "" && dtoObject.ImageFile == nil {
		return collection, nil
	}

	if collection.Image != "" {
		// Download image
		err = s.createImage(gameID, collection.ID, collection.Image)
		if err != nil {
			logger.Warn.Println("Unable to load image. The collection will be saved without an image.", err.Error())
		}
	} else if dtoObject.ImageFile != nil {
		err = s.createImageFromByte(gameID, collection.ID, dtoObject.ImageFile)
		if err != nil {
			logger.Warn.Println("Invalid image. The collection will be saved without an image.", err.Error())
		}
	}

	return collection, nil
}
func (s *CollectionRepository) GetByID(gameID, collectionID string) (*entity.CollectionInfo, error) {
	_, resp, err := s.collection.Get(context.Background(), gameID, collectionID)
	return resp, err
}
func (s *CollectionRepository) GetAll(gameID string) ([]*entity.CollectionInfo, error) {
	return s.collection.List(context.Background(), gameID)
}
func (s *CollectionRepository) Update(gameID, collectionID string, dtoObject *dto.UpdateCollectionDTO) (*entity.CollectionInfo, error) {
	_, oldCollection, err := s.collection.Get(context.Background(), gameID, collectionID)
	if err != nil {
		return nil, err
	}

	var newCollection *entity.CollectionInfo
	if oldCollection.Name != dtoObject.Name {
		// Rename folder
		newCollection, err = s.collection.Move(context.Background(), gameID, oldCollection.Name, dtoObject.Name)
		if err != nil {
			return nil, err
		}
	}

	if oldCollection.Description != dtoObject.Description ||
		oldCollection.Image != dtoObject.Image ||
		dtoObject.ImageFile != nil {
		// Update data
		newCollection, err = s.collection.Update(context.Background(), gameID, dtoObject.Name, dtoObject.Description, dtoObject.Image)
		if err != nil {
			return nil, err
		}
	}

	if newCollection == nil {
		// If nothing has changed
		newCollection = oldCollection
	}

	// If the image has not been changed
	if newCollection.Image == oldCollection.Image && dtoObject.ImageFile == nil {
		return newCollection, nil
	}

	// If image exist, delete
	if data, _, _ := s.GetImage(gameID, newCollection.ID); data != nil {
		err = s.collection.ImageDelete(context.Background(), gameID, collectionID)
		if err != nil {
			return nil, err
		}
	}

	if newCollection.Image == "" && dtoObject.ImageFile == nil {
		return newCollection, nil
	}

	if newCollection.Image != "" {
		// Download image
		err = s.createImage(gameID, newCollection.ID, newCollection.Image)
		if err != nil {
			logger.Warn.Println("Unable to load image. The collection will be saved without an image.", err.Error())
		}
	} else if dtoObject.ImageFile != nil {
		err = s.createImageFromByte(gameID, newCollection.ID, dtoObject.ImageFile)
		if err != nil {
			logger.Warn.Println("Invalid image. The collection will be saved without an image.", err.Error())
		}
	}

	return newCollection, nil
}
func (s *CollectionRepository) DeleteByID(gameID, collectionID string) error {
	return s.collection.Delete(context.Background(), gameID, collectionID)
}
func (s *CollectionRepository) GetImage(gameID, collectionID string) ([]byte, string, error) {
	data, err := s.collection.ImageGet(context.Background(), gameID, collectionID)
	if err != nil {
		return nil, "", err
	}

	imgType, err := images.ValidateImage(data)
	if err != nil {
		return nil, "", err
	}

	return data, imgType, nil
}

func (s *CollectionRepository) createImage(gameID, collectionID, imageURL string) error {
	// Download image
	imageBytes, err := network.DownloadBytes(imageURL)
	if err != nil {
		return err
	}

	return s.createImageFromByte(gameID, collectionID, imageBytes)
}
func (s *CollectionRepository) createImageFromByte(gameID, collectionID string, data []byte) error {
	// Validate image
	_, err := images.ValidateImage(data)
	if err != nil {
		return err
	}

	// Write image to file
	return s.collection.ImageCreate(context.Background(), gameID, collectionID, data)
}
