package card

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/HardDie/fsentry"
	"github.com/stretchr/testify/assert"

	"github.com/HardDie/DeckBuilder/internal/config"
	entitiesCard "github.com/HardDie/DeckBuilder/internal/entities/card"
	er "github.com/HardDie/DeckBuilder/internal/errors"
	repositoriesCore "github.com/HardDie/DeckBuilder/internal/repositories/core"
)

var (
	img = []byte("some_image")
)

type cardEnv struct {
	db   *fsentry.DB
	card *card
}

func initCard(t testing.TB, name string) cardEnv {
	dir, err := os.MkdirTemp("", name)
	if err != nil {
		t.Fatal("error creating temp dir", err)
	}
	t.Cleanup(func() {
		e := os.RemoveAll(dir)
		if e != nil {
			t.Fatal("error RemoveAll", e)
		}
	})

	cfg := config.Get(false, "")
	cfg.SetDataPath(dir)

	db := fsentry.New(cfg.Data, fsentry.WithPretty(), fsentry.WithNoLockFile())
	if err = db.Init(); err != nil {
		t.Fatal("error init db", err)
	}

	core := repositoriesCore.New(db)
	err = core.Init()
	if err != nil {
		t.Fatal("error init core", err)
	}
	t.Cleanup(func() {
		e := core.Drop()
		if e != nil {
			t.Fatal("error drop core", e)
		}
	})

	return cardEnv{
		db:   db,
		card: New(cfg, db).(*card),
	}
}

func (e cardEnv) createGame(t testing.TB, name string) string {
	t.Helper()
	info, err := e.db.CreateFolder[any](name, nil, "games")
	if err != nil {
		t.Fatal("error create game", err)
	}
	return info.ID
}

func (e cardEnv) createCollection(t testing.TB, gameID, name string) string {
	t.Helper()
	info, err := e.db.CreateFolder[any](name, nil, "games", gameID)
	if err != nil {
		t.Fatal("error create collection", err)
	}
	return info.ID
}

func (e cardEnv) createDeck(t testing.TB, gameID, collectionID, name string) string {
	t.Helper()
	info, err := e.db.CreateFolder[any](name, nil, "games", gameID, collectionID)
	if err != nil {
		t.Fatal("error create deck", err)
	}
	_, err = e.db.CreateFolder[any]("cards", nil, "games", gameID, collectionID, info.ID)
	if err != nil {
		t.Fatal("error create cards folder", err)
	}
	return info.ID
}

func (e cardEnv) createParents(t testing.TB, gameName, collectionName, deckName string) (string, string, string) {
	t.Helper()
	gameID := e.createGame(t, gameName)
	collectionID := e.createCollection(t, gameID, collectionName)
	deckID := e.createDeck(t, gameID, collectionID, deckName)
	return gameID, collectionID, deckID
}

