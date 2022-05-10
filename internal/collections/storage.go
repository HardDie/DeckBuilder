package collections

import (
	"log"
	"net/http"
	"time"

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
	// Check ID
	if len(collection.Id) == 0 {
		return nil, errors.BadName.AddMessage(collection.Name)
	}

	// Check if such an object already exists
	if val, _ := s.GetById(gameId, collection.Id); val != nil {
		return nil, errors.CollectionExist
	}

	// Create folder
	if err := fs.CreateFolder(collection.Path(gameId)); err != nil {
		return nil, err
	}

	// Writing info to file
	if err := fs.WriteFile(collection.InfoPath(gameId), collection); err != nil {
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

	collection := CollectionInfo{Id: collectionId}

	// Check if such an object exists
	isExist, err := fs.IsFolderExist(collection.Path(gameId))
	if err != nil {
		return nil, err
	}
	if !isExist {
		return nil, errors.CollectionNotExists
	}

	// Check if such an object exists
	isExist, err = fs.IsFileExist(collection.InfoPath(gameId))
	if err != nil {
		return nil, err
	}
	if !isExist {
		return nil, errors.CollectionInfoNotExists
	}

	// Read info from file
	return fs.ReadFile[CollectionInfo](collection.InfoPath(gameId))
}
func (s *CollectionStorage) GetAll(gameId string) ([]*CollectionInfo, error) {
	// Check if the game exists
	game, err := s.GameService.Item(gameId)
	if err != nil {
		return make([]*CollectionInfo, 0), err
	}

	// Get list of objects
	folders, err := fs.ListOfFolders(game.Path())
	if err != nil {
		return make([]*CollectionInfo, 0), err
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

	// Create collection object
	collection := NewCollectionInfo(dto.Name, dto.Description, dto.Image)
	collection.CreatedAt = oldCollection.CreatedAt
	if len(collection.Id) == 0 {
		return nil, errors.BadName.AddMessage(dto.Name)
	}

	// If the id has been changed, rename the object
	if collection.Id != oldCollection.Id {
		// Check if such an object already exists
		if val, _ := s.GetById(gameId, collection.Id); val != nil {
			return nil, errors.CollectionExist
		}

		// Rename object
		err = fs.MoveFolder(oldCollection.Path(gameId), collection.Path(gameId))
		if err != nil {
			return nil, err
		}
	}

	// If the object has been changed, update the info file
	if !oldCollection.Compare(collection) {
		collection.UpdatedAt = utils.Allocate(time.Now())
		// Writing info to file
		if err = fs.WriteFile(collection.InfoPath(gameId), collection); err != nil {
			return nil, err
		}
	}

	// If the image has been changed
	if collection.Image != oldCollection.Image {
		// If image exist, delete
		if data, _, _ := s.GetImage(gameId, collection.Id); data != nil {
			err = fs.RemoveFile(collection.ImagePath(gameId))
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

	return collection, nil
}
func (s *CollectionStorage) DeleteById(gameId, collectionId string) error {
	collection := CollectionInfo{Id: collectionId}

	// Check if such an object exists
	if val, _ := s.GetById(gameId, collectionId); val == nil {
		return errors.CollectionNotExists.HTTP(http.StatusBadRequest)
	}

	// Remove object
	return fs.RemoveFolder(collection.Path(gameId))
}
func (s *CollectionStorage) GetImage(gameId, collectionId string) ([]byte, string, error) {
	// Check if such an object exists
	collection, err := s.GetById(gameId, collectionId)
	if err != nil {
		return nil, "", err
	}

	// Check if an image exists
	isExist, err := fs.IsFileExist(collection.ImagePath(gameId))
	if err != nil {
		return nil, "", err
	}
	if !isExist {
		return nil, "", errors.CollectionImageNotExists
	}

	// Read an image from a file
	data, err := fs.ReadBinaryFile(collection.ImagePath(gameId))
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
	collection, _ := s.GetById(gameId, collectionId)
	if collection == nil {
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

	// Write image to file
	return fs.WriteBinaryFile(collection.ImagePath(gameId), imageBytes)
}
