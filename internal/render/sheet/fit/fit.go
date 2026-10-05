// Resize one face to the sheet cell.
package fit

import (
	"image"

	"github.com/HardDie/DeckBuilder/third_party/stbresize"
)

// Resize scales img to the cell with stb_image_resize2 (Catmull-Rom).
// An image already at the cell size is returned unchanged.
func Resize(img image.Image, width, height int) image.Image {
	if width == img.Bounds().Dx() && height == img.Bounds().Dy() {
		return img
	}
	return stbresize.Resize(img, width, height)
}
