package deck

import (
	"os"
	"testing"
	"time"

	"github.com/HardDie/fsentry"
	"github.com/stretchr/testify/assert"

	"github.com/HardDie/DeckBuilder/internal/config"
	dbCore "github.com/HardDie/DeckBuilder/internal/db/core"
	entitiesDeck "github.com/HardDie/DeckBuilder/internal/entities/deck"
	er "github.com/HardDie/DeckBuilder/internal/errors"
	repositoriesCollection "github.com/HardDie/DeckBuilder/internal/repositories/collection"
	"github.com/HardDie/DeckBuilder/internal/utils"
)

var (
	img = []byte("some_image")
)

type deckEnv struct {
	db   *fsentry.DB
	deck *deck
}

func initDeck(t testing.TB, name string) deckEnv {
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

	core := dbCore.New(db)
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

	return deckEnv{
		db:   db,
		deck: New(cfg, db, repositoriesCollection.New(cfg, db)).(*deck),
	}
}

func (e deckEnv) createGame(t testing.TB, name string) string {
	t.Helper()
	info, err := e.db.CreateFolder[any](name, nil, "games")
	if err != nil {
		t.Fatal("error create game", err)
	}
	return info.ID
}

func (e deckEnv) createCollection(t testing.TB, gameID, name string) string {
	t.Helper()
	info, err := e.db.CreateFolder[any](name, nil, "games", gameID)
	if err != nil {
		t.Fatal("error create collection", err)
	}
	return info.ID
}

func (e deckEnv) createParents(t testing.TB, gameName, collectionName string) (string, string) {
	t.Helper()
	gameID := e.createGame(t, gameName)
	collectionID := e.createCollection(t, gameID, collectionName)
	return gameID, collectionID
}

func TestDeckCreate(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		wait := &entitiesDeck.Deck{
			Name:         "success",
			Description:  "descrption",
			Image:        "https://some.url/image",
			GameID:       "parent_game",
			CollectionID: "parent_collection",
		}
		wait.ID = utils.NameToID(wait.Name)

		e := initDeck(t, "deck_create__success")
		gameID, collectionID := e.createParents(t, wait.GameID, wait.CollectionID)
		wait.GameID = gameID
		wait.CollectionID = collectionID
		got, err := e.deck.create(wait.GameID, wait.CollectionID, CreateRequest{Name: wait.Name,
			Description: wait.Description,
			Image:       wait.Image,
		})
		assert.NoError(t, err)
		wait.CreatedAt = got.CreatedAt
		wait.UpdatedAt = got.UpdatedAt
		assert.Equal(t, wait, got)
	})

	t.Run("exist", func(t *testing.T) {
		e := initDeck(t, "deck_create__exist")
		gameID, collectionID := e.createParents(t, "parent_game", "parent_collection")
		_, err := e.deck.create(gameID, collectionID, CreateRequest{Name: "exist"})
		assert.NoError(t, err)
		_, err = e.deck.create(gameID, collectionID, CreateRequest{Name: "exist"})
		assert.ErrorIs(t, err, er.DeckExist)
	})

	t.Run("same_name_other_collection", func(t *testing.T) {
		e := initDeck(t, "deck_create__same_name_other_collection")
		gameID, collectionA := e.createParents(t, "parent_game", "collection_a")
		collectionB := e.createCollection(t, gameID, "collection_b")
		_, err := e.deck.create(gameID, collectionA, CreateRequest{Name: "shared"})
		assert.NoError(t, err)
		_, err = e.deck.create(gameID, collectionB, CreateRequest{Name: "shared"})
		assert.NoError(t, err)
	})

	t.Run("bad_name", func(t *testing.T) {
		e := initDeck(t, "deck_create__bad_name")
		gameID, collectionID := e.createParents(t, "parent_game", "parent_collection")
		_, err := e.deck.create(gameID, collectionID, CreateRequest{Name: "---"})
		assert.ErrorIs(t, err, er.BadName)
	})

	t.Run("collection_not_exist", func(t *testing.T) {
		e := initDeck(t, "deck_create__collection_not_exist")
		gameID := e.createGame(t, "parent_game")
		_, err := e.deck.create(gameID, "missing", CreateRequest{Name: "ok"})
		assert.ErrorIs(t, err, er.CollectionNotExists)
	})

	t.Run("game_not_exist", func(t *testing.T) {
		e := initDeck(t, "deck_create__game_not_exist")
		_, err := e.deck.create("missing", "ok", CreateRequest{Name: "ok"})
		assert.ErrorIs(t, err, er.GameNotExists)
	})
}

