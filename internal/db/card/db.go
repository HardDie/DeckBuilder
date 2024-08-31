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

func (d *card) Create(ctx context.Context, gameID, collectionID, deckID, name, description, image string, variables map[string]string, count int) (*CardInfo, error) {
	ctx, list, err := d.rawCardList(ctx, gameID, collectionID, deckID)
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
	_, err = d.db.UpdateFolder("cards", list, d.gamesPath, ctxGameID, ctxCollectionID, ctxDeckID)
	if err != nil {
		if errors.Is(err, fsentry_error.ErrorNotExist) {
			return nil, er.CardNotExists.AddMessage(err.Error())
		} else if errors.Is(err, fsentry_error.ErrorBadName) {
			return nil, er.BadName
		} else {
			return nil, er.InternalError.AddMessage(err.Error())
		}
	}

	return &CardInfo{
		ID:          cardInfo.ID,
		Name:        cardInfo.Name.String(),
		Description: cardInfo.Description.String(),
		Image:       cardInfo.Image.String(),
		Variables:   convertMapQuotedString(cardInfo.Variables),
		Count:       cardInfo.Count,
		CreatedAt:   cardInfo.CreatedAt,
		UpdatedAt:   nil,

		GameID:       gameID,
		CollectionID: collectionID,
		DeckID:       deckID,
	}, nil
}
func (d *card) Get(ctx context.Context, gameID, collectionID, deckID string, cardID int64) (context.Context, *CardInfo, error) {
	ctx, list, err := d.rawCardList(ctx, gameID, collectionID, deckID)
	if err != nil {
		return ctx, nil, err
	}

	card, ok := list[cardID]
	if !ok {
		return ctx, nil, er.CardNotExists.HTTP(http.StatusBadRequest)
	}

	return ctx, &CardInfo{
		ID:          card.ID,
		Name:        card.Name.String(),
		Description: card.Description.String(),
		Image:       card.Image.String(),
		Variables:   convertMapQuotedString(card.Variables),
		Count:       card.Count,
		CreatedAt:   card.CreatedAt,
		UpdatedAt:   card.UpdatedAt,

		GameID:       gameID,
		CollectionID: collectionID,
		DeckID:       deckID,
	}, nil
}
func (d *card) List(ctx context.Context, gameID, collectionID, deckID string) ([]*CardInfo, error) {
	ctx, list, err := d.rawCardList(ctx, gameID, collectionID, deckID)
	if err != nil {
		return nil, err
	}

	var cards []*CardInfo
	for _, item := range list {
		cards = append(cards, &CardInfo{
			ID:          item.ID,
			Name:        item.Name.String(),
			Description: item.Description.String(),
			Image:       item.Image.String(),
			Variables:   convertMapQuotedString(item.Variables),
			Count:       item.Count,
			CreatedAt:   item.CreatedAt,
			UpdatedAt:   item.UpdatedAt,

			GameID:       gameID,
			CollectionID: collectionID,
			DeckID:       deckID,
		})
	}
	return cards, nil
}
func (d *card) Update(ctx context.Context, gameID, collectionID, deckID string, cardID int64, name, description, image string, variables map[string]string, count int) (*CardInfo, error) {
	ctx, list, err := d.rawCardList(ctx, gameID, collectionID, deckID)
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
	_, err = d.db.UpdateFolder("cards", list, d.gamesPath, ctxGameID, ctxCollectionID, ctxDeckID)
	if err != nil {
		if errors.Is(err, fsentry_error.ErrorNotExist) {
			return nil, er.CardNotExists.AddMessage(err.Error())
		} else if errors.Is(err, fsentry_error.ErrorBadName) {
			return nil, er.BadName
		} else {
			return nil, er.InternalError.AddMessage(err.Error())
		}
	}

	return &CardInfo{
		ID:          card.ID,
		Name:        card.Name.String(),
		Description: card.Description.String(),
		Image:       card.Image.String(),
		Variables:   convertMapQuotedString(card.Variables),
		Count:       card.Count,
		CreatedAt:   card.CreatedAt,
		UpdatedAt:   card.UpdatedAt,

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
	ctxGameID := ctx.Value("gameID").(string)
	ctxCollectionID := ctx.Value("collectionID").(string)
	ctxDeckID := ctx.Value("deckID").(string)

	if _, ok := list[cardID]; !ok {
		return er.CardNotExists
	}

	delete(list, cardID)

	// Writing an array of cards to a file again
	_, err = d.db.UpdateFolder("cards", list, d.gamesPath, ctxGameID, ctxCollectionID, ctxDeckID)
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
	ctx, card, err := d.Get(ctx, gameID, collectionID, deckID, cardID)
	if err != nil {
		return err
	}
	ctxGameID := ctx.Value("gameID").(string)
	ctxCollectionID := ctx.Value("collectionID").(string)
	ctxDeckID := ctx.Value("deckID").(string)

	err = d.db.CreateBinary(fmt.Sprintf("%d", card.ID), data, d.gamesPath, ctxGameID, ctxCollectionID, ctxDeckID, "cards")
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
	ctx, card, err := d.Get(ctx, gameID, collectionID, deckID, cardID)
	if err != nil {
		return nil, err
	}
	ctxGameID := ctx.Value("gameID").(string)
	ctxCollectionID := ctx.Value("collectionID").(string)
	ctxDeckID := ctx.Value("deckID").(string)

	data, err := d.db.GetBinary(fmt.Sprintf("%d", card.ID), d.gamesPath, ctxGameID, ctxCollectionID, ctxDeckID, "cards")
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
	ctx, card, err := d.Get(ctx, gameID, collectionID, deckID, cardID)
	if err != nil {
		return err
	}
	ctxGameID := ctx.Value("gameID").(string)
	ctxCollectionID := ctx.Value("collectionID").(string)
	ctxDeckID := ctx.Value("deckID").(string)

	err = d.db.RemoveBinary(fmt.Sprintf("%d", card.ID), d.gamesPath, ctxGameID, ctxCollectionID, ctxDeckID, "cards")
	if err != nil {
		if errors.Is(err, fsentry_error.ErrorNotExist) {
			return er.CardImageNotExists.AddMessage(err.Error())
		} else {
			return er.InternalError.AddMessage(err.Error())
		}
	}
	return nil
}

func (d *card) rawCardList(ctx context.Context, gameID, collectionID, deckID string) (context.Context, map[int64]*cardInfo, error) {
	ctx, deck, err := d.deck.Get(ctx, gameID, collectionID, deckID)
	if err != nil {
		return ctx, nil, err
	}
	ctxGameID := ctx.Value("gameID").(string)
	ctxCollectionID := ctx.Value("collectionID").(string)

	// Get all the cards
	info, err := d.db.GetFolder("cards", d.gamesPath, ctxGameID, ctxCollectionID, deck.ID)
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
