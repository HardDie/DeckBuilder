# 1. Loopback HTTP process; GUI is a separate SPA

* **Status:** Superseded by [013](013-http-images-and-tts.md)
* **Date:** 2026-09-18
* **Authors:** @oleg

---

## Context

DeckBuilder must feel like a desktop app: open a window, edit a catalog, export for Tabletop Simulator. The GUI is Vue (DeckBuilderGUI). The backend is Go.

Options include a single Wails/webview binary, a public web service, or a local HTTP server that serves the built SPA and JSON API from one process.

## Considered options

1. **Wails (or similar) bindings** — Go methods bound into a webview; no HTTP contract for the GUI.
2. **Cloud-hosted API** — accounts, auth, remote storage.
3. **Local `net/http` on loopback** — one binary listens on `127.0.0.1:5000`, embeds `web/dist`, CORS open for GUI `yarn` dev.

## Decision

Use option 3.

The process starts `internal/application`, which binds **only** `127.0.0.1:5000`. Release builds embed the GUI.

The GUI repository stays separate; this repo vendors built assets via `make web-build`, not source.

## Consequences

### Positive

* The GUI can be developed against a stable REST + Swagger contract (`/docs`).
* TTS Lua can `WebRequest.get("http://127.0.0.1:5000/api/tts/data")` without a second protocol.
* No auth or TLS for a single-user loopback tool.

### Negative and risks

* Two repositories must stay in sync on DTO fields and routes.
* Binding loopback means other machines cannot use the API (intentional).

### Neutral

* SPA deep links are registered as mux forwarders to `index.html` (`internal/api/static.go`).
