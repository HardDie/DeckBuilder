package repositories

import (
	"bytes"
	"image"
	"image/jpeg"
	"image/png"
	"net/http"
	"net/http/httptest"
	"sort"
	"testing"

	"github.com/HardDie/fsentry"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/HardDie/DeckBuilder/internal/apperr"
)

// folderLevel is one catalog level for the shared Folder tests.
type folderLevel struct {
	name string
	errs FolderErrors
	// parents creates the parent folders and returns the path below "games".
	parents func(t *testing.T, db *fsentry.DB) []string
	// sibling creates a second parent at the same depth; nil for a game.
	sibling func(t *testing.T, db *fsentry.DB) []string
	// missing is a parent path that does not exist; nil for a game.
	missing    []string
	missingErr *apperr.Error
}

func folderLevels() []folderLevel {
	mkdir := func(t *testing.T, db *fsentry.DB, name string, path ...string) {
		t.Helper()
		_, err := db.CreateFolder[any](name, nil, path...)
		require.NoError(t, err)
	}
	return []folderLevel{
		{
			name: "game",
			errs: FolderErrors{
				Exist: apperr.ErrGameExists, NotExist: apperr.ErrGameNotFound,
				ImageExist: apperr.ErrGameImageExists, ImageNotExist: apperr.ErrGameImageNotFound,
			},
			parents: func(*testing.T, *fsentry.DB) []string { return nil },
		},
		{
			name: "collection",
			errs: FolderErrors{
				Exist: apperr.ErrCollectionExists, NotExist: apperr.ErrCollectionNotFound,
				ImageExist: apperr.ErrCollectionImageExists, ImageNotExist: apperr.ErrCollectionImageNotFound,
			},
			parents: func(t *testing.T, db *fsentry.DB) []string {
				mkdir(t, db, "g", gamesPath)
				return []string{"g"}
			},
			sibling: func(t *testing.T, db *fsentry.DB) []string {
				mkdir(t, db, "g2", gamesPath)
				return []string{"g2"}
			},
			missing:    []string{"nope"},
			missingErr: apperr.ErrGameNotFound,
		},
		{
			name: "deck",
			errs: FolderErrors{
				Exist: apperr.ErrDeckExists, NotExist: apperr.ErrDeckNotFound,
				ImageExist: apperr.ErrDeckImageExists, ImageNotExist: apperr.ErrDeckImageNotFound,
			},
			parents: func(t *testing.T, db *fsentry.DB) []string {
				mkdir(t, db, "g", gamesPath)
				mkdir(t, db, "c", gamesPath, "g")
				return []string{"g", "c"}
			},
			sibling: func(t *testing.T, db *fsentry.DB) []string {
				mkdir(t, db, "c2", gamesPath, "g")
				return []string{"g", "c2"}
			},
			missing:    []string{"g", "nope"},
			missingErr: apperr.ErrCollectionNotFound,
		},
	}
}

func encodedImage(t *testing.T, encode func(*bytes.Buffer, image.Image) error) []byte {
	t.Helper()
	var buf bytes.Buffer
	require.NoError(t, encode(&buf, image.NewRGBA(image.Rect(0, 0, 4, 4))))
	return buf.Bytes()
}

