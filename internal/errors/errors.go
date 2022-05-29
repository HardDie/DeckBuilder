package errors

import (
	"fmt"
	"log"
	"net/http"
)

var (
	// system
	InternalError = NewError("internal error").HTTP(http.StatusInternalServerError)

	// network errors
	NetworkBadURL      = NewError("bad url")
	NetworkBadRequest  = NewError("bad http request")
	NetworkBadResponse = NewError("bad http response")

	BadName = NewError("bad name").HTTP(http.StatusBadRequest)

	// game
	GameExist          = NewError("game exist").HTTP(http.StatusBadRequest)
	GameNotExists      = NewError("game not exists").HTTP(http.StatusNoContent)
	GameInfoNotExists  = NewError("game info not exists")
	GameImageNotExists = NewError("game image not exists").HTTP(http.StatusNoContent)

	// collection
	CollectionExist          = NewError("collection exist").HTTP(http.StatusBadRequest)
	CollectionNotExists      = NewError("collection not exists").HTTP(http.StatusNoContent)
	CollectionInfoNotExists  = NewError("collection info not exists")
	CollectionImageNotExists = NewError("collection image not exists").HTTP(http.StatusNoContent)

	// deck
	DeckExist          = NewError("deck exist").HTTP(http.StatusBadRequest)
	DeckNotExists      = NewError("deck not exists").HTTP(http.StatusNoContent)
	DeckImageNotExists = NewError("deck image not exists").HTTP(http.StatusNoContent)

	// image
	UnknownImageType = NewError("unknown image type").HTTP(http.StatusBadRequest)
)

type Err struct {
	Message string `json:"message"`
	Code    int
	Err     error
}

func NewError(message string) *Err {
	return &Err{
		Message: message,
		Code:    http.StatusBadRequest,
	}
}

func (e Err) Error() string {
	return fmt.Sprintf("HTTP[%d] %s", e.GetCode(), e.GetMessage())
}
func (e Err) Unwrap() error {
	return e.Err
}

func (e *Err) HTTP(code int) *Err {
	return &Err{
		Message: e.Message,
		Code:    code,
		Err:     e,
	}
}
func (e *Err) AddMessage(message string) *Err {
	return &Err{
		Message: message,
		Code:    e.Code,
		Err:     e,
	}
}

func (e *Err) GetCode() int       { return e.Code }
func (e *Err) GetMessage() string { return e.Message }

func IfErrorLog(err error) {
	if err != nil {
		_ = log.Output(2, err.Error())
	}
}
