package collection

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/HardDie/fsentry"
	"github.com/stretchr/testify/assert"

	"github.com/HardDie/DeckBuilder/internal/config"
	dbCore "github.com/HardDie/DeckBuilder/internal/db/core"
	dbGame "github.com/HardDie/DeckBuilder/internal/db/game"
	entitiesCollection "github.com/HardDie/DeckBuilder/internal/entities/collection"
	er "github.com/HardDie/DeckBuilder/internal/errors"
	"github.com/HardDie/DeckBuilder/internal/utils"
)

var (
	img = []byte("some_image")
)

type collectionEnv struct {
	game       dbGame.Game
	collection Collection
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

	game := dbGame.New(db)
	return collectionEnv{
		game:       game,
		collection: New(db, game),
	}
}

func (e collectionEnv) createGame(t testing.TB, ctx context.Context, name string) string {
	t.Helper()
	got, err := e.game.Create(ctx, dbGame.CreateRequest{Name: name})
	if err != nil {
		t.Fatal("error create game", err)
	}
	return got.ID
}

func TestCollectionCreate(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		wait := &entitiesCollection.Collection{
			Name:        "success",
			Description: "descrption",
			Image:       "https://some.url/image",
			GameID:      "parent_game",
		}
		wait.ID = utils.NameToID(wait.Name)

		e := initCollection(t, "collection_create__success")
		e.createGame(t, ctx, wait.GameID)
		got, err := e.collection.Create(ctx, CreateRequest{
			GameID:      wait.GameID,
			Name:        wait.Name,
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
		e.createGame(t, ctx, gameID)
		_, err := e.collection.Create(ctx, CreateRequest{GameID: gameID, Name: "exist"})
		assert.NoError(t, err)
		_, err = e.collection.Create(ctx, CreateRequest{GameID: gameID, Name: "exist"})
		assert.ErrorIs(t, err, er.CollectionExist)
	})

	t.Run("same_name_other_game", func(t *testing.T) {
		e := initCollection(t, "collection_create__same_name_other_game")
		gameA := e.createGame(t, ctx, "game_a")
		gameB := e.createGame(t, ctx, "game_b")
		_, err := e.collection.Create(ctx, CreateRequest{GameID: gameA, Name: "shared"})
		assert.NoError(t, err)
		_, err = e.collection.Create(ctx, CreateRequest{GameID: gameB, Name: "shared"})
		assert.NoError(t, err)
	})

	t.Run("bad_name", func(t *testing.T) {
		gameID := "parent_game"
		e := initCollection(t, "collection_create__bad_name")
		e.createGame(t, ctx, gameID)
		_, err := e.collection.Create(ctx, CreateRequest{GameID: gameID, Name: "---"})
		assert.ErrorIs(t, err, er.BadName)
	})

	t.Run("game_not_exist", func(t *testing.T) {
		e := initCollection(t, "collection_create__game_not_exist")
		_, err := e.collection.Create(ctx, CreateRequest{GameID: "missing", Name: "ok"})
		assert.ErrorIs(t, err, er.GameNotExists)
	})
}

func TestCollectionGet(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		gameID := "parent_game"
		name := "success"
		desc := "descrption"
		image := "https://some.url/image"

		e := initCollection(t, "collection_get__success")
		e.createGame(t, ctx, gameID)
		wait, err := e.collection.Create(ctx, CreateRequest{
			GameID:      gameID,
			Name:        name,
			Description: desc,
			Image:       image,
		})
		assert.NoError(t, err)
		got, err := e.collection.Get(ctx, gameID, name)
		assert.NoError(t, err)
		wait.CreatedAt = got.CreatedAt
		wait.UpdatedAt = got.UpdatedAt
		assert.Equal(t, wait, got)
	})

	t.Run("not_exist", func(t *testing.T) {
		gameID := "parent_game"
		e := initCollection(t, "collection_get__not_exist")
		e.createGame(t, ctx, gameID)
		_, err := e.collection.Get(ctx, gameID, "not_exist")
		assert.ErrorIs(t, err, er.CollectionNotExists)
	})

	t.Run("bad_name", func(t *testing.T) {
		gameID := "parent_game"
		e := initCollection(t, "collection_get__bad_name")
		e.createGame(t, ctx, gameID)
		_, err := e.collection.Get(ctx, gameID, "---")
		assert.ErrorIs(t, err, er.BadName)
	})

	t.Run("game_not_exist", func(t *testing.T) {
		e := initCollection(t, "collection_get__game_not_exist")
		_, err := e.collection.Get(ctx, "missing", "ok")
		assert.ErrorIs(t, err, er.GameNotExists)
	})
}

