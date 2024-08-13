package deck

import "net/http"

type Deck interface {
	AllDecksHandler(w http.ResponseWriter, r *http.Request)
	CreateHandler(w http.ResponseWriter, r *http.Request)
	DeleteHandler(w http.ResponseWriter, r *http.Request)
	ItemHandler(w http.ResponseWriter, r *http.Request)
	ListHandler(w http.ResponseWriter, r *http.Request)
	UpdateHandler(w http.ResponseWriter, r *http.Request)
}
