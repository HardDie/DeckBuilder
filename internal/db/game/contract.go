package game

import (
	"context"

	"github.com/HardDie/DeckBuilder/internal/entity"
)

type Game interface {
	Create(_ context.Context, name, description, image string) (*entity.GameInfo, error)
	Get(ctx context.Context, name string) (context.Context, *entity.GameInfo, error)
	List(ctx context.Context) ([]*entity.GameInfo, error)
	Move(_ context.Context, oldName, newName string) (*entity.GameInfo, error)
	Update(_ context.Context, name, description, image string) (*entity.GameInfo, error)
	Delete(_ context.Context, name string) error
	Duplicate(_ context.Context, srcName, dstName string) (*entity.GameInfo, error)
	UpdateInfo(_ context.Context, name, newName string) error
	ImageCreate(ctx context.Context, gameID string, data []byte) error
	ImageGet(ctx context.Context, gameID string) ([]byte, error)
	ImageDelete(ctx context.Context, gameID string) error
}
