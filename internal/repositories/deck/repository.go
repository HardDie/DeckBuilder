package deck

import (
	"errors"

	"github.com/HardDie/fsentry"

	"github.com/HardDie/DeckBuilder/internal/config"
	entitiesDeck "github.com/HardDie/DeckBuilder/internal/entities/deck"
	er "github.com/HardDie/DeckBuilder/internal/errors"
	"github.com/HardDie/DeckBuilder/internal/images"
	"github.com/HardDie/DeckBuilder/internal/logger"
	"github.com/HardDie/DeckBuilder/internal/repositories"
	"github.com/HardDie/DeckBuilder/internal/utils"
)

type deck struct {
	cfg       *config.Config
	db        *fsentry.DB
	gamesPath string
}

func New(cfg *config.Config, db *fsentry.DB) Deck {
	return &deck{
		cfg:       cfg,
		db:        db,
		gamesPath: "games",
	}
}

func (r *deck) Create(gameID, collectionID string, req CreateRequest) (*entitiesDeck.Deck, error) {
	change, imageErr := repositories.ResolveImage("", req.Image, req.ImageFile)
	req.Image = change.URL

	item, err := r.create(gameID, collectionID, req)
	if err != nil {
		return nil, err
	}
	if change.Data != nil {
		if err = r.imageCreate(gameID, collectionID, item.ID, change.Data); err != nil {
			imageErr = err
		}
	}
	item.ImageError = imageErr
	return item, nil
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
	// Download and validate before any write, so a bad image keeps the old one.
	change, imageErr := repositories.ResolveImage(oldDeck.Image, req.Image, req.ImageFile)

	newDeck := oldDeck
	if oldDeck.Name != req.Name {
		newDeck, err = r.move(gameID, collectionID, oldDeck.Name, req.Name)
		if err != nil {
			return nil, err
		}
	}

	if oldDeck.Description != req.Description ||
		oldDeck.Image != change.URL ||
		change.Data != nil {
		newDeck, err = r.update(gameID, collectionID, updateRequest{
			Name:        req.Name,
			Description: req.Description,
			Image:       change.URL,
		})
		if err != nil {
			return nil, err
		}
	}

	if change.Data != nil || change.Clear {
		err = r.imageDelete(gameID, collectionID, newDeck.ID)
		if err != nil && !errors.Is(err, er.DeckImageNotExists) {
			return nil, err
		}
	}
	if change.Data != nil {
		if err = r.imageCreate(gameID, collectionID, newDeck.ID, change.Data); err != nil {
			imageErr = err
		}
	}
	newDeck.ImageError = imageErr
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

	imgType, err := images.ImageType(data)
	if err != nil {
		return nil, "", err
	}

	return data, imgType, nil
}

