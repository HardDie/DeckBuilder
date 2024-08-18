package game

import (
	"strings"

	"github.com/HardDie/DeckBuilder/internal/config"
	"github.com/HardDie/DeckBuilder/internal/entity"
	"github.com/HardDie/DeckBuilder/internal/network"
	repositoriesGame "github.com/HardDie/DeckBuilder/internal/repositories/game"
	"github.com/HardDie/DeckBuilder/internal/utils"
)

type game struct {
	cfg            *config.Config
	repositoryGame repositoriesGame.Game
}

func New(cfg *config.Config, repositoryGame repositoriesGame.Game) Game {
	return &game{
		cfg:            cfg,
		repositoryGame: repositoryGame,
	}
}

func (s *game) Create(req CreateRequest) (*entity.GameInfo, error) {
	g, err := s.repositoryGame.Create(repositoriesGame.CreateRequest{
		Name:        req.Name,
		Description: req.Description,
		Image:       req.Image,
		ImageFile:   req.ImageFile,
	})
	if err != nil {
		return nil, err
	}
	g.FillCachedImage(s.cfg)
	return g, nil
}
func (s *game) Item(gameID string) (*entity.GameInfo, error) {
	g, err := s.repositoryGame.GetByID(gameID)
	if err != nil {
		return nil, err
	}
	g.FillCachedImage(s.cfg)
	return g, nil
}
func (s *game) List(sortField, search string) ([]*entity.GameInfo, *network.Meta, error) {
	items, err := s.repositoryGame.GetAll()
	if err != nil {
		return make([]*entity.GameInfo, 0), nil, err
	}

	// Filter
	var filteredItems []*entity.GameInfo
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
		filteredItems[i].FillCachedImage(s.cfg)
	}

	// Return empty array if no elements
	if filteredItems == nil {
		filteredItems = make([]*entity.GameInfo, 0)
	}

	meta := &network.Meta{
		Total: len(filteredItems),
	}
	return filteredItems, meta, nil
}
func (s *game) Update(gameID string, req UpdateRequest) (*entity.GameInfo, error) {
	g, err := s.repositoryGame.Update(gameID, repositoriesGame.UpdateRequest{
		Name:        req.Name,
		Description: req.Description,
		Image:       req.Image,
		ImageFile:   req.ImageFile,
	})
	if err != nil {
		return nil, err
	}
	g.FillCachedImage(s.cfg)
	return g, nil
}
func (s *game) Delete(gameID string) error {
	return s.repositoryGame.DeleteByID(gameID)
}
func (s *game) GetImage(gameID string) ([]byte, string, error) {
	return s.repositoryGame.GetImage(gameID)
}
func (s *game) Duplicate(gameID string, req DuplicateRequest) (*entity.GameInfo, error) {
	g, err := s.repositoryGame.Duplicate(gameID, repositoriesGame.DuplicateRequest{
		Name: req.Name,
	})
	if err != nil {
		return nil, err
	}
	g.FillCachedImage(s.cfg)
	return g, nil
}
func (s *game) Export(gameID string) ([]byte, error) {
	return s.repositoryGame.Export(gameID)
}
func (s *game) Import(data []byte, name string) (*entity.GameInfo, error) {
	g, err := s.repositoryGame.Import(data, name)
	if err != nil {
		return nil, err
	}
	g.FillCachedImage(s.cfg)
	return g, nil
}
