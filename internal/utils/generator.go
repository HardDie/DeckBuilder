package utils

import (
	"github.com/HardDie/DeckBuilder/internal/config"
)

// Allows you to find the minimum image size for all cards on the page
func CalculateGridSize(cardsNumber int) (cols, rows int) {
	cols = 10
	rows = 7
	maxCards := cols * rows
	for r := config.MinHeight; r <= config.MaxHeight; r++ {
		for c := config.MinWidth; c <= config.MaxWidth; c++ {
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
func CardIdToPageCoordinates(id, columns int) (column, row int) {
	row = id / columns
	column = id % columns
	return
}
