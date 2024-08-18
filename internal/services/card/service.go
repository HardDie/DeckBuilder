package card

import (
	"strings"

	"github.com/HardDie/DeckBuilder/internal/config"
	"github.com/HardDie/DeckBuilder/internal/entity"
	"github.com/HardDie/DeckBuilder/internal/network"
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

func (s *card) Create(gameID, collectionID, deckID string, req CreateRequest) (*entity.CardInfo, error) {
	if req.Count < 1 {
		req.Count = 1
	}
	c, err := s.repositoryCard.Create(gameID, collectionID, deckID, repositoriesCard.CreateRequest{
		Name:        req.Name,
		Description: req.Description,
		Image:       req.Image,
		Variables:   req.Variables,
		Count:       req.Count,
		ImageFile:   req.ImageFile,
	})
	if err != nil {
		return nil, err
	}
	c.FillCachedImage(s.cfg, gameID, collectionID, deckID)
	return c, nil
}
func (s *card) Item(gameID, collectionID, deckID string, cardID int64) (*entity.CardInfo, error) {
	c, err := s.repositoryCard.GetByID(gameID, collectionID, deckID, cardID)
	if err != nil {
		return nil, err
	}
	c.FillCachedImage(s.cfg, gameID, collectionID, deckID)
	return c, nil
}
func (s *card) List(gameID, collectionID, deckID, sortField, search string) ([]*entity.CardInfo, *network.Meta, error) {
	items, err := s.repositoryCard.GetAll(gameID, collectionID, deckID)
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

	// Sorting
	utils.Sort(&filteredItems, sortField)

	// Generate field cachedImage
	var cardsTotal int
	for i := 0; i < len(filteredItems); i++ {
		cardsTotal += filteredItems[i].Count
		filteredItems[i].FillCachedImage(s.cfg, gameID, collectionID, deckID)
	}

	// Return empty array if no elements
	if filteredItems == nil {
		filteredItems = make([]*entity.CardInfo, 0)
	}

	meta := &network.Meta{
		Total:      len(filteredItems),
		CardsTotal: cardsTotal,
	}
	return filteredItems, meta, nil
}
func (s *card) Update(gameID, collectionID, deckID string, cardID int64, req UpdateRequest) (*entity.CardInfo, error) {
	c, err := s.repositoryCard.Update(gameID, collectionID, deckID, cardID, repositoriesCard.UpdateRequest{
		Name:        req.Name,
		Description: req.Description,
		Image:       req.Image,
		Variables:   req.Variables,
		Count:       req.Count,
		ImageFile:   req.ImageFile,
	})
	if err != nil {
		return nil, err
	}
	c.FillCachedImage(s.cfg, gameID, collectionID, deckID)
	return c, nil
}
func (s *card) Delete(gameID, collectionID, deckID string, cardID int64) error {
	return s.repositoryCard.DeleteByID(gameID, collectionID, deckID, cardID)
}
func (s *card) GetImage(gameID, collectionID, deckID string, cardID int64) ([]byte, string, error) {
	return s.repositoryCard.GetImage(gameID, collectionID, deckID, cardID)
}
