package service

import (
	"strings"

	"github.com/HardDie/DeckBuilder/internal/config"
	"github.com/HardDie/DeckBuilder/internal/dto"
	"github.com/HardDie/DeckBuilder/internal/entity"
	"github.com/HardDie/DeckBuilder/internal/network"
	"github.com/HardDie/DeckBuilder/internal/repository"
	"github.com/HardDie/DeckBuilder/internal/utils"
)

type ICardService interface {
	Create(gameID, collectionID, deckID string, dtoObject *dto.CreateCardDTO) (*entity.CardInfo, error)
	Item(gameID, collectionID, deckID string, cardID int64) (*entity.CardInfo, error)
	List(gameID, collectionID, deckID, sortField, search string) ([]*entity.CardInfo, *network.Meta, error)
	Update(gameID, collectionID, deckID string, cardID int64, dtoObject *dto.UpdateCardDTO) (*entity.CardInfo, error)
	Delete(gameID, collectionID, deckID string, cardID int64) error
	GetImage(gameID, collectionID, deckID string, cardID int64) ([]byte, string, error)
}
type CardService struct {
	cfg            *config.Config
	cardRepository repository.ICardRepository
}

func NewCardService(cfg *config.Config, cardRepository repository.ICardRepository) *CardService {
	return &CardService{
		cfg:            cfg,
		cardRepository: cardRepository,
	}
}

func (s *CardService) Create(gameID, collectionID, deckID string, dtoObject *dto.CreateCardDTO) (*entity.CardInfo, error) {
	if dtoObject.Count < 1 {
		dtoObject.Count = 1
	}
	card, err := s.cardRepository.Create(gameID, collectionID, deckID, dtoObject)
	if err != nil {
		return nil, err
	}
	card.FillCachedImage(s.cfg, gameID, collectionID, deckID)
	return card, nil
}
func (s *CardService) Item(gameID, collectionID, deckID string, cardID int64) (*entity.CardInfo, error) {
	card, err := s.cardRepository.GetByID(gameID, collectionID, deckID, cardID)
	if err != nil {
		return nil, err
	}
	card.FillCachedImage(s.cfg, gameID, collectionID, deckID)
	return card, nil
}
func (s *CardService) List(gameID, collectionID, deckID, sortField, search string) ([]*entity.CardInfo, *network.Meta, error) {
	items, err := s.cardRepository.GetAll(gameID, collectionID, deckID)
	if err != nil {
		return make([]*entity.CardInfo, 0), nil, err
	}

	// Filter
	var filteredItems []*entity.CardInfo
	if search != "" {
		search = strings.ToLower(search)
		for _, item := range items {
			if strings.Contains(strings.ToLower(item.Name), search) {
				filteredItems = append(filteredItems, item)
			}
		}
	} else {
		filteredItems = items
	}

	//Sorting
	utils.Sort(&filteredItems, sortField)

	// Generate field cachedImage
	for i := 0; i < len(filteredItems); i++ {
		filteredItems[i].FillCachedImage(s.cfg, gameID, collectionID, deckID)
	}

	// Return empty array if no elements
	if filteredItems == nil {
		filteredItems = make([]*entity.CardInfo, 0)
	}

	meta := &network.Meta{
		Total: len(filteredItems),
	}
	return filteredItems, meta, nil
}
func (s *CardService) Update(gameID, collectionID, deckID string, cardID int64, dtoObject *dto.UpdateCardDTO) (*entity.CardInfo, error) {
	card, err := s.cardRepository.Update(gameID, collectionID, deckID, cardID, dtoObject)
	if err != nil {
		return nil, err
	}
	card.FillCachedImage(s.cfg, gameID, collectionID, deckID)
	return card, nil
}
func (s *CardService) Delete(gameID, collectionID, deckID string, cardID int64) error {
	return s.cardRepository.DeleteByID(gameID, collectionID, deckID, cardID)
}
func (s *CardService) GetImage(gameID, collectionID, deckID string, cardID int64) ([]byte, string, error) {
	return s.cardRepository.GetImage(gameID, collectionID, deckID, cardID)
}
