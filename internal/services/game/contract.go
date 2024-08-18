package game

import (
	"github.com/HardDie/DeckBuilder/internal/entity"
	"github.com/HardDie/DeckBuilder/internal/network"
)

type Game interface {
	Create(req CreateRequest) (*entity.GameInfo, error)
	Item(gameID string) (*entity.GameInfo, error)
	List(sortField, search string) ([]*entity.GameInfo, *network.Meta, error)
	Update(gameID string, req UpdateRequest) (*entity.GameInfo, error)
	Delete(gameID string) error
	GetImage(gameID string) ([]byte, string, error)
	Duplicate(gameID string, req DuplicateRequest) (*entity.GameInfo, error)
	Export(gameID string) ([]byte, error)
	Import(data []byte, name string) (*entity.GameInfo, error)
}

type CreateRequest struct {
	Name        string
	Description string
	Image       string
	ImageFile   []byte
}

type UpdateRequest struct {
	Name        string
	Description string
	Image       string
	ImageFile   []byte
}

type DuplicateRequest struct {
	Name string
}
