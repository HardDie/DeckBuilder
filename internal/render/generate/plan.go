// Plan a game render: which sheets to paint, where, and the TTS object.
package generate

import (
	"errors"
	"github.com/HardDie/DeckBuilder/internal/apperr"
	entitiesGame "github.com/HardDie/DeckBuilder/internal/entities/game"
	entitiesSettings "github.com/HardDie/DeckBuilder/internal/entities/settings"
	"github.com/HardDie/DeckBuilder/internal/render/generate/catalog"
	"github.com/HardDie/DeckBuilder/internal/render/generate/layout"
	"github.com/HardDie/DeckBuilder/internal/render/generate/script"
	servicesCard "github.com/HardDie/DeckBuilder/internal/services/card"
	servicesDeck "github.com/HardDie/DeckBuilder/internal/services/deck"
	"github.com/HardDie/DeckBuilder/internal/tts_entity"
	"log/slog"
)

// Sheet is one page ready to draw. Faces and back are the original file bytes.
// DeckName names the deck for messages; the images carry no card names.
type Sheet struct {
	DeckName string
	Faces    [][]byte
	Back     []byte
	CellW    int
	CellH    int
	Shadow   bool
	Path     string
}

// File is a non-image result: a raw back PNG.
type File struct {
	Path string
	Body []byte
}

// Plan is everything the drawer and the JSON writer need.
type Plan struct {
	Sheets   []Sheet
	Backs    []File
	JSONPath string
	Root     tts_entity.RootObjects
	Bag      tts_entity.Bag
}

// Prepare loads face and back bytes and measures every page.
// It does not decode, resize, draw, or encode.
func Prepare(
	dir string,
	gameItem *entitiesGame.Game,
	decks map[catalog.Deck][]catalog.Card,
	order []catalog.Deck,
	scale int,
	cfg *entitiesSettings.Settings,
	deckSvc servicesDeck.Deck,
	cardSvc servicesCard.Card,
) (Plan, error) {
	measured := make([]layout.Deck, 0, len(order))
	deckNames := make(map[string]string, len(order))
	for _, deckInfo := range order {
		cards := decks[deckInfo]
		if len(cards) == 0 {
			continue
		}
		back, _, err := deckSvc.GetImage(cards[0].GameID, cards[0].CollectionID, deckInfo.ID)
		if err != nil {
			slog.Error("deck back image not found", "game", cards[0].GameID, "collection", cards[0].CollectionID, "deck", deckInfo.ID, "err", err)
			return Plan{}, err
		}
		faces := make([][]byte, 0, len(cards))
		for _, card := range cards {
			face, _, err := cardSvc.GetImage(card.GameID, card.CollectionID, deckInfo.ID, card.ID)
			if err != nil {
				slog.Error("card image not found", "game", card.GameID, "collection", card.CollectionID, "deck", deckInfo.ID, "card", card.ID, "err", err)
				return Plan{}, err
			}
			faces = append(faces, face)
		}
		measured = append(measured, layout.Deck{ID: deckInfo.ID, Back: back, Faces: faces})
		deckNames[deckInfo.ID] = deckInfo.Name
	}

	pages, backs, err := layout.Pages(dir, measured, scale, cfg.EnableBackShadow)
	var imageErr *layout.ImageError
	if errors.As(err, &imageErr) {
		slog.Warn("deck image size unreadable", "deck", deckNames[imageErr.DeckID], "err", err)
		return Plan{}, UnreadableImage(deckNames[imageErr.DeckID])
	}
	if err != nil {
		return Plan{}, err
	}
	scriptPages := make([]script.Page, len(pages))
	sheets := make([]Sheet, len(pages))
	for i, page := range pages {
		scriptPages[i] = script.Page{
			DeckID: page.DeckID,
			Index:  page.Index,
			Image:  page.SheetPath,
			Back:   page.BackPath,
			Cols:   page.Cols,
			Rows:   page.Rows,
		}
		sheets[i] = Sheet{
			DeckName: deckNames[page.DeckID],
			Faces:    page.Faces,
			Back:     page.Back,
			CellW:    page.CellW,
			CellH:    page.CellH,
			Shadow:   page.Shadow,
			Path:     page.SheetPath,
		}
	}
	doc, err := script.Build(dir, gameItem, decks, order, scriptPages, cfg, cardSvc)
	if err != nil {
		return Plan{}, err
	}
	rawBacks := make([]File, len(backs))
	for i, back := range backs {
		rawBacks[i] = File{Path: back.Path, Body: back.Body}
	}
	return Plan{
		Sheets:   sheets,
		Backs:    rawBacks,
		JSONPath: doc.Path,
		Root:     doc.Root,
		Bag:      doc.Bag,
	}, nil
}

// UnreadableImage is what the window shows when an image in a deck cannot be read.
// The deck is known; which card is not, since sheets hold images, not card names.
func UnreadableImage(deckName string) error {
	return apperr.Withf(apperr.ErrUnsupportedImage,
		"an image in deck %q could not be read. The file may be damaged", deckName)
}
