package repository

import (
	"github.com/HardDie/DeckBuilder/internal/config"
	"github.com/HardDie/DeckBuilder/internal/db"
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
	cfg *config.Config
	db  *db.DB
}

func NewCollectionRepository(cfg *config.Config, db *db.DB) *CollectionRepository {
	return &CollectionRepository{
		cfg: cfg,
		db:  db,
	}
}

func (s *CollectionRepository) Create(gameID string, req *dto.CreateCollectionDTO) (*entity.CollectionInfo, error) {
	collection, err := s.db.CollectionCreate(gameID, req.Name, req.Description, req.Image)
	if err != nil {
		return nil, err
	}

	if collection.Image == "" {
		return collection, nil
	}

	// Download image
	if err := s.createImage(gameID, collection.ID, collection.Image); err != nil {
		logger.Warn.Println("Unable to load image. The collection will be saved without an image.", err.Error())
	}

	return collection, nil
}
func (s *CollectionRepository) GetByID(gameID, collectionID string) (*entity.CollectionInfo, error) {
	return s.db.CollectionGet(gameID, collectionID)
}
func (s *CollectionRepository) GetAll(gameID string) ([]*entity.CollectionInfo, error) {
	return s.db.CollectionList(gameID)
}
func (s *CollectionRepository) Update(gameID, collectionID string, dtoObject *dto.UpdateCollectionDTO) (*entity.CollectionInfo, error) {
	oldCollection, err := s.db.CollectionGet(gameID, collectionID)
	if err != nil {
		return nil, err
	}

	var newCollection *entity.CollectionInfo
	if oldCollection.Name != dtoObject.Name {
		// Rename folder
		newCollection, err = s.db.CollectionMove(gameID, oldCollection.Name, dtoObject.Name)
		if err != nil {
			return nil, err
		}
	}

	if oldCollection.Description != dtoObject.Description ||
		oldCollection.Image != dtoObject.Image {
		// Update data
		newCollection, err = s.db.CollectionUpdate(gameID, dtoObject.Name, dtoObject.Description, dtoObject.Image)
		if err != nil {
			return nil, err
		}
	}

	if newCollection == nil {
		// If nothing has changed
		newCollection = oldCollection
	}

	// If the image has not been changed
	if newCollection.Image == oldCollection.Image {
		return newCollection, nil
	}

	// If image exist, delete
	if data, _, _ := s.GetImage(gameID, newCollection.ID); data != nil {
		err = s.db.CollectionImageDelete(gameID, collectionID)
		if err != nil {
			return nil, err
		}
	}

	if newCollection.Image == "" {
		return newCollection, nil
	}

	// Download image
	if err = s.createImage(gameID, newCollection.ID, newCollection.Image); err != nil {
		logger.Warn.Println("Unable to load image. The collection will be saved without an image.", err.Error())
	}

	return newCollection, nil
}
func (s *CollectionRepository) DeleteByID(gameID, collectionID string) error {
	return s.db.CollectionDelete(gameID, collectionID)
}
func (s *CollectionRepository) GetImage(gameID, collectionID string) ([]byte, string, error) {
	data, err := s.db.CollectionImageGet(gameID, collectionID)
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

	// Validate image
	_, err = images.ValidateImage(imageBytes)
	if err != nil {
		return err
	}

	// Write image to file
	return s.db.CollectionImageCreate(gameID, collectionID, imageBytes)
}
