# 11. Incremental Wails catalog bindings beside HTTP

* **Status:** Accepted
* **Date:** 2026-09-19
* **Authors:** @oleg

---

## Context

[ADR 010](010-wails-shell.md) put a Wails window around the same loopback HTTP stack and said catalog calls stay `fetch('/api/...')`. That avoided a one-shot DTO rewrite. The desktop UI can now take verbs one at a time without dropping HTTP for `cmd/deck_builder`, images, TTS Lua, or Swagger.

The HTTP `internal/servers` layer exists to adapt mux/swagger. When catalog traffic is all Wails, that adapter is not the long-term home for list/read/update.

## Considered options

1. **Keep every catalog call on HTTP** in the Wails window forever.
2. **Replace all `/api` routes with bindings in one pass.**
3. **Share a non-HTTP method on `internal/servers`** used by both `ListHandler` and Wails.
4. **Copy handler mapping into `desktop/app.go`** and call **services** (and `StopQuit`) from bindings. Remove the REST route when that verb is unused.

## Decision

Use option 4.

Migrated desktop verbs (HTTP routes removed; Wails only):

* Games: `ListGames` / `CreateGame` / `ReadGame` / `UpdateGame` / `DeleteGame` / `DuplicateGame`
* Collections: `ListCollections` / `CreateCollection` / `ReadCollection` / `UpdateCollection` / `DeleteCollection`

Each duplicates DTO/`cachedImage` mapping in `desktop/app.go` instead of extending the HTTP server interface. `application.Get` exposes config, game/collection services, and `servers/system`. `cachedImage` URLs still point at `/api/.../image`. Binding errors are toasted in the Vue API layer because they skip the `window.fetch` wrapper.

`GET /api/games/{game}/export`, `POST /api/games/import`, generate, images, TTS, decks, and cards stay on HTTP. `cmd/deck_builder` / `gui/` no longer have game or collection CRUD REST.

This supersedes ADR 010’s rule that bindings are not the catalog API. ADR 010 still describes the window + loopback process. ADR 001 still describes the HTTP contract.

## Consequences

### Positive

* Desktop bindings do not grow the HTTP server interface that we expect to delete with the mux.
* Bindings sit next to services, which is the layer that remains after HTTP adapters go away.

### Negative and risks

* Two copies of mapping exist while a verb still has both a binding and an HTTP handler (game zip import).
* Wails error toasts are a second path next to the fetch wrapper.

### Neutral

* Images, generate, TTS `WebRequest.get`, and not-yet-migrated CRUD stay on HTTP.
