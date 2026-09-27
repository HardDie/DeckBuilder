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
4. **Copy handler mapping into `desktop/bindings/<aggregate>`** and call **services** from bindings. Bind each aggregate struct in `wails.Run`. Remove the REST route when that verb is unused.

## Decision

Use option 4.

Migrated desktop verbs (HTTP routes removed; Wails only):

* `desktop/bindings/game.Game`: `List` / `Create` / `Read` / `Update` / `Delete` / `Duplicate`
* `desktop/bindings/collection.Collection`: `List` / `Create` / `Read` / `Update` / `Delete`
* `desktop/bindings/deck.Deck`: `List` / `ListAllUnique` / `Create` / `Read` / `Update` / `Delete`
* `desktop/bindings/card.Card`: `List` / `Create` / `Read` / `Update` / `Delete`
* `desktop/bindings/system.System`: `GetSettings` / `UpdateSettings` / `Status` / `GetVersion`

`desktop/app.go` is window lifecycle only. Shared write fields and DTO/`cachedImage` mapping live in `desktop/bindings/catalog`. `application.Get` exposes config, catalog services, and the system service. `cachedImage` URLs still point at `/api/.../image`. Binding errors are toasted in the Vue API layer because they skip the `window.fetch` wrapper.

`GET /api/games/{game}/export`, `POST /api/games/import`, generate, images, and TTS stay on HTTP. `cmd/deck_builder` / `gui/` no longer have catalog CRUD REST.

Wails JSON turns a missing `imageFile` into `[]`; HTTP `GetFileFromMultipart` uses `nil`. Bindings must pass `WriteRequest.ImageBytes()` (empty → nil), not the raw slice.

This supersedes ADR 010’s rule that bindings are not the catalog API. ADR 010 still describes the window + loopback process. ADR 001 still describes the HTTP contract.

## Consequences

### Positive

* Desktop bindings do not grow the HTTP server interface that we expect to delete with the mux.
* Each catalog aggregate is its own bound struct under `desktop/bindings/`.

### Negative and risks

* Two copies of mapping exist while a verb still has both a binding and an HTTP handler (game zip import).
* Wails error toasts are a second path next to the fetch wrapper.

### Neutral

* Images, generate, and TTS `WebRequest.get` stay on HTTP.
* Image reads stay the `cachedImage` URL ([ADR 012](012-card-image-urls.md)).
* System settings, status, and version left HTTP.