func (r *deck) GetAllDecksInGame(gameID string) ([]*entitiesDeck.Deck, error) {
	list, err := r.db.List(r.gamesPath, gameID)
	if err != nil {
		if mapped := er.MissingAncestor(err); mapped != nil {
			return make([]*entitiesDeck.Deck, 0), mapped
		}
		return make([]*entitiesDeck.Deck, 0), er.InternalError.AddMessage(err.Error())
	}

	uniqueDecks := make(map[string]struct{})

	decks := make([]*entitiesDeck.Deck, 0)
	for _, folder := range list.Folders {
		collectionDecks, err := r.GetAll(gameID, folder)
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

func (r *deck) create(gameID, collectionID string, req CreateRequest) (*entitiesDeck.Deck, error) {
	info, err := r.db.CreateFolder(req.Name, model{
		Description: fsentry.QuotedString(req.Description),
		Image:       fsentry.QuotedString(req.Image),
	}, r.gamesPath, gameID, collectionID)
	if err != nil {
		if mapped := er.MissingAncestor(err); mapped != nil {
			return nil, mapped
		}
		return nil, repositories.MapFsentry(err, repositories.FsentrySentinels{
			Exist: er.DeckExist,
		})
	}

	_, err = r.db.CreateFolder[any]("cards", nil, r.gamesPath, gameID, collectionID, info.ID)
	if err != nil {
		if mapped := er.MissingAncestor(err); mapped != nil {
			return nil, mapped
		}
		return nil, er.InternalError.AddMessage(err.Error())
	}

	return r.toEntity(info, gameID, collectionID), nil
}

func (r *deck) get(gameID, collectionID, name string) (*entitiesDeck.Deck, error) {
	info, err := r.db.GetFolder[model](name, r.gamesPath, gameID, collectionID)
	if err != nil {
		if mapped := er.MissingAncestor(err); mapped != nil {
			return nil, mapped
		}
		return nil, repositories.MapFsentry(err, repositories.FsentrySentinels{
			NotExist: er.DeckNotExists,
			Message:  true,
		})
	}

	return r.toEntity(info, gameID, collectionID), nil
}

func (r *deck) list(gameID, collectionID string) ([]*entitiesDeck.Deck, error) {
	list, err := r.db.List(r.gamesPath, gameID, collectionID)
	if err != nil {
		if mapped := er.MissingAncestor(err); mapped != nil {
			return nil, mapped
		}
		return nil, er.InternalError.AddMessage(err.Error())
	}

	var decks []*entitiesDeck.Deck
	for _, folder := range list.Folders {
		deck, err := r.get(gameID, collectionID, folder)
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
	info, err := r.db.MoveFolder[model](oldName, newName, r.gamesPath, gameID, collectionID)
	if err != nil {
		if mapped := er.MissingAncestor(err); mapped != nil {
			return nil, mapped
		}
		return nil, repositories.MapFsentry(err, repositories.FsentrySentinels{
			NotExist: er.DeckNotExists,
			Message:  true,
		})
	}

	return r.toEntity(info, gameID, collectionID), nil
}

type updateRequest struct {
	Name        string
	Description string
	Image       string
}

func (r *deck) update(gameID, collectionID string, req updateRequest) (*entitiesDeck.Deck, error) {
	info, err := r.db.UpdateFolder(req.Name, model{
		Description: fsentry.QuotedString(req.Description),
		Image:       fsentry.QuotedString(req.Image),
	}, r.gamesPath, gameID, collectionID)
	if err != nil {
		if mapped := er.MissingAncestor(err); mapped != nil {
			return nil, mapped
		}
		return nil, repositories.MapFsentry(err, repositories.FsentrySentinels{
			NotExist: er.DeckNotExists,
			Message:  true,
		})
	}

	return r.toEntity(info, gameID, collectionID), nil
}

func (r *deck) delete(gameID, collectionID, name string) error {
	err := r.db.RemoveFolder(name, r.gamesPath, gameID, collectionID)
	if err != nil {
		if mapped := er.MissingAncestor(err); mapped != nil {
			return mapped
		}
		return repositories.MapFsentry(err, repositories.FsentrySentinels{
			NotExist: er.DeckNotExists,
			Message:  true,
		})
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
		if mapped := er.MissingAncestor(err); mapped != nil {
			return mapped
		}
		return repositories.MapFsentry(err, repositories.FsentrySentinels{
			Exist:   er.DeckImageExist,
			Message: true,
		})
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
		if mapped := er.MissingAncestor(err); mapped != nil {
			return nil, mapped
		}
		return nil, repositories.MapFsentry(err, repositories.FsentrySentinels{
			NotExist: er.DeckImageNotExists,
			Message:  true,
		})
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
		if mapped := er.MissingAncestor(err); mapped != nil {
			return mapped
		}
		return repositories.MapFsentry(err, repositories.FsentrySentinels{
			NotExist: er.DeckImageNotExists,
			Message:  true,
		})
	}
	return nil
}

func (r *deck) toEntity(info fsentry.FolderInfo[model], gameID, collectionID string) *entitiesDeck.Deck {
	createdAt, updatedAt := utils.NormalizeTimestamps(info.CreatedAt, info.UpdatedAt)
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
