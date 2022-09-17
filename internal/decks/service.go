package decks

import (
	"tts_deck_build/internal/config"
	"tts_deck_build/internal/dto"
	"tts_deck_build/internal/entity"
	"tts_deck_build/internal/repository"
	"tts_deck_build/internal/utils"
)

type DeckService struct {
	rep repository.IDeckRepository
}

func NewService(cfg *config.Config) *DeckService {
	return &DeckService{
		rep: repository.NewDeckRepository(
			cfg,
			repository.NewCollectionRepository(cfg, repository.NewGameRepository(cfg)),
		),
	}
}

func (s *DeckService) Create(gameID, collectionID string, dtoObject *dto.CreateDeckDTO) (*entity.DeckInfo, error) {
	return s.rep.Create(gameID, collectionID, entity.NewDeckInfo(dtoObject.Type, dtoObject.BacksideImage))
}
func (s *DeckService) Item(gameID, collectionID, deckID string) (*entity.DeckInfo, error) {
	return s.rep.GetByID(gameID, collectionID, deckID)
}
func (s *DeckService) List(gameID, collectionID, sortField string) ([]*entity.DeckInfo, error) {
	items, err := s.rep.GetAll(gameID, collectionID)
	if err != nil {
		return make([]*entity.DeckInfo, 0), err
	}
	utils.Sort(&items, sortField)
	return items, nil
}
func (s *DeckService) Update(gameID, collectionID, deckID string, dtoObject *dto.UpdateDeckDTO) (*entity.DeckInfo, error) {
	return s.rep.Update(gameID, collectionID, deckID, dtoObject)
}
func (s *DeckService) Delete(gameID, collectionID, deckID string) error {
	return s.rep.DeleteByID(gameID, collectionID, deckID)
}
func (s *DeckService) GetImage(gameID, collectionID, deckID string) ([]byte, string, error) {
	return s.rep.GetImage(gameID, collectionID, deckID)
}
func (s *DeckService) ListAllUnique(gameID string) ([]*entity.DeckInfo, error) {
	items, err := s.rep.GetAllDecksInGame(gameID)
	if err != nil {
		return make([]*entity.DeckInfo, 0), err
	}
	utils.Sort(&items, "name")
	return items, nil
}
