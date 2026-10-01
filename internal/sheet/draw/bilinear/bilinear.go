// ApproxBiLinear resize, then the same sequential cell draw as seq.
package bilinear

import (
	"bytes"
	"image"

	xdraw "golang.org/x/image/draw"

	"github.com/HardDie/DeckBuilder/internal/images"
	"github.com/HardDie/DeckBuilder/internal/sheet/fit"
	"github.com/HardDie/DeckBuilder/internal/sheet/grid"
)

// JPEG may differ from seq when a face is resampled.
func JPEG(faces []image.Image, back image.Image) ([]byte, error) {
	faces, back = resizeAll(faces, back)
	cols, rows := grid.Size(len(faces) + 1)
	cellW := faces[0].Bounds().Dx()
	cellH := faces[0].Bounds().Dy()
	page := images.CreateImage(cellW*cols, cellH*rows)
	for i, face := range faces {
		col, row := grid.Slot(i, cols)
		images.Draw(page, col, row, face)
	}
	images.Draw(page, cols-1, rows-1, back)
	return encode(page)
}

func resizeAll(faces []image.Image, back image.Image) ([]image.Image, image.Image) {
	b := faces[0].Bounds().Max
	_, size := fit.Apply(b.X, b.Y, 1, fit.Size{})
	out := make([]image.Image, len(faces))
	for i, face := range faces {
		out[i] = scale(face, size.Width, size.Height)
	}
	return out, scale(back, size.Width, size.Height)
}

func scale(img image.Image, width, height int) image.Image {
	if width == img.Bounds().Max.X && height == img.Bounds().Max.Y {
		return img
	}
	dst := image.NewRGBA(image.Rect(0, 0, width, height))
	xdraw.ApproxBiLinear.Scale(dst, dst.Bounds(), img, img.Bounds(), xdraw.Src, nil)
	return dst
}

func encode(img image.Image) ([]byte, error) {
	var buf bytes.Buffer
	if err := images.JpegSaveToWriter(&buf, img); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
