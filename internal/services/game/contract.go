package game

import (
	"github.com/HardDie/DeckBuilder/internal/dto"
	"github.com/HardDie/DeckBuilder/internal/entity"
	"github.com/HardDie/DeckBuilder/internal/network"
)

type Game interface {
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
