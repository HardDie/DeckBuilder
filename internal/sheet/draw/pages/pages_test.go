package pages

import (
	"bytes"
	"image"
	"image/color"
	"testing"

	"github.com/HardDie/DeckBuilder/internal/sheet/draw/seq"
)

func TestEachPageMatchesSeq(t *testing.T) {
	first := Sheet{
		Faces: []image.Image{pattern(4, 6, 1), pattern(4, 6, 2)},
		Back:  pattern(4, 6, 9),
	}
	second := Sheet{
		Faces: []image.Image{pattern(5, 7, 3)},
		Back:  pattern(5, 7, 4),
	}
	got, err := JPEGs([]Sheet{first, second})
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 2 {
		t.Fatalf("got %d pages", len(got))
	}
	want0, err := seq.JPEG(first.Faces, first.Back)
	if err != nil {
		t.Fatal(err)
	}
	want1, err := seq.JPEG(second.Faces, second.Back)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(got[0], want0) || !bytes.Equal(got[1], want1) {
		t.Fatal("page jpeg bytes differ from seq")
	}
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
