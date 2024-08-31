package card

import (
	"context"
	"time"
)

type Card interface {
	Create(ctx context.Context, gameID, collectionID, deckID, name, description, image string, variables map[string]string, count int) (*CardInfo, error)
	Get(ctx context.Context, gameID, collectionID, deckID string, cardID int64) (context.Context, *CardInfo, error)
	List(ctx context.Context, gameID, collectionID, deckID string) ([]*CardInfo, error)
	Update(ctx context.Context, gameID, collectionID, deckID string, cardID int64, name, description, image string, variables map[string]string, count int) (*CardInfo, error)
	Delete(ctx context.Context, gameID, collectionID, deckID string, cardID int64) error
	ImageCreate(ctx context.Context, gameID, collectionID, deckID string, cardID int64, data []byte) error
	ImageGet(ctx context.Context, gameID, collectionID, deckID string, cardID int64) ([]byte, error)
	ImageDelete(ctx context.Context, gameID, collectionID, deckID string, cardID int64) error
}

type CardInfo struct {
	ID          int64             `json:"id"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Image       string            `json:"image"`
	CachedImage string            `json:"cachedImage,omitempty"`
	Variables   map[string]string `json:"variables"`
	Count       int               `json:"count"`
	CreatedAt   *time.Time        `json:"createdAt"`
	UpdatedAt   *time.Time        `json:"updatedAt"`

	// Dynamic fields

	GameID       string `json:"game_id"`
	CollectionID string `json:"collection_id"`
	DeckID       string `json:"deck_id"`
}
