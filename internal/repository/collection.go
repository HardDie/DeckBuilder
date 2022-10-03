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
	"tts_deck_build/internal/logger"
	"tts_deck_build/internal/network"
	"tts_deck_build/internal/utils"
)

type ICollectionRepository interface {
	Create(gameID string, collection *entity.CollectionInfo) (*entity.CollectionInfo, error)
	GetByID(gameID, collectionID string) (*entity.CollectionInfo, error)
	GetAll(gameID string) ([]*entity.CollectionInfo, error)
	Update(gameID, collectionID string, dtoObject *dto.UpdateCollectionDTO) (*entity.CollectionInfo, error)
	DeleteByID(gameID, collectionID string) error
	GetImage(gameID, collectionID string) ([]byte, string, error)
	CreateImage(gameID, collectionID, imageURL string) error
}
type CollectionRepository struct {
	cfg            *config.Config
	gameRepository IGameRepository
}

func NewCollectionRepository(cfg *config.Config, gameRepository IGameRepository) *CollectionRepository {
	return &CollectionRepository{
		cfg:            cfg,
		gameRepository: gameRepository,
	}
}

func (s *CollectionRepository) Create(gameID string, collection *entity.CollectionInfo) (*entity.CollectionInfo, error) {
	// Check ID
	if collection.ID == "" {
		return nil, errors.BadName.AddMessage(collection.Name.String())
	}

	// Check if game exist
	if _, err := s.gameRepository.GetByID(gameID); err != nil {
		return nil, err
	}

	// Check if such an object already exists
	if val, _ := s.GetByID(gameID, collection.ID); val != nil {
		return nil, errors.CollectionExist
	}

	// Create folder
	if err := fs.CreateFolder(collection.Path(gameID, s.cfg)); err != nil {
		return nil, err
	}

	// Quote values before write to file
	collection.SetQuotedOutput()
	defer collection.SetRawOutput()

	// Writing info to file
	if err := fs.CreateAndProcess(collection.InfoPath(gameID, s.cfg), collection, fs.JsonToWriter[*entity.CollectionInfo]); err != nil {
		return nil, err
	}

	if collection.Image == "" {
		return collection, nil
	}

	// Download image
	if err := s.CreateImage(gameID, collection.ID, collection.Image); err != nil {
		return nil, err
	}

	// Writing info to file
	if err := fs.CreateAndProcess(collection.InfoPath(gameID, s.cfg), collection, fs.JsonToWriter[*entity.CollectionInfo]); err != nil {
		return nil, err
	}

	return collection, nil
}
func (s *CollectionRepository) GetByID(gameID, collectionID string) (*entity.CollectionInfo, error) {
	// Check if the game exists
	_, err := s.gameRepository.GetByID(gameID)
	if err != nil {
		return nil, err
	}

	collection := entity.CollectionInfo{ID: collectionID}

	// Check if such an object exists
	isExist, err := fs.IsFolderExist(collection.Path(gameID, s.cfg))
	if err != nil {
		return nil, err
	}
	if !isExist {
		return nil, errors.CollectionNotExists
	}

	// Check if such an object exists
	isExist, err = fs.IsFileExist(collection.InfoPath(gameID, s.cfg))
	if err != nil {
		return nil, err
	}
	if !isExist {
		return nil, errors.CollectionInfoNotExists
	}

	// Read info from file
	retCollection, err := fs.OpenAndProcess(collection.InfoPath(gameID, s.cfg), fs.JsonFromReader[entity.CollectionInfo])
	if err != nil {
		return nil, err
	}

	retCollection.CachedImage = fmt.Sprintf(s.cfg.CollectionImagePath, gameID, collectionID)
	return retCollection, nil
}
func (s *CollectionRepository) GetAll(gameID string) ([]*entity.CollectionInfo, error) {
	// Check if the game exists
	game, err := s.gameRepository.GetByID(gameID)
	if err != nil {
		return make([]*entity.CollectionInfo, 0), err
	}

	// Get list of objects
	folders, err := fs.ListOfFolders(game.Path(s.cfg))
	if err != nil {
		return make([]*entity.CollectionInfo, 0), err
	}

	// Get each collection
	collections := make([]*entity.CollectionInfo, 0)
	for _, collectionID := range folders {
		collection, err := s.GetByID(gameID, collectionID)
		if err != nil {
			logger.Error.Println(err.Error())
			continue
		}
		collections = append(collections, collection)
	}

	return collections, nil
}
func (s *CollectionRepository) Update(gameID, collectionID string, dtoObject *dto.UpdateCollectionDTO) (*entity.CollectionInfo, error) {
	// Get old object
	oldCollection, err := s.GetByID(gameID, collectionID)
	if err != nil {
		return nil, err
	}

	// Create collection object
	if dtoObject.Name == "" {
		dtoObject.Name = oldCollection.Name.String()
	}
	collection := entity.NewCollectionInfo(dtoObject.Name, dtoObject.Description, dtoObject.Image)
	collection.CreatedAt = oldCollection.CreatedAt
	if collection.ID == "" {
		return nil, errors.BadName.AddMessage(dtoObject.Name)
	}

	// If the id has been changed, rename the object
	if collection.ID != oldCollection.ID {
		// Check if such an object already exists
		if val, _ := s.GetByID(gameID, collection.ID); val != nil {
			return nil, errors.CollectionExist
		}

		// Rename object
		err = fs.MoveFolder(oldCollection.Path(gameID, s.cfg), collection.Path(gameID, s.cfg))
		if err != nil {
			return nil, err
		}
	}

	// If the object has been changed, update the info file
	if !oldCollection.Compare(collection) {
		// Quote values before write to file
		collection.SetQuotedOutput()
		defer collection.SetRawOutput()

		collection.UpdatedAt = utils.Allocate(time.Now())
		// Writing info to file
		if err = fs.CreateAndProcess(collection.InfoPath(gameID, s.cfg), collection, fs.JsonToWriter[*entity.CollectionInfo]); err != nil {
			return nil, err
		}
	}

	// If the image has not been changed
	if collection.Image == oldCollection.Image {
		return collection, nil
	}

	// If image exist, delete
	if data, _, _ := s.GetImage(gameID, collection.ID); data != nil {
		err = fs.RemoveFile(collection.ImagePath(gameID, s.cfg))
		if err != nil {
			return nil, err
		}
	}

	if collection.Image == "" {
		return collection, nil
	}

	// Download image
	if err = s.CreateImage(gameID, collection.ID, collection.Image); err != nil {
		return nil, err
	}

	return collection, nil
}
func (s *CollectionRepository) DeleteByID(gameID, collectionID string) error {
	collection := entity.CollectionInfo{ID: collectionID}

	// Check if such an object exists
	if val, _ := s.GetByID(gameID, collectionID); val == nil {
		return errors.CollectionNotExists.HTTP(http.StatusBadRequest)
	}

	// Remove object
	return fs.RemoveFolder(collection.Path(gameID, s.cfg))
}
func (s *CollectionRepository) GetImage(gameID, collectionID string) ([]byte, string, error) {
	// Check if such an object exists
	collection, err := s.GetByID(gameID, collectionID)
	if err != nil {
		return nil, "", err
	}

	// Check if an image exists
	isExist, err := fs.IsFileExist(collection.ImagePath(gameID, s.cfg))
	if err != nil {
		return nil, "", err
	}
	if !isExist {
		return nil, "", errors.CollectionImageNotExists
	}

	// Read an image from a file
	data, err := fs.OpenAndProcess(collection.ImagePath(gameID, s.cfg), fs.BinFromReader)
	if err != nil {
		return nil, "", err
	}

	imgType, err := images.ValidateImage(data)
	if err != nil {
		return nil, "", err
	}

	return data, imgType, nil
}
func (s *CollectionRepository) CreateImage(gameID, collectionID, imageURL string) error {
	// Check if such an object exists
	collection, _ := s.GetByID(gameID, collectionID)
	if collection == nil {
		return errors.CollectionNotExists.HTTP(http.StatusBadRequest)
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
	return fs.CreateAndProcess(collection.ImagePath(gameID, s.cfg), imageBytes, fs.BinToWriter)
}
