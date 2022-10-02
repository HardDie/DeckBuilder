package service

import (
	"tts_deck_build/internal/dto"
	"tts_deck_build/internal/entity"
	"tts_deck_build/internal/repository"
	"tts_deck_build/internal/utils"
)

type IDeckService interface {
	Create(gameID, collectionID string, dtoObject *dto.CreateDeckDTO) (*entity.DeckInfo, error)
	Item(gameID, collectionID, deckID string) (*entity.DeckInfo, error)
	List(gameID, collectionID, sortField string) ([]*entity.DeckInfo, error)
	Update(gameID, collectionID, deckID string, dtoObject *dto.UpdateDeckDTO) (*entity.DeckInfo, error)
	Delete(gameID, collectionID, deckID string) error
	GetImage(gameID, collectionID, deckID string) ([]byte, string, error)
	ListAllUnique(gameID string) ([]*entity.DeckInfo, error)
}
type DeckService struct {
	deckRepository repository.IDeckRepository
}

func NewDeckService(deckRepository repository.IDeckRepository) *DeckService {
	return &DeckService{
		deckRepository: deckRepository,
	}
}

func (s *DeckService) Create(gameID, collectionID string, dtoObject *dto.CreateDeckDTO) (*entity.DeckInfo, error) {
	return s.deckRepository.Create(gameID, collectionID, entity.NewDeckInfo(dtoObject.Name, dtoObject.Image))
}
func (s *DeckService) Item(gameID, collectionID, deckID string) (*entity.DeckInfo, error) {
	return s.deckRepository.GetByID(gameID, collectionID, deckID)
}
func (s *DeckService) List(gameID, collectionID, sortField string) ([]*entity.DeckInfo, error) {
	items, err := s.deckRepository.GetAll(gameID, collectionID)
	if err != nil {
		return make([]*entity.DeckInfo, 0), err
	}
	utils.Sort(&items, sortField)
	return items, nil
}
func (s *DeckService) Update(gameID, collectionID, deckID string, dtoObject *dto.UpdateDeckDTO) (*entity.DeckInfo, error) {
	return s.deckRepository.Update(gameID, collectionID, deckID, dtoObject)
}
func (s *DeckService) Delete(gameID, collectionID, deckID string) error {
	return s.deckRepository.DeleteByID(gameID, collectionID, deckID)
}
func (s *DeckService) GetImage(gameID, collectionID, deckID string) ([]byte, string, error) {
	return s.deckRepository.GetImage(gameID, collectionID, deckID)
}
func (s *DeckService) ListAllUnique(gameID string) ([]*entity.DeckInfo, error) {
	items, err := s.deckRepository.GetAllDecksInGame(gameID)
	if err != nil {
		return make([]*entity.DeckInfo, 0), err
	}
	utils.Sort(&items, "name")
	return items, nil
}
