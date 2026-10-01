package paint

import (
	"image"
	"image/color"
	"testing"
)

func TestCanvasBottomRightIsBack(t *testing.T) {
	t.Parallel()
	face := block(2, 2, color.RGBA{R: 255, A: 255})
	back := block(2, 2, color.RGBA{B: 255, A: 255})
	canvas, cols, rows := Canvas(2, 2, []image.Image{face}, back)
	if cols != 2 || rows != 2 {
		t.Fatalf("grid %dx%d", cols, rows)
	}
	if canvas.RGBAAt(0, 0).R == 0 {
		t.Fatal("face missing from the first cell")
	}
	lastX := canvas.Bounds().Dx() - 1
	lastY := canvas.Bounds().Dy() - 1
	if canvas.RGBAAt(lastX, lastY).B == 0 {
		t.Fatal("back missing from the bottom-right cell")
	}
}

func block(w, h int, c color.RGBA) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			img.SetRGBA(x, y, c)
		}
	}
	return img
}
