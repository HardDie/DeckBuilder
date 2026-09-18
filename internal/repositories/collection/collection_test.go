package collection

import (
	"os"
	"testing"
	"time"

	"github.com/HardDie/fsentry"
	"github.com/stretchr/testify/assert"

	"github.com/HardDie/DeckBuilder/internal/config"
	entitiesCollection "github.com/HardDie/DeckBuilder/internal/entities/collection"
	er "github.com/HardDie/DeckBuilder/internal/errors"
	repositoriesCore "github.com/HardDie/DeckBuilder/internal/repositories/core"
	"github.com/HardDie/DeckBuilder/internal/utils"
)

var (
	img = []byte("some_image")
)

type collectionEnv struct {
	db         *fsentry.DB
	collection *collection
}

func initCollection(t testing.TB, name string) collectionEnv {
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

	return collectionEnv{
		db:         db,
		collection: New(cfg, db).(*collection),
	}
}

func (e collectionEnv) createGame(t testing.TB, name string) string {
	t.Helper()
	info, err := e.db.CreateFolder[any](name, nil, "games")
	if err != nil {
		t.Fatal("error create game", err)
	}
	return info.ID
}

func TestCollectionCreate(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		wait := &entitiesCollection.Collection{
			Name:        "success",
			Description: "descrption",
			Image:       "https://some.url/image",
			GameID:      "parent_game",
		}
		wait.ID = utils.NameToID(wait.Name)

		e := initCollection(t, "collection_create__success")
		e.createGame(t, wait.GameID)
		got, err := e.collection.create(wait.GameID, CreateRequest{Name: wait.Name,
			Description: wait.Description,
			Image:       wait.Image,
		})
		assert.NoError(t, err)
		wait.CreatedAt = got.CreatedAt
		wait.UpdatedAt = got.UpdatedAt
		assert.Equal(t, wait, got)
	})

	t.Run("exist", func(t *testing.T) {
		gameID := "parent_game"
		e := initCollection(t, "collection_create__exist")
		e.createGame(t, gameID)
		_, err := e.collection.create(gameID, CreateRequest{Name: "exist"})
		assert.NoError(t, err)
		_, err = e.collection.create(gameID, CreateRequest{Name: "exist"})
		assert.ErrorIs(t, err, er.CollectionExist)
	})

	t.Run("same_name_other_game", func(t *testing.T) {
		e := initCollection(t, "collection_create__same_name_other_game")
		gameA := e.createGame(t, "game_a")
		gameB := e.createGame(t, "game_b")
		_, err := e.collection.create(gameA, CreateRequest{Name: "shared"})
		assert.NoError(t, err)
		_, err = e.collection.create(gameB, CreateRequest{Name: "shared"})
		assert.NoError(t, err)
	})

	t.Run("bad_name", func(t *testing.T) {
		gameID := "parent_game"
		e := initCollection(t, "collection_create__bad_name")
		e.createGame(t, gameID)
		_, err := e.collection.create(gameID, CreateRequest{Name: "---"})
		assert.ErrorIs(t, err, er.BadName)
	})

	t.Run("game_not_exist", func(t *testing.T) {
		e := initCollection(t, "collection_create__game_not_exist")
		_, err := e.collection.create("missing", CreateRequest{Name: "ok"})
		assert.ErrorIs(t, err, er.GameNotExists)
	})
}

