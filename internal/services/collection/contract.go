package collection

import (
	"github.com/HardDie/DeckBuilder/internal/dto"
	"github.com/HardDie/DeckBuilder/internal/entity"
	"github.com/HardDie/DeckBuilder/internal/network"
)

type Collection interface {
	Create(gameID string, dtoObject *dto.CreateCollectionDTO) (*entity.CollectionInfo, error)
	Item(gameID, collectionID string) (*entity.CollectionInfo, error)
	List(gameID, sortField, search string) ([]*entity.CollectionInfo, *network.Meta, error)
	Update(gameID, collectionID string, dtoObject *dto.UpdateCollectionDTO) (*entity.CollectionInfo, error)
	Delete(gameID, collectionID string) error
	GetImage(gameID, collectionID string) ([]byte, string, error)
}
