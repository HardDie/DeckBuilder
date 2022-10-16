package service

import (
	"crypto/md5"
	"fmt"
	"image"
	"path/filepath"
	"strings"
	"time"

	"github.com/disintegration/imaging"

	"github.com/HardDie/DeckBuilder/internal/config"
	"github.com/HardDie/DeckBuilder/internal/dto"
	"github.com/HardDie/DeckBuilder/internal/entity"
	"github.com/HardDie/DeckBuilder/internal/fs"
	"github.com/HardDie/DeckBuilder/internal/images"
	"github.com/HardDie/DeckBuilder/internal/logger"
	"github.com/HardDie/DeckBuilder/internal/progress"
	"github.com/HardDie/DeckBuilder/internal/tts_entity"
	"github.com/HardDie/DeckBuilder/internal/utils"
)

type IGeneratorService interface {
	GenerateGame(gameID string, dtoObject *dto.GenerateGameDTO) error
}
type GeneratorService struct {
	cfg               *config.Config
	gameService       IGameService
	collectionService ICollectionService
	deckService       IDeckService
	cardService       ICardService
}

func NewGeneratorService(cfg *config.Config, gameService IGameService, collectionService ICollectionService, deckService IDeckService, cardService ICardService) *GeneratorService {
	return &GeneratorService{
		cfg:               cfg,
		gameService:       gameService,
		collectionService: collectionService,
		deckService:       deckService,
		cardService:       cardService,
	}
}

func (s *GeneratorService) GenerateGame(gameID string, dtoObject *dto.GenerateGameDTO) error {
	pr := progress.GetProgress()

	// Check if the game exists
	gameItem, err := s.gameService.Item(gameID)
	if err != nil {
		return err
	}

	deckArray, totalCountOfCards, err := s.getListOfCards(gameItem.ID, dtoObject.SortOrder)
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
	pr.SetStatus(progress.StatusInProgress)
	go func() {
		err = s.generateBody(gameItem, deckArray, totalCountOfCards)
		if err != nil {
			pr.SetStatus(progress.StatusError)
			logger.Error.Println("Generator:", err.Error())
			return
		}
		pr.SetStatus(progress.StatusDone)
	}()

	return nil
}

