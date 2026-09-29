package game

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/HardDie/fsentry"
	"github.com/stretchr/testify/assert"

	"github.com/HardDie/DeckBuilder/internal/config"
	entitiesGame "github.com/HardDie/DeckBuilder/internal/entities/game"
	er "github.com/HardDie/DeckBuilder/internal/errors"
	repositoriesCore "github.com/HardDie/DeckBuilder/internal/repositories/core"
	"github.com/HardDie/DeckBuilder/internal/utils"
)

var (
	img = []byte("some_image")
)

func initGame(t testing.TB, name string) *game {
	// Create temp dir
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

	// Init config with tmp dir
	cfg := config.Get("")
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
			t.Fatal("error drop core", err)
		}
	})

	return New(cfg, db).(*game)
}

func TestGameCreate(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		wait := &entitiesGame.Game{
			Name:        "success",
			Description: "descrption",
			Image:       "https://some.url/image",
		}
		wait.ID = utils.NameToID(wait.Name)

		g := initGame(t, "game_create__success")
		got, err := g.create(CreateRequest{
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
		g := initGame(t, "game_create__exist")
		_, err := g.create(CreateRequest{Name: "exist"})
		assert.NoError(t, err)
		_, err = g.create(CreateRequest{Name: "exist"})
		assert.ErrorIs(t, err, er.GameExist)
	})

	t.Run("bad_name", func(t *testing.T) {
		g := initGame(t, "game_create__bad_name")
		_, err := g.create(CreateRequest{Name: "---"})
		assert.ErrorIs(t, err, er.BadName)
	})
}
func TestGameGet(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		name := "success"
		desc := "descrption"
		img := "https://some.url/image"

		g := initGame(t, "game_get__success")
		wait, err := g.create(CreateRequest{
			Name:        name,
			Description: desc,
			Image:       img,
		})
		assert.NoError(t, err)
		got, err := g.get(name)
		assert.NoError(t, err)
		wait.CreatedAt = got.CreatedAt
		wait.UpdatedAt = got.UpdatedAt
		assert.Equal(t, wait, got)
	})

	t.Run("not_exist", func(t *testing.T) {
		g := initGame(t, "game_get__not_exist")
		_, err := g.get("not_exist")
		assert.ErrorIs(t, err, er.GameNotExists)
	})

	t.Run("bad_name", func(t *testing.T) {
		g := initGame(t, "game_get__bad_name")
		_, err := g.create(CreateRequest{Name: "---"})
		assert.ErrorIs(t, err, er.BadName)
	})
}
func TestGameList(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		name := "success"
		desc := "descrption"
		img := "https://some.url/image"

		g := initGame(t, "game_list__success")
		wait, err := g.create(CreateRequest{
			Name:        name,
			Description: desc,
			Image:       img,
		})
		assert.NoError(t, err)
		got, err := g.list()
		assert.NoError(t, err)
		assert.Len(t, got, 1)
		wait.CreatedAt = got[0].CreatedAt
		wait.UpdatedAt = got[0].UpdatedAt
		assert.Equal(t, []*entitiesGame.Game{wait}, got)
	})

	t.Run("empty", func(t *testing.T) {
		g := initGame(t, "game_list__empty")
		got, err := g.list()
		assert.NoError(t, err)
		assert.Equal(t, []*entitiesGame.Game(nil), got)
	})
}
func TestGameMove(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		oldName := "success_old"
		newName := "success_new"
		desc := "descrption"
		img := "https://some.url/image"

		g := initGame(t, "game_move__success")
		oldGame, err := g.create(CreateRequest{
			Name:        oldName,
			Description: desc,
			Image:       img,
		})
		assert.NoError(t, err)
		newGame, err := g.move(oldName, newName)
		assert.NoError(t, err)

		oldGame.ID = utils.NameToID(newName)
		oldGame.Name = newName
		oldGame.CreatedAt = oldGame.CreatedAt.Truncate(time.Nanosecond)
		oldGame.UpdatedAt = newGame.UpdatedAt
		assert.Equal(t, oldGame, newGame)
	})

	t.Run("not_exist", func(t *testing.T) {
		g := initGame(t, "game_move__not_exist")
		_, err := g.move("not_exist", "new_name")
		assert.ErrorIs(t, err, er.GameNotExists)
	})

	t.Run("bad_name", func(t *testing.T) {
		oldName := "bad_name_old"
		newName := "---"
		g := initGame(t, "game_move__bad_name")
		_, err := g.create(CreateRequest{Name: oldName})
		assert.NoError(t, err)
		_, err = g.move(oldName, newName)
		assert.ErrorIs(t, err, er.BadName)
	})
}
func TestGameUpdate(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		name := "success"
		newDesc := "desc"
		newImage := "img"

		g := initGame(t, "game_update__success")
		wait, err := g.create(CreateRequest{Name: name})
		assert.NoError(t, err)
		got, err := g.update(updateRequest{
			Name:        name,
			Description: newDesc,
			Image:       newImage,
		})
		wait.Description = newDesc
		wait.Image = newImage
		wait.CreatedAt = got.CreatedAt
		wait.UpdatedAt = got.UpdatedAt
		assert.Equal(t, wait, got)
	})

	t.Run("not_exist", func(t *testing.T) {
		g := initGame(t, "game_update__not_exist")
		_, err := g.update(updateRequest{Name: "not_exist"})
		assert.ErrorIs(t, err, er.GameNotExists)
	})

	t.Run("bad_name", func(t *testing.T) {
		g := initGame(t, "game_update__bad_name")
		_, err := g.update(updateRequest{Name: "---"})
		assert.ErrorIs(t, err, er.BadName)
	})
}
func TestGameDelete(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		name := "success"
		g := initGame(t, "game_delete__success")
		_, err := g.create(CreateRequest{Name: name})
		assert.NoError(t, err)
		err = g.delete(name)
		assert.NoError(t, err)
	})

	t.Run("not_exist", func(t *testing.T) {
		g := initGame(t, "game_delete__not_exist")
		err := g.delete("not_exist")
		assert.ErrorIs(t, err, er.GameNotExists)
	})

	t.Run("bad_name", func(t *testing.T) {
		g := initGame(t, "game_delete__bad_name")
		err := g.delete("---")
		assert.ErrorIs(t, err, er.BadName)
	})
}
func TestGameDuplicate(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		srcName := "success_origin"
		dstName := "success_copy"
		g := initGame(t, "game_duplicate__success")
		srcGame, err := g.create(CreateRequest{Name: srcName})
		assert.NoError(t, err)
		srcGame.CreatedAt = srcGame.CreatedAt.Truncate(time.Nanosecond)
		srcGame.UpdatedAt = srcGame.UpdatedAt.Truncate(time.Nanosecond)
		dstGame, err := g.duplicate(srcName, dstName)
		assert.NoError(t, err)
		dstGame.CreatedAt = dstGame.CreatedAt.Truncate(time.Nanosecond)
		dstGame.UpdatedAt = dstGame.UpdatedAt.Truncate(time.Nanosecond)
		assert.NotEqual(t, srcGame, dstGame)
		list, err := g.list()
		assert.NoError(t, err)
		assert.Len(t, list, 2)
		assert.ElementsMatch(t, []*entitiesGame.Game{srcGame, dstGame}, list)
	})

	t.Run("not_exist", func(t *testing.T) {
		g := initGame(t, "game_duplicate__not_exist")
		_, err := g.duplicate("not_exist", "new")
		assert.ErrorIs(t, err, er.GameNotExists)
	})

	t.Run("exist", func(t *testing.T) {
		srcName := "exist_origin"
		dstName := "exist"
		g := initGame(t, "game_duplicate__exist")
		_, err := g.create(CreateRequest{Name: srcName})
		assert.NoError(t, err)
		_, err = g.create(CreateRequest{Name: dstName})
		assert.NoError(t, err)
		_, err = g.duplicate(srcName, dstName)
		assert.ErrorIs(t, err, er.GameExist)
	})

	t.Run("bad_name_1", func(t *testing.T) {
		g := initGame(t, "game_duplicate__bad_name_1")
		_, err := g.duplicate("---", "good")
		assert.ErrorIs(t, err, er.BadName)
	})

	t.Run("bad_name_2", func(t *testing.T) {
		srcName := "good"
		dstName := "---"
		g := initGame(t, "game_duplicate__bad_name_2")
		_, err := g.create(CreateRequest{Name: srcName})
		assert.NoError(t, err)
		_, err = g.duplicate(srcName, dstName)
		assert.ErrorIs(t, err, er.BadName)
	})
}
func TestGameUpdateInfo(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		oldName := "success_old"
		newName := "success_new"

		g := initGame(t, "game_update_info__success")
		created, err := g.create(CreateRequest{Name: oldName})
		assert.NoError(t, err)
		err = g.updateInfo(oldName, newName)
		assert.NoError(t, err)
		got, err := g.get(newName)
		assert.NoError(t, err)
		assert.Equal(t, utils.NameToID(newName), got.ID)
		assert.Equal(t, newName, got.Name)
		assert.Equal(t, created.CreatedAt, got.CreatedAt)
		assert.Equal(t, created.UpdatedAt, got.UpdatedAt)
		_, err = g.get(oldName)
		assert.ErrorIs(t, err, er.GameNotExists)
	})

	t.Run("not_exist", func(t *testing.T) {
		g := initGame(t, "game_update_info__not_exist")
		err := g.updateInfo("not_exist", "good")
		assert.ErrorIs(t, err, fsentry.ErrNotExist)
	})

	t.Run("bad_name_1", func(t *testing.T) {
		g := initGame(t, "game_update_info__bad_name_1")
		err := g.updateInfo("---", "good")
		assert.ErrorIs(t, err, fsentry.ErrBadName)
	})

	t.Run("bad_name_2", func(t *testing.T) {
		name := "good"
		g := initGame(t, "game_update_info__bad_name_2")
		_, err := g.create(CreateRequest{Name: name})
		assert.NoError(t, err)
		err = g.updateInfo(name, "---")
		assert.ErrorIs(t, err, fsentry.ErrBadName)
	})
}
func TestImageCreate(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		name := "success"

		g := initGame(t, "game_image_create__success")
		_, err := g.create(CreateRequest{Name: name})
		assert.NoError(t, err)
		err = g.imageCreate(name, img)
		assert.NoError(t, err)
	})

	t.Run("image_exist", func(t *testing.T) {
		name := "image_exist"

		g := initGame(t, "game_image_create__image_exist")
		_, err := g.create(CreateRequest{Name: name})
		assert.NoError(t, err)
		err = g.imageCreate(name, img)
		assert.NoError(t, err)
		err = g.imageCreate(name, img)
		assert.ErrorIs(t, err, er.GameImageExist)
	})

	t.Run("game_not_exist", func(t *testing.T) {
		name := "game_not_exist"

		g := initGame(t, "game_image_create__game_not_exist")
		err := g.imageCreate(name, img)
		assert.ErrorIs(t, err, er.GameNotExists)
	})
}
func TestImageGet(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		name := "success"

		g := initGame(t, "game_image_get__success")
		_, err := g.create(CreateRequest{Name: name})
		assert.NoError(t, err)
		err = g.imageCreate(name, img)
		assert.NoError(t, err)
		got, err := g.imageGet(name)
		assert.NoError(t, err)
		assert.Equal(t, img, got)
	})

	t.Run("image_not_exist", func(t *testing.T) {
		name := "image_not_exist"

		g := initGame(t, "game_image_get__image_not_exist")
		_, err := g.create(CreateRequest{Name: name})
		assert.NoError(t, err)
		_, err = g.imageGet(name)
		assert.ErrorIs(t, err, er.GameImageNotExists)
	})

	t.Run("game_not_exist", func(t *testing.T) {
		name := "game_not_exist"

		g := initGame(t, "game_image_get__game_not_exist")
		_, err := g.imageGet(name)
		assert.ErrorIs(t, err, er.GameNotExists)
	})
}
func TestImageDelete(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		name := "success"

		g := initGame(t, "game_image_delete__success")
		_, err := g.create(CreateRequest{Name: name})
		assert.NoError(t, err)
		err = g.imageCreate(name, img)
		assert.NoError(t, err)
		err = g.imageDelete(name)
		assert.NoError(t, err)
	})

	t.Run("image_not_exist", func(t *testing.T) {
		name := "image_not_exist"

		g := initGame(t, "game_image_delete__image_not_exist")
		_, err := g.create(CreateRequest{Name: name})
		assert.NoError(t, err)
		err = g.imageDelete(name)
		assert.ErrorIs(t, err, er.GameImageNotExists)
	})

	t.Run("game_not_exist", func(t *testing.T) {
		name := "game_not_exist"

		g := initGame(t, "game_image_delete__game_not_exist")
		err := g.imageDelete(name)
		assert.ErrorIs(t, err, er.GameNotExists)
	})
}

