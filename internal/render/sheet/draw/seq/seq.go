// Sequential sheet: Lanczos, then one cell at a time, then JPEG quality 80.
package seq

import (
	"bytes"
	"image"

	"github.com/HardDie/DeckBuilder/internal/images"
	"github.com/HardDie/DeckBuilder/internal/render/sheet/fit"
	"github.com/HardDie/DeckBuilder/internal/render/sheet/grid"
)

// Image draws one page. It does not encode.
func Image(faces []image.Image, back image.Image) *image.RGBA {
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
	return page
}

// JPEG draws faces left to right, top to bottom, and the back in the bottom-right cell.
func JPEG(faces []image.Image, back image.Image) ([]byte, error) {
	return encode(Image(faces, back))
}

func resizeAll(faces []image.Image, back image.Image) ([]image.Image, image.Image) {
	b := faces[0].Bounds().Max
	_, size := fit.Apply(b.X, b.Y, 1, fit.Size{})
	out := make([]image.Image, len(faces))
	for i, face := range faces {
		out[i] = fit.Resize(face, size.Width, size.Height)
	}
	return out, fit.Resize(back, size.Width, size.Height)
}

func encode(img image.Image) ([]byte, error) {
	var buf bytes.Buffer
	if err := images.JpegSaveToWriter(&buf, img); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
