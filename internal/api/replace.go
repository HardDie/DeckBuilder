package api

import (
	"net/http"

	"github.com/gorilla/mux"
)

type IReplaceServer interface {
	PrepareHandler(w http.ResponseWriter, r *http.Request)
	ReplaceHandler(w http.ResponseWriter, r *http.Request)
}

func RegisterReplaceServer(route *mux.Router, srv IReplaceServer) {
	ReplaceRoute := route.PathPrefix("/api/replace").Subrouter()
	ReplaceRoute.HandleFunc("/prepare", srv.PrepareHandler).Methods(http.MethodPost)
	ReplaceRoute.HandleFunc("", srv.ReplaceHandler).Methods(http.MethodPost)
}

type UnimplementedReplaceServer struct {
}

var (
	// Validation
	_ IReplaceServer = &UnimplementedReplaceServer{}
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
// # List of images
//
//	Consumes:
//	- multipart/form-data
//
//	Produces:
//	- application/json
//
//	Schemes: http
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
// # List of images
//
//	Consumes:
//	- multipart/form-data
//
//	Produces:
//	- application/json
//
//	Schemes: http
//
//	Responses:
//	  200: ResponseReplaceReplace
//	  default: ResponseError
func (s *UnimplementedReplaceServer) ReplaceHandler(w http.ResponseWriter, r *http.Request) {}
