package api

import (
	"net/http"

	"github.com/gorilla/mux"

	serversReplace "github.com/HardDie/DeckBuilder/internal/servers/replace"
)

func RegisterReplaceServer(route *mux.Router, srv serversReplace.Replace) {
	ReplaceRoute := route.PathPrefix("/api/replace").Subrouter()
	ReplaceRoute.HandleFunc("/prepare", srv.PrepareHandler).Methods(http.MethodPost)
	ReplaceRoute.HandleFunc("", srv.ReplaceHandler).Methods(http.MethodPost)
}

type UnimplementedReplaceServer struct {
}

var (
	// Validation
	_ serversReplace.Replace = &UnimplementedReplaceServer{}
)

// swagger:parameters RequestPrepareReplace
type RequestPrepareReplace struct {
	// Json file
	// In: formData
	// Required: true
	File []byte `json:"file"`
}

// swagger:response ResponsePrepareReplace
type ResponsePrepareReplace struct {
	// In: body
	Body struct {
		// Required: true
		Data []byte `json:"data"`
	}
}

// swagger:route POST /api/replace/prepare Replace RequestPrepareReplace
//
// # Map with key image and empty value for URLs
//
// Takes as input the generated resulting json file for a saved TTS object.
// As a result, it returns a map of the files that should be uploaded to the web repository
// and allows you to manually map those files to valid URLs.
//
//	Consumes:
//	- multipart/form-data
//
//	Responses:
//	  200: ResponsePrepareReplace
//	  default: ResponseError
func (s *UnimplementedReplaceServer) PrepareHandler(w http.ResponseWriter, r *http.Request) {}

// swagger:parameters RequestReplaceReplace
type RequestReplaceReplace struct {
	// Json file
	// In: formData
	// Required: true
	File []byte `json:"file"`
	// Json file
	// In: formData
	// Required: true
	Mapping []byte `json:"mapping"`
}

// swagger:response ResponseReplaceReplace
type ResponseReplaceReplace struct {
	// In: body
	Body struct {
		// Required: true
		Data []byte `json:"data"`
	}
}

// swagger:route POST /api/replace Replace RequestReplaceReplace
//
// # Replace all image paths with a mapping file
//
// Accepts as input the generated resulting json file for the saved TTS object and the mapped URL files.
// This returns a saved TTS object json file with local file paths replaced by URLs.
// And with these cards you can already save your table and share it with other users.
//
//	Consumes:
//	- multipart/form-data
//
//	Responses:
//	  200: ResponseReplaceReplace
//	  default: ResponseError
func (s *UnimplementedReplaceServer) ReplaceHandler(w http.ResponseWriter, r *http.Request) {}
