package generator_image

import (
	"crypto/md5"
	"fmt"
	"image"
	"log"
	"path/filepath"

	"tts_deck_build/internal/cards"
	"tts_deck_build/internal/config"
	"tts_deck_build/internal/decks"
	"tts_deck_build/internal/fs"
	"tts_deck_build/internal/images"
	"tts_deck_build/internal/progress"

	"github.com/disintegration/imaging"
)

func GenerateImagesForGame(deckArray *DeckArray, totalCountOfCards int) error {
	pr := progress.GetProgress()
	pr.SetMessage("Reading a list of cards from the disk...")

	var processedCards int
	backsideMark := make(map[string]struct{})

	// Create deck and card service
	deckService := decks.NewService()
	cardService := cards.NewService()

	pr.SetMessage("Generating the resulting image pages...")
	pr.SetProgress(0)
	for deckType, pages := range deckArray.Decks {
		// Processing one type of deck

		log.Printf("Deck: %q\n", deckType.DeckID)
		// If there are many cards in the deck, then one image page may not be enough.
		// Processing each page of the image.
		for pageId, page := range pages.Pages {
			// Getting the first image from the page
			firstCard := page[0]
			// Extracting the size of the image
			imgBin, _, err := cardService.GetImage(firstCard.GameID, firstCard.CollectionID, deckType.DeckID, firstCard.CardID)
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
			// Getting an deck item
			deckItem, err := deckService.Item(firstCard.GameID, firstCard.CollectionID, deckType.DeckID)
			if err != nil {
				return err
			}
			// Getting an image of the backside
			deckBinImg, _, err := deckService.GetImage(firstCard.GameID, firstCard.CollectionID, deckType.DeckID)
			if err != nil {
				return err
			}
			if _, ok := backsideMark[deckItem.ID+deckItem.BacksideImage]; !ok {
				backsideMark[deckItem.ID+deckItem.BacksideImage] = struct{}{}
				hash := md5.Sum([]byte(deckItem.BacksideImage))
				backside := "backside_" + deckItem.ID + "_" + fmt.Sprintf("%x", hash[0:3]) + ".png"
				err = fs.CreateAndProcess(filepath.Join(config.GetConfig().Results(), backside), deckBinImg, fs.BinToWriter)
				if err != nil {
					return err
				}
			}
			deckImg, err := images.ImageFromBinary(deckBinImg)
			if err != nil {
				return err
			}
			// Make the backside image slightly darker than the original image
			darkerDeckImg := imaging.AdjustBrightness(deckImg, -30)
			pr.SetMessage("Drawing cards on the resulting page...")
			for cardId, card := range page {
				// Get image
				imgBin, _, err := cardService.GetImage(card.GameID, card.CollectionID, deckType.DeckID, card.CardID)
				if err != nil {
					return err
				}
				// Converting binary data to image type
				cardImg, err := images.ImageFromBinary(imgBin)
				if err != nil {
					return err
				}
				// Calculate the position of the image on the page
				column, row := cardIdToPageCoordinates(cardId, columns)
				// Draw an image on the page
				images.Draw(pageImage, column, row, cardImg)
				processedCards++
				pr.SetProgress(float32(processedCards) / float32(totalCountOfCards) * 100)
			}
			pr.SetMessage("Drawing backside image on the resulting page...")
			// Draw a picture of the backside in the bottom right position
			// images.Draw(pageImage, columns-1, rows-1, deckImg)
			images.Draw(pageImage, columns-1, rows-1, darkerDeckImg)
			// Build the file name of the result page
			pageFileName := fmt.Sprintf("%s_%d_%d_%dx%d.png", deckType.DeckID, pageId+1, len(page), columns, rows)
			// Save the image page to file
			pr.SetMessage("Saving the resulting page to disk...")
			err = fs.CreateAndProcess[image.Image](filepath.Join(config.GetConfig().Results(), pageFileName), pageImage, images.SaveToWriter)
			if err != nil {
				return err
			}
		}
	}
	pr.SetMessage("All image pages were successfully generated!")
	return nil
}