func (s *GeneratorService) getListOfCards(gameID string, sortField string) (*entity.DeckArray, int, error) {
	deckArray := entity.NewDeckArray()
	totalCountOfCards := 0

	// Get collection list
	collectionItems, err := s.collectionService.List(gameID, sortField)
	if err != nil {
		return nil, 0, err
	}

	// Get a list of decks for each collection
	for _, collectionItem := range collectionItems {
		deckItems, err := s.deckService.List(gameID, collectionItem.ID, sortField)
		if err != nil {
			return nil, 0, err
		}
		// Get a list of cards for each deck
		for _, deckItem := range deckItems {
			// Create a unique description of the deck
			deckArray.SelectDeck(deckItem.ID, deckItem.Image)
			cardItems, err := s.cardService.List(gameID, collectionItem.ID, deckItem.ID, sortField, "")
			if err != nil {
				return nil, 0, err
			}
			for _, cardItem := range cardItems {
				// Add card a card to the linked unique deck
				deckArray.AddCard(gameID, collectionItem.ID, cardItem.ID, cardItem.Count)
				totalCountOfCards++
			}
		}
	}
	return deckArray, totalCountOfCards, nil
}
func (s *GeneratorService) generateBody(gameItem *entity.GameInfo, deckArray *entity.DeckArray, totalCountOfCards int) error {
	pr := progress.GetProgress()
	pr.SetMessage("Reading a list of cards from the disk...")

	var err error

	transform := tts_entity.Transform{
		ScaleX: 1,
		ScaleY: 1,
		ScaleZ: 1,
	}
	bag := tts_entity.Bag{
		Name:      "Bag",
		Nickname:  gameItem.Name,
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
		if len(pages.Pages) == 0 {
			// If deck is empty, skip
			continue
		}

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
					deckItem, err = s.deckService.Item(card.GameID, card.CollectionID, deckInfo.DeckID)
					if err != nil {
						return err
					}
					// Getting an image of the backside
					deckBacksideImage, _, err = s.deckService.GetImage(card.GameID, card.CollectionID, deckInfo.DeckID)
					if err != nil {
						return err
					}
					deckBacksideImg, err := images.ImageFromBinary(deckBacksideImage)
					if err != nil {
						return err
					}
					// Make the backside image slightly darker than the original image
					deckBacksideImageDarker = imaging.AdjustBrightness(deckBacksideImg, -30)

					hash := md5.Sum([]byte(deckItem.Image))
					deckBacksideImageName = "backside_" + deckItem.ID + "_" + fmt.Sprintf("%x", hash[0:3]) + ".png"
					err = fs.CreateAndProcess(filepath.Join(s.cfg.Results(), deckBacksideImageName), deckBacksideImage, fs.BinToWriter)
					if err != nil {
						return err
					}
				}
				if pageInfo == nil {
					// Calculation the optimal proportion of the image.
					// Add one card to the bottom right place for the backside image.
					columns, rows := utils.CalculateGridSize(len(page) + 1)
					// Extracting the size of the card
					imgBin, _, err := s.cardService.GetImage(card.GameID, card.CollectionID, deckInfo.DeckID, card.CardID)
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
						FaceURL:   "file:///" + fs.PathToAbsolutePath(filepath.Join(s.cfg.Results(), pageInfo.Name)),
						BackURL:   "file:///" + fs.PathToAbsolutePath(filepath.Join(s.cfg.Results(), deckBacksideImageName)),
						NumWidth:  pageInfo.Columns,
						NumHeight: pageInfo.Rows,
					}
					deck.CustomDeck[pageId+1] = deckDesc
				}

				// Processing image

				// Get card image
				cardImageBin, _, err := s.cardService.GetImage(card.GameID, card.CollectionID, deckInfo.DeckID, card.CardID)
				if err != nil {
					return err
				}
				// Converting binary data to image type
				cardImg, err := images.ImageFromBinary(cardImageBin)
				if err != nil {
					return err
				}
				// Calculate the position of the image on the page
				column, row := utils.CardIdToPageCoordinates(cardId, pageInfo.Columns)
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
						Nickname: deckItem.Name,
						CustomDeck: map[int]tts_entity.DeckDescription{
							pageId + 1: deckDesc,
						},
						Transform: transform,
					}
				}
				// Get information about the card
				cardItem, err := s.cardService.Item(card.GameID, card.CollectionID, deckInfo.DeckID, card.CardID)
				if err != nil {
					return err
				}
				// Converting lua variables into strings
				var variables []string
				for key, value := range cardItem.Variables {
					variables = append(variables, key+"="+value)
				}

				// Add a card to the deck as many times as set in the count variable
				for i := 0; i < cardItem.Count; i++ {
					// Place the card ID in the list of cards inside the deck object
					deck.DeckIDs = append(deck.DeckIDs, (pageId+1)*100+cardId)

					// Create a card and place it in the list of cards inside the deck
					deck.ContainedObjects = append(deck.ContainedObjects, tts_entity.Card{
						Name:        "Card",
						Nickname:    utils.Allocate(cardItem.Name),
						Description: utils.Allocate(cardItem.Description),
						CardID:      (pageId+1)*100 + cardId,
						LuaScript:   strings.Join(variables, "\n"),
						CustomDeck: map[int]tts_entity.DeckDescription{
							pageId + 1: deckDesc,
						},
						Transform: &transform,
					})
				}

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

	bag.Description = fmt.Sprintf("Created at: %v", time.Now().Format("2006-01-02 15:04:05"))
	root := tts_entity.RootObjects{
		ObjectStates: []tts_entity.Bag{
			bag,
		},
	}
	err = fs.CreateAndProcess(filepath.Join(s.cfg.Results(), gameItem.ID+".json"), root, fs.JsonToWriter[tts_entity.RootObjects])
	if err != nil {
		return err
	}

	pr.SetMessage("All image pages were successfully generated!")
	return nil
}
