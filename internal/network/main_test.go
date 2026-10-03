package network

import (
	stderrors "errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/HardDie/DeckBuilder/internal/apperr"
)

func TestResponseError(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want int
	}{
		{name: "not_found", err: apperr.ErrGameNotFound, want: http.StatusNotFound},
		{name: "wrapped_not_found", err: fmt.Errorf("load: %w", apperr.ErrGameNotFound), want: http.StatusNotFound},
		{name: "already_exists", err: apperr.ErrGameExists, want: http.StatusConflict},
		{name: "busy", err: apperr.ErrRenderInProgress, want: http.StatusConflict},
		{name: "invalid", err: apperr.ErrBadName, want: http.StatusBadRequest},
		{name: "tts_nothing_to_serve", err: apperr.ErrNothingForTTS, want: http.StatusNotFound},
		{name: "plain_error", err: stderrors.New("boom"), want: http.StatusInternalServerError},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			ResponseError(rec, tt.err)
			if rec.Code != tt.want {
				t.Fatalf("code %d, want %d", rec.Code, tt.want)
			}
		})
	}
}
