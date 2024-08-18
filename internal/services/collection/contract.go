package collection

import (
	"github.com/HardDie/DeckBuilder/internal/entity"
	"github.com/HardDie/DeckBuilder/internal/network"
)

type Collection interface {
	Create(gameID string, req CreateRequest) (*entity.CollectionInfo, error)
	Item(gameID, collectionID string) (*entity.CollectionInfo, error)
	List(gameID, sortField, search string) ([]*entity.CollectionInfo, *network.Meta, error)
	Update(gameID, collectionID string, req UpdateRequest) (*entity.CollectionInfo, error)
	Delete(gameID, collectionID string) error
	GetImage(gameID, collectionID string) ([]byte, string, error)
}

type CreateRequest struct {
	Name        string
	Description string
	Image       string
	ImageFile   []byte
}

type UpdateRequest struct {
	Name        string
	Description string
	Image       string
	ImageFile   []byte
}
