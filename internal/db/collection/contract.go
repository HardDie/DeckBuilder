package collection

import (
	"context"

	"github.com/HardDie/DeckBuilder/internal/entity"
)

type Collection interface {
	Create(ctx context.Context, gameID, name, description, image string) (*entity.CollectionInfo, error)
	Get(ctx context.Context, gameID, name string) (context.Context, *entity.CollectionInfo, error)
	List(ctx context.Context, gameID string) ([]*entity.CollectionInfo, error)
	Move(ctx context.Context, gameID, oldName, newName string) (*entity.CollectionInfo, error)
	Update(ctx context.Context, gameID, name, description, image string) (*entity.CollectionInfo, error)
	Delete(ctx context.Context, gameID, name string) error
	ImageCreate(ctx context.Context, gameID, collectionID string, data []byte) error
	ImageGet(ctx context.Context, gameID, collectionID string) ([]byte, error)
	ImageDelete(ctx context.Context, gameID, collectionID string) error
}
