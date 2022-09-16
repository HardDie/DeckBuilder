package server

import (
	"net/http"

	"github.com/gorilla/mux"

	"tts_deck_build/internal/cards"
	"tts_deck_build/internal/collections"
	"tts_deck_build/internal/decks"
	"tts_deck_build/internal/errors"
	"tts_deck_build/internal/fs"
	"tts_deck_build/internal/games"
	"tts_deck_build/internal/network"
)

type ImageServer struct {
}

func NewImageServer() *ImageServer {
	return &ImageServer{}
}

func (s *ImageServer) CardHandler(w http.ResponseWriter, r *http.Request) {
	gameID := mux.Vars(r)["game"]
	collectionID := mux.Vars(r)["collection"]
	deckID := mux.Vars(r)["deck"]
	cardID, e := fs.StringToInt64(mux.Vars(r)["card"])
	if e != nil {
		network.ResponseError(w, e)
		return
	}
	img, imgType, e := cards.NewService().GetImage(gameID, collectionID, deckID, cardID)
	if e != nil {
		network.ResponseError(w, e)
		return
	}
	w.Header().Set("Content-Type", "image/"+imgType)
	if _, err := w.Write(img); err != nil {
		errors.IfErrorLog(err)
	}
}
func (s *ImageServer) CollectionHandler(w http.ResponseWriter, r *http.Request) {
	gameID := mux.Vars(r)["game"]
	collectionID := mux.Vars(r)["collection"]
	img, imgType, e := collections.NewService().GetImage(gameID, collectionID)
	if e != nil {
		network.ResponseError(w, e)
		return
	}
	w.Header().Set("Content-Type", "image/"+imgType)
	if _, err := w.Write(img); err != nil {
		errors.IfErrorLog(err)
	}
}
func (s *ImageServer) DeckHandler(w http.ResponseWriter, r *http.Request) {
	gameID := mux.Vars(r)["game"]
	collectionID := mux.Vars(r)["collection"]
	deckID := mux.Vars(r)["deck"]
	img, imgType, e := decks.NewService().GetImage(gameID, collectionID, deckID)
	if e != nil {
		network.ResponseError(w, e)
		return
	}
	w.Header().Set("Content-Type", "image/"+imgType)
	if _, err := w.Write(img); err != nil {
		errors.IfErrorLog(err)
	}
}
func (s *ImageServer) GameHandler(w http.ResponseWriter, r *http.Request) {
	gameID := mux.Vars(r)["game"]
	img, imgType, e := games.NewService().GetImage(gameID)
	if e != nil {
		network.ResponseError(w, e)
		return
	}
	w.Header().Set("Content-Type", "image/"+imgType)
	if _, err := w.Write(img); err != nil {
		errors.IfErrorLog(err)
	}
}
