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
	"github.com/HardDie/fsentry/pkg/fsentry_error"
	"github.com/HardDie/fsentry/pkg/fsentry_types"

	dbDeck "github.com/HardDie/DeckBuilder/internal/db/deck"
	entitiesCard "github.com/HardDie/DeckBuilder/internal/entities/card"
	er "github.com/HardDie/DeckBuilder/internal/errors"
	"github.com/HardDie/DeckBuilder/internal/utils"
)

type card struct {
	db        fsentry.IFSEntry
	gamesPath string

	deck dbDeck.Deck
}

func New(db fsentry.IFSEntry, deck dbDeck.Deck) Card {
	return &card{
		db:        db,
		gamesPath: "games",

		deck: deck,
	}
}

func (d *card) Create(ctx context.Context, gameID, collectionID, deckID, name, description, image string, variables map[string]string, count int) (*entitiesCard.Card, error) {
	ctx, list, err := d.rawCardList(ctx, gameID, collectionID, deckID)
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
	cardInfo := &model{
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
	_, err = d.db.UpdateFolder("cards", list, d.gamesPath, gameID, collectionID, deckID)
	if err != nil {
		if errors.Is(err, fsentry_error.ErrorNotExist) {
			return nil, er.CardNotExists.AddMessage(err.Error())
		} else if errors.Is(err, fsentry_error.ErrorBadName) {
			return nil, er.BadName
		} else {
			return nil, er.InternalError.AddMessage(err.Error())
		}
	}

	createdAt, updatedAt := d.convertCreateUpdate(cardInfo.CreatedAt, cardInfo.UpdatedAt)
	return &entitiesCard.Card{
		ID:          cardInfo.ID,
		Name:        cardInfo.Name.String(),
		Description: cardInfo.Description.String(),
		Image:       cardInfo.Image.String(),
		Variables:   convertMapQuotedString(cardInfo.Variables),
		Count:       cardInfo.Count,
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,

		GameID:       gameID,
		CollectionID: collectionID,
		DeckID:       deckID,
	}, nil
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

	createdAt, updatedAt := d.convertCreateUpdate(card.CreatedAt, card.UpdatedAt)
	return &entitiesCard.Card{
		ID:          card.ID,
		Name:        card.Name.String(),
		Description: card.Description.String(),
		Image:       card.Image.String(),
		Variables:   convertMapQuotedString(card.Variables),
		Count:       card.Count,
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,

		GameID:       gameID,
		CollectionID: collectionID,
		DeckID:       deckID,
	}, nil
}
func (d *card) List(ctx context.Context, gameID, collectionID, deckID string) ([]*entitiesCard.Card, error) {
	ctx, list, err := d.rawCardList(ctx, gameID, collectionID, deckID)
	if err != nil {
		return nil, err
	}

	var cards []*entitiesCard.Card
	for _, item := range list {
		createdAt, updatedAt := d.convertCreateUpdate(item.CreatedAt, item.UpdatedAt)
		cards = append(cards, &entitiesCard.Card{
			ID:          item.ID,
			Name:        item.Name.String(),
			Description: item.Description.String(),
			Image:       item.Image.String(),
			Variables:   convertMapQuotedString(item.Variables),
			Count:       item.Count,
			CreatedAt:   createdAt,
			UpdatedAt:   updatedAt,

			GameID:       gameID,
			CollectionID: collectionID,
			DeckID:       deckID,
		})
	}
	return cards, nil
}
func (d *card) Update(ctx context.Context, gameID, collectionID, deckID string, cardID int64, name, description, image string, variables map[string]string, count int) (*entitiesCard.Card, error) {
	ctx, list, err := d.rawCardList(ctx, gameID, collectionID, deckID)
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
	_, err = d.db.UpdateFolder("cards", list, d.gamesPath, gameID, collectionID, deckID)
	if err != nil {
		if errors.Is(err, fsentry_error.ErrorNotExist) {
			return nil, er.CardNotExists.AddMessage(err.Error())
		} else if errors.Is(err, fsentry_error.ErrorBadName) {
			return nil, er.BadName
		} else {
			return nil, er.InternalError.AddMessage(err.Error())
		}
	}

	createdAt, updatedAt := d.convertCreateUpdate(card.CreatedAt, card.UpdatedAt)
	return &entitiesCard.Card{
		ID:          card.ID,
		Name:        card.Name.String(),
		Description: card.Description.String(),
		Image:       card.Image.String(),
		Variables:   convertMapQuotedString(card.Variables),
		Count:       card.Count,
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,

		GameID:       gameID,
		CollectionID: collectionID,
		DeckID:       deckID,
	}, nil
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

	// Writing an array of cards to a file again
	_, err = d.db.UpdateFolder("cards", list, d.gamesPath, gameID, collectionID, deckID)
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
func (d *card) ImageCreate(ctx context.Context, gameID, collectionID, deckID string, cardID int64, data []byte) error {
	card, err := d.Get(ctx, gameID, collectionID, deckID, cardID)
	if err != nil {
		return err
	}

	err = d.db.CreateBinary(fmt.Sprintf("%d", card.ID), data, d.gamesPath, gameID, collectionID, deckID, "cards")
	if err != nil {
		if errors.Is(err, fsentry_error.ErrorExist) {
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

	data, err := d.db.GetBinary(fmt.Sprintf("%d", card.ID), d.gamesPath, gameID, collectionID, deckID, "cards")
	if err != nil {
		if errors.Is(err, fsentry_error.ErrorNotExist) {
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
		if errors.Is(err, fsentry_error.ErrorNotExist) {
			return er.CardImageNotExists.AddMessage(err.Error())
		} else {
			return er.InternalError.AddMessage(err.Error())
		}
	}
	return nil
}

func (d *card) rawCardList(ctx context.Context, gameID, collectionID, deckID string) (context.Context, map[int64]*model, error) {
	deck, err := d.deck.Get(ctx, gameID, collectionID, deckID)
	if err != nil {
		return ctx, nil, err
	}

	// Get all the cards
	info, err := d.db.GetFolder("cards", d.gamesPath, gameID, collectionID, deck.ID)
	if err != nil {
		return ctx, nil, er.InternalError.AddMessage(err.Error())
	}

	// Parsing an array of cards from json
	var list map[int64]*model
	err = json.Unmarshal(info.Data, &list)
	if err != nil {
		return ctx, nil, er.InternalError.AddMessage(err.Error())
	}

	if list == nil {
		list = make(map[int64]*model)
	}
	return ctx, list, nil
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

func (d *card) convertCreateUpdate(createdAt, updatedAt *time.Time) (time.Time, time.Time) {
	if createdAt == nil {
		createdAt = utils.Allocate(time.Now())
	}
	if updatedAt == nil {
		updatedAt = createdAt
	}
	return *createdAt, *updatedAt
}
