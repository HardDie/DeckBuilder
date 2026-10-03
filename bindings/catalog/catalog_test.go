package catalog

import (
	"errors"
	"testing"

	er "github.com/HardDie/DeckBuilder/internal/errors"
)

func TestImageWarning(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want string
	}{
		{name: "no_error", want: ""},
		{
			name: "catalog_error_without_http_prefix",
			err:  er.NetworkTimeout.AddMessage("download timed out after 120 s"),
			want: "Image was not saved: download timed out after 120 s",
		},
		{
			name: "plain_error",
			err:  errors.New("disk full"),
			want: "Image was not saved: disk full",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ImageWarning(tt.err); got != tt.want {
				t.Fatalf("got %q, want %q", got, tt.want)
			}
		})
	}
}
