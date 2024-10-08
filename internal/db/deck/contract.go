package deck

import (
	"context"

	entitiesDeck "github.com/HardDie/DeckBuilder/internal/entities/deck"
)

type Deck interface {
	Create(ctx context.Context, req CreateRequest) (*entitiesDeck.Deck, error)
	Get(ctx context.Context, gameID, collectionID, name string) (*entitiesDeck.Deck, error)
	List(ctx context.Context, gameID, collectionID string) ([]*entitiesDeck.Deck, error)
	Move(ctx context.Context, gameID, collectionID, oldName, newName string) (*entitiesDeck.Deck, error)
	Update(ctx context.Context, req UpdateRequest) (*entitiesDeck.Deck, error)
	Delete(ctx context.Context, gameID, collectionID, name string) error
	ImageCreate(ctx context.Context, gameID, collectionID, deckID string, data []byte) error
	ImageGet(ctx context.Context, gameID, collectionID, deckID string) ([]byte, error)
	ImageDelete(ctx context.Context, gameID, collectionID, deckID string) error
}

type CreateRequest struct {
	GameID       string
	CollectionID string
	Name         string
	Description  string
	Image        string
}

type UpdateRequest struct {
	// Select
	GameID       string
	CollectionID string
	// Update
	Name        string
	Description string
	Image       string
}
