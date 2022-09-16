package games

import (
	"log"
	"net/http"
	"path/filepath"
	"time"

	"tts_deck_build/internal/config"
	"tts_deck_build/internal/dto"
	"tts_deck_build/internal/errors"
	"tts_deck_build/internal/fs"
	"tts_deck_build/internal/images"
	"tts_deck_build/internal/network"
	"tts_deck_build/internal/utils"
)

type GameStorage struct {
	Config *config.Config
}

func NewGameStorage(config *config.Config) *GameStorage {
	return &GameStorage{
		Config: config,
	}
}

func (s *GameStorage) Create(game *GameInfo) (*GameInfo, error) {
	// Check ID
	if game.ID == "" {
		return nil, errors.BadName.AddMessage(game.Name.String())
	}

	// Check if such an object already exists
	if val, _ := s.GetByID(game.ID); val != nil {
		return nil, errors.GameExist
	}

	// Create folder
	if err := fs.CreateFolder(game.Path()); err != nil {
		return nil, err
	}

	// Quote values before write to file
	game.SetQuotedOutput()
	defer game.SetRawOutput()

	// Writing info to file
	if err := fs.CreateAndProcess(game.InfoPath(), game, fs.JsonToWriter[*GameInfo]); err != nil {
		return nil, err
	}

	if len(game.Image) > 0 {
		// Download image
		if err := s.CreateImage(game.ID, game.Image); err != nil {
			return nil, err
		}
	}

	return game, nil
}
func (s *GameStorage) GetByID(gameID string) (*GameInfo, error) {
	game := GameInfo{ID: gameID}

	// Check if such an object exists
	isExist, err := fs.IsFolderExist(game.Path())
	if err != nil {
		return nil, err
	}
	if !isExist {
		return nil, errors.GameNotExists
	}

	// Check if such an object exists
	isExist, err = fs.IsFileExist(game.InfoPath())
	if err != nil {
		return nil, err
	}
	if !isExist {
		return nil, errors.GameInfoNotExists
	}

	// Read info from file
	retGame, err := fs.OpenAndProcess(game.InfoPath(), fs.JsonFromReader[GameInfo])
	if err != nil {
		return nil, err
	}

	return retGame, nil
}
func (s *GameStorage) GetAll() ([]*GameInfo, error) {
	isExist, err := fs.IsFolderExist(s.Config.Games())
	if err != nil {
		return make([]*GameInfo, 0), err
	}
	if !isExist {
		return make([]*GameInfo, 0), nil
	}

	// Get list of objects
	folders, err := fs.ListOfFolders(s.Config.Games())
	if err != nil {
		return make([]*GameInfo, 0), err
	}

	// Get each game
	games := make([]*GameInfo, 0)
	for _, gameID := range folders {
		game, err := s.GetByID(gameID)
		if err != nil {
			log.Println(err.Error())
			continue
		}
		games = append(games, game)
	}

	return games, nil
}
func (s *GameStorage) Update(gameID string, dtoObject *dto.UpdateGameDTO) (*GameInfo, error) {
	// Get old object
	oldGame, err := s.GetByID(gameID)
	if err != nil {
		return nil, err
	}

	// Create game object
	if dtoObject.Name == "" {
		dtoObject.Name = oldGame.Name.String()
	}
	game := NewGameInfo(dtoObject.Name, dtoObject.Description, dtoObject.Image)
	game.CreatedAt = oldGame.CreatedAt
	if game.ID == "" {
		return nil, errors.BadName.AddMessage(dtoObject.Name)
	}

	// If the id has been changed, rename the object
	if game.ID != oldGame.ID {
		// Check if such an object already exists
		if val, _ := s.GetByID(game.ID); val != nil {
			return nil, errors.GameExist
		}

		// Rename object
		err = fs.MoveFolder(oldGame.Path(), game.Path())
		if err != nil {
			return nil, err
		}
	}

	// If the object has been changed, update the info file
	if !oldGame.Compare(game) {
		// Quote values before write to file
		game.SetQuotedOutput()
		defer game.SetRawOutput()

		game.UpdatedAt = utils.Allocate(time.Now())
		// Writing info to file
		if err = fs.CreateAndProcess(game.InfoPath(), game, fs.JsonToWriter[*GameInfo]); err != nil {
			return nil, err
		}
	}

	// If the image has been changed
	if game.Image != oldGame.Image {
		// If image exist, delete
		if data, _, _ := s.GetImage(game.ID); data != nil {
			err = fs.RemoveFile(game.ImagePath())
			if err != nil {
				return nil, err
			}
		}

		if len(game.Image) > 0 {
			// Download image
			if err = s.CreateImage(game.ID, game.Image); err != nil {
				return nil, err
			}
		}
	}

	return game, nil
}
func (s *GameStorage) DeleteByID(gameID string) error {
	game := GameInfo{ID: gameID}

	// Check if such an object exists
	if val, _ := s.GetByID(gameID); val == nil {
		return errors.GameNotExists.HTTP(http.StatusBadRequest)
	}

	// Remove object
	return fs.RemoveFolder(game.Path())
}
func (s *GameStorage) GetImage(gameID string) ([]byte, string, error) {
	// Check if such an object exists
	game, err := s.GetByID(gameID)
	if err != nil {
		return nil, "", err
	}

	// Check if an image exists
	isExist, err := fs.IsFileExist(game.ImagePath())
	if err != nil {
		return nil, "", err
	}
	if !isExist {
		return nil, "", errors.GameImageNotExists
	}

	// Read an image from a file
	data, err := fs.OpenAndProcess(game.ImagePath(), fs.BinFromReader)
	if err != nil {
		return nil, "", err
	}

	imgType, err := images.ValidateImage(data)
	if err != nil {
		return nil, "", err
	}

	return data, imgType, nil
}
func (s *GameStorage) CreateImage(gameID, imageURL string) error {
	// Check if such an object exists
	game, _ := s.GetByID(gameID)
	if game == nil {
		return errors.GameNotExists.HTTP(http.StatusBadRequest)
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
	return fs.CreateAndProcess(game.ImagePath(), imageBytes, fs.BinToWriter)
}
func (s *GameStorage) Duplicate(gameID string, dtoObject *dto.DuplicateGameDTO) (*GameInfo, error) {
	// Check if the game exists
	oldGame, _ := s.GetByID(gameID)
	if oldGame == nil {
		return nil, errors.GameNotExists.HTTP(http.StatusBadRequest)
	}

	// New game object
	game := NewGameInfo(dtoObject.Name, oldGame.Description.String(), oldGame.Image)

	// Check ID
	if game.ID == "" {
		return nil, errors.BadName.AddMessage(game.Name.String())
	}

	// Check if such an object already exists
	if val, _ := s.GetByID(game.ID); val != nil {
		return nil, errors.GameExist
	}

	// Create a copy of the game
	err := fs.CopyFolder(oldGame.Path(), game.Path())
	if err != nil {
		return nil, err
	}

	// Quote values before write to file
	game.SetQuotedOutput()
	defer game.SetRawOutput()

	// Writing info to file
	if err = fs.CreateAndProcess(game.InfoPath(), game, fs.JsonToWriter[*GameInfo]); err != nil {
		return nil, err
	}

	return game, nil
}
func (s *GameStorage) Export(gameID string) ([]byte, error) {
	// Check if such an object exists
	game, _ := s.GetByID(gameID)
	if game == nil {
		return nil, errors.GameNotExists.HTTP(http.StatusBadRequest)
	}

	return fs.ArchiveFolder(game.Path(), game.ID)
}
func (s *GameStorage) Import(data []byte, name string) error {
	gameID := utils.NameToID(name)
	if name != "" && gameID == "" {
		return errors.BadName
	}

	// Unpack the archive
	resultGameID, err := fs.UnarchiveFolder(data, gameID)
	if err != nil {
		return err
	}

	// Check if the root folder contains information about the game
	game, err := s.GetByID(resultGameID)
	if err != nil {
		// Build a full relative path to the root game folder
		gameRootPath := filepath.Join(config.GetConfig().Games(), resultGameID)
		// If an error occurs during unzipping, delete the created folder with the game
		errors.IfErrorLog(fs.RemoveFolder(gameRootPath))
		return err
	}

	// If the user skipped passing a new name for the game,
	// but the root folder has a different name than in the game information file.
	// Fix the game information file.
	if name == "" && resultGameID != game.ID {
		gameID = resultGameID
		name = resultGameID
	}

	// If the name has been changed
	if name != "" {
		// Update the title of the game
		game.ID = gameID
		game.Name = utils.NewQuotedString(name)

		// Quote values before write to file
		game.SetQuotedOutput()
		defer game.SetRawOutput()

		// Writing info to file
		if err = fs.CreateAndProcess(game.InfoPath(), game, fs.JsonToWriter[*GameInfo]); err != nil {
			return err
		}
	}

	return nil
}
