package cards

import (
	"net/http"
	"time"

	"tts_deck_build/internal/config"
	"tts_deck_build/internal/decks"
	"tts_deck_build/internal/errors"
	"tts_deck_build/internal/fs"
	"tts_deck_build/internal/utils"
)

type CardStorage struct {
	Config      *config.Config
	DeckService *decks.DeckService
}

func NewCardStorage(config *config.Config, deckService *decks.DeckService) *CardStorage {
	return &CardStorage{
		Config:      config,
		DeckService: deckService,
	}
}

func (s *CardStorage) Create(gameID, collectionID, deckID string, card *CardInfo) (*CardInfo, error) {
	// Check if the deck exists
	deck, err := s.DeckService.Item(gameID, collectionID, deckID)
	if err != nil {
		return nil, err
	}

	// Read info from file
	readCard, err := fs.OpenAndProcess(deck.Path(gameID, collectionID), fs.JsonFromReader[Card])
	if err != nil {
		return nil, err
	}

	// Init map of cards
	if readCard.Cards == nil {
		readCard.Cards = make(map[int64]*CardInfo)
	}

	// Add card to deck
	readCard.Cards[card.ID] = card

	// Quote values before write to file
	defer card.SetRawOutput()
	for key := range readCard.Cards {
		readCard.Cards[key].SetQuotedOutput()
	}

	// Writing info to file
	if err := fs.CreateAndProcess(deck.Path(gameID, collectionID), *readCard, fs.JsonToWriter[Card]); err != nil {
		return nil, err
	}

	return card, nil
}
func (s *CardStorage) GetByID(gameID, collectionID, deckID string, cardID int64) (*CardInfo, error) {
	// Read map of cards
	cardsMap, err := s.getCardsMap(gameID, collectionID, deckID)
	if err != nil {
		return nil, err
	}

	// Check if card exist
	card, ok := cardsMap[cardID]
	if !ok {
		return nil, errors.CardNotExists
	}

	return card, nil
}
func (s *CardStorage) GetAll(gameID, collectionID, deckID string) ([]*CardInfo, error) {
	// Read map of cards
	cardsMap, err := s.getCardsMap(gameID, collectionID, deckID)
	if err != nil {
		return nil, err
	}

	// Convert map to list
	cards := make([]*CardInfo, 0)
	for _, card := range cardsMap {
		cards = append(cards, card)
	}
	return cards, nil
}
func (s *CardStorage) Update(gameID, collectionID, deckID string, cardID int64, dto *UpdateCardDTO) (*CardInfo, error) {
	// Check if the deck exists
	deck, err := s.DeckService.Item(gameID, collectionID, deckID)
	if err != nil {
		return nil, err
	}

	// Read info from file
	readCard, err := fs.OpenAndProcess(deck.Path(gameID, collectionID), fs.JsonFromReader[Card])
	if err != nil {
		return nil, err
	}

	// Check if card exist
	oldCard, ok := readCard.Cards[cardID]
	if !ok {
		return nil, errors.CardNotExists
	}

	// Create card object
	card := NewCardInfo(dto.Title, dto.Description, dto.Image, dto.Variables, dto.Count)
	card.ID = oldCard.ID
	card.CreatedAt = oldCard.CreatedAt

	// If the object has been changed, update the object file
	if !oldCard.Compare(card) {
		card.UpdatedAt = utils.Allocate(time.Now())
		// Replace old card with new one
		readCard.Cards[card.ID] = card

		// Quote values before write to file
		defer card.SetRawOutput()
		for key := range readCard.Cards {
			readCard.Cards[key].SetQuotedOutput()
		}

		// Writing info to file
		if err := fs.CreateAndProcess(deck.Path(gameID, collectionID), *readCard, fs.JsonToWriter[Card]); err != nil {
			return nil, err
		}
	}

	return card, nil
}
func (s *CardStorage) DeleteByID(gameID, collectionID, deckID string, cardID int64) error {
	// Check if the deck exists
	deck, err := s.DeckService.Item(gameID, collectionID, deckID)
	if err != nil {
		return err
	}

	// Read info from file
	readCard, err := fs.OpenAndProcess(deck.Path(gameID, collectionID), fs.JsonFromReader[Card])
	if err != nil {
		return err
	}

	// Check if card exist
	if _, ok := readCard.Cards[cardID]; !ok {
		return errors.CardNotExists.HTTP(http.StatusBadRequest)
	}

	// Delete card from deck
	delete(readCard.Cards, cardID)

	// Writing info to file
	if err := fs.CreateAndProcess(deck.Path(gameID, collectionID), *readCard, fs.JsonToWriter[Card]); err != nil {
		return err
	}
	return nil
}

// Internal function. Get map of cards inside deck
func (s *CardStorage) getCardsMap(gameID, collectionID, deckID string) (map[int64]*CardInfo, error) {
	// Check if the deck exists
	deck, err := s.DeckService.Item(gameID, collectionID, deckID)
	if err != nil {
		return nil, err
	}

	// Read info from file
	readCard, err := fs.OpenAndProcess(deck.Path(gameID, collectionID), fs.JsonFromReader[Card])
	if err != nil {
		return nil, err
	}
	return readCard.Cards, nil
}
