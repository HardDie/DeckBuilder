// In-memory catalog for generator tests.
package fake

import (
	"errors"
	"sync"

	entitiesCard "github.com/HardDie/DeckBuilder/internal/entities/card"
	entitiesCollection "github.com/HardDie/DeckBuilder/internal/entities/collection"
	entitiesDeck "github.com/HardDie/DeckBuilder/internal/entities/deck"
	entitiesGame "github.com/HardDie/DeckBuilder/internal/entities/game"
	entitiesSettings "github.com/HardDie/DeckBuilder/internal/entities/settings"
	servicesCard "github.com/HardDie/DeckBuilder/internal/services/card"
	servicesCollection "github.com/HardDie/DeckBuilder/internal/services/collection"
	servicesDeck "github.com/HardDie/DeckBuilder/internal/services/deck"
	servicesGame "github.com/HardDie/DeckBuilder/internal/services/game"
	servicesSystem "github.com/HardDie/DeckBuilder/internal/services/system"
	servicesTTS "github.com/HardDie/DeckBuilder/internal/services/tts"
)

// World is one game. List order is the order stored here.
type World struct {
	GameID       string
	GameName     string
	Settings     entitiesSettings.Settings
	Collections  []*Coll
	MissingBacks map[string]bool
	MissingFaces map[int64]bool

	mu    sync.Mutex
	sorts []string
	tts   []any
}

// Coll is one collection and its decks.
type Coll struct {
	ID    string
	Name  string
	Decks []*Deck
}

// Deck is one deck, its back bytes, and its cards.
type Deck struct {
	ID    string
	Name  string
	Image string
	Back  []byte
	Cards []*Card
}

// Card is one catalog card and its face bytes.
type Card struct {
	ID          int64
	Name        string
	Description string
	Count       int
	Variables   map[string]string
	Face        []byte
}

// Games is the game service.
type Games struct{ *World }

// Collections is the collection service.
type Collections struct{ *World }

// Decks is the deck service.
type Decks struct{ *World }

// Cards is the card service.
type Cards struct{ *World }

// Systems is the settings service.
type Systems struct{ *World }

// Speech is the TTS service.
type Speech struct{ *World }

func (w *World) noteSort(field string) {
	w.mu.Lock()
	w.sorts = append(w.sorts, field)
	w.mu.Unlock()
}

// Sorts is every sort field passed to List, in call order.
func (w *World) Sorts() []string {
	w.mu.Lock()
	defer w.mu.Unlock()
	out := make([]string, len(w.sorts))
	copy(out, w.sorts)
	return out
}

// TTS is every value passed to SendToTTS.
func (w *World) TTS() []any {
	w.mu.Lock()
	defer w.mu.Unlock()
	out := make([]any, len(w.tts))
	copy(out, w.tts)
	return out
}

func (g Games) Item(gameID string) (*entitiesGame.Game, error) {
	if g.GameID != gameID {
		return nil, errors.New("game not found")
	}
	return &entitiesGame.Game{ID: g.GameID, Name: g.GameName}, nil
}
func (Games) Create(servicesGame.CreateRequest) (*entitiesGame.Game, error) { return nil, nil }
func (Games) List(string, string) ([]*entitiesGame.Game, error)             { return nil, nil }
func (Games) Update(string, servicesGame.UpdateRequest) (*entitiesGame.Game, error) {
	return nil, nil
}
func (Games) Delete(string) error                     { return nil }
func (Games) GetImage(string) ([]byte, string, error) { return nil, "", nil }
func (Games) Duplicate(string, servicesGame.DuplicateRequest) (*entitiesGame.Game, error) {
	return nil, nil
}
func (Games) Export(string) ([]byte, error)                     { return nil, nil }
func (Games) Import([]byte, string) (*entitiesGame.Game, error) { return nil, nil }

func (Collections) Create(string, servicesCollection.CreateRequest) (*entitiesCollection.Collection, error) {
	return nil, nil
}
func (Collections) Item(string, string) (*entitiesCollection.Collection, error) {
	return nil, nil
}
func (c Collections) List(gameID, sortField, _ string) ([]*entitiesCollection.Collection, error) {
	c.noteSort(sortField)
	out := make([]*entitiesCollection.Collection, len(c.Collections))
	for i, col := range c.Collections {
		out[i] = &entitiesCollection.Collection{ID: col.ID, Name: col.Name, GameID: gameID}
	}
	return out, nil
}
func (Collections) Update(string, string, servicesCollection.UpdateRequest) (*entitiesCollection.Collection, error) {
	return nil, nil
}
func (Collections) Delete(string, string) error                     { return nil }
func (Collections) GetImage(string, string) ([]byte, string, error) { return nil, "", nil }

