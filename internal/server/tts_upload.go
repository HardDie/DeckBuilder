package server

import (
	"net/http"

	"github.com/HardDie/DeckBuilder/internal/network"
	"github.com/HardDie/DeckBuilder/internal/service"
)

type TTSServer struct {
	ttsService service.ITTSService
}

func NewTTSServer(ttsService service.ITTSService) *TTSServer {
	return &TTSServer{
		ttsService: ttsService,
	}
}

func (s *TTSServer) DataHandler(w http.ResponseWriter, r *http.Request) {
	resp, err := s.ttsService.DataForTTS()
	if err != nil {
		network.ResponseError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write(resp)
}
