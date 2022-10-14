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

type IDeckRepository interface {
	Create(gameID, collectionID string, req *dto.CreateDeckDTO) (*entity.DeckInfo, error)
	GetByID(gameID, collectionID, deckID string) (*entity.DeckInfo, error)
	GetAll(gameID, collectionID string) ([]*entity.DeckInfo, error)
	Update(gameID, collectionID, deckID string, dtoObject *dto.UpdateDeckDTO) (*entity.DeckInfo, error)
	DeleteByID(gameID, collectionID, deckID string) error
	GetImage(gameID, collectionID, deckID string) ([]byte, string, error)
	GetAllDecksInGame(gameID string) ([]*entity.DeckInfo, error)
}
type DeckRepository struct {
	cfg *config.Config
	db  *db.DB
}

func NewDeckRepository(cfg *config.Config, db *db.DB) *DeckRepository {
	return &DeckRepository{
		cfg: cfg,
		db:  db,
	}
}

func (s *DeckRepository) Create(gameID, collectionID string, req *dto.CreateDeckDTO) (*entity.DeckInfo, error) {
	deck, err := s.db.DeckCreate(gameID, collectionID, req.Name, req.Description, req.Image)
	if err != nil {
		return nil, err
	}

	if deck.Image == "" {
		return deck, nil
	}

	// Download image
	if err := s.createImage(gameID, collectionID, deck.ID, deck.Image); err != nil {
		logger.Warn.Println("Unable to load image. The deck will be saved without an image.", err.Error())
	}

	return deck, nil
}
func (s *DeckRepository) GetByID(gameID, collectionID, deckID string) (*entity.DeckInfo, error) {
	return s.db.DeckGet(gameID, collectionID, deckID)
}
func (s *DeckRepository) GetAll(gameID, collectionID string) ([]*entity.DeckInfo, error) {
	return s.db.DeckList(gameID, collectionID)
}
func (s *DeckRepository) Update(gameID, collectionID, deckID string, dtoObject *dto.UpdateDeckDTO) (*entity.DeckInfo, error) {
	oldDeck, err := s.db.DeckGet(gameID, collectionID, deckID)
	if err != nil {
		return nil, err
	}

	var newDeck *entity.DeckInfo
	if oldDeck.Name != dtoObject.Name {
		// Rename folder
		newDeck, err = s.db.DeckMove(gameID, collectionID, oldDeck.Name, dtoObject.Name)
		if err != nil {
			return nil, err
		}
	}

	if oldDeck.Description != dtoObject.Description ||
		oldDeck.Image != dtoObject.Image {
		// Update data
		newDeck, err = s.db.DeckUpdate(gameID, collectionID, dtoObject.Name, dtoObject.Description, dtoObject.Image)
		if err != nil {
			return nil, err
		}
	}

	if newDeck == nil {
		// If nothing has changed
		newDeck = oldDeck
	}

	// If the image has not been changed
	if newDeck.Image == oldDeck.Image {
		return newDeck, nil
	}

	// If image exist, delete
	if data, _, _ := s.GetImage(gameID, collectionID, newDeck.ID); data != nil {
		err = s.db.DeckImageDelete(gameID, collectionID, deckID)
		if err != nil {
			return nil, err
		}
	}

	if newDeck.Image == "" {
		return newDeck, nil
	}

	// Download image
	if err = s.createImage(gameID, collectionID, newDeck.ID, newDeck.Image); err != nil {
		logger.Warn.Println("Unable to load image. The deck will be saved without an image.", err.Error())
	}

	return newDeck, nil
}
func (s *DeckRepository) DeleteByID(gameID, collectionID, deckID string) error {
	return s.db.DeckDelete(gameID, collectionID, deckID)
}
func (s *DeckRepository) GetImage(gameID, collectionID, deckID string) ([]byte, string, error) {
	data, err := s.db.DeckImageGet(gameID, collectionID, deckID)
	if err != nil {
		return nil, "", err
	}

	imgType, err := images.ValidateImage(data)
	if err != nil {
		return nil, "", err
	}

	return data, imgType, nil
}
func (s *DeckRepository) GetAllDecksInGame(gameID string) ([]*entity.DeckInfo, error) {
	// Get all collections in selected game
	listCollections, err := s.db.CollectionList(gameID)
	if err != nil {
		return make([]*entity.DeckInfo, 0), err
	}

	// Mark unique deck types
	uniqueDecks := make(map[string]struct{})

	// Go through all collections and find unique types of decks
	decks := make([]*entity.DeckInfo, 0)
	for _, collection := range listCollections {
		// Get all decks in selected collection
		collectionDecks, err := s.GetAll(gameID, collection.ID)
		if err != nil {
			return make([]*entity.DeckInfo, 0), err
		}

		// Go through all decks and keep only unique decks
		for _, deck := range collectionDecks {
			if _, ok := uniqueDecks[deck.Name+deck.Image]; ok {
				// If we have already seen such a deck, we skip it
				continue
			}
			// If deck unique, put mark in map
			uniqueDecks[deck.Name+deck.Image] = struct{}{}
			deck.FillCachedImage(s.cfg, gameID, collection.ID)
			decks = append(decks, deck)
		}
	}
	return decks, nil
}

func (s *DeckRepository) createImage(gameID, collectionID, deckID, imageURL string) error {
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
	return s.db.DeckImageCreate(gameID, collectionID, deckID, imageBytes)
}
