package game

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

type game struct {
	cfg  *config.Config
	game dbGame.Game
}

func New(cfg *config.Config, g dbGame.Game) Game {
	return &game{
		cfg:  cfg,
		game: g,
	}
}

func (r *game) Create(dtoObject *dto.CreateGameDTO) (*entity.GameInfo, error) {
	game, err := r.game.Create(context.Background(), dtoObject.Name, dtoObject.Description, dtoObject.Image)
	if err != nil {
		return nil, err
	}

	if game.Image == "" && dtoObject.ImageFile == nil {
		return game, nil
	}

	if game.Image != "" {
		// Download image
		err = r.createImage(game.ID, game.Image)
		if err != nil {
			logger.Warn.Println("Unable to load image. The game will be saved without an image.", err.Error())
		}
	} else if dtoObject.ImageFile != nil {
		err = r.createImageFromByte(game.ID, dtoObject.ImageFile)
		if err != nil {
			logger.Warn.Println("Invalid image. The game will be saved without an image.", err.Error())
		}
	}

	return game, nil
}
func (r *game) GetByID(gameID string) (*entity.GameInfo, error) {
	_, resp, err := r.game.Get(context.Background(), gameID)
	return resp, err
}
func (r *game) GetAll() ([]*entity.GameInfo, error) {
	return r.game.List(context.Background())
}
func (r *game) Update(gameID string, dtoObject *dto.UpdateGameDTO) (*entity.GameInfo, error) {
	_, oldGame, err := r.game.Get(context.Background(), gameID)
	if err != nil {
		return nil, err
	}

	var newGame *entity.GameInfo
	if oldGame.Name != dtoObject.Name {
		// Rename folder
		newGame, err = r.game.Move(context.Background(), oldGame.Name, dtoObject.Name)
		if err != nil {
			return nil, err
		}
	}

	if oldGame.Description != dtoObject.Description ||
		oldGame.Image != dtoObject.Image ||
		dtoObject.ImageFile != nil {
		// Update data
		newGame, err = r.game.Update(context.Background(), dtoObject.Name, dtoObject.Description, dtoObject.Image)
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
	if data, _, _ := r.GetImage(newGame.ID); data != nil {
		err = r.game.ImageDelete(context.Background(), newGame.ID)
		if err != nil {
			return nil, err
		}
	}

	if newGame.Image == "" && dtoObject.ImageFile == nil {
		return newGame, nil
	}

	if newGame.Image != "" {
		// Download image
		err = r.createImage(newGame.ID, newGame.Image)
		if err != nil {
			logger.Warn.Println("Unable to load image. The game will be saved without an image.", err.Error())
		}
	} else if dtoObject.ImageFile != nil {
		err = r.createImageFromByte(newGame.ID, dtoObject.ImageFile)
		if err != nil {
			logger.Warn.Println("Invalid image. The game will be saved without an image.", err.Error())
		}
	}

	return newGame, nil
}
func (r *game) DeleteByID(gameID string) error {
	return r.game.Delete(context.Background(), gameID)
}
func (r *game) GetImage(gameID string) ([]byte, string, error) {
	data, err := r.game.ImageGet(context.Background(), gameID)
	if err != nil {
		return nil, "", err
	}

	imgType, err := images.ValidateImage(data)
	if err != nil {
		return nil, "", err
	}

	return data, imgType, nil
}
func (r *game) Duplicate(gameID string, dtoObject *dto.DuplicateGameDTO) (*entity.GameInfo, error) {
	game, err := r.game.Duplicate(context.Background(), gameID, dtoObject.Name)
	if err != nil {
		return nil, err
	}
	return game, nil
}
func (r *game) Export(gameID string) ([]byte, error) {
	// Check if such an object exists
	game, err := r.GetByID(gameID)
	if err != nil {
		return nil, err
	}

	return fs.ArchiveFolder(game.Path(r.cfg), game.ID)
}
func (r *game) Import(data []byte, name string) (*entity.GameInfo, error) {
	gameID := utils.NameToID(name)
	if name != "" && gameID == "" {
		return nil, errors.BadName
	}

	// Unpack the archive
	resultGameID, err := fs.UnarchiveFolder(data, gameID, r.cfg)
	if err != nil {
		return nil, err
	}

	// Check if the root folder contains information about the game
	game, err := r.GetByID(resultGameID)
	if err != nil {
		// If an error occurs during unzipping, delete the created folder with the game
		errors.IfErrorLog(r.game.Delete(context.Background(), resultGameID))
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

		if err = r.game.UpdateInfo(context.Background(), game.ID, name); err != nil {
			return nil, err
		}
	}

	return game, nil
}

func (r *game) createImage(gameID, imageURL string) error {
	// Download image
	imageBytes, err := network.DownloadBytes(imageURL)
	if err != nil {
		return err
	}

	return r.createImageFromByte(gameID, imageBytes)
}
func (r *game) createImageFromByte(gameID string, data []byte) error {
	// Validate image
	_, err := images.ValidateImage(data)
	if err != nil {
		return err
	}

	// Write image to file
	return r.game.ImageCreate(context.Background(), gameID, data)
}