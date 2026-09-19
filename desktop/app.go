package main

import (
	"context"
	"fmt"
	"net"

	"github.com/HardDie/DeckBuilder/internal/config"
	"github.com/HardDie/DeckBuilder/internal/dto"
	entitiesCollection "github.com/HardDie/DeckBuilder/internal/entities/collection"
	entitiesGame "github.com/HardDie/DeckBuilder/internal/entities/game"
	"github.com/HardDie/DeckBuilder/internal/network"
	servicesCollection "github.com/HardDie/DeckBuilder/internal/services/collection"
	servicesGame "github.com/HardDie/DeckBuilder/internal/services/game"
	"github.com/HardDie/DeckBuilder/internal/utils"
)

type ListGamesResult struct {
	Data []*dto.Game   `json:"data"`
	Meta *network.Meta `json:"meta"`
}

type CreateGameRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Image       string `json:"image"`
	ImageFile   []byte `json:"imageFile"`
}

type GameResult struct {
	Data dto.Game `json:"data"`
}

type ListCollectionsResult struct {
	Data []*dto.Collection `json:"data"`
	Meta *network.Meta     `json:"meta"`
}

type CollectionResult struct {
	Data dto.Collection `json:"data"`
}

// App is the Wails bindings surface. Catalog verbs move here incrementally; HTTP remains.
type App struct {
	ctx        context.Context
	ln         net.Listener
	cfg        config.Config
	game       servicesGame.Game
	collection servicesCollection.Collection
}

func NewApp(
	ln net.Listener,
	cfg config.Config,
	game servicesGame.Game,
	collection servicesCollection.Collection,
) *App {
	return &App{ln: ln, cfg: cfg, game: game, collection: collection}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *App) shutdown(_ context.Context) {
	if a.ln != nil {
		_ = a.ln.Close()
	}
}

func (a *App) ListGames(sort, search string) (*ListGamesResult, error) {
	items, err := a.game.List(sort, search)
	if err != nil {
		return nil, err
	}

	respItems := make([]*dto.Game, 0, len(items))
	for _, item := range items {
		g := a.gameDTO(*item)
		respItems = append(respItems, &g)
	}

	return &ListGamesResult{
		Data: respItems,
		Meta: &network.Meta{Total: len(respItems)},
	}, nil
}

func (a *App) CreateGame(req CreateGameRequest) (*GameResult, error) {
	item, err := a.game.Create(servicesGame.CreateRequest{
		Name:        req.Name,
		Description: req.Description,
		Image:       req.Image,
		ImageFile:   req.ImageFile,
	})
	if err != nil {
		return nil, err
	}

	return &GameResult{Data: a.gameDTO(*item)}, nil
}

func (a *App) ReadGame(gameID string) (*GameResult, error) {
	item, err := a.game.Item(gameID)
	if err != nil {
		return nil, err
	}

	return &GameResult{Data: a.gameDTO(*item)}, nil
}

func (a *App) UpdateGame(gameID string, req CreateGameRequest) (*GameResult, error) {
	item, err := a.game.Update(gameID, servicesGame.UpdateRequest{
		Name:        req.Name,
		Description: req.Description,
		Image:       req.Image,
		ImageFile:   req.ImageFile,
	})
	if err != nil {
		return nil, err
	}

	return &GameResult{Data: a.gameDTO(*item)}, nil
}

func (a *App) DeleteGame(gameID string) error {
	return a.game.Delete(gameID)
}

func (a *App) DuplicateGame(gameID, name string) (*GameResult, error) {
	item, err := a.game.Duplicate(gameID, servicesGame.DuplicateRequest{Name: name})
	if err != nil {
		return nil, err
	}

	return &GameResult{Data: a.gameDTO(*item)}, nil
}

func (a *App) ListCollections(gameID, sort, search string) (*ListCollectionsResult, error) {
	items, err := a.collection.List(gameID, sort, search)
	if err != nil {
		return nil, err
	}

	respItems := make([]*dto.Collection, 0, len(items))
	for _, item := range items {
		c := a.collectionDTO(gameID, *item)
		respItems = append(respItems, &c)
	}

	return &ListCollectionsResult{
		Data: respItems,
		Meta: &network.Meta{Total: len(respItems)},
	}, nil
}

func (a *App) CreateCollection(gameID string, req CreateGameRequest) (*CollectionResult, error) {
	item, err := a.collection.Create(gameID, servicesCollection.CreateRequest{
		Name:        req.Name,
		Description: req.Description,
		Image:       req.Image,
		ImageFile:   req.ImageFile,
	})
	if err != nil {
		return nil, err
	}

	return &CollectionResult{Data: a.collectionDTO(gameID, *item)}, nil
}

func (a *App) ReadCollection(gameID, collectionID string) (*CollectionResult, error) {
	item, err := a.collection.Item(gameID, collectionID)
	if err != nil {
		return nil, err
	}

	return &CollectionResult{Data: a.collectionDTO(gameID, *item)}, nil
}

func (a *App) UpdateCollection(gameID, collectionID string, req CreateGameRequest) (*CollectionResult, error) {
	item, err := a.collection.Update(gameID, collectionID, servicesCollection.UpdateRequest{
		Name:        req.Name,
		Description: req.Description,
		Image:       req.Image,
		ImageFile:   req.ImageFile,
	})
	if err != nil {
		return nil, err
	}

	return &CollectionResult{Data: a.collectionDTO(gameID, *item)}, nil
}

func (a *App) DeleteCollection(gameID, collectionID string) error {
	return a.collection.Delete(gameID, collectionID)
}

func (a *App) gameDTO(item entitiesGame.Game) dto.Game {
	return dto.Game{
		ID:          item.ID,
		Name:        item.Name,
		Description: item.Description,
		Image:       item.Image,
		CachedImage: a.cachedGameImage(item),
		CreatedAt:   item.CreatedAt,
		UpdatedAt:   item.UpdatedAt,
	}
}

func (a *App) cachedGameImage(game entitiesGame.Game) string {
	return fmt.Sprintf(a.cfg.GameImagePath+"?%s", game.ID, utils.HashForTime(&game.UpdatedAt))
}

func (a *App) collectionDTO(gameID string, item entitiesCollection.Collection) dto.Collection {
	return dto.Collection{
		ID:          item.ID,
		Name:        item.Name,
		Description: item.Description,
		Image:       item.Image,
		CachedImage: a.cachedCollectionImage(gameID, item),
		CreatedAt:   item.CreatedAt,
		UpdatedAt:   item.UpdatedAt,
	}
}

func (a *App) cachedCollectionImage(gameID string, collection entitiesCollection.Collection) string {
	return fmt.Sprintf(a.cfg.CollectionImagePath+"?%s", gameID, collection.ID, utils.HashForTime(&collection.UpdatedAt))
}
