package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"image/jpeg"
	"net"
	"path/filepath"
	"strconv"
	"time"

	"github.com/HardDie/DeckBuilder/internal/config"
	"github.com/HardDie/DeckBuilder/internal/dto"
	"github.com/HardDie/DeckBuilder/internal/entity"
	"github.com/HardDie/DeckBuilder/internal/fs"
	"github.com/HardDie/DeckBuilder/internal/images"
	"github.com/HardDie/DeckBuilder/internal/logger"
	pageDrawer "github.com/HardDie/DeckBuilder/internal/page_drawer"
	"github.com/HardDie/DeckBuilder/internal/progress"
	"github.com/HardDie/DeckBuilder/internal/tts_entity"
)

type IGeneratorService interface {
	GenerateGame(gameID string, dtoObject *dto.GenerateGameDTO) error
	DataForTTS() ([]byte, error)
}
type GeneratorService struct {
	cfg               *config.Config
	gameService       IGameService
	collectionService ICollectionService
	deckService       IDeckService
	cardService       ICardService

	dataForTTS []byte
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
		err = s.generateBody(gameItem, deckArray, dtoObject.Scale)
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

func (s *GeneratorService) getListOfCards(gameID string, sortField string) (map[Deck][]Card, error) {
	decks := make(map[Deck][]Card)
	// Get list of collections
	collectionItems, _, err := s.collectionService.List(gameID, sortField, "")
	if err != nil {
		return nil, err
	}
	// Iterate through collections
	for _, collectionItem := range collectionItems {
		// Get list of decks
		deckItems, _, err := s.deckService.List(gameID, collectionItem.ID, sortField, "")
		if err != nil {
			return nil, err
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
			cardItems, _, err := s.cardService.List(gameID, collectionItem.ID, deckItem.ID, sortField, "")
			if err != nil {
				return nil, err
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
	return decks, nil
}

func (s *GeneratorService) generateBody(gameItem *entity.GameInfo, decks map[Deck][]Card, scale int) error {
	pr := progress.GetProgress()
	pr.SetMessage("Reading a list of cards from the disk...")

	// Generate images
	imageMapping, err := s.generateImages(decks, scale)
	if err != nil {
		return err
	}
	// Generate json description
	err = s.generateJson(gameItem, decks, imageMapping)
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

func (s *GeneratorService) generateImages(decks map[Deck][]Card, scale int) (map[string]PageInfo, error) {
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
	for deckInfo, cards := range decks {
		// Create page drawer object
		page := pageDrawer.New(deckInfo.ID, s.cfg.Results(), scale)
		var backsidePath string

		// Iterate through all cards in deck
		for _, card := range cards {
			// Init page drawer with deck information
			if page.IsEmpty() {
				// Getting an image of the backside
				deckBacksideImage, _, err := s.deckService.GetImage(card.GameID, card.CollectionID, deckInfo.ID)
				if err != nil {
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
			}

			// Get card image
			cardImageBin, _, err := s.cardService.GetImage(card.GameID, card.CollectionID, deckInfo.ID, card.ID)
			if err != nil {
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

func (s *GeneratorService) generateJson(gameItem *entity.GameInfo, decks map[Deck][]Card, imageMapping map[string]PageInfo) error {
	bag := tts_entity.NewBag(gameItem.Name)
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

	for deckInfo, cards := range decks {
		deck = tts_entity.NewDeck(deckInfo.Name)
		// Create page drawer object
		page := pageDrawer.New(deckInfo.ID, "", 1)

		pageInfo := imageMapping[deckInfo.ID+"_"+strconv.Itoa(page.GetIndex())]
		deckDescription := tts_entity.DeckDescription{
			FaceURL:   "file:///" + pageInfo.Image,
			BackURL:   "file:///" + pageInfo.Backside,
			NumWidth:  pageInfo.Columns,
			NumHeight: pageInfo.Rows,
		}
		deck.CustomDeck[page.GetIndex()+deckIdOffset] = deckDescription

		var prevCollection string

		// Iterate through all cards in deck
		for _, card := range cards {
			if page.IsEmpty() {
				prevCollection = card.CollectionID + deckInfo.ID
			}

			// Start new page if current is full
			if page.IsFull() {
				page = (&pageDrawer.PageDrawer{}).Inherit(page)

				pageInfo = imageMapping[deckInfo.ID+"_"+strconv.Itoa(page.GetIndex())]
				deckDescription = tts_entity.DeckDescription{
					FaceURL:   "file:///" + pageInfo.Image,
					BackURL:   "file:///" + pageInfo.Backside,
					NumWidth:  pageInfo.Columns,
					NumHeight: pageInfo.Rows,
				}
				deck.CustomDeck[page.GetIndex()+deckIdOffset] = deckDescription
			}

			if card.CollectionID+deckInfo.ID != prevCollection {
				prevCollection = card.CollectionID + deckInfo.ID

				switch {
				case len(deck.ContainedObjects) == 1:
					// We cannot create a deck object with a single card. We must create a card object.
					bag.ContainedObjects = append(bag.ContainedObjects, deck.ContainedObjects[0])
				case len(deck.ContainedObjects) > 1:
					// If there is more than one card in the deck, place the deck in the object list.
					bag.ContainedObjects = append(bag.ContainedObjects, deck)
				}
				deck = tts_entity.NewDeck(deckInfo.Name)
				deck.CustomDeck[page.GetIndex()+deckIdOffset] = deckDescription
			}

			// Add card on page
			err := page.AddImage(dummyImage)
			if err != nil {
				return err
			}

			// Get information about the card
			cardItem, err := s.cardService.Item(card.GameID, card.CollectionID, deckInfo.ID, card.ID)
			if err != nil {
				return err
			}

			cardObject := tts_entity.NewCard(cardItem.Name, cardItem.Description, page.GetIndex()+deckIdOffset, page.Size()-1, cardItem.Variables, deckDescription)
			for i := 0; i < cardItem.Count; i++ {
				// Add a card to the deck as many times as set in the count variable
				deck.AddCard(cardObject)
			}
		}

		if !page.IsEmpty() {
			switch {
			case len(deck.ContainedObjects) == 1:
				// We cannot create a deck object with a single card. We must create a card object.
				bag.ContainedObjects = append(bag.ContainedObjects, deck.ContainedObjects[0])
			case len(deck.ContainedObjects) > 1:
				// If there is more than one card in the deck, place the deck in the object list.
				bag.ContainedObjects = append(bag.ContainedObjects, deck)
			}
		}

		deckIdOffset += page.GetIndex()
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
	conn, err := net.Dial("tcp", "localhost:39999")
	if err != nil {
		logger.Info.Println("Can't connect to TTS via tcp connection 'localhost:39999':", err.Error())
		return nil
	}
	defer func() { conn.Close() }()

	dataForTTS, err := json.Marshal(bag)
	if err != nil {
		logger.Warn.Println("error marshal data for TTS:", err.Error())
		return nil
	}
	s.dataForTTS = dataForTTS

	type Message struct {
		MessageID int    `json:"messageID"`
		GUID      string `json:"guid"`
		Script    string `json:"script"`
	}

	msg := Message{
		MessageID: 3,
		GUID:      "-1",
		Script: `
WebRequest.get("http://localhost:5000/api/generator/data", function(request)
	if request.is_error then
		print('Downloading json error: ', request.error)
		return
	end
	print('JSON were downloaded!')
	local object = spawnObjectJSON({
		json = request.text,
		callback_function = function(spawned_object)
			print('Object were spawned! Done!')
		end
	})
end)`,
	}

	data, err := json.Marshal(msg)
	if err != nil {
		logger.Warn.Println("error marshal msg for TTS:", err.Error())
		return nil
	}

	_, err = conn.Write(data)
	if err != nil {
		logger.Warn.Println("error write message into TTS socket:", err.Error())
		return nil
	}

	return nil
}

func (s *GeneratorService) DataForTTS() ([]byte, error) {
	if s.dataForTTS == nil {
		return nil, errors.New("there is nothing to serve")
	}
	res := s.dataForTTS
	s.dataForTTS = nil
	return res, nil
}
