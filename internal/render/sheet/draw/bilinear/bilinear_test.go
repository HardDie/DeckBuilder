package bilinear

import (
	"bytes"
	"image"
	"image/color"
	"testing"

	"github.com/HardDie/DeckBuilder/internal/render/sheet/draw/seq"
)

func TestJPEGDiffersWhenResampled(t *testing.T) {
	faces := []image.Image{pattern(8, 10), pattern(3, 5)}
	back := pattern(8, 12)
	got, err := JPEG(faces, back)
	if err != nil {
		t.Fatal(err)
	}
	want, err := seq.JPEG(faces, back)
	if err != nil {
		t.Fatal(err)
	}
	if bytes.Equal(got, want) {
		t.Fatal("ApproxBiLinear matched Lanczos")
	}
}

func pattern(w, h int) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			img.SetRGBA(x, y, color.RGBA{R: uint8(x * 40), G: uint8(y * 25), B: uint8(x * y), A: 255})
		}
	}
	return img
}