func TestCollectionList(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		gameID := "parent_game"
		name := "success"
		desc := "descrption"
		image := "https://some.url/image"

		e := initCollection(t, "collection_list__success")
		e.createGame(t, ctx, gameID)
		wait, err := e.collection.Create(ctx, CreateRequest{
			GameID:      gameID,
			Name:        name,
			Description: desc,
			Image:       image,
		})
		assert.NoError(t, err)
		got, err := e.collection.List(ctx, gameID)
		assert.NoError(t, err)
		assert.Len(t, got, 1)
		wait.CreatedAt = got[0].CreatedAt
		wait.UpdatedAt = got[0].UpdatedAt
		assert.Equal(t, []*entitiesCollection.Collection{wait}, got)
	})

	t.Run("empty", func(t *testing.T) {
		gameID := "parent_game"
		e := initCollection(t, "collection_list__empty")
		e.createGame(t, ctx, gameID)
		got, err := e.collection.List(ctx, gameID)
		assert.NoError(t, err)
		assert.Equal(t, []*entitiesCollection.Collection(nil), got)
	})

	t.Run("isolated_by_game", func(t *testing.T) {
		e := initCollection(t, "collection_list__isolated_by_game")
		gameA := e.createGame(t, ctx, "game_a")
		gameB := e.createGame(t, ctx, "game_b")
		_, err := e.collection.Create(ctx, CreateRequest{GameID: gameA, Name: "only_a"})
		assert.NoError(t, err)
		got, err := e.collection.List(ctx, gameB)
		assert.NoError(t, err)
		assert.Equal(t, []*entitiesCollection.Collection(nil), got)
	})

	t.Run("game_not_exist", func(t *testing.T) {
		e := initCollection(t, "collection_list__game_not_exist")
		_, err := e.collection.List(ctx, "missing")
		assert.ErrorIs(t, err, er.GameNotExists)
	})
}

func TestCollectionMove(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		gameID := "parent_game"
		oldName := "success_old"
		newName := "success_new"
		desc := "descrption"
		image := "https://some.url/image"

		e := initCollection(t, "collection_move__success")
		e.createGame(t, ctx, gameID)
		oldCollection, err := e.collection.Create(ctx, CreateRequest{
			GameID:      gameID,
			Name:        oldName,
			Description: desc,
			Image:       image,
		})
		assert.NoError(t, err)
		newCollection, err := e.collection.Move(ctx, gameID, oldName, newName)
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
		e.createGame(t, ctx, gameID)
		_, err := e.collection.Move(ctx, gameID, "not_exist", "new_name")
		assert.ErrorIs(t, err, er.CollectionNotExists)
	})

	t.Run("bad_name", func(t *testing.T) {
		gameID := "parent_game"
		oldName := "bad_name_old"
		e := initCollection(t, "collection_move__bad_name")
		e.createGame(t, ctx, gameID)
		_, err := e.collection.Create(ctx, CreateRequest{GameID: gameID, Name: oldName})
		assert.NoError(t, err)
		_, err = e.collection.Move(ctx, gameID, oldName, "---")
		assert.ErrorIs(t, err, er.BadName)
	})

	t.Run("game_not_exist", func(t *testing.T) {
		e := initCollection(t, "collection_move__game_not_exist")
		_, err := e.collection.Move(ctx, "missing", "old", "new")
		assert.ErrorIs(t, err, er.GameNotExists)
	})
}

