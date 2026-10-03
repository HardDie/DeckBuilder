// The deck back as drawn on the sheet.
package back

import (
	"image"

	"github.com/disintegration/imaging"
)

// Shade darkens by 30 when shadow is on.
// Brightness 0 still converts the image.
func Shade(img image.Image, shadow bool) *image.NRGBA {
	if shadow {
		return imaging.AdjustBrightness(img, -30)
	}
	return imaging.AdjustBrightness(img, 0)
}
