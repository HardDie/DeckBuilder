package card

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/HardDie/fsentry"

	"github.com/HardDie/DeckBuilder/internal/config"
	entitiesCard "github.com/HardDie/DeckBuilder/internal/entities/card"
	er "github.com/HardDie/DeckBuilder/internal/errors"
	"github.com/HardDie/DeckBuilder/internal/images"
	"github.com/HardDie/DeckBuilder/internal/logger"
	"github.com/HardDie/DeckBuilder/internal/network"
	"github.com/HardDie/DeckBuilder/internal/utils"
)

type cardList map[int64]*model

type card struct {
	cfg       *config.Config
	db        *fsentry.DB
	gamesPath string
}

func New(cfg *config.Config, db *fsentry.DB) Card {
	return &card{
		cfg:       cfg,
		db:        db,
		gamesPath: "games",
	}
}

func (r *card) Create(gameID, collectionID, deckID string, req CreateRequest) (*entitiesCard.Card, error) {
	c, err := r.create(gameID, collectionID, deckID, req)
	if err != nil {
		return nil, err
	}

	if c.Image == "" && req.ImageFile == nil {
		return c, nil
	}

	if c.Image != "" {
		err = r.createImage(gameID, collectionID, deckID, c.ID, c.Image)
		if err != nil {
			logger.Warn.Println("Unable to load image. The card will be saved without an image.", err.Error())
		}
	} else if req.ImageFile != nil {
		err = r.createImageFromByte(gameID, collectionID, deckID, c.ID, req.ImageFile)
		if err != nil {
			logger.Warn.Println("Invalid image. The card will be saved without an image.", err.Error())
		}
	}

	return c, nil
}

func (r *card) GetByID(gameID, collectionID, deckID string, cardID int64) (*entitiesCard.Card, error) {
	return r.get(gameID, collectionID, deckID, cardID)
}

func (r *card) GetAll(gameID, collectionID, deckID string) ([]*entitiesCard.Card, error) {
	return r.list(gameID, collectionID, deckID)
}

func (r *card) Update(gameID, collectionID, deckID string, cardID int64, req UpdateRequest) (*entitiesCard.Card, error) {
	oldCard, err := r.get(gameID, collectionID, deckID, cardID)
	if err != nil {
		return nil, err
	}

	var newCard *entitiesCard.Card
	if oldCard.Name != req.Name ||
		oldCard.Description != req.Description ||
		oldCard.Image != req.Image ||
		req.ImageFile != nil ||
		oldCard.Count != req.Count ||
		!utils.CompareMaps(oldCard.Variables, req.Variables) {
		newCard, err = r.update(gameID, collectionID, deckID, cardID, req)
		if err != nil {
			return nil, err
		}
	}

	if newCard == nil {
		newCard = oldCard
	}

	if newCard.Image == oldCard.Image && req.ImageFile == nil {
		return newCard, nil
	}

	if data, _, _ := r.GetImage(gameID, collectionID, deckID, newCard.ID); data != nil {
		err = r.imageDelete(gameID, collectionID, deckID, cardID)
		if err != nil {
			return nil, err
		}
	}

	if newCard.Image == "" && req.ImageFile == nil {
		return newCard, nil
	}

	if newCard.Image != "" {
		if err = r.createImage(gameID, collectionID, deckID, newCard.ID, newCard.Image); err != nil {
			logger.Warn.Println("Unable to load image. The card will be saved without an image.", err.Error())
		}
	} else if req.ImageFile != nil {
		err = r.createImageFromByte(gameID, collectionID, deckID, newCard.ID, req.ImageFile)
		if err != nil {
			logger.Warn.Println("Invalid image. The card will be saved without an image.", err.Error())
		}
	}

	return newCard, nil
}

func (r *card) DeleteByID(gameID, collectionID, deckID string, cardID int64) error {
	err := r.imageDelete(gameID, collectionID, deckID, cardID)
	if err != nil {
		if !errors.Is(err, er.CardImageNotExists) {
			return err
		}
	}
	return r.delete(gameID, collectionID, deckID, cardID)
}

func (r *card) GetImage(gameID, collectionID, deckID string, cardID int64) ([]byte, string, error) {
	data, err := r.imageGet(gameID, collectionID, deckID, cardID)
	if err != nil {
		return nil, "", err
	}

	imgType, err := images.ValidateImage(data)
	if err != nil {
		return nil, "", err
	}

	return data, imgType, nil
}

func (r *card) createImage(gameID, collectionID, deckID string, cardID int64, imageURL string) error {
	imageBytes, err := network.DownloadBytes(imageURL)
	if err != nil {
		return err
	}

	return r.createImageFromByte(gameID, collectionID, deckID, cardID, imageBytes)
}

