# 10. Wails desktop shell beside the loopback HTTP server

* **Status:** Accepted (the “bindings are not the catalog API” rule is superseded by [ADR 011](011-wails-bindings.md))
* **Also:** [ADR 013](013-http-images-and-tts.md) drops the embedded SPA and the browser.
* **Also:** The UI is `frontend/`. There is no `gui/` submodule.
* **Also:** [ADR 014](014-wails-at-module-root.md) moves Wails to the module root.
* **Date:** 2026-09-18
* **Authors:** @oleg

---

## Context

The product should feel like a desktop app. Today [ADR 001](001-loopback-http-embedded-spa.md) uses a loopback HTTP server and opens a browser. A later step may wrap the same Vue UI in a Wails window. Rewriting REST into Wails Go bindings in one change is high risk: mux, swagger, TTS `WebRequest.get`, and the GUI repo all assume HTTP.

## Considered options

1. **Replace HTTP with Wails bindings** in one pass.
2. **No Wails** — keep browser + embed forever.
3. **Scaffold Wails in `desktop/`** as its own module; run the same loopback HTTP stack in that process. Bindings are not the catalog API. The root binary remains the browser/embed entry.

## Decision

Use option 3.

`desktop/` is a Wails v2.16 project (`wails.json`, `app.go`). The Vue UI is a copy of `gui/` under `desktop/frontend`. `fetch('/api/...')` is unchanged. The Wails process calls `wire` / `Listen` / `Serve` (loopback, port retry, no browser). TTS still uses that port. In `wails dev`, GET `/api` can go through Vite to `:5000`; Wails does **not** send POST/PATCH/DELETE to Vite, so `AssetServer.Handler` serves `/api` from the same mux. Bindings are not the catalog API. Production HTTP+SPA remains the root binary. Wails CLI must be v2.16+ (Go 1.27 breaks binding generation on v2.11).

A nested module (`github.com/HardDie/DeckBuilder/desktop`) keeps CGO/webview out of root `go test ./...`.

## Consequences

### Positive

* One `wails dev` / Wails binary is enough for window + API.
* First UI migration is “window + existing HTTP,” not a DTO rewrite.

### Negative and risks

* Do not replace `/api` with Wails bindings in the same step as the window shell.

### Neutral

* [ADR 013](013-http-images-and-tts.md) is the HTTP contract for the window.