func TestDeckGet(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		name := "success"
		desc := "descrption"
		image := "https://some.url/image"

		e := initDeck(t, "deck_get__success")
		gameID, collectionID := e.createParents(t, "parent_game", "parent_collection")
		wait, err := e.deck.create(gameID, collectionID, CreateRequest{Name: name,
			Description: desc,
			Image:       image,
		})
		assert.NoError(t, err)
		got, err := e.deck.get(gameID, collectionID, name)
		assert.NoError(t, err)
		wait.CreatedAt = got.CreatedAt
		wait.UpdatedAt = got.UpdatedAt
		assert.Equal(t, wait, got)
	})

	t.Run("not_exist", func(t *testing.T) {
		e := initDeck(t, "deck_get__not_exist")
		gameID, collectionID := e.createParents(t, "parent_game", "parent_collection")
		_, err := e.deck.get(gameID, collectionID, "not_exist")
		assert.ErrorIs(t, err, er.DeckNotExists)
	})

	t.Run("bad_name", func(t *testing.T) {
		e := initDeck(t, "deck_get__bad_name")
		gameID, collectionID := e.createParents(t, "parent_game", "parent_collection")
		_, err := e.deck.get(gameID, collectionID, "---")
		assert.ErrorIs(t, err, er.BadName)
	})

	t.Run("collection_not_exist", func(t *testing.T) {
		e := initDeck(t, "deck_get__collection_not_exist")
		gameID := e.createGame(t, "parent_game")
		_, err := e.deck.get(gameID, "missing", "ok")
		assert.ErrorIs(t, err, er.CollectionNotExists)
	})

	t.Run("game_not_exist", func(t *testing.T) {
		e := initDeck(t, "deck_get__game_not_exist")
		_, err := e.deck.get("missing", "ok", "ok")
		assert.ErrorIs(t, err, er.GameNotExists)
	})
}

func TestDeckList(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		name := "success"
		desc := "descrption"
		image := "https://some.url/image"

		e := initDeck(t, "deck_list__success")
		gameID, collectionID := e.createParents(t, "parent_game", "parent_collection")
		wait, err := e.deck.create(gameID, collectionID, CreateRequest{Name: name,
			Description: desc,
			Image:       image,
		})
		assert.NoError(t, err)
		got, err := e.deck.list(gameID, collectionID)
		assert.NoError(t, err)
		assert.Len(t, got, 1)
		wait.CreatedAt = got[0].CreatedAt
		wait.UpdatedAt = got[0].UpdatedAt
		assert.Equal(t, []*entitiesDeck.Deck{wait}, got)
	})

	t.Run("empty", func(t *testing.T) {
		e := initDeck(t, "deck_list__empty")
		gameID, collectionID := e.createParents(t, "parent_game", "parent_collection")
		got, err := e.deck.list(gameID, collectionID)
		assert.NoError(t, err)
		assert.Equal(t, []*entitiesDeck.Deck(nil), got)
	})

	t.Run("isolated_by_collection", func(t *testing.T) {
		e := initDeck(t, "deck_list__isolated_by_collection")
		gameID, collectionA := e.createParents(t, "parent_game", "collection_a")
		collectionB := e.createCollection(t, gameID, "collection_b")
		_, err := e.deck.create(gameID, collectionA, CreateRequest{Name: "only_a"})
		assert.NoError(t, err)
		got, err := e.deck.list(gameID, collectionB)
		assert.NoError(t, err)
		assert.Equal(t, []*entitiesDeck.Deck(nil), got)
	})

	t.Run("collection_not_exist", func(t *testing.T) {
		e := initDeck(t, "deck_list__collection_not_exist")
		gameID := e.createGame(t, "parent_game")
		_, err := e.deck.list(gameID, "missing")
		assert.ErrorIs(t, err, er.CollectionNotExists)
	})

	t.Run("game_not_exist", func(t *testing.T) {
		e := initDeck(t, "deck_list__game_not_exist")
		_, err := e.deck.list("missing", "ok")
		assert.ErrorIs(t, err, er.GameNotExists)
	})
}

