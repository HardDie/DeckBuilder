package game

import (
	"github.com/HardDie/DeckBuilder/internal/dto"
	"github.com/HardDie/DeckBuilder/internal/entity"
)

type Game interface {
	Create(req *dto.CreateGameDTO) (*entity.GameInfo, error)
	GetByID(gameID string) (*entity.GameInfo, error)
	GetAll() ([]*entity.GameInfo, error)
	Update(gameID string, dtoObject *dto.UpdateGameDTO) (*entity.GameInfo, error)
	DeleteByID(gameID string) error
	GetImage(gameID string) ([]byte, string, error)
	Duplicate(gameID string, dtoObject *dto.DuplicateGameDTO) (*entity.GameInfo, error)
	Export(gameID string) ([]byte, error)
	Import(data []byte, name string) (*entity.GameInfo, error)
}
