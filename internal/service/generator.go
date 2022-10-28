package service

import (
	"crypto/md5"
	"fmt"
	"image"
	"path/filepath"
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

	deckArray, err := s.getListOfCards(gameItem.ID, dtoObject.SortOrder)
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
		err = s.generateBody(gameItem, deckArray)
		if err != nil {
			pr.SetStatus(progress.StatusError)
			logger.Error.Println("Generator:", err.Error())
			return
		}
		pr.SetStatus(progress.StatusDone)
	}()

	return nil
}

func (s *GeneratorService) getListOfCards(gameID string, sortField string) (*entity.DeckArray, error) {
	deckArray := entity.NewDeckArray()

	// Get collection list
	collectionItems, err := s.collectionService.List(gameID, sortField, "")
	if err != nil {
		return nil, err
	}

	// Get a list of decks for each collection
	for _, collectionItem := range collectionItems {
		deckItems, err := s.deckService.List(gameID, collectionItem.ID, sortField, "")
		if err != nil {
			return nil, err
		}
		// Get a list of cards for each deck
		for _, deckItem := range deckItems {
			// Create a unique description of the deck
			deckArray.SelectDeck(deckItem.ID, deckItem.Image)
			cardItems, err := s.cardService.List(gameID, collectionItem.ID, deckItem.ID, sortField, "")
			if err != nil {
				return nil, err
			}
			for _, cardItem := range cardItems {
				// Add a card to the linked unique deck
				deckArray.AddCard(gameID, collectionItem.ID, cardItem.ID, cardItem.Count)
			}
		}
	}
	return deckArray, nil
}
func (s *GeneratorService) generateBody(gameItem *entity.GameInfo, deckArray *entity.DeckArray) error {
	pr := progress.GetProgress()
	pr.SetMessage("Reading a list of cards from the disk...")

	var err error

	bag := tts_entity.NewBag(gameItem.Name)
	deck := tts_entity.NewDeck("")

	var processedCards int

	pr.SetMessage("Generating the resulting image pages...")
	pr.SetProgress(0)
	for deckInfo, pages := range deckArray.Decks {
		if len(pages.Pages) == 0 {
			// If deck is empty, skip
			continue
		}

		var infoDeck DeckInformation

		for pageId, page := range pages.Pages {
			var infoPage *PageInformation
			pr.SetMessage("Drawing cards on the resulting page...")
			for cardIndex, card := range page {
				// Preparation

				if infoDeck.backside == nil {
					// Getting an deck item
					infoDeck.deckItem, err = s.deckService.Item(card.GameID, card.CollectionID, deckInfo.DeckID)
					if err != nil {
						return err
					}

					infoDeck.backside, err = s.prepareBacksideImage(card.GameID, card.CollectionID, deckInfo.DeckID, card.CardID)
					if err != nil {
						return err
					}
				}
				if infoPage == nil {
					infoPage, err = s.preparePageInfo(pageId, page, card, deckInfo, &infoDeck)
					if err != nil {
						return err
					}
					deck.CustomDeck[pageId+1] = infoDeck.deckDesc

					// Draw a picture of the backside in the bottom right position
					images.Draw(infoPage.image, infoPage.info.Columns-1, infoPage.info.Rows-1, infoDeck.backside.imageDarker)
				}

				// Processing image

				// Draw an image on the page
				err = s.drawImageOnPage(card.GameID, card.CollectionID, deckInfo.DeckID, card.CardID, cardIndex, infoPage)
				if err != nil {
					return err
				}

				// Processing json

				// If the collection on the previous card is different,
				// we move the current deck to the object list and create a new deck
				if infoDeck.collectionType != card.CollectionID+deckInfo.DeckID {
					infoDeck.collectionType = card.CollectionID + deckInfo.DeckID

					switch {
					case len(deck.ContainedObjects) == 1:
						// We cannot create a deck object with a single card. We must create a card object.
						bag.ContainedObjects = append(bag.ContainedObjects, deck.ContainedObjects[0])
					case len(deck.ContainedObjects) > 1:
						// If there is more than one card in the deck, place the deck in the object list.
						bag.ContainedObjects = append(bag.ContainedObjects, deck)
					}

					// Create a new deck object
					deck = tts_entity.NewDeck(infoDeck.deckItem.Name)
					deck.CustomDeck[pageId+1] = infoDeck.deckDesc
				}
				// Get information about the card
				cardItem, err := s.cardService.Item(card.GameID, card.CollectionID, deckInfo.DeckID, card.CardID)
				if err != nil {
					return err
				}

				cardObject := tts_entity.NewCard(cardItem.Name, cardItem.Description, pageId, cardIndex, cardItem.Variables, infoDeck.deckDesc)
				for i := 0; i < cardItem.Count; i++ {
					// Add a card to the deck as many times as set in the count variable
					deck.AddCard(cardObject)
				}

				processedCards++
				pr.SetProgress(float32(processedCards) / float32(deckArray.TotalCount) * 100)
			}
			// Save the image page to file
			pr.SetMessage("Saving the resulting page to disk...")
			err = fs.CreateAndProcess[image.Image](filepath.Join(s.cfg.Results(), infoPage.info.Name), infoPage.image, images.SaveToWriter)
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

func (s *GeneratorService) prepareBacksideImage(gameID, collectionID, deckID string, cardID int64) (*BackSideInformation, error) {
	// Get image of first card
	cardImageBin, _, err := s.cardService.GetImage(gameID, collectionID, deckID, cardID)
	if err != nil {
		return nil, err
	}
	// Get card image resolution
	cardWidth, cardHeight, err := images.ImageSize(cardImageBin)
	if err != nil {
		return nil, err
	}

	// Getting an deck item
	deckItem, err := s.deckService.Item(gameID, collectionID, deckID)
	if err != nil {
		return nil, err
	}
	// Getting an image of the backside
	deckBacksideImage, _, err := s.deckService.GetImage(gameID, collectionID, deckID)
	if err != nil {
		return nil, err
	}

	// Get backside image resolution
	backSideWidth, backSideHeight, err := images.ImageSize(deckBacksideImage)
	if err != nil {
		return nil, err
	}

	deckBacksideImg, err := images.ImageFromBinary(deckBacksideImage)
	if err != nil {
		return nil, err
	}

	// Resize the backside image in case it has a different size from the card
	if cardHeight != backSideWidth ||
		cardWidth != backSideHeight {
		deckBacksideImg = imaging.Resize(deckBacksideImg, cardWidth, cardHeight, imaging.Lanczos)
	}

	hash := md5.Sum(deckBacksideImage)
	name := "backside_" + deckItem.ID + "_" + fmt.Sprintf("%x", hash[0:3]) + ".png"
	err = fs.CreateAndProcess(filepath.Join(s.cfg.Results(), name), deckBacksideImage, fs.BinToWriter)
	if err != nil {
		return nil, err
	}

	backside := &BackSideInformation{
		imaging.AdjustBrightness(deckBacksideImg, -30),
		name,
	}

	return backside, nil
}
func (s *GeneratorService) preparePageInfo(pageId int, page entity.CardPage, card entity.CardDescription, deckInfo entity.DeckType, infoDeck *DeckInformation) (*PageInformation, error) {
	// Calculation the optimal proportion of the image.
	// Add one card to the bottom right place for the backside image.
	columns, rows := utils.CalculateGridSize(len(page) + 1)
	// Extracting the size of the card
	imgBin, _, err := s.cardService.GetImage(card.GameID, card.CollectionID, deckInfo.DeckID, card.CardID)
	if err != nil {
		return nil, err
	}
	// Calculating the resolution of the resulting image
	cardWidth, cardHeight, err := images.ImageSize(imgBin)
	if err != nil {
		return nil, err
	}
	infoPage := &PageInformation{}
	infoPage.info = &entity.PageInfo{
		Columns: columns,
		Rows:    rows,
		Width:   cardWidth * columns,
		Height:  cardHeight * rows,
		Count:   len(page),
		Name:    fmt.Sprintf("%s_%d_%d_%dx%d.png", deckInfo.DeckID, pageId+1, len(page), columns, rows),
	}
	infoPage.image = images.CreateImage(infoPage.info.Width, infoPage.info.Height)

	infoDeck.deckDesc = tts_entity.DeckDescription{
		FaceURL:   "file:///" + fs.PathToAbsolutePath(filepath.Join(s.cfg.Results(), infoPage.info.Name)),
		BackURL:   "file:///" + fs.PathToAbsolutePath(filepath.Join(s.cfg.Results(), infoDeck.backside.imageName)),
		NumWidth:  infoPage.info.Columns,
		NumHeight: infoPage.info.Rows,
	}

	return infoPage, nil
}
func (s *GeneratorService) drawImageOnPage(gameID, collectionID, deckID string, cardID int64, cardIndex int, infoPage *PageInformation) error {
	// Get card image
	cardImageBin, _, err := s.cardService.GetImage(gameID, collectionID, deckID, cardID)
	if err != nil {
		return err
	}
	// Converting binary data to image type
	cardImg, err := images.ImageFromBinary(cardImageBin)
	if err != nil {
		return err
	}
	// Calculate the position of the image on the page
	column, row := utils.CardIdToPageCoordinates(cardIndex, infoPage.info.Columns)
	// Draw an image on the page
	images.Draw(infoPage.image, column, row, cardImg)
	return nil
}

type BackSideInformation struct {
	imageDarker *image.NRGBA
	imageName   string
}
type DeckInformation struct {
	collectionType string
	deckItem       *entity.DeckInfo
	backside       *BackSideInformation
	deckDesc       tts_entity.DeckDescription
}
type PageInformation struct {
	info  *entity.PageInfo
	image *image.RGBA
}
