package service

import (
	"strings"

	"github.com/HardDie/DeckBuilder/internal/config"
	"github.com/HardDie/DeckBuilder/internal/dto"
	"github.com/HardDie/DeckBuilder/internal/entity"
	"github.com/HardDie/DeckBuilder/internal/network"
	repositoriesGame "github.com/HardDie/DeckBuilder/internal/repositories/game"
	"github.com/HardDie/DeckBuilder/internal/utils"
)

type IGameService interface {
	Create(dtoObject *dto.CreateGameDTO) (*entity.GameInfo, error)
	Item(gameID string) (*entity.GameInfo, error)
	List(sortField, search string) ([]*entity.GameInfo, *network.Meta, error)
	Update(gameID string, dtoObject *dto.UpdateGameDTO) (*entity.GameInfo, error)
	Delete(gameID string) error
	GetImage(gameID string) ([]byte, string, error)
	Duplicate(gameID string, dtoObject *dto.DuplicateGameDTO) (*entity.GameInfo, error)
	Export(gameID string) ([]byte, error)
	Import(data []byte, name string) (*entity.GameInfo, error)
}
type GameService struct {
	cfg            *config.Config
	repositoryGame repositoriesGame.Game
}

func NewGameService(cfg *config.Config, repositoryGame repositoriesGame.Game) *GameService {
	return &GameService{
		cfg:            cfg,
		repositoryGame: repositoryGame,
	}
}

func (s *GameService) Create(dtoObject *dto.CreateGameDTO) (*entity.GameInfo, error) {
	game, err := s.repositoryGame.Create(dtoObject)
	if err != nil {
		return nil, err
	}
	game.FillCachedImage(s.cfg)
	return game, nil
}
func (s *GameService) Item(gameID string) (*entity.GameInfo, error) {
	game, err := s.repositoryGame.GetByID(gameID)
	if err != nil {
		return nil, err
	}
	game.FillCachedImage(s.cfg)
	return game, nil
}
func (s *GameService) List(sortField, search string) ([]*entity.GameInfo, *network.Meta, error) {
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
func (s *GameService) Update(gameID string, dtoObject *dto.UpdateGameDTO) (*entity.GameInfo, error) {
	game, err := s.repositoryGame.Update(gameID, dtoObject)
	if err != nil {
		return nil, err
	}
	game.FillCachedImage(s.cfg)
	return game, nil
}
func (s *GameService) Delete(gameID string) error {
	return s.repositoryGame.DeleteByID(gameID)
}
func (s *GameService) GetImage(gameID string) ([]byte, string, error) {
	return s.repositoryGame.GetImage(gameID)
}
func (s *GameService) Duplicate(gameID string, dtoObject *dto.DuplicateGameDTO) (*entity.GameInfo, error) {
	game, err := s.repositoryGame.Duplicate(gameID, dtoObject)
	if err != nil {
		return nil, err
	}
	game.FillCachedImage(s.cfg)
	return game, nil
}
func (s *GameService) Export(gameID string) ([]byte, error) {
	return s.repositoryGame.Export(gameID)
}
func (s *GameService) Import(data []byte, name string) (*entity.GameInfo, error) {
	game, err := s.repositoryGame.Import(data, name)
	if err != nil {
		return nil, err
	}
	game.FillCachedImage(s.cfg)
	return game, nil
}
