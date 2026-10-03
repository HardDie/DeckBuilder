package images

import (
	"bytes"
	stderrors "errors"
	"image"
	"image/draw"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"

	"github.com/HardDie/DeckBuilder/internal/apperr"
	"github.com/HardDie/DeckBuilder/internal/logger"
)

// maxImagePixels caps width × height of an accepted image (128 MP, e.g. 16384×8192).
// Decoding allocates about 4 bytes per pixel.
const maxImagePixels = 16384 * 8192

// ValidateImage checks the header size first, then fully decodes input.
// A tiny file can declare huge dimensions, so the header check comes before decode.
func ValidateImage(input []byte) (string, error) {
	cfg, _, err := image.DecodeConfig(bytes.NewReader(input))
	if stderrors.Is(err, image.ErrFormat) {
		return "", apperr.ErrUnsupportedImage
	}
	if err != nil {
		return "", damaged(err)
	}
	if int64(cfg.Width)*int64(cfg.Height) > maxImagePixels {
		return "", apperr.Withf(apperr.ErrImageTooLarge, "image is too large: %dx%d, max %d megapixels", cfg.Width, cfg.Height, maxImagePixels>>20)
	}
	_, imgType, err := image.Decode(bytes.NewBuffer(input))
	if err != nil {
		return "", damaged(err)
	}
	return imgType, nil
}

// ImageType reads only the header. Use it for stored images that were validated on upload.
func ImageType(input []byte) (string, error) {
	_, imgType, err := image.DecodeConfig(bytes.NewReader(input))
	if err != nil {
		return "", damaged(err)
	}
	return imgType, nil
}
func CreateImage(width, height int) *image.RGBA {
	return image.NewRGBA(image.Rect(0, 0, width, height))
}
func Draw(dst *image.RGBA, col, row int, src image.Image) {
	pos := image.Rect(
		col*src.Bounds().Dx(),                   // Start X
		row*src.Bounds().Dy(),                   // Start Y
		col*src.Bounds().Dx()+src.Bounds().Dx(), // End X
		row*src.Bounds().Dy()+src.Bounds().Dy(), // End Y
	)
	draw.Draw(dst, pos, src, image.Point{}, draw.Src)
}

func ImageFromReader(r io.Reader) (image.Image, error) {
	img, _, err := image.Decode(r)
	if err != nil {
		return nil, err
	}
	return img, nil
}
func ImageFromBinary(data []byte) (image.Image, error) {
	return ImageFromReader(bytes.NewReader(data))
}

func ImageToPng(img image.Image) ([]byte, error) {
	var res []byte
	w := bytes.NewBuffer(res)
	err := png.Encode(w, img)
	if err != nil {
		return nil, err
	}
	return w.Bytes(), nil
}
func ImageToJpeg(img image.Image) ([]byte, error) {
	var res []byte
	w := bytes.NewBuffer(res)
	err := jpeg.Encode(w, img, nil)
	if err != nil {
		return nil, err
	}
	return w.Bytes(), nil
}
func ImageToGif(img image.Image) ([]byte, error) {
	var res []byte
	w := bytes.NewBuffer(res)
	err := gif.Encode(w, img, nil)
	if err != nil {
		return nil, err
	}
	return w.Bytes(), nil
}

// damaged reports a file whose format is known but which does not decode.
// The decoder's text goes to the log; the user reads a plain sentence.
func damaged(err error) error {
	logger.Warn.Println("image decode:", err.Error())
	return apperr.With(apperr.ErrUnsupportedImage, "the image file is damaged or incomplete")
}
