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

func TestWriteAtomic(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "page.jpg")
	if err := os.WriteFile(path, []byte("old"), 0o644); err != nil {
		t.Fatal(err)
	}

	// A failed write keeps the old file and leaves no temporary file.
	err := WriteAtomic(path, func(tmp string) error {
		_ = os.WriteFile(tmp, []byte("half"), 0o644)
		return errors.New("draw failed")
	})
	if err == nil {
		t.Fatal("want the write error")
	}
	if got, _ := os.ReadFile(path); string(got) != "old" {
		t.Fatalf("old file changed to %q", got)
	}
	if FileExists(path + ".tmp") {
		t.Fatal("temporary file left behind")
	}

	// A good write replaces the file.
	if err = WriteAtomic(path, func(tmp string) error { return os.WriteFile(tmp, []byte("new"), 0o644) }); err != nil {
		t.Fatal(err)
	}
	if got, _ := os.ReadFile(path); string(got) != "new" {
		t.Fatalf("got %q, want new", got)
	}
	if !FileExists(path) || FileExists(dir) {
		t.Fatal("FileExists: a file is true, a folder is false")
	}
}
