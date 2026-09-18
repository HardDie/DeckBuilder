package card

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/HardDie/fsentry"

	dbDeck "github.com/HardDie/DeckBuilder/internal/db/deck"
	entitiesCard "github.com/HardDie/DeckBuilder/internal/entities/card"
	er "github.com/HardDie/DeckBuilder/internal/errors"
	"github.com/HardDie/DeckBuilder/internal/utils"
)

type cardList map[int64]*model

type card struct {
	db        *fsentry.DB
	gamesPath string

	deck dbDeck.Deck
}

func New(db *fsentry.DB, deck dbDeck.Deck) Card {
	return &card{
		db:        db,
		gamesPath: "games",

		deck: deck,
	}
}

func (d *card) Create(ctx context.Context, req CreateRequest) (*entitiesCard.Card, error) {
	ctx, list, err := d.rawCardList(ctx, req.GameID, req.CollectionID, req.DeckID)
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

	_, err = d.db.UpdateFolder("cards", list, d.gamesPath, req.GameID, req.CollectionID, req.DeckID)
	if err != nil {
		if errors.Is(err, fsentry.ErrNotExist) {
			return nil, er.CardNotExists.AddMessage(err.Error())
		} else if errors.Is(err, fsentry.ErrBadName) {
			return nil, er.BadName
		} else {
			return nil, er.InternalError.AddMessage(err.Error())
		}
	}

	return d.toEntity(cardInfo, req.GameID, req.CollectionID, req.DeckID), nil
}

func (d *card) Get(ctx context.Context, gameID, collectionID, deckID string, cardID int64) (*entitiesCard.Card, error) {
	ctx, list, err := d.rawCardList(ctx, gameID, collectionID, deckID)
	if err != nil {
		return nil, err
	}

	card, ok := list[cardID]
	if !ok {
		return nil, er.CardNotExists.HTTP(http.StatusBadRequest)
	}

	return d.toEntity(card, gameID, collectionID, deckID), nil
}

func (d *card) List(ctx context.Context, gameID, collectionID, deckID string) ([]*entitiesCard.Card, error) {
	ctx, list, err := d.rawCardList(ctx, gameID, collectionID, deckID)
	if err != nil {
		return nil, err
	}

	var cards []*entitiesCard.Card
	for _, item := range list {
		cards = append(cards, d.toEntity(item, gameID, collectionID, deckID))
	}
	return cards, nil
}

func (d *card) Update(ctx context.Context, req UpdateRequest) (*entitiesCard.Card, error) {
	ctx, list, err := d.rawCardList(ctx, req.GameID, req.CollectionID, req.DeckID)
	if err != nil {
		return nil, err
	}

	card, ok := list[req.CardID]
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

	_, err = d.db.UpdateFolder("cards", list, d.gamesPath, req.GameID, req.CollectionID, req.DeckID)
	if err != nil {
		if errors.Is(err, fsentry.ErrNotExist) {
			return nil, er.CardNotExists.AddMessage(err.Error())
		} else if errors.Is(err, fsentry.ErrBadName) {
			return nil, er.BadName
		} else {
			return nil, er.InternalError.AddMessage(err.Error())
		}
	}

	return d.toEntity(card, req.GameID, req.CollectionID, req.DeckID), nil
}

func (d *card) Delete(ctx context.Context, gameID, collectionID, deckID string, cardID int64) error {
	ctx, list, err := d.rawCardList(ctx, gameID, collectionID, deckID)
	if err != nil {
		return err
	}

	if _, ok := list[cardID]; !ok {
		return er.CardNotExists
	}

	delete(list, cardID)

	_, err = d.db.UpdateFolder("cards", list, d.gamesPath, gameID, collectionID, deckID)
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

func (d *card) ImageCreate(ctx context.Context, gameID, collectionID, deckID string, cardID int64, data []byte) error {
	card, err := d.Get(ctx, gameID, collectionID, deckID, cardID)
	if err != nil {
		return err
	}

	err = d.db.CreateBinary(fmt.Sprintf("%d", card.ID), data, d.gamesPath, gameID, collectionID, deckID, "cards")
	if err != nil {
		if errors.Is(err, fsentry.ErrExist) {
			return er.CardImageExist.AddMessage(err.Error())
		} else {
			return er.InternalError.AddMessage(err.Error())
		}
	}
	return nil
}

func (d *card) ImageGet(ctx context.Context, gameID, collectionID, deckID string, cardID int64) ([]byte, error) {
	card, err := d.Get(ctx, gameID, collectionID, deckID, cardID)
	if err != nil {
		return nil, err
	}

	data, err := d.db.GetBinary(fmt.Sprintf("%d", card.ID), nil, d.gamesPath, gameID, collectionID, deckID, "cards")
	if err != nil {
		if errors.Is(err, fsentry.ErrNotExist) {
			return nil, er.CardImageNotExists.AddMessage(err.Error())
		} else {
			return nil, er.InternalError.AddMessage(err.Error())
		}
	}
	return data, nil
}

func (d *card) ImageDelete(ctx context.Context, gameID, collectionID, deckID string, cardID int64) error {
	card, err := d.Get(ctx, gameID, collectionID, deckID, cardID)
	if err != nil {
		return err
	}

	err = d.db.RemoveBinary(fmt.Sprintf("%d", card.ID), d.gamesPath, gameID, collectionID, deckID, "cards")
	if err != nil {
		if errors.Is(err, fsentry.ErrNotExist) {
			return er.CardImageNotExists.AddMessage(err.Error())
		} else {
			return er.InternalError.AddMessage(err.Error())
		}
	}
	return nil
}

func (d *card) rawCardList(ctx context.Context, gameID, collectionID, deckID string) (context.Context, cardList, error) {
	deck, err := d.deck.Get(ctx, gameID, collectionID, deckID)
	if err != nil {
		return ctx, nil, err
	}

	info, err := d.db.GetFolder[cardList]("cards", d.gamesPath, gameID, collectionID, deck.ID)
	if err != nil {
		return ctx, nil, er.InternalError.AddMessage(err.Error())
	}

	list := info.Data
	if list == nil {
		list = make(cardList)
	}
	return ctx, list, nil
}

func (d *card) toEntity(item *model, gameID, collectionID, deckID string) *entitiesCard.Card {
	createdAt, updatedAt := d.convertCreateUpdate(item.CreatedAt, item.UpdatedAt)
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

func (d *card) convertCreateUpdate(createdAt, updatedAt *time.Time) (time.Time, time.Time) {
	if createdAt == nil {
		createdAt = utils.Allocate(time.Now())
	}
	if updatedAt == nil {
		updatedAt = createdAt
	}
	return *createdAt, *updatedAt
}
