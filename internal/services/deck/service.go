package deck

import (
	"strings"

	"github.com/HardDie/DeckBuilder/internal/config"
	"github.com/HardDie/DeckBuilder/internal/dto"
	"github.com/HardDie/DeckBuilder/internal/entity"
	"github.com/HardDie/DeckBuilder/internal/network"
	repositoriesDeck "github.com/HardDie/DeckBuilder/internal/repositories/deck"
	"github.com/HardDie/DeckBuilder/internal/utils"
)

type deck struct {
	cfg            *config.Config
	repositoryDeck repositoriesDeck.Deck
}

func New(cfg *config.Config, repositoryDeck repositoriesDeck.Deck) Deck {
	return &deck{
		cfg:            cfg,
		repositoryDeck: repositoryDeck,
	}
}

func (s *deck) Create(gameID, collectionID string, dtoObject *dto.CreateDeckDTO) (*entity.DeckInfo, error) {
	deck, err := s.repositoryDeck.Create(gameID, collectionID, dtoObject)
	if err != nil {
		return nil, err
	}
	deck.FillCachedImage(s.cfg, gameID, collectionID)
	return deck, nil
}
func (s *deck) Item(gameID, collectionID, deckID string) (*entity.DeckInfo, error) {
	deck, err := s.repositoryDeck.GetByID(gameID, collectionID, deckID)
	if err != nil {
		return nil, err
	}
	deck.FillCachedImage(s.cfg, gameID, collectionID)
	return deck, nil
}
func (s *deck) List(gameID, collectionID, sortField, search string) ([]*entity.DeckInfo, *network.Meta, error) {
	items, err := s.repositoryDeck.GetAll(gameID, collectionID)
	if err != nil {
		return make([]*entity.DeckInfo, 0), nil, err
	}

	// Filter
	var filteredItems []*entity.DeckInfo
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

	// Sorting
	utils.Sort(&filteredItems, sortField)

	// Generate field cachedImage
	for i := 0; i < len(filteredItems); i++ {
		filteredItems[i].FillCachedImage(s.cfg, gameID, collectionID)
	}

	// Return empty array if no elements
	if filteredItems == nil {
		filteredItems = make([]*entity.DeckInfo, 0)
	}

	meta := &network.Meta{
		Total: len(filteredItems),
	}
	return filteredItems, meta, nil
}
func (s *deck) Update(gameID, collectionID, deckID string, dtoObject *dto.UpdateDeckDTO) (*entity.DeckInfo, error) {
	deck, err := s.repositoryDeck.Update(gameID, collectionID, deckID, dtoObject)
	if err != nil {
		return nil, err
	}
	deck.FillCachedImage(s.cfg, gameID, collectionID)
	return deck, nil
}
func (s *deck) Delete(gameID, collectionID, deckID string) error {
	return s.repositoryDeck.DeleteByID(gameID, collectionID, deckID)
}
func (s *deck) GetImage(gameID, collectionID, deckID string) ([]byte, string, error) {
	return s.repositoryDeck.GetImage(gameID, collectionID, deckID)
}
func (s *deck) ListAllUnique(gameID string) ([]*entity.DeckInfo, error) {
	items, err := s.repositoryDeck.GetAllDecksInGame(gameID)
	if err != nil {
		return make([]*entity.DeckInfo, 0), err
	}
	utils.Sort(&items, "name")
	return items, nil
}
