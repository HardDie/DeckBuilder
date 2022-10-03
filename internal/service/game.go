package service

import (
	"tts_deck_build/internal/dto"
	"tts_deck_build/internal/entity"
	"tts_deck_build/internal/repository"
	"tts_deck_build/internal/utils"
)

type IGameService interface {
	Create(dtoObject *dto.CreateGameDTO) (*entity.GameInfo, error)
	Item(gameID string) (*entity.GameInfo, error)
	List(sortField string) ([]*entity.GameInfo, error)
	Update(gameID string, dtoObject *dto.UpdateGameDTO) (*entity.GameInfo, error)
	Delete(gameID string) error
	GetImage(gameID string) ([]byte, string, error)
	Duplicate(gameID string, dtoObject *dto.DuplicateGameDTO) (*entity.GameInfo, error)
	Export(gameID string) ([]byte, error)
	Import(data []byte, name string) (*entity.GameInfo, error)
}
type GameService struct {
	gameRepository repository.IGameRepository
}

func NewGameService(gameRepository repository.IGameRepository) *GameService {
	return &GameService{
		gameRepository: gameRepository,
	}
}

func (s *GameService) Create(dtoObject *dto.CreateGameDTO) (*entity.GameInfo, error) {
	return s.gameRepository.Create(entity.NewGameInfo(dtoObject.Name, dtoObject.Description, dtoObject.Image))
}
func (s *GameService) Item(gameID string) (*entity.GameInfo, error) {
	return s.gameRepository.GetByID(gameID)
}
func (s *GameService) List(sortField string) ([]*entity.GameInfo, error) {
	items, err := s.gameRepository.GetAll()
	if err != nil {
		return make([]*entity.GameInfo, 0), err
	}
	utils.Sort(&items, sortField)
	return items, nil
}
func (s *GameService) Update(gameID string, dtoObject *dto.UpdateGameDTO) (*entity.GameInfo, error) {
	return s.gameRepository.Update(gameID, dtoObject)
}
func (s *GameService) Delete(gameID string) error {
	return s.gameRepository.DeleteByID(gameID)
}
func (s *GameService) GetImage(gameID string) ([]byte, string, error) {
	return s.gameRepository.GetImage(gameID)
}
func (s *GameService) Duplicate(gameID string, dtoObject *dto.DuplicateGameDTO) (*entity.GameInfo, error) {
	return s.gameRepository.Duplicate(gameID, dtoObject)
}
func (s *GameService) Export(gameID string) ([]byte, error) {
	return s.gameRepository.Export(gameID)
}
func (s *GameService) Import(data []byte, name string) (*entity.GameInfo, error) {
	return s.gameRepository.Import(data, name)
}
