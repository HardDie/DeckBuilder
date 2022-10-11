package db

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/HardDie/fsentry"
	"github.com/HardDie/fsentry/pkg/fsentry_error"
	"github.com/HardDie/fsentry/pkg/fsentry_types"

	"github.com/HardDie/DeckBuilder/internal/entity"
	er "github.com/HardDie/DeckBuilder/internal/errors"
	"github.com/HardDie/DeckBuilder/internal/logger"
)

type DB struct {
	db fsentry.IFSEntry
}

func NewFSEntryDB(db fsentry.IFSEntry) *DB {
	return &DB{
		db: db,
	}
}

// Game

type commonInfo struct {
	Description fsentry_types.QuotedString `json:"description"`
	Image       fsentry_types.QuotedString `json:"image"`
}

func (s *DB) Init() error {
	err := s.db.Init()
	if err != nil {
		return er.InternalError.AddMessage(err.Error())
	}
	return nil
}
func (s *DB) Drop() error {
	err := s.db.Drop()
	if err != nil {
		return er.InternalError.AddMessage(err.Error())
	}
	return nil
}

func (s *DB) GameCreate(name, description, image string) (*entity.GameInfo, error) {
	info, err := s.db.CreateFolder(name, &commonInfo{
		Description: fsentry_types.QS(description),
		Image:       fsentry_types.QS(image),
	})
	if err != nil {
		if errors.Is(err, fsentry_error.ErrorExist) {
			return nil, er.GameExist
		} else if errors.Is(err, fsentry_error.ErrorBadName) {
			return nil, er.BadName
		} else {
			return nil, er.InternalError.AddMessage(err.Error())
		}
	}

	return &entity.GameInfo{
		ID:        info.Id,
		Name:      info.Name.String(),
		CreatedAt: info.CreatedAt,
		UpdatedAt: info.UpdatedAt,

		Description: description,
		Image:       image,
	}, nil
}
func (s *DB) GameGet(name string) (*entity.GameInfo, error) {
	info, err := s.db.GetFolder(name)
	if err != nil {
		if errors.Is(err, fsentry_error.ErrorNotExist) {
			return nil, er.GameNotExists.AddMessage(err.Error()).HTTP(http.StatusBadRequest)
		} else if errors.Is(err, fsentry_error.ErrorBadName) {
			return nil, er.BadName
		} else {
			return nil, er.InternalError.AddMessage(err.Error())
		}
	}
	var gInfo commonInfo

	err = json.Unmarshal(info.Data, &gInfo)
	if err != nil {
		return nil, er.InternalError.AddMessage(err.Error())
	}

	return &entity.GameInfo{
		ID:        info.Id,
		Name:      info.Name.String(),
		CreatedAt: info.CreatedAt,
		UpdatedAt: info.UpdatedAt,

		Description: gInfo.Description.String(),
		Image:       gInfo.Image.String(),
	}, nil
}
func (s *DB) GameList() ([]*entity.GameInfo, error) {
	list, err := s.db.List()
	if err != nil {
		return nil, er.InternalError.AddMessage(err.Error())
	}

	var games []*entity.GameInfo
	for _, folder := range list.Folders {
		game, err := s.GameGet(folder)
		if err != nil {
			logger.Error.Println(folder, err.Error())
			continue
		}
		games = append(games, game)
	}
	return games, nil
}
func (s *DB) GameMove(oldName, newName string) (*entity.GameInfo, error) {
	info, err := s.db.MoveFolder(oldName, newName)
	if err != nil {
		if errors.Is(err, fsentry_error.ErrorNotExist) {
			return nil, er.GameNotExists.AddMessage(err.Error()).HTTP(http.StatusBadRequest)
		} else if errors.Is(err, fsentry_error.ErrorBadName) {
			return nil, er.BadName
		} else {
			return nil, er.InternalError.AddMessage(err.Error())
		}
	}
	var gInfo commonInfo

	err = json.Unmarshal(info.Data, &gInfo)
	if err != nil {
		return nil, er.InternalError.AddMessage(err.Error())
	}

	return &entity.GameInfo{
		ID:        info.Id,
		Name:      info.Name.String(),
		CreatedAt: info.CreatedAt,
		UpdatedAt: info.UpdatedAt,

		Description: gInfo.Description.String(),
		Image:       gInfo.Image.String(),
	}, nil
}
func (s *DB) GameUpdate(name, description, image string) (*entity.GameInfo, error) {
	info, err := s.db.UpdateFolder(name, &commonInfo{
		Description: fsentry_types.QS(description),
		Image:       fsentry_types.QS(image),
	})
	if err != nil {
		if errors.Is(err, fsentry_error.ErrorNotExist) {
			return nil, er.GameNotExists.AddMessage(err.Error()).HTTP(http.StatusBadRequest)
		} else if errors.Is(err, fsentry_error.ErrorBadName) {
			return nil, er.BadName
		} else {
			return nil, er.InternalError.AddMessage(err.Error())
		}
	}
	var gInfo commonInfo

	err = json.Unmarshal(info.Data, &gInfo)
	if err != nil {
		return nil, er.InternalError.AddMessage(err.Error())
	}

	return &entity.GameInfo{
		ID:        info.Id,
		Name:      info.Name.String(),
		CreatedAt: info.CreatedAt,
		UpdatedAt: info.UpdatedAt,

		Description: gInfo.Description.String(),
		Image:       gInfo.Image.String(),
	}, nil
}
func (s *DB) GameDelete(name string) error {
	err := s.db.RemoveFolder(name)
	if err != nil {
		if errors.Is(err, fsentry_error.ErrorNotExist) {
			return er.GameNotExists.AddMessage(err.Error()).HTTP(http.StatusBadRequest)
		} else if errors.Is(err, fsentry_error.ErrorBadName) {
			return er.BadName
		} else {
			return er.InternalError.AddMessage(err.Error())
		}
	}
	return nil
}
func (s *DB) GameDuplicate(srcName, dstName string) (*entity.GameInfo, error) {
	info, err := s.db.DuplicateFolder(srcName, dstName)
	if err != nil {
		if errors.Is(err, fsentry_error.ErrorNotExist) {
			return nil, er.GameNotExists.AddMessage(err.Error()).HTTP(http.StatusBadRequest)
		} else if errors.Is(err, fsentry_error.ErrorExist) {
			return nil, er.GameExist.AddMessage(err.Error()).HTTP(http.StatusBadRequest)
		} else if errors.Is(err, fsentry_error.ErrorBadName) {
			return nil, er.BadName
		} else {
			return nil, er.InternalError.AddMessage(err.Error())
		}
	}
	var gInfo commonInfo

	err = json.Unmarshal(info.Data, &gInfo)
	if err != nil {
		return nil, er.InternalError.AddMessage(err.Error())
	}

	return &entity.GameInfo{
		ID:        info.Id,
		Name:      info.Name.String(),
		CreatedAt: info.CreatedAt,
		UpdatedAt: info.UpdatedAt,

		Description: gInfo.Description.String(),
		Image:       gInfo.Image.String(),
	}, nil
}
func (s *DB) GameUpdateInfo(name, newName string) error {
	return s.db.UpdateFolderNameWithoutTimestamp(name, newName)
}

