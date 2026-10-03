package game

import (
	"bytes"
	"errors"
	"fmt"
	"log/slog"

	"github.com/HardDie/fsentry"

	"github.com/HardDie/DeckBuilder/internal/apperr"
	entitiesGame "github.com/HardDie/DeckBuilder/internal/entities/game"
	"github.com/HardDie/DeckBuilder/internal/repositories"
	"github.com/HardDie/DeckBuilder/internal/utils"
)

type game struct {
	db        *fsentry.DB
	gamesPath string
	folder    *repositories.Folder
}

func New(db *fsentry.DB) Game {
	return &game{
		db:        db,
		gamesPath: "games",
		folder: repositories.NewFolder(db, repositories.FolderErrors{
			Exist:         apperr.ErrGameExists,
			NotExist:      apperr.ErrGameNotFound,
			ImageExist:    apperr.ErrGameImageExists,
			ImageNotExist: apperr.ErrGameImageNotFound,
		}),
	}
}

func (r *game) Create(req CreateRequest) (*entitiesGame.Game, error) {
	saved, err := r.folder.Create(nil, repositories.FolderWrite(req))
	if err != nil {
		return nil, err
	}
	return toEntity(saved.Info, saved.ImageError), nil
}

func (r *game) GetByID(gameID string) (*entitiesGame.Game, error) {
	info, err := r.folder.Get(nil, gameID)
	if err != nil {
		return nil, err
	}
	return toEntity(info, nil), nil
}

func (r *game) GetAll() ([]*entitiesGame.Game, error) {
	infos, err := r.folder.List(nil)
	if err != nil {
		return nil, err
	}
	var games []*entitiesGame.Game
	for _, info := range infos {
		games = append(games, toEntity(info, nil))
	}
	return games, nil
}

func (r *game) Update(gameID string, req UpdateRequest) (*entitiesGame.Game, error) {
	saved, err := r.folder.Update(nil, gameID, repositories.FolderWrite(req))
	if err != nil {
		return nil, err
	}
	return toEntity(saved.Info, saved.ImageError), nil
}

func (r *game) DeleteByID(gameID string) error {
	return r.folder.Delete(nil, gameID)
}

func (r *game) GetImage(gameID string) ([]byte, string, error) {
	return r.folder.Image(nil, gameID)
}

func (r *game) Duplicate(gameID string, req DuplicateRequest) (*entitiesGame.Game, error) {
	info, err := r.db.DuplicateFolder[repositories.FolderModel](gameID, req.Name, r.gamesPath)
	if err != nil {
		return nil, repositories.MapFsentry(err, repositories.FsentrySentinels{
			Exist:    apperr.ErrGameExists,
			NotExist: apperr.ErrGameNotFound,
		})
	}
	return toEntity(info, nil), nil
}

func (r *game) Export(gameID string) ([]byte, error) {
	g, err := r.GetByID(gameID)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	if err = r.db.ExportFolder(&buf, g.ID, r.gamesPath); err != nil {
		return nil, fmt.Errorf("export game %q: %w", g.ID, err)
	}
	return buf.Bytes(), nil
}

func (r *game) Import(data []byte, name string) (*entitiesGame.Game, error) {
	id, err := r.db.ImportFolder(bytes.NewReader(data), name, r.gamesPath)
	if err != nil {
		if errors.Is(err, fsentry.ErrBadArchive) {
			return nil, apperr.ErrBadArchive
		}
		return nil, repositories.MapFsentry(err, repositories.FsentrySentinels{
			Exist: apperr.ErrGameExists,
		})
	}

	g, err := r.GetByID(id)
	if err != nil {
		if delErr := r.DeleteByID(id); delErr != nil {
			slog.Error("remove half-imported game", "game", id, "err", delErr)
		}
		return nil, err
	}
	return g, nil
}

func toEntity(info repositories.FolderInfo, imageErr error) *entitiesGame.Game {
	createdAt, updatedAt := utils.NormalizeTimestamps(info.CreatedAt, info.UpdatedAt)
	return &entitiesGame.Game{
		ID:          info.ID,
		Name:        info.Name,
		Description: info.Data.Description.String(),
		Image:       info.Data.Image.String(),
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
		ImageError:  imageErr,
	}
}
