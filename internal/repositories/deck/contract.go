package deck

import (
	entitiesDeck "github.com/HardDie/DeckBuilder/internal/entities/deck"
	"github.com/HardDie/DeckBuilder/internal/entity"
)

type Deck interface {
	Create(gameID, collectionID string, req CreateRequest) (*entitiesDeck.Deck, error)
	GetByID(gameID, collectionID, deckID string) (*entitiesDeck.Deck, error)
	GetAll(gameID, collectionID string) ([]*entity.DeckInfo, error)
	Update(gameID, collectionID, deckID string, req UpdateRequest) (*entitiesDeck.Deck, error)
	DeleteByID(gameID, collectionID, deckID string) error
	GetImage(gameID, collectionID, deckID string) ([]byte, string, error)
	GetAllDecksInGame(gameID string) ([]*entity.DeckInfo, error)
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
