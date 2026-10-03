package collection

import (
	"testing"

	"github.com/HardDie/fsentry"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/HardDie/DeckBuilder/internal/apperr"
	"github.com/HardDie/DeckBuilder/internal/config"
	repositoriesCore "github.com/HardDie/DeckBuilder/internal/repositories/core"
)

// The shared storage rules are tested in repositories.TestFolder.
// This checks only what the collection adds: its errors and the GameID field.
func TestCollection(t *testing.T) {
	cfg := config.Get("")
	cfg.SetDataPath(t.TempDir())
	db := fsentry.New(cfg.Data, fsentry.WithPretty(), fsentry.WithNoLockFile())
	require.NoError(t, db.Init())
	require.NoError(t, repositoriesCore.New(db).Init())
	_, err := db.CreateFolder[any]("game", nil, "games")
	require.NoError(t, err)

	repo := New(db)

	created, err := repo.Create("game", CreateRequest{Name: "base", Description: "d"})
	require.NoError(t, err)
	assert.Equal(t, "base", created.ID)
	assert.Equal(t, "game", created.GameID)
	assert.NoError(t, created.ImageError)

	_, err = repo.Create("game", CreateRequest{Name: "base"})
	assert.ErrorIs(t, err, apperr.ErrCollectionExists)
	_, err = repo.Create("missing", CreateRequest{Name: "base"})
	assert.ErrorIs(t, err, apperr.ErrGameNotFound)

	got, err := repo.GetByID("game", "base")
	require.NoError(t, err)
	assert.Equal(t, "d", got.Description)
	assert.Equal(t, "game", got.GameID)

	all, err := repo.GetAll("game")
	require.NoError(t, err)
	require.Len(t, all, 1)
	assert.Equal(t, "game", all[0].GameID)

	updated, err := repo.Update("game", "base", UpdateRequest{Name: "dlc", Image: "empty"})
	require.NoError(t, err)
	assert.Equal(t, "dlc", updated.ID)
	assert.Equal(t, "game", updated.GameID)
	assert.Error(t, updated.ImageError, "a bad URL is reported, not applied")

	_, _, err = repo.GetImage("game", "dlc")
	assert.ErrorIs(t, err, apperr.ErrCollectionImageNotFound)

	require.NoError(t, repo.DeleteByID("game", "dlc"))
	_, err = repo.GetByID("game", "dlc")
	assert.ErrorIs(t, err, apperr.ErrCollectionNotFound)
}
