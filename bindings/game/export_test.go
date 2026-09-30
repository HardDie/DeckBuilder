package game

import "testing"

func TestExportFilename(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		gameID string
		want   string
	}{
		{name: "Munchkin", gameID: "munchkin", want: "Munchkin.zip"},
		{name: "a/b:c", gameID: "abc", want: "a_b_c.zip"},
		{name: "  ", gameID: "munchkin", want: "munchkin.zip"},
		{name: "...", gameID: "", want: "game.zip"},
		{name: "", gameID: "", want: "game.zip"},
	}

	for _, tt := range tests {
		t.Run(tt.want+" "+tt.gameID, func(t *testing.T) {
			t.Parallel()
			if got := exportFilename(tt.name, tt.gameID); got != tt.want {
				t.Fatalf("exportFilename(%q, %q) = %q, want %q", tt.name, tt.gameID, got, tt.want)
			}
		})
	}
}
