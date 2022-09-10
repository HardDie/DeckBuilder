package generator_image

import (
	"tts_deck_build/internal/cards"
	"tts_deck_build/internal/collections"
	"tts_deck_build/internal/config"
	"tts_deck_build/internal/decks"
	"tts_deck_build/internal/fs"
	"tts_deck_build/internal/games"
	"tts_deck_build/internal/progress"
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

func Generate(gameID string) error {
	pr := progress.GetProgress()

	deckArray, totalCountOfCards, err := getListCards(gameID)
	if err != nil {
		return err
	}

	// Create folder
	err = fs.CreateFolder(config.GetConfig().Results())
	if err != nil {
		return err
	}

	pr.SetType("Image generation")
	err = GenerateImagesForGame(deckArray, totalCountOfCards)
	if err != nil {
		return err
	}

	pr.SetType("Json generation")
	err = GenerateJsonForTTS(deckArray)
	if err != nil {
		return err
	}

	return nil
}
