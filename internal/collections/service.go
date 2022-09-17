package collections

import (
	"tts_deck_build/internal/config"
	"tts_deck_build/internal/dto"
	"tts_deck_build/internal/entity"
	"tts_deck_build/internal/repository"
	"tts_deck_build/internal/utils"
)

type CollectionService struct {
	rep repository.ICollectionRepository
}

func NewService() *CollectionService {
	cfg := config.GetConfig()
	return &CollectionService{
		rep: repository.NewCollectionRepository(cfg, repository.NewGameRepository(cfg)),
	}
}

func (s *CollectionService) Create(gameID string, dtoObject *dto.CreateCollectionDTO) (*entity.CollectionInfo, error) {
	return s.rep.Create(gameID, entity.NewCollectionInfo(dtoObject.Name, dtoObject.Description, dtoObject.Image))
}
func (s *CollectionService) Item(gameID, collectionID string) (*entity.CollectionInfo, error) {
	return s.rep.GetByID(gameID, collectionID)
}
func (s *CollectionService) List(gameID, sortField string) ([]*entity.CollectionInfo, error) {
	items, err := s.rep.GetAll(gameID)
	if err != nil {
		return make([]*entity.CollectionInfo, 0), err
	}
	utils.Sort(&items, sortField)
	return items, nil
}
func (s *CollectionService) Update(gameID, collectionID string, dtoObject *dto.UpdateCollectionDTO) (*entity.CollectionInfo, error) {
	return s.rep.Update(gameID, collectionID, dtoObject)
}
func (s *CollectionService) Delete(gameID, collectionID string) error {
	return s.rep.DeleteByID(gameID, collectionID)
}
func (s *CollectionService) GetImage(gameID, collectionID string) ([]byte, string, error) {
	return s.rep.GetImage(gameID, collectionID)
}
