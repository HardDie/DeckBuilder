package tts

import "net/http"

type TTS interface {
	DataHandler(w http.ResponseWriter, r *http.Request)
}
