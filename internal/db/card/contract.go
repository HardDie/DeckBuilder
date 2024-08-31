package card

import (
	"context"

	entitiesCard "github.com/HardDie/DeckBuilder/internal/entities/card"
)

type Card interface {
	Create(ctx context.Context, gameID, collectionID, deckID, name, description, image string, variables map[string]string, count int) (*entitiesCard.Card, error)
	Get(ctx context.Context, gameID, collectionID, deckID string, cardID int64) (*entitiesCard.Card, error)
	List(ctx context.Context, gameID, collectionID, deckID string) ([]*entitiesCard.Card, error)
	Update(ctx context.Context, gameID, collectionID, deckID string, cardID int64, name, description, image string, variables map[string]string, count int) (*entitiesCard.Card, error)
	Delete(ctx context.Context, gameID, collectionID, deckID string, cardID int64) error
	ImageCreate(ctx context.Context, gameID, collectionID, deckID string, cardID int64, data []byte) error
	ImageGet(ctx context.Context, gameID, collectionID, deckID string, cardID int64) ([]byte, error)
	ImageDelete(ctx context.Context, gameID, collectionID, deckID string, cardID int64) error
}
