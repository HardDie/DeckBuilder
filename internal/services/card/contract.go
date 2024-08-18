package card

import (
	"github.com/HardDie/DeckBuilder/internal/entity"
	"github.com/HardDie/DeckBuilder/internal/network"
)

type Card interface {
	Create(gameID, collectionID, deckID string, req CreateRequest) (*entity.CardInfo, error)
	Item(gameID, collectionID, deckID string, cardID int64) (*entity.CardInfo, error)
	List(gameID, collectionID, deckID, sortField, search string) ([]*entity.CardInfo, *network.Meta, error)
	Update(gameID, collectionID, deckID string, cardID int64, req UpdateRequest) (*entity.CardInfo, error)
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
