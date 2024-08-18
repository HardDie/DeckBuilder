package collection

import (
	"github.com/HardDie/DeckBuilder/internal/entity"
)

type Collection interface {
	Create(gameID string, req CreateRequest) (*entity.CollectionInfo, error)
	GetByID(gameID, collectionID string) (*entity.CollectionInfo, error)
	GetAll(gameID string) ([]*entity.CollectionInfo, error)
	Update(gameID, collectionID string, req UpdateRequest) (*entity.CollectionInfo, error)
	DeleteByID(gameID, collectionID string) error
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
