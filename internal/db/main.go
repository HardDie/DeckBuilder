package db

import (
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
func (s *DB) GameImageCreate(name string, data []byte) error {
	game, err := s.GameGet(name)
	if err != nil {
		return err
	}

	err = s.db.CreateBinary("image", data, game.ID)
	if err != nil {
		if errors.Is(err, fsentry_error.ErrorExist) {
			return er.GameImageExist.AddMessage(err.Error())
		} else {
			return er.InternalError.AddMessage(err.Error())
		}
	}
	return nil
}
func (s *DB) GameImageGet(name string) ([]byte, error) {
	game, err := s.GameGet(name)
	if err != nil {
		return nil, err
	}

	data, err := s.db.GetBinary("image", game.ID)
	if err != nil {
		if errors.Is(err, fsentry_error.ErrorNotExist) {
			return nil, er.GameImageNotExists.AddMessage(err.Error())
		} else {
			return nil, er.InternalError.AddMessage(err.Error())
		}
	}
	return data, nil
}
func (s *DB) GameImageDelete(name string) error {
	game, err := s.GameGet(name)
	if err != nil {
		return err
	}

	err = s.db.RemoveBinary("image", game.ID)
	if err != nil {
		if errors.Is(err, fsentry_error.ErrorNotExist) {
			return er.GameImageNotExists.AddMessage(err.Error())
		} else {
			return er.InternalError.AddMessage(err.Error())
		}
	}
	return nil
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
func (s *DB) CollectionImageCreate(gameID, name string, data []byte) error {
	game, err := s.GameGet(gameID)
	if err != nil {
		return err
	}
	collection, err := s.CollectionGet(gameID, name)
	if err != nil {
		return err
	}

	err = s.db.CreateBinary("image", data, game.ID, collection.ID)
	if err != nil {
		if errors.Is(err, fsentry_error.ErrorExist) {
			return er.CollectionImageExist.AddMessage(err.Error())
		} else {
			return er.InternalError.AddMessage(err.Error())
		}
	}
	return nil
}
func (s *DB) CollectionImageGet(gameID, name string) ([]byte, error) {
	game, err := s.GameGet(gameID)
	if err != nil {
		return nil, err
	}
	collection, err := s.CollectionGet(gameID, name)
	if err != nil {
		return nil, err
	}

	data, err := s.db.GetBinary("image", game.ID, collection.ID)
	if err != nil {
		if errors.Is(err, fsentry_error.ErrorNotExist) {
			return nil, er.CollectionImageNotExists.AddMessage(err.Error())
		} else {
			return nil, er.InternalError.AddMessage(err.Error())
		}
	}
	return data, nil
}
func (s *DB) CollectionImageDelete(gameID, name string) error {
	game, err := s.GameGet(gameID)
	if err != nil {
		return err
	}
	collection, err := s.CollectionGet(gameID, name)
	if err != nil {
		return err
	}

	err = s.db.RemoveBinary("image", game.ID, collection.ID)
	if err != nil {
		if errors.Is(err, fsentry_error.ErrorNotExist) {
			return er.CollectionImageNotExists.AddMessage(err.Error())
		} else {
			return er.InternalError.AddMessage(err.Error())
		}
	}
	return nil
}

func (s *DB) DeckCreate(gameID, collectionID, name, description, image string) (*entity.DeckInfo, error) {
	game, err := s.GameGet(gameID)
	if err != nil {
		return nil, err
	}
	collection, err := s.CollectionGet(gameID, collectionID)
	if err != nil {
		return nil, err
	}

	info, err := s.db.CreateFolder(name, &commonInfo{
		Description: fsentry_types.QS(description),
		Image:       fsentry_types.QS(image),
	}, game.ID, collection.ID)
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
	_, err = s.db.CreateFolder("cards", nil, game.ID, collection.ID, info.Id)
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
func (s *DB) DeckGet(gameID, collectionID, name string) (*entity.DeckInfo, error) {
	game, err := s.GameGet(gameID)
	if err != nil {
		return nil, err
	}
	collection, err := s.CollectionGet(gameID, collectionID)
	if err != nil {
		return nil, err
	}

	info, err := s.db.GetFolder(name, game.ID, collection.ID)
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
func (s *DB) DeckList(gameID, collectionID string) ([]*entity.DeckInfo, error) {
	game, err := s.GameGet(gameID)
	if err != nil {
		return nil, err
	}
	collection, err := s.CollectionGet(gameID, collectionID)
	if err != nil {
		return nil, err
	}

	list, err := s.db.List(game.ID, collection.ID)
	if err != nil {
		return nil, er.InternalError.AddMessage(err.Error())
	}

	var decks []*entity.DeckInfo
	for _, folder := range list.Folders {
		deck, err := s.DeckGet(game.ID, collection.ID, folder)
		if err != nil {
			logger.Error.Println(folder, err.Error())
			continue
		}
		decks = append(decks, deck)
	}
	return decks, nil
}
func (s *DB) DeckMove(gameID, collectionID, oldName, newName string) (*entity.DeckInfo, error) {
	game, err := s.GameGet(gameID)
	if err != nil {
		return nil, err
	}
	collection, err := s.CollectionGet(gameID, collectionID)
	if err != nil {
		return nil, err
	}

	info, err := s.db.MoveFolder(oldName, newName, game.ID, collection.ID)
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
func (s *DB) DeckUpdate(gameID, collectionID, name, description, image string) (*entity.DeckInfo, error) {
	game, err := s.GameGet(gameID)
	if err != nil {
		return nil, err
	}
	collection, err := s.CollectionGet(gameID, collectionID)
	if err != nil {
		return nil, err
	}

	info, err := s.db.UpdateFolder(name, &commonInfo{
		Description: fsentry_types.QS(description),
		Image:       fsentry_types.QS(image),
	}, game.ID, collection.ID)
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
func (s *DB) DeckDelete(gameID, collectionID, name string) error {
	game, err := s.GameGet(gameID)
	if err != nil {
		return err
	}
	collection, err := s.CollectionGet(gameID, collectionID)
	if err != nil {
		return err
	}

	err = s.db.RemoveFolder(name, game.ID, collection.ID)
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
func (s *DB) DeckImageCreate(gameID, collectionID, name string, data []byte) error {
	game, err := s.GameGet(gameID)
	if err != nil {
		return err
	}
	collection, err := s.CollectionGet(gameID, collectionID)
	if err != nil {
		return err
	}
	deck, err := s.DeckGet(gameID, collectionID, name)
	if err != nil {
		return err
	}

	err = s.db.CreateBinary("image", data, game.ID, collection.ID, deck.ID)
	if err != nil {
		if errors.Is(err, fsentry_error.ErrorExist) {
			return er.DeckImageExist.AddMessage(err.Error())
		} else {
			return er.InternalError.AddMessage(err.Error())
		}
	}
	return nil
}
func (s *DB) DeckImageGet(gameID, collectionID, name string) ([]byte, error) {
	game, err := s.GameGet(gameID)
	if err != nil {
		return nil, err
	}
	collection, err := s.CollectionGet(gameID, collectionID)
	if err != nil {
		return nil, err
	}
	deck, err := s.DeckGet(gameID, collectionID, name)
	if err != nil {
		return nil, err
	}

	data, err := s.db.GetBinary("image", game.ID, collection.ID, deck.ID)
	if err != nil {
		if errors.Is(err, fsentry_error.ErrorNotExist) {
			return nil, er.DeckImageNotExists.AddMessage(err.Error())
		} else {
			return nil, er.InternalError.AddMessage(err.Error())
		}
	}
	return data, nil
}
func (s *DB) DeckImageDelete(gameID, collectionID, name string) error {
	game, err := s.GameGet(gameID)
	if err != nil {
		return err
	}
	collection, err := s.CollectionGet(gameID, collectionID)
	if err != nil {
		return err
	}
	deck, err := s.DeckGet(gameID, collectionID, name)
	if err != nil {
		return err
	}

	err = s.db.RemoveBinary("image", game.ID, collection.ID, deck.ID)
	if err != nil {
		if errors.Is(err, fsentry_error.ErrorNotExist) {
			return er.DeckImageNotExists.AddMessage(err.Error())
		} else {
			return er.InternalError.AddMessage(err.Error())
		}
	}
	return nil
}

func (s *DB) CardCreate(gameID, collectionID, deckID, name, description, image string, variables map[string]string, count int) (*entity.CardInfo, error) {
	game, err := s.GameGet(gameID)
	if err != nil {
		return nil, err
	}
	collection, err := s.CollectionGet(gameID, collectionID)
	if err != nil {
		return nil, err
	}
	deck, err := s.DeckGet(gameID, collectionID, deckID)
	if err != nil {
		return nil, err
	}

	list, err := s.rawCardList(gameID, collectionID, deckID)
	if err != nil {
		return nil, err
	}

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
	_, err = s.db.UpdateFolder("cards", list, game.ID, collection.ID, deck.ID)
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
func (s *DB) CardGet(gameID, collectionID, deckID string, cardID int64) (*entity.CardInfo, error) {
	list, err := s.rawCardList(gameID, collectionID, deckID)
	if err != nil {
		return nil, err
	}

	card, ok := list[cardID]
	if !ok {
		return nil, er.CardNotExists.HTTP(http.StatusBadRequest)
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
func (s *DB) CardList(gameID, collectionID, deckID string) ([]*entity.CardInfo, error) {
	list, err := s.rawCardList(gameID, collectionID, deckID)
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
func (s *DB) CardUpdate(gameID, collectionID, deckID string, cardID int64, name, description, image string, variables map[string]string, count int) (*entity.CardInfo, error) {
	game, err := s.GameGet(gameID)
	if err != nil {
		return nil, err
	}
	collection, err := s.CollectionGet(gameID, collectionID)
	if err != nil {
		return nil, err
	}
	deck, err := s.DeckGet(gameID, collectionID, deckID)
	if err != nil {
		return nil, err
	}

	list, err := s.rawCardList(gameID, collectionID, deckID)
	if err != nil {
		return nil, err
	}

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
	_, err = s.db.UpdateFolder("cards", list, game.ID, collection.ID, deck.ID)
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
func (s *DB) CardDelete(gameID, collectionID, deckID string, cardID int64) error {
	game, err := s.GameGet(gameID)
	if err != nil {
		return err
	}
	collection, err := s.CollectionGet(gameID, collectionID)
	if err != nil {
		return err
	}
	deck, err := s.DeckGet(gameID, collectionID, deckID)
	if err != nil {
		return err
	}

	list, err := s.rawCardList(gameID, collectionID, deckID)
	if err != nil {
		return err
	}

	if _, ok := list[cardID]; !ok {
		return er.CardNotExists
	}

	delete(list, cardID)

	// Writing an array of cards to a file again
	_, err = s.db.UpdateFolder("cards", list, game.ID, collection.ID, deck.ID)
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
func (s *DB) rawCardList(gameID, collectionID, deckID string) (map[int64]*cardInfo, error) {
	game, err := s.GameGet(gameID)
	if err != nil {
		return nil, err
	}
	collection, err := s.CollectionGet(gameID, collectionID)
	if err != nil {
		return nil, err
	}
	deck, err := s.DeckGet(gameID, collectionID, deckID)
	if err != nil {
		return nil, err
	}

	// Get all the cards
	info, err := s.db.GetFolder("cards", game.ID, collection.ID, deck.ID)
	if err != nil {
		return nil, er.InternalError.AddMessage(err.Error())
	}

	// Parsing an array of cards from json
	var list map[int64]*cardInfo
	err = json.Unmarshal(info.Data, &list)
	if err != nil {
		return nil, er.InternalError.AddMessage(err.Error())
	}

	if list == nil {
		list = make(map[int64]*cardInfo)
	}
	return list, nil
}
func (s *DB) CardImageCreate(gameID, collectionID, deckID string, cardID int64, data []byte) error {
	game, err := s.GameGet(gameID)
	if err != nil {
		return err
	}
	collection, err := s.CollectionGet(gameID, collectionID)
	if err != nil {
		return err
	}
	deck, err := s.DeckGet(gameID, collectionID, deckID)
	if err != nil {
		return err
	}
	card, err := s.CardGet(gameID, collectionID, deckID, cardID)
	if err != nil {
		return err
	}

	err = s.db.CreateBinary(fmt.Sprintf("%d", card.ID), data, game.ID, collection.ID, deck.ID, "cards")
	if err != nil {
		if errors.Is(err, fsentry_error.ErrorExist) {
			return er.CardImageExist.AddMessage(err.Error())
		} else {
			return er.InternalError.AddMessage(err.Error())
		}
	}
	return nil
}
func (s *DB) CardImageGet(gameID, collectionID, deckID string, cardID int64) ([]byte, error) {
	game, err := s.GameGet(gameID)
	if err != nil {
		return nil, err
	}
	collection, err := s.CollectionGet(gameID, collectionID)
	if err != nil {
		return nil, err
	}
	deck, err := s.DeckGet(gameID, collectionID, deckID)
	if err != nil {
		return nil, err
	}
	card, err := s.CardGet(gameID, collectionID, deckID, cardID)
	if err != nil {
		return nil, err
	}

	data, err := s.db.GetBinary(fmt.Sprintf("%d", card.ID), game.ID, collection.ID, deck.ID, "cards")
	if err != nil {
		if errors.Is(err, fsentry_error.ErrorNotExist) {
			return nil, er.CardImageNotExists.AddMessage(err.Error())
		} else {
			return nil, er.InternalError.AddMessage(err.Error())
		}
	}
	return data, nil
}
func (s *DB) CardImageDelete(gameID, collectionID, deckID string, cardID int64) error {
	game, err := s.GameGet(gameID)
	if err != nil {
		return err
	}
	collection, err := s.CollectionGet(gameID, collectionID)
	if err != nil {
		return err
	}
	deck, err := s.DeckGet(gameID, collectionID, deckID)
	if err != nil {
		return err
	}
	card, err := s.CardGet(gameID, collectionID, deckID, cardID)
	if err != nil {
		return err
	}

	err = s.db.RemoveBinary(fmt.Sprintf("%d", card.ID), game.ID, collection.ID, deck.ID, "cards")
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
