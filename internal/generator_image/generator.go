package generator_image

import (
	"tts_deck_build/internal/config"
	"tts_deck_build/internal/fs"
	"tts_deck_build/internal/status"
)

func Generate(gameID string) error {
	st := status.GetStatus()

	deckArray, totalCountOfCards, err := getListCards(gameID)
	if err != nil {
		return err
	}

	// Create folder
	if err := fs.CreateFolder(config.GetConfig().Results()); err != nil {
		return err
	}

	st.SetType("Image generation")
	err = GenerateImagesForGame(deckArray, totalCountOfCards)
	if err != nil {
		return err
	}

	return nil
}
