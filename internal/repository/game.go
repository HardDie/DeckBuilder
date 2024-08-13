package repository

import (
	"context"

	"github.com/HardDie/DeckBuilder/internal/config"
	dbGame "github.com/HardDie/DeckBuilder/internal/db/game"
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
	cfg  *config.Config
	game dbGame.Game
}

func NewGameRepository(cfg *config.Config, game dbGame.Game) *GameRepository {
	return &GameRepository{
		cfg:  cfg,
		game: game,
	}
}

func (s *GameRepository) Create(dtoObject *dto.CreateGameDTO) (*entity.GameInfo, error) {
	game, err := s.game.Create(context.Background(), dtoObject.Name, dtoObject.Description, dtoObject.Image)
	if err != nil {
		return nil, err
	}

	if game.Image == "" && dtoObject.ImageFile == nil {
		return game, nil
	}

	if game.Image != "" {
		// Download image
		err = s.createImage(game.ID, game.Image)
		if err != nil {
			logger.Warn.Println("Unable to load image. The game will be saved without an image.", err.Error())
		}
	} else if dtoObject.ImageFile != nil {
		err = s.createImageFromByte(game.ID, dtoObject.ImageFile)
		if err != nil {
			logger.Warn.Println("Invalid image. The game will be saved without an image.", err.Error())
		}
	}

	return game, nil
}
func (s *GameRepository) GetByID(gameID string) (*entity.GameInfo, error) {
	_, resp, err := s.game.Get(context.Background(), gameID)
	return resp, err
}
func (s *GameRepository) GetAll() ([]*entity.GameInfo, error) {
	return s.game.List(context.Background())
}
func (s *GameRepository) Update(gameID string, dtoObject *dto.UpdateGameDTO) (*entity.GameInfo, error) {
	_, oldGame, err := s.game.Get(context.Background(), gameID)
	if err != nil {
		return nil, err
	}

	var newGame *entity.GameInfo
	if oldGame.Name != dtoObject.Name {
		// Rename folder
		newGame, err = s.game.Move(context.Background(), oldGame.Name, dtoObject.Name)
		if err != nil {
			return nil, err
		}
	}

	if oldGame.Description != dtoObject.Description ||
		oldGame.Image != dtoObject.Image ||
		dtoObject.ImageFile != nil {
		// Update data
		newGame, err = s.game.Update(context.Background(), dtoObject.Name, dtoObject.Description, dtoObject.Image)
		if err != nil {
			return nil, err
		}
	}

	if newGame == nil {
		// If nothing has changed
		newGame = oldGame
	}

	// If the image has not been changed
	if newGame.Image == oldGame.Image && dtoObject.ImageFile == nil {
		return newGame, nil
	}

	// If image exist, delete
	if data, _, _ := s.GetImage(newGame.ID); data != nil {
		err = s.game.ImageDelete(context.Background(), newGame.ID)
		if err != nil {
			return nil, err
		}
	}

	if newGame.Image == "" && dtoObject.ImageFile == nil {
		return newGame, nil
	}

	if newGame.Image != "" {
		// Download image
		err = s.createImage(newGame.ID, newGame.Image)
		if err != nil {
			logger.Warn.Println("Unable to load image. The game will be saved without an image.", err.Error())
		}
	} else if dtoObject.ImageFile != nil {
		err = s.createImageFromByte(newGame.ID, dtoObject.ImageFile)
		if err != nil {
			logger.Warn.Println("Invalid image. The game will be saved without an image.", err.Error())
		}
	}

	return newGame, nil
}
func (s *GameRepository) DeleteByID(gameID string) error {
	return s.game.Delete(context.Background(), gameID)
}
func (s *GameRepository) GetImage(gameID string) ([]byte, string, error) {
	data, err := s.game.ImageGet(context.Background(), gameID)
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
	game, err := s.game.Duplicate(context.Background(), gameID, dtoObject.Name)
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
		errors.IfErrorLog(s.game.Delete(context.Background(), resultGameID))
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

		if err = s.game.UpdateInfo(context.Background(), game.ID, name); err != nil {
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

	return s.createImageFromByte(gameID, imageBytes)
}
func (s *GameRepository) createImageFromByte(gameID string, data []byte) error {
	// Validate image
	_, err := images.ValidateImage(data)
	if err != nil {
		return err
	}

	// Write image to file
	return s.game.ImageCreate(context.Background(), gameID, data)
}
