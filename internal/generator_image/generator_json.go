package generator_image

import (
	"crypto/md5"
	"fmt"
	"path/filepath"
	"strings"

	"tts_deck_build/internal/cards"
	"tts_deck_build/internal/config"
	"tts_deck_build/internal/decks"
	"tts_deck_build/internal/fs"
	"tts_deck_build/internal/tts_entity"
)

func allocate[T any](val T) *T {
	return &val
}

func GenerateJsonForTTS(deckArray *DeckArray) error {
	transform := tts_entity.Transform{
		ScaleX: 1,
		ScaleY: 1,
		ScaleZ: 1,
	}

	// Create deck and card service
	deckService := decks.NewService()
	cardService := cards.NewService()

	bag := tts_entity.Bag{
		Name:      "Bag",
		Transform: transform,
	}
	deck := tts_entity.DeckObject{
		CustomDeck: make(map[int]tts_entity.DeckDescription),
		Transform:  transform,
	}

	var collectionType string

	for deckType, pages := range deckArray.Decks {
		for pageId, page := range pages.Pages {
			var deckDesc *tts_entity.DeckDescription
			for cardId, card := range page {
				// If we started the iteration on a new page
				if deckDesc == nil {
					// Calculation the optimal proportion of the image.
					columns, rows := calculateGridSize(len(page) + 1)
					// Build file page name
					pagePath := fmt.Sprintf("file:///home/user/data/%s_%d_%d_%dx%d.png", deckType.DeckID, pageId+1, len(page), columns, rows)
					// Get information about the deck
					deckItem, err := deckService.Item(card.GameID, card.CollectionID, deckType.DeckID)
					if err != nil {
						return err
					}
					hash := md5.Sum([]byte(deckItem.BacksideImage))
					backside := "backside_" + deckItem.ID + "_" + fmt.Sprintf("%x", hash[0:3]) + ".png"
					// Build a description for the image page
					deckDesc = &tts_entity.DeckDescription{
						FaceURL:   pagePath,
						BackURL:   "file:///home/user/data/" + backside,
						NumWidth:  columns,
						NumHeight: rows,
					}
					// Add page information to the deck
					if _, ok := deck.CustomDeck[pageId+1]; !ok {
						deck.CustomDeck[pageId+1] = *deckDesc
					}
				}

				// If the collection on the previous card is different,
				// we move the current deck to the object list and create a new deck
				if collectionType != card.CollectionID+deckType.DeckID {
					collectionType = card.CollectionID + deckType.DeckID

					switch {
					case len(deck.ContainedObjects) == 1:
						// We cannot create a deck object with a single card. We must create a card object.
						bag.ContainedObjects = append(bag.ContainedObjects, deck.ContainedObjects[0])
					case len(deck.ContainedObjects) > 1:
						// If there is more than one card in the deck, place the deck in the object list.
						bag.ContainedObjects = append(bag.ContainedObjects, deck)
					}

					// Get information about the deck
					deckItem, err := deckService.Item(card.GameID, card.CollectionID, deckType.DeckID)
					if err != nil {
						return err
					}

					// Create a new deck object
					deck = tts_entity.DeckObject{
						Name:     "Deck",
						Nickname: deckItem.Type.String(),
						CustomDeck: map[int]tts_entity.DeckDescription{
							pageId + 1: *deckDesc,
						},
						Transform: transform,
					}
				}

				// Get information about the card
				cardItem, err := cardService.Item(card.GameID, card.CollectionID, deckType.DeckID, card.CardID)
				if err != nil {
					return err
				}
				// Place the card ID in the list of cards inside the deck object
				deck.DeckIDs = append(deck.DeckIDs, (pageId+1)*100+cardId)
				// Converting lua variables into strings
				var variables []string
				for key, value := range cardItem.Variables {
					variables = append(variables, key+"="+value)
				}
				// Create a card and place it in the list of cards inside the deck
				deck.ContainedObjects = append(deck.ContainedObjects, tts_entity.Card{
					Name:        "Card",
					Nickname:    allocate(cardItem.Title.String()),
					Description: allocate(cardItem.Description.String()),
					CardID:      (pageId+1)*100 + cardId,
					LuaScript:   strings.Join(variables, "\n"),
					CustomDeck: map[int]tts_entity.DeckDescription{
						pageId + 1: *deckDesc,
					},
					Transform: &transform,
				})
			}
		}
	}

	// If there are cards in the deck, after iterating over all the cards, place the deck in the list of objects
	switch {
	case len(deck.ContainedObjects) == 1:
		// We cannot create a deck object with a single card. We must create a card object.
		bag.ContainedObjects = append(bag.ContainedObjects, deck.ContainedObjects[0])
	case len(deck.ContainedObjects) > 1:
		// If there is more than one card in the deck, place the deck in the object list.
		bag.ContainedObjects = append(bag.ContainedObjects, deck)
	}

	root := tts_entity.RootObjects{
		ObjectStates: []tts_entity.Bag{
			bag,
		},
	}

	err := fs.CreateAndProcess(filepath.Join(config.GetConfig().Results(), "decks.json"), root, fs.JsonToWriter[tts_entity.RootObjects])
	if err != nil {
		return err
	}
	return nil
}
