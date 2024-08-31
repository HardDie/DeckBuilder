package card

import (
	"strings"

	"github.com/HardDie/DeckBuilder/internal/config"
	entitiesCard "github.com/HardDie/DeckBuilder/internal/entities/card"
	repositoriesCard "github.com/HardDie/DeckBuilder/internal/repositories/card"
	"github.com/HardDie/DeckBuilder/internal/utils"
)

type card struct {
	cfg            *config.Config
	repositoryCard repositoriesCard.Card
}

func New(cfg *config.Config, repositoryCard repositoriesCard.Card) Card {
	return &card{
		cfg:            cfg,
		repositoryCard: repositoryCard,
	}
}

func (s *card) Create(gameID, collectionID, deckID string, req CreateRequest) (*entitiesCard.Card, error) {
	if req.Count < 1 {
		req.Count = 1
	}
	return s.repositoryCard.Create(gameID, collectionID, deckID, repositoriesCard.CreateRequest{
		Name:        req.Name,
		Description: req.Description,
		Image:       req.Image,
		Variables:   req.Variables,
		Count:       req.Count,
		ImageFile:   req.ImageFile,
	})
}
func (s *card) Item(gameID, collectionID, deckID string, cardID int64) (*entitiesCard.Card, error) {
	return s.repositoryCard.GetByID(gameID, collectionID, deckID, cardID)
}
func (s *card) List(gameID, collectionID, deckID, sortField, search string) ([]*entitiesCard.Card, error) {
	items, err := s.repositoryCard.GetAll(gameID, collectionID, deckID)
	if err != nil {
		return make([]*entitiesCard.Card, 0), err
	}

	// Filter
	var filteredItems []*entitiesCard.Card
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
		filteredItems = make([]*entitiesCard.Card, 0)
	}

	return filteredItems, nil
}
func (s *card) Update(gameID, collectionID, deckID string, cardID int64, req UpdateRequest) (*entitiesCard.Card, error) {
	return s.repositoryCard.Update(gameID, collectionID, deckID, cardID, repositoriesCard.UpdateRequest{
		Name:        req.Name,
		Description: req.Description,
		Image:       req.Image,
		Variables:   req.Variables,
		Count:       req.Count,
		ImageFile:   req.ImageFile,
	})
}
func (s *card) Delete(gameID, collectionID, deckID string, cardID int64) error {
	return s.repositoryCard.DeleteByID(gameID, collectionID, deckID, cardID)
}
func (s *card) GetImage(gameID, collectionID, deckID string, cardID int64) ([]byte, string, error) {
	return s.repositoryCard.GetImage(gameID, collectionID, deckID, cardID)
}
