package deck

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/HardDie/fsentry"
	"github.com/HardDie/fsentry/pkg/fsentry_error"
	"github.com/HardDie/fsentry/pkg/fsentry_types"

	dbCollection "github.com/HardDie/DeckBuilder/internal/db/collection"
	entitiesDeck "github.com/HardDie/DeckBuilder/internal/entities/deck"
	er "github.com/HardDie/DeckBuilder/internal/errors"
	"github.com/HardDie/DeckBuilder/internal/logger"
	"github.com/HardDie/DeckBuilder/internal/utils"
)

type deck struct {
	db        fsentry.IFSEntry
	gamesPath string

	collection dbCollection.Collection
}

func New(db fsentry.IFSEntry, collection dbCollection.Collection) Deck {
	return &deck{
		db:        db,
		gamesPath: "games",

		collection: collection,
	}
}

func (d *deck) Create(ctx context.Context, req CreateRequest) (*entitiesDeck.Deck, error) {
	collection, err := d.collection.Get(ctx, req.GameID, req.CollectionID)
	if err != nil {
		return nil, err
	}

	info, err := d.db.CreateFolder(req.Name, model{
		Description: fsentry_types.QS(req.Description),
		Image:       fsentry_types.QS(req.Image),
	}, d.gamesPath, req.GameID, collection.ID)
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
	_, err = d.db.CreateFolder("cards", nil, d.gamesPath, req.GameID, collection.ID, info.Id)
	if err != nil {
		return nil, er.InternalError.AddMessage(err.Error())
	}

	createdAt, updatedAt := d.convertCreateUpdate(info.CreatedAt, info.UpdatedAt)
	return &entitiesDeck.Deck{
		ID:          info.Id,
		Name:        info.Name.String(),
		Description: req.Description,
		Image:       req.Image,
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,

		GameID:       req.GameID,
		CollectionID: req.CollectionID,
	}, nil
}
func (d *deck) Get(ctx context.Context, gameID, collectionID, name string) (*entitiesDeck.Deck, error) {
	collection, err := d.collection.Get(ctx, gameID, collectionID)
	if err != nil {
		return nil, err
	}

	info, err := d.db.GetFolder(name, d.gamesPath, gameID, collection.ID)
	if err != nil {
		if errors.Is(err, fsentry_error.ErrorNotExist) {
			return nil, er.DeckNotExists.AddMessage(err.Error()).HTTP(http.StatusBadRequest)
		} else if errors.Is(err, fsentry_error.ErrorBadName) {
			return nil, er.BadName
		} else {
			return nil, er.InternalError.AddMessage(err.Error())
		}
	}

	var dInfo model
	err = json.Unmarshal(info.Data, &dInfo)
	if err != nil {
		return nil, er.InternalError.AddMessage(err.Error())
	}

	createdAt, updatedAt := d.convertCreateUpdate(info.CreatedAt, info.UpdatedAt)
	return &entitiesDeck.Deck{
		ID:          info.Id,
		Name:        info.Name.String(),
		Description: dInfo.Description.String(),
		Image:       dInfo.Image.String(),
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,

		GameID:       gameID,
		CollectionID: collectionID,
	}, nil
}
func (d *deck) List(ctx context.Context, gameID, collectionID string) ([]*entitiesDeck.Deck, error) {
	collection, err := d.collection.Get(ctx, gameID, collectionID)
	if err != nil {
		return nil, err
	}

	list, err := d.db.List(d.gamesPath, gameID, collection.ID)
	if err != nil {
		return nil, er.InternalError.AddMessage(err.Error())
	}

	var decks []*entitiesDeck.Deck
	for _, folder := range list.Folders {
		deck, err := d.Get(ctx, gameID, collection.ID, folder)
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
func (d *deck) Move(ctx context.Context, gameID, collectionID, oldName, newName string) (*entitiesDeck.Deck, error) {
	collection, err := d.collection.Get(ctx, gameID, collectionID)
	if err != nil {
		return nil, err
	}

	info, err := d.db.MoveFolder(oldName, newName, d.gamesPath, gameID, collection.ID)
	if err != nil {
		if errors.Is(err, fsentry_error.ErrorNotExist) {
			return nil, er.DeckNotExists.AddMessage(err.Error()).HTTP(http.StatusBadRequest)
		} else if errors.Is(err, fsentry_error.ErrorBadName) {
			return nil, er.BadName
		} else {
			return nil, er.InternalError.AddMessage(err.Error())
		}
	}

	var dInfo model
	err = json.Unmarshal(info.Data, &dInfo)
	if err != nil {
		return nil, er.InternalError.AddMessage(err.Error())
	}

	createdAt, updatedAt := d.convertCreateUpdate(info.CreatedAt, info.UpdatedAt)
	return &entitiesDeck.Deck{
		ID:          info.Id,
		Name:        info.Name.String(),
		Description: dInfo.Description.String(),
		Image:       dInfo.Image.String(),
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,

		GameID:       gameID,
		CollectionID: collectionID,
	}, nil
}
func (d *deck) Update(ctx context.Context, req UpdateRequest) (*entitiesDeck.Deck, error) {
	collection, err := d.collection.Get(ctx, req.GameID, req.CollectionID)
	if err != nil {
		return nil, err
	}

	info, err := d.db.UpdateFolder(req.Name, model{
		Description: fsentry_types.QS(req.Description),
		Image:       fsentry_types.QS(req.Image),
	}, d.gamesPath, req.GameID, collection.ID)
	if err != nil {
		if errors.Is(err, fsentry_error.ErrorNotExist) {
			return nil, er.DeckNotExists.AddMessage(err.Error()).HTTP(http.StatusBadRequest)
		} else if errors.Is(err, fsentry_error.ErrorBadName) {
			return nil, er.BadName
		} else {
			return nil, er.InternalError.AddMessage(err.Error())
		}
	}

	var dInfo model
	err = json.Unmarshal(info.Data, &dInfo)
	if err != nil {
		return nil, er.InternalError.AddMessage(err.Error())
	}

	createdAt, updatedAt := d.convertCreateUpdate(info.CreatedAt, info.UpdatedAt)
	return &entitiesDeck.Deck{
		ID:          info.Id,
		Name:        info.Name.String(),
		Description: dInfo.Description.String(),
		Image:       dInfo.Image.String(),
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,

		GameID:       req.GameID,
		CollectionID: req.CollectionID,
	}, nil
}
func (d *deck) Delete(ctx context.Context, gameID, collectionID, name string) error {
	collection, err := d.collection.Get(ctx, gameID, collectionID)
	if err != nil {
		return err
	}

	err = d.db.RemoveFolder(name, d.gamesPath, gameID, collection.ID)
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
func (d *deck) ImageCreate(ctx context.Context, gameID, collectionID, deckID string, data []byte) error {
	deck, err := d.Get(ctx, gameID, collectionID, deckID)
	if err != nil {
		return err
	}

	err = d.db.CreateBinary("image", data, d.gamesPath, gameID, collectionID, deck.ID)
	if err != nil {
		if errors.Is(err, fsentry_error.ErrorExist) {
			return er.DeckImageExist.AddMessage(err.Error())
		} else {
			return er.InternalError.AddMessage(err.Error())
		}
	}
	return nil
}
func (d *deck) ImageGet(ctx context.Context, gameID, collectionID, deckID string) ([]byte, error) {
	deck, err := d.Get(ctx, gameID, collectionID, deckID)
	if err != nil {
		return nil, err
	}

	data, err := d.db.GetBinary("image", d.gamesPath, gameID, collectionID, deck.ID)
	if err != nil {
		if errors.Is(err, fsentry_error.ErrorNotExist) {
			return nil, er.DeckImageNotExists.AddMessage(err.Error())
		} else {
			return nil, er.InternalError.AddMessage(err.Error())
		}
	}
	return data, nil
}
func (d *deck) ImageDelete(ctx context.Context, gameID, collectionID, deckID string) error {
	deck, err := d.Get(ctx, gameID, collectionID, deckID)
	if err != nil {
		return err
	}

	err = d.db.RemoveBinary("image", d.gamesPath, gameID, collectionID, deck.ID)
	if err != nil {
		if errors.Is(err, fsentry_error.ErrorNotExist) {
			return er.DeckImageNotExists.AddMessage(err.Error())
		} else {
			return er.InternalError.AddMessage(err.Error())
		}
	}
	return nil
}

func (d *deck) convertCreateUpdate(createdAt, updatedAt *time.Time) (time.Time, time.Time) {
	if createdAt == nil {
		createdAt = utils.Allocate(time.Now())
	}
	if updatedAt == nil {
		updatedAt = createdAt
	}
	return *createdAt, *updatedAt
}
