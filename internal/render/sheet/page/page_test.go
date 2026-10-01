package page

import (
	"bytes"
	"image"
	"image/png"
	"testing"

	entitiesSettings "github.com/HardDie/DeckBuilder/internal/entities/settings"
)

func TestEmptySave(t *testing.T) {
	t.Parallel()
	p := New("deck", t.TempDir(), 1, 1, settings(false))
	path, cols, rows, err := p.Save()
	if err != nil {
		t.Fatal(err)
	}
	if path != "" || cols != 0 || rows != 0 {
		t.Fatalf("empty save = %q %d %d", path, cols, rows)
	}
}

func TestFullPageRejectsAnotherFace(t *testing.T) {
	t.Parallel()
	p := New("deck", t.TempDir(), 1, 1, settings(false))
	face := solidPNG(t, 2, 2)
	if _, err := p.SetBacksideImageAndSave(face); err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 69; i++ {
		if err := p.AddImage(face); err != nil {
			t.Fatalf("face %d: %v", i, err)
		}
	}
	if !p.IsFull() {
		t.Fatal("69 faces is a full page")
	}
	err := p.AddImage(face)
	if err == nil || err.Error() != "page is full" {
		t.Fatalf("full page error = %v", err)
	}
}

func TestBadFace(t *testing.T) {
	t.Parallel()
	p := New("deck", t.TempDir(), 1, 1, settings(false))
	if err := p.AddImage([]byte("not an image")); err == nil {
		t.Fatal("expected decode error")
	}
}

func TestInheritKeepsCellAndDropsFaces(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	first := New("deck", dir, 1, 4, settings(false))
	back := solidPNG(t, 4, 6)
	if _, err := first.SetBacksideImageAndSave(back); err != nil {
		t.Fatal(err)
	}
	if err := first.AddImage(solidPNG(t, 4, 6)); err != nil {
		t.Fatal(err)
	}
	next := (&Page{}).Inherit(first)
	if !next.IsEmpty() || next.GetIndex() != 2 || next.Size() != 0 {
		t.Fatalf("inherit index=%d size=%d empty=%v", next.GetIndex(), next.Size(), next.IsEmpty())
	}
	if err := next.AddImage(solidPNG(t, 8, 8)); err != nil {
		t.Fatal(err)
	}
	if next.size.Width != first.size.Width || next.size.Height != first.size.Height {
		t.Fatalf("cell %dx%d, first page %dx%d", next.size.Width, next.size.Height, first.size.Width, first.size.Height)
	}
}

func settings(shadow bool) *entitiesSettings.Settings {
	cfg := entitiesSettings.Default()
	cfg.EnableBackShadow = shadow
	return &cfg
}

func solidPNG(t *testing.T, w, h int) []byte {
	t.Helper()
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		t.Fatal(err)
	}
	return buf.Bytes()
}
