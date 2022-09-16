package collections

import (
	"tts_deck_build/internal/config"
	"tts_deck_build/internal/dto"
	"tts_deck_build/internal/entity"
	"tts_deck_build/internal/games"
	"tts_deck_build/internal/utils"
)

type CollectionService struct {
	storage *CollectionStorage
}

func NewService() *CollectionService {
	return &CollectionService{
		storage: NewCollectionStorage(config.GetConfig(), games.NewService()),
	}
}

func (s *CollectionService) Create(gameID string, dtoObject *dto.CreateCollectionDTO) (*entity.CollectionInfo, error) {
	return s.storage.Create(gameID, entity.NewCollectionInfo(dtoObject.Name, dtoObject.Description, dtoObject.Image))
}

func (s *CollectionService) Item(gameID, collectionID string) (*entity.CollectionInfo, error) {
	return s.storage.GetByID(gameID, collectionID)
}

func (s *CollectionService) List(gameID, sortField string) ([]*entity.CollectionInfo, error) {
	items, err := s.storage.GetAll(gameID)
	if err != nil {
		return make([]*entity.CollectionInfo, 0), err
	}
	utils.Sort(&items, sortField)
	return items, nil
}

func (s *CollectionService) Update(gameID, collectionID string, dtoObject *dto.UpdateCollectionDTO) (*entity.CollectionInfo, error) {
	return s.storage.Update(gameID, collectionID, dtoObject)
}

func (s *CollectionService) Delete(gameID, collectionID string) error {
	return s.storage.DeleteByID(gameID, collectionID)
}

func (s *CollectionService) GetImage(gameID, collectionID string) ([]byte, string, error) {
	return s.storage.GetImage(gameID, collectionID)
}
