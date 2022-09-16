package games

import (
	"tts_deck_build/internal/config"
	"tts_deck_build/internal/dto"
	"tts_deck_build/internal/entity"
	"tts_deck_build/internal/utils"
)

type GameService struct {
	storage *GameStorage
}

func NewService() *GameService {
	return &GameService{
		storage: NewGameStorage(config.GetConfig()),
	}
}

func (s *GameService) Create(dtoObject *dto.CreateGameDTO) (*entity.GameInfo, error) {
	return s.storage.Create(entity.NewGameInfo(dtoObject.Name, dtoObject.Description, dtoObject.Image))
}

func (s *GameService) Item(gameID string) (*entity.GameInfo, error) {
	return s.storage.GetByID(gameID)
}

func (s *GameService) List(sortField string) ([]*entity.GameInfo, error) {
	items, err := s.storage.GetAll()
	if err != nil {
		return make([]*entity.GameInfo, 0), err
	}
	utils.Sort(&items, sortField)
	return items, nil
}

func (s *GameService) Update(gameID string, dtoObject *dto.UpdateGameDTO) (*entity.GameInfo, error) {
	return s.storage.Update(gameID, dtoObject)
}

func (s *GameService) Delete(gameID string) error {
	return s.storage.DeleteByID(gameID)
}

func (s *GameService) GetImage(gameID string) ([]byte, string, error) {
	return s.storage.GetImage(gameID)
}

func (s *GameService) Duplicate(gameID string, dtoObject *dto.DuplicateGameDTO) (*entity.GameInfo, error) {
	return s.storage.Duplicate(gameID, dtoObject)
}

func (s *GameService) Export(gameID string) ([]byte, error) {
	return s.storage.Export(gameID)
}

func (s *GameService) Import(data []byte, name string) error {
	return s.storage.Import(data, name)
}
