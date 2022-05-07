package games

import (
	"log"
	"net/http"
	"path/filepath"

	"tts_deck_build/internal/config"
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
	// Convert name to ID
	game.Id = utils.NameToId(game.Name)
	if len(game.Id) == 0 {
		return nil, errors.BadName.AddMessage(game.Name)
	}

	// Check if such an object already exists
	if val, _ := s.GetById(game.Id); val != nil {
		return nil, errors.GameExist
	}

	// Build path
	gamePath := filepath.Join(s.Config.Games(), game.Id)

	// Create folder
	if err := fs.CreateFolder(gamePath); err != nil {
		return nil, err
	}

	// Build info path
	gameInfoPath := filepath.Join(gamePath, s.Config.InfoFilename)

	// Writing info to file
	if err := fs.WriteFile(gameInfoPath, game); err != nil {
		return nil, err
	}

	if len(game.Image) > 0 {
		// Download image
		if err := s.CreateImage(game.Id, game.Image); err != nil {
			return nil, err
		}
	}

	return game, nil
}
func (s *GameStorage) GetById(gameId string) (*GameInfo, error) {
	// Build path
	gamePath := filepath.Join(s.Config.Games(), gameId)

	// Check if such an object exists
	isExist, err := fs.IsFolderExist(gamePath)
	if err != nil {
		return nil, err
	}
	if !isExist {
		return nil, errors.GameNotExists
	}

	// Build info path
	gameInfoPath := filepath.Join(gamePath, s.Config.InfoFilename)

	// Check if such an object exists
	isExist, err = fs.IsFileExist(gameInfoPath)
	if err != nil {
		return nil, err
	}
	if !isExist {
		return nil, errors.GameInfoNotExists
	}

	// Read info from file
	return fs.ReadFile[GameInfo](gameInfoPath)
}
func (s *GameStorage) GetAll() ([]*GameInfo, error) {
	// Get list of objects
	folders, err := fs.ListOfFolders(s.Config.Games())
	if err != nil {
		return nil, err
	}

	// Get each game
	games := make([]*GameInfo, 0)
	for _, gameId := range folders {
		game, err := s.GetById(gameId)
		if err != nil {
			log.Println(err.Error())
			continue
		}
		games = append(games, game)
	}

	return games, nil
}
func (s *GameStorage) Update(gameId string, dto *UpdateGameDTO) (*GameInfo, error) {
	// Get old object
	oldGame, err := s.GetById(gameId)
	if err != nil {
		return nil, err
	}

	// Convert name to ID
	newGameId := utils.NameToId(dto.Name)
	if len(newGameId) == 0 {
		return nil, errors.BadName.AddMessage(dto.Name)
	}

	// Create game object
	game := NewGameInfo(newGameId, dto.Name, dto.Description, dto.Image)

	// Build path
	newGamePath := filepath.Join(s.Config.Games(), game.Id)

	// If the id has been changed, rename the object
	if game.Id != oldGame.Id {
		// Check if such an object already exists
		if val, _ := s.GetById(game.Id); val != nil {
			return nil, errors.GameExist
		}

		// Build path
		oldGamePath := filepath.Join(s.Config.Games(), oldGame.Id)

		// Rename object
		err = fs.MoveFolder(oldGamePath, newGamePath)
		if err != nil {
			return nil, err
		}
	}

	// If the image has been changed
	if game.Image != oldGame.Image {
		// Build image path
		gameImagePath := filepath.Join(s.Config.Games(), game.Id, s.Config.ImageFilename)

		// If image exist, delete
		if data, _, _ := s.GetImage(game.Id); data != nil {
			err = fs.RemoveFile(gameImagePath)
			if err != nil {
				return nil, err
			}
		}

		if len(game.Image) > 0 {
			// Download image
			if err = s.CreateImage(game.Id, game.Image); err != nil {
				return nil, err
			}
		}
	}

	// If the object has been changed, update the info file
	if !oldGame.Compare(game) {
		// Build info path
		gameInfoPath := filepath.Join(newGamePath, s.Config.InfoFilename)

		// Writing info to file
		if err = fs.WriteFile(gameInfoPath, game); err != nil {
			return nil, err
		}
		return game, nil
	}

	return oldGame, nil
}
func (s *GameStorage) DeleteById(gameId string) error {
	// Build path
	gamePath := filepath.Join(s.Config.Games(), gameId)

	// Check if such an object exists
	if val, _ := s.GetById(gameId); val == nil {
		return errors.GameNotExists.HTTP(http.StatusBadRequest)
	}

	// Remove object
	return fs.RemoveFolder(gamePath)
}
func (s *GameStorage) GetImage(gameId string) ([]byte, string, error) {
	// Check if such an object exists
	_, err := s.GetById(gameId)
	if err != nil {
		return nil, "", err
	}

	// Build image path
	gameImagePath := filepath.Join(s.Config.Games(), gameId, s.Config.ImageFilename)

	// Check if an image exists
	isExist, err := fs.IsFileExist(gameImagePath)
	if err != nil {
		return nil, "", err
	}
	if !isExist {
		return nil, "", errors.GameImageNotExists
	}

	// Read an image from a file
	data, err := fs.ReadBinaryFile(gameImagePath)
	if err != nil {
		return nil, "", err
	}

	imgType, err := images.ValidateImage(data)
	if err != nil {
		return nil, "", err
	}

	return data, imgType, nil
}
func (s *GameStorage) CreateImage(gameId, imageUrl string) error {
	// Check if such an object exists
	if val, _ := s.GetById(gameId); val == nil {
		return errors.GameNotExists.HTTP(http.StatusBadRequest)
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
	gameImagePath := filepath.Join(s.Config.Games(), gameId, s.Config.ImageFilename)

	// Write image to file
	return fs.WriteBinaryFile(gameImagePath, imageBytes)
}
