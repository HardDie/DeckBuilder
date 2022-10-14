package api

import (
	"net/http"
	"path/filepath"
	"strings"

	"github.com/go-openapi/runtime/middleware"
	"github.com/gorilla/mux"

	"github.com/HardDie/DeckBuilder/internal/logger"
	"github.com/HardDie/DeckBuilder/web"
)

var (
	// In this map, we store the matches between the http path and the file path in the built-in fs structure
	pages = map[string]string{}
)

// This handler is assigned to each file path we have in the page map.
// When we get a call, we return one of the files.
func servePages(w http.ResponseWriter, r *http.Request) {
	page, ok := pages[r.URL.Path]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	file, err := web.GetWeb().ReadFile(page)
	if err != nil {
		logger.Error.Printf("page %s not found in pages cache...", r.RequestURI)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	switch filepath.Ext(page) {
	case ".css":
		w.Header().Set("Content-Type", "text/css")
	case ".js":
		w.Header().Set("Content-Type", "text/javascript")
	case ".ico":
		w.Header().Set("Content-Type", "image/x-icon")
	}
	_, err = w.Write(file)
	if err != nil {
		logger.Error.Println(err.Error())
	}
}

// Go through each file in the embedded fs, convert the file path to a url path, and put the value in map.
func registerFiles(dirName string) {
	files, _ := web.GetWeb().ReadDir(dirName)
	for _, file := range files {
		if file.IsDir() {
			registerFiles(dirName + "/" + file.Name())
		} else {
			fullName := dirName + "/" + file.Name()
			src := strings.TrimPrefix(fullName, "dist")
			if src == "/index.html" {
				pages["/"] = fullName
				continue
			}
			pages[src] = fullName
		}
	}
}

// Workaround. Allow any path to be served as /.
// This will solve the problem when we reload the page at some url and such file is not registered on go mux.
func forwarder(w http.ResponseWriter, r *http.Request) {
	r.URL.Path = "/"
	servePages(w, r)
}

func RegisterStaticServer(route *mux.Router) {
	// Parse files and fill pages map
	registerFiles("dist")

	// Register each file as route
	for page := range pages {
		route.HandleFunc(page, servePages)
	}

	// Workaround: if page we reloaded, forward request to index.html
	route.HandleFunc("/games/{id}", forwarder)
	route.HandleFunc("/games/{id}/collections/{id}", forwarder)
	route.HandleFunc("/games/{id}/collections/{id}/decks/{id}", forwarder)

	// Swagger
	route.HandleFunc("/swagger.json", web.ServeSwagger)

	// Register for /docs page
	redocHandler := middleware.Redoc(middleware.RedocOpts{SpecURL: "/swagger.json"}, nil)
	route.Handle("/docs", redocHandler)
}
