package collection

import (
	"github.com/HardDie/DeckBuilder/bindings/catalog"
	"github.com/HardDie/DeckBuilder/internal/config"
	"github.com/HardDie/DeckBuilder/internal/dto"
	"github.com/HardDie/DeckBuilder/internal/network"
	servicesCollection "github.com/HardDie/DeckBuilder/internal/services/collection"
)

type ListResult struct {
	Data []*dto.Collection `json:"data"`
	Meta *network.Meta     `json:"meta"`
}

type Result struct {
	Data dto.Collection `json:"data"`
}

type Collection struct {
	cfg config.Config
	svc servicesCollection.Collection
}

func New(cfg config.Config, svc servicesCollection.Collection) *Collection {
	return &Collection{cfg: cfg, svc: svc}
}

func (c *Collection) List(gameID, sort, search string) (*ListResult, error) {
	items, err := c.svc.List(gameID, sort, search)
	if err != nil {
		return nil, err
	}

	respItems := make([]*dto.Collection, 0, len(items))
	for _, item := range items {
		d := catalog.CollectionDTO(c.cfg, gameID, *item)
		respItems = append(respItems, &d)
	}

	return &ListResult{
		Data: respItems,
		Meta: &network.Meta{Total: len(respItems)},
	}, nil
}

func (c *Collection) Create(gameID string, req catalog.WriteRequest) (*Result, error) {
	item, err := c.svc.Create(gameID, servicesCollection.CreateRequest{
		Name:        req.Name,
		Description: req.Description,
		Image:       req.Image,
		ImageFile:   req.ImageBytes(),
	})
	if err != nil {
		return nil, err
	}
	return &Result{Data: catalog.CollectionDTO(c.cfg, gameID, *item)}, nil
}

func (c *Collection) Read(gameID, collectionID string) (*Result, error) {
	item, err := c.svc.Item(gameID, collectionID)
	if err != nil {
		return nil, err
	}
	return &Result{Data: catalog.CollectionDTO(c.cfg, gameID, *item)}, nil
}

func (c *Collection) Update(gameID, collectionID string, req catalog.WriteRequest) (*Result, error) {
	item, err := c.svc.Update(gameID, collectionID, servicesCollection.UpdateRequest{
		Name:        req.Name,
		Description: req.Description,
		Image:       req.Image,
		ImageFile:   req.ImageBytes(),
	})
	if err != nil {
		return nil, err
	}
	return &Result{Data: catalog.CollectionDTO(c.cfg, gameID, *item)}, nil
}

func (c *Collection) Delete(gameID, collectionID string) error {
	return c.svc.Delete(gameID, collectionID)
}
