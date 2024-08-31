package collection

import (
	"context"
	"time"
)

type Collection interface {
	Create(ctx context.Context, gameID, name, description, image string) (*CollectionInfo, error)
	Get(ctx context.Context, gameID, name string) (context.Context, *CollectionInfo, error)
	List(ctx context.Context, gameID string) ([]*CollectionInfo, error)
	Move(ctx context.Context, gameID, oldName, newName string) (*CollectionInfo, error)
	Update(ctx context.Context, gameID, name, description, image string) (*CollectionInfo, error)
	Delete(ctx context.Context, gameID, name string) error
	ImageCreate(ctx context.Context, gameID, collectionID string, data []byte) error
	ImageGet(ctx context.Context, gameID, collectionID string) ([]byte, error)
	ImageDelete(ctx context.Context, gameID, collectionID string) error
}

type CollectionInfo struct {
	ID          string     `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Image       string     `json:"image"`
	CreatedAt   *time.Time `json:"createdAt"`
	UpdatedAt   *time.Time `json:"updatedAt"`

	// Dynamic fields

	GameID string `json:"game_id"`
}
