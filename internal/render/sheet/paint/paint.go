// Draw faces and the back into one sheet image.
package paint

import (
	"image"

	"github.com/HardDie/DeckBuilder/internal/images"
	"github.com/HardDie/DeckBuilder/internal/render/sheet/grid"
)

// Canvas places faces left to right and the back in the bottom-right cell.
func Canvas(cellW, cellH int, faces []image.Image, back image.Image) (*image.RGBA, int, int) {
	cols, rows := grid.Size(len(faces) + 1)
	page := images.CreateImage(cellW*cols, cellH*rows)
	for i, face := range faces {
		col, row := grid.Slot(i, cols)
		images.Draw(page, col, row, face)
	}
	images.Draw(page, cols-1, rows-1, back)
	return page, cols, rows
}
