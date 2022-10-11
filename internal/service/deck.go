package service

import (
	"github.com/HardDie/DeckBuilder/internal/config"
	"github.com/HardDie/DeckBuilder/internal/dto"
	"github.com/HardDie/DeckBuilder/internal/entity"
	"github.com/HardDie/DeckBuilder/internal/repository"
	"github.com/HardDie/DeckBuilder/internal/utils"
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
	cfg            *config.Config
	deckRepository repository.IDeckRepository
}

func NewDeckService(cfg *config.Config, deckRepository repository.IDeckRepository) *DeckService {
	return &DeckService{
		cfg:            cfg,
		deckRepository: deckRepository,
	}
}

func (s *DeckService) Create(gameID, collectionID string, dtoObject *dto.CreateDeckDTO) (*entity.DeckInfo, error) {
	deck, err := s.deckRepository.Create(gameID, collectionID, dtoObject)
	if err != nil {
		return nil, err
	}
	deck.FillCachedImage(s.cfg, gameID, collectionID)
	return deck, nil
}
func (s *DeckService) Item(gameID, collectionID, deckID string) (*entity.DeckInfo, error) {
	deck, err := s.deckRepository.GetByID(gameID, collectionID, deckID)
	if err != nil {
		return nil, err
	}
	deck.FillCachedImage(s.cfg, gameID, collectionID)
	return deck, nil
}
func (s *DeckService) List(gameID, collectionID, sortField string) ([]*entity.DeckInfo, error) {
	items, err := s.deckRepository.GetAll(gameID, collectionID)
	if err != nil {
		return make([]*entity.DeckInfo, 0), err
	}
	utils.Sort(&items, sortField)
	for i := 0; i < len(items); i++ {
		items[i].FillCachedImage(s.cfg, gameID, collectionID)
	}
	return items, nil
}
func (s *DeckService) Update(gameID, collectionID, deckID string, dtoObject *dto.UpdateDeckDTO) (*entity.DeckInfo, error) {
	deck, err := s.deckRepository.Update(gameID, collectionID, deckID, dtoObject)
	if err != nil {
		return nil, err
	}
	deck.FillCachedImage(s.cfg, gameID, collectionID)
	return deck, nil
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
