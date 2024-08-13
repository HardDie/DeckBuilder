package tts

import (
	"net/http"

	"github.com/HardDie/DeckBuilder/internal/network"
	servicesTTS "github.com/HardDie/DeckBuilder/internal/services/tts"
)

type tts struct {
	serviceTTS servicesTTS.TTS
}

func New(serviceTTS servicesTTS.TTS) TTS {
	return &tts{
		serviceTTS: serviceTTS,
	}
}

func (s *tts) DataHandler(w http.ResponseWriter, r *http.Request) {
	resp, err := s.serviceTTS.DataForTTS()
	if err != nil {
		network.ResponseError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write(resp)
}
