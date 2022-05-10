package games

import (
	"tts_deck_build/internal/config"
)

type GameService struct {
	storage *GameStorage
}

func NewService() *GameService {
	return &GameService{
		storage: NewGameStorage(config.GetConfig()),
	}
}

func (s *GameService) Create(dto *CreateGameDTO) (*GameInfo, error) {
	return s.storage.Create(NewGameInfo(dto.Name, dto.Description, dto.Image))
}

func (s *GameService) Item(gameId string) (*GameInfo, error) {
	return s.storage.GetById(gameId)
}

func (s *GameService) List(sortField string) ([]*GameInfo, error) {
	items, err := s.storage.GetAll()
	if err != nil {
		return make([]*GameInfo, 0), err
	}
	Sort(&items, sortField)
	return items, nil
}

func (s *GameService) Update(gameId string, dto *UpdateGameDTO) (*GameInfo, error) {
	return s.storage.Update(gameId, dto)
}

func (s *GameService) Delete(gameId string) error {
	return s.storage.DeleteById(gameId)
}

func (s *GameService) GetImage(gameId string) ([]byte, string, error) {
	return s.storage.GetImage(gameId)
}
