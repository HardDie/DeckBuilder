package repository

import (
	"fmt"
	"net/http"
	"time"

	"tts_deck_build/internal/config"
	"tts_deck_build/internal/dto"
	"tts_deck_build/internal/entity"
	"tts_deck_build/internal/errors"
	"tts_deck_build/internal/fs"
	"tts_deck_build/internal/images"
	"tts_deck_build/internal/network"
	"tts_deck_build/internal/utils"
)

type ICardRepository interface {
	Create(gameID, collectionID, deckID string, card *entity.CardInfo) (*entity.CardInfo, error)
	GetByID(gameID, collectionID, deckID string, cardID int64) (*entity.CardInfo, error)
	GetAll(gameID, collectionID, deckID string) ([]*entity.CardInfo, error)
	Update(gameID, collectionID, deckID string, cardID int64, dtoObject *dto.UpdateCardDTO) (*entity.CardInfo, error)
	DeleteByID(gameID, collectionID, deckID string, cardID int64) error
	GetImage(gameID, collectionID, deckID string, cardID int64) ([]byte, string, error)
	CreateImage(gameID, collectionID, deckID string, cardID int64, imageURL string) error
}
type CardRepository struct {
	cfg            *config.Config
	deckRepository IDeckRepository
}

func NewCardRepository(cfg *config.Config, deckRepository IDeckRepository) *CardRepository {
	return &CardRepository{
		cfg:            cfg,
		deckRepository: deckRepository,
	}
}

