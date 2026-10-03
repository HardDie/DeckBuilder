package network

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/HardDie/DeckBuilder/internal/apperr"
	"github.com/HardDie/DeckBuilder/internal/fs"
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

// ResponseError writes e as the JSON envelope. The status comes from the error's kind.
func ResponseError(w http.ResponseWriter, e error) {
	resp := JSONResponse{Error: apperr.Message(e)}

	switch {
	case errors.Is(e, apperr.ErrNotFound):
		_ = response(w, http.StatusNotFound, resp)
	case errors.Is(e, apperr.ErrAlreadyExists), errors.Is(e, apperr.ErrBusy):
		_ = response(w, http.StatusConflict, resp)
	case errors.Is(e, apperr.ErrInvalid):
		_ = response(w, http.StatusBadRequest, resp)
	default:
		slog.Warn("unexpected HTTP error", "err", e)
		_ = response(w, http.StatusInternalServerError, resp)
	}
}