func TestFolder(t *testing.T) {
	pngBytes := encodedImage(t, func(b *bytes.Buffer, i image.Image) error { return png.Encode(b, i) })
	jpegBytes := encodedImage(t, func(b *bytes.Buffer, i image.Image) error { return jpeg.Encode(b, i, nil) })
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/ok.png" {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		_, _ = w.Write(pngBytes)
	}))
	defer srv.Close()
	okURL := srv.URL + "/ok.png"

	cases := []struct {
		name string
		run  func(t *testing.T, f *Folder, parent []string, lv folderLevel)
	}{
		{"create_and_get", func(t *testing.T, f *Folder, parent []string, _ folderLevel) {
			saved, err := f.Create(parent, FolderWrite{Name: "alpha", Description: "d"})
			require.NoError(t, err)
			assert.Equal(t, "alpha", saved.Info.ID)
			assert.NoError(t, saved.ImageError)

			got, err := f.Get(parent, "alpha")
			require.NoError(t, err)
			assert.Equal(t, "d", got.Data.Description.String())
		}},
		{"create_exists", func(t *testing.T, f *Folder, parent []string, lv folderLevel) {
			_, err := f.Create(parent, FolderWrite{Name: "alpha"})
			require.NoError(t, err)
			_, err = f.Create(parent, FolderWrite{Name: "alpha"})
			assert.ErrorIs(t, err, lv.errs.Exist)
		}},
		{"get_not_exist", func(t *testing.T, f *Folder, parent []string, lv folderLevel) {
			_, err := f.Get(parent, "missing")
			assert.ErrorIs(t, err, lv.errs.NotExist)
		}},
		{"missing_parent", func(t *testing.T, f *Folder, _ []string, lv folderLevel) {
			if lv.missing == nil {
				t.Skip("a game has no parent")
			}
			_, err := f.Create(lv.missing, FolderWrite{Name: "alpha"})
			assert.ErrorIs(t, err, lv.missingErr)
			_, err = f.Get(lv.missing, "alpha")
			assert.ErrorIs(t, err, lv.missingErr)
			_, err = f.List(lv.missing)
			assert.ErrorIs(t, err, lv.missingErr)
			_, err = f.Update(lv.missing, "alpha", FolderWrite{Name: "alpha"})
			assert.ErrorIs(t, err, lv.missingErr)
			assert.ErrorIs(t, f.Delete(lv.missing, "alpha"), lv.missingErr)
			_, _, err = f.Image(lv.missing, "alpha")
			assert.ErrorIs(t, err, lv.missingErr)
		}},
		{"bad_name", func(t *testing.T, f *Folder, parent []string, _ folderLevel) {
			_, err := f.Create(parent, FolderWrite{Name: "---"})
			assert.ErrorIs(t, err, apperr.ErrBadName)
			_, err = f.Get(parent, "---")
			assert.ErrorIs(t, err, apperr.ErrBadName)
			assert.ErrorIs(t, f.Delete(parent, "---"), apperr.ErrBadName)

			_, err = f.Create(parent, FolderWrite{Name: "alpha"})
			require.NoError(t, err)
			_, err = f.Update(parent, "alpha", FolderWrite{Name: "---"})
			assert.ErrorIs(t, err, apperr.ErrBadName, "rename to a bad name")
		}},
		{"list_empty", func(t *testing.T, f *Folder, parent []string, _ folderLevel) {
			infos, err := f.List(parent)
			require.NoError(t, err)
			assert.Empty(t, infos)
		}},
		{"isolated_by_parent", func(t *testing.T, f *Folder, parent []string, lv folderLevel) {
			if lv.sibling == nil {
				t.Skip("a game has no parent")
			}
			sibling := lv.sibling(t, f.db)
			_, err := f.Create(parent, FolderWrite{Name: "shared"})
			require.NoError(t, err)
			_, err = f.Create(sibling, FolderWrite{Name: "shared"})
			assert.NoError(t, err, "the same name under another parent")

			_, err = f.Create(parent, FolderWrite{Name: "only_here"})
			require.NoError(t, err)
			infos, err := f.List(sibling)
			require.NoError(t, err)
			require.Len(t, infos, 1)
			assert.Equal(t, "shared", infos[0].ID)
		}},
		{"image_exist", func(t *testing.T, f *Folder, parent []string, lv folderLevel) {
			_, err := f.Create(parent, FolderWrite{Name: "alpha", ImageFile: pngBytes})
			require.NoError(t, err)
			assert.ErrorIs(t, f.imageCreate(parent, "alpha", pngBytes), lv.errs.ImageExist)
			_, err = f.imageGet(parent, "missing")
			assert.ErrorIs(t, err, lv.errs.NotExist, "image of a missing entity")
			assert.ErrorIs(t, f.imageDelete(parent, "missing"), lv.errs.NotExist)
		}},
		{"list", func(t *testing.T, f *Folder, parent []string, _ folderLevel) {
			for _, name := range []string{"beta", "alpha"} {
				_, err := f.Create(parent, FolderWrite{Name: name})
				require.NoError(t, err)
			}
			infos, err := f.List(parent)
			require.NoError(t, err)
			var ids []string
			for _, info := range infos {
				ids = append(ids, info.ID)
			}
			sort.Strings(ids)
			assert.Equal(t, []string{"alpha", "beta"}, ids)
		}},
		{"update_renames_and_saves", func(t *testing.T, f *Folder, parent []string, lv folderLevel) {
			_, err := f.Create(parent, FolderWrite{Name: "alpha", ImageFile: pngBytes})
			require.NoError(t, err)

			saved, err := f.Update(parent, "alpha", FolderWrite{Name: "gamma", Description: "x"})
			require.NoError(t, err)
			assert.Equal(t, "gamma", saved.Info.ID)
			assert.Equal(t, "x", saved.Info.Data.Description.String())

			_, err = f.Get(parent, "alpha")
			assert.ErrorIs(t, err, lv.errs.NotExist)
			_, imgType, err := f.Image(parent, "gamma")
			require.NoError(t, err)
			assert.Equal(t, "png", imgType, "the image moves with the rename")
		}},
		{"update_not_exist", func(t *testing.T, f *Folder, parent []string, lv folderLevel) {
			_, err := f.Update(parent, "missing", FolderWrite{Name: "missing"})
			assert.ErrorIs(t, err, lv.errs.NotExist)
		}},
		{"delete", func(t *testing.T, f *Folder, parent []string, lv folderLevel) {
			_, err := f.Create(parent, FolderWrite{Name: "alpha"})
			require.NoError(t, err)
			require.NoError(t, f.Delete(parent, "alpha"))
			_, err = f.Get(parent, "alpha")
			assert.ErrorIs(t, err, lv.errs.NotExist)
			assert.ErrorIs(t, f.Delete(parent, "alpha"), lv.errs.NotExist)
		}},
		{"image_lifecycle", func(t *testing.T, f *Folder, parent []string, lv folderLevel) {
			_, _, err := f.Image(parent, "alpha")
			assert.ErrorIs(t, err, lv.errs.NotExist, "no entity yet")

			_, err = f.Create(parent, FolderWrite{Name: "alpha", ImageFile: pngBytes})
			require.NoError(t, err)
			_, imgType, err := f.Image(parent, "alpha")
			require.NoError(t, err)
			assert.Equal(t, "png", imgType)

			// A new file replaces the image.
			_, err = f.Update(parent, "alpha", FolderWrite{Name: "alpha", ImageFile: jpegBytes})
			require.NoError(t, err)
			_, imgType, err = f.Image(parent, "alpha")
			require.NoError(t, err)
			assert.Equal(t, "jpeg", imgType)

			// A bad file is not applied; the old image stays.
			saved, err := f.Update(parent, "alpha", FolderWrite{Name: "alpha", ImageFile: []byte("text")})
			require.NoError(t, err)
			assert.ErrorIs(t, saved.ImageError, apperr.ErrUnsupportedImage)
			_, imgType, err = f.Image(parent, "alpha")
			require.NoError(t, err)
			assert.Equal(t, "jpeg", imgType)

			// A URL replaces it and is stored.
			saved, err = f.Update(parent, "alpha", FolderWrite{Name: "alpha", Image: okURL})
			require.NoError(t, err)
			assert.Equal(t, okURL, saved.Info.Data.Image.String())
			_, imgType, err = f.Image(parent, "alpha")
			require.NoError(t, err)
			assert.Equal(t, "png", imgType)

			// Clearing the URL removes the image.
			_, err = f.Update(parent, "alpha", FolderWrite{Name: "alpha"})
			require.NoError(t, err)
			_, _, err = f.Image(parent, "alpha")
			assert.ErrorIs(t, err, lv.errs.ImageNotExist)
		}},
	}

	for _, lv := range folderLevels() {
		for _, tc := range cases {
			t.Run(lv.name+"/"+tc.name, func(t *testing.T) {
				db := fsentry.New(t.TempDir(), fsentry.WithPretty(), fsentry.WithNoLockFile())
				require.NoError(t, db.Init())
				_, err := db.CreateFolder[any](gamesPath, nil)
				require.NoError(t, err)

				parent := lv.parents(t, db)
				tc.run(t, NewFolder(db, lv.errs), parent, lv)
			})
		}
	}
}
