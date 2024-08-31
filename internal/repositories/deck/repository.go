package deck

import (
	"context"
	"time"

	"github.com/HardDie/DeckBuilder/internal/config"
	dbCollection "github.com/HardDie/DeckBuilder/internal/db/collection"
	dbDeck "github.com/HardDie/DeckBuilder/internal/db/deck"
	entitiesDeck "github.com/HardDie/DeckBuilder/internal/entities/deck"
	"github.com/HardDie/DeckBuilder/internal/images"
	"github.com/HardDie/DeckBuilder/internal/logger"
	"github.com/HardDie/DeckBuilder/internal/network"
	"github.com/HardDie/DeckBuilder/internal/utils"
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

func (r *deck) Create(gameID, collectionID string, req CreateRequest) (*entitiesDeck.Deck, error) {
	d, err := r.deck.Create(context.Background(), gameID, collectionID, req.Name, req.Description, req.Image)
	if err != nil {
		return nil, err
	}

	if d.Image == "" && req.ImageFile == nil {
		return r.oldEntityToNew(d), nil
	}

	if d.Image != "" {
		// Download image
		err = r.createImage(gameID, collectionID, d.ID, d.Image)
		if err != nil {
			logger.Warn.Println("Unable to load image. The deck will be saved without an image.", err.Error())
		}
	} else if req.ImageFile != nil {
		err = r.createImageFromByte(gameID, collectionID, d.ID, req.ImageFile)
		if err != nil {
			logger.Warn.Println("Invalid image. The deck will be saved without an image.", err.Error())
		}
	}

	return r.oldEntityToNew(d), nil
}
func (r *deck) GetByID(gameID, collectionID, deckID string) (*entitiesDeck.Deck, error) {
	_, resp, err := r.deck.Get(context.Background(), gameID, collectionID, deckID)
	return r.oldEntityToNew(resp), err
}
func (r *deck) GetAll(gameID, collectionID string) ([]*entitiesDeck.Deck, error) {
	items, err := r.deck.List(context.Background(), gameID, collectionID)
	if err != nil {
		return nil, err
	}
	res := make([]*entitiesDeck.Deck, 0, len(items))
	for _, item := range items {
		res = append(res, r.oldEntityToNew(item))
	}
	return res, nil
}
func (r *deck) Update(gameID, collectionID, deckID string, req UpdateRequest) (*entitiesDeck.Deck, error) {
	_, oldDeck, err := r.deck.Get(context.Background(), gameID, collectionID, deckID)
	if err != nil {
		return nil, err
	}

	var newDeck *dbDeck.DeckInfo
	if oldDeck.Name != req.Name {
		// Rename folder
		newDeck, err = r.deck.Move(context.Background(), gameID, collectionID, oldDeck.Name, req.Name)
		if err != nil {
			return nil, err
		}
	}

	if oldDeck.Description != req.Description ||
		oldDeck.Image != req.Image ||
		req.ImageFile != nil {
		// Update data
		newDeck, err = r.deck.Update(context.Background(), gameID, collectionID, req.Name, req.Description, req.Image)
		if err != nil {
			return nil, err
		}
	}

	if newDeck == nil {
		// If nothing has changed
		newDeck = oldDeck
	}

	// If the image has not been changed
	if newDeck.Image == oldDeck.Image && req.ImageFile == nil {
		return r.oldEntityToNew(newDeck), nil
	}

	// If image exist, delete
	if data, _, _ := r.GetImage(gameID, collectionID, newDeck.ID); data != nil {
		err = r.deck.ImageDelete(context.Background(), gameID, collectionID, deckID)
		if err != nil {
			return nil, err
		}
	}

	if newDeck.Image == "" && req.ImageFile == nil {
		return r.oldEntityToNew(newDeck), nil
	}

	if newDeck.Image != "" {
		// Download image
		err = r.createImage(gameID, collectionID, newDeck.ID, newDeck.Image)
		if err != nil {
			logger.Warn.Println("Unable to load image. The deck will be saved without an image.", err.Error())
		}
	} else if req.ImageFile != nil {
		err = r.createImageFromByte(gameID, collectionID, newDeck.ID, req.ImageFile)
		if err != nil {
			logger.Warn.Println("Invalid image. The deck will be saved without an image.", err.Error())
		}
	}

	return r.oldEntityToNew(newDeck), nil
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
func (r *deck) GetAllDecksInGame(gameID string) ([]*entitiesDeck.Deck, error) {
	// Get all collections in selected game
	listCollections, err := r.collection.List(context.Background(), gameID)
	if err != nil {
		return make([]*entitiesDeck.Deck, 0), err
	}

	// Mark unique deck types
	uniqueDecks := make(map[string]struct{})

	// Go through all collections and find unique types of decks
	decks := make([]*entitiesDeck.Deck, 0)
	for _, collection := range listCollections {
		// Get all decks in selected collection
		collectionDecks, err := r.GetAll(gameID, collection.ID)
		if err != nil {
			return make([]*entitiesDeck.Deck, 0), err
		}

		// Go through all decks and keep only unique decks
		for _, d := range collectionDecks {
			if _, ok := uniqueDecks[d.Name+d.Image]; ok {
				// If we have already seen such a deck, we skip it
				continue
			}
			// If deck unique, put mark in map
			uniqueDecks[d.Name+d.Image] = struct{}{}
			decks = append(decks, d)
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

func (r *deck) oldEntityToNew(g *dbDeck.DeckInfo) *entitiesDeck.Deck {
	if g == nil {
		return nil
	}
	createdAt, updatedAt := r.convertCreateUpdate(g.CreatedAt, g.UpdatedAt)
	return &entitiesDeck.Deck{
		ID:          g.ID,
		Name:        g.Name,
		Description: g.Description,
		Image:       g.Image,
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,

		GameID:       g.GameID,
		CollectionID: g.CollectionID,
	}
}
func (r *deck) convertCreateUpdate(createdAt, updatedAt *time.Time) (time.Time, time.Time) {
	if createdAt == nil {
		createdAt = utils.Allocate(time.Now())
	}
	if updatedAt == nil {
		updatedAt = createdAt
	}
	return *createdAt, *updatedAt
}
