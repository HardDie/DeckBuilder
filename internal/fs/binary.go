package fs

import (
	"io"

	"github.com/HardDie/DeckBuilder/internal/errors"
)

func BinToWriter(w io.Writer, data []byte) error {
	// Write data to file
	_, err := w.Write(data)
	if err != nil {
		errors.IfErrorLog(err)
		return errors.InternalError.AddMessage(err.Error())
	}
	return nil
}
