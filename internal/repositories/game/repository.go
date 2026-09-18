package game

import (
	"errors"
	"net/http"
	"path/filepath"
	"time"

	"github.com/HardDie/fsentry"

	"github.com/HardDie/DeckBuilder/internal/config"
	entitiesGame "github.com/HardDie/DeckBuilder/internal/entities/game"
	er "github.com/HardDie/DeckBuilder/internal/errors"
	"github.com/HardDie/DeckBuilder/internal/fs"
	"github.com/HardDie/DeckBuilder/internal/images"
	"github.com/HardDie/DeckBuilder/internal/logger"
	"github.com/HardDie/DeckBuilder/internal/network"
	"github.com/HardDie/DeckBuilder/internal/utils"
)

type game struct {
	cfg       *config.Config
	db        *fsentry.DB
	gamesPath string
}

func New(cfg *config.Config, db *fsentry.DB) Game {
	return &game{
		cfg:       cfg,
		db:        db,
		gamesPath: "games",
	}
}

func (r *game) Create(req CreateRequest) (*entitiesGame.Game, error) {
	g, err := r.create(req)
	if err != nil {
		return nil, err
	}

	if g.Image == "" && req.ImageFile == nil {
		return g, nil
	}

	if g.Image != "" {
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
	return r.get(gameID)
}

func (r *game) GetAll() ([]*entitiesGame.Game, error) {
	return r.list()
}

func (r *game) Update(gameID string, req UpdateRequest) (*entitiesGame.Game, error) {
	oldGame, err := r.get(gameID)
	if err != nil {
		return nil, err
	}

	var newGame *entitiesGame.Game
	if oldGame.Name != req.Name {
		newGame, err = r.move(oldGame.Name, req.Name)
		if err != nil {
			return nil, err
		}
	}

	if oldGame.Description != req.Description ||
		oldGame.Image != req.Image ||
		req.ImageFile != nil {
		newGame, err = r.update(updateRequest{
			Name:        req.Name,
			Description: req.Description,
			Image:       req.Image,
		})
		if err != nil {
			return nil, err
		}
	}

	if newGame == nil {
		newGame = oldGame
	}

	if newGame.Image == oldGame.Image && req.ImageFile == nil {
		return newGame, nil
	}

	if data, _, _ := r.GetImage(newGame.ID); data != nil {
		err = r.imageDelete(newGame.ID)
		if err != nil {
			return nil, err
		}
	}

	if newGame.Image == "" && req.ImageFile == nil {
		return newGame, nil
	}

	if newGame.Image != "" {
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
	return r.delete(gameID)
}

func (r *game) GetImage(gameID string) ([]byte, string, error) {
	data, err := r.imageGet(gameID)
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
	return r.duplicate(gameID, req.Name)
}

func (r *game) Export(gameID string) ([]byte, error) {
	g, err := r.GetByID(gameID)
	if err != nil {
		return nil, err
	}

	return fs.ArchiveFolder(filepath.Join(r.cfg.Games(), g.ID), g.ID)
}

func (r *game) Import(data []byte, name string) (*entitiesGame.Game, error) {
	gameID := utils.NameToID(name)
	if name != "" && gameID == "" {
		return nil, er.BadName
	}

	resultGameID, err := fs.UnarchiveFolder(data, gameID, r.cfg)
	if err != nil {
		return nil, err
	}

	g, err := r.GetByID(resultGameID)
	if err != nil {
		er.IfErrorLog(r.delete(resultGameID))
		return nil, err
	}

	if name == "" && resultGameID != g.ID {
		gameID = resultGameID
		name = resultGameID
	}

	if name != "" {
		g.ID = gameID
		g.Name = name

		if err = r.updateInfo(g.ID, name); err != nil {
			return nil, err
		}
	}

	return g, nil
}

func (r *game) createImage(gameID, imageURL string) error {
	imageBytes, err := network.DownloadBytes(imageURL)
	if err != nil {
		return err
	}

	return r.createImageFromByte(gameID, imageBytes)
}

func (r *game) createImageFromByte(gameID string, data []byte) error {
	_, err := images.ValidateImage(data)
	if err != nil {
		return err
	}

	return r.imageCreate(gameID, data)
}

func (r *game) create(req CreateRequest) (*entitiesGame.Game, error) {
	info, err := r.db.CreateFolder(req.Name, model{
		Description: fsentry.QuotedString(req.Description),
		Image:       fsentry.QuotedString(req.Image),
	}, r.gamesPath)
	if err != nil {
		if errors.Is(err, fsentry.ErrExist) {
			return nil, er.GameExist
		} else if errors.Is(err, fsentry.ErrBadName) {
			return nil, er.BadName
		} else {
			return nil, er.InternalError.AddMessage(err.Error())
		}
	}

	return r.toEntity(info), nil
}

func (r *game) get(name string) (*entitiesGame.Game, error) {
	info, err := r.db.GetFolder[model](name, r.gamesPath)
	if err != nil {
		if errors.Is(err, fsentry.ErrNotExist) {
			return nil, er.GameNotExists.AddMessage(err.Error()).HTTP(http.StatusBadRequest)
		} else if errors.Is(err, fsentry.ErrBadName) {
			return nil, er.BadName
		} else {
			return nil, er.InternalError.AddMessage(err.Error())
		}
	}

	return r.toEntity(info), nil
}

func (r *game) list() ([]*entitiesGame.Game, error) {
	list, err := r.db.List(r.gamesPath)
	if err != nil {
		return nil, er.InternalError.AddMessage(err.Error())
	}

	var games []*entitiesGame.Game
	for _, folder := range list.Folders {
		game, err := r.get(folder)
		if err != nil {
			logger.Error.Println(folder, err.Error())
			continue
		}
		if folder != game.ID {
			logger.Error.Println("Corrupted game folder:", folder)
			continue
		}
		games = append(games, game)
	}
	return games, nil
}

func (r *game) move(oldName, newName string) (*entitiesGame.Game, error) {
	info, err := r.db.MoveFolder[model](oldName, newName, r.gamesPath)
	if err != nil {
		if errors.Is(err, fsentry.ErrNotExist) {
			return nil, er.GameNotExists.AddMessage(err.Error()).HTTP(http.StatusBadRequest)
		} else if errors.Is(err, fsentry.ErrBadName) {
			return nil, er.BadName
		} else {
			return nil, er.InternalError.AddMessage(err.Error())
		}
	}

	return r.toEntity(info), nil
}

type updateRequest struct {
	Name        string
	Description string
	Image       string
}

func (r *game) update(req updateRequest) (*entitiesGame.Game, error) {
	info, err := r.db.UpdateFolder(req.Name, model{
		Description: fsentry.QuotedString(req.Description),
		Image:       fsentry.QuotedString(req.Image),
	}, r.gamesPath)
	if err != nil {
		if errors.Is(err, fsentry.ErrNotExist) {
			return nil, er.GameNotExists.AddMessage(err.Error()).HTTP(http.StatusBadRequest)
		} else if errors.Is(err, fsentry.ErrBadName) {
			return nil, er.BadName
		} else {
			return nil, er.InternalError.AddMessage(err.Error())
		}
	}

	return r.toEntity(info), nil
}

func (r *game) delete(name string) error {
	err := r.db.RemoveFolder(name, r.gamesPath)
	if err != nil {
		if errors.Is(err, fsentry.ErrNotExist) {
			return er.GameNotExists.AddMessage(err.Error()).HTTP(http.StatusBadRequest)
		} else if errors.Is(err, fsentry.ErrBadName) {
			return er.BadName
		} else {
			return er.InternalError.AddMessage(err.Error())
		}
	}
	return nil
}

func (r *game) duplicate(srcName, dstName string) (*entitiesGame.Game, error) {
	info, err := r.db.DuplicateFolder[model](srcName, dstName, r.gamesPath)
	if err != nil {
		if errors.Is(err, fsentry.ErrNotExist) {
			return nil, er.GameNotExists.AddMessage(err.Error()).HTTP(http.StatusBadRequest)
		} else if errors.Is(err, fsentry.ErrExist) {
			return nil, er.GameExist.AddMessage(err.Error()).HTTP(http.StatusBadRequest)
		} else if errors.Is(err, fsentry.ErrBadName) {
			return nil, er.BadName
		} else {
			return nil, er.InternalError.AddMessage(err.Error())
		}
	}
	return r.toEntity(info), nil
}

func (r *game) updateInfo(name, newName string) error {
	_, err := r.db.UpdateFolderNameWithoutTimestamp[model](name, newName, r.gamesPath)
	return err
}

func (r *game) imageCreate(gameID string, data []byte) error {
	game, err := r.get(gameID)
	if err != nil {
		return err
	}

	err = r.db.CreateBinary("image", data, r.gamesPath, game.ID)
	if err != nil {
		if errors.Is(err, fsentry.ErrExist) {
			return er.GameImageExist.AddMessage(err.Error())
		} else {
			return er.InternalError.AddMessage(err.Error())
		}
	}
	return nil
}

func (r *game) imageGet(gameID string) ([]byte, error) {
	game, err := r.get(gameID)
	if err != nil {
		return nil, err
	}

	data, err := r.db.GetBinary("image", nil, r.gamesPath, game.ID)
	if err != nil {
		if errors.Is(err, fsentry.ErrNotExist) {
			return nil, er.GameImageNotExists.AddMessage(err.Error())
		} else {
			return nil, er.InternalError.AddMessage(err.Error())
		}
	}
	return data, nil
}

func (r *game) imageDelete(gameID string) error {
	game, err := r.get(gameID)
	if err != nil {
		return err
	}

	err = r.db.RemoveBinary("image", r.gamesPath, game.ID)
	if err != nil {
		if errors.Is(err, fsentry.ErrNotExist) {
			return er.GameImageNotExists.AddMessage(err.Error())
		} else {
			return er.InternalError.AddMessage(err.Error())
		}
	}
	return nil
}

func (r *game) toEntity(info fsentry.FolderInfo[model]) *entitiesGame.Game {
	createdAt, updatedAt := r.convertCreateUpdate(info.CreatedAt, info.UpdatedAt)
	return &entitiesGame.Game{
		ID:          info.ID,
		Name:        info.Name,
		Description: info.Data.Description.String(),
		Image:       info.Data.Image.String(),
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
	}
}

func (r *game) convertCreateUpdate(createdAt, updatedAt time.Time) (time.Time, time.Time) {
	if createdAt.IsZero() {
		createdAt = time.Now()
	}
	if updatedAt.IsZero() {
		updatedAt = createdAt
	}
	return createdAt, updatedAt
}
