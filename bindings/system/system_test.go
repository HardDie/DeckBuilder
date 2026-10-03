package system

import (
	"errors"
	"testing"

	"github.com/HardDie/DeckBuilder/internal/apperr"
	"github.com/HardDie/DeckBuilder/internal/dto"
	"github.com/HardDie/DeckBuilder/internal/network"
	renderprogress "github.com/HardDie/DeckBuilder/internal/render/progress"
)

func TestDownloadStatusDTO(t *testing.T) {
	tests := []struct {
		name  string
		state network.DownloadState
		want  dto.DownloadStatus
	}{
		{name: "idle"},
		{
			name:  "half",
			state: network.DownloadState{Active: true, Done: 50, Total: 200},
			want:  dto.DownloadStatus{Active: true, Done: 50, Total: 200, Percent: 25},
		},
		{
			name:  "unknown_size",
			state: network.DownloadState{Active: true, Done: 50},
			want:  dto.DownloadStatus{Active: true, Done: 50},
		},
		{
			name:  "more_than_announced",
			state: network.DownloadState{Active: true, Done: 300, Total: 200},
			want:  dto.DownloadStatus{Active: true, Done: 300, Total: 200, Percent: 100},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := downloadStatusDTO(tt.state); got != tt.want {
				t.Fatalf("got %+v, want %+v", got, tt.want)
			}
		})
	}
}

func TestStatusCarriesTheRenderError(t *testing.T) {
	s := &System{}
	renderprogress.Reset()
	renderprogress.Begin()
	renderprogress.Fail(apperr.Withf(apperr.ErrUnsupportedImage, "an image in deck %q could not be read", "Crew"))

	got := s.Status().Data
	if got.Status != renderprogress.Error || got.Message != `An image in deck "Crew" could not be read` {
		t.Fatalf("status %+v", got)
	}
	if again := s.Status().Data; again.Status != renderprogress.Empty || again.Message != "" {
		t.Fatalf("after reading the error, status %+v", again)
	}

	renderprogress.Begin()
	renderprogress.Fail(errors.New("disk full"))
	if got := s.Status().Data; got.Message != apperr.Unexpected {
		t.Fatalf("unexpected error message %q", got.Message)
	}

	renderprogress.Begin()
	renderprogress.Finish()
	if got := s.Status().Data; got.Status != renderprogress.Done || got.Message != "" {
		t.Fatalf("done status %+v", got)
	}
}
