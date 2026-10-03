package servers

import (
	"log/slog"
	"net"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/HardDie/DeckBuilder/internal/config"
	"github.com/HardDie/DeckBuilder/internal/fs"
	"github.com/HardDie/DeckBuilder/internal/network"
)

type gameImage interface {
	GetImage(gameID string) ([]byte, string, error)
}

type collectionImage interface {
	GetImage(gameID, collectionID string) ([]byte, string, error)
}

type deckImage interface {
	GetImage(gameID, collectionID, deckID string) ([]byte, string, error)
}

type cardImage interface {
	GetImage(gameID, collectionID, deckID string, cardID int64) ([]byte, string, error)
}

type ttsHTTP interface {
	DataForTTS() ([]byte, error)
	SetHTTPPort(port int)
}

type handlers struct {
	game       gameImage
	collection collectionImage
	deck       deckImage
	card       cardImage
	tts        ttsHTTP
}

type Server struct {
	router *mux.Router
	tts    ttsHTTP
}

func New(
	game gameImage,
	collection collectionImage,
	deck deckImage,
	card cardImage,
	tts ttsHTTP,
) *Server {
	route := mux.NewRouter().StrictSlash(false)
	h := &handlers{
		game:       game,
		collection: collection,
		deck:       deck,
		card:       card,
		tts:        tts,
	}
	register(route, h)
	route.Use(corsMiddleware)
	return &Server{router: route, tts: tts}
}

func register(route *mux.Router, h *handlers) {
	gamesRoute := route.PathPrefix("/api/games").Subrouter()
	gamesRoute.HandleFunc("/{game}/image", h.gameHandler).Methods(http.MethodGet)

	collectionsRoute := gamesRoute.PathPrefix("/{game}/collections").Subrouter()
	collectionsRoute.HandleFunc("/{collection}/image", h.collectionHandler).Methods(http.MethodGet)

	decksRoute := collectionsRoute.PathPrefix("/{collection}/decks").Subrouter()
	decksRoute.HandleFunc("/{deck}/image", h.deckHandler).Methods(http.MethodGet)

	cardsRoute := decksRoute.PathPrefix("/{deck}/cards").Subrouter()
	cardsRoute.HandleFunc("/{card}/image", h.cardHandler).Methods(http.MethodGet)

	route.HandleFunc("/api/tts/data", h.ttsHandler).Methods(http.MethodGet)
}

func (s *Server) Handler() http.Handler {
	return s.router
}

func (s *Server) Listen() (net.Listener, int, error) {
	ln, port, err := listenLoopback(config.HTTPHost, config.HTTPPort, config.HTTPPortAttempts)
	if err != nil {
		return nil, 0, err
	}
	s.tts.SetHTTPPort(port)
	slog.Info("HTTP server listening", "host", config.HTTPHost, "port", port)
	return ln, port, nil
}

func (s *Server) Serve(ln net.Listener) error {
	return http.Serve(ln, s.router)
}

func corsSetupHeaders(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET,OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, ContentType")
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		corsSetupHeaders(w)
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func writeImage(w http.ResponseWriter, img []byte, imgType string, err error) {
	if err != nil {
		network.ResponseError(w, err)
		return
	}
	w.Header().Set("Content-Type", "image/"+imgType)
	if _, err := w.Write(img); err != nil {
		slog.Warn("write image response", "err", err)
	}
}

func (h *handlers) gameHandler(w http.ResponseWriter, r *http.Request) {
	img, imgType, err := h.game.GetImage(mux.Vars(r)["game"])
	writeImage(w, img, imgType, err)
}

func (h *handlers) collectionHandler(w http.ResponseWriter, r *http.Request) {
	img, imgType, err := h.collection.GetImage(mux.Vars(r)["game"], mux.Vars(r)["collection"])
	writeImage(w, img, imgType, err)
}

func (h *handlers) deckHandler(w http.ResponseWriter, r *http.Request) {
	img, imgType, err := h.deck.GetImage(mux.Vars(r)["game"], mux.Vars(r)["collection"], mux.Vars(r)["deck"])
	writeImage(w, img, imgType, err)
}

func (h *handlers) cardHandler(w http.ResponseWriter, r *http.Request) {
	cardID, err := fs.StringToInt64(mux.Vars(r)["card"])
	if err != nil {
		network.ResponseError(w, err)
		return
	}
	img, imgType, err := h.card.GetImage(mux.Vars(r)["game"], mux.Vars(r)["collection"], mux.Vars(r)["deck"], cardID)
	writeImage(w, img, imgType, err)
}

func (h *handlers) ttsHandler(w http.ResponseWriter, r *http.Request) {
	resp, err := h.tts.DataForTTS()
	if err != nil {
		network.ResponseError(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write(resp)
}
