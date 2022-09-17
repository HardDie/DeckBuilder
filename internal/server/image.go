package server

import (
	"net/http"

	"github.com/gorilla/mux"

	"tts_deck_build/internal/errors"
	"tts_deck_build/internal/fs"
	"tts_deck_build/internal/network"
	"tts_deck_build/internal/service"
)

type ImageServer struct {
	gameService       service.IGameService
	collectionService service.ICollectionService
	deckService       service.IDeckService
	cardService       service.ICardService
}

func NewImageServer(gameService service.IGameService, collectionService service.ICollectionService, deckService service.IDeckService, cardService service.ICardService) *ImageServer {
	return &ImageServer{
		gameService:       gameService,
		collectionService: collectionService,
		deckService:       deckService,
		cardService:       cardService,
	}
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
	img, imgType, e := s.cardService.GetImage(gameID, collectionID, deckID, cardID)
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
	img, imgType, e := s.collectionService.GetImage(gameID, collectionID)
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
	img, imgType, e := s.deckService.GetImage(gameID, collectionID, deckID)
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
	img, imgType, e := s.gameService.GetImage(gameID)
	if e != nil {
		network.ResponseError(w, e)
		return
	}
	w.Header().Set("Content-Type", "image/"+imgType)
	if _, err := w.Write(img); err != nil {
		errors.IfErrorLog(err)
	}
}
