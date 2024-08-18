package game

import (
	entitiesGame "github.com/HardDie/DeckBuilder/internal/entities/game"
)

type Game interface {
	Create(req CreateRequest) (*entitiesGame.Game, error)
	GetByID(gameID string) (*entitiesGame.Game, error)
	GetAll() ([]*entitiesGame.Game, error)
	Update(gameID string, req UpdateRequest) (*entitiesGame.Game, error)
	DeleteByID(gameID string) error
	GetImage(gameID string) ([]byte, string, error)
	Duplicate(gameID string, req DuplicateRequest) (*entitiesGame.Game, error)
	Export(gameID string) ([]byte, error)
	Import(data []byte, name string) (*entitiesGame.Game, error)
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
