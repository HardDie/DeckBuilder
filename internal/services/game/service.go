package game

import (
	"strings"

	"github.com/HardDie/DeckBuilder/internal/config"
	entitiesGame "github.com/HardDie/DeckBuilder/internal/entities/game"
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

func (s *game) Create(req CreateRequest) (*entitiesGame.Game, error) {
	return s.repositoryGame.Create(repositoriesGame.CreateRequest{
		Name:        req.Name,
		Description: req.Description,
		Image:       req.Image,
		ImageFile:   req.ImageFile,
	})
}
func (s *game) Item(gameID string) (*entitiesGame.Game, error) {
	return s.repositoryGame.GetByID(gameID)
}
func (s *game) List(sortField, search string) ([]*entitiesGame.Game, error) {
	items, err := s.repositoryGame.GetAll()
	if err != nil {
		return make([]*entitiesGame.Game, 0), err
	}

	// Filter
	var filteredItems []*entitiesGame.Game
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
		filteredItems = make([]*entitiesGame.Game, 0)
	}

	return filteredItems, nil
}
func (s *game) Update(gameID string, req UpdateRequest) (*entitiesGame.Game, error) {
	return s.repositoryGame.Update(gameID, repositoriesGame.UpdateRequest{
		Name:        req.Name,
		Description: req.Description,
		Image:       req.Image,
		ImageFile:   req.ImageFile,
	})
}
func (s *game) Delete(gameID string) error {
	return s.repositoryGame.DeleteByID(gameID)
}
func (s *game) GetImage(gameID string) ([]byte, string, error) {
	return s.repositoryGame.GetImage(gameID)
}
func (s *game) Duplicate(gameID string, req DuplicateRequest) (*entitiesGame.Game, error) {
	return s.repositoryGame.Duplicate(gameID, repositoriesGame.DuplicateRequest{
		Name: req.Name,
	})
}
func (s *game) Export(gameID string) ([]byte, error) {
	return s.repositoryGame.Export(gameID)
}
func (s *game) Import(data []byte, name string) (*entitiesGame.Game, error) {
	return s.repositoryGame.Import(data, name)
}
