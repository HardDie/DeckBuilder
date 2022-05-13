package decks

import (
	"tts_deck_build/internal/collections"
	"tts_deck_build/internal/config"
	"tts_deck_build/internal/utils"
)

type DeckService struct {
	storage *DeckStorage
}

func NewService() *DeckService {
	return &DeckService{
		storage: NewDeckStorage(config.GetConfig(), collections.NewService()),
	}
}

func (s *DeckService) Create(gameID, collectionID string, dto *CreateDeckDTO) (*DeckInfo, error) {
	return s.storage.Create(gameID, collectionID, NewDeckInfo(dto.Type, dto.BacksideImage))
}

func (s *DeckService) Item(gameID, collectionID, deckID string) (*DeckInfo, error) {
	return s.storage.GetByID(gameID, collectionID, deckID)
}

func (s *DeckService) List(gameID, collectionID, sortField string) ([]*DeckInfo, error) {
	items, err := s.storage.GetAll(gameID, collectionID)
	if err != nil {
		return make([]*DeckInfo, 0), err
	}
	utils.Sort(&items, sortField)
	return items, nil
}

func (s *DeckService) Update(gameID, collectionID, deckID string, dto *UpdateDeckDTO) (*DeckInfo, error) {
	return s.storage.Update(gameID, collectionID, deckID, dto)
}

func (s *DeckService) Delete(gameID, collectionID, deckID string) error {
	return s.storage.DeleteByID(gameID, collectionID, deckID)
}

func (s *DeckService) GetImage(gameID, collectionID, deckID string) ([]byte, string, error) {
	return s.storage.GetImage(gameID, collectionID, deckID)
}

func (s *DeckService) ListAllUnique(gameID string) ([]*DeckInfo, error) {
	items, err := s.storage.GetAllDecksInGame(gameID)
	if err != nil {
		return make([]*DeckInfo, 0), err
	}
	utils.Sort(&items, "name")
	return items, nil
}
