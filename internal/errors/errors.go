package errors

import (
	"fmt"
	"net/http"

	"github.com/HardDie/DeckBuilder/internal/logger"
)

var (
	// system
	InternalError = NewError("internal error")

	// network errors
	NetworkBadURL      = NewError("bad url", http.StatusBadRequest)
	NetworkBadRequest  = NewError("bad http request", http.StatusBadRequest)
	NetworkBadResponse = NewError("bad http response", http.StatusBadRequest)

	BadName = NewError("bad name", http.StatusBadRequest)
	BadId   = NewError("bad id", http.StatusBadRequest)

	// game
	GameExist          = NewError("game exist", http.StatusBadRequest)
	GameNotExists      = NewError("game not exists", http.StatusBadRequest)
	GameInfoNotExists  = NewError("game info not exists")
	GameImageExist     = NewError("game image already exists", http.StatusBadRequest)
	GameImageNotExists = NewError("game image not exists", http.StatusBadRequest)

	// collection
	CollectionExist          = NewError("collection exist", http.StatusBadRequest)
	CollectionNotExists      = NewError("collection not exists", http.StatusBadRequest)
	CollectionInfoNotExists  = NewError("collection info not exists")
	CollectionImageExist     = NewError("collection image already exists", http.StatusBadRequest)
	CollectionImageNotExists = NewError("collection image not exists", http.StatusBadRequest)

	// deck
	DeckExist          = NewError("deck exist", http.StatusBadRequest)
	DeckNotExists      = NewError("deck not exists", http.StatusBadRequest)
	DeckImageExist     = NewError("deck image already exists", http.StatusBadRequest)
	DeckImageNotExists = NewError("deck image not exists", http.StatusBadRequest)

	// card
	CardExists         = NewError("card exists", http.StatusBadRequest)
	CardNotExists      = NewError("card not exists", http.StatusBadRequest)
	CardImageExist     = NewError("card image already exists", http.StatusBadRequest)
	CardImageNotExists = NewError("card image not exists", http.StatusBadRequest)

	// settings
	SettingsNotExists = NewError("settings file not exists", http.StatusBadRequest)

	// image
	UnknownImageType = NewError("unknown image type").HTTP(http.StatusBadRequest)

	// zip
	BadArchive = NewError("bad zip archive").HTTP(http.StatusBadRequest)
)

type Err struct {
	Message string `json:"message"`
	Code    int
	Err     error
}

func NewError(message string, code ...int) *Err {
	err := &Err{
		Message: message,
		Code:    http.StatusInternalServerError,
	}
	if len(code) > 0 {
		err.Code = code[0]
	}
	return err
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
		_ = logger.Error.Output(2, err.Error())
	}
}
