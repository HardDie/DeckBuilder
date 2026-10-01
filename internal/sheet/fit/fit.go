// Cell size for one face on a TTS sheet.
package fit

import (
	"image"
	"math"

	"github.com/disintegration/imaging"
)

// Size is the locked cell, plus the scale used while choosing it.
type Size struct {
	Inner  float64
	Width  int
	Height int
	Set    bool
}

// Apply repeats the page_drawer cell rule.
// A later face keeps Width and Height once Set is true.
// It still updates Inner.
func Apply(cardW, cardH, scale int, cur Size) (int, Size) {
	inner := cur.Inner
	if cardW*10 > 10_000 {
		inner = 10_000 / 10 / float64(cardW)
		if int(math.Trunc(float64(cardW)*inner)) > 10_000 {
			inner += 0.01
		}
	}
	if math.Trunc(float64(cardH)*inner)*7 > 10_000 {
		inner = 10_000 / 7 / float64(cardH)
		if int(math.Trunc(float64(cardH)*inner)) > 10_000 {
			inner += 0.01
		}
	}
	if scale == 0 {
		scale = 1
	}
	if inner == 0 {
		inner = 1
	}
	out := cur
	out.Inner = inner
	if !cur.Set {
		out.Width = int(math.Trunc(float64(cardW)*inner)) / scale
		out.Height = int(math.Trunc(float64(cardH)*inner)) / scale
		out.Set = true
	}
	return scale, out
}

// Resize uses Lanczos when Bounds().Max differs from the cell.
func Resize(img image.Image, width, height int) image.Image {
	if width != img.Bounds().Max.X || height != img.Bounds().Max.Y {
		return imaging.Resize(img, width, height, imaging.Lanczos)
	}
	return img
}