func TestCardCreate(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		wait := &entitiesCard.Card{
			ID:          1,
			Name:        "success",
			Description: "descrption",
			Image:       "https://some.url/image",
			Variables:   map[string]string{"atk": "2", "name key": "value"},
			Count:       3,
		}

		e := initCard(t, "card_create__success")
		gameID, collectionID, deckID := e.createParents(t, "parent_game", "parent_collection", "parent_deck")
		wait.GameID = gameID
		wait.CollectionID = collectionID
		wait.DeckID = deckID
		got, err := e.card.create(wait.GameID, wait.CollectionID, wait.DeckID, CreateRequest{Name: wait.Name,
			Description: wait.Description,
			Image:       wait.Image,
			Variables:   wait.Variables,
			Count:       wait.Count,
		})
		assert.NoError(t, err)
		wait.CreatedAt = got.CreatedAt
		wait.UpdatedAt = got.UpdatedAt
		assert.Equal(t, wait.CreatedAt, wait.UpdatedAt)
		assert.Equal(t, wait, got)
	})

	t.Run("ids_increment_and_reuse_gap", func(t *testing.T) {
		e := initCard(t, "card_create__ids")
		gameID, collectionID, deckID := e.createParents(t, "parent_game", "parent_collection", "parent_deck")
		first, err := e.card.create(gameID, collectionID, deckID, CreateRequest{Name: "one"})
		assert.NoError(t, err)
		assert.Equal(t, int64(1), first.ID)
		second, err := e.card.create(gameID, collectionID, deckID, CreateRequest{Name: "two"})
		assert.NoError(t, err)
		assert.Equal(t, int64(2), second.ID)
		err = e.card.delete(gameID, collectionID, deckID, second.ID)
		assert.NoError(t, err)
		third, err := e.card.create(gameID, collectionID, deckID, CreateRequest{Name: "three"})
		assert.NoError(t, err)
		assert.Equal(t, int64(2), third.ID)
	})

	t.Run("isolated_by_deck", func(t *testing.T) {
		e := initCard(t, "card_create__isolated_by_deck")
		gameID, collectionID, deckA := e.createParents(t, "parent_game", "parent_collection", "deck_a")
		deckB := e.createDeck(t, gameID, collectionID, "deck_b")
		cardA, err := e.card.create(gameID, collectionID, deckA, CreateRequest{Name: "shared"})
		assert.NoError(t, err)
		cardB, err := e.card.create(gameID, collectionID, deckB, CreateRequest{Name: "shared"})
		assert.NoError(t, err)
		assert.Equal(t, int64(1), cardA.ID)
		assert.Equal(t, int64(1), cardB.ID)
	})

	t.Run("empty_variables", func(t *testing.T) {
		e := initCard(t, "card_create__empty_variables")
		gameID, collectionID, deckID := e.createParents(t, "parent_game", "parent_collection", "parent_deck")
		got, err := e.card.create(gameID, collectionID, deckID, CreateRequest{Name: "empty"})
		assert.NoError(t, err)
		assert.Equal(t, map[string]string{}, got.Variables)
	})

	t.Run("deck_not_exist", func(t *testing.T) {
		e := initCard(t, "card_create__deck_not_exist")
		gameID := e.createGame(t, "parent_game")
		collectionID := e.createCollection(t, gameID, "parent_collection")
		_, err := e.card.create(gameID, collectionID, "missing", CreateRequest{Name: "ok"})
		assert.ErrorIs(t, err, er.DeckNotExists)
	})

	t.Run("collection_not_exist", func(t *testing.T) {
		e := initCard(t, "card_create__collection_not_exist")
		gameID := e.createGame(t, "parent_game")
		_, err := e.card.create(gameID, "missing", "ok", CreateRequest{Name: "ok"})
		assert.ErrorIs(t, err, er.CollectionNotExists)
	})

	t.Run("game_not_exist", func(t *testing.T) {
		e := initCard(t, "card_create__game_not_exist")
		_, err := e.card.create("missing", "ok", "ok", CreateRequest{Name: "ok"})
		assert.ErrorIs(t, err, er.GameNotExists)
	})
}

func TestCardGet(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		e := initCard(t, "card_get__success")
		gameID, collectionID, deckID := e.createParents(t, "parent_game", "parent_collection", "parent_deck")
		wait, err := e.card.create(gameID, collectionID, deckID, CreateRequest{Name: "success",
			Description: "descrption",
			Image:       "https://some.url/image",
			Variables:   map[string]string{"k": "v"},
			Count:       1,
		})
		assert.NoError(t, err)
		got, err := e.card.get(gameID, collectionID, deckID, wait.ID)
		assert.NoError(t, err)
		wait.CreatedAt = got.CreatedAt
		wait.UpdatedAt = got.UpdatedAt
		assert.Equal(t, wait, got)
	})

	t.Run("not_exist", func(t *testing.T) {
		e := initCard(t, "card_get__not_exist")
		gameID, collectionID, deckID := e.createParents(t, "parent_game", "parent_collection", "parent_deck")
		_, err := e.card.get(gameID, collectionID, deckID, 1)
		assert.ErrorIs(t, err, er.CardNotExists)
	})

	t.Run("deck_not_exist", func(t *testing.T) {
		e := initCard(t, "card_get__deck_not_exist")
		gameID := e.createGame(t, "parent_game")
		collectionID := e.createCollection(t, gameID, "parent_collection")
		_, err := e.card.get(gameID, collectionID, "missing", 1)
		assert.ErrorIs(t, err, er.DeckNotExists)
	})

	t.Run("collection_not_exist", func(t *testing.T) {
		e := initCard(t, "card_get__collection_not_exist")
		gameID := e.createGame(t, "parent_game")
		_, err := e.card.get(gameID, "missing", "ok", 1)
		assert.ErrorIs(t, err, er.CollectionNotExists)
	})

	t.Run("game_not_exist", func(t *testing.T) {
		e := initCard(t, "card_get__game_not_exist")
		_, err := e.card.get("missing", "ok", "ok", 1)
		assert.ErrorIs(t, err, er.GameNotExists)
	})
}

