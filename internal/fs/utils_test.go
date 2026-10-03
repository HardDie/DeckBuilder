package fs

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"testing"
)

func TestCreateAndProcess(t *testing.T) {
	writeErr := errors.New("write failed")
	tests := []struct {
		name     string
		dir      string
		cb       func(w io.Writer, in string) error
		wantErr  error
		wantBody string
	}{
		{
			name:     "writes_the_file",
			cb:       writeString,
			wantBody: "body",
		},
		{
			name:    "callback_error_comes_back",
			cb:      func(io.Writer, string) error { return writeErr },
			wantErr: writeErr,
		},
		{
			name:    "missing_folder",
			dir:     "missing",
			cb:      writeString,
			wantErr: os.ErrNotExist,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), tt.dir, "out.bin")
			err := CreateAndProcess(path, "body", tt.cb)
			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("err %v, want %v", err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			got, err := os.ReadFile(path)
			if err != nil {
				t.Fatal(err)
			}
			if string(got) != tt.wantBody {
				t.Fatalf("body %q, want %q", got, tt.wantBody)
			}
		})
	}
}

// writeString writes in to w, for the table above.
func writeString(w io.Writer, in string) error {
	_, err := io.WriteString(w, in)
	return err
}
