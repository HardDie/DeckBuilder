package repository

import (
	"context"

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
	game, err := s.db.GameCreate(context.Background(), req.Name, req.Description, req.Image)
	if err != nil {
		return nil, err
	}

	if game.Image == "" {
		return game, nil
	}

	// Download image
	if err := s.createImage(game.ID, game.Image); err != nil {
		logger.Warn.Println("Unable to load image. The game will be saved without an image.", err.Error())
	}

	return game, nil
}
func (s *GameRepository) GetByID(gameID string) (*entity.GameInfo, error) {
	_, resp, err := s.db.GameGet(context.Background(), gameID)
	return resp, err
}
func (s *GameRepository) GetAll() ([]*entity.GameInfo, error) {
	return s.db.GameList(context.Background())
}
func (s *GameRepository) Update(gameID string, dtoObject *dto.UpdateGameDTO) (*entity.GameInfo, error) {
	_, oldGame, err := s.db.GameGet(context.Background(), gameID)
	if err != nil {
		return nil, err
	}

	var newGame *entity.GameInfo
	if oldGame.Name != dtoObject.Name {
		// Rename folder
		newGame, err = s.db.GameMove(context.Background(), oldGame.Name, dtoObject.Name)
		if err != nil {
			return nil, err
		}
	}

	if oldGame.Description != dtoObject.Description ||
		oldGame.Image != dtoObject.Image {
		// Update data
		newGame, err = s.db.GameUpdate(context.Background(), dtoObject.Name, dtoObject.Description, dtoObject.Image)
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
		err = s.db.GameImageDelete(context.Background(), newGame.ID)
		if err != nil {
			return nil, err
		}
	}

	if newGame.Image == "" {
		return newGame, nil
	}

	// Download image
	if err = s.createImage(newGame.ID, newGame.Image); err != nil {
		logger.Warn.Println("Unable to load image. The game will be saved without an image.", err.Error())
	}

	return newGame, nil
}
func (s *GameRepository) DeleteByID(gameID string) error {
	return s.db.GameDelete(context.Background(), gameID)
}
func (s *GameRepository) GetImage(gameID string) ([]byte, string, error) {
	data, err := s.db.GameImageGet(context.Background(), gameID)
	if err != nil {
		return nil, "", err
	}

	imgType, err := images.ValidateImage(data)
	if err != nil {
		return nil, "", err
	}

	return data, imgType, nil
}
func (s *GameRepository) Duplicate(gameID string, dtoObject *dto.DuplicateGameDTO) (*entity.GameInfo, error) {
	game, err := s.db.GameDuplicate(context.Background(), gameID, dtoObject.Name)
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
		// If an error occurs during unzipping, delete the created folder with the game
		errors.IfErrorLog(s.db.GameDelete(context.Background(), resultGameID))
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

		if err = s.db.GameUpdateInfo(context.Background(), game.ID, name); err != nil {
			return nil, err
		}
	}

	return game, nil
}

func (s *GameRepository) createImage(gameID, imageURL string) error {
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
	return s.db.GameImageCreate(context.Background(), gameID, imageBytes)
}