func TestDeckMove(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		oldName := "success_old"
		newName := "success_new"
		desc := "descrption"
		image := "https://some.url/image"

		e := initDeck(t, "deck_move__success")
		gameID, collectionID := e.createParents(t, "parent_game", "parent_collection")
		oldDeck, err := e.deck.create(gameID, collectionID, CreateRequest{Name: oldName,
			Description: desc,
			Image:       image,
		})
		assert.NoError(t, err)
		newDeck, err := e.deck.move(gameID, collectionID, oldName, newName)
		assert.NoError(t, err)

		oldDeck.ID = utils.NameToID(newName)
		oldDeck.Name = newName
		oldDeck.CreatedAt = oldDeck.CreatedAt.Truncate(time.Nanosecond)
		oldDeck.UpdatedAt = newDeck.UpdatedAt
		assert.Equal(t, oldDeck, newDeck)
	})

	t.Run("not_exist", func(t *testing.T) {
		e := initDeck(t, "deck_move__not_exist")
		gameID, collectionID := e.createParents(t, "parent_game", "parent_collection")
		_, err := e.deck.move(gameID, collectionID, "not_exist", "new_name")
		assert.ErrorIs(t, err, er.DeckNotExists)
	})

	t.Run("bad_name", func(t *testing.T) {
		oldName := "bad_name_old"
		e := initDeck(t, "deck_move__bad_name")
		gameID, collectionID := e.createParents(t, "parent_game", "parent_collection")
		_, err := e.deck.create(gameID, collectionID, CreateRequest{Name: oldName})
		assert.NoError(t, err)
		_, err = e.deck.move(gameID, collectionID, oldName, "---")
		assert.ErrorIs(t, err, er.BadName)
	})

	t.Run("collection_not_exist", func(t *testing.T) {
		e := initDeck(t, "deck_move__collection_not_exist")
		gameID := e.createGame(t, "parent_game")
		_, err := e.deck.move(gameID, "missing", "old", "new")
		assert.ErrorIs(t, err, er.CollectionNotExists)
	})

	t.Run("game_not_exist", func(t *testing.T) {
		e := initDeck(t, "deck_move__game_not_exist")
		_, err := e.deck.move("missing", "ok", "old", "new")
		assert.ErrorIs(t, err, er.GameNotExists)
	})
}

func TestDeckUpdate(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		name := "success"
		newDesc := "desc"
		newImage := "img"

		e := initDeck(t, "deck_update__success")
		gameID, collectionID := e.createParents(t, "parent_game", "parent_collection")
		wait, err := e.deck.create(gameID, collectionID, CreateRequest{Name: name})
		assert.NoError(t, err)
		got, err := e.deck.update(gameID, collectionID, updateRequest{Name: name,
			Description: newDesc,
			Image:       newImage,
		})
		assert.NoError(t, err)
		wait.Description = newDesc
		wait.Image = newImage
		wait.CreatedAt = got.CreatedAt
		wait.UpdatedAt = got.UpdatedAt
		assert.Equal(t, wait, got)
	})

	t.Run("not_exist", func(t *testing.T) {
		e := initDeck(t, "deck_update__not_exist")
		gameID, collectionID := e.createParents(t, "parent_game", "parent_collection")
		_, err := e.deck.update(gameID, collectionID, updateRequest{Name: "not_exist"})
		assert.ErrorIs(t, err, er.DeckNotExists)
	})

	t.Run("bad_name", func(t *testing.T) {
		e := initDeck(t, "deck_update__bad_name")
		gameID, collectionID := e.createParents(t, "parent_game", "parent_collection")
		_, err := e.deck.update(gameID, collectionID, updateRequest{Name: "---"})
		assert.ErrorIs(t, err, er.BadName)
	})

	t.Run("collection_not_exist", func(t *testing.T) {
		e := initDeck(t, "deck_update__collection_not_exist")
		gameID := e.createGame(t, "parent_game")
		_, err := e.deck.update(gameID, "missing", updateRequest{Name: "ok"})
		assert.ErrorIs(t, err, er.CollectionNotExists)
	})

	t.Run("game_not_exist", func(t *testing.T) {
		e := initDeck(t, "deck_update__game_not_exist")
		_, err := e.deck.update("missing", "ok", updateRequest{Name: "ok"})
		assert.ErrorIs(t, err, er.GameNotExists)
	})
}

