package decks

import (
	"log"
	"net/http"
	"path/filepath"
	"sort"

	"tts_deck_build/internal/collections"
	"tts_deck_build/internal/config"
	"tts_deck_build/internal/errors"
	"tts_deck_build/internal/fs"
	"tts_deck_build/internal/images"
	"tts_deck_build/internal/network"
	"tts_deck_build/internal/utils"
)

type DeckStorage struct {
	Config            *config.Config
	CollectionService *collections.CollectionService
}

func NewDeckStorage(config *config.Config, collectionService *collections.CollectionService) *DeckStorage {
	return &DeckStorage{
		Config:            config,
		CollectionService: collectionService,
	}
}

func (s *DeckStorage) Create(gameId, collectionId string, deck *DeckInfo) (*DeckInfo, error) {
	// Convert name to ID
	deck.Id = utils.NameToId(deck.Type)
	if len(deck.Id) == 0 {
		return nil, errors.BadName.AddMessage(deck.Type)
	}

	// Check if such an object already exists
	if val, _ := s.GetById(gameId, collectionId, deck.Id); val != nil {
		return nil, errors.DeckExist
	}

	// Build path
	deckPath := filepath.Join(s.Config.Games(), gameId, collectionId, deck.Id+".json")

	// Writing info to file
	if err := fs.WriteFile(deckPath, deck); err != nil {
		return nil, err
	}

	if len(deck.BacksideImage) > 0 {
		// Download image
		if err := s.CreateImage(gameId, collectionId, deck.Id, deck.BacksideImage); err != nil {
			return nil, err
		}
	}

	return deck, nil
}
func (s *DeckStorage) GetById(gameId, collectionId, deckId string) (*DeckInfo, error) {
	// Check if the collection exists
	_, err := s.CollectionService.Item(gameId, collectionId)
	if err != nil {
		return nil, err
	}

	// Build path
	deckPath := filepath.Join(s.Config.Games(), gameId, collectionId, deckId+".json")

	// Check if such an object exists
	isExist, err := fs.IsFileExist(deckPath)
	if err != nil {
		return nil, err
	}
	if !isExist {
		return nil, errors.DeckNotExists
	}

	// Read info from file
	return fs.ReadFile[DeckInfo](deckPath)
}
func (s *DeckStorage) GetAll(gameId, collectionId string) ([]*DeckInfo, error) {
	decks := make([]*DeckInfo, 0)

	// Check if the collection exists
	_, err := s.CollectionService.Item(gameId, collectionId)
	if err != nil {
		return decks, err
	}

	// Get list of objects
	folders, err := fs.ListOfFiles(filepath.Join(s.Config.Games(), gameId, collectionId))
	if err != nil {
		return decks, err
	}

	// Get each deck
	for _, deckFileName := range folders {
		deckId := fs.GetFilenameWithoutExt(deckFileName)
		deck, err := s.GetById(gameId, collectionId, deckId)
		if err != nil {
			log.Println(err.Error())
			continue
		}
		decks = append(decks, deck)
	}

	return decks, nil
}
func (s *DeckStorage) Update(gameId, collectionId, deckId string, dto *UpdateDeckDTO) (*DeckInfo, error) {
	// Get old object
	oldDeck, err := s.GetById(gameId, collectionId, deckId)
	if err != nil {
		return nil, err
	}

	// Convert name to ID
	newDeckId := utils.NameToId(dto.Type)
	if len(newDeckId) == 0 {
		return nil, errors.BadName.AddMessage(dto.Type)
	}

	// Create deck object
	deck := NewDeckInfo(newDeckId, dto.Type, dto.BacksideImage)

	// Build path
	newDeckPath := filepath.Join(s.Config.Games(), gameId, collectionId, deck.Id+".json")

	// If the id has been changed, rename the object
	if deck.Id != oldDeck.Id {
		// Check if such an object already exists
		if val, _ := s.GetById(gameId, collectionId, deck.Id); val != nil {
			return nil, errors.DeckExist
		}

		// If image exist, rename
		if data, _, _ := s.GetImage(gameId, collectionId, oldDeck.Id); data != nil {
			oldDeckImagePath := filepath.Join(s.Config.Games(), gameId, collectionId, oldDeck.Id+".bin")
			newDeckImagePath := filepath.Join(s.Config.Games(), gameId, collectionId, deck.Id+".bin")
			err = fs.MoveFolder(oldDeckImagePath, newDeckImagePath)
			if err != nil {
				return nil, err
			}
		}

		// Build path
		oldDeckPath := filepath.Join(s.Config.Games(), gameId, collectionId, oldDeck.Id+".json")

		// Rename object
		err = fs.MoveFolder(oldDeckPath, newDeckPath)
		if err != nil {
			return nil, err
		}
	}

	// If the image has been changed
	if deck.BacksideImage != oldDeck.BacksideImage {
		// Build image path
		deckImagePath := filepath.Join(s.Config.Games(), gameId, collectionId, deck.Id+".bin")

		// If image exist, delete
		if data, _, _ := s.GetImage(gameId, collectionId, deck.Id); data != nil {
			err = fs.RemoveFile(deckImagePath)
			if err != nil {
				return nil, err
			}
		}

		if len(deck.BacksideImage) > 0 {
			// Download image
			if err = s.CreateImage(gameId, collectionId, deck.Id, deck.BacksideImage); err != nil {
				return nil, err
			}
		}
	}

	// If the object has been changed, update the object file
	if !oldDeck.Compare(deck) {
		// Writing info to file
		if err = fs.WriteFile(newDeckPath, deck); err != nil {
			return nil, err
		}
		return deck, nil
	}

	return oldDeck, nil
}
func (s *DeckStorage) DeleteById(gameId, collectionId, deckId string) error {
	// Build path
	deckPath := filepath.Join(s.Config.Games(), gameId, collectionId, deckId+".json")

	// Check if such an object exists
	if val, _ := s.GetById(gameId, collectionId, deckId); val == nil {
		return errors.DeckNotExists.HTTP(http.StatusBadRequest)
	}

	// Remove object
	if err := fs.RemoveFile(deckPath); err != nil {
		return err
	}

	// Remove image
	deckImagePath := filepath.Join(s.Config.Games(), gameId, collectionId, deckId+".bin")
	if err := fs.RemoveFile(deckImagePath); err != nil {
		return err
	}

	return nil
}
func (s *DeckStorage) GetImage(gameId, collectionId, deckId string) ([]byte, string, error) {
	// Check if such an object exists
	_, err := s.GetById(gameId, collectionId, deckId)
	if err != nil {
		return nil, "", err
	}

	// Build image path
	deckImagePath := filepath.Join(s.Config.Games(), gameId, collectionId, deckId+".bin")

	// Check if an image exists
	isExist, err := fs.IsFileExist(deckImagePath)
	if err != nil {
		return nil, "", err
	}
	if !isExist {
		return nil, "", errors.DeckImageNotExists
	}

	// Read an image from a file
	data, err := fs.ReadBinaryFile(deckImagePath)
	if err != nil {
		return nil, "", err
	}

	imgType, err := images.ValidateImage(data)
	if err != nil {
		return nil, "", err
	}

	return data, imgType, nil
}
func (s *DeckStorage) CreateImage(gameId, collectionId, deckId, imageUrl string) error {
	// Check if such an object exists
	if val, _ := s.GetById(gameId, collectionId, deckId); val == nil {
		return errors.DeckNotExists.HTTP(http.StatusBadRequest)
	}

	// Download image
	imageBytes, err := network.DownloadBytes(imageUrl)
	if err != nil {
		return err
	}

	// Validate image
	_, err = images.ValidateImage(imageBytes)
	if err != nil {
		return err
	}

	// Build image path
	deckImagePath := filepath.Join(s.Config.Games(), gameId, collectionId, deckId+".bin")

	// Write image to file
	return fs.WriteBinaryFile(deckImagePath, imageBytes)
}
func (s *DeckStorage) GetAllDecksInGame(gameId string) ([]*DeckInfo, error) {
	decks := make([]*DeckInfo, 0)

	// Get all collections in selected game
	listCollections, err := s.CollectionService.List(gameId)
	if err != nil {
		return decks, err
	}

	// Mark unique deck types
	uniqueDecks := make(map[string]struct{})

	// Go through all collections and find unique types of decks
	for _, collection := range listCollections {
		// Get all decks in selected collection
		collectionDecks, err := s.GetAll(gameId, collection.Id)
		if err != nil {
			return decks, err
		}

		// Go through all decks and keep only unique decks
		for _, deck := range collectionDecks {
			if _, ok := uniqueDecks[deck.Type+deck.BacksideImage]; ok {
				// If we have already seen such a deck, we skip it
				continue
			}
			// If deck unique, put mark in map
			uniqueDecks[deck.Type+deck.BacksideImage] = struct{}{}
			decks = append(decks, deck)
		}
	}

	// Sort decks in result
	sort.SliceStable(decks, func(i, j int) bool {
		return decks[i].Type < decks[j].Type
	})
	return decks, nil
}
