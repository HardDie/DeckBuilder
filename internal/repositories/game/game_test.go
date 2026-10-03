package game

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/HardDie/fsentry"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/HardDie/DeckBuilder/internal/config"
	entitiesGame "github.com/HardDie/DeckBuilder/internal/entities/game"
	er "github.com/HardDie/DeckBuilder/internal/errors"
	repositoriesCore "github.com/HardDie/DeckBuilder/internal/repositories/core"
)

var (
	img = []byte("some_image")
)

type gameEnv struct {
	*game
	cfg *config.Config
}

func initGame(t testing.TB) gameEnv {
	t.Helper()
	cfg := config.Get("")
	cfg.SetDataPath(t.TempDir())

	db := fsentry.New(cfg.Data, fsentry.WithPretty(), fsentry.WithNoLockFile())
	require.NoError(t, db.Init())
	require.NoError(t, repositoriesCore.New(db).Init())

	return gameEnv{game: New(db).(*game), cfg: cfg}
}

// writeImage stores raw bytes as the game image, skipping validation.
func (g gameEnv) writeImage(t testing.TB, gameID string, data []byte) {
	t.Helper()
	require.NoError(t, g.db.CreateBinary("image", data, "games", gameID))
}

func (g gameEnv) readImage(t testing.TB, gameID string) []byte {
	t.Helper()
	data, err := g.db.GetBinary("image", nil, "games", gameID)
	require.NoError(t, err)
	return data
}

// The shared storage rules are tested in repositories.TestFolder.
// This checks the game wiring: its errors and entity fields.
func TestGame(t *testing.T) {
	g := initGame(t)

	created, err := g.Create(CreateRequest{Name: "Munchkin", Description: "d"})
	require.NoError(t, err)
	assert.Equal(t, "munchkin", created.ID)
	assert.Equal(t, "Munchkin", created.Name)
	assert.NoError(t, created.ImageError)

	_, err = g.Create(CreateRequest{Name: "Munchkin"})
	assert.ErrorIs(t, err, er.GameExist)

	all, err := g.GetAll()
	require.NoError(t, err)
	require.Len(t, all, 1)
	assert.Equal(t, "d", all[0].Description)

	updated, err := g.Update("munchkin", UpdateRequest{Name: "Munchkin 2", Image: "empty"})
	require.NoError(t, err)
	assert.Equal(t, "munchkin_2", updated.ID)
	assert.Error(t, updated.ImageError, "a bad URL is reported, not applied")

	_, _, err = g.GetImage("munchkin_2")
	assert.ErrorIs(t, err, er.GameImageNotExists)

	require.NoError(t, g.DeleteByID("munchkin_2"))
	_, err = g.GetByID("munchkin_2")
	assert.ErrorIs(t, err, er.GameNotExists)
}

func TestGameDuplicate(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		g := initGame(t)
		srcGame, err := g.Create(CreateRequest{Name: "success_origin"})
		require.NoError(t, err)
		dstGame, err := g.Duplicate("success_origin", DuplicateRequest{Name: "success_copy"})
		require.NoError(t, err)
		assert.Equal(t, "success_copy", dstGame.ID)

		list, err := g.GetAll()
		require.NoError(t, err)
		truncate := func(e *entitiesGame.Game) *entitiesGame.Game {
			e.CreatedAt = e.CreatedAt.Truncate(time.Nanosecond)
			e.UpdatedAt = e.UpdatedAt.Truncate(time.Nanosecond)
			return e
		}
		for _, e := range list {
			truncate(e)
		}
		assert.ElementsMatch(t, []*entitiesGame.Game{truncate(srcGame), truncate(dstGame)}, list)
	})

	t.Run("not_exist", func(t *testing.T) {
		g := initGame(t)
		_, err := g.Duplicate("not_exist", DuplicateRequest{Name: "new"})
		assert.ErrorIs(t, err, er.GameNotExists)
	})

	t.Run("exist", func(t *testing.T) {
		g := initGame(t)
		_, err := g.Create(CreateRequest{Name: "exist_origin"})
		require.NoError(t, err)
		_, err = g.Create(CreateRequest{Name: "exist"})
		require.NoError(t, err)
		_, err = g.Duplicate("exist_origin", DuplicateRequest{Name: "exist"})
		assert.ErrorIs(t, err, er.GameExist)
	})

	t.Run("bad_source_name", func(t *testing.T) {
		g := initGame(t)
		_, err := g.Duplicate("---", DuplicateRequest{Name: "good"})
		assert.ErrorIs(t, err, er.BadName)
	})

	t.Run("bad_target_name", func(t *testing.T) {
		g := initGame(t)
		_, err := g.Create(CreateRequest{Name: "good"})
		require.NoError(t, err)
		_, err = g.Duplicate("good", DuplicateRequest{Name: "---"})
		assert.ErrorIs(t, err, er.BadName)
	})
}

