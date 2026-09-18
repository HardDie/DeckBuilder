package deck

import (
	"errors"
	"net/http"
	"time"

	"github.com/HardDie/fsentry"

	"github.com/HardDie/DeckBuilder/internal/config"
	entitiesDeck "github.com/HardDie/DeckBuilder/internal/entities/deck"
	er "github.com/HardDie/DeckBuilder/internal/errors"
	"github.com/HardDie/DeckBuilder/internal/images"
	"github.com/HardDie/DeckBuilder/internal/logger"
	"github.com/HardDie/DeckBuilder/internal/network"
	repositoriesCollection "github.com/HardDie/DeckBuilder/internal/repositories/collection"
)

type deck struct {
	cfg        *config.Config
	db         *fsentry.DB
	gamesPath  string
	collection repositoriesCollection.Collection
}

func New(cfg *config.Config, db *fsentry.DB, c repositoriesCollection.Collection) Deck {
	return &deck{
		cfg:        cfg,
		db:         db,
		gamesPath:  "games",
		collection: c,
	}
}

func (r *deck) Create(gameID, collectionID string, req CreateRequest) (*entitiesDeck.Deck, error) {
	d, err := r.create(gameID, collectionID, req)
	if err != nil {
		return nil, err
	}

	if d.Image == "" && req.ImageFile == nil {
		return d, nil
	}

	if d.Image != "" {
		err = r.createImage(gameID, collectionID, d.ID, d.Image)
		if err != nil {
			logger.Warn.Println("Unable to load image. The deck will be saved without an image.", err.Error())
		}
	} else if req.ImageFile != nil {
		err = r.createImageFromByte(gameID, collectionID, d.ID, req.ImageFile)
		if err != nil {
			logger.Warn.Println("Invalid image. The deck will be saved without an image.", err.Error())
		}
	}

	return d, nil
}

func (r *deck) GetByID(gameID, collectionID, deckID string) (*entitiesDeck.Deck, error) {
	return r.get(gameID, collectionID, deckID)
}

func (r *deck) GetAll(gameID, collectionID string) ([]*entitiesDeck.Deck, error) {
	return r.list(gameID, collectionID)
}

func (r *deck) Update(gameID, collectionID, deckID string, req UpdateRequest) (*entitiesDeck.Deck, error) {
	oldDeck, err := r.get(gameID, collectionID, deckID)
	if err != nil {
		return nil, err
	}

	var newDeck *entitiesDeck.Deck
	if oldDeck.Name != req.Name {
		newDeck, err = r.move(gameID, collectionID, oldDeck.Name, req.Name)
		if err != nil {
			return nil, err
		}
	}

	if oldDeck.Description != req.Description ||
		oldDeck.Image != req.Image ||
		req.ImageFile != nil {
		newDeck, err = r.update(gameID, collectionID, updateRequest{
			Name:        req.Name,
			Description: req.Description,
			Image:       req.Image,
		})
		if err != nil {
			return nil, err
		}
	}

	if newDeck == nil {
		newDeck = oldDeck
	}

	if newDeck.Image == oldDeck.Image && req.ImageFile == nil {
		return newDeck, nil
	}

	if data, _, _ := r.GetImage(gameID, collectionID, newDeck.ID); data != nil {
		err = r.imageDelete(gameID, collectionID, deckID)
		if err != nil {
			return nil, err
		}
	}

	if newDeck.Image == "" && req.ImageFile == nil {
		return newDeck, nil
	}

	if newDeck.Image != "" {
		err = r.createImage(gameID, collectionID, newDeck.ID, newDeck.Image)
		if err != nil {
			logger.Warn.Println("Unable to load image. The deck will be saved without an image.", err.Error())
		}
	} else if req.ImageFile != nil {
		err = r.createImageFromByte(gameID, collectionID, newDeck.ID, req.ImageFile)
		if err != nil {
			logger.Warn.Println("Invalid image. The deck will be saved without an image.", err.Error())
		}
	}

	return newDeck, nil
}

func (r *deck) DeleteByID(gameID, collectionID, deckID string) error {
	return r.delete(gameID, collectionID, deckID)
}

func (r *deck) GetImage(gameID, collectionID, deckID string) ([]byte, string, error) {
	data, err := r.imageGet(gameID, collectionID, deckID)
	if err != nil {
		return nil, "", err
	}

	imgType, err := images.ValidateImage(data)
	if err != nil {
		return nil, "", err
	}

	return data, imgType, nil
}

