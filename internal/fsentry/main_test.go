package fsentry

import (
	"errors"
	"path"
	"testing"

	"github.com/HardDie/DeckBuilder/internal/fsentry/pkg/fsentry_error"
)

func TestFolder(t *testing.T) {
	folderDB := NewFSEntry("test")
	err := folderDB.Init()
	if err != nil {
		t.Fatal(err)
	}
	defer folderDB.Drop()

	t.Run("create", func(t *testing.T) {
		t.Parallel()

		db := NewFSEntry(path.Join("test", "test_folder_create"))
		err := db.Init()
		if err != nil {
			t.Fatal(err)
		}

		// Try to create folder with empty name
		_, err = db.CreateFolder("", nil)
		if !errors.Is(err, fsentry_error.ErrorBadName) {
			t.Fatal("Bad name")
		}

		// Try to create folder in not exist subdirectory
		_, err = db.CreateFolder("bad_path", nil, "not_exist_folder")
		if !errors.Is(err, fsentry_error.ErrorBadPath) {
			t.Fatal("Bad path for folder")
		}

		// Create directory
		_, err = db.CreateFolder("some_folder", nil)
		if err != nil {
			t.Fatal(err)
		}

		// Create subdirectory
		_, err = db.CreateFolder("some_inner_folder", nil, "some_folder")
		if err != nil {
			t.Fatal(err)
		}

		// Try to create duplicate
		_, err = db.CreateFolder("some_folder", nil)
		if !errors.Is(err, fsentry_error.ErrorExist) {
			t.Fatal("Folder already exist")
		}

		err = db.Drop()
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("get", func(t *testing.T) {
		t.Parallel()

		db := NewFSEntry(path.Join("test", "test_folder_get"))
		err := db.Init()
		if err != nil {
			t.Fatal(err)
		}

		// Try to get bad name
		_, err = db.GetFolder("")
		if !errors.Is(err, fsentry_error.ErrorBadName) {
			t.Fatal("Bad name")
		}

		// Try to get not exist folder
		_, err = db.GetFolder("not_exist_folder")
		if !errors.Is(err, fsentry_error.ErrorNotExist) {
			t.Fatal("Folder not exist")
		}

		// Try to get from bad path
		_, err = db.GetFolder("not_exist_folder", "bad_path")
		if !errors.Is(err, fsentry_error.ErrorBadPath) {
			t.Fatal("Bad path")
		}

		_, err = db.CreateFolder("some_folder", nil)
		if err != nil {
			t.Fatal(err)
		}

		folder, err := db.GetFolder("some_folder")
		if err != nil {
			t.Fatal(err)
		}

		if folder.Name != "some_folder" {
			t.Fatal("Bad name")
		}

		err = db.Drop()
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("move", func(t *testing.T) {
		t.Parallel()

		db := NewFSEntry(path.Join("test", "test_folder_move"))
		err := db.Init()
		if err != nil {
			t.Fatal(err)
		}

		// Try to move not exist folder
		_, err = db.MoveFolder("not_exist", "new_not_exist")
		if !errors.Is(err, fsentry_error.ErrorNotExist) {
			t.Fatal("Folder not exist")
		}

		// Try to move bad path
		_, err = db.MoveFolder("not_exist", "new_not_exist", "bad_path")
		if !errors.Is(err, fsentry_error.ErrorBadPath) {
			t.Fatal("Bad path")
		}

		_, err = db.CreateFolder("first_folder", nil)
		if err != nil {
			t.Fatal(err)
		}

		_, err = db.CreateFolder("second_folder", nil)
		if err != nil {
			t.Fatal(err)
		}

		// Try to move bad name
		_, err = db.MoveFolder("", "new_not_exist")
		if !errors.Is(err, fsentry_error.ErrorBadName) {
			t.Fatal("Bad name")
		}

		// Try to move bad name
		_, err = db.MoveFolder("first_folder", "")
		if !errors.Is(err, fsentry_error.ErrorBadName) {
			t.Fatal("Bad name")
		}

		// Try to move into exist folder
		_, err = db.MoveFolder("first_folder", "second_folder")
		if !errors.Is(err, fsentry_error.ErrorExist) {
			t.Fatal("Folder already exist")
		}

		_, err = db.MoveFolder("first_folder", "new_first_folder")
		if err != nil {
			t.Fatal(err)
		}

		err = db.Drop()
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("update", func(t *testing.T) {
		t.Parallel()

		db := NewFSEntry(path.Join("test", "test_folder_update"))
		err := db.Init()
		if err != nil {
			t.Fatal(err)
		}

		// Try to update bad name
		_, err = db.UpdateFolder("", nil)
		if !errors.Is(err, fsentry_error.ErrorBadName) {
			t.Fatal("Bad name")
		}

		// Try to update not exist folder
		_, err = db.UpdateFolder("not_exist_folder", nil)
		if !errors.Is(err, fsentry_error.ErrorNotExist) {
			t.Fatal("Folder not exist")
		}

		// Try to update bad path folder
		_, err = db.UpdateFolder("not_exist_folder", nil, "bad_path")
		if !errors.Is(err, fsentry_error.ErrorBadPath) {
			t.Fatal("Bad path")
		}

		_, err = db.CreateFolder("some_folder", 5)
		if err != nil {
			t.Fatal(err)
		}

		_, err = db.UpdateFolder("some_folder", 15)
		if err != nil {
			t.Fatal(err)
		}

		err = db.Drop()
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("remove", func(t *testing.T) {
		t.Parallel()

		db := NewFSEntry(path.Join("test", "test_folder_remove"))
		err := db.Init()
		if err != nil {
			t.Fatal(err)
		}

		// Try to remove bad name
		err = db.RemoveFolder("")
		if !errors.Is(err, fsentry_error.ErrorBadName) {
			t.Fatal("Bad name")
		}

		// Try to remove not exist folder
		err = db.RemoveFolder("not_exist_folder")
		if !errors.Is(err, fsentry_error.ErrorNotExist) {
			t.Fatal("Folder not exist")
		}

		_, err = db.CreateFolder("some_folder", nil)
		if err != nil {
			t.Fatal(err)
		}

		err = db.RemoveFolder("some_folder")
		if err != nil {
			t.Fatal(err)
		}

		// Try to remove twice
		err = db.RemoveFolder("some_folder")
		if !errors.Is(err, fsentry_error.ErrorNotExist) {
			t.Fatal("Folder not exist")
		}

		err = db.Drop()
		if err != nil {
			t.Fatal(err)
		}
	})
}

func TestEntry(t *testing.T) {
	folderDB := NewFSEntry("test")
	err := folderDB.Init()
	if err != nil {
		t.Fatal(err)
	}
	defer folderDB.Drop()

	t.Run("create", func(t *testing.T) {
		t.Parallel()

		db := NewFSEntry(path.Join("test", "test_entry_create"))
		err := db.Init()
		if err != nil {
			t.Fatal(err)
		}

		// Try to create entry with bad name
		err = db.CreateEntry("", nil)
		if !errors.Is(err, fsentry_error.ErrorBadName) {
			t.Fatal("Bad name")
		}

		// Try to create entry in not exist subdirectory
		err = db.CreateEntry("bad_path", nil, "bad")
		if !errors.Is(err, fsentry_error.ErrorBadPath) {
			t.Fatal("Bad path for folder")
		}

		// Create entry
		err = db.CreateEntry("some_entry", nil)
		if err != nil {
			t.Fatal(err)
		}

		// Try to create duplicate
		err = db.CreateEntry("some_entry", nil)
		if !errors.Is(err, fsentry_error.ErrorExist) {
			t.Fatal("Entry already exist")
		}

		err = db.Drop()
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("get", func(t *testing.T) {
		t.Parallel()

		db := NewFSEntry(path.Join("test", "test_entry_get"))
		err := db.Init()
		if err != nil {
			t.Fatal(err)
		}

		// Try to get bad name
		_, err = db.GetEntry("")
		if !errors.Is(err, fsentry_error.ErrorBadName) {
			t.Fatal("Bad name")
		}

		// Try to get not exist entry
		_, err = db.GetEntry("not_exist_entry")
		if !errors.Is(err, fsentry_error.ErrorNotExist) {
			t.Fatal("Entry not exist")
		}

		// Try to get from bad path
		_, err = db.GetEntry("not_exist_entry", "bad_path")
		if !errors.Is(err, fsentry_error.ErrorBadPath) {
			t.Fatal("Bad path")
		}

		err = db.CreateEntry("some_entry", nil)
		if err != nil {
			t.Fatal(err)
		}

		entry, err := db.GetEntry("some_entry")
		if err != nil {
			t.Fatal(err)
		}

		if entry.Name != "some_entry" {
			t.Fatal("Bad name")
		}

		err = db.Drop()
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("move", func(t *testing.T) {
		t.Parallel()

		db := NewFSEntry(path.Join("test", "test_entry_move"))
		err := db.Init()
		if err != nil {
			t.Fatal(err)
		}

		// Try to move not exist entry
		err = db.MoveEntry("not_exist", "new_not_exist")
		if !errors.Is(err, fsentry_error.ErrorNotExist) {
			t.Fatal("Entry not exist")
		}

		// Try to move bad path
		err = db.MoveEntry("not_exist", "new_not_exist", "bad_path")
		if !errors.Is(err, fsentry_error.ErrorBadPath) {
			t.Fatal("Bad path")
		}

		err = db.CreateEntry("first_entry", nil)
		if err != nil {
			t.Fatal(err)
		}

		err = db.CreateEntry("second_entry", nil)
		if err != nil {
			t.Fatal(err)
		}

		// Try to move bad name
		err = db.MoveEntry("", "new_not_exist")
		if !errors.Is(err, fsentry_error.ErrorBadName) {
			t.Fatal("Bad name")
		}

		// Try to move bad name
		err = db.MoveEntry("first_entry", "")
		if !errors.Is(err, fsentry_error.ErrorBadName) {
			t.Fatal("Bad name")
		}

		// Try to move into exist folder
		err = db.MoveEntry("first_entry", "second_entry")
		if !errors.Is(err, fsentry_error.ErrorExist) {
			t.Fatal("Entry already exist")
		}

		err = db.MoveEntry("first_entry", "new_first_entry")
		if err != nil {
			t.Fatal(err)
		}

		err = db.Drop()
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("update", func(t *testing.T) {
		t.Parallel()

		db := NewFSEntry(path.Join("test", "test_entry_update"))
		err := db.Init()
		if err != nil {
			t.Fatal(err)
		}

		// Try to update bad name
		err = db.UpdateEntry("", nil)
		if !errors.Is(err, fsentry_error.ErrorBadName) {
			t.Fatal("Bad name")
		}

		// Try to update not exist entry
		err = db.UpdateEntry("not_exist_entry", nil)
		if !errors.Is(err, fsentry_error.ErrorNotExist) {
			t.Fatal("Entry not exist")
		}

		// Try to update bad path entry
		err = db.UpdateEntry("not_exist_entry", nil, "bad_path")
		if !errors.Is(err, fsentry_error.ErrorBadPath) {
			t.Fatal("Bad path")
		}

		err = db.CreateEntry("some_entry", 5)
		if err != nil {
			t.Fatal(err)
		}

		err = db.UpdateEntry("some_entry", 15)
		if err != nil {
			t.Fatal(err)
		}

		err = db.Drop()
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("remove", func(t *testing.T) {
		t.Parallel()

		db := NewFSEntry(path.Join("test", "test_entry_remove"))
		err := db.Init()
		if err != nil {
			t.Fatal(err)
		}

		// Try to remove bad name
		err = db.RemoveEntry("")
		if !errors.Is(err, fsentry_error.ErrorBadName) {
			t.Fatal("Bad name")
		}

		// Try to remove not exist entry
		err = db.RemoveEntry("not_exist_entry")
		if !errors.Is(err, fsentry_error.ErrorNotExist) {
			t.Fatal("Entry not exist")
		}

		err = db.CreateEntry("some_entry", nil)
		if err != nil {
			t.Fatal(err)
		}

		err = db.RemoveEntry("some_entry")
		if err != nil {
			t.Fatal(err)
		}

		// Try to remove twice
		err = db.RemoveEntry("some_entry")
		if !errors.Is(err, fsentry_error.ErrorNotExist) {
			t.Fatal("Entry not exist")
		}

		err = db.Drop()
		if err != nil {
			t.Fatal(err)
		}
	})
}

func TestBinary(t *testing.T) {
	folderDB := NewFSEntry("test")
	err := folderDB.Init()
	if err != nil {
		t.Fatal(err)
	}
	defer folderDB.Drop()

	t.Run("create", func(t *testing.T) {
		t.Parallel()

		db := NewFSEntry(path.Join("test", "test_binary_create"))
		err := db.Init()
		if err != nil {
			t.Fatal(err)
		}

		// Try to create entry with bad name
		err = db.CreateBinary("", []byte("data"))
		if !errors.Is(err, fsentry_error.ErrorBadName) {
			t.Fatal("Bad name")
		}

		// Try to create entry in not exist subdirectory
		err = db.CreateBinary("bad_path", []byte("data"), "bad")
		if !errors.Is(err, fsentry_error.ErrorBadPath) {
			t.Fatal("Bad path for folder")
		}

		// Create binary
		err = db.CreateBinary("some_binary", []byte("data"))
		if err != nil {
			t.Fatal(err)
		}

		// Try to create duplicate
		err = db.CreateBinary("some_binary", []byte("data"))
		if !errors.Is(err, fsentry_error.ErrorExist) {
			t.Fatal("Entry already exist")
		}

		err = db.Drop()
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("get", func(t *testing.T) {
		t.Parallel()

		db := NewFSEntry(path.Join("test", "test_binary_get"))
		err := db.Init()
		if err != nil {
			t.Fatal(err)
		}

		// Try to get bad name
		_, err = db.GetBinary("")
		if !errors.Is(err, fsentry_error.ErrorBadName) {
			t.Fatal("Bad name")
		}

		// Try to get not exist binary
		_, err = db.GetBinary("not_exist_binary")
		if !errors.Is(err, fsentry_error.ErrorNotExist) {
			t.Fatal("Binary not exist")
		}

		// Try to get from bad path
		_, err = db.GetBinary("not_exist_binary", "bad_path")
		if !errors.Is(err, fsentry_error.ErrorBadPath) {
			t.Fatal("Bad path")
		}

		err = db.CreateBinary("some_entry", []byte("data"))
		if err != nil {
			t.Fatal(err)
		}

		data, err := db.GetBinary("some_entry")
		if err != nil {
			t.Fatal(err)
		}

		if string(data) != "data" {
			t.Fatal("Bad data")
		}

		err = db.Drop()
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("move", func(t *testing.T) {
		t.Parallel()

		db := NewFSEntry(path.Join("test", "test_binary_move"))
		err := db.Init()
		if err != nil {
			t.Fatal(err)
		}

		// Try to move not exist binary
		err = db.MoveBinary("not_exist", "new_not_exist")
		if !errors.Is(err, fsentry_error.ErrorNotExist) {
			t.Fatal("Binary not exist")
		}

		// Try to move bad path
		err = db.MoveBinary("not_exist", "new_not_exist", "bad_path")
		if !errors.Is(err, fsentry_error.ErrorBadPath) {
			t.Fatal("Bad path")
		}

		err = db.CreateBinary("first_binary", nil)
		if err != nil {
			t.Fatal(err)
		}

		err = db.CreateBinary("second_binary", nil)
		if err != nil {
			t.Fatal(err)
		}

		// Try to move bad name
		err = db.MoveBinary("", "new_not_exist")
		if !errors.Is(err, fsentry_error.ErrorBadName) {
			t.Fatal("Bad name")
		}

		// Try to move bad name
		err = db.MoveBinary("first_binary", "")
		if !errors.Is(err, fsentry_error.ErrorBadName) {
			t.Fatal("Bad name")
		}

		// Try to move into exist binary
		err = db.MoveBinary("first_binary", "second_binary")
		if !errors.Is(err, fsentry_error.ErrorExist) {
			t.Fatal("Binary already exist")
		}

		err = db.MoveBinary("first_binary", "new_first_binary")
		if err != nil {
			t.Fatal(err)
		}

		err = db.Drop()
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("update", func(t *testing.T) {
		t.Parallel()

		db := NewFSEntry(path.Join("test", "test_binary_update"))
		err := db.Init()
		if err != nil {
			t.Fatal(err)
		}

		// Try to update bad name
		err = db.UpdateBinary("", []byte("data"))
		if !errors.Is(err, fsentry_error.ErrorBadName) {
			t.Fatal("Bad name")
		}

		// Try to update not exist binary
		err = db.UpdateBinary("not_exist_binary", []byte("data"))
		if !errors.Is(err, fsentry_error.ErrorNotExist) {
			t.Fatal("Binary not exist")
		}

		// Try to update bad path binary
		err = db.UpdateBinary("not_exist_binary", []byte("data"), "bad_path")
		if !errors.Is(err, fsentry_error.ErrorBadPath) {
			t.Fatal("Bad path")
		}

		err = db.CreateBinary("some_binary", []byte("data"))
		if err != nil {
			t.Fatal(err)
		}

		data, err := db.GetBinary("some_binary")
		if err != nil {
			t.Fatal(err)
		}
		if string(data) != "data" {
			t.Fatal("Bad data")
		}

		err = db.UpdateBinary("some_binary", []byte("other"))
		if err != nil {
			t.Fatal(err)
		}

		data, err = db.GetBinary("some_binary")
		if err != nil {
			t.Fatal(err)
		}
		if string(data) != "other" {
			t.Fatal("Bad data")
		}

		err = db.Drop()
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("remove", func(t *testing.T) {
		t.Parallel()

		db := NewFSEntry(path.Join("test", "test_binary_remove"))
		err := db.Init()
		if err != nil {
			t.Fatal(err)
		}

		// Try to remove bad name
		err = db.RemoveBinary("")
		if !errors.Is(err, fsentry_error.ErrorBadName) {
			t.Fatal("Bad name")
		}

		// Try to remove not exist binary
		err = db.RemoveBinary("not_exist_binary")
		if !errors.Is(err, fsentry_error.ErrorNotExist) {
			t.Fatal("Binary not exist")
		}

		err = db.CreateBinary("some_binary", []byte("data"))
		if err != nil {
			t.Fatal(err)
		}

		err = db.RemoveBinary("some_binary")
		if err != nil {
			t.Fatal(err)
		}

		// Try to remove twice
		err = db.RemoveEntry("some_binary")
		if !errors.Is(err, fsentry_error.ErrorNotExist) {
			t.Fatal("Binary not exist")
		}

		err = db.Drop()
		if err != nil {
			t.Fatal(err)
		}
	})
}
