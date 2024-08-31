package card

import (
	entitiesCard "github.com/HardDie/DeckBuilder/internal/entities/card"
)

type Card interface {
	Create(gameID, collectionID, deckID string, req CreateRequest) (*entitiesCard.Card, error)
	Item(gameID, collectionID, deckID string, cardID int64) (*entitiesCard.Card, error)
	List(gameID, collectionID, deckID, sortField, search string) ([]*entitiesCard.Card, error)
	Update(gameID, collectionID, deckID string, cardID int64, req UpdateRequest) (*entitiesCard.Card, error)
	Delete(gameID, collectionID, deckID string, cardID int64) error
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
