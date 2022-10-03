package repository

import (
	"fmt"
	"net/http"
	"path/filepath"
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

type IGameRepository interface {
	Create(game *entity.GameInfo) (*entity.GameInfo, error)
	GetByID(gameID string) (*entity.GameInfo, error)
	GetAll() ([]*entity.GameInfo, error)
	Update(gameID string, dtoObject *dto.UpdateGameDTO) (*entity.GameInfo, error)
	DeleteByID(gameID string) error
	GetImage(gameID string) ([]byte, string, error)
	CreateImage(gameID, imageURL string) error
	Duplicate(gameID string, dtoObject *dto.DuplicateGameDTO) (*entity.GameInfo, error)
	Export(gameID string) ([]byte, error)
	Import(data []byte, name string) error
}
type GameRepository struct {
	cfg *config.Config
}

func NewGameRepository(cfg *config.Config) *GameRepository {
	return &GameRepository{
		cfg: cfg,
	}
}

func (s *GameRepository) Create(game *entity.GameInfo) (*entity.GameInfo, error) {
	// Check ID
	if game.ID == "" {
		return nil, errors.BadName.AddMessage(game.Name.String())
	}

	// Check if such an object already exists
	if val, _ := s.GetByID(game.ID); val != nil {
		return nil, errors.GameExist
	}

	// Create folder
	if err := fs.CreateFolder(game.Path(s.cfg)); err != nil {
		return nil, err
	}

	// Quote values before write to file
	game.SetQuotedOutput()
	defer game.SetRawOutput()

	// Writing info to file
	if err := fs.CreateAndProcess(game.InfoPath(s.cfg), game, fs.JsonToWriter[*entity.GameInfo]); err != nil {
		return nil, err
	}

	if game.Image == "" {
		return game, nil
	}

	// Download image
	if err := s.CreateImage(game.ID, game.Image); err != nil {
		return nil, err
	}

	game.CachedImage = fmt.Sprintf(s.cfg.GameImagePath, game.ID)

	// Writing info to file
	if err := fs.CreateAndProcess(game.InfoPath(s.cfg), game, fs.JsonToWriter[*entity.GameInfo]); err != nil {
		return nil, err
	}

	return game, nil
}
func (s *GameRepository) GetByID(gameID string) (*entity.GameInfo, error) {
	game := entity.GameInfo{ID: gameID}

	// Check if such an object exists
	isExist, err := fs.IsFolderExist(game.Path(s.cfg))
	if err != nil {
		return nil, err
	}
	if !isExist {
		return nil, errors.GameNotExists
	}

	// Check if such an object exists
	isExist, err = fs.IsFileExist(game.InfoPath(s.cfg))
	if err != nil {
		return nil, err
	}
	if !isExist {
		return nil, errors.GameInfoNotExists
	}

	// Read info from file
	retGame, err := fs.OpenAndProcess(game.InfoPath(s.cfg), fs.JsonFromReader[entity.GameInfo])
	if err != nil {
		return nil, err
	}

	return retGame, nil
}
func (s *GameRepository) GetAll() ([]*entity.GameInfo, error) {
	isExist, err := fs.IsFolderExist(s.cfg.Games())
	if err != nil {
		return make([]*entity.GameInfo, 0), err
	}
	if !isExist {
		return make([]*entity.GameInfo, 0), nil
	}

	// Get list of objects
	folders, err := fs.ListOfFolders(s.cfg.Games())
	if err != nil {
		return make([]*entity.GameInfo, 0), err
	}

	// Get each game
	games := make([]*entity.GameInfo, 0)
	for _, gameID := range folders {
		game, err := s.GetByID(gameID)
		if err != nil {
			logger.Error.Println(err.Error())
			continue
		}
		games = append(games, game)
	}

	return games, nil
}
func (s *GameRepository) Update(gameID string, dtoObject *dto.UpdateGameDTO) (*entity.GameInfo, error) {
	// Get old object
	oldGame, err := s.GetByID(gameID)
	if err != nil {
		return nil, err
	}

	// Create game object
	if dtoObject.Name == "" {
		dtoObject.Name = oldGame.Name.String()
	}
	game := entity.NewGameInfo(dtoObject.Name, dtoObject.Description, dtoObject.Image)
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
		err = fs.MoveFolder(oldGame.Path(s.cfg), game.Path(s.cfg))
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
		if err = fs.CreateAndProcess(game.InfoPath(s.cfg), game, fs.JsonToWriter[*entity.GameInfo]); err != nil {
			return nil, err
		}
	}

	// If the image has been changed
	if game.Image == oldGame.Image {
		return game, nil
	}

	// If image exist, delete
	if data, _, _ := s.GetImage(game.ID); data != nil {
		err = fs.RemoveFile(game.ImagePath(s.cfg))
		if err != nil {
			return nil, err
		}
	}

	if game.Image == "" {
		return game, nil
	}

	// Download image
	if err = s.CreateImage(game.ID, game.Image); err != nil {
		return nil, err
	}

	game.CachedImage = fmt.Sprintf(s.cfg.GameImagePath, game.ID)

	// Writing info to file
	if err = fs.CreateAndProcess(game.InfoPath(s.cfg), game, fs.JsonToWriter[*entity.GameInfo]); err != nil {
		return nil, err
	}

	return game, nil
}
func (s *GameRepository) DeleteByID(gameID string) error {
	game := entity.GameInfo{ID: gameID}

	// Check if such an object exists
	if val, _ := s.GetByID(gameID); val == nil {
		return errors.GameNotExists.HTTP(http.StatusBadRequest)
	}

	// Remove object
	return fs.RemoveFolder(game.Path(s.cfg))
}
func (s *GameRepository) GetImage(gameID string) ([]byte, string, error) {
	// Check if such an object exists
	game, err := s.GetByID(gameID)
	if err != nil {
		return nil, "", err
	}

	// Check if an image exists
	isExist, err := fs.IsFileExist(game.ImagePath(s.cfg))
	if err != nil {
		return nil, "", err
	}
	if !isExist {
		return nil, "", errors.GameImageNotExists
	}

	// Read an image from a file
	data, err := fs.OpenAndProcess(game.ImagePath(s.cfg), fs.BinFromReader)
	if err != nil {
		return nil, "", err
	}

	imgType, err := images.ValidateImage(data)
	if err != nil {
		return nil, "", err
	}

	return data, imgType, nil
}
func (s *GameRepository) CreateImage(gameID, imageURL string) error {
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
	return fs.CreateAndProcess(game.ImagePath(s.cfg), imageBytes, fs.BinToWriter)
}
func (s *GameRepository) Duplicate(gameID string, dtoObject *dto.DuplicateGameDTO) (*entity.GameInfo, error) {
	// Check if the game exists
	oldGame, _ := s.GetByID(gameID)
	if oldGame == nil {
		return nil, errors.GameNotExists.HTTP(http.StatusBadRequest)
	}

	// New game object
	game := entity.NewGameInfo(dtoObject.Name, oldGame.Description.String(), oldGame.Image)

	// Check ID
	if game.ID == "" {
		return nil, errors.BadName.AddMessage(game.Name.String())
	}

	// Check if such an object already exists
	if val, _ := s.GetByID(game.ID); val != nil {
		return nil, errors.GameExist
	}

	// Create a copy of the game
	err := fs.CopyFolder(oldGame.Path(s.cfg), game.Path(s.cfg))
	if err != nil {
		return nil, err
	}

	// Quote values before write to file
	game.SetQuotedOutput()
	defer game.SetRawOutput()

	// Writing info to file
	if err = fs.CreateAndProcess(game.InfoPath(s.cfg), game, fs.JsonToWriter[*entity.GameInfo]); err != nil {
		return nil, err
	}

	return game, nil
}
func (s *GameRepository) Export(gameID string) ([]byte, error) {
	// Check if such an object exists
	game, _ := s.GetByID(gameID)
	if game == nil {
		return nil, errors.GameNotExists.HTTP(http.StatusBadRequest)
	}

	return fs.ArchiveFolder(game.Path(s.cfg), game.ID)
}
func (s *GameRepository) Import(data []byte, name string) error {
	gameID := utils.NameToID(name)
	if name != "" && gameID == "" {
		return errors.BadName
	}

	// Unpack the archive
	resultGameID, err := fs.UnarchiveFolder(data, gameID, s.cfg)
	if err != nil {
		return err
	}

	// Check if the root folder contains information about the game
	game, err := s.GetByID(resultGameID)
	if err != nil {
		// Build a full relative path to the root game folder
		gameRootPath := filepath.Join(s.cfg.Games(), resultGameID)
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
		if err = fs.CreateAndProcess(game.InfoPath(s.cfg), game, fs.JsonToWriter[*entity.GameInfo]); err != nil {
			return err
		}
	}

	return nil
}
