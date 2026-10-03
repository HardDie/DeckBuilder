package game

import (
	"bytes"
	"errors"

	"github.com/HardDie/fsentry"

	"github.com/HardDie/DeckBuilder/internal/config"
	entitiesGame "github.com/HardDie/DeckBuilder/internal/entities/game"
	er "github.com/HardDie/DeckBuilder/internal/errors"
	"github.com/HardDie/DeckBuilder/internal/images"
	"github.com/HardDie/DeckBuilder/internal/logger"
	"github.com/HardDie/DeckBuilder/internal/repositories"
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
	change, imageErr := repositories.ResolveImage("", req.Image, req.ImageFile)
	req.Image = change.URL

	g, err := r.create(req)
	if err != nil {
		return nil, err
	}
	if change.Data != nil {
		if err = r.imageCreate(g.ID, change.Data); err != nil {
			imageErr = err
		}
	}
	g.ImageError = imageErr
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
	// Download and validate before any write, so a bad image keeps the old one.
	change, imageErr := repositories.ResolveImage(oldGame.Image, req.Image, req.ImageFile)

	newGame := oldGame
	if oldGame.Name != req.Name {
		newGame, err = r.move(oldGame.Name, req.Name)
		if err != nil {
			return nil, err
		}
	}

	if oldGame.Description != req.Description ||
		oldGame.Image != change.URL ||
		change.Data != nil {
		newGame, err = r.update(updateRequest{
			Name:        req.Name,
			Description: req.Description,
			Image:       change.URL,
		})
		if err != nil {
			return nil, err
		}
	}

	if change.Data != nil || change.Clear {
		err = r.imageDelete(newGame.ID)
		if err != nil && !errors.Is(err, er.GameImageNotExists) {
			return nil, err
		}
	}
	if change.Data != nil {
		if err = r.imageCreate(newGame.ID, change.Data); err != nil {
			imageErr = err
		}
	}
	newGame.ImageError = imageErr
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

	imgType, err := images.ImageType(data)
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

	var buf bytes.Buffer
	if err = r.db.ExportFolder(&buf, g.ID, r.gamesPath); err != nil {
		return nil, er.InternalError.AddMessage(err.Error())
	}
	return buf.Bytes(), nil
}

func (r *game) Import(data []byte, name string) (*entitiesGame.Game, error) {
	id, err := r.db.ImportFolder(bytes.NewReader(data), name, r.gamesPath)
	if err != nil {
		if errors.Is(err, fsentry.ErrBadArchive) {
			return nil, er.BadArchive.AddMessage(err.Error())
		}
		return nil, repositories.MapFsentry(err, repositories.FsentrySentinels{
			Exist:   er.GameExist,
			Message: true,
		})
	}

	g, err := r.GetByID(id)
	if err != nil {
		er.IfErrorLog(r.delete(id))
		return nil, err
	}
	return g, nil
}

func (r *game) create(req CreateRequest) (*entitiesGame.Game, error) {
	info, err := r.db.CreateFolder(req.Name, model{
		Description: fsentry.QuotedString(req.Description),
		Image:       fsentry.QuotedString(req.Image),
	}, r.gamesPath)
	if err != nil {
		return nil, repositories.MapFsentry(err, repositories.FsentrySentinels{
			Exist: er.GameExist,
		})
	}

	return r.toEntity(info), nil
}

func (r *game) get(name string) (*entitiesGame.Game, error) {
	info, err := r.db.GetFolder[model](name, r.gamesPath)
	if err != nil {
		return nil, repositories.MapFsentry(err, repositories.FsentrySentinels{
			NotExist: er.GameNotExists,
			Message:  true,
		})
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
		return nil, repositories.MapFsentry(err, repositories.FsentrySentinels{
			NotExist: er.GameNotExists,
			Message:  true,
		})
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
		return nil, repositories.MapFsentry(err, repositories.FsentrySentinels{
			NotExist: er.GameNotExists,
			Message:  true,
		})
	}

	return r.toEntity(info), nil
}

func (r *game) delete(name string) error {
	err := r.db.RemoveFolder(name, r.gamesPath)
	if err != nil {
		return repositories.MapFsentry(err, repositories.FsentrySentinels{
			NotExist: er.GameNotExists,
			Message:  true,
		})
	}
	return nil
}

func (r *game) duplicate(srcName, dstName string) (*entitiesGame.Game, error) {
	info, err := r.db.DuplicateFolder[model](srcName, dstName, r.gamesPath)
	if err != nil {
		return nil, repositories.MapFsentry(err, repositories.FsentrySentinels{
			Exist:    er.GameExist,
			NotExist: er.GameNotExists,
			Message:  true,
		})
	}
	return r.toEntity(info), nil
}

func (r *game) imageCreate(gameID string, data []byte) error {
	game, err := r.get(gameID)
	if err != nil {
		return err
	}

	err = r.db.CreateBinary("image", data, r.gamesPath, game.ID)
	if err != nil {
		return repositories.MapFsentry(err, repositories.FsentrySentinels{
			Exist:   er.GameImageExist,
			Message: true,
		})
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
		return nil, repositories.MapFsentry(err, repositories.FsentrySentinels{
			NotExist: er.GameImageNotExists,
			Message:  true,
		})
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
		return repositories.MapFsentry(err, repositories.FsentrySentinels{
			NotExist: er.GameImageNotExists,
			Message:  true,
		})
	}
	return nil
}

func (r *game) toEntity(info fsentry.FolderInfo[model]) *entitiesGame.Game {
	createdAt, updatedAt := utils.NormalizeTimestamps(info.CreatedAt, info.UpdatedAt)
	return &entitiesGame.Game{
		ID:          info.ID,
		Name:        info.Name,
		Description: info.Data.Description.String(),
		Image:       info.Data.Image.String(),
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
	}
}
