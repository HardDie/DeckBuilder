package deck

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/HardDie/fsentry"
	"github.com/HardDie/fsentry/pkg/fsentry_error"
	"github.com/HardDie/fsentry/pkg/fsentry_types"

	dbCollection "github.com/HardDie/DeckBuilder/internal/db/collection"
	dbCommon "github.com/HardDie/DeckBuilder/internal/db/common"
	"github.com/HardDie/DeckBuilder/internal/entity"
	er "github.com/HardDie/DeckBuilder/internal/errors"
	"github.com/HardDie/DeckBuilder/internal/logger"
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

func (d *deck) Create(ctx context.Context, gameID, collectionID, name, description, image string) (*entity.DeckInfo, error) {
	ctx, collection, err := d.collection.Get(ctx, gameID, collectionID)
	if err != nil {
		return nil, err
	}
	ctxGameID := ctx.Value("gameID").(string)

	info, err := d.db.CreateFolder(name, &dbCommon.Info{
		Description: fsentry_types.QS(description),
		Image:       fsentry_types.QS(image),
	}, d.gamesPath, ctxGameID, collection.ID)
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
	_, err = d.db.CreateFolder("cards", nil, d.gamesPath, ctxGameID, collection.ID, info.Id)
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
func (d *deck) Get(ctx context.Context, gameID, collectionID, name string) (context.Context, *entity.DeckInfo, error) {
	ctx, collection, err := d.collection.Get(ctx, gameID, collectionID)
	if err != nil {
		return ctx, nil, err
	}
	ctxGameID := ctx.Value("gameID").(string)

	info, err := d.db.GetFolder(name, d.gamesPath, ctxGameID, collection.ID)
	if err != nil {
		if errors.Is(err, fsentry_error.ErrorNotExist) {
			return ctx, nil, er.DeckNotExists.AddMessage(err.Error()).HTTP(http.StatusBadRequest)
		} else if errors.Is(err, fsentry_error.ErrorBadName) {
			return ctx, nil, er.BadName
		} else {
			return ctx, nil, er.InternalError.AddMessage(err.Error())
		}
	}
	var dInfo dbCommon.Info

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
func (d *deck) List(ctx context.Context, gameID, collectionID string) ([]*entity.DeckInfo, error) {
	ctx, collection, err := d.collection.Get(ctx, gameID, collectionID)
	if err != nil {
		return nil, err
	}
	ctxGameID := ctx.Value("gameID").(string)

	list, err := d.db.List(d.gamesPath, ctxGameID, collection.ID)
	if err != nil {
		return nil, er.InternalError.AddMessage(err.Error())
	}

	var decks []*entity.DeckInfo
	for _, folder := range list.Folders {
		_, deck, err := d.Get(ctx, ctxGameID, collection.ID, folder)
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
func (d *deck) Move(ctx context.Context, gameID, collectionID, oldName, newName string) (*entity.DeckInfo, error) {
	ctx, collection, err := d.collection.Get(ctx, gameID, collectionID)
	if err != nil {
		return nil, err
	}
	ctxGameID := ctx.Value("gameID").(string)

	info, err := d.db.MoveFolder(oldName, newName, d.gamesPath, ctxGameID, collection.ID)
	if err != nil {
		if errors.Is(err, fsentry_error.ErrorNotExist) {
			return nil, er.DeckNotExists.AddMessage(err.Error()).HTTP(http.StatusBadRequest)
		} else if errors.Is(err, fsentry_error.ErrorBadName) {
			return nil, er.BadName
		} else {
			return nil, er.InternalError.AddMessage(err.Error())
		}
	}
	var dInfo dbCommon.Info

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
func (d *deck) Update(ctx context.Context, gameID, collectionID, name, description, image string) (*entity.DeckInfo, error) {
	ctx, collection, err := d.collection.Get(ctx, gameID, collectionID)
	if err != nil {
		return nil, err
	}
	ctxGameID := ctx.Value("gameID").(string)

	info, err := d.db.UpdateFolder(name, &dbCommon.Info{
		Description: fsentry_types.QS(description),
		Image:       fsentry_types.QS(image),
	}, d.gamesPath, ctxGameID, collection.ID)
	if err != nil {
		if errors.Is(err, fsentry_error.ErrorNotExist) {
			return nil, er.DeckNotExists.AddMessage(err.Error()).HTTP(http.StatusBadRequest)
		} else if errors.Is(err, fsentry_error.ErrorBadName) {
			return nil, er.BadName
		} else {
			return nil, er.InternalError.AddMessage(err.Error())
		}
	}
	var dInfo dbCommon.Info

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
func (d *deck) Delete(ctx context.Context, gameID, collectionID, name string) error {
	ctx, collection, err := d.collection.Get(ctx, gameID, collectionID)
	if err != nil {
		return err
	}
	ctxGameID := ctx.Value("gameID").(string)

	err = d.db.RemoveFolder(name, d.gamesPath, ctxGameID, collection.ID)
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
	ctx, deck, err := d.Get(ctx, gameID, collectionID, deckID)
	if err != nil {
		return err
	}
	ctxGameID := ctx.Value("gameID").(string)
	ctxCollectionID := ctx.Value("collectionID").(string)

	err = d.db.CreateBinary("image", data, d.gamesPath, ctxGameID, ctxCollectionID, deck.ID)
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
	ctx, deck, err := d.Get(ctx, gameID, collectionID, deckID)
	if err != nil {
		return nil, err
	}
	ctxGameID := ctx.Value("gameID").(string)
	ctxCollectionID := ctx.Value("collectionID").(string)

	data, err := d.db.GetBinary("image", d.gamesPath, ctxGameID, ctxCollectionID, deck.ID)
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
	ctx, deck, err := d.Get(ctx, gameID, collectionID, deckID)
	if err != nil {
		return err
	}
	ctxGameID := ctx.Value("gameID").(string)
	ctxCollectionID := ctx.Value("collectionID").(string)

	err = d.db.RemoveBinary("image", d.gamesPath, ctxGameID, ctxCollectionID, deck.ID)
	if err != nil {
		if errors.Is(err, fsentry_error.ErrorNotExist) {
			return er.DeckImageNotExists.AddMessage(err.Error())
		} else {
			return er.InternalError.AddMessage(err.Error())
		}
	}
	return nil
}
