package layout_test

import (
	"bytes"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"path/filepath"
	"testing"

	"github.com/HardDie/DeckBuilder/internal/render/generate/layout"
)

func TestPagesNamesAndCell(t *testing.T) {
	face := pngBytes(t, 8, 12, color.RGBA{R: 200, A: 255})
	back := pngBytes(t, 8, 12, color.RGBA{B: 80, A: 255})
	dir := t.TempDir()
	pages, backs, err := layout.Pages(dir, []layout.Deck{
		{ID: "bandits", Back: back, Faces: [][]byte{face}},
		{ID: "crew", Back: back, Faces: [][]byte{face, face}},
	}, 1, true)
	if err != nil {
		t.Fatal(err)
	}
	if len(backs) != 2 || len(pages) != 2 {
		t.Fatalf("backs %d pages %d", len(backs), len(pages))
	}
	if filepath.Base(pages[0].SheetPath) != "1_bandits_1_1_2x2.jpg" {
		t.Fatalf("first %s", pages[0].SheetPath)
	}
	if filepath.Base(pages[1].SheetPath) != "2_crew_1_2_2x2.jpg" {
		t.Fatalf("second %s", pages[1].SheetPath)
	}
	if pages[0].CellW != 8 || pages[0].CellH != 12 || pages[0].Cols != 2 || pages[0].Rows != 2 {
		t.Fatalf("cell %dx%d grid %dx%d", pages[0].CellW, pages[0].CellH, pages[0].Cols, pages[0].Rows)
	}
	if !pages[0].Shadow || !bytes.Equal(backs[0].Body, back) {
		t.Fatal("back bytes or shadow")
	}
}

func TestPagesSecondPageKeepsCell(t *testing.T) {
	face := pngBytes(t, 4, 4, color.RGBA{A: 255})
	faces := make([][]byte, 70)
	for i := range faces {
		faces[i] = face
	}
	pages, _, err := layout.Pages(t.TempDir(), []layout.Deck{{ID: "pile", Back: face, Faces: faces}}, 1, false)
	if err != nil {
		t.Fatal(err)
	}
	if len(pages) != 2 {
		t.Fatalf("pages %d", len(pages))
	}
	if filepath.Base(pages[0].SheetPath) != "1_pile_1_69_10x7.jpg" || pages[0].Cols != 10 || pages[0].Rows != 7 {
		t.Fatalf("first %+v", pages[0].SheetPath)
	}
	if filepath.Base(pages[1].SheetPath) != "2_pile_2_1_2x2.jpg" || pages[1].Cols != 2 || pages[1].Rows != 2 {
		t.Fatalf("second %+v", pages[1].SheetPath)
	}
	if pages[1].CellW != 4 || pages[1].CellH != 4 {
		t.Fatalf("cell %dx%d", pages[1].CellW, pages[1].CellH)
	}
}

func TestCellScaleAndWideHeader(t *testing.T) {
	face := pngBytes(t, 8, 12, color.RGBA{G: 10, A: 255})
	w, h := cell(t, face, 2)
	if w != 4 || h != 6 {
		t.Fatalf("scaled %dx%d", w, h)
	}
	w, h = cell(t, pngHeader(2000, 100), 1)
	if w != 1000 || h != 50 {
		t.Fatalf("wide %dx%d", w, h)
	}
}

func TestJPEGHeader(t *testing.T) {
	img := image.NewRGBA(image.Rect(0, 0, 3, 5))
	var buf bytes.Buffer
	if err := jpeg.Encode(&buf, img, &jpeg.Options{Quality: 80}); err != nil {
		t.Fatal(err)
	}
	w, h := cell(t, buf.Bytes(), 1)
	if w != 3 || h != 5 {
		t.Fatalf("jpeg %dx%d", w, h)
	}
}

func cell(t *testing.T, raw []byte, scale int) (int, int) {
	t.Helper()
	pages, _, err := layout.Pages(t.TempDir(), []layout.Deck{{
		ID: "d", Back: raw, Faces: [][]byte{raw},
	}}, scale, false)
	if err != nil {
		t.Fatal(err)
	}
	return pages[0].CellW, pages[0].CellH
}

func pngHeader(w, h int) []byte {
	b := make([]byte, 24)
	copy(b, []byte{0x89, 'P', 'N', 'G', '\r', '\n', 0x1a, '\n'})
	copy(b[12:16], []byte("IHDR"))
	b[11] = 13
	put := func(off, v int) {
		b[off] = byte(v >> 24)
		b[off+1] = byte(v >> 16)
		b[off+2] = byte(v >> 8)
		b[off+3] = byte(v)
	}
	put(16, w)
	put(20, h)
	return b
}

func pngBytes(t *testing.T, w, h int, c color.RGBA) []byte {
	t.Helper()
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			img.SetRGBA(x, y, c)
		}
	}
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		t.Fatal(err)
	}
	return buf.Bytes()
}
