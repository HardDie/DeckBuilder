package fit

import (
	"image"
	"image/color"
	"math"
	"testing"
)

func TestResize(t *testing.T) {
	t.Parallel()
	for name, resize := range map[string]func(image.Image, int, int) image.Image{
		"stb":     Resize,
		"lanczos": ResizeLanczos,
	} {
		src := image.NewNRGBA(image.Rect(0, 0, 20, 30))
		if got := resize(src, 20, 30); got != image.Image(src) {
			t.Fatalf("%s: same size should return the input unchanged", name)
		}
		if got := resize(src, 10, 15); got.Bounds() != image.Rect(0, 0, 10, 15) {
			t.Fatalf("%s: bounds %v, want 10x15", name, got.Bounds())
		}
	}
}

// The live resize stays visually identical to the Lanczos it replaced.
func TestResizeMatchesLanczos(t *testing.T) {
	t.Parallel()
	src := image.NewNRGBA(image.Rect(0, 0, 262, 192))
	for y := 0; y < 192; y++ {
		for x := 0; x < 262; x++ {
			src.SetNRGBA(x, y, color.NRGBA{uint8(x), uint8(y), uint8(x ^ y), 255})
		}
	}
	a := Resize(src, 200, 147).(*image.NRGBA)
	b := ResizeLanczos(src, 200, 147).(*image.NRGBA)
	var sum float64
	for i := range a.Pix {
		d := float64(a.Pix[i]) - float64(b.Pix[i])
		sum += d * d
	}
	if p := 10 * math.Log10(255*255/(sum/float64(len(a.Pix)))); p < 40 {
		t.Fatalf("PSNR %.1f dB vs Lanczos, want >= 40", p)
	}
}
