package replace

import (
	"net/http"

	er "github.com/HardDie/DeckBuilder/internal/errors"
	"github.com/HardDie/DeckBuilder/internal/network"
	servicesReplace "github.com/HardDie/DeckBuilder/internal/services/replace"
	"github.com/HardDie/DeckBuilder/internal/utils"
)

type replace struct {
	serviceReplace servicesReplace.Replace
}

func New(serviceReplace servicesReplace.Replace) Replace {
	return &replace{
		serviceReplace: serviceReplace,
	}
}

func (s *replace) PrepareHandler(w http.ResponseWriter, r *http.Request) {
	e := r.ParseMultipartForm(0)
	if e != nil {
		er.IfErrorLog(e)
		e = er.InternalError.HTTP(http.StatusBadRequest).AddMessage(e.Error())
		network.ResponseError(w, e)
		return
	}

	data, e := utils.GetFileFromMultipart("file", r)
	if e != nil {
		network.ResponseError(w, e)
		return
	}
	if data == nil {
		e = er.BadArchive.AddMessage("The file must be passed as an argument")
		network.ResponseError(w, e)
		return
	}

	resp, e := s.serviceReplace.Prepare(data)
	if e != nil {
		network.ResponseError(w, e)
		return
	}
	network.Response(w, resp)
}
func (s *replace) ReplaceHandler(w http.ResponseWriter, r *http.Request) {
	e := r.ParseMultipartForm(0)
	if e != nil {
		er.IfErrorLog(e)
		e = er.InternalError.HTTP(http.StatusBadRequest).AddMessage(e.Error())
		network.ResponseError(w, e)
		return
	}

	data, e := utils.GetFileFromMultipart("file", r)
	if e != nil {
		network.ResponseError(w, e)
		return
	}
	if data == nil {
		e = er.BadArchive.AddMessage("The file must be passed as an argument")
		network.ResponseError(w, e)
		return
	}

	mapping, e := utils.GetFileFromMultipart("mapping", r)
	if e != nil {
		network.ResponseError(w, e)
		return
	}
	if mapping == nil {
		e = er.BadArchive.AddMessage("The mapping must be passed as an argument")
		network.ResponseError(w, e)
		return
	}

	resp, e := s.serviceReplace.Replace(data, mapping)
	if e != nil {
		network.ResponseError(w, e)
		return
	}
	network.Response(w, resp)
}