func (r *card) createImageFromByte(gameID, collectionID, deckID string, cardID int64, data []byte) error {
	_, err := images.ValidateImage(data)
	if err != nil {
		return err
	}

	return r.imageCreate(gameID, collectionID, deckID, cardID, data)
}

func (r *card) create(gameID, collectionID, deckID string, req CreateRequest) (*entitiesCard.Card, error) {
	list, err := r.rawCardList(gameID, collectionID, deckID)
	if err != nil {
		return nil, err
	}

	maxID := int64(1)
	for _, card := range list {
		if card.ID >= maxID {
			maxID = card.ID + 1
		}
	}

	cardInfo := &model{
		ID:          maxID,
		Name:        fsentry.QuotedString(req.Name),
		Description: fsentry.QuotedString(req.Description),
		Image:       fsentry.QuotedString(req.Image),
		Variables:   convertMapString(req.Variables),
		Count:       req.Count,
		CreatedAt:   utils.Allocate(time.Now()),
		UpdatedAt:   nil,
	}

	list[cardInfo.ID] = cardInfo

	_, err = r.db.UpdateFolder("cards", list, r.gamesPath, gameID, collectionID, deckID)
	if err != nil {
		if errors.Is(err, fsentry.ErrNotExist) {
			return nil, er.CardNotExists.AddMessage(err.Error())
		} else if errors.Is(err, fsentry.ErrBadName) {
			return nil, er.BadName
		} else {
			return nil, er.InternalError.AddMessage(err.Error())
		}
	}

	return r.toEntity(cardInfo, gameID, collectionID, deckID), nil
}

func (r *card) get(gameID, collectionID, deckID string, cardID int64) (*entitiesCard.Card, error) {
	list, err := r.rawCardList(gameID, collectionID, deckID)
	if err != nil {
		return nil, err
	}

	card, ok := list[cardID]
	if !ok {
		return nil, er.CardNotExists.HTTP(http.StatusBadRequest)
	}

	return r.toEntity(card, gameID, collectionID, deckID), nil
}

func (r *card) list(gameID, collectionID, deckID string) ([]*entitiesCard.Card, error) {
	list, err := r.rawCardList(gameID, collectionID, deckID)
	if err != nil {
		return nil, err
	}

	var cards []*entitiesCard.Card
	for _, item := range list {
		cards = append(cards, r.toEntity(item, gameID, collectionID, deckID))
	}
	return cards, nil
}

func (r *card) update(gameID, collectionID, deckID string, cardID int64, req UpdateRequest) (*entitiesCard.Card, error) {
	list, err := r.rawCardList(gameID, collectionID, deckID)
	if err != nil {
		return nil, err
	}

	card, ok := list[cardID]
	if !ok {
		return nil, er.CardNotExists
	}

	card.Name = fsentry.QuotedString(req.Name)
	card.Description = fsentry.QuotedString(req.Description)
	card.Image = fsentry.QuotedString(req.Image)
	card.Variables = convertMapString(req.Variables)
	card.Count = req.Count
	card.UpdatedAt = utils.Allocate(time.Now())

	list[card.ID] = card

	_, err = r.db.UpdateFolder("cards", list, r.gamesPath, gameID, collectionID, deckID)
	if err != nil {
		if errors.Is(err, fsentry.ErrNotExist) {
			return nil, er.CardNotExists.AddMessage(err.Error())
		} else if errors.Is(err, fsentry.ErrBadName) {
			return nil, er.BadName
		} else {
			return nil, er.InternalError.AddMessage(err.Error())
		}
	}

	return r.toEntity(card, gameID, collectionID, deckID), nil
}

func (r *card) delete(gameID, collectionID, deckID string, cardID int64) error {
	list, err := r.rawCardList(gameID, collectionID, deckID)
	if err != nil {
		return err
	}

	if _, ok := list[cardID]; !ok {
		return er.CardNotExists
	}

	delete(list, cardID)

	_, err = r.db.UpdateFolder("cards", list, r.gamesPath, gameID, collectionID, deckID)
	if err != nil {
		if errors.Is(err, fsentry.ErrNotExist) {
			return er.CardNotExists.AddMessage(err.Error())
		} else if errors.Is(err, fsentry.ErrBadName) {
			return er.BadName
		} else {
			return er.InternalError.AddMessage(err.Error())
		}
	}

	return nil
}

func (r *card) imageCreate(gameID, collectionID, deckID string, cardID int64, data []byte) error {
	card, err := r.get(gameID, collectionID, deckID, cardID)
	if err != nil {
		return err
	}

	err = r.db.CreateBinary(fmt.Sprintf("%d", card.ID), data, r.gamesPath, gameID, collectionID, deckID, "cards")
	if err != nil {
		if errors.Is(err, fsentry.ErrExist) {
			return er.CardImageExist.AddMessage(err.Error())
		} else {
			return er.InternalError.AddMessage(err.Error())
		}
	}
	return nil
}

