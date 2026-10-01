package servers

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/gorilla/mux"

	er "github.com/HardDie/DeckBuilder/internal/errors"
)

type imageStub struct {
	img  []byte
	kind string
	err  error
	got  []string
}

func (s *imageStub) GetImage(ids ...string) ([]byte, string, error) {
	s.got = ids
	return s.img, s.kind, s.err
}

type gameStub struct{ imageStub }

func (s *gameStub) GetImage(gameID string) ([]byte, string, error) {
	return s.imageStub.GetImage(gameID)
}

type collectionStub struct{ imageStub }

func (s *collectionStub) GetImage(gameID, collectionID string) ([]byte, string, error) {
	return s.imageStub.GetImage(gameID, collectionID)
}

type deckStub struct{ imageStub }

func (s *deckStub) GetImage(gameID, collectionID, deckID string) ([]byte, string, error) {
	return s.imageStub.GetImage(gameID, collectionID, deckID)
}

type cardStub struct{ imageStub }

func (s *cardStub) GetImage(gameID, collectionID, deckID string, cardID int64) ([]byte, string, error) {
	s.got = []string{gameID, collectionID, deckID, strconv.FormatInt(cardID, 10)}
	return s.img, s.kind, s.err
}

type ttsStub struct {
	data []byte
	err  error
}

func (s *ttsStub) DataForTTS() ([]byte, error) {
	if s.err != nil {
		return nil, s.err
	}
	if s.data == nil {
		return nil, errors.New("there is nothing to serve")
	}
	res := s.data
	s.data = nil
	return res, nil
}

func TestRoutes(t *testing.T) {
	game := &gameStub{imageStub{img: []byte("png-bytes"), kind: "png"}}
	collection := &collectionStub{imageStub{img: []byte("jpeg-bytes"), kind: "jpeg"}}
	deck := &deckStub{imageStub{err: er.DeckImageNotExists}}
	card := &cardStub{imageStub{img: []byte("gif-bytes"), kind: "gif"}}
	tts := &ttsStub{data: []byte(`{"Name":"Bag"}`)}

	route := mux.NewRouter()
	Register(route, game, collection, deck, card, tts)

	cases := []struct {
		path        string
		status      int
		contentType string
		body        string
	}{
		{"/api/games/four_souls/image", http.StatusOK, "image/png", "png-bytes"},
		{"/api/games/four_souls/collections/base/image", http.StatusOK, "image/jpeg", "jpeg-bytes"},
		{"/api/games/four_souls/collections/base/decks/loot/image", http.StatusBadRequest, "application/json; charset=utf-8", ""},
		{"/api/games/four_souls/collections/base/decks/loot/cards/7/image", http.StatusOK, "image/gif", "gif-bytes"},
		{"/api/games/four_souls/collections/base/decks/loot/cards/nope/image", http.StatusBadRequest, "application/json; charset=utf-8", ""},
		{"/api/tts/data", http.StatusOK, "application/json; charset=utf-8", `{"Name":"Bag"}`},
	}
	for _, tc := range cases {
		rec := httptest.NewRecorder()
		route.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, tc.path, nil))
		if rec.Code != tc.status {
			t.Fatalf("%s status %d, want %d, body %s", tc.path, rec.Code, tc.status, rec.Body.String())
		}
		if got := rec.Header().Get("Content-Type"); got != tc.contentType {
			t.Fatalf("%s content type %q, want %q", tc.path, got, tc.contentType)
		}
		if tc.body != "" && rec.Body.String() != tc.body {
			t.Fatalf("%s body %q, want %q", tc.path, rec.Body.String(), tc.body)
		}
	}

	if len(game.got) != 1 || game.got[0] != "four_souls" {
		t.Fatalf("game ids %#v", game.got)
	}
	if len(collection.got) != 2 || collection.got[0] != "four_souls" || collection.got[1] != "base" {
		t.Fatalf("collection ids %#v", collection.got)
	}
	if len(card.got) != 4 || card.got[3] != "7" {
		t.Fatalf("card ids %#v", card.got)
	}

	rec := httptest.NewRecorder()
	route.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/api/tts/data", nil))
	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("second tts status %d", rec.Code)
	}
}

func TestTTSEmpty(t *testing.T) {
	route := mux.NewRouter()
	Register(route, &gameStub{}, &collectionStub{}, &deckStub{}, &cardStub{}, &ttsStub{err: errors.New("there is nothing to serve")})
	rec := httptest.NewRecorder()
	route.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/api/tts/data", nil))
	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("status %d", rec.Code)
	}
}
