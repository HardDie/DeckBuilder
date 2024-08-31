package deck

import (
	"strings"

	"github.com/HardDie/DeckBuilder/internal/config"
	entitiesDeck "github.com/HardDie/DeckBuilder/internal/entities/deck"
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

func (s *deck) Create(gameID, collectionID string, req CreateRequest) (*entitiesDeck.Deck, error) {
	return s.repositoryDeck.Create(gameID, collectionID, repositoriesDeck.CreateRequest{
		Name:        req.Name,
		Description: req.Description,
		Image:       req.Image,
		ImageFile:   req.ImageFile,
	})
}
func (s *deck) Item(gameID, collectionID, deckID string) (*entitiesDeck.Deck, error) {
	return s.repositoryDeck.GetByID(gameID, collectionID, deckID)
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
func (s *deck) Update(gameID, collectionID, deckID string, req UpdateRequest) (*entitiesDeck.Deck, error) {
	return s.repositoryDeck.Update(gameID, collectionID, deckID, repositoriesDeck.UpdateRequest{
		Name:        req.Name,
		Description: req.Description,
		Image:       req.Image,
		ImageFile:   req.ImageFile,
	})
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
