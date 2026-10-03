package network

import (
	stderrors "errors"
	"net/http"

	"github.com/HardDie/DeckBuilder/internal/errors"
	"github.com/HardDie/DeckBuilder/internal/fs"
	"github.com/HardDie/DeckBuilder/internal/logger"
)

type Meta struct {
	Total      int `json:"total"`
	CardsTotal int `json:"cardsTotal,omitempty"`
	//Limit int `json:"limit"`
	//Page  int `json:"page"`
}
type JSONResponse struct {
	// Body
	Data interface{} `json:"data,omitempty"`
	// Meta
	Meta *Meta `json:"meta,omitempty"`
	// Error information
	Error interface{} `json:"error,omitempty"`
}

func response(w http.ResponseWriter, httpCode int, data interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(httpCode)
	return fs.JsonToWriter(w, data)
}
func ResponseError(w http.ResponseWriter, e error) {
	resp := JSONResponse{
		Error: e,
	}

	httpCode := http.StatusInternalServerError
	var val *errors.Err
	if stderrors.As(e, &val) {
		if val.GetCode() > 0 {
			httpCode = val.GetCode()
		}
	} else {
		logger.Warn.Println("unhandled error: " + e.Error())
	}

	_ = response(w, httpCode, resp)
}
