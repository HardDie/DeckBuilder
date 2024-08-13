package image

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/HardDie/DeckBuilder/internal/errors"
	"github.com/HardDie/DeckBuilder/internal/fs"
	"github.com/HardDie/DeckBuilder/internal/network"
	servicesCard "github.com/HardDie/DeckBuilder/internal/services/card"
	servicesCollection "github.com/HardDie/DeckBuilder/internal/services/collection"
	servicesDeck "github.com/HardDie/DeckBuilder/internal/services/deck"
	servicesGame "github.com/HardDie/DeckBuilder/internal/services/game"
)

type image struct {
	serviceGame       servicesGame.Game
	serviceCollection servicesCollection.Collection
	serviceDeck       servicesDeck.Deck
	serviceCard       servicesCard.Card
}

func New(
	serviceGame servicesGame.Game,
	serviceCollection servicesCollection.Collection,
	serviceDeck servicesDeck.Deck,
	serviceCard servicesCard.Card,
) Image {
	return &image{
		serviceGame:       serviceGame,
		serviceCollection: serviceCollection,
		serviceDeck:       serviceDeck,
		serviceCard:       serviceCard,
	}
}

func (s *image) CardHandler(w http.ResponseWriter, r *http.Request) {
	gameID := mux.Vars(r)["game"]
	collectionID := mux.Vars(r)["collection"]
	deckID := mux.Vars(r)["deck"]
	cardID, e := fs.StringToInt64(mux.Vars(r)["card"])
	if e != nil {
		network.ResponseError(w, e)
		return
	}
	img, imgType, e := s.serviceCard.GetImage(gameID, collectionID, deckID, cardID)
	if e != nil {
		network.ResponseError(w, e)
		return
	}
	w.Header().Set("Content-Type", "image/"+imgType)
	if _, err := w.Write(img); err != nil {
		errors.IfErrorLog(err)
	}
}
func (s *image) CollectionHandler(w http.ResponseWriter, r *http.Request) {
	gameID := mux.Vars(r)["game"]
	collectionID := mux.Vars(r)["collection"]
	img, imgType, e := s.serviceCollection.GetImage(gameID, collectionID)
	if e != nil {
		network.ResponseError(w, e)
		return
	}
	w.Header().Set("Content-Type", "image/"+imgType)
	if _, err := w.Write(img); err != nil {
		errors.IfErrorLog(err)
	}
}
func (s *image) DeckHandler(w http.ResponseWriter, r *http.Request) {
	gameID := mux.Vars(r)["game"]
	collectionID := mux.Vars(r)["collection"]
	deckID := mux.Vars(r)["deck"]
	img, imgType, e := s.serviceDeck.GetImage(gameID, collectionID, deckID)
	if e != nil {
		network.ResponseError(w, e)
		return
	}
	w.Header().Set("Content-Type", "image/"+imgType)
	if _, err := w.Write(img); err != nil {
		errors.IfErrorLog(err)
	}
}
func (s *image) GameHandler(w http.ResponseWriter, r *http.Request) {
	gameID := mux.Vars(r)["game"]
	img, imgType, e := s.serviceGame.GetImage(gameID)
	if e != nil {
		network.ResponseError(w, e)
		return
	}
	w.Header().Set("Content-Type", "image/"+imgType)
	if _, err := w.Write(img); err != nil {
		errors.IfErrorLog(err)
	}
}
