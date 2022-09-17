package cards

import (
	"tts_deck_build/internal/config"
	"tts_deck_build/internal/dto"
	"tts_deck_build/internal/entity"
	"tts_deck_build/internal/repository"
	"tts_deck_build/internal/utils"
)

type CardService struct {
	rep repository.ICardRepository
}

func NewService(cfg *config.Config) *CardService {
	return &CardService{
		rep: repository.NewCardRepository(cfg,
			repository.NewDeckRepository(cfg,
				repository.NewCollectionRepository(cfg, repository.NewGameRepository(cfg)),
			),
		),
	}
}

func (s *CardService) Create(gameID, collectionID, deckID string, dtoObject *dto.CreateCardDTO) (*entity.CardInfo, error) {
	card := entity.NewCardInfo(dtoObject.Title, dtoObject.Description, dtoObject.Image, dtoObject.Variables, dtoObject.Count)

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

	return s.rep.Create(gameID, collectionID, deckID, card)
}
func (s *CardService) Item(gameID, collectionID, deckID string, cardID int64) (*entity.CardInfo, error) {
	return s.rep.GetByID(gameID, collectionID, deckID, cardID)
}
func (s *CardService) List(gameID, collectionID, deckID, sortField string) ([]*entity.CardInfo, error) {
	items, err := s.rep.GetAll(gameID, collectionID, deckID)
	if err != nil {
		return make([]*entity.CardInfo, 0), err
	}
	utils.Sort(&items, sortField)
	return items, nil
}
func (s *CardService) Update(gameID, collectionID, deckID string, cardID int64, dtoObject *dto.UpdateCardDTO) (*entity.CardInfo, error) {
	return s.rep.Update(gameID, collectionID, deckID, cardID, dtoObject)
}
func (s *CardService) Delete(gameID, collectionID, deckID string, cardID int64) error {
	return s.rep.DeleteByID(gameID, collectionID, deckID, cardID)
}
func (s *CardService) GetImage(gameID, collectionID, deckID string, cardID int64) ([]byte, string, error) {
	return s.rep.GetImage(gameID, collectionID, deckID, cardID)
}
