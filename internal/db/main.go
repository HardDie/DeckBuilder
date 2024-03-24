package db

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/HardDie/fsentry"
	"github.com/HardDie/fsentry/pkg/fsentry_error"
	"github.com/HardDie/fsentry/pkg/fsentry_types"

	"github.com/HardDie/DeckBuilder/internal/entity"
	er "github.com/HardDie/DeckBuilder/internal/errors"
	"github.com/HardDie/DeckBuilder/internal/logger"
	"github.com/HardDie/DeckBuilder/internal/utils"
)

type DB struct {
	db         fsentry.IFSEntry
	gamesPath  string
	resultPath string
}

func NewFSEntryDB(db fsentry.IFSEntry) *DB {
	return &DB{
		db:         db,
		gamesPath:  "games",
		resultPath: "result",
	}
}

// Game

type commonInfo struct {
	Description fsentry_types.QuotedString `json:"description"`
	Image       fsentry_types.QuotedString `json:"image"`
}
type cardInfo struct {
	ID          int64                                 `json:"id"`
	Name        fsentry_types.QuotedString            `json:"name"`
	Description fsentry_types.QuotedString            `json:"description"`
	Image       fsentry_types.QuotedString            `json:"image"`
	Variables   map[string]fsentry_types.QuotedString `json:"variables"`
	Count       int                                   `json:"count"`
	CreatedAt   *time.Time                            `json:"createdAt"`
	UpdatedAt   *time.Time                            `json:"updatedAt"`
}