func TestCollectionUpdate(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		gameID := "parent_game"
		name := "success"
		newDesc := "desc"
		newImage := "img"

		e := initCollection(t, "collection_update__success")
		e.createGame(t, ctx, gameID)
		wait, err := e.collection.Create(ctx, CreateRequest{GameID: gameID, Name: name})
		assert.NoError(t, err)
		got, err := e.collection.Update(ctx, UpdateRequest{
			GameID:      gameID,
			Name:        name,
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
		e.createGame(t, ctx, gameID)
		_, err := e.collection.Update(ctx, UpdateRequest{GameID: gameID, Name: "not_exist"})
		assert.ErrorIs(t, err, er.CollectionNotExists)
	})

	t.Run("bad_name", func(t *testing.T) {
		gameID := "parent_game"
		e := initCollection(t, "collection_update__bad_name")
		e.createGame(t, ctx, gameID)
		_, err := e.collection.Update(ctx, UpdateRequest{GameID: gameID, Name: "---"})
		assert.ErrorIs(t, err, er.BadName)
	})

	t.Run("game_not_exist", func(t *testing.T) {
		e := initCollection(t, "collection_update__game_not_exist")
		_, err := e.collection.Update(ctx, UpdateRequest{GameID: "missing", Name: "ok"})
		assert.ErrorIs(t, err, er.GameNotExists)
	})
}

func TestCollectionDelete(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		gameID := "parent_game"
		name := "success"
		e := initCollection(t, "collection_delete__success")
		e.createGame(t, ctx, gameID)
		_, err := e.collection.Create(ctx, CreateRequest{GameID: gameID, Name: name})
		assert.NoError(t, err)
		err = e.collection.Delete(ctx, gameID, name)
		assert.NoError(t, err)
		_, err = e.collection.Get(ctx, gameID, name)
		assert.ErrorIs(t, err, er.CollectionNotExists)
	})

	t.Run("not_exist", func(t *testing.T) {
		gameID := "parent_game"
		e := initCollection(t, "collection_delete__not_exist")
		e.createGame(t, ctx, gameID)
		err := e.collection.Delete(ctx, gameID, "not_exist")
		assert.ErrorIs(t, err, er.CollectionNotExists)
	})

	t.Run("bad_name", func(t *testing.T) {
		gameID := "parent_game"
		e := initCollection(t, "collection_delete__bad_name")
		e.createGame(t, ctx, gameID)
		err := e.collection.Delete(ctx, gameID, "---")
		assert.ErrorIs(t, err, er.BadName)
	})

	t.Run("game_not_exist", func(t *testing.T) {
		e := initCollection(t, "collection_delete__game_not_exist")
		err := e.collection.Delete(ctx, "missing", "ok")
		assert.ErrorIs(t, err, er.GameNotExists)
	})
}

func TestCollectionImageCreate(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		gameID := "parent_game"
		name := "success"

		e := initCollection(t, "collection_image_create__success")
		e.createGame(t, ctx, gameID)
		_, err := e.collection.Create(ctx, CreateRequest{GameID: gameID, Name: name})
		assert.NoError(t, err)
		err = e.collection.ImageCreate(ctx, gameID, name, img)
		assert.NoError(t, err)
	})

	t.Run("image_exist", func(t *testing.T) {
		gameID := "parent_game"
		name := "image_exist"

		e := initCollection(t, "collection_image_create__image_exist")
		e.createGame(t, ctx, gameID)
		_, err := e.collection.Create(ctx, CreateRequest{GameID: gameID, Name: name})
		assert.NoError(t, err)
		err = e.collection.ImageCreate(ctx, gameID, name, img)
		assert.NoError(t, err)
		err = e.collection.ImageCreate(ctx, gameID, name, img)
		assert.ErrorIs(t, err, er.CollectionImageExist)
	})

	t.Run("collection_not_exist", func(t *testing.T) {
		gameID := "parent_game"
		e := initCollection(t, "collection_image_create__collection_not_exist")
		e.createGame(t, ctx, gameID)
		err := e.collection.ImageCreate(ctx, gameID, "missing", img)
		assert.ErrorIs(t, err, er.CollectionNotExists)
	})

	t.Run("game_not_exist", func(t *testing.T) {
		e := initCollection(t, "collection_image_create__game_not_exist")
		err := e.collection.ImageCreate(ctx, "missing", "ok", img)
		assert.ErrorIs(t, err, er.GameNotExists)
	})
}

