package card

import (
	"github.com/HardDie/DeckBuilder/bindings/catalog"
	"github.com/HardDie/DeckBuilder/internal/config"
	"github.com/HardDie/DeckBuilder/internal/dto"
	"github.com/HardDie/DeckBuilder/internal/network"
	servicesCard "github.com/HardDie/DeckBuilder/internal/services/card"
)

type ListResult struct {
	Data []*dto.Card   `json:"data"`
	Meta *network.Meta `json:"meta"`
}

type Result struct {
	Data dto.Card `json:"data"`
	// Warning explains why a new image was not applied. The save itself succeeded.
	Warning string `json:"warning,omitempty"`
}

type Card struct {
	cfg config.Config
	svc servicesCard.Card
}

func New(cfg config.Config, svc servicesCard.Card) *Card {
	return &Card{cfg: cfg, svc: svc}
}

func (c *Card) List(gameID, collectionID, deckID, sort, search string) (*ListResult, error) {
	items, err := c.svc.List(gameID, collectionID, deckID, sort, search)
	if err != nil {
		return nil, err
	}

	respItems := make([]*dto.Card, 0, len(items))
	var cardsTotal int
	for _, item := range items {
		cardsTotal += item.Count
		d := catalog.CardDTO(c.cfg, *item)
		respItems = append(respItems, &d)
	}

	return &ListResult{
		Data: respItems,
		Meta: &network.Meta{Total: len(respItems), CardsTotal: cardsTotal},
	}, nil
}

func (c *Card) Create(gameID, collectionID, deckID string, req catalog.WriteRequest) (*Result, error) {
	item, err := c.svc.Create(gameID, collectionID, deckID, servicesCard.CreateRequest{
		Name:        req.Name,
		Description: req.Description,
		Image:       req.Image,
		Variables:   req.Variables,
		Count:       req.Count,
		ImageFile:   req.ImageBytes(),
	})
	if err != nil {
		return nil, err
	}
	return &Result{Data: catalog.CardDTO(c.cfg, *item), Warning: catalog.ImageWarning(item.ImageError)}, nil
}

func (c *Card) Read(gameID, collectionID, deckID string, cardID int64) (*Result, error) {
	item, err := c.svc.Item(gameID, collectionID, deckID, cardID)
	if err != nil {
		return nil, err
	}
	return &Result{Data: catalog.CardDTO(c.cfg, *item)}, nil
}

func (c *Card) Update(gameID, collectionID, deckID string, cardID int64, req catalog.WriteRequest) (*Result, error) {
	item, err := c.svc.Update(gameID, collectionID, deckID, cardID, servicesCard.UpdateRequest{
		Name:        req.Name,
		Description: req.Description,
		Image:       req.Image,
		Variables:   req.Variables,
		Count:       req.Count,
		ImageFile:   req.ImageBytes(),
	})
	if err != nil {
		return nil, err
	}
	return &Result{Data: catalog.CardDTO(c.cfg, *item), Warning: catalog.ImageWarning(item.ImageError)}, nil
}

func (c *Card) Delete(gameID, collectionID, deckID string, cardID int64) error {
	return c.svc.Delete(gameID, collectionID, deckID, cardID)
}