func TestCardList(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		e := initCard(t, "card_list__success")
		gameID, collectionID, deckID := e.createParents(t, "parent_game", "parent_collection", "parent_deck")
		first, err := e.card.create(gameID, collectionID, deckID, CreateRequest{Name: "one", Count: 1})
		assert.NoError(t, err)
		second, err := e.card.create(gameID, collectionID, deckID, CreateRequest{Name: "two", Count: 2})
		assert.NoError(t, err)
		got, err := e.card.list(gameID, collectionID, deckID)
		assert.NoError(t, err)
		assert.Len(t, got, 2)
		first.CreatedAt = first.CreatedAt.Truncate(time.Nanosecond)
		first.UpdatedAt = first.UpdatedAt.Truncate(time.Nanosecond)
		second.CreatedAt = second.CreatedAt.Truncate(time.Nanosecond)
		second.UpdatedAt = second.UpdatedAt.Truncate(time.Nanosecond)
		for _, item := range got {
			item.CreatedAt = item.CreatedAt.Truncate(time.Nanosecond)
			item.UpdatedAt = item.UpdatedAt.Truncate(time.Nanosecond)
		}
		assert.ElementsMatch(t, []*entitiesCard.Card{first, second}, got)
	})

	t.Run("empty", func(t *testing.T) {
		e := initCard(t, "card_list__empty")
		gameID, collectionID, deckID := e.createParents(t, "parent_game", "parent_collection", "parent_deck")
		got, err := e.card.list(gameID, collectionID, deckID)
		assert.NoError(t, err)
		assert.Equal(t, []*entitiesCard.Card(nil), got)
	})

	t.Run("isolated_by_deck", func(t *testing.T) {
		e := initCard(t, "card_list__isolated_by_deck")
		gameID, collectionID, deckA := e.createParents(t, "parent_game", "parent_collection", "deck_a")
		deckB := e.createDeck(t, gameID, collectionID, "deck_b")
		_, err := e.card.create(gameID, collectionID, deckA, CreateRequest{Name: "only_a"})
		assert.NoError(t, err)
		got, err := e.card.list(gameID, collectionID, deckB)
		assert.NoError(t, err)
		assert.Equal(t, []*entitiesCard.Card(nil), got)
	})

	t.Run("deck_not_exist", func(t *testing.T) {
		e := initCard(t, "card_list__deck_not_exist")
		gameID := e.createGame(t, "parent_game")
		collectionID := e.createCollection(t, gameID, "parent_collection")
		_, err := e.card.list(gameID, collectionID, "missing")
		assert.ErrorIs(t, err, er.DeckNotExists)
	})

	t.Run("collection_not_exist", func(t *testing.T) {
		e := initCard(t, "card_list__collection_not_exist")
		gameID := e.createGame(t, "parent_game")
		_, err := e.card.list(gameID, "missing", "ok")
		assert.ErrorIs(t, err, er.CollectionNotExists)
	})

	t.Run("game_not_exist", func(t *testing.T) {
		e := initCard(t, "card_list__game_not_exist")
		_, err := e.card.list("missing", "ok", "ok")
		assert.ErrorIs(t, err, er.GameNotExists)
	})
}