func (s *DB) Init() error {
	err := s.db.Init()
	if err != nil {
		return er.InternalError.AddMessage(err.Error())
	}
	_, err = s.db.CreateFolder(s.gamesPath, nil)
	if err != nil {
		if !errors.Is(err, fsentry_error.ErrorExist) {
			return er.InternalError.AddMessage(err.Error())
		}
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

func (s *DB) GameCreate(_ context.Context, name, description, image string) (*entity.GameInfo, error) {
	info, err := s.db.CreateFolder(name, &commonInfo{
		Description: fsentry_types.QS(description),
		Image:       fsentry_types.QS(image),
	}, s.gamesPath)
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
func (s *DB) GameGet(ctx context.Context, name string) (context.Context, *entity.GameInfo, error) {
	info, err := s.db.GetFolder(name, s.gamesPath)
	if err != nil {
		if errors.Is(err, fsentry_error.ErrorNotExist) {
			return ctx, nil, er.GameNotExists.AddMessage(err.Error()).HTTP(http.StatusBadRequest)
		} else if errors.Is(err, fsentry_error.ErrorBadName) {
			return ctx, nil, er.BadName
		} else {
			return ctx, nil, er.InternalError.AddMessage(err.Error())
		}
	}
	var gInfo commonInfo

	err = json.Unmarshal(info.Data, &gInfo)
	if err != nil {
		return ctx, nil, er.InternalError.AddMessage(err.Error())
	}

	ctx = context.WithValue(ctx, "gameID", info.Id)
	return ctx, &entity.GameInfo{
		ID:        info.Id,
		Name:      info.Name.String(),
		CreatedAt: info.CreatedAt,
		UpdatedAt: info.UpdatedAt,

		Description: gInfo.Description.String(),
		Image:       gInfo.Image.String(),
	}, nil
}
func (s *DB) GameList(ctx context.Context) ([]*entity.GameInfo, error) {
	list, err := s.db.List(s.gamesPath)
	if err != nil {
		return nil, er.InternalError.AddMessage(err.Error())
	}

	var games []*entity.GameInfo
	for _, folder := range list.Folders {
		_, game, err := s.GameGet(ctx, folder)
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
func (s *DB) GameMove(_ context.Context, oldName, newName string) (*entity.GameInfo, error) {
	info, err := s.db.MoveFolder(oldName, newName, s.gamesPath)
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
func (s *DB) GameUpdate(_ context.Context, name, description, image string) (*entity.GameInfo, error) {
	info, err := s.db.UpdateFolder(name, &commonInfo{
		Description: fsentry_types.QS(description),
		Image:       fsentry_types.QS(image),
	}, s.gamesPath)
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
func (s *DB) GameDelete(_ context.Context, name string) error {
	err := s.db.RemoveFolder(name, s.gamesPath)
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
func (s *DB) GameDuplicate(_ context.Context, srcName, dstName string) (*entity.GameInfo, error) {
	info, err := s.db.DuplicateFolder(srcName, dstName, s.gamesPath)
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
func (s *DB) GameUpdateInfo(_ context.Context, name, newName string) error {
	return s.db.UpdateFolderNameWithoutTimestamp(name, newName, s.gamesPath)
}
func (s *DB) GameImageCreate(ctx context.Context, gameID string, data []byte) error {
	_, game, err := s.GameGet(ctx, gameID)
	if err != nil {
		return err
	}

	err = s.db.CreateBinary("image", data, s.gamesPath, game.ID)
	if err != nil {
		if errors.Is(err, fsentry_error.ErrorExist) {
			return er.GameImageExist.AddMessage(err.Error())
		} else {
			return er.InternalError.AddMessage(err.Error())
		}
	}
	return nil
}
func (s *DB) GameImageGet(ctx context.Context, gameID string) ([]byte, error) {
	_, game, err := s.GameGet(ctx, gameID)
	if err != nil {
		return nil, err
	}

	data, err := s.db.GetBinary("image", s.gamesPath, game.ID)
	if err != nil {
		if errors.Is(err, fsentry_error.ErrorNotExist) {
			return nil, er.GameImageNotExists.AddMessage(err.Error())
		} else {
			return nil, er.InternalError.AddMessage(err.Error())
		}
	}
	return data, nil
}
func (s *DB) GameImageDelete(ctx context.Context, gameID string) error {
	_, game, err := s.GameGet(ctx, gameID)
	if err != nil {
		return err
	}

	err = s.db.RemoveBinary("image", s.gamesPath, game.ID)
	if err != nil {
		if errors.Is(err, fsentry_error.ErrorNotExist) {
			return er.GameImageNotExists.AddMessage(err.Error())
		} else {
			return er.InternalError.AddMessage(err.Error())
		}
	}
	return nil
}

func (s *DB) CollectionCreate(ctx context.Context, gameID, name, description, image string) (*entity.CollectionInfo, error) {
	_, game, err := s.GameGet(ctx, gameID)
	if err != nil {
		return nil, err
	}

	info, err := s.db.CreateFolder(name, &commonInfo{
		Description: fsentry_types.QS(description),
		Image:       fsentry_types.QS(image),
	}, s.gamesPath, game.ID)
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
func (s *DB) CollectionGet(ctx context.Context, gameID, name string) (context.Context, *entity.CollectionInfo, error) {
	ctx, game, err := s.GameGet(ctx, gameID)
	if err != nil {
		return ctx, nil, err
	}

	info, err := s.db.GetFolder(name, s.gamesPath, game.ID)
	if err != nil {
		if errors.Is(err, fsentry_error.ErrorNotExist) {
			return ctx, nil, er.CollectionNotExists.AddMessage(err.Error()).HTTP(http.StatusBadRequest)
		} else if errors.Is(err, fsentry_error.ErrorBadName) {
			return ctx, nil, er.BadName
		} else {
			return ctx, nil, er.InternalError.AddMessage(err.Error())
		}
	}
	var cInfo commonInfo

	err = json.Unmarshal(info.Data, &cInfo)
	if err != nil {
		return ctx, nil, er.InternalError.AddMessage(err.Error())
	}

	ctx = context.WithValue(ctx, "collectionID", info.Id)
	return ctx, &entity.CollectionInfo{
		ID:        info.Id,
		Name:      info.Name.String(),
		CreatedAt: info.CreatedAt,
		UpdatedAt: info.UpdatedAt,

		Description: cInfo.Description.String(),
		Image:       cInfo.Image.String(),
	}, nil
}
func (s *DB) CollectionList(ctx context.Context, gameID string) ([]*entity.CollectionInfo, error) {
	_, game, err := s.GameGet(ctx, gameID)
	if err != nil {
		return nil, err
	}

	list, err := s.db.List(s.gamesPath, game.ID)
	if err != nil {
		return nil, er.InternalError.AddMessage(err.Error())
	}

	var collections []*entity.CollectionInfo
	for _, folder := range list.Folders {
		_, collection, err := s.CollectionGet(ctx, game.ID, folder)
		if err != nil {
			logger.Error.Println(folder, err.Error())
			continue
		}
		if folder != collection.ID {
			logger.Error.Println("Corrupted collection folder:", folder)
			continue
		}
		collections = append(collections, collection)
	}
	return collections, nil
}
func (s *DB) CollectionMove(ctx context.Context, gameID, oldName, newName string) (*entity.CollectionInfo, error) {
	_, game, err := s.GameGet(ctx, gameID)
	if err != nil {
		return nil, err
	}

	info, err := s.db.MoveFolder(oldName, newName, s.gamesPath, game.ID)
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
func (s *DB) CollectionUpdate(ctx context.Context, gameID, name, description, image string) (*entity.CollectionInfo, error) {
	_, game, err := s.GameGet(ctx, gameID)
	if err != nil {
		return nil, err
	}

	info, err := s.db.UpdateFolder(name, &commonInfo{
		Description: fsentry_types.QS(description),
		Image:       fsentry_types.QS(image),
	}, s.gamesPath, game.ID)
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
func (s *DB) CollectionDelete(ctx context.Context, gameID, name string) error {
	_, game, err := s.GameGet(ctx, gameID)
	if err != nil {
		return err
	}

	err = s.db.RemoveFolder(name, s.gamesPath, game.ID)
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
func (s *DB) CollectionImageCreate(ctx context.Context, gameID, collectionID string, data []byte) error {
	ctx, collection, err := s.CollectionGet(ctx, gameID, collectionID)
	if err != nil {
		return err
	}
	ctxGameID := ctx.Value("gameID").(string)

	err = s.db.CreateBinary("image", data, s.gamesPath, ctxGameID, collection.ID)
	if err != nil {
		if errors.Is(err, fsentry_error.ErrorExist) {
			return er.CollectionImageExist.AddMessage(err.Error())
		} else {
			return er.InternalError.AddMessage(err.Error())
		}
	}
	return nil
}
func (s *DB) CollectionImageGet(ctx context.Context, gameID, collectionID string) ([]byte, error) {
	ctx, collection, err := s.CollectionGet(ctx, gameID, collectionID)
	if err != nil {
		return nil, err
	}
	ctxGameID := ctx.Value("gameID").(string)

	data, err := s.db.GetBinary("image", s.gamesPath, ctxGameID, collection.ID)
	if err != nil {
		if errors.Is(err, fsentry_error.ErrorNotExist) {
			return nil, er.CollectionImageNotExists.AddMessage(err.Error())
		} else {
			return nil, er.InternalError.AddMessage(err.Error())
		}
	}
	return data, nil
}
func (s *DB) CollectionImageDelete(ctx context.Context, gameID, collectionID string) error {
	ctx, collection, err := s.CollectionGet(ctx, gameID, collectionID)
	if err != nil {
		return err
	}
	ctxGameID := ctx.Value("gameID").(string)

	err = s.db.RemoveBinary("image", s.gamesPath, ctxGameID, collection.ID)
	if err != nil {
		if errors.Is(err, fsentry_error.ErrorNotExist) {
			return er.CollectionImageNotExists.AddMessage(err.Error())
		} else {
			return er.InternalError.AddMessage(err.Error())
		}
	}
	return nil
}

func (s *DB) DeckCreate(ctx context.Context, gameID, collectionID, name, description, image string) (*entity.DeckInfo, error) {
	ctx, collection, err := s.CollectionGet(ctx, gameID, collectionID)
	if err != nil {
		return nil, err
	}
	ctxGameID := ctx.Value("gameID").(string)

	info, err := s.db.CreateFolder(name, &commonInfo{
		Description: fsentry_types.QS(description),
		Image:       fsentry_types.QS(image),
	}, s.gamesPath, ctxGameID, collection.ID)
	if err != nil {
		if errors.Is(err, fsentry_error.ErrorExist) {
			return nil, er.DeckExist
		} else if errors.Is(err, fsentry_error.ErrorBadName) {
			return nil, er.BadName
		} else {
			return nil, er.InternalError.AddMessage(err.Error())
		}
	}

	// Create folder for cards
	_, err = s.db.CreateFolder("cards", nil, s.gamesPath, ctxGameID, collection.ID, info.Id)
	if err != nil {
		return nil, er.InternalError.AddMessage(err.Error())
	}

	return &entity.DeckInfo{
		ID:        info.Id,
		Name:      info.Name.String(),
		CreatedAt: info.CreatedAt,
		UpdatedAt: info.UpdatedAt,

		Description: description,
		Image:       image,
	}, nil
}
func (s *DB) DeckGet(ctx context.Context, gameID, collectionID, name string) (context.Context, *entity.DeckInfo, error) {
	ctx, collection, err := s.CollectionGet(ctx, gameID, collectionID)
	if err != nil {
		return ctx, nil, err
	}
	ctxGameID := ctx.Value("gameID").(string)

	info, err := s.db.GetFolder(name, s.gamesPath, ctxGameID, collection.ID)
	if err != nil {
		if errors.Is(err, fsentry_error.ErrorNotExist) {
			return ctx, nil, er.DeckNotExists.AddMessage(err.Error()).HTTP(http.StatusBadRequest)
		} else if errors.Is(err, fsentry_error.ErrorBadName) {
			return ctx, nil, er.BadName
		} else {
			return ctx, nil, er.InternalError.AddMessage(err.Error())
		}
	}
	var dInfo commonInfo

	err = json.Unmarshal(info.Data, &dInfo)
	if err != nil {
		return ctx, nil, er.InternalError.AddMessage(err.Error())
	}

	ctx = context.WithValue(ctx, "deckID", info.Id)
	return ctx, &entity.DeckInfo{
		ID:        info.Id,
		Name:      info.Name.String(),
		CreatedAt: info.CreatedAt,
		UpdatedAt: info.UpdatedAt,

		Description: dInfo.Description.String(),
		Image:       dInfo.Image.String(),
	}, nil
}
func (s *DB) DeckList(ctx context.Context, gameID, collectionID string) ([]*entity.DeckInfo, error) {
	ctx, collection, err := s.CollectionGet(ctx, gameID, collectionID)
	if err != nil {
		return nil, err
	}
	ctxGameID := ctx.Value("gameID").(string)

	list, err := s.db.List(s.gamesPath, ctxGameID, collection.ID)
	if err != nil {
		return nil, er.InternalError.AddMessage(err.Error())
	}

	var decks []*entity.DeckInfo
	for _, folder := range list.Folders {
		_, deck, err := s.DeckGet(ctx, ctxGameID, collection.ID, folder)
		if err != nil {
			logger.Error.Println(folder, err.Error())
			continue
		}
		if folder != deck.ID {
			logger.Error.Println("Corrupted deck folder:", folder)
			continue
		}
		decks = append(decks, deck)
	}
	return decks, nil
}
func (s *DB) DeckMove(ctx context.Context, gameID, collectionID, oldName, newName string) (*entity.DeckInfo, error) {
	ctx, collection, err := s.CollectionGet(ctx, gameID, collectionID)
	if err != nil {
		return nil, err
	}
	ctxGameID := ctx.Value("gameID").(string)

	info, err := s.db.MoveFolder(oldName, newName, s.gamesPath, ctxGameID, collection.ID)
	if err != nil {
		if errors.Is(err, fsentry_error.ErrorNotExist) {
			return nil, er.DeckNotExists.AddMessage(err.Error()).HTTP(http.StatusBadRequest)
		} else if errors.Is(err, fsentry_error.ErrorBadName) {
			return nil, er.BadName
		} else {
			return nil, er.InternalError.AddMessage(err.Error())
		}
	}
	var dInfo commonInfo

	err = json.Unmarshal(info.Data, &dInfo)
	if err != nil {
		return nil, er.InternalError.AddMessage(err.Error())
	}

	return &entity.DeckInfo{
		ID:        info.Id,
		Name:      info.Name.String(),
		CreatedAt: info.CreatedAt,
		UpdatedAt: info.UpdatedAt,

		Description: dInfo.Description.String(),
		Image:       dInfo.Image.String(),
	}, nil
}
func (s *DB) DeckUpdate(ctx context.Context, gameID, collectionID, name, description, image string) (*entity.DeckInfo, error) {
	ctx, collection, err := s.CollectionGet(ctx, gameID, collectionID)
	if err != nil {
		return nil, err
	}
	ctxGameID := ctx.Value("gameID").(string)

	info, err := s.db.UpdateFolder(name, &commonInfo{
		Description: fsentry_types.QS(description),
		Image:       fsentry_types.QS(image),
	}, s.gamesPath, ctxGameID, collection.ID)
	if err != nil {
		if errors.Is(err, fsentry_error.ErrorNotExist) {
			return nil, er.DeckNotExists.AddMessage(err.Error()).HTTP(http.StatusBadRequest)
		} else if errors.Is(err, fsentry_error.ErrorBadName) {
			return nil, er.BadName
		} else {
			return nil, er.InternalError.AddMessage(err.Error())
		}
	}
	var dInfo commonInfo

	err = json.Unmarshal(info.Data, &dInfo)
	if err != nil {
		return nil, er.InternalError.AddMessage(err.Error())
	}

	return &entity.DeckInfo{
		ID:        info.Id,
		Name:      info.Name.String(),
		CreatedAt: info.CreatedAt,
		UpdatedAt: info.UpdatedAt,

		Description: dInfo.Description.String(),
		Image:       dInfo.Image.String(),
	}, nil
}
func (s *DB) DeckDelete(ctx context.Context, gameID, collectionID, name string) error {
	ctx, collection, err := s.CollectionGet(ctx, gameID, collectionID)
	if err != nil {
		return err
	}
	ctxGameID := ctx.Value("gameID").(string)

	err = s.db.RemoveFolder(name, s.gamesPath, ctxGameID, collection.ID)
	if err != nil {
		if errors.Is(err, fsentry_error.ErrorNotExist) {
			return er.DeckNotExists.AddMessage(err.Error()).HTTP(http.StatusBadRequest)
		} else if errors.Is(err, fsentry_error.ErrorBadName) {
			return er.BadName
		} else {
			return er.InternalError.AddMessage(err.Error())
		}
	}
	return nil
}
func (s *DB) DeckImageCreate(ctx context.Context, gameID, collectionID, deckID string, data []byte) error {
	ctx, deck, err := s.DeckGet(ctx, gameID, collectionID, deckID)
	if err != nil {
		return err
	}
	ctxGameID := ctx.Value("gameID").(string)
	ctxCollectionID := ctx.Value("collectionID").(string)

	err = s.db.CreateBinary("image", data, s.gamesPath, ctxGameID, ctxCollectionID, deck.ID)
	if err != nil {
		if errors.Is(err, fsentry_error.ErrorExist) {
			return er.DeckImageExist.AddMessage(err.Error())
		} else {
			return er.InternalError.AddMessage(err.Error())
		}
	}
	return nil
}
func (s *DB) DeckImageGet(ctx context.Context, gameID, collectionID, deckID string) ([]byte, error) {
	ctx, deck, err := s.DeckGet(ctx, gameID, collectionID, deckID)
	if err != nil {
		return nil, err
	}
	ctxGameID := ctx.Value("gameID").(string)
	ctxCollectionID := ctx.Value("collectionID").(string)

	data, err := s.db.GetBinary("image", s.gamesPath, ctxGameID, ctxCollectionID, deck.ID)
	if err != nil {
		if errors.Is(err, fsentry_error.ErrorNotExist) {
			return nil, er.DeckImageNotExists.AddMessage(err.Error())
		} else {
			return nil, er.InternalError.AddMessage(err.Error())
		}
	}
	return data, nil
}
func (s *DB) DeckImageDelete(ctx context.Context, gameID, collectionID, deckID string) error {
	ctx, deck, err := s.DeckGet(ctx, gameID, collectionID, deckID)
	if err != nil {
		return err
	}
	ctxGameID := ctx.Value("gameID").(string)
	ctxCollectionID := ctx.Value("collectionID").(string)

	err = s.db.RemoveBinary("image", s.gamesPath, ctxGameID, ctxCollectionID, deck.ID)
	if err != nil {
		if errors.Is(err, fsentry_error.ErrorNotExist) {
			return er.DeckImageNotExists.AddMessage(err.Error())
		} else {
			return er.InternalError.AddMessage(err.Error())
		}
	}
	return nil
}

func (s *DB) CardCreate(ctx context.Context, gameID, collectionID, deckID, name, description, image string, variables map[string]string, count int) (*entity.CardInfo, error) {
	ctx, list, err := s.rawCardList(ctx, gameID, collectionID, deckID)
	if err != nil {
		return nil, err
	}
	ctxGameID := ctx.Value("gameID").(string)
	ctxCollectionID := ctx.Value("collectionID").(string)
	ctxDeckID := ctx.Value("deckID").(string)

	// Search for the largest card ID
	maxID := int64(1)
	for _, card := range list {
		if card.ID >= maxID {
			maxID = card.ID + 1
		}
	}

	// Create a card with the found identifier
	cardInfo := &cardInfo{
		ID:          maxID,
		Name:        fsentry_types.QS(name),
		Description: fsentry_types.QS(description),
		Image:       fsentry_types.QS(image),
		Variables:   convertMapString(variables),
		Count:       count,
		CreatedAt:   utils.Allocate(time.Now()),
		UpdatedAt:   nil,
	}

	// Add a card to the card array
	list[cardInfo.ID] = cardInfo

	// Writing an array of cards to a file again
	_, err = s.db.UpdateFolder("cards", list, s.gamesPath, ctxGameID, ctxCollectionID, ctxDeckID)
	if err != nil {
		if errors.Is(err, fsentry_error.ErrorNotExist) {
			return nil, er.CardNotExists.AddMessage(err.Error())
		} else if errors.Is(err, fsentry_error.ErrorBadName) {
			return nil, er.BadName
		} else {
			return nil, er.InternalError.AddMessage(err.Error())
		}
	}

	return &entity.CardInfo{
		ID:          cardInfo.ID,
		Name:        cardInfo.Name.String(),
		Description: cardInfo.Description.String(),
		Image:       cardInfo.Image.String(),
		Variables:   convertMapQuotedString(cardInfo.Variables),
		Count:       cardInfo.Count,
		CreatedAt:   cardInfo.CreatedAt,
		UpdatedAt:   nil,
	}, nil
}
func (s *DB) CardGet(ctx context.Context, gameID, collectionID, deckID string, cardID int64) (context.Context, *entity.CardInfo, error) {
	ctx, list, err := s.rawCardList(ctx, gameID, collectionID, deckID)
	if err != nil {
		return ctx, nil, err
	}

	card, ok := list[cardID]
	if !ok {
		return ctx, nil, er.CardNotExists.HTTP(http.StatusBadRequest)
	}

	return ctx, &entity.CardInfo{
		ID:          card.ID,
		Name:        card.Name.String(),
		Description: card.Description.String(),
		Image:       card.Image.String(),
		Variables:   convertMapQuotedString(card.Variables),
		Count:       card.Count,
		CreatedAt:   card.CreatedAt,
		UpdatedAt:   card.UpdatedAt,
	}, nil
}
func (s *DB) CardList(ctx context.Context, gameID, collectionID, deckID string) ([]*entity.CardInfo, error) {
	ctx, list, err := s.rawCardList(ctx, gameID, collectionID, deckID)
	if err != nil {
		return nil, err
	}

	var cards []*entity.CardInfo
	for _, item := range list {
		cards = append(cards, &entity.CardInfo{
			ID:          item.ID,
			Name:        item.Name.String(),
			Description: item.Description.String(),
			Image:       item.Image.String(),
			Variables:   convertMapQuotedString(item.Variables),
			Count:       item.Count,
			CreatedAt:   item.CreatedAt,
			UpdatedAt:   item.UpdatedAt,
		})
	}
	return cards, nil
}
func (s *DB) CardUpdate(ctx context.Context, gameID, collectionID, deckID string, cardID int64, name, description, image string, variables map[string]string, count int) (*entity.CardInfo, error) {
	ctx, list, err := s.rawCardList(ctx, gameID, collectionID, deckID)
	if err != nil {
		return nil, err
	}
	ctxGameID := ctx.Value("gameID").(string)
	ctxCollectionID := ctx.Value("collectionID").(string)
	ctxDeckID := ctx.Value("deckID").(string)

	card, ok := list[cardID]
	if !ok {
		return nil, er.CardNotExists
	}

	card.Name = fsentry_types.QS(name)
	card.Description = fsentry_types.QS(description)
	card.Image = fsentry_types.QS(image)
	card.Variables = convertMapString(variables)
	card.Count = count
	card.UpdatedAt = utils.Allocate(time.Now())

	list[card.ID] = card

	// Writing an array of cards to a file again
	_, err = s.db.UpdateFolder("cards", list, s.gamesPath, ctxGameID, ctxCollectionID, ctxDeckID)
	if err != nil {
		if errors.Is(err, fsentry_error.ErrorNotExist) {
			return nil, er.CardNotExists.AddMessage(err.Error())
		} else if errors.Is(err, fsentry_error.ErrorBadName) {
			return nil, er.BadName
		} else {
			return nil, er.InternalError.AddMessage(err.Error())
		}
	}

	return &entity.CardInfo{
		ID:          card.ID,
		Name:        card.Name.String(),
		Description: card.Description.String(),
		Image:       card.Image.String(),
		Variables:   convertMapQuotedString(card.Variables),
		Count:       card.Count,
		CreatedAt:   card.CreatedAt,
		UpdatedAt:   card.UpdatedAt,
	}, nil
}
func (s *DB) CardDelete(ctx context.Context, gameID, collectionID, deckID string, cardID int64) error {
	ctx, list, err := s.rawCardList(ctx, gameID, collectionID, deckID)
	if err != nil {
		return err
	}
	ctxGameID := ctx.Value("gameID").(string)
	ctxCollectionID := ctx.Value("collectionID").(string)
	ctxDeckID := ctx.Value("deckID").(string)

	if _, ok := list[cardID]; !ok {
		return er.CardNotExists
	}

	delete(list, cardID)

	// Writing an array of cards to a file again
	_, err = s.db.UpdateFolder("cards", list, s.gamesPath, ctxGameID, ctxCollectionID, ctxDeckID)
	if err != nil {
		if errors.Is(err, fsentry_error.ErrorNotExist) {
			return er.CardNotExists.AddMessage(err.Error())
		} else if errors.Is(err, fsentry_error.ErrorBadName) {
			return er.BadName
		} else {
			return er.InternalError.AddMessage(err.Error())
		}
	}

	return nil
}
func (s *DB) rawCardList(ctx context.Context, gameID, collectionID, deckID string) (context.Context, map[int64]*cardInfo, error) {
	ctx, deck, err := s.DeckGet(ctx, gameID, collectionID, deckID)
	if err != nil {
		return ctx, nil, err
	}
	ctxGameID := ctx.Value("gameID").(string)
	ctxCollectionID := ctx.Value("collectionID").(string)

	// Get all the cards
	info, err := s.db.GetFolder("cards", s.gamesPath, ctxGameID, ctxCollectionID, deck.ID)
	if err != nil {
		return ctx, nil, er.InternalError.AddMessage(err.Error())
	}

	// Parsing an array of cards from json
	var list map[int64]*cardInfo
	err = json.Unmarshal(info.Data, &list)
	if err != nil {
		return ctx, nil, er.InternalError.AddMessage(err.Error())
	}

	if list == nil {
		list = make(map[int64]*cardInfo)
	}
	return ctx, list, nil
}
func (s *DB) CardImageCreate(ctx context.Context, gameID, collectionID, deckID string, cardID int64, data []byte) error {
	ctx, card, err := s.CardGet(ctx, gameID, collectionID, deckID, cardID)
	if err != nil {
		return err
	}
	ctxGameID := ctx.Value("gameID").(string)
	ctxCollectionID := ctx.Value("collectionID").(string)
	ctxDeckID := ctx.Value("deckID").(string)

	err = s.db.CreateBinary(fmt.Sprintf("%d", card.ID), data, s.gamesPath, ctxGameID, ctxCollectionID, ctxDeckID, "cards")
	if err != nil {
		if errors.Is(err, fsentry_error.ErrorExist) {
			return er.CardImageExist.AddMessage(err.Error())
		} else {
			return er.InternalError.AddMessage(err.Error())
		}
	}
	return nil
}
func (s *DB) CardImageGet(ctx context.Context, gameID, collectionID, deckID string, cardID int64) ([]byte, error) {
	ctx, card, err := s.CardGet(ctx, gameID, collectionID, deckID, cardID)
	if err != nil {
		return nil, err
	}
	ctxGameID := ctx.Value("gameID").(string)
	ctxCollectionID := ctx.Value("collectionID").(string)
	ctxDeckID := ctx.Value("deckID").(string)

	data, err := s.db.GetBinary(fmt.Sprintf("%d", card.ID), s.gamesPath, ctxGameID, ctxCollectionID, ctxDeckID, "cards")
	if err != nil {
		if errors.Is(err, fsentry_error.ErrorNotExist) {
			return nil, er.CardImageNotExists.AddMessage(err.Error())
		} else {
			return nil, er.InternalError.AddMessage(err.Error())
		}
	}
	return data, nil
}
func (s *DB) CardImageDelete(ctx context.Context, gameID, collectionID, deckID string, cardID int64) error {
	ctx, card, err := s.CardGet(ctx, gameID, collectionID, deckID, cardID)
	if err != nil {
		return err
	}
	ctxGameID := ctx.Value("gameID").(string)
	ctxCollectionID := ctx.Value("collectionID").(string)
	ctxDeckID := ctx.Value("deckID").(string)

	err = s.db.RemoveBinary(fmt.Sprintf("%d", card.ID), s.gamesPath, ctxGameID, ctxCollectionID, ctxDeckID, "cards")
	if err != nil {
		if errors.Is(err, fsentry_error.ErrorNotExist) {
			return er.CardImageNotExists.AddMessage(err.Error())
		} else {
			return er.InternalError.AddMessage(err.Error())
		}
	}
	return nil
}

func (s *DB) SettingsGet() (*entity.SettingInfo, error) {
	info, err := s.db.GetEntry("settings")
	if err != nil {
		if errors.Is(err, fsentry_error.ErrorNotExist) {
			return nil, er.SettingsNotExists.AddMessage(err.Error())
		} else {
			return nil, er.InternalError.AddMessage(err.Error())
		}
	}
	setting := &entity.SettingInfo{}

	err = json.Unmarshal(info.Data, setting)
	if err != nil {
		return nil, er.InternalError.AddMessage(err.Error())
	}

	return setting, nil
}
func (s *DB) SettingsSet(data *entity.SettingInfo) error {
	err := s.db.CreateEntry("settings", data)
	if err == nil {
		return nil
	}
	if !errors.Is(err, fsentry_error.ErrorExist) {
		return err
	}
	err = s.db.UpdateEntry("settings", data)
	if err != nil {
		return er.InternalError.AddMessage(err.Error())
	}
	return nil
}

func convertMapString(in map[string]string) map[string]fsentry_types.QuotedString {
	res := make(map[string]fsentry_types.QuotedString)
	for key, val := range in {
		keyJson, _ := json.Marshal(strconv.Quote(key))
		res[string(keyJson)] = fsentry_types.QS(val)
	}
	return res
}
func convertMapQuotedString(in map[string]fsentry_types.QuotedString) map[string]string {
	res := make(map[string]string)
	for keyJson, val := range in {
		var key string
		_ = json.Unmarshal([]byte(keyJson), &key)
		key, _ = strconv.Unquote(key)
		res[key] = val.String()
	}
	return res
}
