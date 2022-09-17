package service

import (
	"tts_deck_build/internal/dto"
	"tts_deck_build/internal/entity"
	"tts_deck_build/internal/repository"
	"tts_deck_build/internal/utils"
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
	collectionRepository repository.ICollectionRepository
}

func NewCollectionService(collectionRepository repository.ICollectionRepository) *CollectionService {
	return &CollectionService{
		collectionRepository: collectionRepository,
	}
}

func (s *CollectionService) Create(gameID string, dtoObject *dto.CreateCollectionDTO) (*entity.CollectionInfo, error) {
	return s.collectionRepository.Create(gameID, entity.NewCollectionInfo(dtoObject.Name, dtoObject.Description, dtoObject.Image))
}
func (s *CollectionService) Item(gameID, collectionID string) (*entity.CollectionInfo, error) {
	return s.collectionRepository.GetByID(gameID, collectionID)
}
func (s *CollectionService) List(gameID, sortField string) ([]*entity.CollectionInfo, error) {
	items, err := s.collectionRepository.GetAll(gameID)
	if err != nil {
		return make([]*entity.CollectionInfo, 0), err
	}
	utils.Sort(&items, sortField)
	return items, nil
}
func (s *CollectionService) Update(gameID, collectionID string, dtoObject *dto.UpdateCollectionDTO) (*entity.CollectionInfo, error) {
	return s.collectionRepository.Update(gameID, collectionID, dtoObject)
}
func (s *CollectionService) Delete(gameID, collectionID string) error {
	return s.collectionRepository.DeleteByID(gameID, collectionID)
}
func (s *CollectionService) GetImage(gameID, collectionID string) ([]byte, string, error) {
	return s.collectionRepository.GetImage(gameID, collectionID)
}
