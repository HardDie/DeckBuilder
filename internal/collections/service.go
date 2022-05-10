package collections

import (
	"tts_deck_build/internal/config"
	"tts_deck_build/internal/games"
)

type CollectionService struct {
	storage *CollectionStorage
}

func NewService() *CollectionService {
	return &CollectionService{
		storage: NewCollectionStorage(config.GetConfig(), games.NewService()),
	}
}

func (s *CollectionService) Create(gameId string, dto *CreateCollectionDTO) (*CollectionInfo, error) {
	return s.storage.Create(gameId, NewCollectionInfo(dto.Name, dto.Description, dto.Image))
}

func (s *CollectionService) Item(gameId, collectionId string) (*CollectionInfo, error) {
	return s.storage.GetById(gameId, collectionId)
}

func (s *CollectionService) List(gameId, sortField string) ([]*CollectionInfo, error) {
	items, err := s.storage.GetAll(gameId)
	if err != nil {
		return make([]*CollectionInfo, 0), err
	}
	Sort(&items, sortField)
	return items, nil
}

func (s *CollectionService) Update(gameId, collectionId string, dto *UpdateCollectionDTO) (*CollectionInfo, error) {
	return s.storage.Update(gameId, collectionId, dto)
}

func (s *CollectionService) Delete(gameId, collectionId string) error {
	return s.storage.DeleteById(gameId, collectionId)
}

func (s *CollectionService) GetImage(gameId, collectionId string) ([]byte, string, error) {
	return s.storage.GetImage(gameId, collectionId)
}
