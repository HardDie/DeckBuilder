package repository

import (
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

func (s *CardRepository) Create(gameID, collectionID, deckID string, req *dto.CreateCardDTO) (*entity.CardInfo, error) {
	card, err := s.db.CardCreate(gameID, collectionID, deckID, req.Name, req.Description, req.Image, req.Variables, req.Count)
	if err != nil {
		return nil, err
	}

	if card.Image == "" {
		return card, nil
	}

	// Download image
	if err := s.createImage(gameID, collectionID, deckID, card.ID, card.Image); err != nil {
		logger.Warn.Println("Unable to load image. The card will be saved without an image.", err.Error())
	}

	return card, nil
}
func (s *CardRepository) GetByID(gameID, collectionID, deckID string, cardID int64) (*entity.CardInfo, error) {
	return s.db.CardGet(gameID, collectionID, deckID, cardID)
}
func (s *CardRepository) GetAll(gameID, collectionID, deckID string) ([]*entity.CardInfo, error) {
	return s.db.CardList(gameID, collectionID, deckID)
}
func (s *CardRepository) Update(gameID, collectionID, deckID string, cardID int64, req *dto.UpdateCardDTO) (*entity.CardInfo, error) {
	oldCard, err := s.db.CardGet(gameID, collectionID, deckID, cardID)
	if err != nil {
		return nil, err
	}

	var newCard *entity.CardInfo
	if oldCard.Name != req.Name ||
		oldCard.Description != req.Description ||
		oldCard.Image != req.Image ||
		oldCard.Count != req.Count ||
		!utils.CompareMaps(oldCard.Variables, req.Variables) {
		// Update data
		newCard, err = s.db.CardUpdate(gameID, collectionID, deckID, cardID, req.Name, req.Description, req.Image, req.Variables, req.Count)
		if err != nil {
			return nil, err
		}
	}

	if newCard == nil {
		// If nothing has changed
		newCard = oldCard
	}

	// If the image has not been changed
	if newCard.Image == oldCard.Image {
		return newCard, nil
	}

	// If image exist, delete
	if data, _, _ := s.GetImage(gameID, collectionID, deckID, newCard.ID); data != nil {
		err = s.db.CardImageDelete(gameID, collectionID, deckID, cardID)
		if err != nil {
			return nil, err
		}
	}

	if newCard.Image == "" {
		return newCard, nil
	}

	// Download image
	if err = s.createImage(gameID, collectionID, deckID, newCard.ID, newCard.Image); err != nil {
		logger.Warn.Println("Unable to load image. The card will be saved without an image.", err.Error())
	}

	return newCard, nil
}
func (s *CardRepository) DeleteByID(gameID, collectionID, deckID string, cardID int64) error {
	err := s.db.CardImageDelete(gameID, collectionID, deckID, cardID)
	if err != nil {
		// Skip if image not exist
		if !errors.Is(err, er.CardImageNotExists) {
			return err
		}
	}
	return s.db.CardDelete(gameID, collectionID, deckID, cardID)
}
func (s *CardRepository) GetImage(gameID, collectionID, deckID string, cardID int64) ([]byte, string, error) {
	data, err := s.db.CardImageGet(gameID, collectionID, deckID, cardID)
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

	// Validate image
	_, err = images.ValidateImage(imageBytes)
	if err != nil {
		return err
	}

	// Write image to file
	return s.db.CardImageCreate(gameID, collectionID, deckID, cardID, imageBytes)
}
