package game

import (
	"github.com/HardDie/DeckBuilder/desktop/bindings/catalog"
	"github.com/HardDie/DeckBuilder/internal/config"
	"github.com/HardDie/DeckBuilder/internal/dto"
	"github.com/HardDie/DeckBuilder/internal/network"
	servicesGame "github.com/HardDie/DeckBuilder/internal/services/game"
)

type ListResult struct {
	Data []*dto.Game   `json:"data"`
	Meta *network.Meta `json:"meta"`
}

type Result struct {
	Data dto.Game `json:"data"`
}

type Game struct {
	cfg config.Config
	svc servicesGame.Game
}

func New(cfg config.Config, svc servicesGame.Game) *Game {
	return &Game{cfg: cfg, svc: svc}
}

func (g *Game) List(sort, search string) (*ListResult, error) {
	items, err := g.svc.List(sort, search)
	if err != nil {
		return nil, err
	}

	respItems := make([]*dto.Game, 0, len(items))
	for _, item := range items {
		d := catalog.GameDTO(g.cfg, *item)
		respItems = append(respItems, &d)
	}

	return &ListResult{
		Data: respItems,
		Meta: &network.Meta{Total: len(respItems)},
	}, nil
}

func (g *Game) Create(req catalog.WriteRequest) (*Result, error) {
	item, err := g.svc.Create(servicesGame.CreateRequest{
		Name:        req.Name,
		Description: req.Description,
		Image:       req.Image,
		ImageFile:   req.ImageBytes(),
	})
	if err != nil {
		return nil, err
	}
	return &Result{Data: catalog.GameDTO(g.cfg, *item)}, nil
}

func (g *Game) Read(gameID string) (*Result, error) {
	item, err := g.svc.Item(gameID)
	if err != nil {
		return nil, err
	}
	return &Result{Data: catalog.GameDTO(g.cfg, *item)}, nil
}

func (g *Game) Update(gameID string, req catalog.WriteRequest) (*Result, error) {
	item, err := g.svc.Update(gameID, servicesGame.UpdateRequest{
		Name:        req.Name,
		Description: req.Description,
		Image:       req.Image,
		ImageFile:   req.ImageBytes(),
	})
	if err != nil {
		return nil, err
	}
	return &Result{Data: catalog.GameDTO(g.cfg, *item)}, nil
}

func (g *Game) Delete(gameID string) error {
	return g.svc.Delete(gameID)
}

func (g *Game) Duplicate(gameID, name string) (*Result, error) {
	item, err := g.svc.Duplicate(gameID, servicesGame.DuplicateRequest{Name: name})
	if err != nil {
		return nil, err
	}
	return &Result{Data: catalog.GameDTO(g.cfg, *item)}, nil
}