func (r *card) imageGet(gameID, collectionID, deckID string, cardID int64) ([]byte, error) {
	card, err := r.get(gameID, collectionID, deckID, cardID)
	if err != nil {
		return nil, err
	}

	data, err := r.db.GetBinary(fmt.Sprintf("%d", card.ID), nil, r.gamesPath, gameID, collectionID, deckID, "cards")
	if err != nil {
		if errors.Is(err, fsentry.ErrNotExist) {
			return nil, er.CardImageNotExists.AddMessage(err.Error())
		} else {
			return nil, er.InternalError.AddMessage(err.Error())
		}
	}
	return data, nil
}

func (r *card) imageDelete(gameID, collectionID, deckID string, cardID int64) error {
	card, err := r.get(gameID, collectionID, deckID, cardID)
	if err != nil {
		return err
	}

	err = r.db.RemoveBinary(fmt.Sprintf("%d", card.ID), r.gamesPath, gameID, collectionID, deckID, "cards")
	if err != nil {
		if errors.Is(err, fsentry.ErrNotExist) {
			return er.CardImageNotExists.AddMessage(err.Error())
		} else {
			return er.InternalError.AddMessage(err.Error())
		}
	}
	return nil
}

func (r *card) getDeck(gameID, collectionID, deckID string) (string, error) {
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

	collectionInfo, err := r.db.GetFolder[any](collectionID, r.gamesPath, gameInfo.ID)
	if err != nil {
		if errors.Is(err, fsentry.ErrNotExist) {
			return "", er.CollectionNotExists.AddMessage(err.Error()).HTTP(http.StatusBadRequest)
		} else if errors.Is(err, fsentry.ErrBadName) {
			return "", er.BadName
		} else {
			return "", er.InternalError.AddMessage(err.Error())
		}
	}

	info, err := r.db.GetFolder[any](deckID, r.gamesPath, gameInfo.ID, collectionInfo.ID)
	if err != nil {
		if errors.Is(err, fsentry.ErrNotExist) {
			return "", er.DeckNotExists.AddMessage(err.Error()).HTTP(http.StatusBadRequest)
		} else if errors.Is(err, fsentry.ErrBadName) {
			return "", er.BadName
		} else {
			return "", er.InternalError.AddMessage(err.Error())
		}
	}
	return info.ID, nil
}

func (r *card) rawCardList(gameID, collectionID, deckID string) (cardList, error) {
	deckIDResolved, err := r.getDeck(gameID, collectionID, deckID)
	if err != nil {
		return nil, err
	}

	info, err := r.db.GetFolder[cardList]("cards", r.gamesPath, gameID, collectionID, deckIDResolved)
	if err != nil {
		return nil, er.InternalError.AddMessage(err.Error())
	}

	list := info.Data
	if list == nil {
		list = make(cardList)
	}
	return list, nil
}

func (r *card) toEntity(item *model, gameID, collectionID, deckID string) *entitiesCard.Card {
	createdAt, updatedAt := r.convertCreateUpdate(item.CreatedAt, item.UpdatedAt)
	return &entitiesCard.Card{
		ID:           item.ID,
		Name:         item.Name.String(),
		Description:  item.Description.String(),
		Image:        item.Image.String(),
		Variables:    convertMapQuotedString(item.Variables),
		Count:        item.Count,
		CreatedAt:    createdAt,
		UpdatedAt:    updatedAt,
		GameID:       gameID,
		CollectionID: collectionID,
		DeckID:       deckID,
	}
}

func convertMapString(in map[string]string) map[string]fsentry.QuotedString {
	res := make(map[string]fsentry.QuotedString)
	for key, val := range in {
		keyJSON, _ := json.Marshal(strconv.Quote(key))
		res[string(keyJSON)] = fsentry.QuotedString(val)
	}
	return res
}

func convertMapQuotedString(in map[string]fsentry.QuotedString) map[string]string {
	res := make(map[string]string)
	for keyJSON, val := range in {
		var key string
		_ = json.Unmarshal([]byte(keyJSON), &key)
		key, _ = strconv.Unquote(key)
		res[key] = val.String()
	}
	return res
}

func (r *card) convertCreateUpdate(createdAt, updatedAt *time.Time) (time.Time, time.Time) {
	if createdAt == nil {
		createdAt = utils.Allocate(time.Now())
	}
	if updatedAt == nil {
		updatedAt = createdAt
	}
	return *createdAt, *updatedAt
}
