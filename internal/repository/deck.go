package repository

import (
	"net/http"
	"time"

	"tts_deck_build/internal/config"
	"tts_deck_build/internal/dto"
	"tts_deck_build/internal/entity"
	"tts_deck_build/internal/errors"
	"tts_deck_build/internal/fs"
	"tts_deck_build/internal/images"
	"tts_deck_build/internal/logger"
	"tts_deck_build/internal/network"
	"tts_deck_build/internal/utils"
)

type IDeckRepository interface {
	Create(gameID, collectionID string, deck *entity.DeckInfo) (*entity.DeckInfo, error)
	GetByID(gameID, collectionID, deckID string) (*entity.DeckInfo, error)
	GetAll(gameID, collectionID string) ([]*entity.DeckInfo, error)
	Update(gameID, collectionID, deckID string, dtoObject *dto.UpdateDeckDTO) (*entity.DeckInfo, error)
	DeleteByID(gameID, collectionID, deckID string) error
	GetImage(gameID, collectionID, deckID string) ([]byte, string, error)
	CreateImage(gameID, collectionID, deckID, imageURL string) error
	GetAllDecksInGame(gameID string) ([]*entity.DeckInfo, error)
}
type DeckRepository struct {
	cfg                  *config.Config
	collectionRepository ICollectionRepository
}

func NewDeckRepository(cfg *config.Config, collectionRepository ICollectionRepository) *DeckRepository {
	return &DeckRepository{
		cfg:                  cfg,
		collectionRepository: collectionRepository,
	}
}

