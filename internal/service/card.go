package service

import (
	"github.com/HardDie/DeckBuilder/internal/dto"
	"github.com/HardDie/DeckBuilder/internal/entity"
	"github.com/HardDie/DeckBuilder/internal/repository"
	"github.com/HardDie/DeckBuilder/internal/utils"
)

type ICardService interface {
	Create(gameID, collectionID, deckID string, dtoObject *dto.CreateCardDTO) (*entity.CardInfo, error)
	Item(gameID, collectionID, deckID string, cardID int64) (*entity.CardInfo, error)
	List(gameID, collectionID, deckID, sortField string) ([]*entity.CardInfo, error)
	Update(gameID, collectionID, deckID string, cardID int64, dtoObject *dto.UpdateCardDTO) (*entity.CardInfo, error)
	Delete(gameID, collectionID, deckID string, cardID int64) error
	GetImage(gameID, collectionID, deckID string, cardID int64) ([]byte, string, error)
}
type CardService struct {
	cardRepository repository.ICardRepository
}

func NewCardService(cardRepository repository.ICardRepository) *CardService {
	return &CardService{
		cardRepository: cardRepository,
	}
}

func (s *CardService) Create(gameID, collectionID, deckID string, dtoObject *dto.CreateCardDTO) (*entity.CardInfo, error) {
	card := entity.NewCardInfo(dtoObject.Name, dtoObject.Description, dtoObject.Image, dtoObject.Variables, dtoObject.Count)

	// Get all cards in deck
	allCards, err := s.List(gameID, collectionID, deckID, "")
	if err != nil {
		return nil, err
	}

	// Find current biggest index
	var maxID int64
	for _, currentCard := range allCards {
		if currentCard.ID > maxID {
			maxID = currentCard.ID
		}
	}

	// Increase value
	maxID += 1
	// Set new max value to the new card
	card.ID = maxID

	return s.cardRepository.Create(gameID, collectionID, deckID, card)
}
func (s *CardService) Item(gameID, collectionID, deckID string, cardID int64) (*entity.CardInfo, error) {
	return s.cardRepository.GetByID(gameID, collectionID, deckID, cardID)
}
func (s *CardService) List(gameID, collectionID, deckID, sortField string) ([]*entity.CardInfo, error) {
	items, err := s.cardRepository.GetAll(gameID, collectionID, deckID)
	if err != nil {
		return make([]*entity.CardInfo, 0), err
	}
	utils.Sort(&items, sortField)
	return items, nil
}
func (s *CardService) Update(gameID, collectionID, deckID string, cardID int64, dtoObject *dto.UpdateCardDTO) (*entity.CardInfo, error) {
	return s.cardRepository.Update(gameID, collectionID, deckID, cardID, dtoObject)
}
func (s *CardService) Delete(gameID, collectionID, deckID string, cardID int64) error {
	return s.cardRepository.DeleteByID(gameID, collectionID, deckID, cardID)
}
func (s *CardService) GetImage(gameID, collectionID, deckID string, cardID int64) ([]byte, string, error) {
	return s.cardRepository.GetImage(gameID, collectionID, deckID, cardID)
}