func TestCardUpdate(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		e := initCard(t, "card_update__success")
		gameID, collectionID, deckID := e.createParents(t, "parent_game", "parent_collection", "parent_deck")
		wait, err := e.card.create(gameID, collectionID, deckID, CreateRequest{Name: "success",
			Count: 1,
		})
		assert.NoError(t, err)
		got, err := e.card.update(gameID, collectionID, deckID, wait.ID, UpdateRequest{Name: "renamed",
			Description: "desc",
			Image:       "img",
			Variables:   map[string]string{"hp": "10"},
			Count:       4,
		})
		assert.NoError(t, err)
		wait.Name = "renamed"
		wait.Description = "desc"
		wait.Image = "img"
		wait.Variables = map[string]string{"hp": "10"}
		wait.Count = 4
		wait.CreatedAt = got.CreatedAt
		wait.UpdatedAt = got.UpdatedAt
		assert.Equal(t, wait, got)
		assert.True(t, got.UpdatedAt.After(got.CreatedAt) || got.UpdatedAt.Equal(got.CreatedAt))
	})

	t.Run("not_exist", func(t *testing.T) {
		e := initCard(t, "card_update__not_exist")
		gameID, collectionID, deckID := e.createParents(t, "parent_game", "parent_collection", "parent_deck")
		_, err := e.card.update(gameID, collectionID, deckID, 1, UpdateRequest{Name: "ok"})
		assert.ErrorIs(t, err, er.CardNotExists)
	})

	t.Run("deck_not_exist", func(t *testing.T) {
		e := initCard(t, "card_update__deck_not_exist")
		gameID := e.createGame(t, "parent_game")
		collectionID := e.createCollection(t, gameID, "parent_collection")
		_, err := e.card.update(gameID, collectionID, "missing", 1, UpdateRequest{Name: "ok"})
		assert.ErrorIs(t, err, er.DeckNotExists)
	})

	t.Run("collection_not_exist", func(t *testing.T) {
		e := initCard(t, "card_update__collection_not_exist")
		gameID := e.createGame(t, "parent_game")
		_, err := e.card.update(gameID, "missing", "ok", 1, UpdateRequest{Name: "ok"})
		assert.ErrorIs(t, err, er.CollectionNotExists)
	})

	t.Run("game_not_exist", func(t *testing.T) {
		e := initCard(t, "card_update__game_not_exist")
		_, err := e.card.update("missing", "ok", "ok", 1, UpdateRequest{Name: "ok"})
		assert.ErrorIs(t, err, er.GameNotExists)
	})
}

func TestCardDelete(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		e := initCard(t, "card_delete__success")
		gameID, collectionID, deckID := e.createParents(t, "parent_game", "parent_collection", "parent_deck")
		created, err := e.card.create(gameID, collectionID, deckID, CreateRequest{Name: "success"})
		assert.NoError(t, err)
		err = e.card.delete(gameID, collectionID, deckID, created.ID)
		assert.NoError(t, err)
		_, err = e.card.get(gameID, collectionID, deckID, created.ID)
		assert.ErrorIs(t, err, er.CardNotExists)
	})

	t.Run("not_exist", func(t *testing.T) {
		e := initCard(t, "card_delete__not_exist")
		gameID, collectionID, deckID := e.createParents(t, "parent_game", "parent_collection", "parent_deck")
		err := e.card.delete(gameID, collectionID, deckID, 1)
		assert.ErrorIs(t, err, er.CardNotExists)
	})

	t.Run("deck_not_exist", func(t *testing.T) {
		e := initCard(t, "card_delete__deck_not_exist")
		gameID := e.createGame(t, "parent_game")
		collectionID := e.createCollection(t, gameID, "parent_collection")
		err := e.card.delete(gameID, collectionID, "missing", 1)
		assert.ErrorIs(t, err, er.DeckNotExists)
	})

	t.Run("collection_not_exist", func(t *testing.T) {
		e := initCard(t, "card_delete__collection_not_exist")
		gameID := e.createGame(t, "parent_game")
		err := e.card.delete(gameID, "missing", "ok", 1)
		assert.ErrorIs(t, err, er.CollectionNotExists)
	})

	t.Run("game_not_exist", func(t *testing.T) {
		e := initCard(t, "card_delete__game_not_exist")
		err := e.card.delete("missing", "ok", "ok", 1)
		assert.ErrorIs(t, err, er.GameNotExists)
	})
}

