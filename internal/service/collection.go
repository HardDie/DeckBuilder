package service

import (
	"github.com/HardDie/DeckBuilder/internal/config"
	"github.com/HardDie/DeckBuilder/internal/dto"
	"github.com/HardDie/DeckBuilder/internal/entity"
	"github.com/HardDie/DeckBuilder/internal/repository"
	"github.com/HardDie/DeckBuilder/internal/utils"
)

type ICollectionService interface {
	Create(gameID string, dtoObject *dto.CreateCollectionDTO) (*entity.CollectionInfo, error)
	Item(gameID, collectionID string) (*entity.CollectionInfo, error)
	List(gameID, sortField string) ([]*entity.CollectionInfo, error)
	Update(gameID, collectionID string, dtoObject *dto.UpdateCollectionDTO) (*entity.CollectionInfo, error)
	Delete(gameID, collectionID string) error
	GetImage(gameID, collectionID string) ([]byte, string, error)
}
type CollectionService struct {
	cfg                  *config.Config
	collectionRepository repository.ICollectionRepository
}

func NewCollectionService(cfg *config.Config, collectionRepository repository.ICollectionRepository) *CollectionService {
	return &CollectionService{
		cfg:                  cfg,
		collectionRepository: collectionRepository,
	}
}

func (s *CollectionService) Create(gameID string, dtoObject *dto.CreateCollectionDTO) (*entity.CollectionInfo, error) {
	collection, err := s.collectionRepository.Create(gameID, dtoObject)
	if err != nil {
		return nil, err
	}
	collection.FillCachedImage(s.cfg, gameID)
	return collection, nil
}
func (s *CollectionService) Item(gameID, collectionID string) (*entity.CollectionInfo, error) {
	collection, err := s.collectionRepository.GetByID(gameID, collectionID)
	if err != nil {
		return nil, err
	}
	collection.FillCachedImage(s.cfg, gameID)
	return collection, nil
}
func (s *CollectionService) List(gameID, sortField string) ([]*entity.CollectionInfo, error) {
	items, err := s.collectionRepository.GetAll(gameID)
	if err != nil {
		return make([]*entity.CollectionInfo, 0), err
	}
	utils.Sort(&items, sortField)
	for i := 0; i < len(items); i++ {
		items[i].FillCachedImage(s.cfg, gameID)
	}
	if items == nil {
		items = make([]*entity.CollectionInfo, 0)
	}
	return items, nil
}
func (s *CollectionService) Update(gameID, collectionID string, dtoObject *dto.UpdateCollectionDTO) (*entity.CollectionInfo, error) {
	collection, err := s.collectionRepository.Update(gameID, collectionID, dtoObject)
	if err != nil {
		return nil, err
	}
	collection.FillCachedImage(s.cfg, gameID)
	return collection, nil
}
func (s *CollectionService) Delete(gameID, collectionID string) error {
	return s.collectionRepository.DeleteByID(gameID, collectionID)
}
func (s *CollectionService) GetImage(gameID, collectionID string) ([]byte, string, error) {
	return s.collectionRepository.GetImage(gameID, collectionID)
}