func (r *deck) GetAllDecksInGame(gameID string) ([]*entitiesDeck.Deck, error) {
	listCollections, err := r.collection.GetAll(gameID)
	if err != nil {
		return make([]*entitiesDeck.Deck, 0), err
	}

	uniqueDecks := make(map[string]struct{})

	decks := make([]*entitiesDeck.Deck, 0)
	for _, collection := range listCollections {
		collectionDecks, err := r.GetAll(gameID, collection.ID)
		if err != nil {
			return make([]*entitiesDeck.Deck, 0), err
		}

		for _, d := range collectionDecks {
			if _, ok := uniqueDecks[d.Name+d.Image]; ok {
				continue
			}
			uniqueDecks[d.Name+d.Image] = struct{}{}
			decks = append(decks, d)
		}
	}
	return decks, nil
}

func (r *deck) createImage(gameID, collectionID, deckID, imageURL string) error {
	imageBytes, err := network.DownloadBytes(imageURL)
	if err != nil {
		return err
	}

	return r.createImageFromByte(gameID, collectionID, deckID, imageBytes)
}

func (r *deck) createImageFromByte(gameID, collectionID, deckID string, data []byte) error {
	_, err := images.ValidateImage(data)
	if err != nil {
		return err
	}

	return r.imageCreate(gameID, collectionID, deckID, data)
}

