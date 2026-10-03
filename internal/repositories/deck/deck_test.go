package deck

import (
	"bytes"
	"image"
	pngenc "image/png"
	"testing"

	"github.com/HardDie/fsentry"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/HardDie/DeckBuilder/internal/config"
	er "github.com/HardDie/DeckBuilder/internal/errors"
	"github.com/HardDie/DeckBuilder/internal/repositories"
	repositoriesCore "github.com/HardDie/DeckBuilder/internal/repositories/core"
)

func newDeckDB(t *testing.T, collections ...string) *fsentry.DB {
	t.Helper()
	cfg := config.Get("")
	cfg.SetDataPath(t.TempDir())
	db := fsentry.New(cfg.Data, fsentry.WithPretty(), fsentry.WithNoLockFile())
	require.NoError(t, db.Init())
	require.NoError(t, repositoriesCore.New(db).Init())
	_, err := db.CreateFolder[any]("game", nil, "games")
	require.NoError(t, err)
	for _, c := range collections {
		_, err = db.CreateFolder[any](c, nil, "games", "game")
		require.NoError(t, err)
	}
	return db
}

// The shared storage rules are tested in repositories.TestFolder.
// This checks only what the deck adds: its errors, its parent fields,
// and the "cards" folder.
func TestDeck(t *testing.T) {
	db := newDeckDB(t, "base")
	repo := New(db)

	created, err := repo.Create("game", "base", CreateRequest{Name: "monster", Description: "d"})
	require.NoError(t, err)
	assert.Equal(t, "monster", created.ID)
	assert.Equal(t, "game", created.GameID)
	assert.Equal(t, "base", created.CollectionID)
	assert.NoError(t, created.ImageError)

	_, err = db.GetFolder[any]("cards", "games", "game", "base", "monster")
	assert.NoError(t, err, "create makes the cards folder")

	_, err = repo.Create("game", "base", CreateRequest{Name: "monster"})
	assert.ErrorIs(t, err, er.DeckExist)
	_, err = repo.Create("game", "missing", CreateRequest{Name: "monster"})
	assert.ErrorIs(t, err, er.CollectionNotExists)

	got, err := repo.GetByID("game", "base", "monster")
	require.NoError(t, err)
	assert.Equal(t, "d", got.Description)
	assert.Equal(t, "base", got.CollectionID)

	all, err := repo.GetAll("game", "base")
	require.NoError(t, err)
	require.Len(t, all, 1)
	assert.Equal(t, "base", all[0].CollectionID)

	updated, err := repo.Update("game", "base", "monster", UpdateRequest{Name: "monsters", Image: "empty"})
	require.NoError(t, err)
	assert.Equal(t, "monsters", updated.ID)
	assert.Equal(t, "base", updated.CollectionID)
	assert.Error(t, updated.ImageError, "a bad URL is reported, not applied")
	_, err = db.GetFolder[any]("cards", "games", "game", "base", "monsters")
	assert.NoError(t, err, "the cards folder moves with the rename")

	_, _, err = repo.GetImage("game", "base", "monsters")
	assert.ErrorIs(t, err, er.DeckImageNotExists)

	require.NoError(t, repo.DeleteByID("game", "base", "monsters"))
	_, err = repo.GetByID("game", "base", "monsters")
	assert.ErrorIs(t, err, er.DeckNotExists)
}

func TestGetAllDecksInGame(t *testing.T) {
	t.Run("unique_by_name_and_image", func(t *testing.T) {
		db := newDeckDB(t, "collection_a", "collection_b")
		// Write the folders directly: the image URLs are kept without a download.
		write := func(collection, name, image string) {
			_, err := db.CreateFolder(name, repositories.FolderModel{
				Image: fsentry.QuotedString(image),
			}, "games", "game", collection)
			require.NoError(t, err)
		}
		write("collection_a", "shared", "https://img")
		write("collection_b", "shared", "https://img")
		write("collection_b", "shared_other_image", "https://other")
		write("collection_b", "other", "")

		got, err := New(db).GetAllDecksInGame("game")
		require.NoError(t, err)
		var names []string
		for _, d := range got {
			names = append(names, d.Name)
		}
		assert.ElementsMatch(t, []string{"shared", "shared_other_image", "other"}, names)
	})

	t.Run("no_collections", func(t *testing.T) {
		db := newDeckDB(t)
		got, err := New(db).GetAllDecksInGame("game")
		require.NoError(t, err)
		assert.NotNil(t, got, "an empty list, not nil")
		assert.Empty(t, got)
	})

	t.Run("game_not_exist", func(t *testing.T) {
		db := newDeckDB(t)
		_, err := New(db).GetAllDecksInGame("missing")
		assert.ErrorIs(t, err, er.GameNotExists)
	})
}

func TestDeckHasImage(t *testing.T) {
	db := newDeckDB(t, "base")
	repo := New(db)
	var png bytes.Buffer
	require.NoError(t, pngenc.Encode(&png, image.NewRGBA(image.Rect(0, 0, 2, 2))))

	plain, err := repo.Create("game", "base", CreateRequest{Name: "plain"})
	require.NoError(t, err)
	assert.False(t, plain.HasImage)
	withBack, err := repo.Create("game", "base", CreateRequest{Name: "with_back", ImageFile: png.Bytes()})
	require.NoError(t, err)
	assert.True(t, withBack.HasImage)

	all, err := repo.GetAll("game", "base")
	require.NoError(t, err)
	got := map[string]bool{}
	for _, d := range all {
		got[d.ID] = d.HasImage
	}
	assert.Equal(t, map[string]bool{"plain": false, "with_back": true}, got)
}

func TestDeckCardsMissingImage(t *testing.T) {
	db := newDeckDB(t, "base")
	repo := New(db)
	_, err := repo.Create("game", "base", CreateRequest{Name: "crew"})
	require.NoError(t, err)
	cardsPath := []string{"games", "game", "base", "crew"}
	missing := func() bool {
		t.Helper()
		got, err := repo.GetByID("game", "base", "crew")
		require.NoError(t, err)
		return got.CardsMissingImage
	}

	assert.False(t, missing(), "no cards")

	// Two cards in the list, an image file for card 1 only.
	_, err = db.UpdateFolder("cards", map[string]any{"1": map[string]any{}, "2": map[string]any{}}, cardsPath...)
	require.NoError(t, err)
	require.NoError(t, db.CreateBinary("1", []byte("png"), append(cardsPath, "cards")...))
	assert.True(t, missing(), "card 2 has no image")

	require.NoError(t, db.CreateBinary("2", []byte("png"), append(cardsPath, "cards")...))
	assert.False(t, missing(), "every card has an image")
}