func (s *DB) CollectionCreate(gameID, name, description, image string) (*entity.CollectionInfo, error) {
	game, err := s.GameGet(gameID)
	if err != nil {
		return nil, err
	}

	info, err := s.db.CreateFolder(name, &commonInfo{
		Description: fsentry_types.QS(description),
		Image:       fsentry_types.QS(image),
	}, game.ID)
	if err != nil {
		if errors.Is(err, fsentry_error.ErrorExist) {
			return nil, er.CollectionExist
		} else if errors.Is(err, fsentry_error.ErrorBadName) {
			return nil, er.BadName
		} else {
			return nil, er.InternalError.AddMessage(err.Error())
		}
	}

	return &entity.CollectionInfo{
		ID:        info.Id,
		Name:      info.Name.String(),
		CreatedAt: info.CreatedAt,
		UpdatedAt: info.UpdatedAt,

		Description: description,
		Image:       image,
	}, nil
}
func (s *DB) CollectionGet(gameID, name string) (*entity.CollectionInfo, error) {
	game, err := s.GameGet(gameID)
	if err != nil {
		return nil, err
	}

	info, err := s.db.GetFolder(name, game.ID)
	if err != nil {
		if errors.Is(err, fsentry_error.ErrorNotExist) {
			return nil, er.CollectionNotExists.AddMessage(err.Error()).HTTP(http.StatusBadRequest)
		} else if errors.Is(err, fsentry_error.ErrorBadName) {
			return nil, er.BadName
		} else {
			return nil, er.InternalError.AddMessage(err.Error())
		}
	}
	var cInfo commonInfo

	err = json.Unmarshal(info.Data, &cInfo)
	if err != nil {
		return nil, er.InternalError.AddMessage(err.Error())
	}

	return &entity.CollectionInfo{
		ID:        info.Id,
		Name:      info.Name.String(),
		CreatedAt: info.CreatedAt,
		UpdatedAt: info.UpdatedAt,

		Description: cInfo.Description.String(),
		Image:       cInfo.Image.String(),
	}, nil
}
func (s *DB) CollectionList(gameID string) ([]*entity.CollectionInfo, error) {
	game, err := s.GameGet(gameID)
	if err != nil {
		return nil, err
	}

	list, err := s.db.List(game.ID)
	if err != nil {
		return nil, er.InternalError.AddMessage(err.Error())
	}

	var collections []*entity.CollectionInfo
	for _, folder := range list.Folders {
		collection, err := s.CollectionGet(game.ID, folder)
		if err != nil {
			logger.Error.Println(folder, err.Error())
			continue
		}
		collections = append(collections, collection)
	}
	return collections, nil
}
func (s *DB) CollectionMove(gameID, oldName, newName string) (*entity.CollectionInfo, error) {
	game, err := s.GameGet(gameID)
	if err != nil {
		return nil, err
	}

	info, err := s.db.MoveFolder(oldName, newName, game.ID)
	if err != nil {
		if errors.Is(err, fsentry_error.ErrorNotExist) {
			return nil, er.CollectionNotExists.AddMessage(err.Error()).HTTP(http.StatusBadRequest)
		} else if errors.Is(err, fsentry_error.ErrorBadName) {
			return nil, er.BadName
		} else {
			return nil, er.InternalError.AddMessage(err.Error())
		}
	}
	var cInfo commonInfo

	err = json.Unmarshal(info.Data, &cInfo)
	if err != nil {
		return nil, er.InternalError.AddMessage(err.Error())
	}

	return &entity.CollectionInfo{
		ID:        info.Id,
		Name:      info.Name.String(),
		CreatedAt: info.CreatedAt,
		UpdatedAt: info.UpdatedAt,

		Description: cInfo.Description.String(),
		Image:       cInfo.Image.String(),
	}, nil
}
func (s *DB) CollectionUpdate(gameID, name, description, image string) (*entity.CollectionInfo, error) {
	game, err := s.GameGet(gameID)
	if err != nil {
		return nil, err
	}

	info, err := s.db.UpdateFolder(name, &commonInfo{
		Description: fsentry_types.QS(description),
		Image:       fsentry_types.QS(image),
	}, game.ID)
	if err != nil {
		if errors.Is(err, fsentry_error.ErrorNotExist) {
			return nil, er.CollectionNotExists.AddMessage(err.Error()).HTTP(http.StatusBadRequest)
		} else if errors.Is(err, fsentry_error.ErrorBadName) {
			return nil, er.BadName
		} else {
			return nil, er.InternalError.AddMessage(err.Error())
		}
	}
	var cInfo commonInfo

	err = json.Unmarshal(info.Data, &cInfo)
	if err != nil {
		return nil, er.InternalError.AddMessage(err.Error())
	}

	return &entity.CollectionInfo{
		ID:        info.Id,
		Name:      info.Name.String(),
		CreatedAt: info.CreatedAt,
		UpdatedAt: info.UpdatedAt,

		Description: cInfo.Description.String(),
		Image:       cInfo.Image.String(),
	}, nil
}
func (s *DB) CollectionDelete(gameID, name string) error {
	game, err := s.GameGet(gameID)
	if err != nil {
		return err
	}

	err = s.db.RemoveFolder(name, game.ID)
	if err != nil {
		if errors.Is(err, fsentry_error.ErrorNotExist) {
			return er.CollectionNotExists.AddMessage(err.Error()).HTTP(http.StatusBadRequest)
		} else if errors.Is(err, fsentry_error.ErrorBadName) {
			return er.BadName
		} else {
			return er.InternalError.AddMessage(err.Error())
		}
	}
	return nil
}
