package api

import (
	"net/http"

	"github.com/gorilla/mux"
)

type ITTSServer interface {
	DataHandler(w http.ResponseWriter, r *http.Request)
}

func RegisterTTSServer(route *mux.Router, srv ITTSServer) {
	route.HandleFunc("/api/generator/data", srv.DataHandler).Methods(http.MethodGet)
}

type UnimplementedTTSServer struct {
}

var (
	// Validation
	_ ITTSServer = &UnimplementedTTSServer{}
)

// swagger:parameters RequestDataTTS
type RequestDataTTS struct {
}

// swagger:response ResponseDataTTS
type ResponseDataTTS struct {
	// In: body
	Body []byte
}

// swagger:route GET /api/generator/data TTS RequestDataTTS
//
// # Get json file from last generator
//
// API for TTS for downloading JSON file inside game
//
//	Responses:
//	  200: ResponseDataTTS
//	  default: ResponseError
func (s *UnimplementedTTSServer) DataHandler(w http.ResponseWriter, r *http.Request) {}
