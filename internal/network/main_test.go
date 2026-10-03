package network

import (
	stderrors "errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/HardDie/DeckBuilder/internal/errors"
)

func TestResponseError(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want int
	}{
		{name: "catalog_error", err: errors.GameNotExists, want: http.StatusBadRequest},
		{name: "wrapped_catalog_error", err: fmt.Errorf("load: %w", error(errors.GameNotExists)), want: http.StatusBadRequest},
		{name: "tts_nothing_to_serve", err: errors.TTSNothingToServe, want: http.StatusNotFound},
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
