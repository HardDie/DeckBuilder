package card

import (
	"github.com/HardDie/DeckBuilder/internal/entity"
)

type Card interface {
	Create(gameID, collectionID, deckID string, req CreateRequest) (*entity.CardInfo, error)
	GetByID(gameID, collectionID, deckID string, cardID int64) (*entity.CardInfo, error)
	GetAll(gameID, collectionID, deckID string) ([]*entity.CardInfo, error)
	Update(gameID, collectionID, deckID string, cardID int64, req UpdateRequest) (*entity.CardInfo, error)
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
