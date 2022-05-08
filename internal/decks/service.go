package decks

import (
	"tts_deck_build/internal/collections"
	"tts_deck_build/internal/config"
)

type DeckService struct {
	storage *DeckStorage
}

func NewService() *DeckService {
	return &DeckService{
		storage: NewDeckStorage(config.GetConfig(), collections.NewService()),
	}
}

func (s *DeckService) Create(gameId, collectionId string, dto *CreateDeckDTO) (*DeckInfo, error) {
	return s.storage.Create(gameId, collectionId, NewDeckInfo("", dto.Type, dto.BacksideImage))
}

func (s *DeckService) Item(gameId, collectionId, deckId string) (*DeckInfo, error) {
	return s.storage.GetById(gameId, collectionId, deckId)
}

func (s *DeckService) List(gameId, collectionId string) ([]*DeckInfo, error) {
	return s.storage.GetAll(gameId, collectionId)
}

func (s *DeckService) Update(gameId, collectionId, deckId string, dto *UpdateDeckDTO) (*DeckInfo, error) {
	return s.storage.Update(gameId, collectionId, deckId, dto)
}

func (s *DeckService) Delete(gameId, collectionId, deckId string) error {
	return s.storage.DeleteById(gameId, collectionId, deckId)
}

func (s *DeckService) GetImage(gameId, collectionId, deckId string) ([]byte, string, error) {
	return s.storage.GetImage(gameId, collectionId, deckId)
}

func (s *DeckService) ListAllUnique(gameId string) ([]*DeckInfo, error) {
	return s.storage.GetAllDecksInGame(gameId)
}
