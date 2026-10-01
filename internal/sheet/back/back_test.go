package back

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"testing"
)

func TestWriteKeepsOriginalBytes(t *testing.T) {
	t.Parallel()
	raw := pngBytes(t, 3, 3, color.RGBA{G: 255, A: 255})
	dir := t.TempDir()
	abs, err := Write(dir, "deck", raw)
	if err != nil {
		t.Fatal(err)
	}
	got, err := os.ReadFile(abs)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(got, raw) {
		t.Fatal("back file bytes differ from the input")
	}
	if filepath.Base(abs) == "backside_deck_.png" {
		t.Fatal("missing hash in the file name")
	}
}

func TestShadeShadowChangesPixels(t *testing.T) {
	t.Parallel()
	img, err := png.Decode(bytes.NewReader(pngBytes(t, 2, 2, color.RGBA{R: 200, G: 200, B: 200, A: 255})))
	if err != nil {
		t.Fatal(err)
	}
	plain := Shade(img, false)
	dark := Shade(img, true)
	if plain.NRGBAAt(0, 0) == dark.NRGBAAt(0, 0) {
		t.Fatal("shadow did not change the pixel")
	}
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
