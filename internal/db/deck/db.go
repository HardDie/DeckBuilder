package deck

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/HardDie/fsentry"

	entitiesDeck "github.com/HardDie/DeckBuilder/internal/entities/deck"
	er "github.com/HardDie/DeckBuilder/internal/errors"
	"github.com/HardDie/DeckBuilder/internal/logger"
)

type deck struct {
	db        *fsentry.DB
	gamesPath string
}

func New(db *fsentry.DB) Deck {
	return &deck{
		db:        db,
		gamesPath: "games",
	}
}

func (d *deck) getCollection(_ context.Context, gameID, collectionID string) (string, error) {
	gameInfo, err := d.db.GetFolder[any](gameID, d.gamesPath)
	if err != nil {
		if errors.Is(err, fsentry.ErrNotExist) {
			return "", er.GameNotExists.AddMessage(err.Error()).HTTP(http.StatusBadRequest)
		} else if errors.Is(err, fsentry.ErrBadName) {
			return "", er.BadName
		} else {
			return "", er.InternalError.AddMessage(err.Error())
		}
	}

	info, err := d.db.GetFolder[any](collectionID, d.gamesPath, gameInfo.ID)
	if err != nil {
		if errors.Is(err, fsentry.ErrNotExist) {
			return "", er.CollectionNotExists.AddMessage(err.Error()).HTTP(http.StatusBadRequest)
		} else if errors.Is(err, fsentry.ErrBadName) {
			return "", er.BadName
		} else {
			return "", er.InternalError.AddMessage(err.Error())
		}
	}
	return info.ID, nil
}

func (d *deck) Create(ctx context.Context, req CreateRequest) (*entitiesDeck.Deck, error) {
	collectionIDResolved, err := d.getCollection(ctx, req.GameID, req.CollectionID)
	if err != nil {
		return nil, err
	}

	info, err := d.db.CreateFolder(req.Name, model{
		Description: fsentry.QuotedString(req.Description),
		Image:       fsentry.QuotedString(req.Image),
	}, d.gamesPath, req.GameID, collectionIDResolved)
	if err != nil {
		if errors.Is(err, fsentry.ErrExist) {
			return nil, er.DeckExist
		} else if errors.Is(err, fsentry.ErrBadName) {
			return nil, er.BadName
		} else {
			return nil, er.InternalError.AddMessage(err.Error())
		}
	}

	_, err = d.db.CreateFolder[any]("cards", nil, d.gamesPath, req.GameID, collectionIDResolved, info.ID)
	if err != nil {
		return nil, er.InternalError.AddMessage(err.Error())
	}

	return d.toEntity(info, req.GameID, req.CollectionID), nil
}

func (d *deck) Get(ctx context.Context, gameID, collectionID, name string) (*entitiesDeck.Deck, error) {
	collectionIDResolved, err := d.getCollection(ctx, gameID, collectionID)
	if err != nil {
		return nil, err
	}

	info, err := d.db.GetFolder[model](name, d.gamesPath, gameID, collectionIDResolved)
	if err != nil {
		if errors.Is(err, fsentry.ErrNotExist) {
			return nil, er.DeckNotExists.AddMessage(err.Error()).HTTP(http.StatusBadRequest)
		} else if errors.Is(err, fsentry.ErrBadName) {
			return nil, er.BadName
		} else {
			return nil, er.InternalError.AddMessage(err.Error())
		}
	}

	return d.toEntity(info, gameID, collectionID), nil
}

func (d *deck) List(ctx context.Context, gameID, collectionID string) ([]*entitiesDeck.Deck, error) {
	collectionIDResolved, err := d.getCollection(ctx, gameID, collectionID)
	if err != nil {
		return nil, err
	}

	list, err := d.db.List(d.gamesPath, gameID, collectionIDResolved)
	if err != nil {
		return nil, er.InternalError.AddMessage(err.Error())
	}

	var decks []*entitiesDeck.Deck
	for _, folder := range list.Folders {
		deck, err := d.Get(ctx, gameID, collectionIDResolved, folder)
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
	collectionIDResolved, err := d.getCollection(ctx, gameID, collectionID)
	if err != nil {
		return nil, err
	}

	info, err := d.db.MoveFolder[model](oldName, newName, d.gamesPath, gameID, collectionIDResolved)
	if err != nil {
		if errors.Is(err, fsentry.ErrNotExist) {
			return nil, er.DeckNotExists.AddMessage(err.Error()).HTTP(http.StatusBadRequest)
		} else if errors.Is(err, fsentry.ErrBadName) {
			return nil, er.BadName
		} else {
			return nil, er.InternalError.AddMessage(err.Error())
		}
	}

	return d.toEntity(info, gameID, collectionID), nil
}

func (d *deck) Update(ctx context.Context, req UpdateRequest) (*entitiesDeck.Deck, error) {
	collectionIDResolved, err := d.getCollection(ctx, req.GameID, req.CollectionID)
	if err != nil {
		return nil, err
	}

	info, err := d.db.UpdateFolder(req.Name, model{
		Description: fsentry.QuotedString(req.Description),
		Image:       fsentry.QuotedString(req.Image),
	}, d.gamesPath, req.GameID, collectionIDResolved)
	if err != nil {
		if errors.Is(err, fsentry.ErrNotExist) {
			return nil, er.DeckNotExists.AddMessage(err.Error()).HTTP(http.StatusBadRequest)
		} else if errors.Is(err, fsentry.ErrBadName) {
			return nil, er.BadName
		} else {
			return nil, er.InternalError.AddMessage(err.Error())
		}
	}

	return d.toEntity(info, req.GameID, req.CollectionID), nil
}

func (d *deck) Delete(ctx context.Context, gameID, collectionID, name string) error {
	collectionIDResolved, err := d.getCollection(ctx, gameID, collectionID)
	if err != nil {
		return err
	}

	err = d.db.RemoveFolder(name, d.gamesPath, gameID, collectionIDResolved)
	if err != nil {
		if errors.Is(err, fsentry.ErrNotExist) {
			return er.DeckNotExists.AddMessage(err.Error()).HTTP(http.StatusBadRequest)
		} else if errors.Is(err, fsentry.ErrBadName) {
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
		if errors.Is(err, fsentry.ErrExist) {
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

	data, err := d.db.GetBinary("image", nil, d.gamesPath, gameID, collectionID, deck.ID)
	if err != nil {
		if errors.Is(err, fsentry.ErrNotExist) {
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
		if errors.Is(err, fsentry.ErrNotExist) {
			return er.DeckImageNotExists.AddMessage(err.Error())
		} else {
			return er.InternalError.AddMessage(err.Error())
		}
	}
	return nil
}

func (d *deck) toEntity(info fsentry.FolderInfo[model], gameID, collectionID string) *entitiesDeck.Deck {
	createdAt, updatedAt := d.convertCreateUpdate(info.CreatedAt, info.UpdatedAt)
	return &entitiesDeck.Deck{
		ID:           info.ID,
		Name:         info.Name,
		Description:  info.Data.Description.String(),
		Image:        info.Data.Image.String(),
		CreatedAt:    createdAt,
		UpdatedAt:    updatedAt,
		GameID:       gameID,
		CollectionID: collectionID,
	}
}

func (d *deck) convertCreateUpdate(createdAt, updatedAt time.Time) (time.Time, time.Time) {
	if createdAt.IsZero() {
		createdAt = time.Now()
	}
	if updatedAt.IsZero() {
		updatedAt = createdAt
	}
	return createdAt, updatedAt
}
