package generator

import (
	"tts_deck_build/internal/entity"
)

// Allows you to find the minimum image size for all cards on the page
func calculateGridSize(cardsNumber int) (cols, rows int) {
	cols = 10
	rows = 7
	maxCards := cols * rows
	for r := entity.MinHeight; r <= entity.MaxHeight; r++ {
		for c := entity.MinWidth; c <= entity.MaxWidth; c++ {
			possible := c * r
			if possible < maxCards && possible >= cardsNumber {
				maxCards = possible
				cols = c
				rows = r
			}
		}
	}
	return
}

// Allows you to calculate the position of the card on the page by its identifier
func cardIdToPageCoordinates(id, columns int) (column, row int) {
	row = id / columns
	column = id % columns
	return
}
