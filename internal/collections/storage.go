package collections

import (
	"log"
	"net/http"
	"path/filepath"

	"tts_deck_build/internal/config"
	"tts_deck_build/internal/errors"
	"tts_deck_build/internal/fs"
	"tts_deck_build/internal/games"
	"tts_deck_build/internal/images"
	"tts_deck_build/internal/network"
	"tts_deck_build/internal/utils"
)

type CollectionStorage struct {
	Config      *config.Config
	GameService *games.GameService
}

func NewCollectionStorage(config *config.Config, gameService *games.GameService) *CollectionStorage {
	return &CollectionStorage{
		Config:      config,
		GameService: gameService,
	}
}

func (s *CollectionStorage) Create(gameId string, collection *CollectionInfo) (*CollectionInfo, error) {
	// Convert name to ID
	collection.Id = utils.NameToId(collection.Name)
	if len(collection.Id) == 0 {
		return nil, errors.BadName.AddMessage(collection.Name)
	}

	// Check if such an object already exists
	if val, _ := s.GetById(gameId, collection.Id); val != nil {
		return nil, errors.CollectionExist
	}

	// Build path
	collectionPath := filepath.Join(s.Config.Games(), gameId, collection.Id)

	// Create folder
	if err := fs.CreateFolder(collectionPath); err != nil {
		return nil, err
	}

	// Build info path
	gameInfoPath := filepath.Join(collectionPath, s.Config.InfoFilename)

	// Writing info to file
	if err := fs.WriteFile(gameInfoPath, collection); err != nil {
		return nil, err
	}

	if len(collection.Image) > 0 {
		// Download image
		if err := s.CreateImage(gameId, collection.Id, collection.Image); err != nil {
			return nil, err
		}
	}

	return collection, nil
}
func (s *CollectionStorage) GetById(gameId, collectionId string) (*CollectionInfo, error) {
	// Check if the game exists
	_, err := s.GameService.Item(gameId)
	if err != nil {
		return nil, err
	}

	// Build path
	collectionPath := filepath.Join(s.Config.Games(), gameId, collectionId)

	// Check if such an object exists
	isExist, err := fs.IsFolderExist(collectionPath)
	if err != nil {
		return nil, err
	}
	if !isExist {
		return nil, errors.CollectionNotExists
	}

	// Build info path
	collectionInfoPath := filepath.Join(collectionPath, s.Config.InfoFilename)

	// Check if such an object exists
	isExist, err = fs.IsFileExist(collectionInfoPath)
	if err != nil {
		return nil, err
	}
	if !isExist {
		return nil, errors.CollectionInfoNotExists
	}

	// Read info from file
	return fs.ReadFile[CollectionInfo](collectionInfoPath)
}
func (s *CollectionStorage) GetAll(gameId string) ([]*CollectionInfo, error) {
	// Check if the game exists
	_, err := s.GameService.Item(gameId)
	if err != nil {
		return nil, err
	}

	// Get list of objects
	folders, err := fs.ListOfFolders(filepath.Join(s.Config.Games(), gameId))
	if err != nil {
		return nil, err
	}

	// Get each collection
	collections := make([]*CollectionInfo, 0)
	for _, collectionId := range folders {
		collection, err := s.GetById(gameId, collectionId)
		if err != nil {
			log.Println(err.Error())
			continue
		}
		collections = append(collections, collection)
	}

	return collections, nil
}
func (s *CollectionStorage) Update(gameId, collectionId string, dto *UpdateCollectionDTO) (*CollectionInfo, error) {
	// Get old object
	oldCollection, err := s.GetById(gameId, collectionId)
	if err != nil {
		return nil, err
	}

	// Convert name to ID
	newCollectionId := utils.NameToId(dto.Name)
	if len(newCollectionId) == 0 {
		return nil, errors.BadName.AddMessage(dto.Name)
	}

	// Create collection object
	collection := NewCollectionInfo(newCollectionId, dto.Name, dto.Description, dto.Image)

	// Build path
	newCollectionPath := filepath.Join(s.Config.Games(), gameId, collection.Id)

	// If the id has been changed, rename the object
	if collection.Id != oldCollection.Id {
		// Check if such an object already exists
		if val, _ := s.GetById(gameId, collection.Id); val != nil {
			return nil, errors.CollectionExist
		}

		// Build path
		oldCollectionPath := filepath.Join(s.Config.Games(), gameId, oldCollection.Id)

		// Rename object
		err = fs.MoveFolder(oldCollectionPath, newCollectionPath)
		if err != nil {
			return nil, err
		}
	}

	// If the image has been changed
	if collection.Image != oldCollection.Image {
		// Build image path
		collectionImagePath := filepath.Join(s.Config.Games(), gameId, collection.Id, s.Config.ImageFilename)

		// If image exist, delete
		if data, _, _ := s.GetImage(gameId, collection.Id); data != nil {
			err = fs.RemoveFile(collectionImagePath)
			if err != nil {
				return nil, err
			}
		}

		if len(collection.Image) > 0 {
			// Download image
			if err = s.CreateImage(gameId, collection.Id, collection.Image); err != nil {
				return nil, err
			}
		}
	}

	// If the object has been changed, update the info file
	if !oldCollection.Compare(collection) {
		// Build info path
		collectionInfoPath := filepath.Join(newCollectionPath, s.Config.InfoFilename)

		// Writing info to file
		if err = fs.WriteFile(collectionInfoPath, collection); err != nil {
			return nil, err
		}
		return collection, nil
	}

	return oldCollection, nil
}
func (s *CollectionStorage) DeleteById(gameId, collectionId string) error {
	// Build path
	collectionPath := filepath.Join(s.Config.Games(), gameId, collectionId)

	// Check if such an object exists
	if val, _ := s.GetById(gameId, collectionId); val == nil {
		return errors.CollectionNotExists.HTTP(http.StatusBadRequest)
	}

	// Remove object
	return fs.RemoveFolder(collectionPath)
}
func (s *CollectionStorage) GetImage(gameId, collectionId string) ([]byte, string, error) {
	// Check if such an object exists
	_, err := s.GetById(gameId, collectionId)
	if err != nil {
		return nil, "", err
	}

	// Build image path
	collectionImagePath := filepath.Join(s.Config.Games(), gameId, collectionId, s.Config.ImageFilename)

	// Check if an image exists
	isExist, err := fs.IsFileExist(collectionImagePath)
	if err != nil {
		return nil, "", err
	}
	if !isExist {
		return nil, "", errors.CollectionImageNotExists
	}

	// Read an image from a file
	data, err := fs.ReadBinaryFile(collectionImagePath)
	if err != nil {
		return nil, "", err
	}

	imgType, err := images.ValidateImage(data)
	if err != nil {
		return nil, "", err
	}

	return data, imgType, nil
}
func (s *CollectionStorage) CreateImage(gameId, collectionId, imageUrl string) error {
	// Check if such an object exists
	if val, _ := s.GetById(gameId, collectionId); val == nil {
		return errors.CollectionNotExists.HTTP(http.StatusBadRequest)
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
	collectionImagePath := filepath.Join(s.Config.Games(), gameId, collectionId, s.Config.ImageFilename)

	// Write image to file
	return fs.WriteBinaryFile(collectionImagePath, imageBytes)
}
