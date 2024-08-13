package collection

import (
	"github.com/HardDie/DeckBuilder/internal/dto"
	"github.com/HardDie/DeckBuilder/internal/entity"
)

type Collection interface {
	Create(gameID string, req *dto.CreateCollectionDTO) (*entity.CollectionInfo, error)
	GetByID(gameID, collectionID string) (*entity.CollectionInfo, error)
	GetAll(gameID string) ([]*entity.CollectionInfo, error)
	Update(gameID, collectionID string, dtoObject *dto.UpdateCollectionDTO) (*entity.CollectionInfo, error)
	DeleteByID(gameID, collectionID string) error
	GetImage(gameID, collectionID string) ([]byte, string, error)
}