func (Decks) Create(string, string, servicesDeck.CreateRequest) (*entitiesDeck.Deck, error) {
	return nil, nil
}
func (d Decks) Item(gameID, collectionID, deckID string) (*entitiesDeck.Deck, error) {
	deck := d.findDeck(collectionID, deckID)
	if deck == nil {
		return nil, errors.New("deck not found")
	}
	return &entitiesDeck.Deck{
		ID:           deck.ID,
		Name:         deck.Name,
		Image:        deck.Image,
		GameID:       gameID,
		CollectionID: collectionID,
		HasImage:     !d.MissingBacks[deckID],
	}, nil
}
func (d Decks) List(gameID, collectionID, sortField, _ string) ([]*entitiesDeck.Deck, error) {
	d.noteSort(sortField)
	for _, col := range d.Collections {
		if col.ID != collectionID {
			continue
		}
		out := make([]*entitiesDeck.Deck, len(col.Decks))
		for i, deck := range col.Decks {
			out[i] = &entitiesDeck.Deck{
				ID:           deck.ID,
				Name:         deck.Name,
				Image:        deck.Image,
				GameID:       gameID,
				CollectionID: collectionID,
				HasImage:     !d.MissingBacks[deck.ID],
			}
		}
		return out, nil
	}
	return nil, nil
}
func (Decks) Update(string, string, string, servicesDeck.UpdateRequest) (*entitiesDeck.Deck, error) {
	return nil, nil
}
func (Decks) Delete(string, string, string) error { return nil }
func (d Decks) GetImage(_, collectionID, deckID string) ([]byte, string, error) {
	if d.MissingBacks[deckID] {
		return nil, "", errors.New("back missing")
	}
	deck := d.findDeck(collectionID, deckID)
	if deck == nil {
		return nil, "", errors.New("back missing")
	}
	return deck.Back, "image/png", nil
}
func (Decks) ListAllUnique(string) ([]*entitiesDeck.Deck, error) { return nil, nil }

func (Cards) Create(string, string, string, servicesCard.CreateRequest) (*entitiesCard.Card, error) {
	return nil, nil
}
func (c Cards) Item(gameID, collectionID, deckID string, cardID int64) (*entitiesCard.Card, error) {
	card := c.findCard(collectionID, deckID, cardID)
	if card == nil {
		return nil, errors.New("card missing")
	}
	return cardEntity(gameID, collectionID, deckID, card, !c.MissingFaces[card.ID]), nil
}
func (c Cards) List(gameID, collectionID, deckID, sortField, _ string) ([]*entitiesCard.Card, error) {
	c.noteSort(sortField)
	deck := c.findDeck(collectionID, deckID)
	if deck == nil {
		return nil, nil
	}
	out := make([]*entitiesCard.Card, len(deck.Cards))
	for i, card := range deck.Cards {
		out[i] = cardEntity(gameID, collectionID, deckID, card, !c.MissingFaces[card.ID])
	}
	return out, nil
}
func (Cards) Update(string, string, string, int64, servicesCard.UpdateRequest) (*entitiesCard.Card, error) {
	return nil, nil
}
func (Cards) Delete(string, string, string, int64) error { return nil }
func (c Cards) GetImage(_, collectionID, deckID string, cardID int64) ([]byte, string, error) {
	if c.MissingFaces[cardID] {
		return nil, "", errors.New("face missing")
	}
	card := c.findCard(collectionID, deckID, cardID)
	if card == nil {
		return nil, "", errors.New("face missing")
	}
	return card.Face, "image/png", nil
}

func (s Systems) GetSettings() (*entitiesSettings.Settings, error) {
	cfg := s.Settings
	return &cfg, nil
}
func (Systems) UpdateSettings(servicesSystem.UpdateSettingsRequest) (*entitiesSettings.Settings, error) {
	return nil, nil
}

func (Speech) SetHTTPPort(int) {}
func (s Speech) SendToTTS(data any) {
	s.mu.Lock()
	s.tts = append(s.tts, data)
	s.mu.Unlock()
}
func (Speech) DataForTTS() ([]byte, error) { return nil, nil }

func cardEntity(gameID, collectionID, deckID string, card *Card, hasImage bool) *entitiesCard.Card {
	return &entitiesCard.Card{
		ID:           card.ID,
		Name:         card.Name,
		Description:  card.Description,
		Variables:    card.Variables,
		Count:        card.Count,
		GameID:       gameID,
		CollectionID: collectionID,
		DeckID:       deckID,
		HasImage:     hasImage,
	}
}

var (
	_ servicesGame.Game             = Games{}
	_ servicesCollection.Collection = Collections{}
	_ servicesDeck.Deck             = Decks{}
	_ servicesCard.Card             = Cards{}
	_ servicesSystem.System         = Systems{}
	_ servicesTTS.TTS               = Speech{}
)

func (w *World) findDeck(collectionID, deckID string) *Deck {
	for _, col := range w.Collections {
		if col.ID != collectionID {
			continue
		}
		for _, deck := range col.Decks {
			if deck.ID == deckID {
				return deck
			}
		}
	}
	return nil
}

func (w *World) findCard(collectionID, deckID string, cardID int64) *Card {
	deck := w.findDeck(collectionID, deckID)
	if deck == nil {
		return nil
	}
	for _, card := range deck.Cards {
		if card.ID == cardID {
			return card
		}
	}
	return nil
}
