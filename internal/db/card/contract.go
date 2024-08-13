package card

import (
	"context"

	"github.com/HardDie/DeckBuilder/internal/entity"
)

type Card interface {
	Create(ctx context.Context, gameID, collectionID, deckID, name, description, image string, variables map[string]string, count int) (*entity.CardInfo, error)
	Get(ctx context.Context, gameID, collectionID, deckID string, cardID int64) (context.Context, *entity.CardInfo, error)
	List(ctx context.Context, gameID, collectionID, deckID string) ([]*entity.CardInfo, error)
	Update(ctx context.Context, gameID, collectionID, deckID string, cardID int64, name, description, image string, variables map[string]string, count int) (*entity.CardInfo, error)
	Delete(ctx context.Context, gameID, collectionID, deckID string, cardID int64) error
	ImageCreate(ctx context.Context, gameID, collectionID, deckID string, cardID int64, data []byte) error
	ImageGet(ctx context.Context, gameID, collectionID, deckID string, cardID int64) ([]byte, error)
	ImageDelete(ctx context.Context, gameID, collectionID, deckID string, cardID int64) error
}
