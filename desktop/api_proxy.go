package main

import (
	"net/http"
	"strings"
)

// Wails AssetServer only reverse-proxies GET to the Vite dev server; POST/PATCH/DELETE
// (and GET when the file is missing from assets) go through Handler. Serve /api from
// the same mux as the loopback HTTP server.
func apiHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasPrefix(r.URL.Path, "/api") {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			return
		}
		h.ServeHTTP(w, r)
	})
}
