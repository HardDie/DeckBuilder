package fs

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/HardDie/DeckBuilder/internal/logger"
)

const (
	DirPerm = 0755
)

func CreateFolder(path string) error {
	err := os.MkdirAll(path, DirPerm)
	if err != nil {
		return fmt.Errorf("create folder %q: %w", path, err)
	}
	return nil
}
func RemoveFolder(path string) error {
	err := os.RemoveAll(path)
	if err != nil {
		return fmt.Errorf("remove folder %q: %w", path, err)
	}
	return nil
}

// CreateAndProcess creates path and lets cb write into it.
// A Close error is returned too: some disks report a failed write only there.
func CreateAndProcess[T any](path string, in T, cb func(w io.Writer, in T) error) (err error) {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer func() {
		if closeErr := file.Close(); closeErr != nil && err == nil {
			err = closeErr
		}
	}()

	return cb(file, in)
}

// WriteAtomic lets write fill a temporary file next to path, then renames it to path.
// A crash or an error leaves no half-written file at path.
func WriteAtomic(path string, write func(tmp string) error) error {
	tmp := path + ".tmp"
	if err := write(tmp); err != nil {
		_ = os.Remove(tmp)
		return err
	}
	if err := os.Rename(tmp, path); err != nil {
		_ = os.Remove(tmp)
		return err
	}
	return nil
}

// FileExists reports whether path is an existing regular file.
func FileExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.Mode().IsRegular()
}

func PathToAbsolutePath(path string) string {
	res, err := filepath.Abs(path)
	if err != nil {
		logger.Error.Printf("Can't transform path %q to absolute path. %q", path, err.Error())
		return path
	}
	return res
}