func TestDeckDelete(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		name := "success"
		e := initDeck(t, "deck_delete__success")
		gameID, collectionID := e.createParents(t, "parent_game", "parent_collection")
		_, err := e.deck.create(gameID, collectionID, CreateRequest{Name: name})
		assert.NoError(t, err)
		err = e.deck.delete(gameID, collectionID, name)
		assert.NoError(t, err)
		_, err = e.deck.get(gameID, collectionID, name)
		assert.ErrorIs(t, err, er.DeckNotExists)
	})

	t.Run("not_exist", func(t *testing.T) {
		e := initDeck(t, "deck_delete__not_exist")
		gameID, collectionID := e.createParents(t, "parent_game", "parent_collection")
		err := e.deck.delete(gameID, collectionID, "not_exist")
		assert.ErrorIs(t, err, er.DeckNotExists)
	})

	t.Run("bad_name", func(t *testing.T) {
		e := initDeck(t, "deck_delete__bad_name")
		gameID, collectionID := e.createParents(t, "parent_game", "parent_collection")
		err := e.deck.delete(gameID, collectionID, "---")
		assert.ErrorIs(t, err, er.BadName)
	})

	t.Run("collection_not_exist", func(t *testing.T) {
		e := initDeck(t, "deck_delete__collection_not_exist")
		gameID := e.createGame(t, "parent_game")
		err := e.deck.delete(gameID, "missing", "ok")
		assert.ErrorIs(t, err, er.CollectionNotExists)
	})

	t.Run("game_not_exist", func(t *testing.T) {
		e := initDeck(t, "deck_delete__game_not_exist")
		err := e.deck.delete("missing", "ok", "ok")
		assert.ErrorIs(t, err, er.GameNotExists)
	})
}

func TestDeckImageCreate(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		name := "success"
		e := initDeck(t, "deck_image_create__success")
		gameID, collectionID := e.createParents(t, "parent_game", "parent_collection")
		_, err := e.deck.create(gameID, collectionID, CreateRequest{Name: name})
		assert.NoError(t, err)
		err = e.deck.imageCreate(gameID, collectionID, name, img)
		assert.NoError(t, err)
	})

	t.Run("image_exist", func(t *testing.T) {
		name := "image_exist"
		e := initDeck(t, "deck_image_create__image_exist")
		gameID, collectionID := e.createParents(t, "parent_game", "parent_collection")
		_, err := e.deck.create(gameID, collectionID, CreateRequest{Name: name})
		assert.NoError(t, err)
		err = e.deck.imageCreate(gameID, collectionID, name, img)
		assert.NoError(t, err)
		err = e.deck.imageCreate(gameID, collectionID, name, img)
		assert.ErrorIs(t, err, er.DeckImageExist)
	})

	t.Run("deck_not_exist", func(t *testing.T) {
		e := initDeck(t, "deck_image_create__deck_not_exist")
		gameID, collectionID := e.createParents(t, "parent_game", "parent_collection")
		err := e.deck.imageCreate(gameID, collectionID, "missing", img)
		assert.ErrorIs(t, err, er.DeckNotExists)
	})

	t.Run("collection_not_exist", func(t *testing.T) {
		e := initDeck(t, "deck_image_create__collection_not_exist")
		gameID := e.createGame(t, "parent_game")
		err := e.deck.imageCreate(gameID, "missing", "ok", img)
		assert.ErrorIs(t, err, er.CollectionNotExists)
	})

	t.Run("game_not_exist", func(t *testing.T) {
		e := initDeck(t, "deck_image_create__game_not_exist")
		err := e.deck.imageCreate("missing", "ok", "ok", img)
		assert.ErrorIs(t, err, er.GameNotExists)
	})
}