func (r *deck) getCollection(gameID, collectionID string) (string, error) {
	gameInfo, err := r.db.GetFolder[any](gameID, r.gamesPath)
	if err != nil {
		if errors.Is(err, fsentry.ErrNotExist) {
			return "", er.GameNotExists.AddMessage(err.Error()).HTTP(http.StatusBadRequest)
		} else if errors.Is(err, fsentry.ErrBadName) {
			return "", er.BadName
		} else {
			return "", er.InternalError.AddMessage(err.Error())
		}
	}

	info, err := r.db.GetFolder[any](collectionID, r.gamesPath, gameInfo.ID)
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

func (r *deck) create(gameID, collectionID string, req CreateRequest) (*entitiesDeck.Deck, error) {
	collectionIDResolved, err := r.getCollection(gameID, collectionID)
	if err != nil {
		return nil, err
	}

	info, err := r.db.CreateFolder(req.Name, model{
		Description: fsentry.QuotedString(req.Description),
		Image:       fsentry.QuotedString(req.Image),
	}, r.gamesPath, gameID, collectionIDResolved)
	if err != nil {
		if errors.Is(err, fsentry.ErrExist) {
			return nil, er.DeckExist
		} else if errors.Is(err, fsentry.ErrBadName) {
			return nil, er.BadName
		} else {
			return nil, er.InternalError.AddMessage(err.Error())
		}
	}

	_, err = r.db.CreateFolder[any]("cards", nil, r.gamesPath, gameID, collectionIDResolved, info.ID)
	if err != nil {
		return nil, er.InternalError.AddMessage(err.Error())
	}

	return r.toEntity(info, gameID, collectionID), nil
}

func (r *deck) get(gameID, collectionID, name string) (*entitiesDeck.Deck, error) {
	collectionIDResolved, err := r.getCollection(gameID, collectionID)
	if err != nil {
		return nil, err
	}

	info, err := r.db.GetFolder[model](name, r.gamesPath, gameID, collectionIDResolved)
	if err != nil {
		if errors.Is(err, fsentry.ErrNotExist) {
			return nil, er.DeckNotExists.AddMessage(err.Error()).HTTP(http.StatusBadRequest)
		} else if errors.Is(err, fsentry.ErrBadName) {
			return nil, er.BadName
		} else {
			return nil, er.InternalError.AddMessage(err.Error())
		}
	}

	return r.toEntity(info, gameID, collectionID), nil
}

func (r *deck) list(gameID, collectionID string) ([]*entitiesDeck.Deck, error) {
	collectionIDResolved, err := r.getCollection(gameID, collectionID)
	if err != nil {
		return nil, err
	}

	list, err := r.db.List(r.gamesPath, gameID, collectionIDResolved)
	if err != nil {
		return nil, er.InternalError.AddMessage(err.Error())
	}

	var decks []*entitiesDeck.Deck
	for _, folder := range list.Folders {
		deck, err := r.get(gameID, collectionIDResolved, folder)
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

func (r *deck) move(gameID, collectionID, oldName, newName string) (*entitiesDeck.Deck, error) {
	collectionIDResolved, err := r.getCollection(gameID, collectionID)
	if err != nil {
		return nil, err
	}

	info, err := r.db.MoveFolder[model](oldName, newName, r.gamesPath, gameID, collectionIDResolved)
	if err != nil {
		if errors.Is(err, fsentry.ErrNotExist) {
			return nil, er.DeckNotExists.AddMessage(err.Error()).HTTP(http.StatusBadRequest)
		} else if errors.Is(err, fsentry.ErrBadName) {
			return nil, er.BadName
		} else {
			return nil, er.InternalError.AddMessage(err.Error())
		}
	}

	return r.toEntity(info, gameID, collectionID), nil
}

type updateRequest struct {
	Name        string
	Description string
	Image       string
}

func (r *deck) update(gameID, collectionID string, req updateRequest) (*entitiesDeck.Deck, error) {
	collectionIDResolved, err := r.getCollection(gameID, collectionID)
	if err != nil {
		return nil, err
	}

	info, err := r.db.UpdateFolder(req.Name, model{
		Description: fsentry.QuotedString(req.Description),
		Image:       fsentry.QuotedString(req.Image),
	}, r.gamesPath, gameID, collectionIDResolved)
	if err != nil {
		if errors.Is(err, fsentry.ErrNotExist) {
			return nil, er.DeckNotExists.AddMessage(err.Error()).HTTP(http.StatusBadRequest)
		} else if errors.Is(err, fsentry.ErrBadName) {
			return nil, er.BadName
		} else {
			return nil, er.InternalError.AddMessage(err.Error())
		}
	}

	return r.toEntity(info, gameID, collectionID), nil
}

func (r *deck) delete(gameID, collectionID, name string) error {
	collectionIDResolved, err := r.getCollection(gameID, collectionID)
	if err != nil {
		return err
	}

	err = r.db.RemoveFolder(name, r.gamesPath, gameID, collectionIDResolved)
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

func (r *deck) imageCreate(gameID, collectionID, deckID string, data []byte) error {
	deck, err := r.get(gameID, collectionID, deckID)
	if err != nil {
		return err
	}

	err = r.db.CreateBinary("image", data, r.gamesPath, gameID, collectionID, deck.ID)
	if err != nil {
		if errors.Is(err, fsentry.ErrExist) {
			return er.DeckImageExist.AddMessage(err.Error())
		} else {
			return er.InternalError.AddMessage(err.Error())
		}
	}
	return nil
}

func (r *deck) imageGet(gameID, collectionID, deckID string) ([]byte, error) {
	deck, err := r.get(gameID, collectionID, deckID)
	if err != nil {
		return nil, err
	}

	data, err := r.db.GetBinary("image", nil, r.gamesPath, gameID, collectionID, deck.ID)
	if err != nil {
		if errors.Is(err, fsentry.ErrNotExist) {
			return nil, er.DeckImageNotExists.AddMessage(err.Error())
		} else {
			return nil, er.InternalError.AddMessage(err.Error())
		}
	}
	return data, nil
}

func (r *deck) imageDelete(gameID, collectionID, deckID string) error {
	deck, err := r.get(gameID, collectionID, deckID)
	if err != nil {
		return err
	}

	err = r.db.RemoveBinary("image", r.gamesPath, gameID, collectionID, deck.ID)
	if err != nil {
		if errors.Is(err, fsentry.ErrNotExist) {
			return er.DeckImageNotExists.AddMessage(err.Error())
		} else {
			return er.InternalError.AddMessage(err.Error())
		}
	}
	return nil
}

func (r *deck) toEntity(info fsentry.FolderInfo[model], gameID, collectionID string) *entitiesDeck.Deck {
	createdAt, updatedAt := r.convertCreateUpdate(info.CreatedAt, info.UpdatedAt)
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

func (r *deck) convertCreateUpdate(createdAt, updatedAt time.Time) (time.Time, time.Time) {
	if createdAt.IsZero() {
		createdAt = time.Now()
	}
	if updatedAt.IsZero() {
		updatedAt = createdAt
	}
	return createdAt, updatedAt
}
