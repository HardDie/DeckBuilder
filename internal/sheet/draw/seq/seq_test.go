package seq

import (
	"bytes"
	"image"
	"image/color"
	"image/jpeg"
	"testing"
)

func TestJPEGGrid(t *testing.T) {
	face := pattern(4, 6)
	back := pattern(4, 6)
	got, err := JPEG([]image.Image{face}, back)
	if err != nil {
		t.Fatal(err)
	}
	img, err := jpeg.Decode(bytes.NewReader(got))
	if err != nil {
		t.Fatal(err)
	}
	if img.Bounds().Dx() != 8 || img.Bounds().Dy() != 12 {
		t.Fatalf("sheet is %dx%d, want 8x12", img.Bounds().Dx(), img.Bounds().Dy())
	}
}

func pattern(w, h int) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			img.SetRGBA(x, y, color.RGBA{R: uint8(x * 30), G: uint8(y * 20), B: 40, A: 255})
		}
	}
	return img
}
