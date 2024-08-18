package game

import (
	"context"
	"time"
)

type Game interface {
	Create(_ context.Context, name, description, image string) (*GameInfo, error)
	Get(ctx context.Context, name string) (context.Context, *GameInfo, error)
	List(ctx context.Context) ([]*GameInfo, error)
	Move(_ context.Context, oldName, newName string) (*GameInfo, error)
	Update(_ context.Context, name, description, image string) (*GameInfo, error)
	Delete(_ context.Context, name string) error
	Duplicate(_ context.Context, srcName, dstName string) (*GameInfo, error)
	UpdateInfo(_ context.Context, name, newName string) error
	ImageCreate(ctx context.Context, gameID string, data []byte) error
	ImageGet(ctx context.Context, gameID string) ([]byte, error)
	ImageDelete(ctx context.Context, gameID string) error
}

type GameInfo struct {
	ID          string     `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Image       string     `json:"image"`
	CreatedAt   *time.Time `json:"createdAt"`
	UpdatedAt   *time.Time `json:"updatedAt"`
}