func TestCardImageCreate(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		e := initCard(t, "card_image_create__success")
		gameID, collectionID, deckID := e.createParents(t, "parent_game", "parent_collection", "parent_deck")
		created, err := e.card.create(gameID, collectionID, deckID, CreateRequest{Name: "success"})
		assert.NoError(t, err)
		err = e.card.imageCreate(gameID, collectionID, deckID, created.ID, img)
		assert.NoError(t, err)
	})

	t.Run("image_exist", func(t *testing.T) {
		e := initCard(t, "card_image_create__image_exist")
		gameID, collectionID, deckID := e.createParents(t, "parent_game", "parent_collection", "parent_deck")
		created, err := e.card.create(gameID, collectionID, deckID, CreateRequest{Name: "image_exist"})
		assert.NoError(t, err)
		err = e.card.imageCreate(gameID, collectionID, deckID, created.ID, img)
		assert.NoError(t, err)
		err = e.card.imageCreate(gameID, collectionID, deckID, created.ID, img)
		assert.ErrorIs(t, err, er.CardImageExist)
	})

	t.Run("card_not_exist", func(t *testing.T) {
		e := initCard(t, "card_image_create__card_not_exist")
		gameID, collectionID, deckID := e.createParents(t, "parent_game", "parent_collection", "parent_deck")
		err := e.card.imageCreate(gameID, collectionID, deckID, 1, img)
		assert.ErrorIs(t, err, er.CardNotExists)
	})

	t.Run("deck_not_exist", func(t *testing.T) {
		e := initCard(t, "card_image_create__deck_not_exist")
		gameID := e.createGame(t, "parent_game")
		collectionID := e.createCollection(t, gameID, "parent_collection")
		err := e.card.imageCreate(gameID, collectionID, "missing", 1, img)
		assert.ErrorIs(t, err, er.DeckNotExists)
	})

	t.Run("collection_not_exist", func(t *testing.T) {
		e := initCard(t, "card_image_create__collection_not_exist")
		gameID := e.createGame(t, "parent_game")
		err := e.card.imageCreate(gameID, "missing", "ok", 1, img)
		assert.ErrorIs(t, err, er.CollectionNotExists)
	})

	t.Run("game_not_exist", func(t *testing.T) {
		e := initCard(t, "card_image_create__game_not_exist")
		err := e.card.imageCreate("missing", "ok", "ok", 1, img)
		assert.ErrorIs(t, err, er.GameNotExists)
	})
}

func TestCardImageGet(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		e := initCard(t, "card_image_get__success")
		gameID, collectionID, deckID := e.createParents(t, "parent_game", "parent_collection", "parent_deck")
		created, err := e.card.create(gameID, collectionID, deckID, CreateRequest{Name: "success"})
		assert.NoError(t, err)
		err = e.card.imageCreate(gameID, collectionID, deckID, created.ID, img)
		assert.NoError(t, err)
		got, err := e.card.imageGet(gameID, collectionID, deckID, created.ID)
		assert.NoError(t, err)
		assert.Equal(t, img, got)
	})

	t.Run("image_not_exist", func(t *testing.T) {
		e := initCard(t, "card_image_get__image_not_exist")
		gameID, collectionID, deckID := e.createParents(t, "parent_game", "parent_collection", "parent_deck")
		created, err := e.card.create(gameID, collectionID, deckID, CreateRequest{Name: "image_not_exist"})
		assert.NoError(t, err)
		_, err = e.card.imageGet(gameID, collectionID, deckID, created.ID)
		assert.ErrorIs(t, err, er.CardImageNotExists)
	})

	t.Run("card_not_exist", func(t *testing.T) {
		e := initCard(t, "card_image_get__card_not_exist")
		gameID, collectionID, deckID := e.createParents(t, "parent_game", "parent_collection", "parent_deck")
		_, err := e.card.imageGet(gameID, collectionID, deckID, 1)
		assert.ErrorIs(t, err, er.CardNotExists)
	})

	t.Run("deck_not_exist", func(t *testing.T) {
		e := initCard(t, "card_image_get__deck_not_exist")
		gameID := e.createGame(t, "parent_game")
		collectionID := e.createCollection(t, gameID, "parent_collection")
		_, err := e.card.imageGet(gameID, collectionID, "missing", 1)
		assert.ErrorIs(t, err, er.DeckNotExists)
	})

	t.Run("collection_not_exist", func(t *testing.T) {
		e := initCard(t, "card_image_get__collection_not_exist")
		gameID := e.createGame(t, "parent_game")
		_, err := e.card.imageGet(gameID, "missing", "ok", 1)
		assert.ErrorIs(t, err, er.CollectionNotExists)
	})

	t.Run("game_not_exist", func(t *testing.T) {
		e := initCard(t, "card_image_get__game_not_exist")
		_, err := e.card.imageGet("missing", "ok", "ok", 1)
		assert.ErrorIs(t, err, er.GameNotExists)
	})
}

