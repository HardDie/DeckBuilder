package generator_image

import (
	"fmt"
	"image"
	"log"
	"path/filepath"

	"tts_deck_build/internal/cards"
	"tts_deck_build/internal/collections"
	"tts_deck_build/internal/config"
	"tts_deck_build/internal/decks"
	"tts_deck_build/internal/fs"
	"tts_deck_build/internal/games"
	"tts_deck_build/internal/images"

	"github.com/disintegration/imaging"
)

func getListCards(gameID string) (*DeckArray, int, error) {
	deckArray := NewDeckArray()
	totalCountOfCards := 0

	// Check if the game exists
	gameService := games.NewService()
	gameItem, err := gameService.Item(gameID)
	if err != nil {
		return nil, 0, err
	}

	// Get collection list
	collectionService := collections.NewService()
	collectionItems, err := collectionService.List(gameItem.ID, "")
	if err != nil {
		return nil, 0, err
	}

	// Create deck and card service
	deckService := decks.NewService()
	cardService := cards.NewService()

	// Get a list of decks for each collection
	for _, collectionItem := range collectionItems {
		deckItems, err := deckService.List(gameItem.ID, collectionItem.ID, "")
		if err != nil {
			return nil, 0, err
		}
		// Get a list of cards for each deck
		for _, deckItem := range deckItems {
			// Create a unique description of the deck
			deckArray.SelectDeck(deckItem.ID, deckItem.BacksideImage)
			cardItems, err := cardService.List(gameItem.ID, collectionItem.ID, deckItem.ID, "")
			if err != nil {
				return nil, 0, err
			}
			for _, cardItem := range cardItems {
				// Add card a card to the linked unique deck
				deckArray.AddCard(gameItem.ID, collectionItem.ID, cardItem.ID)
				totalCountOfCards++
			}
		}
	}
	return deckArray, totalCountOfCards, nil
}

func GenerateImagesForGame(gameID string) error {
	deckArray, _, err := getListCards(gameID)
	if err != nil {
		return err
	}

	// Create deck and card service
	deckService := decks.NewService()
	cardService := cards.NewService()

	for deckType, pages := range deckArray.Decks {
		// Processing one type of deck

		log.Printf("Deck: %q\n", deckType.Title)
		// If there are many cards in the deck, then one image page may not be enough.
		// Processing each page of the image.
		for i, page := range pages.Pages {
			// Getting the first image from the page
			firstCard := page[0]
			// Extracting the size of the image
			imgBin, _, err := cardService.GetImage(firstCard.GameID, firstCard.CollectionID, deckType.Title, firstCard.CardID)
			if err != nil {
				return err
			}
			cardWidth, cardHeight, err := images.ImageSize(imgBin)
			if err != nil {
				return err
			}
			// Calculation the optimal proportion of the image.
			// Add one card to the bottom right place for the backside image.
			columns, rows := calculateGridSize(len(page) + 1)
			// Calculating the resolution of the resulting image
			resultImageWidth := cardWidth * columns
			resultImageHeight := cardHeight * rows
			// Creating a page image
			pageImage := images.CreateImage(resultImageWidth, resultImageHeight)
			// Getting an image of the backside
			deckBinImg, _, err := deckService.GetImage(firstCard.GameID, firstCard.CollectionID, deckType.Title)
			if err != nil {
				return err
			}
			deckImg, err := images.ImageFromBinary(deckBinImg)
			if err != nil {
				return err
			}
			// Make the backside image slightly darker than the original image
			darkerDeckImg := imaging.AdjustBrightness(deckImg, -30)
			for i, card := range page {
				// Get image
				imgBin, _, err := cardService.GetImage(card.GameID, card.CollectionID, deckType.Title, card.CardID)
				if err != nil {
					return err
				}
				// Converting binary data to image type
				cardImg, err := images.ImageFromBinary(imgBin)
				if err != nil {
					return err
				}
				// Calculate the position of the image on the page
				column, row := cardIdToPageCoordinates(i, columns)
				// Draw an image on the page
				images.Draw(pageImage, column, row, cardImg)
			}
			// Draw a picture of the backside in the bottom right position
			// images.Draw(pageImage, columns-1, rows-1, deckImg)
			images.Draw(pageImage, columns-1, rows-1, darkerDeckImg)
			// Build the file name of the result page
			pageFileName := fmt.Sprintf("%s_%d_%d_%dx%d.png", deckType.Title, i+1, len(page), columns, rows)
			// Save the image page to file
			err = fs.CreateAndProcess[image.Image](filepath.Join(config.GetConfig().Results(), pageFileName), pageImage, images.SaveToWriter)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
