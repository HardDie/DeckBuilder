package game

import (
	"context"

	entitiesGame "github.com/HardDie/DeckBuilder/internal/entities/game"
)

type Game interface {
	Create(ctx context.Context, name, description, image string) (*entitiesGame.Game, error)
	Get(ctx context.Context, name string) (*entitiesGame.Game, error)
	List(ctx context.Context) ([]*entitiesGame.Game, error)
	Move(ctx context.Context, oldName, newName string) (*entitiesGame.Game, error)
	Update(ctx context.Context, name, description, image string) (*entitiesGame.Game, error)
	Delete(ctx context.Context, name string) error
	Duplicate(ctx context.Context, srcName, dstName string) (*entitiesGame.Game, error)
	UpdateInfo(ctx context.Context, name, newName string) error
	ImageCreate(ctx context.Context, gameID string, data []byte) error
	ImageGet(ctx context.Context, gameID string) ([]byte, error)
	ImageDelete(ctx context.Context, gameID string) error
}
