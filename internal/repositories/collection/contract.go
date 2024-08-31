package collection

import (
	entitiesCollection "github.com/HardDie/DeckBuilder/internal/entities/collection"
)

type Collection interface {
	Create(gameID string, req CreateRequest) (*entitiesCollection.Collection, error)
	GetByID(gameID, collectionID string) (*entitiesCollection.Collection, error)
	GetAll(gameID string) ([]*entitiesCollection.Collection, error)
	Update(gameID, collectionID string, req UpdateRequest) (*entitiesCollection.Collection, error)
	DeleteByID(gameID, collectionID string) error
	GetImage(gameID, collectionID string) ([]byte, string, error)
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
