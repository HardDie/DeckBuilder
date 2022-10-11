package repository

import (
	"net/http"
	"path/filepath"

	"github.com/HardDie/DeckBuilder/internal/config"
	"github.com/HardDie/DeckBuilder/internal/db"
	"github.com/HardDie/DeckBuilder/internal/dto"
	"github.com/HardDie/DeckBuilder/internal/entity"
	"github.com/HardDie/DeckBuilder/internal/errors"
	"github.com/HardDie/DeckBuilder/internal/fs"
	"github.com/HardDie/DeckBuilder/internal/images"
	"github.com/HardDie/DeckBuilder/internal/logger"
	"github.com/HardDie/DeckBuilder/internal/network"
	"github.com/HardDie/DeckBuilder/internal/utils"
)

type IGameRepository interface {
	Create(req *dto.CreateGameDTO) (*entity.GameInfo, error)
	GetByID(gameID string) (*entity.GameInfo, error)
	GetAll() ([]*entity.GameInfo, error)
	Update(gameID string, dtoObject *dto.UpdateGameDTO) (*entity.GameInfo, error)
	DeleteByID(gameID string) error
	GetImage(gameID string) ([]byte, string, error)
	CreateImage(gameID, imageURL string) error
	Duplicate(gameID string, dtoObject *dto.DuplicateGameDTO) (*entity.GameInfo, error)
	Export(gameID string) ([]byte, error)
	Import(data []byte, name string) (*entity.GameInfo, error)
}
type GameRepository struct {
	cfg *config.Config
	db  *db.DB
}

func NewGameRepository(cfg *config.Config, db *db.DB) *GameRepository {
	return &GameRepository{
		cfg: cfg,
		db:  db,
	}
}

func (s *GameRepository) Create(req *dto.CreateGameDTO) (*entity.GameInfo, error) {
	game, err := s.db.GameCreate(req.Name, req.Description, req.Image)
	if err != nil {
		return nil, err
	}

	if game.Image == "" {
		return game, nil
	}

	// Download image
	if err := s.CreateImage(game.ID, game.Image); err != nil {
		logger.Warn.Println("Unable to load image. The game will be saved without an image.", err.Error())
	}

	return game, nil
}
func (s *GameRepository) GetByID(gameID string) (*entity.GameInfo, error) {
	return s.db.GameGet(gameID)
}
func (s *GameRepository) GetAll() ([]*entity.GameInfo, error) {
	return s.db.GameList()
}
func (s *GameRepository) Update(gameID string, dtoObject *dto.UpdateGameDTO) (*entity.GameInfo, error) {
	oldGame, err := s.db.GameGet(gameID)
	if err != nil {
		return nil, err
	}

	var newGame *entity.GameInfo
	if oldGame.Name != dtoObject.Name {
		// Rename folder
		newGame, err = s.db.GameMove(oldGame.Name, dtoObject.Name)
		if err != nil {
			return nil, err
		}
	}

	if oldGame.Description != dtoObject.Description ||
		oldGame.Image != dtoObject.Image {
		// Update data
		newGame, err = s.db.GameUpdate(dtoObject.Name, dtoObject.Description, dtoObject.Image)
		if err != nil {
			return nil, err
		}
	}

	if newGame == nil {
		// If nothing has changed
		newGame = oldGame
	}

	// If the image has not been changed
	if newGame.Image == oldGame.Image {
		return newGame, nil
	}

	// If image exist, delete
	if data, _, _ := s.GetImage(newGame.ID); data != nil {
		err = fs.RemoveFile(newGame.ImagePath(s.cfg))
		if err != nil {
			return nil, err
		}
	}

	if newGame.Image == "" {
		return newGame, nil
	}

	// Download image
	if err = s.CreateImage(newGame.ID, newGame.Image); err != nil {
		logger.Warn.Println("Unable to load image. The game will be saved without an image.", err.Error())
	}

	return newGame, nil
}
func (s *GameRepository) DeleteByID(gameID string) error {
	return s.db.GameDelete(gameID)
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
	game, err := s.db.GameDuplicate(gameID, dtoObject.Name)
	if err != nil {
		return nil, err
	}
	return game, nil
}
func (s *GameRepository) Export(gameID string) ([]byte, error) {
	// Check if such an object exists
	game, err := s.GetByID(gameID)
	if err != nil {
		return nil, err
	}

	return fs.ArchiveFolder(game.Path(s.cfg), game.ID)
}
func (s *GameRepository) Import(data []byte, name string) (*entity.GameInfo, error) {
	gameID := utils.NameToID(name)
	if name != "" && gameID == "" {
		return nil, errors.BadName
	}

	// Unpack the archive
	resultGameID, err := fs.UnarchiveFolder(data, gameID, s.cfg)
	if err != nil {
		return nil, err
	}

	// Check if the root folder contains information about the game
	game, err := s.GetByID(resultGameID)
	if err != nil {
		// Build a full relative path to the root game folder
		gameRootPath := filepath.Join(s.cfg.Games(), resultGameID)
		// If an error occurs during unzipping, delete the created folder with the game
		errors.IfErrorLog(fs.RemoveFolder(gameRootPath))
		return nil, err
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
		game.Name = name

		if err = s.db.GameUpdateInfo(game.ID, name); err != nil {
			return nil, err
		}
	}

	return game, nil
}
