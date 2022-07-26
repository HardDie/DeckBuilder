package games

import (
	"tts_deck_build/internal/config"
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

func (s *GameService) Create(dto *CreateGameDTO) (*GameInfo, error) {
	return s.storage.Create(NewGameInfo(dto.Name, dto.Description, dto.Image))
}

func (s *GameService) Item(gameID string) (*GameInfo, error) {
	return s.storage.GetByID(gameID)
}

func (s *GameService) List(sortField string) ([]*GameInfo, error) {
	items, err := s.storage.GetAll()
	if err != nil {
		return make([]*GameInfo, 0), err
	}
	utils.Sort(&items, sortField)
	return items, nil
}

func (s *GameService) Update(gameID string, dto *UpdateGameDTO) (*GameInfo, error) {
	return s.storage.Update(gameID, dto)
}

func (s *GameService) Delete(gameID string) error {
	return s.storage.DeleteByID(gameID)
}

func (s *GameService) GetImage(gameID string) ([]byte, string, error) {
	return s.storage.GetImage(gameID)
}

func (s *GameService) Export(gameID string) ([]byte, error) {
	return s.storage.Export(gameID)
}

func (s *GameService) Import(data []byte, name string) error {
	return s.storage.Import(data, name)
}
