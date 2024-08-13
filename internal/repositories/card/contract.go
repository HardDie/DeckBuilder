package card

import (
	"github.com/HardDie/DeckBuilder/internal/dto"
	"github.com/HardDie/DeckBuilder/internal/entity"
)

type Card interface {
	Create(gameID, collectionID, deckID string, req *dto.CreateCardDTO) (*entity.CardInfo, error)
	GetByID(gameID, collectionID, deckID string, cardID int64) (*entity.CardInfo, error)
	GetAll(gameID, collectionID, deckID string) ([]*entity.CardInfo, error)
	Update(gameID, collectionID, deckID string, cardID int64, dtoObject *dto.UpdateCardDTO) (*entity.CardInfo, error)
	DeleteByID(gameID, collectionID, deckID string, cardID int64) error
	GetImage(gameID, collectionID, deckID string, cardID int64) ([]byte, string, error)
}
