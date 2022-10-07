package fs

import (
	"encoding/json"
	"io"

	"github.com/HardDie/DeckBuilder/internal/errors"
)

func JsonFromReader[T any](r io.Reader) (data *T, err error) {
	err = json.NewDecoder(r).Decode(&data)
	if err != nil {
		errors.IfErrorLog(err)
		return nil, errors.InternalError.AddMessage(err.Error())
	}
	return data, nil
}
func JsonToWriter[T any](w io.Writer, data T) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "	")
	if err := enc.Encode(data); err != nil {
		errors.IfErrorLog(err)
		return errors.InternalError.AddMessage(err.Error())
	}
	return nil
}
