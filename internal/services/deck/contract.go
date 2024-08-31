package deck

import (
	entitiesDeck "github.com/HardDie/DeckBuilder/internal/entities/deck"
)

type Deck interface {
	Create(gameID, collectionID string, req CreateRequest) (*entitiesDeck.Deck, error)
	Item(gameID, collectionID, deckID string) (*entitiesDeck.Deck, error)
	List(gameID, collectionID, sortField, search string) ([]*entitiesDeck.Deck, error)
	Update(gameID, collectionID, deckID string, req UpdateRequest) (*entitiesDeck.Deck, error)
	Delete(gameID, collectionID, deckID string) error
	GetImage(gameID, collectionID, deckID string) ([]byte, string, error)
	ListAllUnique(gameID string) ([]*entitiesDeck.Deck, error)
}

type CreateRequest struct {
	Name        string
	Description string
	Image       string
	ImageFile   []byte
}

type UpdateRequest struct {
	Name        string
	Description string
	Image       string
	ImageFile   []byte
}
