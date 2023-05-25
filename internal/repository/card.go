package repository

import (
	"context"
	"errors"

	"github.com/HardDie/DeckBuilder/internal/config"
	"github.com/HardDie/DeckBuilder/internal/db"
	"github.com/HardDie/DeckBuilder/internal/dto"
	"github.com/HardDie/DeckBuilder/internal/entity"
	er "github.com/HardDie/DeckBuilder/internal/errors"
	"github.com/HardDie/DeckBuilder/internal/images"
	"github.com/HardDie/DeckBuilder/internal/logger"
	"github.com/HardDie/DeckBuilder/internal/network"
	"github.com/HardDie/DeckBuilder/internal/utils"
)

type ICardRepository interface {
	Create(gameID, collectionID, deckID string, req *dto.CreateCardDTO) (*entity.CardInfo, error)
	GetByID(gameID, collectionID, deckID string, cardID int64) (*entity.CardInfo, error)
	GetAll(gameID, collectionID, deckID string) ([]*entity.CardInfo, error)
	Update(gameID, collectionID, deckID string, cardID int64, dtoObject *dto.UpdateCardDTO) (*entity.CardInfo, error)
	DeleteByID(gameID, collectionID, deckID string, cardID int64) error
	GetImage(gameID, collectionID, deckID string, cardID int64) ([]byte, string, error)
}
type CardRepository struct {
	cfg *config.Config
	db  *db.DB
}

func NewCardRepository(cfg *config.Config, db *db.DB) *CardRepository {
	return &CardRepository{
		cfg: cfg,
		db:  db,
	}
}

func (s *CardRepository) Create(gameID, collectionID, deckID string, dtoObject *dto.CreateCardDTO) (*entity.CardInfo, error) {
	card, err := s.db.CardCreate(context.Background(), gameID, collectionID, deckID, dtoObject.Name,
		dtoObject.Description, dtoObject.Image, dtoObject.Variables, dtoObject.Count)
	if err != nil {
		return nil, err
	}

	if card.Image == "" && dtoObject.ImageFile == nil {
		return card, nil
	}

	if card.Image != "" {
		// Download image
		err = s.createImage(gameID, collectionID, deckID, card.ID, card.Image)
		if err != nil {
			logger.Warn.Println("Unable to load image. The card will be saved without an image.", err.Error())
		}
	} else if dtoObject.ImageFile != nil {
		err = s.createImageFromByte(gameID, collectionID, deckID, card.ID, dtoObject.ImageFile)
		if err != nil {
			logger.Warn.Println("Invalid image. The card will be saved without an image.", err.Error())
		}
	}

	return card, nil
}
func (s *CardRepository) GetByID(gameID, collectionID, deckID string, cardID int64) (*entity.CardInfo, error) {
	_, resp, err := s.db.CardGet(context.Background(), gameID, collectionID, deckID, cardID)
	return resp, err
}
func (s *CardRepository) GetAll(gameID, collectionID, deckID string) ([]*entity.CardInfo, error) {
	return s.db.CardList(context.Background(), gameID, collectionID, deckID)
}
func (s *CardRepository) Update(gameID, collectionID, deckID string, cardID int64, dtoObject *dto.UpdateCardDTO) (*entity.CardInfo, error) {
	_, oldCard, err := s.db.CardGet(context.Background(), gameID, collectionID, deckID, cardID)
	if err != nil {
		return nil, err
	}

	var newCard *entity.CardInfo
	if oldCard.Name != dtoObject.Name ||
		oldCard.Description != dtoObject.Description ||
		oldCard.Image != dtoObject.Image ||
		dtoObject.ImageFile != nil ||
		oldCard.Count != dtoObject.Count ||
		!utils.CompareMaps(oldCard.Variables, dtoObject.Variables) {
		// Update data
		newCard, err = s.db.CardUpdate(context.Background(), gameID, collectionID, deckID, cardID, dtoObject.Name, dtoObject.Description, dtoObject.Image, dtoObject.Variables, dtoObject.Count)
		if err != nil {
			return nil, err
		}
	}

	if newCard == nil {
		// If nothing has changed
		newCard = oldCard
	}

	// If the image has not been changed
	if newCard.Image == oldCard.Image && dtoObject.ImageFile == nil {
		return newCard, nil
	}

	// If image exist, delete
	if data, _, _ := s.GetImage(gameID, collectionID, deckID, newCard.ID); data != nil {
		err = s.db.CardImageDelete(context.Background(), gameID, collectionID, deckID, cardID)
		if err != nil {
			return nil, err
		}
	}

	if newCard.Image == "" && dtoObject.ImageFile == nil {
		return newCard, nil
	}

	if newCard.Image != "" {
		// Download image
		if err = s.createImage(gameID, collectionID, deckID, newCard.ID, newCard.Image); err != nil {
			logger.Warn.Println("Unable to load image. The card will be saved without an image.", err.Error())
		}
	} else if dtoObject.ImageFile != nil {
		err = s.createImageFromByte(gameID, collectionID, deckID, newCard.ID, dtoObject.ImageFile)
		if err != nil {
			logger.Warn.Println("Invalid image. The card will be saved without an image.", err.Error())
		}
	}

	return newCard, nil
}
func (s *CardRepository) DeleteByID(gameID, collectionID, deckID string, cardID int64) error {
	err := s.db.CardImageDelete(context.Background(), gameID, collectionID, deckID, cardID)
	if err != nil {
		// Skip if image not exist
		if !errors.Is(err, er.CardImageNotExists) {
			return err
		}
	}
	return s.db.CardDelete(context.Background(), gameID, collectionID, deckID, cardID)
}
func (s *CardRepository) GetImage(gameID, collectionID, deckID string, cardID int64) ([]byte, string, error) {
	data, err := s.db.CardImageGet(context.Background(), gameID, collectionID, deckID, cardID)
	if err != nil {
		return nil, "", err
	}

	imgType, err := images.ValidateImage(data)
	if err != nil {
		return nil, "", err
	}

	return data, imgType, nil
}

func (s *CardRepository) createImage(gameID, collectionID, deckID string, cardID int64, imageURL string) error {
	// Download image
	imageBytes, err := network.DownloadBytes(imageURL)
	if err != nil {
		return err
	}

	return s.createImageFromByte(gameID, collectionID, deckID, cardID, imageBytes)
}
func (s *CardRepository) createImageFromByte(gameID, collectionID, deckID string, cardID int64, data []byte) error {
	// Validate image
	_, err := images.ValidateImage(data)
	if err != nil {
		return err
	}

	// Write image to file
	return s.db.CardImageCreate(context.Background(), gameID, collectionID, deckID, cardID, data)
}
