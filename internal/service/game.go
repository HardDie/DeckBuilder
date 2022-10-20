package service

import (
	"strings"

	"github.com/HardDie/DeckBuilder/internal/config"
	"github.com/HardDie/DeckBuilder/internal/dto"
	"github.com/HardDie/DeckBuilder/internal/entity"
	"github.com/HardDie/DeckBuilder/internal/repository"
	"github.com/HardDie/DeckBuilder/internal/utils"
)

type IGameService interface {
	Create(dtoObject *dto.CreateGameDTO) (*entity.GameInfo, error)
	Item(gameID string) (*entity.GameInfo, error)
	List(sortField, search string) ([]*entity.GameInfo, error)
	Update(gameID string, dtoObject *dto.UpdateGameDTO) (*entity.GameInfo, error)
	Delete(gameID string) error
	GetImage(gameID string) ([]byte, string, error)
	Duplicate(gameID string, dtoObject *dto.DuplicateGameDTO) (*entity.GameInfo, error)
	Export(gameID string) ([]byte, error)
	Import(data []byte, name string) (*entity.GameInfo, error)
}
type GameService struct {
	cfg            *config.Config
	gameRepository repository.IGameRepository
}

func NewGameService(cfg *config.Config, gameRepository repository.IGameRepository) *GameService {
	return &GameService{
		cfg:            cfg,
		gameRepository: gameRepository,
	}
}

func (s *GameService) Create(dtoObject *dto.CreateGameDTO) (*entity.GameInfo, error) {
	game, err := s.gameRepository.Create(dtoObject)
	if err != nil {
		return nil, err
	}
	game.FillCachedImage(s.cfg)
	return game, nil
}
func (s *GameService) Item(gameID string) (*entity.GameInfo, error) {
	game, err := s.gameRepository.GetByID(gameID)
	if err != nil {
		return nil, err
	}
	game.FillCachedImage(s.cfg)
	return game, nil
}
func (s *GameService) List(sortField, search string) ([]*entity.GameInfo, error) {
	items, err := s.gameRepository.GetAll()
	if err != nil {
		return make([]*entity.GameInfo, 0), err
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

	//Sorting
	utils.Sort(&filteredItems, sortField)

	// Generate field cachedImage
	for i := 0; i < len(filteredItems); i++ {
		filteredItems[i].FillCachedImage(s.cfg)
	}

	// Return empty array if no elements
	if filteredItems == nil {
		filteredItems = make([]*entity.GameInfo, 0)
	}
	return filteredItems, nil
}
func (s *GameService) Update(gameID string, dtoObject *dto.UpdateGameDTO) (*entity.GameInfo, error) {
	game, err := s.gameRepository.Update(gameID, dtoObject)
	if err != nil {
		return nil, err
	}
	game.FillCachedImage(s.cfg)
	return game, nil
}
func (s *GameService) Delete(gameID string) error {
	return s.gameRepository.DeleteByID(gameID)
}
func (s *GameService) GetImage(gameID string) ([]byte, string, error) {
	return s.gameRepository.GetImage(gameID)
}
func (s *GameService) Duplicate(gameID string, dtoObject *dto.DuplicateGameDTO) (*entity.GameInfo, error) {
	game, err := s.gameRepository.Duplicate(gameID, dtoObject)
	if err != nil {
		return nil, err
	}
	game.FillCachedImage(s.cfg)
	return game, nil
}
func (s *GameService) Export(gameID string) ([]byte, error) {
	return s.gameRepository.Export(gameID)
}
func (s *GameService) Import(data []byte, name string) (*entity.GameInfo, error) {
	game, err := s.gameRepository.Import(data, name)
	if err != nil {
		return nil, err
	}
	game.FillCachedImage(s.cfg)
	return game, nil
}