func TestCollectionGet(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		gameID := "parent_game"
		name := "success"
		desc := "descrption"
		image := "https://some.url/image"

		e := initCollection(t, "collection_get__success")
		e.createGame(t, gameID)
		wait, err := e.collection.create(gameID, CreateRequest{Name: name,
			Description: desc,
			Image:       image,
		})
		assert.NoError(t, err)
		got, err := e.collection.get(gameID, name)
		assert.NoError(t, err)
		wait.CreatedAt = got.CreatedAt
		wait.UpdatedAt = got.UpdatedAt
		assert.Equal(t, wait, got)
	})

	t.Run("not_exist", func(t *testing.T) {
		gameID := "parent_game"
		e := initCollection(t, "collection_get__not_exist")
		e.createGame(t, gameID)
		_, err := e.collection.get(gameID, "not_exist")
		assert.ErrorIs(t, err, er.CollectionNotExists)
	})

	t.Run("bad_name", func(t *testing.T) {
		gameID := "parent_game"
		e := initCollection(t, "collection_get__bad_name")
		e.createGame(t, gameID)
		_, err := e.collection.get(gameID, "---")
		assert.ErrorIs(t, err, er.BadName)
	})

	t.Run("game_not_exist", func(t *testing.T) {
		e := initCollection(t, "collection_get__game_not_exist")
		_, err := e.collection.get("missing", "ok")
		assert.ErrorIs(t, err, er.GameNotExists)
	})
}

func TestCollectionList(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		gameID := "parent_game"
		name := "success"
		desc := "descrption"
		image := "https://some.url/image"

		e := initCollection(t, "collection_list__success")
		e.createGame(t, gameID)
		wait, err := e.collection.create(gameID, CreateRequest{Name: name,
			Description: desc,
			Image:       image,
		})
		assert.NoError(t, err)
		got, err := e.collection.list(gameID)
		assert.NoError(t, err)
		assert.Len(t, got, 1)
		wait.CreatedAt = got[0].CreatedAt
		wait.UpdatedAt = got[0].UpdatedAt
		assert.Equal(t, []*entitiesCollection.Collection{wait}, got)
	})

	t.Run("empty", func(t *testing.T) {
		gameID := "parent_game"
		e := initCollection(t, "collection_list__empty")
		e.createGame(t, gameID)
		got, err := e.collection.list(gameID)
		assert.NoError(t, err)
		assert.Equal(t, []*entitiesCollection.Collection(nil), got)
	})

	t.Run("isolated_by_game", func(t *testing.T) {
		e := initCollection(t, "collection_list__isolated_by_game")
		gameA := e.createGame(t, "game_a")
		gameB := e.createGame(t, "game_b")
		_, err := e.collection.create(gameA, CreateRequest{Name: "only_a"})
		assert.NoError(t, err)
		got, err := e.collection.list(gameB)
		assert.NoError(t, err)
		assert.Equal(t, []*entitiesCollection.Collection(nil), got)
	})

	t.Run("game_not_exist", func(t *testing.T) {
		e := initCollection(t, "collection_list__game_not_exist")
		_, err := e.collection.list("missing")
		assert.ErrorIs(t, err, er.GameNotExists)
	})
}

func TestCollectionMove(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		gameID := "parent_game"
		oldName := "success_old"
		newName := "success_new"
		desc := "descrption"
		image := "https://some.url/image"

		e := initCollection(t, "collection_move__success")
		e.createGame(t, gameID)
		oldCollection, err := e.collection.create(gameID, CreateRequest{Name: oldName,
			Description: desc,
			Image:       image,
		})
		assert.NoError(t, err)
		newCollection, err := e.collection.move(gameID, oldName, newName)
		assert.NoError(t, err)

		oldCollection.ID = utils.NameToID(newName)
		oldCollection.Name = newName
		oldCollection.CreatedAt = oldCollection.CreatedAt.Truncate(time.Nanosecond)
		oldCollection.UpdatedAt = newCollection.UpdatedAt
		assert.Equal(t, oldCollection, newCollection)
	})

	t.Run("not_exist", func(t *testing.T) {
		gameID := "parent_game"
		e := initCollection(t, "collection_move__not_exist")
		e.createGame(t, gameID)
		_, err := e.collection.move(gameID, "not_exist", "new_name")
		assert.ErrorIs(t, err, er.CollectionNotExists)
	})

	t.Run("bad_name", func(t *testing.T) {
		gameID := "parent_game"
		oldName := "bad_name_old"
		e := initCollection(t, "collection_move__bad_name")
		e.createGame(t, gameID)
		_, err := e.collection.create(gameID, CreateRequest{Name: oldName})
		assert.NoError(t, err)
		_, err = e.collection.move(gameID, oldName, "---")
		assert.ErrorIs(t, err, er.BadName)
	})

	t.Run("game_not_exist", func(t *testing.T) {
		e := initCollection(t, "collection_move__game_not_exist")
		_, err := e.collection.move("missing", "old", "new")
		assert.ErrorIs(t, err, er.GameNotExists)
	})
}

