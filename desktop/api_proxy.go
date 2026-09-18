package main

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

// Loopback API used by cmd/deck_builder. Wails AssetServer only reverse-proxies
// GET to the Vite dev server; POST/PATCH/DELETE need this Handler (see
// wails v2 ExternalAssetsHandler).
const loopbackAPI = "http://127.0.0.1:5000"

func apiProxy() http.Handler {
	target, err := url.Parse(loopbackAPI)
	if err != nil {
		panic(err)
	}
	proxy := httputil.NewSingleHostReverseProxy(target)
	baseDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		baseDirector(req)
		req.Host = target.Host
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasPrefix(r.URL.Path, "/api") {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			return
		}
		proxy.ServeHTTP(w, r)
	})
}
