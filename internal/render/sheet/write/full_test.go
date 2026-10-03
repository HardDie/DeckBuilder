//go:build integration

package write

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/HardDie/DeckBuilder/internal/render/sheet/draw/libjpeg"
)

// A full 10×7 page from the 1312×962 inputs matches the sequential paint path byte for byte.
func TestIntegrationFullPageMatchesSequential(t *testing.T) {
	dir := filepath.Join("..", "bench", "testdata", "input")
	read := func(name string) []byte {
		body, err := os.ReadFile(filepath.Join(dir, name))
		if err != nil {
			t.Fatal(err)
		}
		return body
	}
	names := []string{"face_a.png", "face_b.png", "face_wide.png"}
	faces := make([][]byte, 69)
	for i := range faces {
		faces[i] = read(names[i%len(names)])
	}
	rawBack := read("back.png")
	const cellW, cellH = 1000, 733

	path := filepath.Join(t.TempDir(), "sheet.jpg")
	if err := Draw(faces, rawBack, cellW, cellH, true, path); err != nil {
		t.Fatal(err)
	}
	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	want, err := libjpeg.Encode(painted(t, faces, rawBack, cellW, cellH, true))
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(got, want) {
		t.Fatalf("parallel %d bytes, sequential %d bytes", len(got), len(want))
	}
}
