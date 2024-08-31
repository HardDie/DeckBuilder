package deck

import (
	"context"
	"time"
)

type Deck interface {
	Create(ctx context.Context, gameID, collectionID, name, description, image string) (*DeckInfo, error)
	Get(ctx context.Context, gameID, collectionID, name string) (context.Context, *DeckInfo, error)
	List(ctx context.Context, gameID, collectionID string) ([]*DeckInfo, error)
	Move(ctx context.Context, gameID, collectionID, oldName, newName string) (*DeckInfo, error)
	Update(ctx context.Context, gameID, collectionID, name, description, image string) (*DeckInfo, error)
	Delete(ctx context.Context, gameID, collectionID, name string) error
	ImageCreate(ctx context.Context, gameID, collectionID, deckID string, data []byte) error
	ImageGet(ctx context.Context, gameID, collectionID, deckID string) ([]byte, error)
	ImageDelete(ctx context.Context, gameID, collectionID, deckID string) error
}

type DeckInfo struct {
	ID          string     `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Image       string     `json:"image"`
	CachedImage string     `json:"cachedImage,omitempty"`
	CreatedAt   *time.Time `json:"createdAt"`
	UpdatedAt   *time.Time `json:"updatedAt"`

	// Dynamic fields

	GameID       string `json:"game_id"`
	CollectionID string `json:"collection_id"`
}
