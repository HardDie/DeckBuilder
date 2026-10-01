// Quality-80 JPEG through jpegli.
package jpegli

import (
	"bytes"
	"image"

	jli "github.com/gen2brain/jpegli"
)

// Encode writes img at quality 80, 4:2:0, sequential.
func Encode(img image.Image) ([]byte, error) {
	var buf bytes.Buffer
	err := jli.Encode(&buf, img, &jli.EncodingOptions{
		Quality:           80,
		ChromaSubsampling: image.YCbCrSubsampleRatio420,
		DCTMethod:         jli.DCTISlow,
	})
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
