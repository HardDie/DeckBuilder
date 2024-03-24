package images

import (
	"bytes"
	"image"
	"image/draw"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"

	"github.com/HardDie/DeckBuilder/internal/errors"
)

func ValidateImage(input []byte) (string, error) {
	_, imgType, err := image.Decode(bytes.NewBuffer(input))
	if err != nil {
		return "", errors.UnknownImageType.AddMessage(err.Error())
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
func ImageSize(data []byte) (width, height int, err error) {
	img, err := ImageFromBinary(data)
	if err != nil {
		return
	}
	bound := img.Bounds().Max
	width = bound.X
	height = bound.Y
	return
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
func SaveToWriter(w io.Writer, img image.Image) error {
	err := png.Encode(w, img)
	if err != nil {
		return err
	}
	return nil
}
func JpegSaveToWriter(w io.Writer, img image.Image) error {
	err := jpeg.Encode(w, img, &jpeg.Options{
		Quality: 60,
	})
	if err != nil {
		return err
	}
	return nil
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
