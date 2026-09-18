# 10. Wails desktop shell beside the loopback HTTP server

* **Status:** Accepted
* **Date:** 2026-09-18
* **Authors:** @oleg

---

## Context

The product should feel like a desktop app. Today [ADR 001](001-loopback-http-embedded-spa.md) uses a loopback HTTP server and opens a browser. A later step may wrap the same Vue UI in a Wails window. Rewriting REST into Wails Go bindings in one change is high risk: mux, swagger, TTS `WebRequest.get`, and the GUI repo all assume HTTP.

## Considered options

1. **Replace HTTP with Wails bindings** in one pass.
2. **No Wails** — keep browser + embed forever.
3. **Scaffold Wails in `desktop/`** as its own module; keep `cmd/deck_builder` as the HTTP process. Later, the window loads `http://127.0.0.1:<port>/` (same API). Bindings are not the catalog API.

## Decision

Use option 3.

`desktop/` is a Wails v2.16 project (`wails.json`, `app.go`). The Vue UI is a copy of `gui/` under `desktop/frontend`. `fetch('/api/...')` is unchanged. In `wails dev`, GET `/api` is proxied by Vite to `127.0.0.1:5000`; Wails does **not** send POST/PATCH/DELETE to Vite, so `AssetServer.Handler` reverse-proxies those to the same loopback server. Bindings are not the catalog API. Production entry remains `cmd/deck_builder`. Wails CLI must be v2.16+ (Go 1.27 breaks binding generation on v2.11).

A nested module (`github.com/HardDie/DeckBuilder/desktop`) keeps CGO/webview out of root `go test ./...`.

## Consequences

### Positive

* Wails can be developed without blocking HTTP/swagger/TTS.
* First UI migration can be “window + existing port,” not a DTO rewrite.

### Negative and risks

* Two binaries until the shell is wired. Do not invent a second HTTP stack inside Wails.
* Do not replace `/api` with Wails bindings in the same step as the window shell.

### Neutral

* ADR 001 stays in force for the HTTP contract. This ADR does not supersede it.