func TestDeckImageGet(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		name := "success"
		e := initDeck(t, "deck_image_get__success")
		gameID, collectionID := e.createParents(t, "parent_game", "parent_collection")
		_, err := e.deck.create(gameID, collectionID, CreateRequest{Name: name})
		assert.NoError(t, err)
		err = e.deck.imageCreate(gameID, collectionID, name, img)
		assert.NoError(t, err)
		got, err := e.deck.imageGet(gameID, collectionID, name)
		assert.NoError(t, err)
		assert.Equal(t, img, got)
	})

	t.Run("image_not_exist", func(t *testing.T) {
		name := "image_not_exist"
		e := initDeck(t, "deck_image_get__image_not_exist")
		gameID, collectionID := e.createParents(t, "parent_game", "parent_collection")
		_, err := e.deck.create(gameID, collectionID, CreateRequest{Name: name})
		assert.NoError(t, err)
		_, err = e.deck.imageGet(gameID, collectionID, name)
		assert.ErrorIs(t, err, er.DeckImageNotExists)
	})

	t.Run("deck_not_exist", func(t *testing.T) {
		e := initDeck(t, "deck_image_get__deck_not_exist")
		gameID, collectionID := e.createParents(t, "parent_game", "parent_collection")
		_, err := e.deck.imageGet(gameID, collectionID, "missing")
		assert.ErrorIs(t, err, er.DeckNotExists)
	})

	t.Run("collection_not_exist", func(t *testing.T) {
		e := initDeck(t, "deck_image_get__collection_not_exist")
		gameID := e.createGame(t, "parent_game")
		_, err := e.deck.imageGet(gameID, "missing", "ok")
		assert.ErrorIs(t, err, er.CollectionNotExists)
	})

	t.Run("game_not_exist", func(t *testing.T) {
		e := initDeck(t, "deck_image_get__game_not_exist")
		_, err := e.deck.imageGet("missing", "ok", "ok")
		assert.ErrorIs(t, err, er.GameNotExists)
	})
}

func TestDeckImageDelete(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		name := "success"
		e := initDeck(t, "deck_image_delete__success")
		gameID, collectionID := e.createParents(t, "parent_game", "parent_collection")
		_, err := e.deck.create(gameID, collectionID, CreateRequest{Name: name})
		assert.NoError(t, err)
		err = e.deck.imageCreate(gameID, collectionID, name, img)
		assert.NoError(t, err)
		err = e.deck.imageDelete(gameID, collectionID, name)
		assert.NoError(t, err)
		_, err = e.deck.imageGet(gameID, collectionID, name)
		assert.ErrorIs(t, err, er.DeckImageNotExists)
	})

	t.Run("image_not_exist", func(t *testing.T) {
		name := "image_not_exist"
		e := initDeck(t, "deck_image_delete__image_not_exist")
		gameID, collectionID := e.createParents(t, "parent_game", "parent_collection")
		_, err := e.deck.create(gameID, collectionID, CreateRequest{Name: name})
		assert.NoError(t, err)
		err = e.deck.imageDelete(gameID, collectionID, name)
		assert.ErrorIs(t, err, er.DeckImageNotExists)
	})

	t.Run("deck_not_exist", func(t *testing.T) {
		e := initDeck(t, "deck_image_delete__deck_not_exist")
		gameID, collectionID := e.createParents(t, "parent_game", "parent_collection")
		err := e.deck.imageDelete(gameID, collectionID, "missing")
		assert.ErrorIs(t, err, er.DeckNotExists)
	})

	t.Run("collection_not_exist", func(t *testing.T) {
		e := initDeck(t, "deck_image_delete__collection_not_exist")
		gameID := e.createGame(t, "parent_game")
		err := e.deck.imageDelete(gameID, "missing", "ok")
		assert.ErrorIs(t, err, er.CollectionNotExists)
	})

	t.Run("game_not_exist", func(t *testing.T) {
		e := initDeck(t, "deck_image_delete__game_not_exist")
		err := e.deck.imageDelete("missing", "ok", "ok")
		assert.ErrorIs(t, err, er.GameNotExists)
	})
}