func (s *CardRepository) Create(gameID, collectionID, deckID string, card *entity.CardInfo) (*entity.CardInfo, error) {
	// Check if the deck exists
	deck, err := s.deckRepository.GetByID(gameID, collectionID, deckID)
	if err != nil {
		return nil, err
	}

	// Read info from file
	readCard, err := fs.OpenAndProcess(deck.Path(gameID, collectionID, s.cfg), fs.JsonFromReader[entity.Card])
	if err != nil {
		return nil, err
	}

	// Init map of cards
	if readCard.Cards == nil {
		readCard.Cards = make(map[int64]*entity.CardInfo)
	}

	if card.Image != "" {
		card.CachedImage = fmt.Sprintf(s.cfg.CardImagePath, gameID, collectionID, deckID, card.ID)
	} else {
		card.CachedImage = ""
	}

	// Add card to deck
	readCard.Cards[card.ID] = card

	// Quote values before write to file
	defer card.SetRawOutput()
	for key := range readCard.Cards {
		readCard.Cards[key].SetQuotedOutput()
	}

	// Writing info to file
	if err := fs.CreateAndProcess(deck.Path(gameID, collectionID, s.cfg), *readCard, fs.JsonToWriter[entity.Card]); err != nil {
		return nil, err
	}

	if len(card.Image) > 0 {
		// Download image
		if err := s.CreateImage(gameID, collectionID, deck.ID, card.ID, card.Image); err != nil {
			return nil, err
		}
	}

	return card, nil
}
func (s *CardRepository) GetByID(gameID, collectionID, deckID string, cardID int64) (*entity.CardInfo, error) {
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
func (s *CardRepository) GetAll(gameID, collectionID, deckID string) ([]*entity.CardInfo, error) {
	// Read map of cards
	cardsMap, err := s.getCardsMap(gameID, collectionID, deckID)
	if err != nil {
		return nil, err
	}

	// Convert map to list
	cards := make([]*entity.CardInfo, 0)
	for _, card := range cardsMap {
		cards = append(cards, card)
	}
	return cards, nil
}
func (s *CardRepository) Update(gameID, collectionID, deckID string, cardID int64, dtoObject *dto.UpdateCardDTO) (*entity.CardInfo, error) {
	// Check if the deck exists
	deck, err := s.deckRepository.GetByID(gameID, collectionID, deckID)
	if err != nil {
		return nil, err
	}

	// Read info from file
	readCard, err := fs.OpenAndProcess(deck.Path(gameID, collectionID, s.cfg), fs.JsonFromReader[entity.Card])
	if err != nil {
		return nil, err
	}

	// Check if card exist
	oldCard, ok := readCard.Cards[cardID]
	if !ok {
		return nil, errors.CardNotExists
	}

	// Create card object
	card := entity.NewCardInfo(dtoObject.Name, dtoObject.Description, dtoObject.Image, dtoObject.Variables, dtoObject.Count)
	card.ID = oldCard.ID
	card.CreatedAt = oldCard.CreatedAt

	if card.Image != "" {
		card.CachedImage = fmt.Sprintf(s.cfg.CardImagePath, gameID, collectionID, deckID, card.ID)
	} else {
		card.CachedImage = ""
	}

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
		if err := fs.CreateAndProcess(deck.Path(gameID, collectionID, s.cfg), *readCard, fs.JsonToWriter[entity.Card]); err != nil {
			return nil, err
		}
	}

	// If the image has been changed
	if card.Image != oldCard.Image {
		// If image exist, delete
		if data, _, _ := s.GetImage(gameID, collectionID, deckID, card.ID); data != nil {
			err = fs.RemoveFile(card.ImagePath(gameID, collectionID, deckID, s.cfg))
			if err != nil {
				return nil, err
			}
		}

		if len(card.Image) > 0 {
			// Download image
			if err = s.CreateImage(gameID, collectionID, deckID, card.ID, card.Image); err != nil {
				return nil, err
			}
		}
	}

	return card, nil
}
func (s *CardRepository) DeleteByID(gameID, collectionID, deckID string, cardID int64) error {
	// Check if the deck exists
	deck, err := s.deckRepository.GetByID(gameID, collectionID, deckID)
	if err != nil {
		return err
	}

	// Read info from file
	readCard, err := fs.OpenAndProcess(deck.Path(gameID, collectionID, s.cfg), fs.JsonFromReader[entity.Card])
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
	if err := fs.CreateAndProcess(deck.Path(gameID, collectionID, s.cfg), *readCard, fs.JsonToWriter[entity.Card]); err != nil {
		return err
	}
	return nil
}
func (s *CardRepository) GetImage(gameID, collectionID, deckID string, cardID int64) ([]byte, string, error) {
	// Check if such an object exists
	card, err := s.GetByID(gameID, collectionID, deckID, cardID)
	if err != nil {
		return nil, "", err
	}

	// Check if an image exists
	isExist, err := fs.IsFileExist(card.ImagePath(gameID, collectionID, deckID, s.cfg))
	if err != nil {
		return nil, "", err
	}
	if !isExist {
		return nil, "", errors.CardImageNotExists
	}

	// Read an image from a file
	data, err := fs.OpenAndProcess(card.ImagePath(gameID, collectionID, deckID, s.cfg), fs.BinFromReader)
	if err != nil {
		return nil, "", err
	}

	imgType, err := images.ValidateImage(data)
	if err != nil {
		return nil, "", err
	}

	return data, imgType, nil
}
func (s *CardRepository) CreateImage(gameID, collectionID, deckID string, cardID int64, imageURL string) error {
	// Check if such an object exists
	card, _ := s.GetByID(gameID, collectionID, deckID, cardID)
	if card == nil {
		return errors.CardNotExists.HTTP(http.StatusBadRequest)
	}

	// Download image
	imageBytes, err := network.DownloadBytes(imageURL)
	if err != nil {
		return err
	}

	// Validate image
	_, err = images.ValidateImage(imageBytes)
	if err != nil {
		return err
	}

	// Write image to file
	return fs.CreateAndProcess(card.ImagePath(gameID, collectionID, deckID, s.cfg), imageBytes, fs.BinToWriter)
}

// Internal function. Get map of cards inside deck
func (s *CardRepository) getCardsMap(gameID, collectionID, deckID string) (map[int64]*entity.CardInfo, error) {
	// Check if the deck exists
	deck, err := s.deckRepository.GetByID(gameID, collectionID, deckID)
	if err != nil {
		return nil, err
	}

	// Read info from file
	readCard, err := fs.OpenAndProcess(deck.Path(gameID, collectionID, s.cfg), fs.JsonFromReader[entity.Card])
	if err != nil {
		return nil, err
	}
	return readCard.Cards, nil
}
