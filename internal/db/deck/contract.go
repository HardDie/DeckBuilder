package deck

import (
	"context"

	"github.com/HardDie/DeckBuilder/internal/entity"
)

type Deck interface {
	Create(ctx context.Context, gameID, collectionID, name, description, image string) (*entity.DeckInfo, error)
	Get(ctx context.Context, gameID, collectionID, name string) (context.Context, *entity.DeckInfo, error)
	List(ctx context.Context, gameID, collectionID string) ([]*entity.DeckInfo, error)
	Move(ctx context.Context, gameID, collectionID, oldName, newName string) (*entity.DeckInfo, error)
	Update(ctx context.Context, gameID, collectionID, name, description, image string) (*entity.DeckInfo, error)
	Delete(ctx context.Context, gameID, collectionID, name string) error
	ImageCreate(ctx context.Context, gameID, collectionID, deckID string, data []byte) error
	ImageGet(ctx context.Context, gameID, collectionID, deckID string) ([]byte, error)
	ImageDelete(ctx context.Context, gameID, collectionID, deckID string) error
}
