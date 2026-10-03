package layout_test

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"hash/crc32"
	"image"
	"image/color"
	"image/color/palette"
	"image/gif"
	"image/jpeg"
	"image/png"
	"path/filepath"
	"regexp"
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
	if !sheetName("bandits", 1).MatchString(filepath.Base(pages[0].SheetPath)) {
		t.Fatalf("first %s", pages[0].SheetPath)
	}
	if !sheetName("crew", 1).MatchString(filepath.Base(pages[1].SheetPath)) {
		t.Fatalf("second %s", pages[1].SheetPath)
	}
	if !regexp.MustCompile(`^backside_bandits_[0-9a-f]{16}\.png$`).MatchString(filepath.Base(backs[0].Path)) {
		t.Fatalf("back %s", backs[0].Path)
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
	if !sheetName("pile", 1).MatchString(filepath.Base(pages[0].SheetPath)) || pages[0].Cols != 10 || pages[0].Rows != 7 {
		t.Fatalf("first %+v", pages[0].SheetPath)
	}
	if !sheetName("pile", 2).MatchString(filepath.Base(pages[1].SheetPath)) || pages[1].Cols != 2 || pages[1].Rows != 2 {
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

// pngHeader is a PNG with only its signature and a complete IHDR chunk:
// enough for image.DecodeConfig, with no pixel data.
func pngHeader(w, h int) []byte {
	var ihdr bytes.Buffer
	ihdr.WriteString("IHDR")
	_ = binary.Write(&ihdr, binary.BigEndian, uint32(w))
	_ = binary.Write(&ihdr, binary.BigEndian, uint32(h))
	ihdr.Write([]byte{8, 6, 0, 0, 0}) // 8-bit RGBA, no interlace

	var out bytes.Buffer
	out.WriteString("\x89PNG\r\n\x1a\n")
	_ = binary.Write(&out, binary.BigEndian, uint32(ihdr.Len()-4))
	out.Write(ihdr.Bytes())
	_ = binary.Write(&out, binary.BigEndian, crc32.ChecksumIEEE(ihdr.Bytes()))
	return out.Bytes()
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

func sheetName(deck string, page int) *regexp.Regexp {
	return regexp.MustCompile(fmt.Sprintf(`^%s_%d_[0-9a-f]{16}\.jpg$`, deck, page))
}

// The sheet name changes exactly when something drawn on the page changes.
func TestPagesNameFollowsContent(t *testing.T) {
	red := pngBytes(t, 8, 12, color.RGBA{R: 200, A: 255})
	blue := pngBytes(t, 8, 12, color.RGBA{B: 200, A: 255})
	back := pngBytes(t, 8, 12, color.RGBA{G: 80, A: 255})
	otherBack := pngBytes(t, 8, 12, color.RGBA{G: 90, A: 255})
	big := pngBytes(t, 16, 24, color.RGBA{R: 200, A: 255})

	name := func(faces [][]byte, back []byte, scale int, shadow bool) string {
		t.Helper()
		pages, _, err := layout.Pages(t.TempDir(), []layout.Deck{{ID: "crew", Back: back, Faces: faces}}, scale, shadow)
		if err != nil {
			t.Fatal(err)
		}
		return filepath.Base(pages[0].SheetPath)
	}
	base := name([][]byte{red, blue}, back, 1, true)

	if got := name([][]byte{red, blue}, back, 1, true); got != base {
		t.Fatalf("same input, different name: %s vs %s", got, base)
	}
	tests := []struct {
		name string
		got  string
	}{
		{"order", name([][]byte{blue, red}, back, 1, true)},
		{"face", name([][]byte{red, red}, back, 1, true)},
		{"face_count", name([][]byte{red, blue, red}, back, 1, true)},
		{"back", name([][]byte{red, blue}, otherBack, 1, true)},
		{"shadow", name([][]byte{red, blue}, back, 1, false)},
		{"cell", name([][]byte{big, blue}, back, 1, true)},
	}
	for _, tt := range tests {
		if tt.got == base {
			t.Errorf("%s change kept the name %s", tt.name, base)
		}
	}
	if name([][]byte{big, blue}, back, 2, true) == name([][]byte{big, blue}, back, 1, true) {
		t.Error("scale change kept the name")
	}
}

func TestPagesReadsEveryAcceptedFormat(t *testing.T) {
	var gifBuf bytes.Buffer
	if err := gif.Encode(&gifBuf, image.NewPaletted(image.Rect(0, 0, 8, 12), palette.Plan9), nil); err != nil {
		t.Fatal(err)
	}
	if w, h := cell(t, gifBuf.Bytes(), 1); w != 8 || h != 12 {
		t.Fatalf("gif %dx%d", w, h)
	}

	_, _, err := layout.Pages(t.TempDir(), []layout.Deck{{ID: "crew", Faces: [][]byte{[]byte("not an image")}}}, 1, false)
	var imageErr *layout.ImageError
	if !errors.As(err, &imageErr) || imageErr.DeckID != "crew" {
		t.Fatalf("err %v, want an ImageError for deck crew", err)
	}
}
