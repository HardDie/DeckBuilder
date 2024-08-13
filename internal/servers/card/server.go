package card

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/HardDie/DeckBuilder/internal/dto"
	er "github.com/HardDie/DeckBuilder/internal/errors"
	"github.com/HardDie/DeckBuilder/internal/fs"
	"github.com/HardDie/DeckBuilder/internal/network"
	serversSystem "github.com/HardDie/DeckBuilder/internal/servers/system"
	servicesCard "github.com/HardDie/DeckBuilder/internal/services/card"
	"github.com/HardDie/DeckBuilder/internal/utils"
)

type card struct {
	serviceCard  servicesCard.Card
	serverSystem serversSystem.System
}

func New(serviceCard servicesCard.Card, serverSystem serversSystem.System) Card {
	return &card{
		serviceCard:  serviceCard,
		serverSystem: serverSystem,
	}
}

func (s *card) CreateHandler(w http.ResponseWriter, r *http.Request) {
	gameID := mux.Vars(r)["game"]
	collectionID := mux.Vars(r)["collection"]
	deckID := mux.Vars(r)["deck"]

	e := r.ParseMultipartForm(0)
	if e != nil {
		er.IfErrorLog(e)
		e = er.InternalError.AddMessage(e.Error())
		network.ResponseError(w, e)
		return
	}

	data, e := utils.GetFileFromMultipart("imageFile", r)
	if e != nil {
		network.ResponseError(w, e)
		return
	}

	dtoObject := &dto.CreateCardDTO{
		Name:        r.FormValue("name"),
		Description: r.FormValue("description"),
		Image:       r.FormValue("image"),
		Variables:   nil,
		Count:       fs.StringToInt(r.FormValue("count")),
		ImageFile:   data,
	}

	variablesJson := r.FormValue("variables")
	if variablesJson != "" {
		e = json.Unmarshal([]byte(variablesJson), &dtoObject.Variables)
		if e != nil {
			er.IfErrorLog(e)
			e = er.InternalError.HTTP(http.StatusBadRequest).AddMessage("Bad variables json")
			network.ResponseError(w, e)
			return
		}
	}

	item, e := s.serviceCard.Create(gameID, collectionID, deckID, dtoObject)
	if e != nil {
		network.ResponseError(w, e)
		return
	}
	network.Response(w, item)
}
func (s *card) DeleteHandler(w http.ResponseWriter, r *http.Request) {
	gameID := mux.Vars(r)["game"]
	collectionID := mux.Vars(r)["collection"]
	deckID := mux.Vars(r)["deck"]
	cardID, e := fs.StringToInt64(mux.Vars(r)["card"])
	if e != nil {
		network.ResponseError(w, e)
		return
	}
	e = s.serviceCard.Delete(gameID, collectionID, deckID, cardID)
	if e != nil {
		network.ResponseError(w, e)
	}
}
func (s *card) ItemHandler(w http.ResponseWriter, r *http.Request) {
	gameID := mux.Vars(r)["game"]
	collectionID := mux.Vars(r)["collection"]
	deckID := mux.Vars(r)["deck"]
	cardID, e := fs.StringToInt64(mux.Vars(r)["card"])
	if e != nil {
		network.ResponseError(w, e)
		return
	}
	item, e := s.serviceCard.Item(gameID, collectionID, deckID, cardID)
	if e != nil {
		network.ResponseError(w, e)
		return
	}
	network.Response(w, item)
}
func (s *card) ListHandler(w http.ResponseWriter, r *http.Request) {
	s.serverSystem.StopQuit()

	gameID := mux.Vars(r)["game"]
	collectionID := mux.Vars(r)["collection"]
	deckID := mux.Vars(r)["deck"]
	sort := r.URL.Query().Get("sort")
	search := r.URL.Query().Get("search")
	items, meta, e := s.serviceCard.List(gameID, collectionID, deckID, sort, search)
	if e != nil {
		network.ResponseError(w, e)
		return
	}
	network.ResponseWithMeta(w, items, meta)
}
func (s *card) UpdateHandler(w http.ResponseWriter, r *http.Request) {
	gameID := mux.Vars(r)["game"]
	collectionID := mux.Vars(r)["collection"]
	deckID := mux.Vars(r)["deck"]
	cardID, e := fs.StringToInt64(mux.Vars(r)["card"])
	if e != nil {
		network.ResponseError(w, e)
		return
	}

	e = r.ParseMultipartForm(0)
	if e != nil {
		er.IfErrorLog(e)
		e = er.InternalError.AddMessage(e.Error())
		network.ResponseError(w, e)
		return
	}

	data, e := utils.GetFileFromMultipart("imageFile", r)
	if e != nil {
		network.ResponseError(w, e)
		return
	}

	dtoObject := &dto.UpdateCardDTO{
		Name:        r.FormValue("name"),
		Description: r.FormValue("description"),
		Image:       r.FormValue("image"),
		Variables:   nil,
		Count:       fs.StringToInt(r.FormValue("count")),
		ImageFile:   data,
	}

	variablesJson := r.FormValue("variables")
	if variablesJson != "" {
		e = json.Unmarshal([]byte(variablesJson), &dtoObject.Variables)
		if e != nil {
			er.IfErrorLog(e)
			e = er.InternalError.HTTP(http.StatusBadRequest).AddMessage("Bad variables json")
			network.ResponseError(w, e)
			return
		}
	}

	item, e := s.serviceCard.Update(gameID, collectionID, deckID, cardID, dtoObject)
	if e != nil {
		network.ResponseError(w, e)
		return
	}
	network.Response(w, item)
}
