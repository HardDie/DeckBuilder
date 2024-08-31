package card

import (
	"context"

	entitiesCard "github.com/HardDie/DeckBuilder/internal/entities/card"
)

type Card interface {
	Create(ctx context.Context, req CreateRequest) (*entitiesCard.Card, error)
	Get(ctx context.Context, gameID, collectionID, deckID string, cardID int64) (*entitiesCard.Card, error)
	List(ctx context.Context, gameID, collectionID, deckID string) ([]*entitiesCard.Card, error)
	Update(ctx context.Context, req UpdateRequest) (*entitiesCard.Card, error)
	Delete(ctx context.Context, gameID, collectionID, deckID string, cardID int64) error
	ImageCreate(ctx context.Context, gameID, collectionID, deckID string, cardID int64, data []byte) error
	ImageGet(ctx context.Context, gameID, collectionID, deckID string, cardID int64) ([]byte, error)
	ImageDelete(ctx context.Context, gameID, collectionID, deckID string, cardID int64) error
}

type CreateRequest struct {
	GameID       string
	CollectionID string
	DeckID       string
	Name         string
	Description  string
	Image        string
	Variables    map[string]string
	Count        int
}

type UpdateRequest struct {
	// Select
	GameID       string
	CollectionID string
	DeckID       string
	CardID       int64
	// Update
	Name        string
	Description string
	Image       string
	Variables   map[string]string
	Count       int
}
