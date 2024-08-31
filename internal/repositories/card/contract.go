package card

import (
	entitiesCard "github.com/HardDie/DeckBuilder/internal/entities/card"
)

type Card interface {
	Create(gameID, collectionID, deckID string, req CreateRequest) (*entitiesCard.Card, error)
	GetByID(gameID, collectionID, deckID string, cardID int64) (*entitiesCard.Card, error)
	GetAll(gameID, collectionID, deckID string) ([]*entitiesCard.Card, error)
	Update(gameID, collectionID, deckID string, cardID int64, req UpdateRequest) (*entitiesCard.Card, error)
	DeleteByID(gameID, collectionID, deckID string, cardID int64) error
	GetImage(gameID, collectionID, deckID string, cardID int64) ([]byte, string, error)
}

type CreateRequest struct {
	Name        string
	Description string
	Image       string
	Variables   map[string]string
	Count       int
	ImageFile   []byte
}

type UpdateRequest struct {
	Name        string
	Description string
	Image       string
	Variables   map[string]string
	Count       int
	ImageFile   []byte
}
