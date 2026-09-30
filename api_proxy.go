package main

import (
	"net/http"
	"strings"
)

// Image GETs use /api/.../image.
// Wails does not put those files in the embed.
// Handler serves them from the same mux as the loopback server.
// In wails dev, Vite still proxies GET /api to that server.
func apiHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasPrefix(r.URL.Path, "/api") {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			return
		}
		h.ServeHTTP(w, r)
	})
}
