package card

import (
	"context"
	"errors"

	"github.com/HardDie/DeckBuilder/internal/config"
	dbCard "github.com/HardDie/DeckBuilder/internal/db/card"
	"github.com/HardDie/DeckBuilder/internal/entity"
	er "github.com/HardDie/DeckBuilder/internal/errors"
	"github.com/HardDie/DeckBuilder/internal/images"
	"github.com/HardDie/DeckBuilder/internal/logger"
	"github.com/HardDie/DeckBuilder/internal/network"
	"github.com/HardDie/DeckBuilder/internal/utils"
)

type card struct {
	cfg  *config.Config
	card dbCard.Card
}

func New(cfg *config.Config, c dbCard.Card) Card {
	return &card{
		cfg:  cfg,
		card: c,
	}
}

func (r *card) Create(gameID, collectionID, deckID string, req CreateRequest) (*entity.CardInfo, error) {
	c, err := r.card.Create(context.Background(), gameID, collectionID, deckID, req.Name,
		req.Description, req.Image, req.Variables, req.Count)
	if err != nil {
		return nil, err
	}

	if c.Image == "" && req.ImageFile == nil {
		return c, nil
	}

	if c.Image != "" {
		// Download image
		err = r.createImage(gameID, collectionID, deckID, c.ID, c.Image)
		if err != nil {
			logger.Warn.Println("Unable to load image. The card will be saved without an image.", err.Error())
		}
	} else if req.ImageFile != nil {
		err = r.createImageFromByte(gameID, collectionID, deckID, c.ID, req.ImageFile)
		if err != nil {
			logger.Warn.Println("Invalid image. The card will be saved without an image.", err.Error())
		}
	}

	return c, nil
}
func (r *card) GetByID(gameID, collectionID, deckID string, cardID int64) (*entity.CardInfo, error) {
	_, resp, err := r.card.Get(context.Background(), gameID, collectionID, deckID, cardID)
	return resp, err
}
func (r *card) GetAll(gameID, collectionID, deckID string) ([]*entity.CardInfo, error) {
	return r.card.List(context.Background(), gameID, collectionID, deckID)
}
func (r *card) Update(gameID, collectionID, deckID string, cardID int64, req UpdateRequest) (*entity.CardInfo, error) {
	_, oldCard, err := r.card.Get(context.Background(), gameID, collectionID, deckID, cardID)
	if err != nil {
		return nil, err
	}

	var newCard *entity.CardInfo
	if oldCard.Name != req.Name ||
		oldCard.Description != req.Description ||
		oldCard.Image != req.Image ||
		req.ImageFile != nil ||
		oldCard.Count != req.Count ||
		!utils.CompareMaps(oldCard.Variables, req.Variables) {
		// Update data
		newCard, err = r.card.Update(context.Background(), gameID, collectionID, deckID, cardID, req.Name, req.Description, req.Image, req.Variables, req.Count)
		if err != nil {
			return nil, err
		}
	}

	if newCard == nil {
		// If nothing has changed
		newCard = oldCard
	}

	// If the image has not been changed
	if newCard.Image == oldCard.Image && req.ImageFile == nil {
		return newCard, nil
	}

	// If image exist, delete
	if data, _, _ := r.GetImage(gameID, collectionID, deckID, newCard.ID); data != nil {
		err = r.card.ImageDelete(context.Background(), gameID, collectionID, deckID, cardID)
		if err != nil {
			return nil, err
		}
	}

	if newCard.Image == "" && req.ImageFile == nil {
		return newCard, nil
	}

	if newCard.Image != "" {
		// Download image
		if err = r.createImage(gameID, collectionID, deckID, newCard.ID, newCard.Image); err != nil {
			logger.Warn.Println("Unable to load image. The card will be saved without an image.", err.Error())
		}
	} else if req.ImageFile != nil {
		err = r.createImageFromByte(gameID, collectionID, deckID, newCard.ID, req.ImageFile)
		if err != nil {
			logger.Warn.Println("Invalid image. The card will be saved without an image.", err.Error())
		}
	}

	return newCard, nil
}
func (r *card) DeleteByID(gameID, collectionID, deckID string, cardID int64) error {
	err := r.card.ImageDelete(context.Background(), gameID, collectionID, deckID, cardID)
	if err != nil {
		// Skip if image not exist
		if !errors.Is(err, er.CardImageNotExists) {
			return err
		}
	}
	return r.card.Delete(context.Background(), gameID, collectionID, deckID, cardID)
}
func (r *card) GetImage(gameID, collectionID, deckID string, cardID int64) ([]byte, string, error) {
	data, err := r.card.ImageGet(context.Background(), gameID, collectionID, deckID, cardID)
	if err != nil {
		return nil, "", err
	}

	imgType, err := images.ValidateImage(data)
	if err != nil {
		return nil, "", err
	}

	return data, imgType, nil
}

func (r *card) createImage(gameID, collectionID, deckID string, cardID int64, imageURL string) error {
	// Download image
	imageBytes, err := network.DownloadBytes(imageURL)
	if err != nil {
		return err
	}

	return r.createImageFromByte(gameID, collectionID, deckID, cardID, imageBytes)
}
func (r *card) createImageFromByte(gameID, collectionID, deckID string, cardID int64, data []byte) error {
	// Validate image
	_, err := images.ValidateImage(data)
	if err != nil {
		return err
	}

	// Write image to file
	return r.card.ImageCreate(context.Background(), gameID, collectionID, deckID, cardID, data)
}
