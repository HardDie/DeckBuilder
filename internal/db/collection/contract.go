package collection

import (
	"context"

	entitiesCollection "github.com/HardDie/DeckBuilder/internal/entities/collection"
)

type Collection interface {
	Create(ctx context.Context, gameID, name, description, image string) (*entitiesCollection.Collection, error)
	Get(ctx context.Context, gameID, name string) (*entitiesCollection.Collection, error)
	List(ctx context.Context, gameID string) ([]*entitiesCollection.Collection, error)
	Move(ctx context.Context, gameID, oldName, newName string) (*entitiesCollection.Collection, error)
	Update(ctx context.Context, gameID, name, description, image string) (*entitiesCollection.Collection, error)
	Delete(ctx context.Context, gameID, name string) error
	ImageCreate(ctx context.Context, gameID, collectionID string, data []byte) error
	ImageGet(ctx context.Context, gameID, collectionID string) ([]byte, error)
	ImageDelete(ctx context.Context, gameID, collectionID string) error
}
