package game

import (
	"context"
	"path/filepath"

	"github.com/HardDie/DeckBuilder/internal/config"
	dbGame "github.com/HardDie/DeckBuilder/internal/db/game"
	entitiesGame "github.com/HardDie/DeckBuilder/internal/entities/game"
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

func (r *game) Create(req CreateRequest) (*entitiesGame.Game, error) {
	g, err := r.game.Create(context.Background(), dbGame.CreateRequest{
		Name:        req.Name,
		Description: req.Description,
		Image:       req.Image,
	})
	if err != nil {
		return nil, err
	}

	if g.Image == "" && req.ImageFile == nil {
		return g, nil
	}

	if g.Image != "" {
		// Download image
		err = r.createImage(g.ID, g.Image)
		if err != nil {
			logger.Warn.Println("Unable to load image. The game will be saved without an image.", err.Error())
		}
	} else if req.ImageFile != nil {
		err = r.createImageFromByte(g.ID, req.ImageFile)
		if err != nil {
			logger.Warn.Println("Invalid image. The game will be saved without an image.", err.Error())
		}
	}

	return g, nil
}
func (r *game) GetByID(gameID string) (*entitiesGame.Game, error) {
	return r.game.Get(context.Background(), gameID)
}
func (r *game) GetAll() ([]*entitiesGame.Game, error) {
	return r.game.List(context.Background())
}
func (r *game) Update(gameID string, req UpdateRequest) (*entitiesGame.Game, error) {
	oldGame, err := r.game.Get(context.Background(), gameID)
	if err != nil {
		return nil, err
	}

	var newGame *entitiesGame.Game
	if oldGame.Name != req.Name {
		// Rename folder
		newGame, err = r.game.Move(context.Background(), oldGame.Name, req.Name)
		if err != nil {
			return nil, err
		}
	}

	if oldGame.Description != req.Description ||
		oldGame.Image != req.Image ||
		req.ImageFile != nil {
		// Update data
		newGame, err = r.game.Update(context.Background(), dbGame.UpdateRequest{
			Name:        req.Name,
			Description: req.Description,
			Image:       req.Image,
		})
		if err != nil {
			return nil, err
		}
	}

	if newGame == nil {
		// If nothing has changed
		newGame = oldGame
	}

	// If the image has not been changed
	if newGame.Image == oldGame.Image && req.ImageFile == nil {
		return newGame, nil
	}

	// If image exist, delete
	if data, _, _ := r.GetImage(newGame.ID); data != nil {
		err = r.game.ImageDelete(context.Background(), newGame.ID)
		if err != nil {
			return nil, err
		}
	}

	if newGame.Image == "" && req.ImageFile == nil {
		return newGame, nil
	}

	if newGame.Image != "" {
		// Download image
		err = r.createImage(newGame.ID, newGame.Image)
		if err != nil {
			logger.Warn.Println("Unable to load image. The game will be saved without an image.", err.Error())
		}
	} else if req.ImageFile != nil {
		err = r.createImageFromByte(newGame.ID, req.ImageFile)
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
func (r *game) Duplicate(gameID string, req DuplicateRequest) (*entitiesGame.Game, error) {
	g, err := r.game.Duplicate(context.Background(), gameID, req.Name)
	if err != nil {
		return nil, err
	}
	return g, nil
}
func (r *game) Export(gameID string) ([]byte, error) {
	// Check if such an object exists
	g, err := r.GetByID(gameID)
	if err != nil {
		return nil, err
	}

	return fs.ArchiveFolder(filepath.Join(r.cfg.Games(), g.ID), g.ID)
}
func (r *game) Import(data []byte, name string) (*entitiesGame.Game, error) {
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
	g, err := r.GetByID(resultGameID)
	if err != nil {
		// If an error occurs during unzipping, delete the created folder with the game
		errors.IfErrorLog(r.game.Delete(context.Background(), resultGameID))
		return nil, err
	}

	// If the user skipped passing a new name for the game,
	// but the root folder has a different name than in the game information file.
	// Fix the game information file.
	if name == "" && resultGameID != g.ID {
		gameID = resultGameID
		name = resultGameID
	}

	// If the name has been changed
	if name != "" {
		// Update the title of the game
		g.ID = gameID
		g.Name = name

		if err = r.game.UpdateInfo(context.Background(), g.ID, name); err != nil {
			return nil, err
		}
	}

	return g, nil
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