func TestGameExportImport(t *testing.T) {
	t.Run("round_trip", func(t *testing.T) {
		g := initGame(t)
		created, err := g.Create(CreateRequest{Name: "My Game", Description: "desc"})
		require.NoError(t, err)
		_, err = g.Create(CreateRequest{Name: "Other"})
		require.NoError(t, err)
		g.writeImage(t, created.ID, img)

		data, err := g.Export(created.ID)
		require.NoError(t, err)
		zr, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
		require.NoError(t, err)
		assert.NotEmpty(t, zr.File)
		for _, f := range zr.File {
			assert.True(t, strings.HasPrefix(f.Name, created.ID+"/"), f.Name)
			assert.NotContains(t, f.Name, "other")
		}

		dst := initGame(t)
		got, err := dst.Import(data, "")
		require.NoError(t, err)
		assert.Equal(t, created.ID, got.ID)
		assert.Equal(t, created.Name, got.Name)
		assert.Equal(t, created.Description, got.Description)
		assert.True(t, got.CreatedAt.Equal(created.CreatedAt))
		assert.True(t, got.UpdatedAt.Equal(created.UpdatedAt))
		assert.Equal(t, img, dst.readImage(t, created.ID))
		list, err := dst.GetAll()
		require.NoError(t, err)
		assert.Len(t, list, 1)
	})

	t.Run("rename_keeps_source", func(t *testing.T) {
		g := initGame(t)
		created, err := g.Create(CreateRequest{Name: "My Game", Description: "desc"})
		require.NoError(t, err)
		data, err := g.Export(created.ID)
		require.NoError(t, err)

		got, err := g.Import(data, "Other Title")
		require.NoError(t, err)
		assert.Equal(t, "other_title", got.ID)
		assert.Equal(t, "Other Title", got.Name)
		assert.Equal(t, created.Description, got.Description)
		assert.True(t, got.CreatedAt.Equal(created.CreatedAt))
		assert.True(t, got.UpdatedAt.Equal(created.UpdatedAt))

		kept, err := g.GetByID(created.ID)
		require.NoError(t, err)
		assert.Equal(t, created.Description, kept.Description)
	})

	t.Run("existing_destination", func(t *testing.T) {
		g := initGame(t)
		created, err := g.Create(CreateRequest{Name: "My Game", Description: "original"})
		require.NoError(t, err)
		data, err := g.Export(created.ID)
		require.NoError(t, err)
		_, err = g.Update(created.ID, UpdateRequest{Name: created.Name, Description: "edited"})
		require.NoError(t, err)

		_, err = g.Import(data, "")
		assert.ErrorIs(t, err, er.GameExist)
		kept, err := g.GetByID(created.ID)
		require.NoError(t, err)
		assert.Equal(t, "edited", kept.Description)
	})

	t.Run("bad_archive", func(t *testing.T) {
		g := initGame(t)
		_, err := g.Import([]byte("not a zip"), "")
		assert.ErrorIs(t, err, er.BadArchive)
	})

	t.Run("bad_name", func(t *testing.T) {
		g := initGame(t)
		created, err := g.Create(CreateRequest{Name: "My Game"})
		require.NoError(t, err)
		data, err := g.Export(created.ID)
		require.NoError(t, err)
		_, err = g.Import(data, "---")
		assert.ErrorIs(t, err, er.BadName)
	})

	t.Run("not_exist", func(t *testing.T) {
		g := initGame(t)
		_, err := g.Export("missing")
		assert.ErrorIs(t, err, er.GameNotExists)
	})
}

func TestGameLegacyTimestamps(t *testing.T) {
	g := initGame(t)
	created, err := g.Create(CreateRequest{Name: "legacy"})
	require.NoError(t, err)

	path := filepath.Join(g.cfg.Games(), created.ID, ".info.json")
	raw, err := os.ReadFile(path)
	require.NoError(t, err)
	var info map[string]json.RawMessage
	require.NoError(t, json.Unmarshal(raw, &info))
	info["createdAt"] = json.RawMessage("null")
	info["updatedAt"] = json.RawMessage("null")
	out, err := json.Marshal(info)
	require.NoError(t, err)
	require.NoError(t, os.WriteFile(path, out, 0o644))

	got, err := g.GetByID(created.ID)
	require.NoError(t, err)
	assert.False(t, got.CreatedAt.IsZero())
	assert.False(t, got.UpdatedAt.IsZero())

	raw, err = os.ReadFile(path)
	require.NoError(t, err)
	var stored struct {
		CreatedAt *time.Time `json:"createdAt"`
		UpdatedAt *time.Time `json:"updatedAt"`
	}
	require.NoError(t, json.Unmarshal(raw, &stored))
	assert.Nil(t, stored.CreatedAt, "reading does not write timestamps back")
	assert.Nil(t, stored.UpdatedAt)
}
