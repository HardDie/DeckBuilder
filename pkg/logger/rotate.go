package logger

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// rotate renames path to path.1 when it is larger than maxSize,
// shifting older files up to path.<keep> and deleting the oldest.
// The number goes before the extension: app.1.log.
// maxSize 0 or keep 0 never rotates.
func rotate(path string, maxSize int64, keep int) error {
	if maxSize <= 0 || keep <= 0 {
		return nil
	}
	info, err := os.Stat(path)
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	if err != nil {
		return err
	}
	if info.Size() <= maxSize {
		return nil
	}
	if err := os.Remove(numbered(path, keep)); err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}
	for i := keep - 1; i >= 1; i-- {
		if err := os.Rename(numbered(path, i), numbered(path, i+1)); err != nil && !errors.Is(err, os.ErrNotExist) {
			return err
		}
	}
	return os.Rename(path, numbered(path, 1))
}

func numbered(path string, i int) string {
	ext := filepath.Ext(path)
	return fmt.Sprintf("%s.%d%s", strings.TrimSuffix(path, ext), i, ext)
}
