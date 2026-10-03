// Resize one face to the sheet cell.
package fit

import (
	"image"

	"github.com/disintegration/imaging"

	"github.com/HardDie/DeckBuilder/internal/render/sheet/stbresize"
)

// Resize scales img to the cell with stb_image_resize2 (Catmull-Rom).
// An image already at the cell size is returned unchanged.
func Resize(img image.Image, width, height int) image.Image {
	if width == img.Bounds().Dx() && height == img.Bounds().Dy() {
		return img
	}
	return stbresize.Resize(img, width, height)
}

// ResizeLanczos is the previous resize: pure Go Lanczos, about 5× slower.
//
// Deprecated: Use Resize. Kept as a fallback until stb_image_resize2
// is verified in release builds on every target (review item S15, ADR 022).
func ResizeLanczos(img image.Image, width, height int) image.Image {
	if width == img.Bounds().Dx() && height == img.Bounds().Dy() {
		return img
	}
	return imaging.Resize(img, width, height, imaging.Lanczos)
}
