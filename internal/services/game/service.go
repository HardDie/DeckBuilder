package game

import (
	"os"
	"path/filepath"

	"github.com/HardDie/DeckBuilder/internal/config"
	entitiesGame "github.com/HardDie/DeckBuilder/internal/entities/game"
	"github.com/HardDie/DeckBuilder/internal/fs"
	"github.com/HardDie/DeckBuilder/internal/logger"
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

	filteredItems := utils.FilterByName(items, search)

	// Sorting
	utils.Sort(&filteredItems, sortField)

	// Return empty array if no elements
	if filteredItems == nil {
		filteredItems = make([]*entitiesGame.Game, 0)
	}

	return filteredItems, nil
}
func (s *game) Update(gameID string, req UpdateRequest) (*entitiesGame.Game, error) {
	item, err := s.repositoryGame.Update(gameID, repositoriesGame.UpdateRequest{
		Name:        req.Name,
		Description: req.Description,
		Image:       req.Image,
		ImageFile:   req.ImageFile,
	})
	if err != nil {
		return nil, err
	}
	if item.ID != gameID {
		s.moveResults(gameID, item.ID)
	}
	return item, nil
}
func (s *game) Delete(gameID string) error {
	if err := s.repositoryGame.DeleteByID(gameID); err != nil {
		return err
	}
	// The game's last render goes with it.
	if err := fs.RemoveFolder(s.results(gameID)); err != nil {
		logger.Warn.Println("Unable to remove the game's result folder:", err.Error())
	}
	return nil
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

// results is the folder of the game's last render.
func (s *game) results(gameID string) string {
	return filepath.Join(s.cfg.Results(), gameID)
}

// moveResults follows a rename, so the next render can reuse its pages.
// It is best effort: the rename itself already succeeded.
func (s *game) moveResults(oldID, newID string) {
	src := s.results(oldID)
	if _, err := os.Stat(src); err != nil {
		return
	}
	dst := s.results(newID)
	if err := fs.RemoveFolder(dst); err != nil {
		logger.Warn.Println("Unable to clear the renamed game's result folder:", err.Error())
		return
	}
	if err := os.Rename(src, dst); err != nil {
		logger.Warn.Println("Unable to move the game's result folder:", err.Error())
	}
}
