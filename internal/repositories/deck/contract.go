package deck

import (
	"github.com/HardDie/DeckBuilder/internal/dto"
	"github.com/HardDie/DeckBuilder/internal/entity"
)

type Deck interface {
	Create(gameID, collectionID string, req *dto.CreateDeckDTO) (*entity.DeckInfo, error)
	GetByID(gameID, collectionID, deckID string) (*entity.DeckInfo, error)
	GetAll(gameID, collectionID string) ([]*entity.DeckInfo, error)
	Update(gameID, collectionID, deckID string, dtoObject *dto.UpdateDeckDTO) (*entity.DeckInfo, error)
	DeleteByID(gameID, collectionID, deckID string) error
	GetImage(gameID, collectionID, deckID string) ([]byte, string, error)
	GetAllDecksInGame(gameID string) ([]*entity.DeckInfo, error)
}
