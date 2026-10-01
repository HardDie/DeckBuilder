// Columns and rows for one TTS custom-deck page.
package grid

import "github.com/HardDie/DeckBuilder/internal/config"

// Size is the smallest columns×rows that can hold n cells, from 2×2 to 10×7.
func Size(n int) (cols, rows int) {
	cols = 10
	rows = 7
	maxCards := cols * rows
	for r := config.MinHeight; r <= config.MaxHeight; r++ {
		for c := config.MinWidth; c <= config.MaxWidth; c++ {
			possible := c * r
			if possible < maxCards && possible >= n {
				maxCards = possible
				cols = c
				rows = r
			}
		}
	}
	return
}

// Slot is left to right, top to bottom.
func Slot(index, cols int) (col, row int) {
	row = index / cols
	col = index % cols
	return
}
