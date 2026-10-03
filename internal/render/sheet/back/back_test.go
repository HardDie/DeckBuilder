package back

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"testing"
)

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