func TestCollectionUpdate(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		gameID := "parent_game"
		name := "success"
		newDesc := "desc"
		newImage := "img"

		e := initCollection(t, "collection_update__success")
		e.createGame(t, gameID)
		wait, err := e.collection.create(gameID, CreateRequest{Name: name})
		assert.NoError(t, err)
		got, err := e.collection.update(gameID, updateRequest{Name: name,
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
		gameID := "parent_game"
		e := initCollection(t, "collection_update__not_exist")
		e.createGame(t, gameID)
		_, err := e.collection.update(gameID, updateRequest{Name: "not_exist"})
		assert.ErrorIs(t, err, er.CollectionNotExists)
	})

	t.Run("bad_name", func(t *testing.T) {
		gameID := "parent_game"
		e := initCollection(t, "collection_update__bad_name")
		e.createGame(t, gameID)
		_, err := e.collection.update(gameID, updateRequest{Name: "---"})
		assert.ErrorIs(t, err, er.BadName)
	})

	t.Run("game_not_exist", func(t *testing.T) {
		e := initCollection(t, "collection_update__game_not_exist")
		_, err := e.collection.update("missing", updateRequest{Name: "ok"})
		assert.ErrorIs(t, err, er.GameNotExists)
	})
}

func TestCollectionDelete(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		gameID := "parent_game"
		name := "success"
		e := initCollection(t, "collection_delete__success")
		e.createGame(t, gameID)
		_, err := e.collection.create(gameID, CreateRequest{Name: name})
		assert.NoError(t, err)
		err = e.collection.delete(gameID, name)
		assert.NoError(t, err)
		_, err = e.collection.get(gameID, name)
		assert.ErrorIs(t, err, er.CollectionNotExists)
	})

	t.Run("not_exist", func(t *testing.T) {
		gameID := "parent_game"
		e := initCollection(t, "collection_delete__not_exist")
		e.createGame(t, gameID)
		err := e.collection.delete(gameID, "not_exist")
		assert.ErrorIs(t, err, er.CollectionNotExists)
	})

	t.Run("bad_name", func(t *testing.T) {
		gameID := "parent_game"
		e := initCollection(t, "collection_delete__bad_name")
		e.createGame(t, gameID)
		err := e.collection.delete(gameID, "---")
		assert.ErrorIs(t, err, er.BadName)
	})

	t.Run("game_not_exist", func(t *testing.T) {
		e := initCollection(t, "collection_delete__game_not_exist")
		err := e.collection.delete("missing", "ok")
		assert.ErrorIs(t, err, er.GameNotExists)
	})
}

func TestCollectionImageCreate(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		gameID := "parent_game"
		name := "success"

		e := initCollection(t, "collection_image_create__success")
		e.createGame(t, gameID)
		_, err := e.collection.create(gameID, CreateRequest{Name: name})
		assert.NoError(t, err)
		err = e.collection.imageCreate(gameID, name, img)
		assert.NoError(t, err)
	})

	t.Run("image_exist", func(t *testing.T) {
		gameID := "parent_game"
		name := "image_exist"

		e := initCollection(t, "collection_image_create__image_exist")
		e.createGame(t, gameID)
		_, err := e.collection.create(gameID, CreateRequest{Name: name})
		assert.NoError(t, err)
		err = e.collection.imageCreate(gameID, name, img)
		assert.NoError(t, err)
		err = e.collection.imageCreate(gameID, name, img)
		assert.ErrorIs(t, err, er.CollectionImageExist)
	})

	t.Run("collection_not_exist", func(t *testing.T) {
		gameID := "parent_game"
		e := initCollection(t, "collection_image_create__collection_not_exist")
		e.createGame(t, gameID)
		err := e.collection.imageCreate(gameID, "missing", img)
		assert.ErrorIs(t, err, er.CollectionNotExists)
	})

	t.Run("game_not_exist", func(t *testing.T) {
		e := initCollection(t, "collection_image_create__game_not_exist")
		err := e.collection.imageCreate("missing", "ok", img)
		assert.ErrorIs(t, err, er.GameNotExists)
	})
}

