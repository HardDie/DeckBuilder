package game

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/HardDie/fsentry"
	"github.com/HardDie/fsentry/pkg/fsentry_error"
	"github.com/stretchr/testify/assert"

	"github.com/HardDie/DeckBuilder/internal/config"
	dbCore "github.com/HardDie/DeckBuilder/internal/db/core"
	"github.com/HardDie/DeckBuilder/internal/entity"
	er "github.com/HardDie/DeckBuilder/internal/errors"
	"github.com/HardDie/DeckBuilder/internal/utils"
)

var (
	img = []byte("some_image")
)

func initGame(t testing.TB, name string) Game {
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
	cfg := config.Get(false, "")
	cfg.SetDataPath(dir)

	// Init fsentry object
	fs := fsentry.NewFSEntry(cfg.Data, fsentry.WithPretty())

	// Init core directory
	core := dbCore.New(fs)
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

	return New(fs)
}

func TestGameCreate(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		wait := &entity.GameInfo{
			Name:        "success",
			Description: "descrption",
			Image:       "https://some.url/image",
		}
		wait.ID = utils.NameToID(wait.Name)

		g := initGame(t, "game_create__success")
		got, err := g.Create(ctx, wait.Name, wait.Description, wait.Image)
		assert.NoError(t, err)
		wait.CreatedAt = got.CreatedAt
		wait.UpdatedAt = got.UpdatedAt
		assert.Equal(t, wait, got)
	})

	t.Run("exist", func(t *testing.T) {
		g := initGame(t, "game_create__exist")
		_, err := g.Create(ctx, "exist", "", "")
		assert.NoError(t, err)
		_, err = g.Create(ctx, "exist", "", "")
		assert.ErrorIs(t, err, er.GameExist)
	})

	t.Run("bad_name", func(t *testing.T) {
		g := initGame(t, "game_create__bad_name")
		_, err := g.Create(ctx, "---", "", "")
		assert.ErrorIs(t, err, er.BadName)
	})
}
func TestGameGet(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		name := "success"
		desc := "descrption"
		img := "https://some.url/image"

		g := initGame(t, "game_get__success")
		wait, err := g.Create(ctx, name, desc, img)
		assert.NoError(t, err)
		ctx, got, err := g.Get(ctx, name)
		assert.NoError(t, err)
		wait.CreatedAt = got.CreatedAt
		assert.Equal(t, wait, got)
		// TODO: remove context
		gameID := ctx.Value("gameID").(string)
		assert.Equal(t, wait.ID, gameID)
	})

	t.Run("not_exist", func(t *testing.T) {
		g := initGame(t, "game_get__not_exist")
		_, _, err := g.Get(ctx, "not_exist")
		assert.ErrorIs(t, err, er.GameNotExists)
	})

	t.Run("bad_name", func(t *testing.T) {
		g := initGame(t, "game_get__bad_name")
		_, err := g.Create(ctx, "---", "", "")
		assert.ErrorIs(t, err, er.BadName)
	})
}
func TestGameList(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		name := "success"
		desc := "descrption"
		img := "https://some.url/image"

		g := initGame(t, "game_list__success")
		wait, err := g.Create(ctx, name, desc, img)
		assert.NoError(t, err)
		got, err := g.List(ctx)
		assert.NoError(t, err)
		assert.Len(t, got, 1)
		wait.CreatedAt = got[0].CreatedAt
		assert.Equal(t, []*entity.GameInfo{wait}, got)
	})

	t.Run("empty", func(t *testing.T) {
		g := initGame(t, "game_list__empty")
		got, err := g.List(ctx)
		assert.NoError(t, err)
		assert.Equal(t, []*entity.GameInfo(nil), got)
	})
}
func TestGameMove(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		oldName := "success_old"
		newName := "success_new"
		desc := "descrption"
		img := "https://some.url/image"

		g := initGame(t, "game_move__success")
		oldGame, err := g.Create(ctx, oldName, desc, img)
		assert.NoError(t, err)
		newGame, err := g.Move(ctx, oldName, newName)
		assert.NoError(t, err)

		oldGame.ID = utils.NameToID(newName)
		oldGame.Name = newName
		oldGame.CreatedAt = utils.Allocate(oldGame.CreatedAt.Truncate(time.Nanosecond))
		oldGame.UpdatedAt = newGame.UpdatedAt
		assert.Equal(t, oldGame, newGame)
	})

	t.Run("not_exist", func(t *testing.T) {
		g := initGame(t, "game_move__not_exist")
		_, err := g.Move(ctx, "not_exist", "new_name")
		assert.ErrorIs(t, err, er.GameNotExists)
	})

	t.Run("bad_name", func(t *testing.T) {
		oldName := "bad_name_old"
		newName := "---"
		g := initGame(t, "game_move__bad_name")
		_, err := g.Create(ctx, oldName, "", "")
		assert.NoError(t, err)
		_, err = g.Move(ctx, oldName, newName)
		assert.ErrorIs(t, err, er.BadName)
	})
}
func TestGameUpdate(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		name := "success"
		newDesc := "desc"
		newImage := "img"

		g := initGame(t, "game_update__success")
		wait, err := g.Create(ctx, name, "", "")
		assert.NoError(t, err)
		got, err := g.Update(ctx, name, newDesc, newImage)
		wait.Description = newDesc
		wait.Image = newImage
		wait.CreatedAt = got.CreatedAt
		wait.UpdatedAt = got.UpdatedAt
		assert.Equal(t, wait, got)
	})

	t.Run("not_exist", func(t *testing.T) {
		g := initGame(t, "game_update__not_exist")
		_, err := g.Update(ctx, "not_exist", "", "")
		assert.ErrorIs(t, err, er.GameNotExists)
	})

	t.Run("bad_name", func(t *testing.T) {
		g := initGame(t, "game_update__bad_name")
		_, err := g.Update(ctx, "---", "", "")
		assert.ErrorIs(t, err, er.BadName)
	})
}
func TestGameDelete(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		name := "success"
		g := initGame(t, "game_delete__success")
		_, err := g.Create(ctx, name, "", "")
		assert.NoError(t, err)
		err = g.Delete(ctx, name)
		assert.NoError(t, err)
	})

	t.Run("not_exist", func(t *testing.T) {
		g := initGame(t, "game_delete__not_exist")
		err := g.Delete(ctx, "not_exist")
		assert.ErrorIs(t, err, er.GameNotExists)
	})

	t.Run("bad_name", func(t *testing.T) {
		g := initGame(t, "game_delete__bad_name")
		err := g.Delete(ctx, "---")
		assert.ErrorIs(t, err, er.BadName)
	})
}
func TestGameDuplicate(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		srcName := "success_origin"
		dstName := "success_copy"
		g := initGame(t, "game_duplicate__success")
		srcGame, err := g.Create(ctx, srcName, "", "")
		assert.NoError(t, err)
		srcGame.CreatedAt = utils.Allocate(srcGame.CreatedAt.Truncate(time.Nanosecond))
		dstGame, err := g.Duplicate(ctx, srcName, dstName)
		assert.NoError(t, err)
		dstGame.CreatedAt = utils.Allocate(dstGame.CreatedAt.Truncate(time.Nanosecond))
		assert.NotEqual(t, srcGame, dstGame)
		list, err := g.List(ctx)
		assert.NoError(t, err)
		assert.Len(t, list, 2)
		assert.Equal(t, []*entity.GameInfo{srcGame, dstGame}, list)
	})

	t.Run("not_exist", func(t *testing.T) {
		g := initGame(t, "game_duplicate__not_exist")
		_, err := g.Duplicate(ctx, "not_exist", "new")
		assert.ErrorIs(t, err, er.GameNotExists)
	})

	t.Run("exist", func(t *testing.T) {
		srcName := "exist_origin"
		dstName := "exist"
		g := initGame(t, "game_duplicate__exist")
		_, err := g.Create(ctx, srcName, "", "")
		assert.NoError(t, err)
		_, err = g.Create(ctx, dstName, "", "")
		assert.NoError(t, err)
		_, err = g.Duplicate(ctx, srcName, dstName)
		assert.ErrorIs(t, err, er.GameExist)
	})

	t.Run("bad_name_1", func(t *testing.T) {
		g := initGame(t, "game_duplicate__bad_name_1")
		_, err := g.Duplicate(ctx, "---", "good")
		assert.ErrorIs(t, err, er.BadName)
	})

	t.Run("bad_name_2", func(t *testing.T) {
		srcName := "good"
		dstName := "---"
		g := initGame(t, "game_duplicate__bad_name_2")
		_, err := g.Create(ctx, srcName, "", "")
		assert.NoError(t, err)
		_, err = g.Duplicate(ctx, srcName, dstName)
		assert.ErrorIs(t, err, er.BadName)
	})
}
func TestGameUpdateInfo(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		oldName := "success_1"
		newName := "success_2"

		g := initGame(t, "game_update_info__success")
		_, err := g.Create(ctx, oldName, "", "")
		assert.NoError(t, err)
		err = g.UpdateInfo(ctx, oldName, newName)
		assert.NoError(t, err)
		list, err := g.List(ctx)
		assert.NoError(t, err)
		assert.Len(t, list, 0)
		err = g.UpdateInfo(ctx, oldName, oldName)
		assert.NoError(t, err)
		list, err = g.List(ctx)
		assert.NoError(t, err)
		assert.Len(t, list, 1)
	})

	t.Run("not_exist", func(t *testing.T) {
		g := initGame(t, "game_update_info__not_exist")
		err := g.UpdateInfo(ctx, "not_exist", "good")
		assert.ErrorIs(t, err, fsentry_error.ErrorNotExist)
	})

	t.Run("bad_name_1", func(t *testing.T) {
		g := initGame(t, "game_update_info__bad_name_1")
		err := g.UpdateInfo(ctx, "---", "good")
		assert.ErrorIs(t, err, fsentry_error.ErrorBadName)
	})

	t.Run("bad_name_2", func(t *testing.T) {
		name := "good"
		g := initGame(t, "game_update_info__bad_name_2")
		_, err := g.Create(ctx, name, "", "")
		assert.NoError(t, err)
		err = g.UpdateInfo(ctx, name, "---")
		assert.ErrorIs(t, err, fsentry_error.ErrorBadName)
	})
}
func TestImageCreate(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		name := "success"

		g := initGame(t, "game_image_create__success")
		_, err := g.Create(ctx, name, "", "")
		assert.NoError(t, err)
		err = g.ImageCreate(ctx, name, img)
		assert.NoError(t, err)
	})

	t.Run("image_exist", func(t *testing.T) {
		name := "image_exist"

		g := initGame(t, "game_image_create__image_exist")
		_, err := g.Create(ctx, name, "", "")
		assert.NoError(t, err)
		err = g.ImageCreate(ctx, name, img)
		assert.NoError(t, err)
		err = g.ImageCreate(ctx, name, img)
		assert.ErrorIs(t, err, er.GameImageExist)
	})

	t.Run("game_not_exist", func(t *testing.T) {
		name := "game_not_exist"

		g := initGame(t, "game_image_create__game_not_exist")
		err := g.ImageCreate(ctx, name, img)
		assert.ErrorIs(t, err, er.GameNotExists)
	})
}
func TestImageGet(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		name := "success"

		g := initGame(t, "game_image_get__success")
		_, err := g.Create(ctx, name, "", "")
		assert.NoError(t, err)
		err = g.ImageCreate(ctx, name, img)
		assert.NoError(t, err)
		got, err := g.ImageGet(ctx, name)
		assert.NoError(t, err)
		assert.Equal(t, img, got)
	})

	t.Run("image_not_exist", func(t *testing.T) {
		name := "image_not_exist"

		g := initGame(t, "game_image_get__image_not_exist")
		_, err := g.Create(ctx, name, "", "")
		assert.NoError(t, err)
		_, err = g.ImageGet(ctx, name)
		assert.ErrorIs(t, err, er.GameImageNotExists)
	})

	t.Run("game_not_exist", func(t *testing.T) {
		name := "game_not_exist"

		g := initGame(t, "game_image_get__game_not_exist")
		_, err := g.ImageGet(ctx, name)
		assert.ErrorIs(t, err, er.GameNotExists)
	})
}
func TestImageDelete(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		name := "success"

		g := initGame(t, "game_image_delete__success")
		_, err := g.Create(ctx, name, "", "")
		assert.NoError(t, err)
		err = g.ImageCreate(ctx, name, img)
		assert.NoError(t, err)
		err = g.ImageDelete(ctx, name)
		assert.NoError(t, err)
	})

	t.Run("image_not_exist", func(t *testing.T) {
		name := "image_not_exist"

		g := initGame(t, "game_image_delete__image_not_exist")
		_, err := g.Create(ctx, name, "", "")
		assert.NoError(t, err)
		err = g.ImageDelete(ctx, name)
		assert.ErrorIs(t, err, er.GameImageNotExists)
	})

	t.Run("game_not_exist", func(t *testing.T) {
		name := "game_not_exist"

		g := initGame(t, "game_image_delete__game_not_exist")
		err := g.ImageDelete(ctx, name)
		assert.ErrorIs(t, err, er.GameNotExists)
	})
}
