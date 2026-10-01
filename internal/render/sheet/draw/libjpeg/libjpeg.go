// Quality-80 JPEG through libjpeg-turbo.
package libjpeg

import (
	"bytes"
	"image"

	turbojpeg "github.com/pixiv/go-libjpeg/jpeg"
)

// Encode writes img at quality 80.
// The image must be *image.RGBA, *image.YCbCr, or *image.Gray.
func Encode(img image.Image) ([]byte, error) {
	var buf bytes.Buffer
	err := turbojpeg.Encode(&buf, img, &turbojpeg.EncoderOptions{
		Quality:   80,
		DCTMethod: turbojpeg.DCTISlow,
	})
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
