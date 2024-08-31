package api

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/HardDie/DeckBuilder/internal/dto"
	"github.com/HardDie/DeckBuilder/internal/network"
	serversCollection "github.com/HardDie/DeckBuilder/internal/servers/collection"
)

func RegisterCollectionServer(route *mux.Router, srv serversCollection.Collection) {
	CollectionsRoute := route.PathPrefix("/api/games/{game}/collections").Subrouter()
	CollectionsRoute.HandleFunc("", srv.ListHandler).Methods(http.MethodGet)
	CollectionsRoute.HandleFunc("", srv.CreateHandler).Methods(http.MethodPost)
	CollectionsRoute.HandleFunc("/{collection}", srv.DeleteHandler).Methods(http.MethodDelete)
	CollectionsRoute.HandleFunc("/{collection}", srv.ItemHandler).Methods(http.MethodGet)
	CollectionsRoute.HandleFunc("/{collection}", srv.UpdateHandler).Methods(http.MethodPatch)
}

type UnimplementedCollectionServer struct {
}

var (
	// Validation
	_ serversCollection.Collection = &UnimplementedCollectionServer{}
)

// Request to create a collection
//
// swagger:parameters RequestCreateCollection
type RequestCreateCollection struct {
	// In: path
	// Required: true
	Game string `json:"game"`
	// In: formData
	// Required: true
	Name string `json:"name"`
	// In: formData
	// Required: true
	Description string `json:"description"`
	// In: formData
	// Required: false
	Image string `json:"image"`
	// In: formData
	// Required: false
	ImageFile []byte `json:"imageFile"`
}

// Status of collection creation
//
// swagger:response ResponseCreateCollection
type ResponseCreateCollection struct {
	// In: body
	// Required: true
	Body struct {
		// Required: true
		Data dto.Collection `json:"data"`
	}
}

// swagger:route POST /api/games/{game}/collections Collections RequestCreateCollection
//
// # Create collection
//
// Allows you to create a new collection
//
//	Consumes:
//	- multipart/form-data
//
//	Responses:
//	  200: ResponseCreateCollection
//	  default: ResponseError
func (s *UnimplementedCollectionServer) CreateHandler(w http.ResponseWriter, r *http.Request) {}

// Request to delete a collection
//
// swagger:parameters RequestDeleteCollection
type RequestDeleteCollection struct {
	// In: path
	// Required: true
	Game string `json:"game"`
	// In: path
	// Required: true
	Collection string `json:"collection"`
}

// Collection deletion status
//
// swagger:response ResponseDeleteCollection
type ResponseDeleteCollection struct {
}

// swagger:route DELETE /api/games/{game}/collections/{collection} Collections RequestDeleteCollection
//
// # Delete collection
//
// Allows you to delete an existing collection
//
//	Responses:
//	  200: ResponseDeleteCollection
//	  default: ResponseError
func (s *UnimplementedCollectionServer) DeleteHandler(w http.ResponseWriter, r *http.Request) {}

// Requesting an existing collection
//
// swagger:parameters RequestCollection
type RequestCollection struct {
	// In: path
	// Required: true
	Game string `json:"game"`
	// In: path
	// Required: true
	Collection string `json:"collection"`
}

// Collection
//
// swagger:response ResponseCollection
type ResponseCollection struct {
	// In: body
	Body struct {
		// Required: true
		Data dto.Collection `json:"data"`
	}
}

// swagger:route GET /api/games/{game}/collections/{collection} Collections RequestCollection
//
// # Get collection
//
// Get an existing collection
//
//	Responses:
//	  200: ResponseCollection
//	  default: ResponseError
func (s *UnimplementedCollectionServer) ItemHandler(w http.ResponseWriter, r *http.Request) {}

// Requesting a list of existing collections
//
// swagger:parameters RequestListOfCollections
type RequestListOfCollections struct {
	// In: path
	// Required: true
	Game string `json:"game"`
	// In: query
	// Required: false
	Sort string `json:"sort"`
	// In: query
	// Required: false
	Search string `json:"search"`
}

// List of collections
//
// swagger:response ResponseListOfCollections
type ResponseListOfCollections struct {
	// In: body
	// Required: true
	Body struct {
		// Required: true
		Data []*dto.Collection `json:"data"`
		// Required: true
		Meta *network.Meta `json:"meta"`
	}
}

// swagger:route GET /api/games/{game}/collections Collections RequestListOfCollections
//
// # Get collections list
//
// Get a list of existing collections
// Sort values: name, name_desc, created, created_desc
//
//	Responses:
//	  200: ResponseListOfCollections
//	  default: ResponseError
func (s *UnimplementedCollectionServer) ListHandler(w http.ResponseWriter, r *http.Request) {}

// Request to update a collection
//
// swagger:parameters RequestUpdateCollection
type RequestUpdateCollection struct {
	// In: path
	// Required: true
	Game string `json:"game"`
	// In: path
	// Required: true
	Collection string `json:"collection"`
	// In: formData
	// Required: true
	Name string `json:"name"`
	// In: formData
	// Required: true
	Description string `json:"description"`
	// In: formData
	// Required: false
	Image string `json:"image"`
	// In: formData
	// Required: false
	ImageFile []byte `json:"imageFile"`
}

// Status of collection update
//
// swagger:response ResponseUpdateCollection
type ResponseUpdateCollection struct {
	// In: body
	// Required: true
	Body struct {
		// Required: true
		Data dto.Collection `json:"data"`
	}
}

// swagger:route PATCH /api/games/{game}/collections/{collection} Collections RequestUpdateCollection
//
// # Update collection
//
// Allows you to update an existing collection
//
//	Consumes:
//	- multipart/form-data
//
//	Responses:
//	  200: ResponseUpdateCollection
//	  default: ResponseError
func (s *UnimplementedCollectionServer) UpdateHandler(w http.ResponseWriter, r *http.Request) {}
