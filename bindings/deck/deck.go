package deck

import (
	"github.com/HardDie/DeckBuilder/bindings/catalog"
	"github.com/HardDie/DeckBuilder/internal/config"
	"github.com/HardDie/DeckBuilder/internal/dto"
	"github.com/HardDie/DeckBuilder/internal/network"
	servicesDeck "github.com/HardDie/DeckBuilder/internal/services/deck"
)

type ListResult struct {
	Data []*dto.Deck   `json:"data"`
	Meta *network.Meta `json:"meta,omitempty"`
}

type Result struct {
	Data dto.Deck `json:"data"`
}

type Deck struct {
	cfg config.Config
	svc servicesDeck.Deck
}

func New(cfg config.Config, svc servicesDeck.Deck) *Deck {
	return &Deck{cfg: cfg, svc: svc}
}

func (d *Deck) List(gameID, collectionID, sort, search string) (*ListResult, error) {
	items, err := d.svc.List(gameID, collectionID, sort, search)
	if err != nil {
		return nil, err
	}

	respItems := make([]*dto.Deck, 0, len(items))
	for _, item := range items {
		deck := catalog.DeckDTO(d.cfg, *item)
		respItems = append(respItems, &deck)
	}

	return &ListResult{
		Data: respItems,
		Meta: &network.Meta{Total: len(respItems)},
	}, nil
}

func (d *Deck) ListAllUnique(gameID string) (*ListResult, error) {
	items, err := d.svc.ListAllUnique(gameID)
	if err != nil {
		return nil, err
	}

	respItems := make([]*dto.Deck, 0, len(items))
	for _, item := range items {
		deck := catalog.DeckDTO(d.cfg, *item)
		respItems = append(respItems, &deck)
	}

	return &ListResult{Data: respItems}, nil
}

func (d *Deck) Create(gameID, collectionID string, req catalog.WriteRequest) (*Result, error) {
	item, err := d.svc.Create(gameID, collectionID, servicesDeck.CreateRequest{
		Name:        req.Name,
		Description: req.Description,
		Image:       req.Image,
		ImageFile:   req.ImageBytes(),
	})
	if err != nil {
		return nil, err
	}
	return &Result{Data: catalog.DeckDTO(d.cfg, *item)}, nil
}

func (d *Deck) Read(gameID, collectionID, deckID string) (*Result, error) {
	item, err := d.svc.Item(gameID, collectionID, deckID)
	if err != nil {
		return nil, err
	}
	return &Result{Data: catalog.DeckDTO(d.cfg, *item)}, nil
}

func (d *Deck) Update(gameID, collectionID, deckID string, req catalog.WriteRequest) (*Result, error) {
	item, err := d.svc.Update(gameID, collectionID, deckID, servicesDeck.UpdateRequest{
		Name:        req.Name,
		Description: req.Description,
		Image:       req.Image,
		ImageFile:   req.ImageBytes(),
	})
	if err != nil {
		return nil, err
	}
	return &Result{Data: catalog.DeckDTO(d.cfg, *item)}, nil
}

func (d *Deck) Delete(gameID, collectionID, deckID string) error {
	return d.svc.Delete(gameID, collectionID, deckID)
}
