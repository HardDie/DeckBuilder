package server

import (
	"net/http"

	"github.com/HardDie/DeckBuilder/internal/network"
	servicesTTS "github.com/HardDie/DeckBuilder/internal/services/tts"
)

type TTSServer struct {
	serviceTTS servicesTTS.TTS
}

func NewTTSServer(serviceTTS servicesTTS.TTS) *TTSServer {
	return &TTSServer{
		serviceTTS: serviceTTS,
	}
}

func (s *TTSServer) DataHandler(w http.ResponseWriter, r *http.Request) {
	resp, err := s.serviceTTS.DataForTTS()
	if err != nil {
		network.ResponseError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write(resp)
}
