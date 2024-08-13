package deck

import (
	"github.com/HardDie/DeckBuilder/internal/dto"
	"github.com/HardDie/DeckBuilder/internal/entity"
	"github.com/HardDie/DeckBuilder/internal/network"
)

type Deck interface {
	Create(gameID, collectionID string, dtoObject *dto.CreateDeckDTO) (*entity.DeckInfo, error)
	Item(gameID, collectionID, deckID string) (*entity.DeckInfo, error)
	List(gameID, collectionID, sortField, search string) ([]*entity.DeckInfo, *network.Meta, error)
	Update(gameID, collectionID, deckID string, dtoObject *dto.UpdateDeckDTO) (*entity.DeckInfo, error)
	Delete(gameID, collectionID, deckID string) error
	GetImage(gameID, collectionID, deckID string) ([]byte, string, error)
	ListAllUnique(gameID string) ([]*entity.DeckInfo, error)
}
