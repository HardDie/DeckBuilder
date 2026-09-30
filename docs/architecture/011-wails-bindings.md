# 11. Incremental Wails catalog bindings beside HTTP

* **Status:** Accepted
* **Also:** [ADR 013](013-http-images-and-tts.md) moves replace off HTTP and drops the SPA.
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
4. **Copy handler mapping into `bindings/<aggregate>`** and call **services** from bindings. Bind each aggregate struct in `wails.Run`. Remove the REST route when that verb is unused.

## Decision

Use option 4.

Migrated desktop verbs (HTTP routes removed; Wails only):

* `bindings/game.Game`: `List` / `Create` / `Read` / `Update` / `Delete` / `Duplicate` / `Export` / `Import`
* `bindings/collection.Collection`: `List` / `Create` / `Read` / `Update` / `Delete`
* `bindings/deck.Deck`: `List` / `ListAllUnique` / `Create` / `Read` / `Update` / `Delete`
* `bindings/card.Card`: `List` / `Create` / `Read` / `Update` / `Delete`
* `bindings/system.System`: `GetSettings` / `UpdateSettings` / `Status` / `GetVersion`
* `bindings/search.Search`: `Root` / `Game` / `Collection`
* `bindings/generator.Generator`: `Game`
* `bindings/replace.Replace`: `Prepare` / `Replace`

`app.go` is window lifecycle only. Shared write fields and DTO/`cachedImage` mapping live in `bindings/catalog`. `application.Get` exposes config, catalog services, the system service, the search service, and the generator service. `cachedImage` URLs still point at `/api/.../image`. Binding errors are toasted in the Vue API layer because they skip the `window.fetch` wrapper.

Images and TTS stay on HTTP. Catalog CRUD is Wails only.

* `game.Export` opens a native save dialog.
* The binding writes the zip to the chosen path.
* Cancel writes nothing.
* `game.Import` takes zip bytes and an optional name.
* It returns the created game.
* `generator.Game` starts generate and returns.
* Drawing continues in a goroutine.
* Scale below 1 becomes 1.

Wails JSON turns a missing `imageFile` into `[]`; HTTP `GetFileFromMultipart` uses `nil`. Bindings must pass `WriteRequest.ImageBytes()` (empty → nil), not the raw slice.

This supersedes ADR 010’s rule that bindings are not the catalog API. ADR 010 still describes the window + loopback process. [ADR 013](013-http-images-and-tts.md) describes the HTTP contract.

## Consequences

### Positive

* Desktop bindings do not grow the HTTP server interface that we expect to delete with the mux.
* Each catalog aggregate is its own bound struct under `bindings/`.

### Negative and risks

* Wails error toasts are a second path next to the fetch wrapper.

### Neutral

* Images and TTS `WebRequest.get` stay on HTTP.
* Image reads stay the `cachedImage` URL ([ADR 012](012-card-image-urls.md)).
* System settings, status, and version left HTTP.
* Recursive search left HTTP.
