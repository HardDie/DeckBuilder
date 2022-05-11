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

func (s *CollectionService) Create(gameID string, dto *CreateCollectionDTO) (*CollectionInfo, error) {
	return s.storage.Create(gameID, NewCollectionInfo(dto.Name, dto.Description, dto.Image))
}

func (s *CollectionService) Item(gameID, collectionID string) (*CollectionInfo, error) {
	return s.storage.GetByID(gameID, collectionID)
}

func (s *CollectionService) List(gameID, sortField string) ([]*CollectionInfo, error) {
	items, err := s.storage.GetAll(gameID)
	if err != nil {
		return make([]*CollectionInfo, 0), err
	}
	Sort(&items, sortField)
	return items, nil
}

func (s *CollectionService) Update(gameID, collectionID string, dto *UpdateCollectionDTO) (*CollectionInfo, error) {
	return s.storage.Update(gameID, collectionID, dto)
}

func (s *CollectionService) Delete(gameID, collectionID string) error {
	return s.storage.DeleteByID(gameID, collectionID)
}

func (s *CollectionService) GetImage(gameID, collectionID string) ([]byte, string, error) {
	return s.storage.GetImage(gameID, collectionID)
}
