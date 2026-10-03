package fit

import (
	"image"
	"image/color"
	"math"
	"testing"

	"github.com/disintegration/imaging"
)

func TestResize(t *testing.T) {
	t.Parallel()
	src := image.NewNRGBA(image.Rect(0, 0, 20, 30))
	if got := Resize(src, 20, 30); got != image.Image(src) {
		t.Fatal("same size should return the input unchanged")
	}
	if got := Resize(src, 10, 15); got.Bounds() != image.Rect(0, 0, 10, 15) {
		t.Fatalf("bounds %v, want 10x15", got.Bounds())
	}
}

// The resize stays visually identical to the imaging Lanczos it replaced (ADR 022).
func TestResizeMatchesLanczos(t *testing.T) {
	t.Parallel()
	src := image.NewNRGBA(image.Rect(0, 0, 262, 192))
	for y := 0; y < 192; y++ {
		for x := 0; x < 262; x++ {
			src.SetNRGBA(x, y, color.NRGBA{uint8(x), uint8(y), uint8(x ^ y), 255})
		}
	}
	a := Resize(src, 200, 147).(*image.NRGBA)
	b := imaging.Resize(src, 200, 147, imaging.Lanczos)
	var sum float64
	for i := range a.Pix {
		d := float64(a.Pix[i]) - float64(b.Pix[i])
		sum += d * d
	}
	if p := 10 * math.Log10(255*255/(sum/float64(len(a.Pix)))); p < 40 {
		t.Fatalf("PSNR %.1f dB vs Lanczos, want >= 40", p)
	}
}
