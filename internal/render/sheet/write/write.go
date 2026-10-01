// Draw one measured page and write a libjpeg-turbo quality-80 JPEG.
package write

import (
	"image"

	"github.com/HardDie/DeckBuilder/internal/fs"
	"github.com/HardDie/DeckBuilder/internal/images"
	"github.com/HardDie/DeckBuilder/internal/render/sheet/back"
	"github.com/HardDie/DeckBuilder/internal/render/sheet/draw/libjpeg"
	"github.com/HardDie/DeckBuilder/internal/render/sheet/fit"
	"github.com/HardDie/DeckBuilder/internal/render/sheet/paint"
)

// Draw paints faces and the back into path.
// Cell width and height are already chosen. Faces and back are the original file bytes.
func Draw(faces [][]byte, rawBack []byte, cellW, cellH int, shadow bool, path string) error {
	drawn := make([]image.Image, len(faces))
	for i, raw := range faces {
		img, err := images.ImageFromBinary(raw)
		if err != nil {
			return err
		}
		drawn[i] = fit.Resize(img, cellW, cellH)
	}
	decoded, err := images.ImageFromBinary(rawBack)
	if err != nil {
		return err
	}
	shaded := fit.Resize(back.Shade(decoded, shadow), cellW, cellH)
	sheet, _, _ := paint.Canvas(cellW, cellH, drawn, shaded)
	body, err := libjpeg.Encode(sheet)
	if err != nil {
		return err
	}
	return fs.CreateAndProcess(path, body, fs.BinToWriter)
}
