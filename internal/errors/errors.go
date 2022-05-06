package errors

import (
	"fmt"
	"log"
	"net/http"
)

var (
	InternalError  = NewError("internal error").HTTP(http.StatusInternalServerError)
	BarURL         = NewError("bad url")
	BadHTTPRequest = NewError("bad http request")

	BadName = NewError("bad name").HTTP(http.StatusBadRequest)

	GameExist          = NewError("game exist").HTTP(http.StatusBadRequest)
	GameNotExists      = NewError("game not exists").HTTP(http.StatusNoContent)
	GameInvalid        = NewError("game data invalid")
	GameInfoNotExists  = NewError("game info not exists").HTTP(http.StatusInternalServerError)
	GameImageNotExists = NewError("game image not exists").HTTP(http.StatusNoContent)

	CollectionExist          = NewError("collection exist")
	CollectionNotExists      = NewError("collection not exists").HTTP(http.StatusNoContent)
	CollectionInvalid        = NewError("collection data invalid")
	CollectionInfoNotExists  = NewError("collection info not exists").HTTP(http.StatusInternalServerError)
	CollectionImageNotExists = NewError("collection image not exists").HTTP(http.StatusNoContent)

	DeckExist     = NewError("deck exist")
	DeckNotExists = NewError("deck not exists").HTTP(http.StatusNoContent)

	UnknownImageType = NewError("unknown image type")
)

type Err struct {
	Message string `json:"message"`
	code    int
}

func NewError(message string) *Err {
	return &Err{
		Message: message,
		code:    http.StatusBadRequest,
	}
}

func (e Err) HTTP(code int) *Err {
	e.code = code
	return &e
}
func (e Err) AddMessage(message string) *Err {
	e.Message += ": " + message
	return &e
}
func (e Err) Error() string {
	return fmt.Sprintf("HTTP[%d] %s", e.GetCode(), e.GetMessage())
}

func (e *Err) GetCode() int       { return e.code }
func (e *Err) GetMessage() string { return e.Message }

func IfErrorLog(err error) {
	if err != nil {
		_ = log.Output(2, err.Error())
	}
}
