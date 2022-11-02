package server

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/HardDie/DeckBuilder/internal/dto"
	er "github.com/HardDie/DeckBuilder/internal/errors"
	"github.com/HardDie/DeckBuilder/internal/fs"
	"github.com/HardDie/DeckBuilder/internal/network"
	"github.com/HardDie/DeckBuilder/internal/service"
	"github.com/HardDie/DeckBuilder/internal/utils"
)

type CardServer struct {
	cardService  service.ICardService
	systemServer *SystemServer
}

func NewCardServer(cardService service.ICardService, systemServer *SystemServer) *CardServer {
	return &CardServer{
		cardService:  cardService,
		systemServer: systemServer,
	}
}

func (s *CardServer) CreateHandler(w http.ResponseWriter, r *http.Request) {
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

	item, e := s.cardService.Create(gameID, collectionID, deckID, dtoObject)
	if e != nil {
		network.ResponseError(w, e)
		return
	}
	network.Response(w, item)
}
func (s *CardServer) DeleteHandler(w http.ResponseWriter, r *http.Request) {
	gameID := mux.Vars(r)["game"]
	collectionID := mux.Vars(r)["collection"]
	deckID := mux.Vars(r)["deck"]
	cardID, e := fs.StringToInt64(mux.Vars(r)["card"])
	if e != nil {
		network.ResponseError(w, e)
		return
	}
	e = s.cardService.Delete(gameID, collectionID, deckID, cardID)
	if e != nil {
		network.ResponseError(w, e)
	}
}
func (s *CardServer) ItemHandler(w http.ResponseWriter, r *http.Request) {
	gameID := mux.Vars(r)["game"]
	collectionID := mux.Vars(r)["collection"]
	deckID := mux.Vars(r)["deck"]
	cardID, e := fs.StringToInt64(mux.Vars(r)["card"])
	if e != nil {
		network.ResponseError(w, e)
		return
	}
	item, e := s.cardService.Item(gameID, collectionID, deckID, cardID)
	if e != nil {
		network.ResponseError(w, e)
		return
	}
	network.Response(w, item)
}
func (s *CardServer) ListHandler(w http.ResponseWriter, r *http.Request) {
	s.systemServer.StopQuit()

	gameID := mux.Vars(r)["game"]
	collectionID := mux.Vars(r)["collection"]
	deckID := mux.Vars(r)["deck"]
	sort := r.URL.Query().Get("sort")
	search := r.URL.Query().Get("search")
	items, meta, e := s.cardService.List(gameID, collectionID, deckID, sort, search)
	if e != nil {
		network.ResponseError(w, e)
		return
	}
	network.ResponseWithMeta(w, items, meta)
}
func (s *CardServer) UpdateHandler(w http.ResponseWriter, r *http.Request) {
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

	item, e := s.cardService.Update(gameID, collectionID, deckID, cardID, dtoObject)
	if e != nil {
		network.ResponseError(w, e)
		return
	}
	network.Response(w, item)
}
