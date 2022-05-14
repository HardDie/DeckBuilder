package fs

import (
	"io"

	"tts_deck_build/internal/errors"
)

func BinFromReader(r io.Reader) ([]byte, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		errors.IfErrorLog(err)
		return nil, errors.InternalError.AddMessage(err.Error())
	}
	return data, nil
}
func BinToWriter(w io.Writer, data []byte) error {
	// Write data to file
	_, err := w.Write(data)
	if err != nil {
		errors.IfErrorLog(err)
		return errors.InternalError.AddMessage(err.Error())
	}
	return nil
}
