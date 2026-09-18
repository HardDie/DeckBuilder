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
3. **Scaffold Wails in `desktop/`** as its own module; run the same `internal/application` HTTP stack in that process. Bindings are not the catalog API. `cmd/deck_builder` remains the browser/embed entry.

## Decision

Use option 3.

`desktop/` is a Wails v2.16 project (`wails.json`, `app.go`). The Vue UI is a copy of `gui/` under `desktop/frontend`. `fetch('/api/...')` is unchanged. The Wails process calls `application.Get` / `Listen` / `Serve` (loopback, port retry, no browser). TTS still uses that port. In `wails dev`, GET `/api` can go through Vite to `:5000`; Wails does **not** send POST/PATCH/DELETE to Vite, so `AssetServer.Handler` serves `/api` from the same mux. Bindings are not the catalog API. Production HTTP+SPA remains `cmd/deck_builder`. Wails CLI must be v2.16+ (Go 1.27 breaks binding generation on v2.11).

A nested module (`github.com/HardDie/DeckBuilder/desktop`) keeps CGO/webview out of root `go test ./...`.

## Consequences

### Positive

* One `wails dev` / Wails binary is enough for window + API (same data dir as `cmd/deck_builder`).
* First UI migration is “window + existing HTTP,” not a DTO rewrite.

### Negative and risks

* Do not run `cmd/deck_builder` and Wails at the same time if both want `:5000` (Wails will take the next free port; Vite GET still targets `:5000`).
* Do not replace `/api` with Wails bindings in the same step as the window shell.

### Neutral

* ADR 001 stays in force for the HTTP contract. This ADR does not supersede it.