func TestCollectionImageGet(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		gameID := "parent_game"
		name := "success"

		e := initCollection(t, "collection_image_get__success")
		e.createGame(t, gameID)
		_, err := e.collection.create(gameID, CreateRequest{Name: name})
		assert.NoError(t, err)
		err = e.collection.imageCreate(gameID, name, img)
		assert.NoError(t, err)
		got, err := e.collection.imageGet(gameID, name)
		assert.NoError(t, err)
		assert.Equal(t, img, got)
	})

	t.Run("image_not_exist", func(t *testing.T) {
		gameID := "parent_game"
		name := "image_not_exist"

		e := initCollection(t, "collection_image_get__image_not_exist")
		e.createGame(t, gameID)
		_, err := e.collection.create(gameID, CreateRequest{Name: name})
		assert.NoError(t, err)
		_, err = e.collection.imageGet(gameID, name)
		assert.ErrorIs(t, err, er.CollectionImageNotExists)
	})

	t.Run("collection_not_exist", func(t *testing.T) {
		gameID := "parent_game"
		e := initCollection(t, "collection_image_get__collection_not_exist")
		e.createGame(t, gameID)
		_, err := e.collection.imageGet(gameID, "missing")
		assert.ErrorIs(t, err, er.CollectionNotExists)
	})

	t.Run("game_not_exist", func(t *testing.T) {
		e := initCollection(t, "collection_image_get__game_not_exist")
		_, err := e.collection.imageGet("missing", "ok")
		assert.ErrorIs(t, err, er.GameNotExists)
	})
}

func TestCollectionImageDelete(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		gameID := "parent_game"
		name := "success"

		e := initCollection(t, "collection_image_delete__success")
		e.createGame(t, gameID)
		_, err := e.collection.create(gameID, CreateRequest{Name: name})
		assert.NoError(t, err)
		err = e.collection.imageCreate(gameID, name, img)
		assert.NoError(t, err)
		err = e.collection.imageDelete(gameID, name)
		assert.NoError(t, err)
		_, err = e.collection.imageGet(gameID, name)
		assert.ErrorIs(t, err, er.CollectionImageNotExists)
	})

	t.Run("image_not_exist", func(t *testing.T) {
		gameID := "parent_game"
		name := "image_not_exist"

		e := initCollection(t, "collection_image_delete__image_not_exist")
		e.createGame(t, gameID)
		_, err := e.collection.create(gameID, CreateRequest{Name: name})
		assert.NoError(t, err)
		err = e.collection.imageDelete(gameID, name)
		assert.ErrorIs(t, err, er.CollectionImageNotExists)
	})

	t.Run("collection_not_exist", func(t *testing.T) {
		gameID := "parent_game"
		e := initCollection(t, "collection_image_delete__collection_not_exist")
		e.createGame(t, gameID)
		err := e.collection.imageDelete(gameID, "missing")
		assert.ErrorIs(t, err, er.CollectionNotExists)
	})

	t.Run("game_not_exist", func(t *testing.T) {
		e := initCollection(t, "collection_image_delete__game_not_exist")
		err := e.collection.imageDelete("missing", "ok")
		assert.ErrorIs(t, err, er.GameNotExists)
	})
}
