package generator

import (
	"bytes"
	"fmt"
	"image/jpeg"
	"path/filepath"
	"sort"
	"strconv"
	"time"

	"github.com/HardDie/DeckBuilder/internal/config"
	entitiesGame "github.com/HardDie/DeckBuilder/internal/entities/game"
	entitiesSettings "github.com/HardDie/DeckBuilder/internal/entities/settings"
	"github.com/HardDie/DeckBuilder/internal/fs"
	"github.com/HardDie/DeckBuilder/internal/images"
	"github.com/HardDie/DeckBuilder/internal/logger"
	pageDrawer "github.com/HardDie/DeckBuilder/internal/page_drawer"
	"github.com/HardDie/DeckBuilder/internal/progress"
	servicesCard "github.com/HardDie/DeckBuilder/internal/services/card"
	servicesCollection "github.com/HardDie/DeckBuilder/internal/services/collection"
	servicesDeck "github.com/HardDie/DeckBuilder/internal/services/deck"
	servicesGame "github.com/HardDie/DeckBuilder/internal/services/game"
	servicesSystem "github.com/HardDie/DeckBuilder/internal/services/system"
	servicesTTS "github.com/HardDie/DeckBuilder/internal/services/tts"
	"github.com/HardDie/DeckBuilder/internal/tts_entity"
	"github.com/HardDie/DeckBuilder/internal/utils"
)

type generator struct {
	cfg               *config.Config
	serviceGame       servicesGame.Game
	serviceCollection servicesCollection.Collection
	serviceDeck       servicesDeck.Deck
	serviceCard       servicesCard.Card
	serviceSystem     servicesSystem.System
	serviceTTS        servicesTTS.TTS
}

func New(
	cfg *config.Config,
	serviceGame servicesGame.Game,
	serviceCollection servicesCollection.Collection,
	serviceDeck servicesDeck.Deck,
	serviceCard servicesCard.Card,
	serviceSystem servicesSystem.System,
	serviceTTS servicesTTS.TTS,
) Generator {
	return &generator{
		cfg:               cfg,
		serviceGame:       serviceGame,
		serviceCollection: serviceCollection,
		serviceDeck:       serviceDeck,
		serviceCard:       serviceCard,
		serviceSystem:     serviceSystem,
		serviceTTS:        serviceTTS,
	}
}

func (s *generator) GenerateGame(gameID string, req GenerateGameRequest) error {
	cfg, err := s.serviceSystem.GetSettings()
	if err != nil {
		logger.Error.Printf("can't get config")
		return err
	}

	pr := progress.GetProgress()

	// Check if the game exists
	gameItem, err := s.serviceGame.Item(gameID)
	if err != nil {
		return err
	}

	deckArray, order, err := s.getListOfCards(gameItem.ID, req.SortOrder)
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
		err = s.generateBody(gameItem, deckArray, order, req.Scale, cfg)
		if err != nil {
			pr.SetStatus(progress.StatusError)
			logger.Error.Println("Generator:", err.Error())
			return
		}
		pr.SetStatus(progress.StatusDone)
	}()

	return nil
}

type Deck struct {
	ID    string
	Name  string
	Image string
}
type Card struct {
	ID           int64
	GameID       string
	CollectionID string
	Count        int
}

func (s *generator) getListOfCards(gameID string, sortField string) (map[Deck][]Card, []Deck, error) {
	decks := make(map[Deck][]Card)
	// Get list of collections
	collectionItems, err := s.serviceCollection.List(gameID, sortField, "")
	if err != nil {
		return nil, nil, err
	}
	// Iterate through collections
	for _, collectionItem := range collectionItems {
		// Get list of decks
		deckItems, err := s.serviceDeck.List(gameID, collectionItem.ID, sortField, "")
		if err != nil {
			return nil, nil, err
		}
		// Iterate through decks
		for _, deckItem := range deckItems {
			// Create deck object
			deck := Deck{
				ID:    deckItem.ID,
				Name:  deckItem.Name,
				Image: deckItem.Image,
			}
			// Get list of cards
			cardItems, err := s.serviceCard.List(gameID, collectionItem.ID, deckItem.ID, sortField, "")
			if err != nil {
				return nil, nil, err
			}
			// Iterate through cards
			for _, cardItem := range cardItems {
				// Add card into deck
				decks[deck] = append(decks[deck], Card{
					ID:           cardItem.ID,
					GameID:       gameID,
					CollectionID: collectionItem.ID,
					Count:        cardItem.Count,
				})
			}
		}
	}

	var order []Deck
	for deck := range decks {
		order = append(order, deck)
	}
	sort.SliceStable(order, func(i, j int) bool {
		return order[i].Name < order[j].Name
	})
	return decks, order, nil
}

func (s *generator) generateBody(
	gameItem *entitiesGame.Game,
	decks map[Deck][]Card,
	order []Deck,
	scale int,
	cfg *entitiesSettings.Settings,
) error {
	pr := progress.GetProgress()
	pr.SetMessage("Reading a list of cards from the disk...")

	// Generate images
	imageMapping, err := s.generateImages(decks, order, scale, cfg)
	if err != nil {
		return err
	}
	// Generate json description
	err = s.generateJson(gameItem, decks, order, imageMapping, cfg)
	if err != nil {
		return err
	}

	return nil
}