func TestCollectionImageGet(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		gameID := "parent_game"
		name := "success"

		e := initCollection(t, "collection_image_get__success")
		e.createGame(t, ctx, gameID)
		_, err := e.collection.Create(ctx, CreateRequest{GameID: gameID, Name: name})
		assert.NoError(t, err)
		err = e.collection.ImageCreate(ctx, gameID, name, img)
		assert.NoError(t, err)
		got, err := e.collection.ImageGet(ctx, gameID, name)
		assert.NoError(t, err)
		assert.Equal(t, img, got)
	})

	t.Run("image_not_exist", func(t *testing.T) {
		gameID := "parent_game"
		name := "image_not_exist"

		e := initCollection(t, "collection_image_get__image_not_exist")
		e.createGame(t, ctx, gameID)
		_, err := e.collection.Create(ctx, CreateRequest{GameID: gameID, Name: name})
		assert.NoError(t, err)
		_, err = e.collection.ImageGet(ctx, gameID, name)
		assert.ErrorIs(t, err, er.CollectionImageNotExists)
	})

	t.Run("collection_not_exist", func(t *testing.T) {
		gameID := "parent_game"
		e := initCollection(t, "collection_image_get__collection_not_exist")
		e.createGame(t, ctx, gameID)
		_, err := e.collection.ImageGet(ctx, gameID, "missing")
		assert.ErrorIs(t, err, er.CollectionNotExists)
	})

	t.Run("game_not_exist", func(t *testing.T) {
		e := initCollection(t, "collection_image_get__game_not_exist")
		_, err := e.collection.ImageGet(ctx, "missing", "ok")
		assert.ErrorIs(t, err, er.GameNotExists)
	})
}

func TestCollectionImageDelete(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		gameID := "parent_game"
		name := "success"

		e := initCollection(t, "collection_image_delete__success")
		e.createGame(t, ctx, gameID)
		_, err := e.collection.Create(ctx, CreateRequest{GameID: gameID, Name: name})
		assert.NoError(t, err)
		err = e.collection.ImageCreate(ctx, gameID, name, img)
		assert.NoError(t, err)
		err = e.collection.ImageDelete(ctx, gameID, name)
		assert.NoError(t, err)
		_, err = e.collection.ImageGet(ctx, gameID, name)
		assert.ErrorIs(t, err, er.CollectionImageNotExists)
	})

	t.Run("image_not_exist", func(t *testing.T) {
		gameID := "parent_game"
		name := "image_not_exist"

		e := initCollection(t, "collection_image_delete__image_not_exist")
		e.createGame(t, ctx, gameID)
		_, err := e.collection.Create(ctx, CreateRequest{GameID: gameID, Name: name})
		assert.NoError(t, err)
		err = e.collection.ImageDelete(ctx, gameID, name)
		assert.ErrorIs(t, err, er.CollectionImageNotExists)
	})

	t.Run("collection_not_exist", func(t *testing.T) {
		gameID := "parent_game"
		e := initCollection(t, "collection_image_delete__collection_not_exist")
		e.createGame(t, ctx, gameID)
		err := e.collection.ImageDelete(ctx, gameID, "missing")
		assert.ErrorIs(t, err, er.CollectionNotExists)
	})

	t.Run("game_not_exist", func(t *testing.T) {
		e := initCollection(t, "collection_image_delete__game_not_exist")
		err := e.collection.ImageDelete(ctx, "missing", "ok")
		assert.ErrorIs(t, err, er.GameNotExists)
	})
}
