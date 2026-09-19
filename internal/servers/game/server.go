package game

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/HardDie/DeckBuilder/internal/config"
	"github.com/HardDie/DeckBuilder/internal/dto"
	entitiesGame "github.com/HardDie/DeckBuilder/internal/entities/game"
	er "github.com/HardDie/DeckBuilder/internal/errors"
	"github.com/HardDie/DeckBuilder/internal/network"
	servicesGame "github.com/HardDie/DeckBuilder/internal/services/game"
	"github.com/HardDie/DeckBuilder/internal/utils"
)

type game struct {
	cfg         config.Config
	serviceGame servicesGame.Game
}

func New(cfg config.Config, serviceGame servicesGame.Game) Game {
	return &game{
		cfg:         cfg,
		serviceGame: serviceGame,
	}
}

func (s *game) ExportHandler(w http.ResponseWriter, r *http.Request) {
	gameID := mux.Vars(r)["game"]
	archive, e := s.serviceGame.Export(gameID)
	if e != nil {
		network.ResponseError(w, e)
		return
	}
	w.Header().Set("Content-Type", "application/zip")
	if _, err := w.Write(archive); err != nil {
		er.IfErrorLog(err)
	}
}

func (s *game) ImportHandler(w http.ResponseWriter, r *http.Request) {
	e := r.ParseMultipartForm(0)
	if e != nil {
		er.IfErrorLog(e)
		e = er.InternalError.HTTP(http.StatusBadRequest).AddMessage(e.Error())
		network.ResponseError(w, e)
		return
	}

	name := r.FormValue("name")

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

	item, e := s.serviceGame.Import(data, name)
	if e != nil {
		network.ResponseError(w, e)
		return
	}

	network.Response(w, dto.Game{
		ID:          item.ID,
		Name:        item.Name,
		Description: item.Description,
		Image:       item.Image,
		CachedImage: s.calculateCachedImage(*item),
		CreatedAt:   item.CreatedAt,
		UpdatedAt:   item.UpdatedAt,
	})
}

func (s *game) calculateCachedImage(game entitiesGame.Game) string {
	return fmt.Sprintf(s.cfg.GameImagePath+"?%s", game.ID, utils.HashForTime(&game.UpdatedAt))
}