type PageInfo struct {
	Image    string
	Backside string
	Columns  int
	Rows     int
}

// input:
//   - array of decks with array of cards
//
// output:
//   - images in result folder
//   - map[file_name] = info{ path_to_image, path_to_backside, width, height }
func (s *generator) generateImages(
	decks map[Deck][]Card,
	order []Deck,
	scale int,
	cfg *entitiesSettings.Settings,
) (map[string]PageInfo, error) {
	pr := progress.GetProgress()

	// Count total amount of cards
	var processedCards int
	var totalCount int
	for _, cards := range decks {
		totalCount += len(cards)
	}

	pr.SetMessage("Generating the resulting image pages...")
	pr.SetProgress(0)

	images := make(map[string]PageInfo)

	pr.SetMessage("Drawing cards on the page...")
	var commonIndex int
	for _, deckInfo := range order {
		cards := decks[deckInfo]
		commonIndex++

		// Create page drawer object
		page := pageDrawer.New(deckInfo.ID, s.cfg.Results(), scale, commonIndex, cfg)
		var backsidePath string

		// Iterate through all cards in deck
		for _, card := range cards {
			// Init page drawer with deck information
			if page.IsEmpty() {
				// Getting an image of the backside
				deckBacksideImage, _, err := s.serviceDeck.GetImage(card.GameID, card.CollectionID, deckInfo.ID)
				if err != nil {
					logger.Error.Printf("backside not found for: %s.%s.%s", card.GameID, card.CollectionID, deckInfo.ID)
					return nil, err
				}
				// Set backside image
				savePath, err := page.SetBacksideImageAndSave(deckBacksideImage)
				if err != nil {
					return nil, err
				}
				backsidePath = savePath
			}

			// Start new page if current is full
			if page.IsFull() {
				pr.SetMessage("Saving the resulting page to disk...")
				savePath, columns, rows, err := page.Save()
				if err != nil {
					return nil, err
				}
				images[deckInfo.ID+"_"+strconv.Itoa(page.GetIndex())] = PageInfo{
					Image:    savePath,
					Backside: backsidePath,
					Columns:  columns,
					Rows:     rows,
				}
				pr.SetMessage("Drawing cards on the page...")
				page = (&pageDrawer.PageDrawer{}).Inherit(page)
				commonIndex++
			}

			// Get card image
			cardImageBin, _, err := s.serviceCard.GetImage(card.GameID, card.CollectionID, deckInfo.ID, card.ID)
			if err != nil {
				logger.Error.Printf("card image not found for: %s.%s.%s.%d", card.GameID, card.CollectionID, deckInfo.ID, card.ID)
				return nil, err
			}
			// Add card on page
			err = page.AddImage(cardImageBin)
			if err != nil {
				return nil, err
			}

			// Progress
			processedCards++
			pr.SetProgress(float32(processedCards) / float32(totalCount) * 100)
		}

		if !page.IsEmpty() {
			pr.SetMessage("Saving the resulting page to disk...")
			savePath, columns, rows, err := page.Save()
			if err != nil {
				return nil, err
			}
			images[deckInfo.ID+"_"+strconv.Itoa(page.GetIndex())] = PageInfo{
				Image:    savePath,
				Backside: backsidePath,
				Columns:  columns,
				Rows:     rows,
			}
			pr.SetMessage("Drawing cards on the page...")
		}
	}
	pr.SetMessage("All image pages were successfully generated!")
	return images, nil
}

