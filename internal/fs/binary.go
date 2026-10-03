package fs

import (
	"fmt"
	"io"
)

func BinToWriter(w io.Writer, data []byte) error {
	// Write data to file
	_, err := w.Write(data)
	if err != nil {
		return fmt.Errorf("write file: %w", err)
	}
	return nil
}
