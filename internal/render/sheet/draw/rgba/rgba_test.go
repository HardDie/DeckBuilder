package rgba

import (
	"bytes"
	"image"
	"image/color"
	"testing"

	"github.com/HardDie/DeckBuilder/internal/render/sheet/draw/seq"
)

func TestMatchesSeqJPEG(t *testing.T) {
	faces, back := deck()
	got, err := JPEG(faces, back)
	if err != nil {
		t.Fatal(err)
	}
	want, err := seq.JPEG(faces, back)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(got, want) {
		t.Fatalf("jpeg %d bytes, seq %d bytes", len(got), len(want))
	}
}

func deck() ([]image.Image, image.Image) {
	faces := make([]image.Image, 69)
	for i := range faces {
		faces[i] = pattern(4, 6, uint8(i))
	}
	faces[68] = pattern(9, 11, 200)
	return faces, pattern(5, 8, 10)
}

func pattern(w, h int, blue uint8) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			img.SetRGBA(x, y, color.RGBA{R: uint8(x * 20), G: uint8(y * 15), B: blue, A: 255})
		}
	}
	return img
}
