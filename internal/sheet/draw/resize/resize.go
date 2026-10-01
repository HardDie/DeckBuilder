// Lanczos-resize every face and the back in parallel, then draw in order.
package resize

import (
	"bytes"
	"image"
	"sync"

	"github.com/HardDie/DeckBuilder/internal/images"
	"github.com/HardDie/DeckBuilder/internal/sheet/fit"
	"github.com/HardDie/DeckBuilder/internal/sheet/grid"
)

// JPEG matches seq pixels. Draw stays one cell at a time.
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
	var wg sync.WaitGroup
	wg.Add(len(faces) + 1)
	for i, face := range faces {
		i, face := i, face
		go func() {
			defer wg.Done()
			out[i] = fit.Resize(face, size.Width, size.Height)
		}()
	}
	var shaded image.Image
	go func() {
		defer wg.Done()
		shaded = fit.Resize(back, size.Width, size.Height)
	}()
	wg.Wait()
	return out, shaded
}

func encode(img image.Image) ([]byte, error) {
	var buf bytes.Buffer
	if err := images.JpegSaveToWriter(&buf, img); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
