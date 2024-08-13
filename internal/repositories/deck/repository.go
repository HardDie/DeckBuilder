package deck

import (
	"context"

	"github.com/HardDie/DeckBuilder/internal/config"
	dbCollection "github.com/HardDie/DeckBuilder/internal/db/collection"
	dbDeck "github.com/HardDie/DeckBuilder/internal/db/deck"
	"github.com/HardDie/DeckBuilder/internal/dto"
	"github.com/HardDie/DeckBuilder/internal/entity"
	"github.com/HardDie/DeckBuilder/internal/images"
	"github.com/HardDie/DeckBuilder/internal/logger"
	"github.com/HardDie/DeckBuilder/internal/network"
)

type deck struct {
	cfg        *config.Config
	collection dbCollection.Collection
	deck       dbDeck.Deck
}

func New(cfg *config.Config, c dbCollection.Collection, d dbDeck.Deck) Deck {
	return &deck{
		cfg:        cfg,
		collection: c,
		deck:       d,
	}
}

func (r *deck) Create(gameID, collectionID string, dtoObject *dto.CreateDeckDTO) (*entity.DeckInfo, error) {
	deck, err := r.deck.Create(context.Background(), gameID, collectionID, dtoObject.Name, dtoObject.Description, dtoObject.Image)
	if err != nil {
		return nil, err
	}

	if deck.Image == "" && dtoObject.ImageFile == nil {
		return deck, nil
	}

	if deck.Image != "" {
		// Download image
		err = r.createImage(gameID, collectionID, deck.ID, deck.Image)
		if err != nil {
			logger.Warn.Println("Unable to load image. The deck will be saved without an image.", err.Error())
		}
	} else if dtoObject.ImageFile != nil {
		err = r.createImageFromByte(gameID, collectionID, deck.ID, dtoObject.ImageFile)
		if err != nil {
			logger.Warn.Println("Invalid image. The deck will be saved without an image.", err.Error())
		}
	}

	return deck, nil
}
func (r *deck) GetByID(gameID, collectionID, deckID string) (*entity.DeckInfo, error) {
	_, resp, err := r.deck.Get(context.Background(), gameID, collectionID, deckID)
	return resp, err
}
func (r *deck) GetAll(gameID, collectionID string) ([]*entity.DeckInfo, error) {
	return r.deck.List(context.Background(), gameID, collectionID)
}
func (r *deck) Update(gameID, collectionID, deckID string, dtoObject *dto.UpdateDeckDTO) (*entity.DeckInfo, error) {
	_, oldDeck, err := r.deck.Get(context.Background(), gameID, collectionID, deckID)
	if err != nil {
		return nil, err
	}

	var newDeck *entity.DeckInfo
	if oldDeck.Name != dtoObject.Name {
		// Rename folder
		newDeck, err = r.deck.Move(context.Background(), gameID, collectionID, oldDeck.Name, dtoObject.Name)
		if err != nil {
			return nil, err
		}
	}

	if oldDeck.Description != dtoObject.Description ||
		oldDeck.Image != dtoObject.Image ||
		dtoObject.ImageFile != nil {
		// Update data
		newDeck, err = r.deck.Update(context.Background(), gameID, collectionID, dtoObject.Name, dtoObject.Description, dtoObject.Image)
		if err != nil {
			return nil, err
		}
	}

	if newDeck == nil {
		// If nothing has changed
		newDeck = oldDeck
	}

	// If the image has not been changed
	if newDeck.Image == oldDeck.Image && dtoObject.ImageFile == nil {
		return newDeck, nil
	}

	// If image exist, delete
	if data, _, _ := r.GetImage(gameID, collectionID, newDeck.ID); data != nil {
		err = r.deck.ImageDelete(context.Background(), gameID, collectionID, deckID)
		if err != nil {
			return nil, err
		}
	}

	if newDeck.Image == "" && dtoObject.ImageFile == nil {
		return newDeck, nil
	}

	if newDeck.Image != "" {
		// Download image
		err = r.createImage(gameID, collectionID, newDeck.ID, newDeck.Image)
		if err != nil {
			logger.Warn.Println("Unable to load image. The deck will be saved without an image.", err.Error())
		}
	} else if dtoObject.ImageFile != nil {
		err = r.createImageFromByte(gameID, collectionID, newDeck.ID, dtoObject.ImageFile)
		if err != nil {
			logger.Warn.Println("Invalid image. The deck will be saved without an image.", err.Error())
		}
	}

	return newDeck, nil
}
func (r *deck) DeleteByID(gameID, collectionID, deckID string) error {
	return r.deck.Delete(context.Background(), gameID, collectionID, deckID)
}
func (r *deck) GetImage(gameID, collectionID, deckID string) ([]byte, string, error) {
	data, err := r.deck.ImageGet(context.Background(), gameID, collectionID, deckID)
	if err != nil {
		return nil, "", err
	}

	imgType, err := images.ValidateImage(data)
	if err != nil {
		return nil, "", err
	}

	return data, imgType, nil
}
func (r *deck) GetAllDecksInGame(gameID string) ([]*entity.DeckInfo, error) {
	// Get all collections in selected game
	listCollections, err := r.collection.List(context.Background(), gameID)
	if err != nil {
		return make([]*entity.DeckInfo, 0), err
	}

	// Mark unique deck types
	uniqueDecks := make(map[string]struct{})

	// Go through all collections and find unique types of decks
	decks := make([]*entity.DeckInfo, 0)
	for _, collection := range listCollections {
		// Get all decks in selected collection
		collectionDecks, err := r.GetAll(gameID, collection.ID)
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
			deck.FillCachedImage(r.cfg, gameID, collection.ID)
			decks = append(decks, deck)
		}
	}
	return decks, nil
}

func (r *deck) createImage(gameID, collectionID, deckID, imageURL string) error {
	// Download image
	imageBytes, err := network.DownloadBytes(imageURL)
	if err != nil {
		return err
	}

	return r.createImageFromByte(gameID, collectionID, deckID, imageBytes)
}
func (r *deck) createImageFromByte(gameID, collectionID, deckID string, data []byte) error {
	// Validate image
	_, err := images.ValidateImage(data)
	if err != nil {
		return err
	}

	// Write image to file
	return r.deck.ImageCreate(context.Background(), gameID, collectionID, deckID, data)
}