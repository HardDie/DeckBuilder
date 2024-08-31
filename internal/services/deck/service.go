package deck

import (
	"strings"

	"github.com/HardDie/DeckBuilder/internal/config"
	entitiesDeck "github.com/HardDie/DeckBuilder/internal/entities/deck"
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
func (s *deck) List(gameID, collectionID, sortField, search string) ([]*entitiesDeck.Deck, error) {
	items, err := s.repositoryDeck.GetAll(gameID, collectionID)
	if err != nil {
		return make([]*entitiesDeck.Deck, 0), err
	}

	// Filter
	var filteredItems []*entitiesDeck.Deck
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

	// Return empty array if no elements
	if filteredItems == nil {
		filteredItems = make([]*entitiesDeck.Deck, 0)
	}

	return filteredItems, nil
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
func (s *deck) ListAllUnique(gameID string) ([]*entitiesDeck.Deck, error) {
	items, err := s.repositoryDeck.GetAllDecksInGame(gameID)
	if err != nil {
		return make([]*entitiesDeck.Deck, 0), err
	}
	utils.Sort(&items, "name")
	return items, nil
}