func (s *generator) generateJson(
	gameItem *entitiesGame.Game,
	decks map[Deck][]Card,
	order []Deck,
	imageMapping map[string]PageInfo,
	cfg *entitiesSettings.Settings,
) error {
	bag := tts_entity.NewBag(gameItem.Name)
	collectionBags := make(map[string]*tts_entity.Bag)
	var deck tts_entity.DeckObject

	var dummyImage []byte
	{
		buf := new(bytes.Buffer)
		img := images.CreateImage(10, 10)
		err := jpeg.Encode(buf, img, nil)
		if err != nil {
			return err
		}
		dummyImage = buf.Bytes()
	}

	var deckIdOffset int

	var commonIndex int
	for _, deckInfo := range order {
		cards := decks[deckInfo]
		commonIndex++

		deck = tts_entity.NewDeck(
			deckInfo.Name,
			tts_entity.Transform{
				ScaleX: cfg.CardSize.ScaleX,
				ScaleY: cfg.CardSize.ScaleY,
				ScaleZ: cfg.CardSize.ScaleZ,
			},
		)
		// Create page drawer object
		page := pageDrawer.New(deckInfo.ID, "", 1, commonIndex, cfg)

		pageInfo := imageMapping[deckInfo.ID+"_"+strconv.Itoa(page.GetIndex())]
		deckDescription := tts_entity.DeckDescription{
			FaceURL:   "file:///" + pageInfo.Image,
			BackURL:   "file:///" + pageInfo.Backside,
			NumWidth:  pageInfo.Columns,
			NumHeight: pageInfo.Rows,
		}
		deck.CustomDeck[page.GetIndex()+deckIdOffset] = deckDescription

		var prevCollection string
		var prevCollectionDeck string

		// Iterate through all cards in deck
		for _, card := range cards {
			if page.IsEmpty() {
				prevCollection = card.CollectionID
				prevCollectionDeck = card.CollectionID + deckInfo.ID
			}

			// Start new page if current is full
			if page.IsFull() {
				page = (&pageDrawer.PageDrawer{}).Inherit(page)
				commonIndex++

				pageInfo = imageMapping[deckInfo.ID+"_"+strconv.Itoa(page.GetIndex())]
				deckDescription = tts_entity.DeckDescription{
					FaceURL:   "file:///" + pageInfo.Image,
					BackURL:   "file:///" + pageInfo.Backside,
					NumWidth:  pageInfo.Columns,
					NumHeight: pageInfo.Rows,
				}
				deck.CustomDeck[page.GetIndex()+deckIdOffset] = deckDescription
			}

			if card.CollectionID+deckInfo.ID != prevCollectionDeck {
				prevCollectionDeck = card.CollectionID + deckInfo.ID

				if _, ok := collectionBags[prevCollection]; !ok {
					collectionBags[prevCollection] = utils.Allocate(tts_entity.NewBag(prevCollection))
				}
				switch {
				case len(deck.ContainedObjects) == 1:
					// We cannot create a deck object with a single card. We must create a card object.
					collectionBags[prevCollection].ContainedObjects = append(collectionBags[prevCollection].ContainedObjects, deck.ContainedObjects[0])
					// bag.ContainedObjects = append(bag.ContainedObjects, deck.ContainedObjects[0])
				case len(deck.ContainedObjects) > 1:
					// If there is more than one card in the deck, place the deck in the object list.
					// bag.ContainedObjects = append(bag.ContainedObjects, deck)
					collectionBags[prevCollection].ContainedObjects = append(collectionBags[prevCollection].ContainedObjects, deck)
				}
				prevCollection = card.CollectionID
				deck = tts_entity.NewDeck(
					deckInfo.Name,
					tts_entity.Transform{
						ScaleX: cfg.CardSize.ScaleX,
						ScaleY: cfg.CardSize.ScaleY,
						ScaleZ: cfg.CardSize.ScaleZ,
					},
				)
				deck.CustomDeck[page.GetIndex()+deckIdOffset] = deckDescription
			}

			// Add card on page
			err := page.AddImage(dummyImage)
			if err != nil {
				return err
			}

			// Get information about the card
			cardItem, err := s.serviceCard.Item(card.GameID, card.CollectionID, deckInfo.ID, card.ID)
			if err != nil {
				return err
			}

			cardGUID := fmt.Sprintf("%06d", commonIndex)
			commonIndex++
			cardObject := tts_entity.NewCard(
				cardGUID,
				cardItem.Name,
				cardItem.Description,
				page.GetIndex()+deckIdOffset,
				page.Size()-1,
				cardItem.Variables,
				deckDescription,
				tts_entity.Transform{
					ScaleX: cfg.CardSize.ScaleX,
					ScaleY: cfg.CardSize.ScaleY,
					ScaleZ: cfg.CardSize.ScaleZ,
				},
			)
			for i := 0; i < cardItem.Count; i++ {
				// Add a card to the deck as many times as set in the count variable
				deck.AddCard(cardObject)
			}
		}

		if !page.IsEmpty() {
			if _, ok := collectionBags[prevCollection]; !ok {
				collectionBags[prevCollection] = utils.Allocate(tts_entity.NewBag(prevCollection))
			}
			switch {
			case len(deck.ContainedObjects) == 1:
				// We cannot create a deck object with a single card. We must create a card object.
				// bag.ContainedObjects = append(bag.ContainedObjects, deck.ContainedObjects[0])
				collectionBags[prevCollection].ContainedObjects = append(collectionBags[prevCollection].ContainedObjects, deck.ContainedObjects[0])
			case len(deck.ContainedObjects) > 1:
				// If there is more than one card in the deck, place the deck in the object list.
				// bag.ContainedObjects = append(bag.ContainedObjects, deck)
				collectionBags[prevCollection].ContainedObjects = append(collectionBags[prevCollection].ContainedObjects, deck)
			}
		}

		deckIdOffset += page.GetIndex()
	}

	// Add all collection bags into game bag
	for _, collectionBag := range collectionBags {
		bag.ContainedObjects = append(bag.ContainedObjects, collectionBag)
	}

	bag.Description = fmt.Sprintf("Created at: %v", time.Now().Format("2006-01-02 15:04:05"))
	root := tts_entity.RootObjects{
		ObjectStates: []tts_entity.Bag{
			bag,
		},
	}
	err := fs.CreateAndProcess(filepath.Join(s.cfg.Results(), gameItem.ID+".json"), root, fs.JsonToWriter[tts_entity.RootObjects])
	if err != nil {
		return err
	}

	// Try to upload to TTS if it's possible
	s.serviceTTS.SendToTTS(bag)

	return nil
}
