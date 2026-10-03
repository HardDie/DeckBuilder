package card

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/HardDie/fsentry"

	"github.com/HardDie/DeckBuilder/internal/config"
	entitiesCard "github.com/HardDie/DeckBuilder/internal/entities/card"
	er "github.com/HardDie/DeckBuilder/internal/errors"
	"github.com/HardDie/DeckBuilder/internal/images"
	"github.com/HardDie/DeckBuilder/internal/repositories"
	"github.com/HardDie/DeckBuilder/internal/utils"
)

type cardList map[int64]*model

type card struct {
	cfg       *config.Config
	db        *fsentry.DB
	gamesPath string

	// mu serializes the read-modify-write of a deck's cards list.
	// fsentry locks each call, not the read and write pair.
	mu sync.Mutex
}

func New(cfg *config.Config, db *fsentry.DB) Card {
	return &card{
		cfg:       cfg,
		db:        db,
		gamesPath: "games",
	}
}

func (r *card) Create(gameID, collectionID, deckID string, req CreateRequest) (*entitiesCard.Card, error) {
	change, imageErr := repositories.ResolveImage("", req.Image, req.ImageFile)
	req.Image = change.URL

	c, err := r.create(gameID, collectionID, deckID, req)
	if err != nil {
		return nil, err
	}
	if change.Data != nil {
		if err = r.imageCreate(gameID, collectionID, deckID, c.ID, change.Data); err != nil {
			imageErr = err
		}
	}
	c.ImageError = imageErr
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
	// Download and validate before any write, so a bad image keeps the old one.
	change, imageErr := repositories.ResolveImage(oldCard.Image, req.Image, req.ImageFile)
	req.Image = change.URL

	newCard := oldCard
	if oldCard.Name != req.Name ||
		oldCard.Description != req.Description ||
		oldCard.Image != change.URL ||
		change.Data != nil ||
		oldCard.Count != req.Count ||
		!utils.CompareMaps(oldCard.Variables, req.Variables) {
		newCard, err = r.update(gameID, collectionID, deckID, cardID, req)
		if err != nil {
			return nil, err
		}
	}

	if change.Data != nil || change.Clear {
		err = r.imageDelete(gameID, collectionID, deckID, newCard.ID)
		if err != nil && !errors.Is(err, er.CardImageNotExists) {
			return nil, err
		}
	}
	if change.Data != nil {
		if err = r.imageCreate(gameID, collectionID, deckID, newCard.ID, change.Data); err != nil {
			imageErr = err
		}
	}
	newCard.ImageError = imageErr
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

	imgType, err := images.ImageType(data)
	if err != nil {
		return nil, "", err
	}

	return data, imgType, nil
}

func (r *card) create(gameID, collectionID, deckID string, req CreateRequest) (*entitiesCard.Card, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

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

	now := time.Now()
	cardInfo := &model{
		ID:          maxID,
		Name:        fsentry.QuotedString(req.Name),
		Description: fsentry.QuotedString(req.Description),
		Image:       fsentry.QuotedString(req.Image),
		Variables:   convertMapString(req.Variables),
		Count:       req.Count,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	list[cardInfo.ID] = cardInfo

	_, err = r.db.UpdateFolder("cards", list, r.gamesPath, gameID, collectionID, deckID)
	if err != nil {
		return nil, repositories.MapFsentry(err, repositories.FsentrySentinels{
			NotExist: er.CardNotExists,
			Message:  true,
		})
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
	r.mu.Lock()
	defer r.mu.Unlock()

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
	card.UpdatedAt = time.Now()

	list[card.ID] = card

	_, err = r.db.UpdateFolder("cards", list, r.gamesPath, gameID, collectionID, deckID)
	if err != nil {
		return nil, repositories.MapFsentry(err, repositories.FsentrySentinels{
			NotExist: er.CardNotExists,
			Message:  true,
		})
	}

	return r.toEntity(card, gameID, collectionID, deckID), nil
}

func (r *card) delete(gameID, collectionID, deckID string, cardID int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()

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
		return repositories.MapFsentry(err, repositories.FsentrySentinels{
			NotExist: er.CardNotExists,
			Message:  true,
		})
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
		return repositories.MapFsentry(err, repositories.FsentrySentinels{
			Exist:   er.CardImageExist,
			Message: true,
		})
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
		return nil, repositories.MapFsentry(err, repositories.FsentrySentinels{
			NotExist: er.CardImageNotExists,
			Message:  true,
		})
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
		return repositories.MapFsentry(err, repositories.FsentrySentinels{
			NotExist: er.CardImageNotExists,
			Message:  true,
		})
	}
	return nil
}

func (r *card) rawCardList(gameID, collectionID, deckID string) (cardList, error) {
	info, err := r.db.GetFolder[cardList]("cards", r.gamesPath, gameID, collectionID, deckID)
	if err != nil {
		return nil, repositories.MapFsentry(err, repositories.FsentrySentinels{})
	}

	list := info.Data
	if list == nil {
		list = make(cardList)
	}
	return list, nil
}

func (r *card) toEntity(item *model, gameID, collectionID, deckID string) *entitiesCard.Card {
	createdAt, updatedAt := utils.NormalizeTimestamps(item.CreatedAt, item.UpdatedAt)
	// Old or imported data may hold a count below 1; read it as 1.
	count := max(item.Count, 1)
	return &entitiesCard.Card{
		ID:           item.ID,
		Name:         item.Name.String(),
		Description:  item.Description.String(),
		Image:        item.Image.String(),
		Variables:    convertMapQuotedString(item.Variables),
		Count:        count,
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