func TestGameLegacyTimestamps(t *testing.T) {
	g := initGame(t, "game_legacy_timestamps")
	created, err := g.create(CreateRequest{Name: "legacy"})
	assert.NoError(t, err)

	path := filepath.Join(g.cfg.Games(), created.ID, ".info.json")
	raw, err := os.ReadFile(path)
	assert.NoError(t, err)
	var info map[string]json.RawMessage
	assert.NoError(t, json.Unmarshal(raw, &info))
	info["createdAt"] = json.RawMessage("null")
	info["updatedAt"] = json.RawMessage("null")
	out, err := json.Marshal(info)
	assert.NoError(t, err)
	assert.NoError(t, os.WriteFile(path, out, 0o644))

	got, err := g.get(created.ID)
	assert.NoError(t, err)
	assert.False(t, got.CreatedAt.IsZero())
	assert.False(t, got.UpdatedAt.IsZero())

	raw, err = os.ReadFile(path)
	assert.NoError(t, err)
	var stored struct {
		CreatedAt *time.Time `json:"createdAt"`
		UpdatedAt *time.Time `json:"updatedAt"`
	}
	assert.NoError(t, json.Unmarshal(raw, &stored))
	assert.Nil(t, stored.CreatedAt)
	assert.Nil(t, stored.UpdatedAt)
}