func TestCardImageDelete(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		e := initCard(t, "card_image_delete__success")
		gameID, collectionID, deckID := e.createParents(t, "parent_game", "parent_collection", "parent_deck")
		created, err := e.card.create(gameID, collectionID, deckID, CreateRequest{Name: "success"})
		assert.NoError(t, err)
		err = e.card.imageCreate(gameID, collectionID, deckID, created.ID, img)
		assert.NoError(t, err)
		err = e.card.imageDelete(gameID, collectionID, deckID, created.ID)
		assert.NoError(t, err)
		_, err = e.card.imageGet(gameID, collectionID, deckID, created.ID)
		assert.ErrorIs(t, err, er.CardImageNotExists)
	})

	t.Run("image_not_exist", func(t *testing.T) {
		e := initCard(t, "card_image_delete__image_not_exist")
		gameID, collectionID, deckID := e.createParents(t, "parent_game", "parent_collection", "parent_deck")
		created, err := e.card.create(gameID, collectionID, deckID, CreateRequest{Name: "image_not_exist"})
		assert.NoError(t, err)
		err = e.card.imageDelete(gameID, collectionID, deckID, created.ID)
		assert.ErrorIs(t, err, er.CardImageNotExists)
	})

	t.Run("card_not_exist", func(t *testing.T) {
		e := initCard(t, "card_image_delete__card_not_exist")
		gameID, collectionID, deckID := e.createParents(t, "parent_game", "parent_collection", "parent_deck")
		err := e.card.imageDelete(gameID, collectionID, deckID, 1)
		assert.ErrorIs(t, err, er.CardNotExists)
	})

	t.Run("deck_not_exist", func(t *testing.T) {
		e := initCard(t, "card_image_delete__deck_not_exist")
		gameID := e.createGame(t, "parent_game")
		collectionID := e.createCollection(t, gameID, "parent_collection")
		err := e.card.imageDelete(gameID, collectionID, "missing", 1)
		assert.ErrorIs(t, err, er.DeckNotExists)
	})

	t.Run("collection_not_exist", func(t *testing.T) {
		e := initCard(t, "card_image_delete__collection_not_exist")
		gameID := e.createGame(t, "parent_game")
		err := e.card.imageDelete(gameID, "missing", "ok", 1)
		assert.ErrorIs(t, err, er.CollectionNotExists)
	})

	t.Run("game_not_exist", func(t *testing.T) {
		e := initCard(t, "card_image_delete__game_not_exist")
		err := e.card.imageDelete("missing", "ok", "ok", 1)
		assert.ErrorIs(t, err, er.GameNotExists)
	})
}

func TestCardLegacyTimestamps(t *testing.T) {
	e := initCard(t, "card_legacy_timestamps")
	gameID, collectionID, deckID := e.createParents(t, "parent_game", "parent_collection", "parent_deck")
	created, err := e.card.create(gameID, collectionID, deckID, CreateRequest{Name: "legacy"})
	assert.NoError(t, err)

	path := filepath.Join(e.card.cfg.Games(), gameID, collectionID, deckID, "cards", ".info.json")
	raw, err := os.ReadFile(path)
	assert.NoError(t, err)
	var info struct {
		ID        string                     `json:"id"`
		Name      json.RawMessage            `json:"name"`
		CreatedAt json.RawMessage            `json:"createdAt"`
		UpdatedAt json.RawMessage            `json:"updatedAt"`
		Data      map[string]json.RawMessage `json:"data"`
	}
	assert.NoError(t, json.Unmarshal(raw, &info))
	var cardPayload map[string]json.RawMessage
	assert.NoError(t, json.Unmarshal(info.Data["1"], &cardPayload))
	cardPayload["createdAt"] = json.RawMessage("null")
	cardPayload["updatedAt"] = json.RawMessage("null")
	cardRaw, err := json.Marshal(cardPayload)
	assert.NoError(t, err)
	info.Data["1"] = cardRaw
	out, err := json.Marshal(info)
	assert.NoError(t, err)
	assert.NoError(t, os.WriteFile(path, out, 0o644))

	got, err := e.card.get(gameID, collectionID, deckID, created.ID)
	assert.NoError(t, err)
	assert.False(t, got.CreatedAt.IsZero())
	assert.False(t, got.UpdatedAt.IsZero())

	raw, err = os.ReadFile(path)
	assert.NoError(t, err)
	var stored struct {
		Data map[string]struct {
			CreatedAt *time.Time `json:"createdAt"`
			UpdatedAt *time.Time `json:"updatedAt"`
		} `json:"data"`
	}
	assert.NoError(t, json.Unmarshal(raw, &stored))
	assert.Nil(t, stored.Data["1"].CreatedAt)
	assert.Nil(t, stored.Data["1"].UpdatedAt)
}