func (s *DeckRepository) Create(gameID, collectionID string, deck *entity.DeckInfo) (*entity.DeckInfo, error) {
	// Check ID
	if deck.ID == "" {
		return nil, errors.BadName.AddMessage(deck.Name.String())
	}

	// Check if collection exist
	if _, err := s.collectionRepository.GetByID(gameID, collectionID); err != nil {
		return nil, err
	}

	// Check if such an object already exists
	if val, _ := s.GetByID(gameID, collectionID, deck.ID); val != nil {
		return nil, errors.DeckExist
	}

	// Quote values before write to file
	deck.SetQuotedOutput()
	defer deck.SetRawOutput()

	// Writing info to file
	if err := fs.CreateAndProcess(deck.Path(gameID, collectionID, s.cfg), entity.Deck{Deck: deck}, fs.JsonToWriter[entity.Deck]); err != nil {
		return nil, err
	}

	// Create folder for card images
	if err := fs.CreateFolder(deck.CardImagesPath(gameID, collectionID, s.cfg)); err != nil {
		return nil, err
	}

	if len(deck.Image) > 0 {
		// Download image
		if err := s.CreateImage(gameID, collectionID, deck.ID, deck.Image); err != nil {
			return nil, err
		}
	}

	return deck, nil
}
func (s *DeckRepository) GetByID(gameID, collectionID, deckID string) (*entity.DeckInfo, error) {
	deck, err := s.getDeck(gameID, collectionID, deckID)
	if err != nil {
		return nil, err
	}
	return deck.Deck, nil
}
func (s *DeckRepository) GetAll(gameID, collectionID string) ([]*entity.DeckInfo, error) {
	decks := make([]*entity.DeckInfo, 0)

	// Check if the collection exists
	collection, err := s.collectionRepository.GetByID(gameID, collectionID)
	if err != nil {
		return decks, err
	}

	// Get list of objects
	folders, err := fs.ListOfFiles(collection.Path(gameID, s.cfg))
	if err != nil {
		return decks, err
	}

	// Get each deck
	for _, deckFileName := range folders {
		deckID := fs.GetFilenameWithoutExt(deckFileName)
		deck, err := s.GetByID(gameID, collectionID, deckID)
		if err != nil {
			logger.Error.Println(err.Error())
			continue
		}
		if deck == nil {
			logger.Warn.Println("Invalid deck file:", deckFileName)
			continue
		}
		decks = append(decks, deck)
	}

	return decks, nil
}
func (s *DeckRepository) Update(gameID, collectionID, deckID string, dtoObject *dto.UpdateDeckDTO) (*entity.DeckInfo, error) {
	// Get old object
	oldDeck, err := s.getDeck(gameID, collectionID, deckID)
	if err != nil {
		return nil, err
	}

	// Create deck object
	if dtoObject.Name == "" {
		dtoObject.Name = oldDeck.Deck.Name.String()
	}
	deck := entity.NewDeckInfo(dtoObject.Name, dtoObject.Image)
	deck.CreatedAt = oldDeck.Deck.CreatedAt
	if deck.ID == "" {
		return nil, errors.BadName.AddMessage(dtoObject.Name)
	}

	// If the id has been changed, rename the object
	if deck.ID != oldDeck.Deck.ID {
		// Check if such an object already exists
		if val, _ := s.GetByID(gameID, collectionID, deck.ID); val != nil {
			return nil, errors.DeckExist
		}

		// If image exist, rename
		if data, _, _ := s.GetImage(gameID, collectionID, oldDeck.Deck.ID); data != nil {
			err = fs.MoveFolder(oldDeck.Deck.ImagePath(gameID, collectionID, s.cfg), deck.ImagePath(gameID, collectionID, s.cfg))
			if err != nil {
				return nil, err
			}
		}

		// Rename object
		err = fs.MoveFolder(oldDeck.Deck.Path(gameID, collectionID, s.cfg), deck.Path(gameID, collectionID, s.cfg))
		if err != nil {
			return nil, err
		}

		// Rename card images folder
		err = fs.MoveFolder(oldDeck.Deck.CardImagesPath(gameID, collectionID, s.cfg), deck.CardImagesPath(gameID, collectionID, s.cfg))
		if err != nil {
			return nil, err
		}
	}

	// If the object has been changed, update the object file
	if !oldDeck.Deck.Compare(deck) {
		// Quote values before write to file
		deck.SetQuotedOutput()
		defer deck.SetRawOutput()

		deck.UpdatedAt = utils.Allocate(time.Now())
		// Writing info to file
		if err := fs.CreateAndProcess(deck.Path(gameID, collectionID, s.cfg), entity.Deck{Deck: deck, Cards: oldDeck.Cards}, fs.JsonToWriter[entity.Deck]); err != nil {
			return nil, err
		}
	}

	// If the image has been changed
	if deck.Image != oldDeck.Deck.Image {
		// If image exist, delete
		if data, _, _ := s.GetImage(gameID, collectionID, deck.ID); data != nil {
			err = fs.RemoveFile(deck.ImagePath(gameID, collectionID, s.cfg))
			if err != nil {
				return nil, err
			}
		}

		if len(deck.Image) > 0 {
			// Download image
			if err = s.CreateImage(gameID, collectionID, deck.ID, deck.Image); err != nil {
				return nil, err
			}
		}
	}

	return deck, nil
}
func (s *DeckRepository) DeleteByID(gameID, collectionID, deckID string) error {
	deck := entity.DeckInfo{ID: deckID}

	// Check if such an object exists
	val, _ := s.GetByID(gameID, collectionID, deckID)
	if val == nil {
		return errors.DeckNotExists.HTTP(http.StatusBadRequest)
	}

	// Remove object
	if err := fs.RemoveFile(deck.Path(gameID, collectionID, s.cfg)); err != nil {
		return err
	}

	// Remove card images
	if err := fs.RemoveFolder(deck.CardImagesPath(gameID, collectionID, s.cfg)); err != nil {
		return err
	}

	// Remove image
	if val.Image != "" {
		return fs.RemoveFile(deck.ImagePath(gameID, collectionID, s.cfg))
	}
	return nil
}
func (s *DeckRepository) GetImage(gameID, collectionID, deckID string) ([]byte, string, error) {
	// Check if such an object exists
	deck, err := s.GetByID(gameID, collectionID, deckID)
	if err != nil {
		return nil, "", err
	}

	// Check if an image exists
	isExist, err := fs.IsFileExist(deck.ImagePath(gameID, collectionID, s.cfg))
	if err != nil {
		return nil, "", err
	}
	if !isExist {
		return nil, "", errors.DeckImageNotExists
	}

	// Read an image from a file
	data, err := fs.OpenAndProcess(deck.ImagePath(gameID, collectionID, s.cfg), fs.BinFromReader)
	if err != nil {
		return nil, "", err
	}

	imgType, err := images.ValidateImage(data)
	if err != nil {
		return nil, "", err
	}

	return data, imgType, nil
}
func (s *DeckRepository) CreateImage(gameID, collectionID, deckID, imageURL string) error {
	// Check if such an object exists
	deck, _ := s.GetByID(gameID, collectionID, deckID)
	if deck == nil {
		return errors.DeckNotExists.HTTP(http.StatusBadRequest)
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
	return fs.CreateAndProcess(deck.ImagePath(gameID, collectionID, s.cfg), imageBytes, fs.BinToWriter)
}
func (s *DeckRepository) GetAllDecksInGame(gameID string) ([]*entity.DeckInfo, error) {
	// Get all collections in selected game
	listCollections, err := s.collectionRepository.GetAll(gameID)
	if err != nil {
		return make([]*entity.DeckInfo, 0), err
	}

	// Mark unique deck types
	uniqueDecks := make(map[string]struct{})

	// Go through all collections and find unique types of decks
	decks := make([]*entity.DeckInfo, 0)
	for _, collection := range listCollections {
		// Get all decks in selected collection
		collectionDecks, err := s.GetAll(gameID, collection.ID)
		if err != nil {
			return make([]*entity.DeckInfo, 0), err
		}

		// Go through all decks and keep only unique decks
		for _, deck := range collectionDecks {
			if _, ok := uniqueDecks[deck.Name.String()+deck.Image]; ok {
				// If we have already seen such a deck, we skip it
				continue
			}
			// If deck unique, put mark in map
			uniqueDecks[deck.Name.String()+deck.Image] = struct{}{}
			decks = append(decks, deck)
		}
	}
	return decks, nil
}

func (s *DeckRepository) getDeck(gameID, collectionID, deckID string) (*entity.Deck, error) {
	// Check if the collection exists
	_, err := s.collectionRepository.GetByID(gameID, collectionID)
	if err != nil {
		return nil, err
	}

	deck := entity.DeckInfo{ID: deckID}

	// Check if such an object exists
	isExist, err := fs.IsFileExist(deck.Path(gameID, collectionID, s.cfg))
	if err != nil {
		return nil, err
	}
	if !isExist {
		return nil, errors.DeckNotExists
	}

	// Read info from file
	readDeck, err := fs.OpenAndProcess(deck.Path(gameID, collectionID, s.cfg), fs.JsonFromReader[entity.Deck])
	if err != nil {
		return nil, err
	}
	return readDeck, nil
}
