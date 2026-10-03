package system

import (
	"testing"

	"github.com/HardDie/DeckBuilder/internal/dto"
	"github.com/HardDie/DeckBuilder/internal/network"
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
