package card

import (
	"github.com/HardDie/DeckBuilder/internal/dto"
	"github.com/HardDie/DeckBuilder/internal/entity"
	"github.com/HardDie/DeckBuilder/internal/network"
)

type Card interface {
	Create(gameID, collectionID, deckID string, dtoObject *dto.CreateCardDTO) (*entity.CardInfo, error)
	Item(gameID, collectionID, deckID string, cardID int64) (*entity.CardInfo, error)
	List(gameID, collectionID, deckID, sortField, search string) ([]*entity.CardInfo, *network.Meta, error)
	Update(gameID, collectionID, deckID string, cardID int64, dtoObject *dto.UpdateCardDTO) (*entity.CardInfo, error)
	Delete(gameID, collectionID, deckID string, cardID int64) error
	GetImage(gameID, collectionID, deckID string, cardID int64) ([]byte, string, error)
}
