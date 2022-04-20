package errors

import (
	"log"
	"net/http"
)

var (
	InternalError = NewError("internal error").HTTP(http.StatusInternalServerError)
	DataInvalid   = NewError("game data invalid").HTTP(http.StatusNoContent)

	GameExist         = NewError("game exist")
	GameNotExists     = NewError("game not exists")
	GameInvalid       = NewError("game data invalid")
	GameInfoNotExists = NewError("game info not exists")

	CollectionExist         = NewError("collection exist")
	CollectionNotExists     = NewError("collection not exists")
	CollectionInvalid       = NewError("collection data invalid")
	CollectionInfoNotExists = NewError("collection info not exists")

	DeckNotExists = NewError("deck not exists")
)

type Error struct {
	Message string `json:"message"`
	code    int
}

func NewError(message string) *Error {
	return &Error{
		Message: message,
		code:    http.StatusBadRequest,
	}
}

func (e Error) HTTP(code int) *Error {
	e.code = code
	return &e
}
func (e Error) AddMessage(message string) *Error {
	e.Message += ": " + message
	return &e
}

func (e *Error) GetCode() int       { return e.code }
func (e *Error) GetMessage() string { return e.Message }

func IfErrorLog(err error) {
	if err != nil {
		log.Output(2, err.Error())
	}
}
