package write

import (
	"bytes"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"testing"

	"github.com/HardDie/DeckBuilder/internal/images"
	"github.com/HardDie/DeckBuilder/internal/render/sheet/back"
	"github.com/HardDie/DeckBuilder/internal/render/sheet/draw/libjpeg"
	"github.com/HardDie/DeckBuilder/internal/render/sheet/fit"
	"github.com/HardDie/DeckBuilder/internal/render/sheet/paint"
)

func TestDrawMatchesLibjpeg(t *testing.T) {
	face := pngBytes(t, 8, 12, color.RGBA{R: 220, G: 20, B: 40, A: 255})
	rawBack := pngBytes(t, 8, 12, color.RGBA{G: 180, B: 40, A: 255})
	path := filepath.Join(t.TempDir(), "sheet.jpg")
	if err := Draw([][]byte{face}, rawBack, 8, 12, false, path); err != nil {
		t.Fatal(err)
	}
	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	sheet := painted(t, [][]byte{face}, rawBack, 8, 12, false)
	want, err := libjpeg.Encode(sheet)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(got, want) {
		t.Fatalf("file %d bytes, libjpeg %d bytes", len(got), len(want))
	}
	var std bytes.Buffer
	if err := jpeg.Encode(&std, sheet, &jpeg.Options{Quality: 80}); err != nil {
		t.Fatal(err)
	}
	if bytes.Equal(got, std.Bytes()) {
		t.Fatal("sheet jpeg matches image/jpeg")
	}
}

func TestDrawShadowChangesPixels(t *testing.T) {
	face := pngBytes(t, 8, 12, color.RGBA{R: 200, A: 255})
	rawBack := pngBytes(t, 8, 12, color.RGBA{B: 200, A: 255})
	plain := filepath.Join(t.TempDir(), "a.jpg")
	shaded := filepath.Join(t.TempDir(), "b.jpg")
	if err := Draw([][]byte{face}, rawBack, 8, 12, false, plain); err != nil {
		t.Fatal(err)
	}
	if err := Draw([][]byte{face}, rawBack, 8, 12, true, shaded); err != nil {
		t.Fatal(err)
	}
	a, err := os.ReadFile(plain)
	if err != nil {
		t.Fatal(err)
	}
	b, err := os.ReadFile(shaded)
	if err != nil {
		t.Fatal(err)
	}
	if bytes.Equal(a, b) {
		t.Fatal("shadow did not change the sheet")
	}
}

func TestDrawUsesTheGivenCell(t *testing.T) {
	face := pngBytes(t, 8, 12, color.RGBA{G: 180, A: 255})
	rawBack := pngBytes(t, 8, 12, color.RGBA{R: 180, A: 255})
	path := filepath.Join(t.TempDir(), "sheet.jpg")
	if err := Draw([][]byte{face}, rawBack, 4, 6, false, path); err != nil {
		t.Fatal(err)
	}
	f, err := os.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	img, err := jpeg.Decode(f)
	if err != nil {
		t.Fatal(err)
	}
	if img.Bounds().Dx() != 8 || img.Bounds().Dy() != 12 {
		t.Fatalf("sheet %v", img.Bounds())
	}
}

func painted(t *testing.T, faces [][]byte, rawBack []byte, cellW, cellH int, shadow bool) image.Image {
	t.Helper()
	drawn := make([]image.Image, len(faces))
	for i, raw := range faces {
		img, err := images.ImageFromBinary(raw)
		if err != nil {
			t.Fatal(err)
		}
		drawn[i] = fit.Resize(img, cellW, cellH)
	}
	decoded, err := images.ImageFromBinary(rawBack)
	if err != nil {
		t.Fatal(err)
	}
	sheet, _, _ := paint.Canvas(cellW, cellH, drawn, fit.Resize(back.Shade(decoded, shadow), cellW, cellH))
	return sheet
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
