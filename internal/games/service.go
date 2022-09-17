package games

import (
	"tts_deck_build/internal/config"
	"tts_deck_build/internal/dto"
	"tts_deck_build/internal/entity"
	"tts_deck_build/internal/repository"
	"tts_deck_build/internal/utils"
)

type GameService struct {
	rep repository.IGameRepository
}

func NewService() *GameService {
	return &GameService{
		rep: repository.NewGameRepository(config.GetConfig()),
	}
}

func (s *GameService) Create(dtoObject *dto.CreateGameDTO) (*entity.GameInfo, error) {
	return s.rep.Create(entity.NewGameInfo(dtoObject.Name, dtoObject.Description, dtoObject.Image))
}
func (s *GameService) Item(gameID string) (*entity.GameInfo, error) {
	return s.rep.GetByID(gameID)
}
func (s *GameService) List(sortField string) ([]*entity.GameInfo, error) {
	items, err := s.rep.GetAll()
	if err != nil {
		return make([]*entity.GameInfo, 0), err
	}
	utils.Sort(&items, sortField)
	return items, nil
}
func (s *GameService) Update(gameID string, dtoObject *dto.UpdateGameDTO) (*entity.GameInfo, error) {
	return s.rep.Update(gameID, dtoObject)
}
func (s *GameService) Delete(gameID string) error {
	return s.rep.DeleteByID(gameID)
}
func (s *GameService) GetImage(gameID string) ([]byte, string, error) {
	return s.rep.GetImage(gameID)
}
func (s *GameService) Duplicate(gameID string, dtoObject *dto.DuplicateGameDTO) (*entity.GameInfo, error) {
	return s.rep.Duplicate(gameID, dtoObject)
}
func (s *GameService) Export(gameID string) ([]byte, error) {
	return s.rep.Export(gameID)
}
func (s *GameService) Import(data []byte, name string) error {
	return s.rep.Import(data, name)
}
