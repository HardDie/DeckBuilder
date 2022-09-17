package generator

import (
	"crypto/md5"
	"fmt"
	"image"
	"path/filepath"
	"strings"

	"github.com/disintegration/imaging"

	"tts_deck_build/internal/cards"
	"tts_deck_build/internal/collections"
	"tts_deck_build/internal/config"
	"tts_deck_build/internal/decks"
	"tts_deck_build/internal/dto"
	"tts_deck_build/internal/entity"
	"tts_deck_build/internal/fs"
	"tts_deck_build/internal/games"
	"tts_deck_build/internal/images"
	"tts_deck_build/internal/progress"
	"tts_deck_build/internal/tts_entity"
	"tts_deck_build/internal/utils"
)

type GeneratorService struct {
	cfg *config.Config
}

func NewService(cfg *config.Config) *GeneratorService {
	return &GeneratorService{
		cfg: cfg,
	}
}

func (s *GeneratorService) GenerateGame(gameID string, dtoObject *dto.GenerateGameDTO) error {
	pr := progress.GetProgress()

	deckArray, totalCountOfCards, err := s.getListOfCards(gameID, dtoObject.SortOrder)
	if err != nil {
		return err
	}

	// Cleanup before generation
	err = fs.RemoveFolder(s.cfg.Results())
	if err != nil {
		return err
	}

	// Create result folder
	err = fs.CreateFolder(s.cfg.Results())
	if err != nil {
		return err
	}

	pr.SetType("Image generation")
	err = s.generateBody(deckArray, totalCountOfCards)
	if err != nil {
		return err
	}

	return nil
}
func (s *GeneratorService) getListOfCards(gameID string, sortField string) (*entity.DeckArray, int, error) {
	deckArray := entity.NewDeckArray()
	totalCountOfCards := 0

	// Check if the game exists
	gameService := games.NewService(s.cfg)
	gameItem, err := gameService.Item(gameID)
	if err != nil {
		return nil, 0, err
	}

	// Get collection list
	collectionService := collections.NewService(s.cfg)
	collectionItems, err := collectionService.List(gameItem.ID, sortField)
	if err != nil {
		return nil, 0, err
	}

	// Create deck and card service
	deckService := decks.NewService(s.cfg)
	cardService := cards.NewService(s.cfg)

	// Get a list of decks for each collection
	for _, collectionItem := range collectionItems {
		deckItems, err := deckService.List(gameItem.ID, collectionItem.ID, sortField)
		if err != nil {
			return nil, 0, err
		}
		// Get a list of cards for each deck
		for _, deckItem := range deckItems {
			// Create a unique description of the deck
			deckArray.SelectDeck(deckItem.ID, deckItem.BacksideImage)
			cardItems, err := cardService.List(gameItem.ID, collectionItem.ID, deckItem.ID, sortField)
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
func (s *GeneratorService) generateBody(deckArray *entity.DeckArray, totalCountOfCards int) error {
	pr := progress.GetProgress()
	pr.SetMessage("Reading a list of cards from the disk...")

	var err error

	// Create deck and card service
	deckService := decks.NewService(s.cfg)
	cardService := cards.NewService(s.cfg)

	transform := tts_entity.Transform{
		ScaleX: 1,
		ScaleY: 1,
		ScaleZ: 1,
	}
	bag := tts_entity.Bag{
		Name:      "Bag",
		Transform: transform,
	}
	deck := tts_entity.DeckObject{
		CustomDeck: make(map[int]tts_entity.DeckDescription),
		Transform:  transform,
	}

	var processedCards int

	pr.SetMessage("Generating the resulting image pages...")
	pr.SetProgress(0)
	for deckInfo, pages := range deckArray.Decks {
		var collectionType string
		var deckItem *entity.DeckInfo
		var deckBacksideImage []byte
		var deckBacksideImageDarker *image.NRGBA
		var deckBacksideImageName string
		var deckDesc tts_entity.DeckDescription

		for pageId, page := range pages.Pages {
			var pageInfo *entity.PageInfo
			var pageImage *image.RGBA
			pr.SetMessage("Drawing cards on the resulting page...")
			for cardId, card := range page {
				// Preparation

				if deckBacksideImage == nil {
					// Getting an deck item
					deckItem, err = deckService.Item(card.GameID, card.CollectionID, deckInfo.DeckID)
					if err != nil {
						return err
					}
					// Getting an image of the backside
					deckBacksideImage, _, err = deckService.GetImage(card.GameID, card.CollectionID, deckInfo.DeckID)
					if err != nil {
						return err
					}
					deckBacksideImg, err := images.ImageFromBinary(deckBacksideImage)
					if err != nil {
						return err
					}
					// Make the backside image slightly darker than the original image
					deckBacksideImageDarker = imaging.AdjustBrightness(deckBacksideImg, -30)

					hash := md5.Sum([]byte(deckItem.BacksideImage))
					deckBacksideImageName = "backside_" + deckItem.ID + "_" + fmt.Sprintf("%x", hash[0:3]) + ".png"
					err = fs.CreateAndProcess(filepath.Join(s.cfg.Results(), deckBacksideImageName), deckBacksideImage, fs.BinToWriter)
					if err != nil {
						return err
					}
				}
				if pageInfo == nil {
					// Calculation the optimal proportion of the image.
					// Add one card to the bottom right place for the backside image.
					columns, rows := calculateGridSize(len(page) + 1)
					// Extracting the size of the card
					imgBin, _, err := cardService.GetImage(card.GameID, card.CollectionID, deckInfo.DeckID, card.CardID)
					if err != nil {
						return err
					}
					// Calculating the resolution of the resulting image
					cardWidth, cardHeight, err := images.ImageSize(imgBin)
					if err != nil {
						return err
					}
					pageInfo = &entity.PageInfo{
						Columns: columns,
						Rows:    rows,
						Width:   cardWidth * columns,
						Height:  cardHeight * rows,
						Count:   len(page),
						Name:    fmt.Sprintf("%s_%d_%d_%dx%d.png", deckInfo.DeckID, pageId+1, len(page), columns, rows),
					}
					pageImage = images.CreateImage(pageInfo.Width, pageInfo.Height)

					deckDesc = tts_entity.DeckDescription{
						FaceURL:   "file:///" + filepath.Join("home", "user", "data", pageInfo.Name),
						BackURL:   "file:///" + filepath.Join("home", "user", "data", deckBacksideImageName),
						NumWidth:  pageInfo.Columns,
						NumHeight: pageInfo.Rows,
					}
					deck.CustomDeck[pageId+1] = deckDesc
				}

				// Processing image

				// Get card image
				cardImageBin, _, err := cardService.GetImage(card.GameID, card.CollectionID, deckInfo.DeckID, card.CardID)
				if err != nil {
					return err
				}
				// Converting binary data to image type
				cardImg, err := images.ImageFromBinary(cardImageBin)
				if err != nil {
					return err
				}
				// Calculate the position of the image on the page
				column, row := cardIdToPageCoordinates(cardId, pageInfo.Columns)
				// Draw an image on the page
				images.Draw(pageImage, column, row, cardImg)

				// Processing json

				// If the collection on the previous card is different,
				// we move the current deck to the object list and create a new deck
				if collectionType != card.CollectionID+deckInfo.DeckID {
					collectionType = card.CollectionID + deckInfo.DeckID

					switch {
					case len(deck.ContainedObjects) == 1:
						// We cannot create a deck object with a single card. We must create a card object.
						bag.ContainedObjects = append(bag.ContainedObjects, deck.ContainedObjects[0])
					case len(deck.ContainedObjects) > 1:
						// If there is more than one card in the deck, place the deck in the object list.
						bag.ContainedObjects = append(bag.ContainedObjects, deck)
					}

					// Create a new deck object
					deck = tts_entity.DeckObject{
						Name:     "Deck",
						Nickname: deckItem.Type.String(),
						CustomDeck: map[int]tts_entity.DeckDescription{
							pageId + 1: deckDesc,
						},
						Transform: transform,
					}
				}
				// Get information about the card
				cardItem, err := cardService.Item(card.GameID, card.CollectionID, deckInfo.DeckID, card.CardID)
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
					Nickname:    utils.Allocate(cardItem.Title.String()),
					Description: utils.Allocate(cardItem.Description.String()),
					CardID:      (pageId+1)*100 + cardId,
					LuaScript:   strings.Join(variables, "\n"),
					CustomDeck: map[int]tts_entity.DeckDescription{
						pageId + 1: deckDesc,
					},
					Transform: &transform,
				})

				processedCards++
				pr.SetProgress(float32(processedCards) / float32(totalCountOfCards) * 100)
			}
			// Draw a picture of the backside in the bottom right position
			pr.SetMessage("Drawing backside image on the resulting page...")
			images.Draw(pageImage, pageInfo.Columns-1, pageInfo.Rows-1, deckBacksideImageDarker)
			// Save the image page to file
			pr.SetMessage("Saving the resulting page to disk...")
			err = fs.CreateAndProcess[image.Image](filepath.Join(s.cfg.Results(), pageInfo.Name), pageImage, images.SaveToWriter)
			if err != nil {
				return err
			}
		}

		switch {
		case len(deck.ContainedObjects) == 1:
			// We cannot create a deck object with a single card. We must create a card object.
			bag.ContainedObjects = append(bag.ContainedObjects, deck.ContainedObjects[0])
			deck = tts_entity.DeckObject{CustomDeck: make(map[int]tts_entity.DeckDescription)}
		case len(deck.ContainedObjects) > 1:
			// If there is more than one card in the deck, place the deck in the object list.
			bag.ContainedObjects = append(bag.ContainedObjects, deck)
			deck = tts_entity.DeckObject{CustomDeck: make(map[int]tts_entity.DeckDescription)}
		}
	}

	root := tts_entity.RootObjects{
		ObjectStates: []tts_entity.Bag{
			bag,
		},
	}
	err = fs.CreateAndProcess(filepath.Join(s.cfg.Results(), "decks.json"), root, fs.JsonToWriter[tts_entity.RootObjects])
	if err != nil {
		return err
	}

	pr.SetMessage("All image pages were successfully generated!")
	return nil
}
