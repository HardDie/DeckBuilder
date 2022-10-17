package utils

import (
	"errors"
	"io"
	"net/http"

	er "github.com/HardDie/DeckBuilder/internal/errors"
)

func GetFileFromMultipart(name string, r *http.Request) ([]byte, error) {
	f, _, err := r.FormFile(name)
	if err != nil {
		if errors.Is(err, http.ErrMissingFile) {
			return nil, nil
		} else {
			er.IfErrorLog(err)
			err = er.InternalError.AddMessage(err.Error())
			return nil, err
		}
	}

	data, err := io.ReadAll(f)
	if err != nil {
		er.IfErrorLog(err)
		err = er.InternalError.AddMessage(err.Error())
		return nil, err
	}
	return data, nil
}
